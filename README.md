# Matrix Image Manipulation
[<img alt="Build &amp; Test Go" src="https://github.com/MM-coder/matrix-image-manipulation/actions/workflows/go.yml/badge.svg"/>](https://github.com/MM-coder/matrix-image-manipulation/actions/workflows/go.yml)

> This project wishes to provide the programmatic implementation for a paper that provides the mathematical basis of common operations in computer vision. It provides a CLI, written in Golang, that allows for basic image manipulation operations.

# Running the project
## Installing
### Precompiled Binaries
Download the `.zip` corresponding to your operating system from the [Releases](https://github.com/MM-coder/matrix-image-manipulation/releases/latest) page and extract to a directory of choice.
Open the directory and run the executable as normal in your operating system.

### Building from source

These instructions are OS-agnostic, for OS-specific queries please refer to the Go documentation

If not installed, download [Go](https://go.dev/) and [Git](https://git-scm.com/), open a terminal window and run the following commands in the directory you would like to download the project to and run the following in order.

```bash
git clone https://github.com/MM-coder/matrix-image-manipulation.git
cd matrix-image-manipulation
go build
```
This will create the `matrix-image-manipulation` executable, in this case a `.exe`, to run it, you can call it in your terminal or double-click the executable

```bash
.\matrix-image-manipulation.exe
```

This will bring up the CLI.
 
## Usage Guide

The CLI is currently in Portuguese, it will first ask you for a file path, for this guide we will use the provided `gnome.png` file that's located in the repository root. **This file path must be a valid PNG file**.

```
Qual o path do ficheiro: gnome.png
```

You must next pick the operation you'd like to do on the image:

```
Escolha uma operação:
1: Filtro Gaussiano
2: Converter para Grayscale
3: Alterar Contraste
4: Alterar Luminosidade
Escolha (1, 2, 3 ou 4):
```
### Gaussian Filter

This will blur the image with a Gaussian filter with `σ = 2.5` and a kernel size of 7 and save the output to `{filename}_new.png`.

### Greyscale

This will convert the image to greyscale and save the output to `{filename}_new.png`.

### Contrast

This will ask you for the value of `m` and `b` in the equation `G(u) = mu+b` and save the output to `{filename}_new.png`.

```
Qual o path do ficheiro: gnome.png
Escolha uma operação:
1: Filtro Gaussiano
2: Converter para Grayscale
3: Alterar Contraste
4: Alterar Luminosidade
Escolha (1, 2, 3 ou 4): 3
3.1: Insira o valor de m: 1
3.2: Insira o valor de b: -1
```
TIP: If you insert `m=-1` and `b=1` you will invert the image

### Brightness

This will ask you for the value of `b` in the equation `G(u) = u+b` and save the output to `{filename}_new.png`.

```
Qual o path do ficheiro: gnome.png
Escolha uma operação:
1: Filtro Gaussiano
2: Converter para Grayscale
3: Alterar Contraste
4: Alterar Luminosidade
Escolha (1, 2, 3 ou 4): 4
4.1: Insira o valor de b:
```
TIP: You'll only notice a difference for larger values of `b` like `b=50`