package core

type Entity interface {
	isConstruction() bool
	setConstruction(bool)
}
