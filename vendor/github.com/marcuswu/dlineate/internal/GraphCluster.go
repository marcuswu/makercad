package core

import (
	"fmt"

	"github.com/marcuswu/dlineate/internal/constraint"
	el "github.com/marcuswu/dlineate/internal/element"
	"github.com/marcuswu/dlineate/internal/solver"
	"github.com/marcuswu/dlineate/utils"
	"github.com/rs/zerolog"
)

// Constraint a convenient alias for cosntraint.Constraint
type Constraint = constraint.Constraint

// GraphCluster A cluster within a Graph
type GraphCluster struct {
	id          int
	constraints map[uint]*Constraint
	// others      []*GraphCluster
	elements   map[uint]el.SketchElement
	eToC       map[uint][]*Constraint
	solveOrder []uint
	solved     *utils.Set
}

// NewGraphCluster constructs a new GraphCluster
func NewGraphCluster(id int) *GraphCluster {
	g := new(GraphCluster)
	g.id = id
	g.constraints = make(map[uint]*Constraint, 0)
	// g.others = make([]*GraphCluster, 0)
	g.elements = make(map[uint]el.SketchElement, 0)
	g.eToC = make(map[uint][]*Constraint)
	g.solved = utils.NewSet()
	g.solveOrder = make([]uint, 0)
	return g
}

func (g *GraphCluster) GetID() int {
	return g.id
}

func (g *GraphCluster) AddElement(e el.SketchElement) {
	if _, ok := g.elements[e.GetID()]; ok {
		return
	}
	utils.Logger.Debug().
		Int("cluster", g.id).
		Uint("element id", e.GetID()).
		Msg("Cluster adding element")
	g.elements[e.GetID()] = el.CopySketchElement(e)
	g.solveOrder = append(g.solveOrder, e.GetID())
}

// AddConstraint adds a constraint to the cluster
func (g *GraphCluster) AddConstraint(c *Constraint) {
	gc, ok := g.constraints[c.GetID()]
	if ok {
		if c.Solved && !gc.Solved {
			gc.Solved = c.Solved
		}
	} else {
		gc = constraint.CopyConstraint(c)
		g.constraints[gc.GetID()] = gc
	}
	if _, ok := g.elements[gc.Element1.GetID()]; !ok {
		g.elements[gc.Element1.GetID()] = gc.Element1
	} else {
		gc.Element1 = g.elements[gc.Element1.GetID()]
	}

	if _, ok := g.elements[gc.Element2.GetID()]; !ok {
		g.elements[gc.Element2.GetID()] = gc.Element2
	} else {
		gc.Element2 = g.elements[gc.Element2.GetID()]
	}

	if _, ok := g.eToC[gc.Element1.GetID()]; !ok {
		g.eToC[gc.Element1.GetID()] = make([]*Constraint, 0)
	}
	if _, ok := g.eToC[gc.Element2.GetID()]; !ok {
		g.eToC[gc.Element2.GetID()] = make([]*Constraint, 0)
	}
	g.eToC[gc.Element1.GetID()] = append(g.eToC[gc.Element1.GetID()], gc)
	g.eToC[gc.Element2.GetID()] = append(g.eToC[gc.Element2.GetID()], gc)
}

// HasElementIDDirect returns whether this cluster directly contains an element ID
func (g *GraphCluster) HasElementIDImmediate(eID uint) bool {
	_, e := g.elements[eID]
	return e
}

// HasElementID returns whether this cluster contains an element ID
func (g *GraphCluster) HasElementID(eID uint) bool {
	if _, e := g.elements[eID]; e {
		return true
	}
	return false
}

// HasElement returns whether this cluster contains an element
func (g *GraphCluster) HasElement(e el.SketchElement) bool {
	if e == nil {
		return true
	}
	return g.HasElementID(e.GetID())
}

// GetElement returns the copy of an element represented in this cluster
func (g *GraphCluster) GetElement(eID uint) (el.SketchElement, bool) {
	if element, ok := g.elements[eID]; ok {
		return element, ok
	}
	return nil, false
}

// SharedElements returns the shared elements between this and another cluster
func (g *GraphCluster) SharedElements(gc *GraphCluster) *utils.Set {
	var shared *utils.Set = utils.NewSet()

	for elementID := range g.elements {
		if gc.HasElementID(elementID) {
			shared.Add(elementID)
		}
	}

	return shared
}

