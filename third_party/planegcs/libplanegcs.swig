%module(directors="1") planegcs

%{#include <GCS.h>%}

%feature("director") System;
%feature("director") Point;
%feature("director") Line;
%feature("director") Circle;
%feature("director") Arc;

%include <GCS.h>
%include <Geo.h>