package handler

import (
	"encoding/json"
	"gohtmx/lib"
	"html/template"
	"io"
	"log"
	"net/http"
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

func (r *Route) Render(writer http.ResponseWriter, payload interface{}, files ...string){
    temp, err := template.ParseFiles(files...)
    if err != nil{
        log.Println(err)
        return 
    }
    if err := temp.Execute(writer, payload); err != nil{
        log.Println(err)
        return
    }
}

