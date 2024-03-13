package handler

import (
	"encoding/json"
	"errors"
	"gohtmx/lib"
	"html/template"
	"io"
	"log"
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
    ContentType string
    Test string
}

func NewRoute(response http.ResponseWriter, request *http.Request) *Route {
	contentType := request.Header.Get("Content-Type")
	var urlEncoded map[string]string = map[string]string{}
	if contentType == "application/x-www-form-urlencoded" {
		urlEncoded = lib.ReadBody(request.Body)
	}	
    return &Route{
        ContentType: contentType,
		Request:    request,
        Test: "",
		Response:   response,
		Params:     map[string]string{},
		UrlEncoded: urlEncoded,
	}
}

func DecodeJson(jsonStruct any, body io.ReadCloser){
    decoder := json.NewDecoder(body)
    err := decoder.Decode(jsonStruct)
    if err != nil{
        log.Println("error decoding")
    }

}

func (r *Route) Post(jsonStruct any, handleFunc func()) {
	if r.Request.Method == http.MethodPost {
        if r.ContentType == "application/json" && jsonStruct != nil{
            DecodeJson(jsonStruct, r.Request.Body)  
        }
		handleFunc()
	}
}

func (r *Route) Patch(jsonStruct any, handleFunc func()) {
	if r.Request.Method == http.MethodPatch {
        if r.ContentType == "application/json" && jsonStruct != nil{
            DecodeJson(jsonStruct, r.Request.Body)  
        }

		handleFunc()
	}
}

func (r *Route) Delete(jsonStruct any, handleFunc func()) {
    if r.Request.Method == http.MethodDelete {
        if r.ContentType == "application/json" && jsonStruct != nil{
            DecodeJson(jsonStruct, r.Request.Body)  
        }

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
