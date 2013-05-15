package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    //"strings"
)

type Page struct {
	Title string
	Body []byte
}

func loadPage(folder, title string) (*Page, error) {
    filename := folder + "/" + title
    body, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Page{Title: title, Body: body}, nil
}

func fileResponseHandler(folder string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
			fmt.Println(r.ParseForm())
			var p *Page
			var err error
			if len(r.URL.Path) == 1 {
				p, err = loadPage("templates", "index.html")
			} else {
				p, err = loadPage(folder, r.URL.Path[1:])
			}
			if p != nil {
    			fmt.Fprintf(w, string(p.Body))
    		} else {
    			fmt.Println(err)
    		}
		}
}

func main() {
	staticHandler := fileResponseHandler("static")
    http.HandleFunc("/", fileResponseHandler("templates"))
    http.HandleFunc("/css/", staticHandler)
    http.HandleFunc("/js/", staticHandler)
    http.HandleFunc("/img/", staticHandler)
    http.HandleFunc("/favicon.ico", fileResponseHandler("static/img"))
    http.ListenAndServe(":8080", nil)
}
