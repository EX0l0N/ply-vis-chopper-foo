# ply-vis-chopper

### ply-vis-chopper may help you updating a "vis" file created by [colmap](https://demuc.de/colmap/)

This is a very specialised tool which is only interesting for people who play around with colmap.  
If you don't know what colmap is made for, or you don't use it – then you probably won't need this tool and you are wasting your time by reading on.

## Introduction

Colmap basically operates in two steps.

In the first step camera positions are extrapolated from the set of input images. Very few points are generated in this step.
In the second step the camera positions from the first step are used in combination with the input images to generate a rich point cloud of the model.

This second step always generates a file called `fused.ply` and another file named `fused.ply.vis`.

While the first file is a standard PLY file that you should be able to open and edit in almost every tool which is capable of loading and maybe editing vertex mesh based 3D models, the second file is something you will probably not be able to open.

**At least for me this was a problem, which I tried to solve by writing this small helper.**

Colmap comes with two meshers, that try do extrapolate a 3D mesh out of the point-cloud that is created by the second step (densification).

Both of these meshers take the `fused.ply` file as input.
Another tool called [OpenMSV](https://github.com/cdcseacave/openMVS) is also able to work with these files.

A common problem at this stage is, that colmap has no idea if the points it was able to create belong to you model, or not.  
Makeing the meshers work on that data will always create geometry that you don't want in your final result.

That alone won't be that problematic, if only model and background would be separable in any way.
Unfortunately the background noise does even have an effect on your model, as it alters surrounding normals and thus bends your models surfaces.

You may now come to the conclusion that it would be a good idea to remove all background noise from the `fused.ply` file using another tool (like [MeshLab](https://www.meshlab.net/)) before running the meshers on it (as I did).

This will only succeed if you use the internal **poisson** mesher.
You may as well use an external mesher like the reconstruction filters of MeshLab.

But - and this is where **ply-vis-chopper** finally comes into play - if you use the internal colmap delaunay mesher, or try to process the data with OpenMVS, you won't be able to edit the data before using the mesher (or OpenMVS) because the process will fail.

**This is due to a dependency between that mysterious `fused.ply.vis` and the `fused.ply` file.**

The `vis` file simply keeps a list of all the images a point in the `fused.ply` can be seen from.  
The image lists are saved in the exact order the points appear in the `ply` file, so deleting any point in the `ply` will get those files out of sync and the `vis` file becomes useless.

This is what **ply-vis-chopper** can fix … at least in some cases.

What it does is comparing the points, their normals and their colors from your edited `fused.ply` with the original one (that means you need to keep the old one - don't overwrite it) and repositions the image lists in the output `vis` file.

**That means you may only delete points (or change their order), or the matching process will fail**

The points need to preserve their location, their color and their normals to be matched.  
You may not scale, translate or move your model!

## How to run

ply-vis-chopper is a command line binary,  
but you don't need to compile anything.

Instead just download a binary suitable for you operating system from the `/bin` folder of this project.

### There are command line binaries for:

- [macos-amd64](../../raw/main/bin/macos-amd64/ply-vis-chopper)
- [linux-amd64](../../raw/main/bin/linux-amd64/ply-vis-chopper)
- [linux-arm](../../raw/main/bin/linux-arm/ply-vis-chopper)
- [windows-amd64](../../raw/main/bin/windows-amd64/ply-vis-chopper.exe)

The binary takes exactly four parameters.
Failure to provide these parameters will make the tool bark out debug info to your terminal at the moment.

Here's an example:

```
# generic form:
./ply-vis-chopper <ORIGINAL-FUSED.PLY> <EDITED_FUSED.PLY> <ORIGINAL-FUSED.PLY.VIS> <OUPUT.PLY.VIS>

# example invocation:
./ply-vis-chopper fused.ply cleaned.ply fused.ply.vis cleaned.ply.vis
```

The last parameter obviously references a file that is being written to and will be created if needed.  
If it already exists, it will be overwritten.

(Don't use a leading `./` on Windows.)

Please note that you will need to move the original `fused` files out of the way and rename the edited ones to `fused` afterwards.

### And there is even a rudimentary GUI:

![](../../raw/main/img/gui.png)

- [macos-amd64](../../raw/main/bin/macos-amd64/ply-vis-chopper-gui)
- [linux-amd64](../../raw/main/bin/linux-amd64/ply-vis-chopper-gui)
- [windows-amd64](../../raw/main/bin/windows-amd64/ply-vis-chopper-gui.exe)

## But I _want_ to compile

That's ok.  
The `ply-vis-chopper` has no dependencies.

( In case of the GUI version you will need to

- fetch [fyne](https://fyne.io/): `go get fyne.io/fyne`.  
- have a `C` compiler installed for your OS to glue the GL bindings )

Besides that, you need to make your `$GOPATH` point to the projects folder and then build.

This will be

```
export GOPATH=$PWD
go build ply-vis-chopper(-gui).go
```
in most cases.
