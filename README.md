## structexplorer

[![Go](https://github.com/emicklei/structexplorer/actions/workflows/go.yml/badge.svg)](https://github.com/emicklei/structexplorer/actions/workflows/go.yml)
[![GoDoc](https://pkg.go.dev/badge/github.com/emicklei/structexplorer)](https://pkg.go.dev/github.com/emicklei/structexplorer)
[![codecov](https://codecov.io/gh/emicklei/structexplorer/branch/main/graph/badge.svg)](https://codecov.io/gh/emicklei/structexplorer)

A Go Struct Explorer Service (http.Handler) that offers remote inspection of any Go struct and its references.

### example of exploring a [yaegi](https://github.com/traefik/yaegi) program

![program](./doc/explore_yaegi.png "Yaegi explore")

## install

    go get github.com/emicklei/structexplorer

## usage

    structexplorer.NewService("some structure", yourStruct).Start()

or as HTTP Handler:

    s := structexplorer.NewService("some structure", yourStruct)
    http.ListenAndServe(":5656", s)

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

Note: if the list contains just one structural value then selecting it can be skipped for both ⇊ and ⇉.

## explore while debugging

Currently, the standard Go debugger `delve` stops all goroutines while in a debugging session.
This means that if you have started the `structexplorer` service in your program, it will not respond to any HTTP requests during that session.

The explorer can also be asked to dump an HTML page with the current state of values to a file.

    s := structexplorer.NewService()
    s.Explore("some structure", yourStruct, "some field", yourStruct.Field).Dump()

Another method is to use a special test case which starts and explorer at the end of a test and then run it with a longer acceptable timeout.

## examples

See folder `examples` for simple programs demonstrating each feature.
