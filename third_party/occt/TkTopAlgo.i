%module(directors="1") occt
%{#include <Standard_Transient.hxx>%}
%{#include <Geom_Geometry.hxx>%}
%{#include <Geom_Curve.hxx>%}
%{#include <BRepBuilderAPI_MakeEdge.hxx>%}

%feature("director") Standard_Transient;
%feature("director") Geom_Geometry;
%feature("director") Geom_Curve;

%include <Standard_Transient.hxx>
%include <Geom_Geometry.hxx>
%include <Geom_Curve.hxx>
%include <BRepBuilderAPI_MakeEdge.hxx>
