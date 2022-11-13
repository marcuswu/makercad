%module occt

%typedef double Standard_Real;
%typedef void* Standard_Address;
%typedef void* uintptr_t;

// %typemap Standard_Real "float64"
// %insert(go_wrapper) %{
// type Standard_Real float64
// %}
%{
    #include <Standard_Macro.hxx>
    #include <Standard_ErrorHandler.hxx>
    #include <Standard.hxx>
    #include <Standard_DefineAlloc.hxx>
    #include <Standard_Handle.hxx>
    #include <Standard_Type.hxx>
    #include <Standard_Transient.hxx>
    using namespace opencascade;
%}

%include <Standard_Macro.hxx>
%include <Standard_DefineAlloc.hxx>
%include <Standard_Transient.hxx>
%include <Standard_Handle.hxx>
%include <Standard_Type.hxx>
%include <Standard_Real.hxx>
// %include <Standard.hxx>

%template(HandleStandardType) handle<Standard_Type>;
%template(HandleStandardTransient) handle<Standard_Transient>;