func (g *GraphCluster) immediateSharedElements(gc *GraphCluster) *utils.Set {
	var shared *utils.Set = utils.NewSet()

	for elementID := range g.elements {
		if gc.HasElementIDImmediate(elementID) {
			shared.Add(elementID)
		}
	}

	return shared
}

// Translate translates all elements in the cluster by an x and y value
func (g *GraphCluster) Translate(xDist float64, yDist float64) {
	for _, e := range g.elements {
		e.Translate(xDist, yDist)
	}
}

// Rotate rotates all elements in the cluster around a point by an angle
func (g *GraphCluster) Rotate(origin *el.SketchPoint, angle float64) {
	v := el.Vector{X: origin.GetX(), Y: origin.GetY()}
	for _, e := range g.elements {
		e.Translate(-v.X, -v.Y)
		e.Rotate(angle)
		e.Translate(v.X, v.Y)
	}
}

func (g *GraphCluster) rebuildMap() {
	g.elements = make(map[uint]el.SketchElement, 0)

	for _, c := range g.constraints {
		g.elements[c.Element1.GetID()] = c.Element1
		g.elements[c.Element2.GetID()] = c.Element2
	}
}

func (g *GraphCluster) solvedConstraintsFor(eID uint) []*Constraint {
	constraints := g.eToC[eID]
	var solvedC = make([]*Constraint, 0)
	for _, c := range constraints {
		if g.solved.Contains(c.GetID()) {
			solvedC = append(solvedC, c)
		}
	}
	return solvedC
}

func (g *GraphCluster) unsolvedConstraintsFor(eID uint) constraint.ConstraintList {
	var constraints = g.eToC[eID]
	var unsolved = make([]*Constraint, 0)
	for _, c := range constraints {
		if g.solved.Contains(c.GetID()) {
			continue
		}
		unsolved = append(unsolved, c)
	}

	return unsolved
}

// LocalSolve attempts to solve the constraints in the cluster, returns solution state
func (g *GraphCluster) localSolve() solver.SolveState {
	// solver changes element instances in constraints, so rebuild the element map
	defer g.rebuildMap()
	// Order constraints for element 0
	if len(g.solveOrder) < 2 {
		return solver.NonConvergent
	}
	g.logElements(zerolog.InfoLevel)

	state := solver.Solved

	e1 := g.solveOrder[0]
	e2 := g.solveOrder[1]
	g.solveOrder = g.solveOrder[2:]
	utils.Logger.Debug().Msg("Local Solve Step 0")
	for _, c := range g.constraints {
		if !c.HasElements(e1, e2) {
			continue
		}
		utils.Logger.Info().
			Uint("constraint", c.GetID()).
			Uint("element 1", e1).
			Uint("element 2", e2).
			Msg("Solving constraint betw first two elements")
		state = solver.SolveConstraint(c)
		utils.Logger.Trace().
			Str("state", state.String()).
			Msg("State")
		g.solved.Add(c.GetID())
		break
	}

	/*
		1. Look for point w/ 2 constraints to solved elements -- fall back to point w/ fewest unsolved constraints
		2. Solve the element by those 2 constraints
		3. If there are unsolved elements, go to step 1

		An element is considered solved when it has at least two solved constraints.
		A constraint needs a solved flag or a structure to track solved state
		Need to be able to get constraints for an element
		Need to be able to filter constraint list by solved / unsolved (get by state?)
		Need to be able to quickly determine if an element is solved

		solved = Set of constraint
		map[elementID][constraint]
		isElementSolved(elementID)
	*/

	// Pick 2 from constraintList and solve. If only 1 in constraintList, solve just the one

	for len(g.solveOrder) > 0 {
		// Step 1
		utils.Logger.Debug().Msg("Local Solve Step 1")
		utils.Logger.Trace().
			Str("elements", fmt.Sprintf("%v", g.solveOrder)).
			Msg("Solve Order")
		e := g.solveOrder[0]
		g.solveOrder = g.solveOrder[1:]
		c := g.unsolvedConstraintsFor(e)

		if len(g.solvedConstraintsFor(e)) >= 2 {
			utils.Logger.Trace().
				Uint("element", e).
				Msg("Element already solved. Continuing.")
			continue
		}

		utils.Logger.Debug().
			Uint("element", e).
			Msg("Solving for element")
		utils.Logger.Trace().
			Uint("element", e).
			Array("constraints", c).
			Msg("Element's eligible constraints")
		if len(c) < 2 {
			utils.Logger.Error().
				Int("unsolved constraints", len(g.constraints)-g.solved.Count()).
				Msg("Could not find a constraint to solve")
			state = solver.NonConvergent
			break
		}

		// Step 2
		utils.Logger.Debug().Msg("Local Solve Step 2")
		utils.Logger.Debug().
			Uint("constraint 1", c[0].GetID()).
			Uint("constraint 2", c[1].GetID()).
			Msg("Solving constraints")
		if s := solver.SolveConstraints(c[0], c[1], g.elements[e]); state == solver.Solved {
			utils.Logger.Trace().
				Str("state", s.String()).
				Msg("solve state changed")
			utils.Logger.Debug().
				Str("element", g.elements[e].String()).
				Msg("solved element")
			element, _ := c[0].Element(e)
			utils.Logger.Trace().
				Str("element", element.String()).
				Msg("solved element in constraint 1")
			element, _ = c[1].Element(e)
			utils.Logger.Trace().
				Str("element", element.String()).
				Msg("solved element in constraint 2")
			state = s
			utils.Logger.Trace().
				Str("state", state.String()).
				Msg("State")
		}
		g.solved.Add(c[0].GetID())
		g.solved.Add(c[1].GetID())

		utils.Logger.Info().
			Str("solve ratio", fmt.Sprintf("%d / %d", g.solved.Count(), len(g.constraints))).
			Msg("Local Solve Step 3 (check for completion)")
	}

	utils.Logger.Info().
		Str("state", state.String()).
		Msg("finished")
	g.logElements(zerolog.InfoLevel)
	return state
}

