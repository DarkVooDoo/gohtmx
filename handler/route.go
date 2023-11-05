package handler

import (
	"errors"
	"gohtmx/lib"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
)

type RouteInterface interface {
	Get()
	Post()
	Delete()
	Patch()
}

type Route struct {
	RouteInterface
	Params     map[string]string
	Request    *http.Request
	Response   http.ResponseWriter
	UrlEncoded map[string]string
}

func NewRoute(response http.ResponseWriter, request *http.Request) *Route {
	contentType := request.Header.Get("Content-Type")
	var urlEncoded map[string]string = map[string]string{}
	if contentType == "application/x-www-form-urlencoded" {
		urlEncoded = lib.ReadBody(request.Body)
	}
	return &Route{
		Request:    request,
		Response:   response,
		Params:     map[string]string{},
		UrlEncoded: urlEncoded,
	}
}

func (r *Route) Post(handleFunc func()) {
	if r.Request.Method == http.MethodPost {
		handleFunc()
	}
}

func (r *Route) Patch(handleFunc func()) {
	if r.Request.Method == http.MethodPatch {
		handleFunc()
	}
}

func (r *Route) Delete(handleFunc func()) {
	if r.Request.Method == http.MethodDelete {
		handleFunc()
	}
}

func (r *Route) Get(handlerFunc func()) {
	if r.Request.Method == http.MethodGet {
		handlerFunc()
	}
}

func (r *Route) GetParams(pattern string, keys ...string) error {
	isValid, _ := filepath.Match(pattern, r.Request.URL.Path)
	if !isValid {
		temp, _ := template.ParseFiles("route/nofound.html")
		temp.Execute(r.Response, nil)
		return errors.New("error")
	}
	for _, key := range keys {
		r.Params[key] = strings.Replace(r.Request.URL.Path, pattern[:len(pattern)-1], "", 1)
	}
	return nil
}
