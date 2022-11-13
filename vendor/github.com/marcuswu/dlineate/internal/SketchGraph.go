package core

import (
	"fmt"
	"sort"

	"github.com/marcuswu/dlineate/internal/constraint"
	el "github.com/marcuswu/dlineate/internal/element"
	"github.com/marcuswu/dlineate/internal/solver"
	iutils "github.com/marcuswu/dlineate/internal/utils"
	"github.com/marcuswu/dlineate/utils"
	"github.com/rs/zerolog"
)

// SketchGraph A graph representing a set of 2D sketch elements and constraints
type SketchGraph struct {
	constraints map[uint]*constraint.Constraint
	elements    map[uint]el.SketchElement
	eToC        map[uint][]*constraint.Constraint
	clusters    []*GraphCluster
	baseCluster *GraphCluster
	freeNodes   *utils.Set
	usedNodes   *utils.Set

	state            solver.SolveState
	degreesOfFreedom uint
	conflicting      *utils.Set
}

// NewSketch creates a new sketch for solving
func NewSketch() *SketchGraph {
	g := new(SketchGraph)
	g.eToC = make(map[uint][]*Constraint, 0)
	g.constraints = make(map[uint]*Constraint, 0)
	g.elements = make(map[uint]el.SketchElement, 0)
	g.clusters = make([]*GraphCluster, 0)
	g.freeNodes = utils.NewSet()
	g.usedNodes = utils.NewSet()
	g.state = solver.None
	g.degreesOfFreedom = 6
	g.conflicting = utils.NewSet()

	g.baseCluster = NewGraphCluster(0)
	c := NewGraphCluster(0)
	g.clusters = append(g.clusters, c)
	return g
}

// GetElement gets an element from the graph
func (g *SketchGraph) GetElement(id uint) (el.SketchElement, bool) {
	e, ok := g.elements[id]
	return e, ok
}

// GetConstraint gets a constraint from the graph
func (g *SketchGraph) GetConstraint(id uint) (*constraint.Constraint, bool) {
	c, ok := g.constraints[id]
	return c, ok
}

func (g *SketchGraph) FindConstraints(elementId uint) []*constraint.Constraint {
	return g.eToC[elementId]
}

func (g *SketchGraph) MakeFixed(e el.SketchElement) {
	g.addElementToCluster(g.clusters[0], e)
	g.addElementToCluster(g.baseCluster, e)
}

// AddPoint adds a point to the sketch
func (g *SketchGraph) AddPoint(x float64, y float64) el.SketchElement {
	elementID := uint(len(g.elements))
	utils.Logger.Debug().
		Float64("X", x).
		Float64("Y", y).
		Uint("id", elementID).
		Msg("Adding point")
	p := el.NewSketchPoint(elementID, x, y)
	g.freeNodes.Add(elementID)
	g.elements[elementID] = p
	return p
}

// AddLine adds a line to the sketch
func (g *SketchGraph) AddLine(a float64, b float64, c float64) el.SketchElement {
	elementID := uint(len(g.elements))
	utils.Logger.Debug().
		Float64("A", a).
		Float64("B", b).
		Float64("C", c).
		Uint("id", elementID).
		Msg("Adding line")
	l := el.NewSketchLine(elementID, a, b, c)
	g.freeNodes.Add(elementID)
	g.elements[elementID] = l
	return l
}

func (g *SketchGraph) AddOrigin(x float64, y float64) el.SketchElement {
	elementID := uint(len(g.elements))
	utils.Logger.Debug().
		Float64("X", x).
		Float64("Y", y).
		Uint("id", elementID).
		Msgf("Adding origin")
	ax := el.NewSketchPoint(elementID, x, y)
	g.freeNodes.Add(elementID)
	g.elements[elementID] = ax

	g.MakeFixed(ax)

	return ax
}

func (g *SketchGraph) AddAxis(a float64, b float64, c float64) el.SketchElement {
	elementID := uint(len(g.elements))
	utils.Logger.Debug().
		Float64("A", a).
		Float64("B", b).
		Float64("C", c).
		Uint("id", elementID).
		Msg("Adding axis")
	ax := el.NewSketchLine(elementID, a, b, c)
	g.freeNodes.Add(elementID)
	g.elements[elementID] = ax

	g.MakeFixed(ax)

	return ax
}

