## structexplorer

[![Go](https://github.com/emicklei/structexplorer/actions/workflows/go.yml/badge.svg)](https://github.com/emicklei/structexplorer/actions/workflows/go.yml)
[![GoDoc](https://pkg.go.dev/badge/github.com/emicklei/structexplorer)](https://pkg.go.dev/github.com/emicklei/structexplorer)
[![codecov](https://codecov.io/gh/emicklei/structexplorer/branch/master/graph/badge.svg)](https://codecov.io/gh/emicklei/structexplorer)

A Go Struct Explorer Service that offers a remote (HTTP) inspection of any Go struct.

## usage

    structexplorer.NewService("some structure", aStruct).Start()

then a HTTP service will be started

    INFO starting go struct explorer at http://localhost:5656

## syntax

- if a value is a pointer to a standard type then the display value has a "*" prefix
- if a value is a reflect.Value then the display value has a "~" prefix

## buttons

- ⇊ : explore one or more selected values from the list and put them below
- ⇉ : explore one or more selected values from the list and put them on the right
- z : show or hide fields which currently have zero value ("",0,nil,false)
- x : remove the struct from the page


## Example of exploring a [yaegi](https://github.com/traefik/yaegi) program

![program](./doc/explore_yaegi.png "Yaegi explore")