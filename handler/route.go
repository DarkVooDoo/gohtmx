package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"gohtmx/lib"
	"html/template"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
)

type RouteInterface interface {
	Get()
	Post()
	Delete()
	Patch()
}

type Multipart struct{
    Body map[string]string
    File io.Reader
}

type ErrToken error

type Route struct {
	RouteInterface
    Multipart Multipart
	Params     map[string]string
	Request    *http.Request
	Response   http.ResponseWriter
	UrlEncoded map[string]string
    ContentType string
}

func NewRoute(response http.ResponseWriter, request *http.Request) (*Route, ErrToken) {

    var multip Multipart = Multipart{Body: map[string]string{}}
	contentType := request.Header.Get("Content-Type")
    token, err := request.Cookie("x-auth")
    err = VerifyToken(token.String())
    ct, params, _ := mime.ParseMediaType(request.Header.Get("Content-Type"))
	var urlEncoded map[string]string = map[string]string{}
	if contentType == "application/x-www-form-urlencoded" {
		urlEncoded = lib.ReadBody(request.Body)
    }else if ct == "multipart/form-data"{
        reader := multipart.NewReader(request.Body, params["boundary"])
        body := make([]byte, 1024)
        var payload string
        Parts:
        for{
            part, err := reader.NextPart()
            if err == io.EOF{
                break Parts
            }
            defer part.Close()
            var fileBuffer bytes.Buffer
            Exit:
            for {
                bytesReaded, err := part.Read(body)
                fileBuffer.Write(body[:bytesReaded])
                payload = string(body[:bytesReaded])
                if err == io.EOF {break Exit}
            }
            
            if part.FileName() != ""{
                multip.File = &fileBuffer
                continue
            }
            multip.Body[part.FormName()] = payload
            
        }
    }	
    return &Route{
        ContentType: contentType,
        Multipart: multip,
		Request:    request,
		Response:   response,
		Params:     map[string]string{},
		UrlEncoded: urlEncoded,
	}, err
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

func VerifyToken(token string)error{
    if token == ""{
        return errors.New("unauthorized")
    }
    return nil
}
