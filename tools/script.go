package tools

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/dop251/goja"
)

var v8vm *goja.Runtime
var proxyCallback func(string, string, string, map[string]string, map[string]string) string

func loadScript() (string, error) {
	bytes, err := os.ReadFile(Env.DYNAMIC_TARGET)
	if err != nil {
		return "", err
	}

	script := string(bytes)
	script = regexp.MustCompile(`export\s+(function|const)`).ReplaceAllString(script, "$1")

	return script, nil
}

func LoadDynamicTarget() {
	if Env.DYNAMIC_TARGET == "" {
		return
	}

	script, err := loadScript()
	if err != nil {
		log.Fatal(err)
	}

	v8vm = goja.New()

	_, err = v8vm.RunString(script)
	if err != nil {
		log.Fatal(err)
	}

	err = v8vm.ExportTo(v8vm.Get("onProxy"), &proxyCallback)
	if err != nil {
		log.Fatal(errors.New("unable to detect `onProxy` function (" + Env.DYNAMIC_TARGET + ")"))
	}
}

func CallDynamicTarget(req *http.Request) string {
	host := req.Host

	path := req.URL.EscapedPath()
	if path != "/" {
		path = strings.TrimRight(path, "/")
	}

	url := fmt.Sprintf("https://%s%s", host, path)

	query := map[string]string{}
	for k, v := range req.URL.Query() {
		query[k] = v[0]
	}

	headers := map[string]string{}
	for k, v := range req.Header {
		headers[k] = v[0]
	}

	return proxyCallback(url, host, path, query, headers)
}
