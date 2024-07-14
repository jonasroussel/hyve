package tools

import (
	"encoding/json"
	"io"
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