func (g *SketchGraph) IsElementSolved(e el.SketchElement) bool {
	constraints := g.eToC[e.GetID()]
	if len(constraints) < 2 {
		return false
	}

	numSolved := 0
	for _, c := range constraints {
		if c.Solved {
			numSolved++
		}
	}

	return numSolved > 1
}

func (g *SketchGraph) CombinePoints(e1 el.SketchElement, e2 el.SketchElement) el.SketchElement {
	utils.Logger.Debug().
		Uint("element 1", e1.GetID()).
		Uint("removing element 2", e2.GetID()).
		Msg("Combining elements")
	newE2 := e1
	newE1 := e1
	if g.clusters[0].HasElement(e1) {
		newE2 = el.CopySketchElement(e1)
	}
	if g.clusters[0].HasElement(e2) {
		newE1 = el.CopySketchElement(e2)
	}
	// Look for any constraints referencing e2, replace with e1
	for _, constraint := range g.constraints {
		if constraint.Element1.GetID() == e2.GetID() {
			constraint.Element1 = newE2
		}
		if constraint.Element2.GetID() == e2.GetID() {
			constraint.Element2 = newE2
		}
		if constraint.Element1.GetID() == e1.GetID() && !e1.Is(newE1) {
			constraint.Element1 = newE1
		}
		if constraint.Element2.GetID() == e1.GetID() && !e1.Is(newE1) {
			constraint.Element2 = newE1
		}
	}
	// remove e2 from freenodes, elements
	g.eToC[newE1.GetID()] = append(g.eToC[newE1.GetID()], g.eToC[e2.GetID()]...)
	delete(g.eToC, e2.GetID())
	g.freeNodes.Remove(e2.GetID())
	delete(g.elements, e2.GetID())
	return newE1
}

// AddConstraint adds a constraint to sketch elements
func (g *SketchGraph) AddConstraint(t constraint.Type, e1 el.SketchElement, e2 el.SketchElement, value float64) *constraint.Constraint {
	constraintID := uint(len(g.constraints))
	cType := "Distance"
	if t != constraint.Distance {
		cType = "Angle"
	}
	utils.Logger.Debug().
		Str("type", cType).
		Float64("value", value).
		Uint("constraint id", constraintID).
		Msg("Adding constraint")
	constraint := constraint.NewConstraint(constraintID, t, e1, e2, value, false)
	if g.clusters[0].HasElement(e1) && g.clusters[0].HasElement(e2) {
		g.clusters[0].AddConstraint(constraint)
		g.baseCluster.AddConstraint(constraint)
		g.constraints[constraintID] = constraint
		g.eToC[e1.GetID()] = append(g.eToC[e1.GetID()], constraint)
		g.eToC[e2.GetID()] = append(g.eToC[e2.GetID()], constraint)
		return constraint
	}
	g.constraints[constraintID] = constraint
	if _, ok := g.eToC[e1.GetID()]; !ok {
		g.eToC[e1.GetID()] = make([]*Constraint, 0)
	}
	if _, ok := g.eToC[e2.GetID()]; !ok {
		g.eToC[e2.GetID()] = make([]*Constraint, 0)
	}
	g.eToC[e1.GetID()] = append(g.eToC[e1.GetID()], constraint)
	g.eToC[e2.GetID()] = append(g.eToC[e2.GetID()], constraint)
	g.freeNodes.Add(e1.GetID())
	g.freeNodes.Add(e2.GetID())
	return constraint
}

