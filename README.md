_This project has been created as part of the 42 curriculum by lalegran_

# Scop

## Description

The goal is to build a small application to display a 3d object stored in a `.obj`.

## Instructions

This project need `go@1.26.10`, it can be isntalled with [mise](https://mise.jdx.dev/).

(optionnal) Install dependency:

```shell
mise i
```

Compile with:

```shell
make
```

Run:

```shell
./scop [file.obj]
```

### Commands

Move model: `W`, `S`, `A`, `D`, `C`, `V`
Rotate model: `↑`,`↓`, `→`,`←`,`Z`,`X`
Toggle texture: `T`
Enable auto-rotate: `1`, `2`, `3`
Stop auto-rotate: `[space]`
Reset model position: `R`

## Resources

https://opengl.developpez.com/tutoriels/apprendre-opengl/
https://www.songho.ca/opengl/gl_projectionmatrix.html
https://antongerdelan.net/opengl/vertexbuffers.html
https://webglfundamentals.org/webgl/lessons/webgl-planar-projection-mapping.html
