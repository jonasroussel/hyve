package servers

import (
	"net/http"
	"strings"

	"github.com/jonasroussel/proxbee/acme"
)

func HTTP01ChallengeSolver(handler *http.ServeMux) {
	handler.HandleFunc("GET /.well-known/acme-challenge/{token}", func(w http.ResponseWriter, r *http.Request) {
		token := r.PathValue("token")

		w.Header().Set("Content-Type", "text/plain")

		chal, exists := acme.HTTP01Provider.GetChallenge(token)
		if !exists {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("CHALLENGE_NOT_FOUND"))
			return
		}

		if !strings.HasPrefix(r.Host, chal.Domain) {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("DOMAIN_MISMATCH"))
			return
		}

		_, err := w.Write([]byte(chal.KeyAuth))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("INTERNAL_SERVER_ERROR"))
			return
		}
	})
}
