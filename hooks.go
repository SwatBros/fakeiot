package fakeiot

type Hook interface {
	CheckCondition(g *Generator) bool
	Run(g *Generator) error
}
