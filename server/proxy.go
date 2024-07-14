package server

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/jonasroussel/proxbee/config"
)

var TARGET_URL, _ = url.Parse(config.TARGET)

func ForwardProxy(handler *http.ServeMux) {
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Host == config.ADMIN_DOMAIN {
			http.NotFound(w, r)
			return
		}

		proxy(w, r)
	})
}

func proxy(w http.ResponseWriter, r *http.Request) {
	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(TARGET_URL)
			r.SetXForwarded()
			r.Out.Host = TARGET_URL.Host
		},
	}

	proxy.ServeHTTP(w, r)
}
