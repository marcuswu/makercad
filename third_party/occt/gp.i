%module(directors="1") occt
%{#include <gp_Trsf.hxx>%}

// // %include <typemaps.i>
// // %include "std_string.i"
// // %include "std_vector.i"
%feature("director") gp_Pnt;

%include <Standard_Macro.hxx>
%include <gp_Pnt.hxx>
%include <gp_Trsf.hxx>
%include <gp_Vec.hxx>
%include <gp_Dir.hxx>
%include <gp_Ax1.hxx>
%include <gp_Ax3.hxx>