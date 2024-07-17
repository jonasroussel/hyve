package servers

import (
	"net/http"

	"github.com/jonasroussel/proxbee/acme"
	"github.com/jonasroussel/proxbee/stores"
	"github.com/jonasroussel/proxbee/tools"
)

type BodyData struct {
	Domain string `json:"domain"`
}

func AdminAPI(handler *http.ServeMux) {
	if tools.Env.AdminDomain == "" || tools.Env.AdminKey == "" {
		return
	}

	handler.HandleFunc("POST /api/add", func(w http.ResponseWriter, r *http.Request) {
		if r.TLS.ServerName != tools.Env.AdminDomain {
			proxy(w, r)
			return
		}

		var data BodyData
		err := tools.ParseBody(r.Body, &data)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		cert, err := stores.Active.GetCertificate(data.Domain)
		if err != nil && err != stores.ErrNotFound {
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else if cert != nil {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("DOMAIN_ALREADY_REGISTERED"))
			return
		}

		err = acme.RegisterDomain(data.Domain)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write([]byte("OK"))
	})

	handler.HandleFunc("POST /api/renew", func(w http.ResponseWriter, r *http.Request) {
		if r.TLS.ServerName != tools.Env.AdminDomain {
			proxy(w, r)
			return
		}

		// TODO

		w.Write([]byte("TODO"))
	})

	handler.HandleFunc("POST /api/remove", func(w http.ResponseWriter, r *http.Request) {
		if r.TLS.ServerName != tools.Env.AdminDomain {
			proxy(w, r)
			return
		}

		var data BodyData
		err := tools.ParseBody(r.Body, &data)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = stores.Active.RemoveCertificate(data.Domain)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write([]byte("OK"))
	})
}
