# MakerCAD Roadmap

This is a very high level list of the things I want to do with MakerCAD. If you have something specific that you would like to be a part of MakerCAD, you can take one of these options:

* Join us on [discord](https://discord.gg/9EwXdfuJhw) and suggest it there
* Open an [issue on GitHub](https://github.com/marcuswu/makercad/issues)

As a fairly early project, much of the actual work is unknown, but there are categories of work I know need to be handled. When I do separate smaller segments of work, those will be tracked elsewhere. This is for large concerns and a broad view of where the project(s) are headed.

## Constraint Solver Consistency
Improve the consistency of solving sketch constraints. 

- [x] Add a numeric solver to handle cases the graph based approach does not merge into a final solution

## MakerCAD API Improvements & Consistency
I am aware there are some things that need work with the API. There are some consistency issues, some things that should be easier to do, and functionality that is missing.

## UI Development
Code based CAD is great, but it is not for everyone. To reach more people, I am developing UIs to partner with MakerCAD.

- [x] VSCode plugin - Provide an easy way to write code and see the results a la OpenSCAD. This is published, but could use improvement and more features.
- [ ] Interactive UI - Provide a more traditional point and click interface for MakerCAD. The UI should generate code with annotated metadata. When reloading an existing project, use the annotations to rebuild UI steps. Manual editing of the code is unsupported and checksums should be used to provide a UI warning if manually edited code is loaded. Loading of a manually created project is unsupported initially, but some limited support for rendering the generated model may be in the cards.