## structexplorer

A Go Struct Explorer Service that offers a remote (HTTP) inspection of any Go struct.

## usage

    structexplorer.NewService("some structure", aStruct).Start()

then a HTTP service will be started

    INFO starting go struct explorer at http://localhost:5656

## syntax

- if a value is a pointer to a standard type then the display value has a "*" prefix
- if a value is a reflect.Value then the display value has a "~" prefix

## exploring a [yaegi](https://github.com/traefik/yaegi) program

![program](./doc/explore_yaegi.png "Yaegi explore")