func (g *GraphCluster) logElements(level zerolog.Level) {
	for _, e := range g.elements {
		g.logElement(e, level)
	}
}

func (g *GraphCluster) logElement(e el.SketchElement, level zerolog.Level) {
	utils.Logger.WithLevel(level).Msg(e.String())
}

// MergeOne resolves merging one solved child clusters to this one
func (g *GraphCluster) mergeOne(other *GraphCluster, mergeConstraints bool) solver.SolveState {
	if mergeConstraints {
		defer g.mergeConstraints(other, nil)
	}
	sharedElements := g.immediateSharedElements(other).Contents()

	if g.id == 0 && other.id == 1 && len(sharedElements) > 2 {
		sharedElements = []uint{0, 1}
	}

	if len(sharedElements) != 2 {
		return solver.NonConvergent
	}

	// Solve two shared elements
	utils.Logger.Debug().Msg("Initial configuration:")
	utils.Logger.Debug().
		Str("elements", fmt.Sprintf("%v", sharedElements)).
		Msg("Shared elements")
	g.logElements(zerolog.DebugLevel)
	utils.Logger.Debug().Msg("")
	other.logElements(zerolog.DebugLevel)
	utils.Logger.Debug().Msg("")

	first := sharedElements[0]
	second := sharedElements[1]

	if g.elements[first].GetType() == el.Line {
		first, second = second, first
	}

	// If both elements are lines, nonconvergent (I think)
	if g.elements[first].GetType() == el.Line {
		utils.Logger.Error().Msg("In a merge one and both shared elements are line type")
		return solver.NonConvergent
	}

	p1 := g.elements[first]
	p2 := other.elements[first]

	// If there's a line, first rotate the lines into the same angle, then match first element
	if g.elements[second].GetType() == el.Line {
		angle := other.elements[second].AsLine().AngleToLine(g.elements[second].AsLine())
		other.Rotate(p1.AsPoint(), angle)
		utils.Logger.Trace().Msg("Rotated to make line the same angle")
	}

	// Match up the first point
	utils.Logger.Trace().Msg("matching up the first point")
	direction := p1.VectorTo(p2)
	other.Translate(direction.X, direction.Y)

	// If both are points, rotate other to match the element in g
	if g.elements[second].GetType() == el.Point {
		utils.Logger.Trace().Msg("both elements were points, rotating to match the points together")
		v1 := g.elements[second].VectorTo(g.elements[first])
		v2 := other.elements[second].VectorTo(other.elements[first])
		angle := v1.AngleTo(v2)
		other.Rotate(p1.AsPoint(), angle)
	}

	return solver.Solved
}

