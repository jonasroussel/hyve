package servers

import "net/http"

func RedirectToHTTPS(handler *http.ServeMux) {
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://"+r.Host+r.URL.Path, http.StatusMovedPermanently)
	})
}
