package servers

import (
	"net/http"
	"strings"

	"github.com/jonasroussel/hyve/acme"
	"github.com/jonasroussel/hyve/stores"
	"github.com/jonasroussel/hyve/tools"
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

		if !verifyAdminKey(r) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var data BodyData
		err := tools.ParseBody(r.Body, &data)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if data.Domain == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !tools.IsDNSValid(data.Domain) {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("DOMAIN_NOT_CONFIGURED"))
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

	handler.HandleFunc("POST /api/remove", func(w http.ResponseWriter, r *http.Request) {
		if r.TLS.ServerName != tools.Env.AdminDomain {
			proxy(w, r)
			return
		}

		if !verifyAdminKey(r) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var data BodyData
		err := tools.ParseBody(r.Body, &data)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if data.Domain == "" {
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

func verifyAdminKey(r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return false
	}

	parts := strings.Split(strings.Trim(authHeader, " "), " ")
	if len(parts) != 2 {
		return false
	}

	if parts[0] != "Bearer" {
		return false
	}

	return parts[1] == tools.Env.AdminKey
}
