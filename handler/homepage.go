package handler

import (
	"html/template"
	"log"
	"net/http"
)

var HomepageRoute = func(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		temp, _ := template.ParseFiles("route/template.html", "route/notfound.html")
		temp.Execute(res, nil)
    }
    route := NewRoute(res, req)
	route.Get(func() {
		temp, _ := template.ParseFiles("route/template.html", "route/page.html")
        
		temp.Execute(route.Response, nil)
	})
    type Test struct{
        Name string `json:"name"`
    }
    var postPayload Test = Test{
        Name: "",
    }
	route.Post(&postPayload, func() {
        log.Println(postPayload.Name)
		route.Response.Write([]byte("Hello World"))
	})
}
