# libMakerCAD

cd third_party/occt
swig -I/usr/local/include/opencascade -go -c++ -intgosize 64 libocct.swig
cd ../planegcs
swig -I/usr/local/include/planegcs -go -c++ -intgosize 64 libplanegcs.swig