func (g *GraphCluster) mergeConstraints(c1 *GraphCluster, c2 *GraphCluster) {
	if c1 != nil {
		for _, c := range c1.constraints {
			g.AddConstraint(c)
		}
	}
	if c2 != nil {
		for _, c := range c2.constraints {
			g.AddConstraint(c)
		}
	}
}

// SolveMerge resolves merging two solved child clusters to this one
/* TODO: Rewrite this. I originally wrote this when I couldn't solve for a line and had to
solve lines separately and then solve for a point. Now I can solve for a line.

1. Find elements in g shared with c1 and c2
2. Solve c1 and c2 shared elements (moving them to g)
3. Find element shared between c1 and c2 -- this is what we're solving for
4. Construct two constraints from g to c1 and g to c2 based on c1 and c2's shared element
5. Solve the constraint and rotate c1 and c2 to match
*/
func (g *GraphCluster) solveMerge(c1 *GraphCluster, c2 *GraphCluster) solver.SolveState {
	if c2 == nil {
		utils.Logger.Info().Msg("Beginning one cluster merge")
		return g.mergeOne(c1, true)
	}
	// Move constraints / elements from c1, c2 to g when we're done
	defer g.mergeConstraints(c1, c2)
	utils.Logger.Info().Msg("")
	utils.Logger.Info().Msg("Beginning cluster merge")
	solve := g.IsSolved()
	utils.Logger.Info().Msgf("Checking g solved: %v", solve)
	solve = c1.IsSolved()
	utils.Logger.Info().Msgf("Checking c1 solved: %v", solve)
	solve = c2.IsSolved()
	utils.Logger.Info().Msgf("Checking c2 solved: %v", solve)
	utils.Logger.Info().Msgf("")
	utils.Logger.Debug().Msg("Pre-merge state:")
	utils.Logger.Debug().Msg("g:")
	g.logElements(zerolog.DebugLevel)
	utils.Logger.Debug().Msg("c1:")
	c1.logElements(zerolog.DebugLevel)
	utils.Logger.Debug().Msg("c2:")
	c2.logElements(zerolog.DebugLevel)
	clusters := []*GraphCluster{g, c1, c2}
	sharedSet := g.immediateSharedElements(c1)
	sharedSet.AddSet(g.immediateSharedElements(c2))
	sharedSet.AddSet(c1.immediateSharedElements(c2))
	sharedElements := sharedSet.Contents()
	utils.Logger.Trace().
		Str("elements", fmt.Sprintf("%v", sharedElements)).
		Msg("Solving for shared elements")

	orderClustersFor := func(e uint) []*GraphCluster {
		matching := make([]*GraphCluster, 0)
		for _, c := range clusters {
			if _, ok := c.elements[e]; !ok {
				continue
			}
			matching = append(matching, c)
		}
		return matching
	}

	if len(sharedElements) != 3 {
		return solver.NonConvergent
	}

	numSharedLines := func(g *GraphCluster) int {
		lines := 0
		for _, se := range sharedElements {
			if e, ok := g.elements[se]; ok && e.GetType() == el.Line {
				lines++
			}
		}
		return lines
	}

	// Find root cluster
	// Prefer keeping lines on the root cluster (solve lines first)
	rootCluster := g
	sharedLines := numSharedLines(g)
	c1SharedLines := numSharedLines(c1)
	c2SharedLines := numSharedLines(c2)
	if c1SharedLines > sharedLines {
		rootCluster = c1
		sharedLines = c1SharedLines
	}
	if c2SharedLines > sharedLines {
		rootCluster = c2
	}

	// Solve two of the elements
	final := sharedElements[0]
	finalIndex := 0
	for i, ec := range clusters {
		if ec == rootCluster {
			finalIndex = i
			break
		}
	}
	utils.Logger.Trace().
		Int("cluster", finalIndex).
		Msg("root cluster")

	for _, se := range sharedElements {
		parents := orderClustersFor(se)
		if len(parents) != 2 {
			utils.Logger.Error().
				Uint("element", se).
				Int("number of parents", len(parents)).
				Msg("Shared element has too many parents. Returning Non-Convergent")
			return solver.NonConvergent
		}

		if parents[0] != rootCluster && parents[1] != rootCluster {
			final = se
			continue
		}
		eType := parents[0].elements[se].GetType()
		utils.Logger.Trace().
			Uint("element", se).
			Str("type", eType.String()).
			Msg("Solving for element")

		// Solve element
		// if element is a line, rotate it into place first
		other := parents[0]
		if other == rootCluster {
			other = parents[1]
		}
		ec1 := other.elements[se]
		ec2 := rootCluster.elements[se]
		var translation *el.Vector
		if eType == el.Line {
			other.logElements(zerolog.TraceLevel)
			utils.Logger.Trace().Msg("")
			angle := ec1.AsLine().AngleToLine(ec2.AsLine())
			other.Rotate(ec1.AsLine().PointNearestOrigin(), angle)
			translation = ec1.VectorTo(ec2)
		} else {
			translation = ec2.VectorTo(ec1)
		}

		// translate element into place
		other.Translate(translation.X, translation.Y)

		utils.Logger.Trace().
			Uint("element", se).
			Msg("Solved for element")
		utils.Logger.Trace().Msg("g:")
		g.logElements(zerolog.TraceLevel)
		utils.Logger.Trace().Msg("c1:")
		c1.logElements(zerolog.TraceLevel)
		utils.Logger.Trace().Msg("c2:")
		c2.logElements(zerolog.TraceLevel)
		utils.Logger.Trace().Msg("")
	}

	var e = [2]uint{sharedElements[0], sharedElements[1]}
	if e[0] == final {
		e[0] = sharedElements[2]
	}
	if e[1] == final {
		e[1] = sharedElements[2]
	}
	utils.Logger.Trace().
		Uint("element 1", e[0]).
		Uint("element 2", e[1]).
		Uint("final unsolved element", final).
		Msg("Solved two elmements")
	g.logElements(zerolog.TraceLevel)
	utils.Logger.Trace().Msg("")
	c1.logElements(zerolog.TraceLevel)
	utils.Logger.Trace().Msg("")
	c2.logElements(zerolog.TraceLevel)
	utils.Logger.Trace().Msg("")

	// Solve the third element in relation to the other two
	parents := orderClustersFor(final)
	finalE := [2]el.SketchElement{parents[0].elements[final], parents[1].elements[final]}
	// p0Final := parents[0].elements[final]
	// p1Final := parents[1].elements[final]
	e2Type := finalE[0].GetType()
	utils.Logger.Trace().
		Str("type", e2Type.String()).
		Msgf("Final element type")
	if e2Type == el.Line {
		// We avoid e2 being a line, so if it is one, the other two are also lines.
		// This means e2 should already be placed correctly since the other two are.
		state := solver.Solved
		finalE[0] = parents[0].elements[final]
		finalE[1] = parents[1].elements[final]
		if !finalE[0].AsLine().IsEquivalent(finalE[1].AsLine()) {
			utils.Logger.Error().
				Str("line 1", finalE[0].String()).
				Str("line 2", finalE[1].String()).
				Msg("Lines are not equivalent: ")
			state = solver.NonConvergent
		}

		return state
	}

	// var constraint1, constraint2 *Constraint
	// var e1, e2 el.SketchElement
	others := [2]el.SketchElement{nil, nil}
	constraints := [2]*constraint.Constraint{nil, nil}
	for pi := range parents {
		for ei := range e {
			finalElement := finalE[pi]
			otherElement, ok := parents[pi].elements[e[ei]]
			if !ok {
				continue
			}
			others[pi] = otherElement
			dist := finalElement.DistanceTo(otherElement)
			constraints[pi] = constraint.NewConstraint(0, constraint.Distance, finalElement, otherElement, dist, false)
			utils.Logger.Trace().
				Uint("element 1", finalElement.GetID()).
				Uint("element 2", otherElement.GetID()).
				Float64("distance", dist).
				Msg("Creating constraint")
		}
	}

	newE3, state := solver.ConstraintResult(constraints[0], constraints[1], finalE[0])
	newP3 := newE3.AsPoint()

	if state != solver.Solved {
		utils.Logger.Error().Msg("Final element solve failed")
		return state
	}

	utils.Logger.Trace().
		Float64("X", newP3.X).
		Float64("Y", newP3.Y).
		Msg("Desired merge point c1 and c2")

	moveCluster := func(c *GraphCluster, pivot el.SketchElement, from *el.SketchPoint, to *el.SketchPoint) {
		if pivot.GetType() == el.Line {
			move := from.VectorTo(to)
			c.Translate(-move.X, -move.Y)
			return
		}

		current, desired := pivot.VectorTo(from), pivot.VectorTo(to)
		angle := current.AngleTo(desired)
		c.Rotate(pivot.AsPoint(), angle)
	}

	utils.Logger.Trace().
		Uint("pivot", others[0].GetID()).
		Str("from", finalE[0].String()).
		Str("to", newP3.String()).
		Msg("Pivoting c0")
	moveCluster(parents[0], others[0], finalE[0].AsPoint(), newP3)
	utils.Logger.Trace().
		Str("parent 0 final", finalE[0].String()).
		Msgf("parent 0 moved")
	utils.Logger.Trace().
		Uint("pivot", others[1].GetID()).
		Str("from", finalE[1].String()).
		Str("to", newP3.String()).
		Msg("Pivoting c1")
	moveCluster(parents[1], others[1], finalE[1].AsPoint(), newP3)
	utils.Logger.Trace().
		Str("parent 1 final", finalE[1].String()).
		Msgf("parent 1 moved")

	utils.Logger.Info().Msg("Completed cluster merge")
	utils.Logger.Info().Msg("")
	utils.Logger.Info().Msg("g:")
	g.logElements(zerolog.InfoLevel)
	utils.Logger.Info().Msg("c1:")
	c1.logElements(zerolog.InfoLevel)
	utils.Logger.Info().Msg("c2:")
	c2.logElements(zerolog.InfoLevel)
	utils.Logger.Info().Msg("")

	if !g.SharedElementsEquivalent(c1) || !g.SharedElementsEquivalent(c2) || !c1.SharedElementsEquivalent(c2) {
		utils.Logger.Info().Msg("Returning Non-convergent due to element inequivalancy after merge")
		return solver.NonConvergent
	}

	return solver.Solved
}

