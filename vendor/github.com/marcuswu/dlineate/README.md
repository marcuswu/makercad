# dlineate
A 2D geometric constraint solver written in Go. Dlineate is written using a graph solution based on The [Electronic Primer on Geometric Constraint Solving](https://www.cs.purdue.edu/homes/cmh/electrobook/intro.html) published by Purdue.

## Features
 * Graph based approach
 * Entities
   * Arc
   * Circle
   * Line Segment
   * Point
 * Constraints
   * Angle
   * Distance
   * Coincident
   * Equal
   * Midpoint
   * Parallel
   * Perpendicular
   * Ratio
   * Tangent
   * Horizontal
   * Vertical

## Installation

```
go get -u github.com/marcuswu/dlineate
```

## Getting Started

### Simple Rectangular Sketch Example

```go
    package main

    import (
        "fmt"

        dlineate "github.com/marcuswu/dlineate/pkg"
    )

	sketch := dlineate.NewSketch()
	// Add elements
	l1 := sketch.AddLine(0.1, -0.2, 1.1, 0.1)
	l2 := sketch.AddLine(1.01, 0.2, 1.1, 0.9)
	l3 := sketch.AddLine(1.1, 1.2, 0.1, 1.1)
	l4 := sketch.AddLine(-0.1, 1.2, 0.1, 0.1)

	// Add constraints
	sketch.AddCoincidentConstraint(sketch.Origin, l1.Start())
	sketch.AddParallelConstraint(sketch.XAxis, l1)
	sketch.AddCoincidentConstraint(l2.Start(), l1.End())
	sketch.AddCoincidentConstraint(l3.Start(), l2.End())
	sketch.AddCoincidentConstraint(l4.Start(), l3.End())
	sketch.AddCoincidentConstraint(l1.Start(), l4.End())
	sketch.AddPerpendicularConstraint(l1, l2)
	sketch.AddParallelConstraint(l1, l3)
	sketch.AddDistanceConstraint(l1, nil, 1.0)
	sketch.AddDistanceConstraint(l2, nil, 1.0)
	sketch.AddDistanceConstraint(l3, nil, 1.0)

	// Solve
	err := sketch.Solve()

	// Output results
	if err != nil {
		fmt.Printf("Solve error %s\n", err)
	}

	// Export Image for solved sketch
	sketch.ExportImage("squareExample.svg")

	fmt.Println("Solved sketch: ")
	values := l1.Values(sketch)
	fmt.Printf("line 1: (%f, %f) to (%f, %f)\n", values[0], values[1], values[2], values[3])
	values = l2.Values(sketch)
	fmt.Printf("line 2: (%f, %f) to (%f, %f)\n", values[0], values[1], values[2], values[3])
	values = l3.Values(sketch)
	fmt.Printf("line 3: (%f, %f) to (%f, %f)\n", values[0], values[1], values[2], values[3])
	values = l4.Values(sketch)
	fmt.Printf("line 4: (%f, %f) to (%f, %f)\n", values[0], values[1], values[2], values[3])
```
> **dlineate** uses [zerolog](https://github.com/rs/zerolog) for logging. To change the logger or change logging level use utils.Logger.
