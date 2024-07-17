package servers

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/jonasroussel/hyve/tools"
)

func ReverseProxy(handler *http.ServeMux) {
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.TLS.ServerName != "" && r.TLS.ServerName == tools.Env.AdminDomain {
			http.NotFound(w, r)
			return
		}

		proxy(w, r)
	})
}

func proxy(w http.ResponseWriter, r *http.Request) {
	targetURL, _ := url.Parse(tools.Env.Target)

	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(targetURL)
			r.SetXForwarded()
			r.Out.Host = targetURL.Host
		},
	}

	proxy.ServeHTTP(w, r)
}