func (g *SketchGraph) findStartConstraint() uint {
	constraints := make([]uint, 0)
	for constraintId, constraint := range g.constraints {
		// If we have a constraint where both elements are used, that's our constraint
		if g.usedNodes.Contains(constraint.Element1.GetID()) &&
			g.usedNodes.Contains(constraint.Element2.GetID()) {
			return constraintId
		}
		if g.usedNodes.Contains(constraint.Element1.GetID()) ||
			g.usedNodes.Contains(constraint.Element2.GetID()) {
			constraints = append(constraints, constraintId)
		}
	}

	// Check unused elements in constraints for highest constraint count
	var retVal uint
	ccount := 0
	sort.Sort(iutils.IdList(constraints))
	for _, constraintId := range constraints {
		eId := g.constraints[constraintId].Element1.GetID()
		if g.usedNodes.Contains(eId) {
			eId = g.constraints[constraintId].Element2.GetID()
		}
		if len(g.eToC[eId]) < ccount {
			continue
		}

		retVal = constraintId
		ccount = len(g.eToC[eId])
	}

	return retVal
}

// find a pair of free constraints which is connected to a single element (might be in another cluster)
// and each constraint shares an element with the cluster we're creating
func (g *SketchGraph) findConstraints(c *GraphCluster) ([]uint, uint, bool) {
	// First, find free constraints connected to the cluster, grouped by the element not in the cluster
	constraints := make(map[uint]*utils.Set)
	for _, constraint := range g.constraints {
		// Skip constraints with no connection to the cluster
		if !c.HasElementID(constraint.First().GetID()) && !c.HasElementID(constraint.Second().GetID()) {
			continue
		}
		// Skip constraints completely contained in the cluster
		if c.HasElementID(constraint.First().GetID()) && c.HasElementID(constraint.Second().GetID()) {
			continue
		}
		other := constraint.First().GetID()
		if c.HasElementID(other) {
			other = constraint.Second().GetID()
		}

		if _, ok := constraints[other]; !ok {
			constraints[other] = utils.NewSet()
		}
		// fmt.Printf("findConstraints: element %d adding constraint %d\n", other, constraint.GetID())
		constraints[other].Add(constraint.GetID())
	}

	var first bool = true
	var element uint
	var targetConstraints []uint
	for eId, cs := range constraints {
		if cs.Count() < 2 {
			continue
		}
		if cs.Count() == 2 {
			return cs.Contents(), eId, true
		}
		if !first && cs.Count() < len(targetConstraints) {
			continue
		}

		element = eId
		targetConstraints = cs.Contents()
		first = false
	}

	return targetConstraints, element, len(targetConstraints) > 1
}

func (g *SketchGraph) addElementToCluster(c *GraphCluster, e el.SketchElement) {
	c.AddElement(e)
	g.freeNodes.Remove(e.GetID())
	g.usedNodes.Add(e.GetID())
}

func (g *SketchGraph) addConstraintToCluster(c *GraphCluster, constraint *constraint.Constraint) {
	g.addElementToCluster(c, constraint.Element1)
	g.addElementToCluster(c, constraint.Element2)
	c.AddConstraint(constraint)
	delete(g.constraints, constraint.GetID())
}

func (g *SketchGraph) createCluster(first uint, id int) *GraphCluster {
	c := NewGraphCluster(id)

	// Add elements connected to other elements in the cluster by two constraints
	clusterNum := len(g.clusters)
	oc, ok := g.GetConstraint(first)
	if !ok {
		utils.Logger.Error().
			Int("cluster", clusterNum).
			Uint("constraint", first).
			Msgf("createCluster(%d): Failed to find initial constraint", clusterNum)
		return nil
	}
	g.addConstraintToCluster(c, oc)

	// find a pair of free constraints which is connected to an element (might be in another cluster)
	// and each constraint shares an element with the cluster we're creating
	for cIds, eId, ok := g.findConstraints(c); ok; cIds, eId, ok = g.findConstraints(c) {
		utils.Logger.Debug().
			Int("cluster", clusterNum).
			Uint("element", eId).
			Int("constraint count", len(cIds)).
			Bool("found ok", ok).
			Msgf("createCluster(%d): adding element", clusterNum)
		level := el.FullyConstrained
		if len(cIds) > 2 {
			level = el.OverConstrained
			// These constraints are conflicting, add them to the conflicting list
			g.conflicting.AddList(cIds)
		}
		g.elements[eId].SetConstraintLevel(level)
		c.AddElement(g.elements[eId])
		for _, cId := range cIds[:2] {
			utils.Logger.Debug().
				Int("cluster", clusterNum).
				Uint("constraint", cId).
				Msgf("createCluster(%d): adding constraint", clusterNum)
			oc, _ = g.GetConstraint(cId)
			g.addConstraintToCluster(c, oc)
		}
	}

	utils.Logger.Info().
		Int("cluster", clusterNum).
		Int("element count", len(c.elements)).
		Int("constraint count", len(c.constraints)).
		Msgf("createCluster(%d) completed building cluster", clusterNum)
	g.clusters = append(g.clusters, c)

	return c
}

