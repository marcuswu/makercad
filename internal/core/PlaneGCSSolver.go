package core

import "libmakercad/third_party/planegcs"

/*
   PlaneGCSSolver needs to track fixed geometry vs unfixed geometry so we can declare a list of unknown geometry to GCS
	GCSsys.invalidatedDiagnosis();
	GCSsys.declareUnknowns(Parameters);
	GCSsys.declareDrivenParams(DrivenParameters);
    GCSsys.initSolution(defaultSolverRedundant);

	For diagnosis:
	GCSsys.getConflicting(Conflicting);
	GCSsys.getRedundant(Redundant);
	GCSsys.getPartiallyRedundant (PartiallyRedundant);
	GCSsys.getDependentParams(pDependentParametersList);
	GCSsys.dofsNumber();

	Solve:
	GCSsys.solve(isFine, GCS::DogLeg or BFGS or LevenbergMarquardt)
	GCSsys.applySolution();

	If not solved, fall back to other solvers...
*/
type PlaneGCSSolver struct {
	system planegcs.System
}

func NewPlaneGCSSolver() *PlaneGCSSolver {
	return &PlaneGCSSolver{planegcs.NewSystem()}
}

func (s *PlaneGCSSolver) CreatePoint(x float64, y float64) *Point {
	return NewPoint(x, y)
}
