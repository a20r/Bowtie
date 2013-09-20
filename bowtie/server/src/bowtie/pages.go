package main

import (
    // sys pkgs
    "flag"

    // io pkgs
    "fmt"
    "io/ioutil"

    "code.google.com/p/go.net/websocket"

    // network pkgs
    "net/http"

    // string pkgs
    "encoding/json"

    //ADTs
    "time"
)

// Represents an file loaded
type Page struct {
    Title string
    Body []byte
}

// JSON response mapping
type Response map[string]interface{}

// Converts the JSON to strings
// to be sent as a response
func (r Response) String() (s string) {
    b, err := json.Marshal(r)
    if err != nil {
        s = ""
        return
    }
    s = string(b)
    return
}

// Opens a file and returns it represented
// as a Page.
func loadPage(folder, title string) (*Page, error) {
    filename := folder + "/" + title
    body, err := ioutil.ReadFile(filename)

    if err != nil {
        return nil, err
    }

    return &Page{Title: title, Body: body}, nil
}

func timePrinter(message string) {
    fmt.Println(message + "\t: " + time.Now().String())
}

// Creates a function that will be used as a handler
// for static and template responses. See Usage!
func fileResponseCreator(folder string) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        var p *Page
        var err error

        timePrinter("GET\t" + r.URL.Path)

        if len(r.URL.Path) == 1 {
            // In case the path is just '/'
            p, err = loadPage("templates", "index.html")
        } else {
            p, err = loadPage(folder, r.URL.Path[1:])
        }

        if p != nil {
            w.Write(p.Body)
        } else {
            timePrinter("ERROR\t" + err.Error())
        }
    }
}

// Handles all in coming http requests
func UIHandler() {
    staticHandler := fileResponseCreator("static")
    http.HandleFunc("/", fileResponseCreator("templates"))
    http.HandleFunc("/css/", staticHandler)
    http.HandleFunc("/js/", staticHandler)
    http.HandleFunc("/img/", staticHandler)
    http.HandleFunc("/client/", fileResponseCreator("../"))
    http.HandleFunc("/examples/", fileResponseCreator("../../"))
    http.HandleFunc("/favicon.ico", fileResponseCreator("static/img"))
}

// MAIN EXECUTION FLOW
func main() {

    UIHandler()

    http.HandleFunc("/sensors/", restfulSensorsHandler)
    http.HandleFunc("/nodes/", restfulNodesHandler)
    http.HandleFunc("/media/", restfulMediaHandler)

    // Handle webcam stream requests
    http.Handle(
        "/websocket/", 
        websocket.Handler(websocketHandler),
    )

    var addr_flag = flag.String(
        "addr", 
        "localhost", 
        "Address the http server binds to",
    )

    var port_flag = flag.String(
        "port", 
        "8080", 
        "Port used for http server",
    )

    flag.Parse()

    timePrinter("Running server on " + *addr_flag + ":" + *port_flag)
    http.ListenAndServe(*addr_flag + ":" + *port_flag, nil)
}
