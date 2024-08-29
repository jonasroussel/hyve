package tools

import (
	"encoding/json"
	"io"
	"strings"
)

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

func PredictWildcard(domain string) string {
	parts := strings.Split(domain, ".")
	return "*." + strings.Join(parts[1:], ".")
}
