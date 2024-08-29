package servers

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/jonasroussel/hyve/tools"
)

var proxy = &httputil.ReverseProxy{
	Rewrite: func(r *httputil.ProxyRequest) {
		var target *url.URL
		var err error

		if tools.Env.DYNAMIC_TARGET != "" {
			target, err = url.Parse(tools.CallDynamicTarget(r.In))
		} else {
			target, err = url.Parse(tools.Env.Target)
		}

		if err != nil {
			log.Printf("[WARNING] Failed to parse target URL: %s", err)
			return
		}

		r.SetURL(target)
		r.SetXForwarded()
		r.Out.Host = target.Host
	},
}

func ReverseProxy(handler *http.ServeMux) {
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.TLS.ServerName != "" && r.TLS.ServerName == tools.Env.AdminDomain {
			http.NotFound(w, r)
			return
		}

		proxy.ServeHTTP(w, r)
	})
}
