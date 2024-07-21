package test

import (
	"io"
	"net/http"
	"strings"
)

func ReadRespBody(resp http.Response) string {
	bytes, _ := io.ReadAll(resp.Body)
	return strings.TrimSuffix(string(bytes), "\n")
}
