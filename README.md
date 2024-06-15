## structexplorer

A Go Struct Explorer Service that offers a remote (HTTP) inspection of any Go struct.

## usage

    structexplorer.NewService("some structure", aStruct).Start()

then a HTTP service will be started

    INFO starting go struct explorer at http://localhost:5656