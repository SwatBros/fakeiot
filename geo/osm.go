package geo

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"runtime"
	"slices"
	"sort"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmpbf"
)

type TagFilter struct {
	Tag   string
	Value string
}

// Create a FeatureCollection with Multipolygons representing the area of filtered cities (admin level 8)
//
// filepath File path to an osm.pbf file
//
// names The tag and value list used to filter cities
func ListAndFilterCities(filepath string, names []TagFilter) {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var wayIds = []osm.WayID{}

	{
		scanner := osmpbf.New(context.Background(), file, runtime.GOMAXPROCS(-1))
		defer scanner.Close()

		scanner.SkipNodes = true
		scanner.SkipWays = true

		scanner.FilterRelation = func(r *osm.Relation) bool {
			return r.Tags.Find("admin_level") == "8" && slices.ContainsFunc(names, func(filter TagFilter) bool {
				return r.Tags.Find(filter.Tag) == filter.Value
			})
		}

		for scanner.Scan() {
			switch o := scanner.Object().(type) {
			case *osm.Relation:
				for _, member := range o.Members {
					if member.Type == osm.TypeWay {
						wayIds = append(wayIds, osm.WayID(member.Ref))
					}
				}
			}
		}
	}

	file.Seek(0, io.SeekStart)

	var wayResults = []WayResult{}

	{
		scanner := osmpbf.New(context.Background(), file, runtime.GOMAXPROCS(-1))
		defer scanner.Close()

		scanner.SkipNodes = true
		scanner.SkipRelations = true

		scanner.FilterWay = func(r *osm.Way) bool {
			return slices.Contains(wayIds, r.ID)
		}

		for scanner.Scan() {
			switch o := scanner.Object().(type) {
			case *osm.Node:
				break
			case *osm.Way:
				var nodeIds []osm.NodeID
				for _, node := range o.Nodes {
					nodeIds = append(nodeIds, node.ID)
				}

				wayResults = append(wayResults, WayResult{
					WayId:   o.ID,
					NodeIds: nodeIds,
				})
			case *osm.Relation:
				break
			}
		}
	}

	// Sort ways based on Relation member order
	sort.Slice(wayResults, func(i, j int) bool {
		a := wayResults[i]
		b := wayResults[j]
		return slices.Index(wayIds, a.WayId) < slices.Index(wayIds, b.WayId)
	})

	ringGroups := []*RingResult{
		{WayResult: []WayResult{wayResults[0]}},
	}

	// Group ways into rings
	for _, wayResult := range wayResults[1:] {
		found := false

		for _, ringGroup := range ringGroups {
			lastRingResult := ringGroup.WayResult[len(ringGroup.WayResult)-1]
			firstRingNode := lastRingResult.NodeIds[0]
			lastRingNode := lastRingResult.NodeIds[len(lastRingResult.NodeIds)-1]

			firstResultNode := wayResult.NodeIds[0]
			lastResultNode := wayResult.NodeIds[len(wayResult.NodeIds)-1]

			if firstRingNode == firstResultNode ||
				firstRingNode == lastResultNode ||
				lastRingNode == firstResultNode ||
				lastRingNode == lastResultNode {
				ringGroup.WayResult = append(ringGroup.WayResult, wayResult)
				found = true
				break
			}
		}

		if !found {
			ringGroups = append(ringGroups, &RingResult{WayResult: []WayResult{wayResult}})
		}
	}

	// Flip ways within each ring
	for _, ringGroup := range ringGroups {
		ringGroup.flipWays()
	}

	file.Seek(0, io.SeekStart)

	{
		scanner := osmpbf.New(context.Background(), file, runtime.GOMAXPROCS(-1))
		defer scanner.Close()
		scanner.SkipWays = true
		scanner.SkipRelations = true

		scanner.FilterNode = func(r *osm.Node) bool {
			for _, result := range wayResults {
				if slices.Contains(result.NodeIds, r.ID) {
					return true
				}
			}
			return false
		}

		for scanner.Scan() {
			switch o := scanner.Object().(type) {
			case *osm.Node:
				for _, ringGroup := range ringGroups {
					for _, result := range ringGroup.WayResult {
						if slices.Contains(result.NodeIds, o.ID) {
							ringGroup.Points = append(ringGroup.Points, PointResult{
								Point:  orb.Point{o.Lon, o.Lat},
								NodeId: o.ID,
								WayId:  result.WayId,
							})
						}
					}
				}
			case *osm.Way:
				break
			case *osm.Relation:
				break
			}
		}
	}

	collection := geojson.NewFeatureCollection()

	poly := orb.MultiPolygon{}
	for _, ringGroup := range ringGroups {
		poly = append(poly, orb.Polygon{ringGroup.OrmRing()})
	}

	feature := geojson.NewFeature(poly)

	collection.Features = append(collection.Features, feature)

	data, err := json.Marshal(collection)
	if err != nil {
		panic(err)
	}

	outfile := "./out.json"
	file, err = os.OpenFile(outfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}

	_, err = file.Write(data)
	if err != nil {
		panic(err)
	}
}

// A Point along a ring
//
// # Point Coordinates
//
// # NodeId The OSM NodeId
//
// WayId The OSM way the node belongs to
type PointResult struct {
	Point  orb.Point
	NodeId osm.NodeID
	WayId  osm.WayID
}

// An OSM way and its nodes
type WayResult struct {
	WayId   osm.WayID
	NodeIds []osm.NodeID
}

// Informations about a Ring
//
// # WayResult List of ways that compose the ring
//
// Points List of points that compose the ring
type RingResult struct {
	WayResult []WayResult
	Points    []PointResult
}

// Flip the ways so that they all point the same direction
//
// This is needed so that a polygon perimeter can be navigated
// A - B -> B - C
func (r *RingResult) flipWays() {
	if len(r.WayResult) < 2 {
		return
	}

	prev := &r.WayResult[0]
	next := &r.WayResult[1]

	lastNodeOfPrev := prev.NodeIds[len(prev.NodeIds)-1]
	firstNodeOfNext := next.NodeIds[0]
	lastNodeOfNext := next.NodeIds[len(next.NodeIds)-1]

	if lastNodeOfPrev != firstNodeOfNext && lastNodeOfNext != lastNodeOfPrev {
		slices.Reverse(prev.NodeIds)
	}

	for _, result := range r.WayResult[1:] {
		if prev.NodeIds[len(prev.NodeIds)-1] != result.NodeIds[0] {
			slices.Reverse(result.NodeIds)
		}
		prev = &result
	}
}

// Sort the points along a way
func (r *RingResult) sortPoints() {
	wayIds := make([]osm.WayID, len(r.WayResult))
	for i, result := range r.WayResult {
		wayIds[i] = result.WayId
	}

	sort.Slice(r.Points, func(i, j int) bool {
		if r.Points[i].WayId != r.Points[j].WayId {
			return slices.Index(wayIds, r.Points[i].WayId) < slices.Index(wayIds, r.Points[j].WayId)
		}

		var wayResult *WayResult
		for _, result := range r.WayResult {
			if result.WayId == r.Points[i].WayId {
				wayResult = &result
				break
			}
		}

		return slices.Index(wayResult.NodeIds, r.Points[i].NodeId) < slices.Index(wayResult.NodeIds, r.Points[j].NodeId)
	})
}

// Build the ORM ring sorting the points
func (r *RingResult) OrmRing() orb.Ring {
	r.sortPoints()

	points := make([]orb.Point, len(r.Points))
	for i, point := range r.Points {
		points[i] = point.Point
	}

	return orb.Ring(points)
}