func (g *SketchGraph) createClusters() {
	id := 1
	utils.Logger.Info().
		Int("unassigned constraints", len(g.constraints)).
		Msg("Creating clusters")
	for lastLen := len(g.constraints) + 1; len(g.constraints) > 0 && lastLen != len(g.constraints); {
		lastLen = len(g.constraints)
		// Find constraint to begin new cluster
		g.createCluster(g.findStartConstraint(), id)
		id++
		utils.Logger.Info().Msgf("%d unassigned constraints left\n", len(g.constraints))
	}
	utils.Logger.Info().
		Int("unassigned constraints", len(g.constraints)).
		Msg("Created clusters")
}

func (g *SketchGraph) logElements(level zerolog.Level) {
	utils.Logger.WithLevel(level).Msg("Elements: ")
	for _, e := range g.elements {
		utils.Logger.WithLevel(level).Msgf("%v", e)
	}
	utils.Logger.WithLevel(level).Msg("")
}

func (g *SketchGraph) logConstraints(level zerolog.Level) {
	utils.Logger.WithLevel(level).Msg("Constraints: ")
	for _, c := range g.constraints {
		utils.Logger.WithLevel(level).Msgf("%v", c)
	}
	utils.Logger.WithLevel(level).Msg("")
}

func (g *SketchGraph) logConstraintsElements(level zerolog.Level) {
	g.logElements(level)
	g.logConstraints(level)
	utils.Logger.WithLevel(level).Msg("")
}

func (g *SketchGraph) updateElements(c *GraphCluster) {
	for eId, e := range c.elements {
		g.elements[eId] = e
	}
}

func (g *SketchGraph) addClusterConstraints(c *GraphCluster) {
	for _, constraint := range c.constraints {
		g.constraints[constraint.GetID()] = constraint
	}
}

// Restore constraints and elements in the graph
func (g *SketchGraph) ResetClusters() {
	for _, cluster := range g.clusters {
		g.addClusterConstraints(cluster)
	}
	g.clusters = make([]*GraphCluster, 0)
	g.freeNodes.AddSet(g.usedNodes)
	g.usedNodes.Clear()
	g.state = solver.None

	// Recreate cluster 0
	// This isn't great -- if cluster 0 changes over time, this will need to be updated
	// Instead we should track a pristine version of cluster 0 to reset to
	c := NewGraphCluster(0)
	for _, element := range g.baseCluster.elements {
		c.AddElement(element)
	}
	for _, constraint := range g.baseCluster.constraints {
		c.AddConstraint(constraint)
	}
	g.clusters = append(g.clusters, c)
}

func (g *SketchGraph) BuildClusters() {
	for _, c := range g.clusters[0].constraints {
		g.addElementToCluster(g.clusters[0], c.First())
		g.addElementToCluster(g.clusters[0], c.Second())
		delete(g.constraints, c.GetID())
	}
	elements := utils.NewSet()
	for _, c := range g.constraints {
		elements.AddList(c.ElementIDs())
	}
	// Add back in constraints from cluster where both elements
	// are referenced in constraints not in the cluster
	if len(g.clusters) == 1 {
		for _, c := range g.clusters[0].constraints {
			if elements.Contains(c.Element1.GetID()) && elements.Contains(c.Element2.GetID()) {
				utils.Logger.Trace().
					Uint("constraint", c.GetID()).
					Uint("element 1", c.First().GetID()).
					Uint("element 2", c.Second().GetID()).
					Msg("Adding back in constraint")
				g.constraints[c.GetID()] = c
			}
		}
	}
	g.logConstraintsElements(zerolog.InfoLevel)
	if len(g.clusters) == 1 {
		g.createClusters()
	}
}