func (g *GraphCluster) SharedElementsEquivalent(o *GraphCluster) bool {
	compareElement := func(e1 el.SketchElement, e2 el.SketchElement) bool {
		if e1.GetType() != e2.GetType() {
			return false
		}

		if e1.AsLine() != nil {
			l1 := e1.AsLine()
			l2 := e2.AsLine()
			return utils.StandardFloatCompare(l1.GetA(), l2.GetA()) == 0 &&
				utils.StandardFloatCompare(l1.GetB(), l2.GetB()) == 0 &&
				utils.StandardFloatCompare(l1.GetC(), l2.GetC()) == 0
		}

		p1 := e1.AsPoint()
		p2 := e2.AsPoint()

		return utils.StandardFloatCompare(p1.X, p2.X) == 0 &&
			utils.StandardFloatCompare(p1.Y, p2.Y) == 0
	}
	equal := true
	shared := g.SharedElements(o)
	for _, e := range shared.Contents() {
		e1 := g.elements[e]
		e2 := o.elements[e]
		equal = equal && compareElement(e1, e2)
	}

	return equal
}

// Solve solves the cluster and any child clusters associated with it
func (g *GraphCluster) Solve() solver.SolveState {
	utils.Logger.Info().
		Int("cluster", g.id).
		Msg("Solving cluster")
	state := g.localSolve()
	return state
}

func (c *GraphCluster) IsSolved() bool {
	solved := true
	for _, c := range c.constraints {
		if c.IsMet() {
			continue
		}

		utils.Logger.Trace().
			Str("constraint", c.String()).
			Msg("Failed to meet")
		solved = false
	}

	return solved
}

func (c *GraphCluster) ToGraphViz() string {
	edges := ""
	elements := ""
	for _, constraint := range c.constraints {
		edges = edges + constraint.ToGraphViz(c.id)
		elements = elements + constraint.Element1.ToGraphViz(c.id)
		elements = elements + constraint.Element2.ToGraphViz(c.id)
	}
	return fmt.Sprintf(`subgraph cluster_%d {
		label = "Cluster %d"
		%s
		%s
	}`, c.id, c.id, edges, elements)
}
