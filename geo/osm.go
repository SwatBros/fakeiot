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

func ListAndFilterCities(filepath string, names []string) {
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
			return r.Tags.Find("admin_level") == "8" && slices.Contains(names, r.Tags.Find("name"))
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

	sort.Slice(wayResults, func(i, j int) bool {
		a := wayResults[i]
		b := wayResults[j]
		return slices.Index(wayIds, a.WayId) < slices.Index(wayIds, b.WayId)
	})

	prev := &wayResults[0]
	next := &wayResults[1]

	lastNodeOfPrev := prev.NodeIds[len(prev.NodeIds)-1]
	firstNodeOfNext := next.NodeIds[0]
	lastNodeOfNext := next.NodeIds[len(next.NodeIds)-1]

	if lastNodeOfPrev != firstNodeOfNext && lastNodeOfNext != lastNodeOfPrev {
		slices.Reverse(prev.NodeIds)
	}

	for _, result := range wayResults[1:] {
		if prev.NodeIds[len(prev.NodeIds)-1] != result.NodeIds[0] {
			slices.Reverse(result.NodeIds)
		}
		prev = &result
	}

	file.Seek(0, io.SeekStart)

	var points = []PointResult{}

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
				for _, result := range wayResults {
					if slices.Contains(result.NodeIds, o.ID) {
						points = append(points, PointResult{
							Point:  orb.Point{o.Lon, o.Lat},
							NodeId: o.ID,
							WayId:  result.WayId,
						})
					}
				}
			case *osm.Way:
				break
			case *osm.Relation:
				break
			}
		}
	}

	// Sort the ring following wayIds, nodeIds order
	sort.Slice(points, func(i, j int) bool {
		if points[i].WayId != points[j].WayId {
			return slices.Index(wayIds, points[i].WayId) < slices.Index(wayIds, points[j].WayId)
		}

		var wayResult *WayResult
		for _, result := range wayResults {
			if result.WayId == points[i].WayId {
				wayResult = &result
				break
			}
		}

		return slices.Index(wayResult.NodeIds, points[i].NodeId) < slices.Index(wayResult.NodeIds, points[j].NodeId)
	})

	collection := geojson.NewFeatureCollection()
	lines := map[osm.WayID]orb.LineString{}

	for _, point := range points {
		lines[point.WayId] = append(lines[point.WayId], point.Point)
	}

	ring := orb.Ring{}
	for _, point := range points {
		ring = append(ring, point.Point)
	}

	poly := orb.Polygon{ring}
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

type PointResult struct {
	Point  orb.Point
	NodeId osm.NodeID
	WayId  osm.WayID
}

type WayResult struct {
	WayId   osm.WayID
	NodeIds []osm.NodeID
}
