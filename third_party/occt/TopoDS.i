%module(directors="1") occt
%{#include <TopoDS_Edge.hxx>%}

// // %include <typemaps.i>
// // %include "std_string.i"
// // %include "std_vector.i"
%feature("director") TopoDS_Edge;

%include <TopoDS_Edge.hxx>