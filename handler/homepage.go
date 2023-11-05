package handler

import (
	"html/template"
	"net/http"
)

var HomepageRoute = func(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		temp, _ := template.ParseFiles("route/template.html", "route/notfound.html")
		temp.Execute(res, nil)
		return
	}
	route := NewRoute(res, req)
	route.Get(func() {
		temp, _ := template.ParseFiles("route/template.html", "route/page.html")
		temp.Execute(route.Response, nil)
	})

	route.Post(func() {
		route.Response.Write([]byte(route.UrlEncoded["name"]))
	})
}
