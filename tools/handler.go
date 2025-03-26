package tools

import (
	"net/http"
	"net/url"
	pathutil "path"
	"strings"
)

type HTTPHandler struct {
	Routes   map[string]http.Handler
	Fallback http.HandlerFunc
}

func NewHTTPHandler() *HTTPHandler {
	return &HTTPHandler{
		Routes: make(map[string]http.Handler),
	}
}

func (handler *HTTPHandler) Handle(path string, handle http.Handler) {
	path = forceRootedPath(path)
	path = pathutil.Clean(path)

	handler.Routes[path] = handle
}

func (handler *HTTPHandler) HandleFallback(handle func(w http.ResponseWriter, r *http.Request)) {
	handler.Fallback = handle
}

func (handler *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "*" {
		if r.ProtoAtLeast(1, 1) {
			w.Header().Set("Connection", "close")
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	path := r.URL.EscapedPath()
	path = forceRootedPath(path)
	path = pathutil.Clean(path)

	for i := 1; i <= len(path); {
		prefix := path[:i]

		if h, ok := handler.Routes[prefix]; ok {
			stripPrefix(prefix, h).ServeHTTP(w, r)
			return
		}

		if i >= len(path) {
			break
		}

		var moveBy int
		if path[i] == '/' {
			moveBy = strings.Index(path[i+1:], "/") + 1
		} else {
			moveBy = strings.Index(path[i:], "/")
		}

		if moveBy < 1 {
			i = len(path)
		} else {
			i += moveBy
		}
	}

	if handler.Fallback != nil {
		handler.Fallback(w, r)
		return
	}

	http.NotFound(w, r)
}

func stripPrefix(prefix string, handler http.Handler) http.Handler {
	if prefix == "" {
		return handler
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := pathutil.Clean(forceRootedPath(strings.TrimPrefix(r.URL.Path, prefix)))
		rp := strings.TrimPrefix(r.URL.RawPath, prefix)

		if r.URL.RawPath != "" {
			rp = pathutil.Clean(forceRootedPath(rp))
		}

		r2 := new(http.Request)
		*r2 = *r
		r2.URL = new(url.URL)
		*r2.URL = *r.URL
		r2.URL.Path = p
		r2.URL.RawPath = rp

		handler.ServeHTTP(w, r2)
	})
}

func forceRootedPath(path string) string {
	if len(path) == 0 {
		return "/"
	}
	if path[0] != '/' {
		return "/" + path
	}
	return path
}
