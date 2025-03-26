package tools

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type JSON map[string]any

func ParseBody(body io.ReadCloser, out any) error {
	bytes, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	defer body.Close()

	err = json.Unmarshal(bytes, out)
	if err != nil {
		return err
	}

	return nil
}

func DetailedError(w http.ResponseWriter, status int, err string, details ...any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	var jsonError JSON = JSON{"error": err}
	if len(details) > 0 {
		switch details[0].(type) {
		case string:
			jsonError["message"] = details[0]
			if len(details) > 1 {
				jsonError["details"] = details[1:]
			}
		default:
			jsonError["details"] = details
		}
	}

	json.NewEncoder(w).Encode(jsonError)
}

func JSONResp(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func PredictWildcard(domain string) string {
	parts := strings.Split(domain, ".")
	return "*." + strings.Join(parts[1:], ".")
}
