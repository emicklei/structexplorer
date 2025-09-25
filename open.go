package structexplorer

import (
	"fmt"
	"os/exec"
	"runtime"
)

// Open calls the OS default program for uri
func open(uri string) error {
	switch {
	case "windows" == runtime.GOOS:
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", uri).Start()
	case "darwin" == runtime.GOOS:
		return exec.Command("open", uri).Start()
	case "linux" == runtime.GOOS:
		return exec.Command("xdg-open", uri).Start()
	default:
		return fmt.Errorf("unable to open uri:%v on:%v", uri, runtime.GOOS)
	}
}
