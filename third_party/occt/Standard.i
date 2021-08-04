%module occt

%typedef double Standard_Real;
%typedef void* Standard_Address;

// %typemap Standard_Real "float64"
// %insert(go_wrapper) %{
// type Standard_Real float64
// %}
%{#include <Standard_ErrorHandler.hxx>%}
%{#include <Standard.hxx>%}
%{#include <Standard_DefineAlloc.hxx>%}
%include <Standard_DefineAlloc.hxx>
// %include <Standard.hxx>