package ui

import (
	"embed"
	"io/fs"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/jonasroussel/hyve/tools"
)

//go:embed all:dist
var dist embed.FS

var DistFS, _ = fs.Sub(dist, "dist")

// ServeFiles binds the UI files to the given handler.
// If the HYVE_UI_PROXY_PORT env variable is set, it will serve the files through a reverse proxy
// to the Vite server running on port the port specified in the env variable.
func ServeFiles(handler *tools.HTTPHandler) {
	if strings.HasPrefix(os.Args[0], os.TempDir()) {
		handler.HandleFallback(func(w http.ResponseWriter, r *http.Request) {
			proxy := &httputil.ReverseProxy{
				Rewrite: func(r *httputil.ProxyRequest) {
					target, err := url.Parse("http://localhost:8081")
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}

					r.SetURL(target)
					r.SetXForwarded()
					r.Out.Host = target.Host
				},
			}

			proxy.ServeHTTP(w, r)
		})
	} else {
		handler.HandleFallback(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.EscapedPath()

			// Remove the leading slash
			if path[0] == '/' {
				path = path[1:]
			}

			if _, err := DistFS.Open(path); err != nil {
				path = "index.html"
			}

			http.ServeFileFS(w, r, DistFS, path)
		})
	}
}