// Solve builds the graph and solves the sketch
func (g *SketchGraph) Solve() solver.SolveState {
	defer g.logElements(zerolog.DebugLevel)

	utils.Logger.Info().
		Int("cluster count", len(g.clusters)).
		Int("constraint count", len(g.constraints)).
		Msg("Beginning cluster solves")
	for i, c := range g.clusters {
		if i == 0 {
			continue
		}
		utils.Logger.Info().
			Int("cluster", i).
			Msg("Starting cluster solve")
		clusterState := c.Solve()
		g.updateElements(c)
		g.addClusterConstraints(c)
		utils.Logger.Info().
			Int("cluster", i).
			Str("cluster state", clusterState.String()).
			Str("graph state", g.state.String()).
			Msg("Solved cluster")
		c.logElements(zerolog.TraceLevel)
		if g.state == solver.None || (g.state != clusterState && !(g.state != solver.Solved && clusterState == solver.Solved)) {
			utils.Logger.Info().
				Int("cluster", i).
				Str("new state", clusterState.String()).
				Msg("Updating graph state after cluster solve")
			g.state = clusterState
		}
		utils.Logger.Debug().
			Str("graph state", g.state.String()).
			Msg("Current graph solve state")
	}
	// Merge clusters
	utils.Logger.Info().Msg("Starting Cluster Merges")
	removeCluster := func(g *SketchGraph, cIndex int) {
		last := len(g.clusters) - 1
		g.clusters[cIndex], g.clusters[last] = g.clusters[last], g.clusters[cIndex]
		g.clusters = g.clusters[:last]
	}
	for first, second, third := g.findMerge(); first > 0 && second > 0; first, second, third = g.findMerge() {
		utils.Logger.Debug().
			Int("first cluster", first).
			Int("second cluster", second).
			Int("third cluster", third).
			Msg("Found merge")
		c1 := g.clusters[first]
		c2 := g.clusters[second]
		var c3 *GraphCluster = nil
		if third > 0 {
			c3 = g.clusters[third]
		}
		mergeState := c1.solveMerge(c2, c3)
		utils.Logger.Debug().
			Str("state", fmt.Sprintf("%v", mergeState)).
			Msg("Completed merge")
		for _, c := range g.clusters {
			c.IsSolved()
		}
		g.updateElements(c1)
		g.addClusterConstraints(c1)
		for _, c := range g.clusters {
			c.IsSolved()
		}
		if g.state != mergeState && mergeState != solver.Solved {
			utils.Logger.Debug().
				Str("graph state", mergeState.String()).
				Msg("Updating state after cluster merge")
			g.state = mergeState
		}
		// Remove second and third clusters
		ordered := []int{second, third}
		if second < third {
			ordered[0], ordered[1] = ordered[1], ordered[0]
		}
		removeCluster(g, ordered[0])
		if third > 0 {
			removeCluster(g, ordered[1])
		}
	}
	utils.Logger.Info().Msg("Merging with origin and X & Y axes")
	if len(g.clusters) < 2 {
		return g.state
	}
	mergeState := g.clusters[0].mergeOne(g.clusters[1], false)
	g.updateElements(g.clusters[0])
	g.updateElements(g.clusters[1])
	g.addClusterConstraints(g.clusters[0])
	g.addClusterConstraints(g.clusters[1])
	if g.state != mergeState && mergeState != solver.Solved {
		utils.Logger.Debug().
			Str("graph state", mergeState.String()).
			Msg("Updating state after cluster merge")
		g.state = mergeState
	}
	utils.Logger.Debug().
		Str("graph state", g.state.String()).
		Msg("Final graph state")

	if !g.IsSolved() {
		g.state = solver.NonConvergent
	}

	return g.state
}

