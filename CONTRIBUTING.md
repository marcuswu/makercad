# Contributing #

## Adding support for new features, modeling functionality or otherwise improving MakerCAD.
Thank you for your interest in improving MakerCAD.

We appreciate all help making this project better, so contributions are welcome. Below, you will find how to contribute in different ways. If none of those matches your situation, please reach out (see below).

### New to MakerCAD
We would love your feedback on getting started with MakerCAD. Please let us know of any difficulties, confusion, etc so we can make it easier for future users. Please [open a Github issue](https://github.com/marcuswu/makercad/issues) with your questions or your can also get in touch with us directly on the [MakerCAD Discord Server](https://discord.gg/9EwXdfuJhw).

### Something in MakerCAD does not work as expected
Please [open a Github issue](https://github.com/marcuswu/makercad/issues) describing your problem and we will be happy to assist.

### Something in that you want/need does not appear to be in MakerCAD
We may not have implemented it yet. Please [open a Github issue](https://github.com/marcuswu/makercad/issues) to ensure there is no duplication of effort. A pull request adding the functionality to MakerCAD would be greatly appreciated.

### How to use the Github repository
The main branch has the latest version of MakerCAD. Versions will be tagged using semantic versioning. Active development for the next release will take place in the dev branch.

### Here is how to contribute back code or documentation:
* Fork the repository
* Create a feature branch off of the dev branch
* Make your changes
* Format your changes using gofmt
* Make sure tests still pass
* Submit a pull request against the dev branch
* Be kind
* Please rebase instead of merge from the dev branch if your PR needs to include changes that occured after your feature branch was created. This can be accomplished via the git command line with:
```
git checkout dev
git pull --rebase origin dev
git checkout my-feature-branch
git rebase dev
git push myfork my-feature-branch -f
```
