package lib

import (
	"bytes"
	"io"
	"net/url"
)

func ReadBody(body io.ReadCloser) map[string]string {
	var result map[string]string = map[string]string{}
	buf, _ := io.ReadAll(body)
	payload := bytes.Split(buf, []byte("&"))
	for _, value := range payload {
		keyValue := bytes.Split(value, []byte("="))
		key, _ := url.QueryUnescape(string(keyValue[0]))
		value, _ := url.QueryUnescape(string(keyValue[1]))
		result[key] = value
	}
	return result
}