func (g *SketchGraph) findMergeForCluster(c *GraphCluster) (int, int) {
	connectedClusters := func(g *SketchGraph, c *GraphCluster) map[int][]uint {
		connected := make(map[int][]uint)
		for i, other := range g.clusters {
			if other.id == c.id {
				continue
			}
			shared := c.SharedElements(other).Contents()
			if len(shared) < 1 {
				continue
			}
			connected[i] = shared
		}
		return connected
	}

	// These are the clusters connected to c
	// We want to find either one connected to c by two elements
	// or two within this list connected to each other by one
	connected := connectedClusters(g, c)
	for ci, shared := range connected {
		if ci == 0 {
			continue
		}
		utils.Logger.Debug().
			Int("cluster 1", c.id).
			Int("cluster 2", g.clusters[ci].id).
			Msg("Looking for merge")
		if len(shared) == 2 {
			utils.Logger.Debug().
				Int("cluster", g.clusters[ci].id).
				Msg("Found connected cluster for merge")
			return ci, -1
		}

		if len(shared) == 1 {
			// Find another cluster in connected that is connected to g.clusters[ci]
			for oi, oshared := range connected {
				if oi == 0 || ci == oi || len(oshared) != 1 || oshared[0] == shared[0] {
					continue
				}
				utils.Logger.Debug().
					Int("cluster 0", c.id).
					Int("cluster 1", g.clusters[ci].id).
					Int("cluster 2", g.clusters[oi].id).
					Msg("Testing for valid merge for clusters")
				ciOiShared := g.clusters[ci].SharedElements(g.clusters[oi])
				if ciOiShared.Count() == 1 && !ciOiShared.Contains(shared[0]) && !ciOiShared.Contains(oshared[0]) {
					utils.Logger.Debug().
						Int("cluster 0", c.id).
						Int("cluster 1", g.clusters[ci].id).
						Int("cluster 2", g.clusters[oi].id).
						Msg("Found connected clusters for merge")
					return ci, oi
				}
			}
		}
	}

	return -1, -1
}

// Find and return clusters which can be merged.
// This can either be:
//   - Two clusters each sharing an element with g and sharing an element with each other
//   - One cluster sharing two elements with g
//
// Returns the index(es) of the mergable clusters
func (g *SketchGraph) findMerge() (int, int, int) {
	for i, c := range g.clusters {
		// Merge cluster 0 last manually
		if i == 0 {
			continue
		}
		utils.Logger.Debug().
			Int("start cluster", c.id).
			Msg("Looking for merge")
		c1, c2 := g.findMergeForCluster(c)
		if c1 >= 0 {
			return i, c1, c2
		}
	}

	utils.Logger.Debug().Msg("No merge found")
	return -1, -1, -1
}

func (g *SketchGraph) IsSolved() bool {
	solved := true
	for _, c := range g.constraints {
		if c.IsMet() {
			continue
		}

		utils.Logger.Trace().
			Str("constraint", c.String()).
			Msg("Failed to meet constraint")
		solved = false
	}

	return solved
}

func (g *SketchGraph) Conflicting() *utils.Set {
	return g.conflicting
}

func (g *SketchGraph) ToGraphViz() string {
	edges := ""
	uniqueSharedElements := make(map[string]interface{})
	sharedElements := ""

	// Output clusters
	for _, c := range g.clusters {
		edges = edges + c.ToGraphViz()
		for _, other := range g.clusters {
			if c.id == other.id {
				continue
			}
			shared := c.SharedElements(other)
			if shared.Count() == 0 {
				continue
			}
			first := c.id
			second := other.id
			if second < first {
				first, second = second, first
			}
			for _, e := range shared.Contents() {
				key := fmt.Sprintf("\t\"%d-%d\" -- \"%d-%d\"\n", first, e, second, e)
				if _, ok := uniqueSharedElements[key]; ok {
					continue
				}
				sharedElements = sharedElements + key
				uniqueSharedElements[key] = 0
			}
		}
	}

	// Output free constraints
	for _, c := range g.constraints {
		edges = edges + c.ToGraphViz(-1)
	}

	// Output free elements
	for _, eId := range g.freeNodes.Contents() {
		edges = edges + g.elements[eId].ToGraphViz(-1)
	}

	return fmt.Sprintf(`
	graph {
		compound=true
		%s
		%s
	}`, edges, sharedElements)
}
