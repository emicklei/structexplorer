package structexplorer

import (
	"net/http"
	"path"
)

// Options can be used to configure a Service on startup.
type Options struct {
	// Uses 5656 as the default
	HTTPPort int
	// Uses http.DefaultServeMux as default
	ServeMux *http.ServeMux
	// Uses "/" as default
	HTTPBasePath string
}

func (o *Options) rootPath() string {
	if o.HTTPBasePath == "" {
		return "/"
	}
	return path.Join("/", o.HTTPBasePath)
}

func (o *Options) httpPort() int {
	if o.HTTPPort == 0 {
		return 5656
	}
	return o.HTTPPort
}

func (o *Options) serveMux() *http.ServeMux {
	if o.ServeMux == nil {
		return http.DefaultServeMux
	}
	return o.ServeMux
}
