package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    //"os"
    "strings"
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

func fileResponseCreator(folder string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
			var p *Page
			var err error
			if len(r.URL.Path) == 1 {
				p, err = loadPage("templates", "index.html")
			} else {
				p, err = loadPage(folder, r.URL.Path[1:])
			}
			if p != nil {
    			w.Write(p.Body)
    		} else {
    			fmt.Println(err)
    		}
		}
}

func dataSentHandler(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    urlVars := strings.Split(r.URL.Path[1:], "/")
    cpu_id, node_id := urlVars[1], urlVars[2]
    fmt.Println(r.Form, cpu_id, node_id)
}

func displayHandler() {
    staticHandler := fileResponseCreator("static")
    http.HandleFunc("/", fileResponseCreator("templates"))
    http.HandleFunc("/css/", staticHandler)
    http.HandleFunc("/js/", staticHandler)
    http.HandleFunc("/img/", staticHandler)
    http.HandleFunc("/favicon.ico", fileResponseCreator("static/img"))
}

func main() {
    displayHandler()
    http.HandleFunc("/checked/", dataSentHandler)
    http.ListenAndServe(":8080", nil)
}
