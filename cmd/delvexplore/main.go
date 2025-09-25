package main

import (
	"fmt"
	"log/slog"

	"github.com/go-delve/delve/service/rpc2"
)

// https://github.com/go-delve/delve/blob/master/service/rpc2/client.go
// https://github.com/aarzilli/delve_client_testing
func main() {
	client := rpc2.NewClient("127.0.0.1:51551")
	slog.Info("connected")
	if !client.AttachedToExistingProcess() {
		slog.Error("cannot attach")
	}

	state, err := client.GetState()
	if err != nil {
		slog.Error("unable to get state", "err", err)
		return
	}
	fmt.Printf("%#v", state)
}
