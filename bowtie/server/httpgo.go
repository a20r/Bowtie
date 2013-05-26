package main

// Imports
import (
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "strings"
    "encoding/json"
)

// JSON response mapping
type Response map[string]interface{}

// Type definition for disambiguation. Holds the sensor data
type SensorData Response

// Represents an file loaded
type Page struct {
	Title string
	Body []byte
}

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

// Creates a function that will be used as a handler
// for static and template responses. See Usage!
func fileResponseCreator(folder string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
                fmt.Println("GET\t" + r.URL.Path)
    			var p *Page
    			var err error
    			if len(r.URL.Path) == 1 {
                    // In case the path is just '/'
    				p, err = loadPage("templates", "index.html")
    			} else {
    				p, err = loadPage(folder, r.URL.Path[1:])
    			}
    			if p != nil {
        			w.Write(p.Body)
        		} else {
        			fmt.Println("ERROR\t" + err.Error())
        		}
		}
}

// Removes the JSON data once the node stops
// sending sensor data
func dataRemoveHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("GET\t" + r.URL.Path)
    urlVars := strings.Split(r.URL.Path[1:], "/")
    cpu_id, node_id := urlVars[1], urlVars[2]
    err := os.Remove("json_data/" + cpu_id + "/" + node_id + ".json")
    if err != nil {
        fmt.Println("ERROR\t" + err.Error())
    }
}

// Handler called when data is sent
// to the server from a node
func dataSentHandler(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    fmt.Println("POST\t" + r.URL.Path)
    urlVars := strings.Split(r.URL.Path[1:], "/")
    cpu_id, node_id := urlVars[1], urlVars[2]
    path := "./json_data/" + cpu_id + "/"
    os.Mkdir(path, os.ModePerm | os.ModeType)
    file, err := os.Create(path + node_id + ".json")
    if err != nil {
        fmt.Println("ERROR\t" + err.Error())
        return
    }
    file.Write([]byte(r.Form["sensor_data"][0]))
    file.Close()
}

// Responds to the GET request from a client.
// Used for the visualization and for APIs 
// for users to query the data.
func dataGetHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("GET\t" + r.URL.Path)
    urlVars := strings.Split(r.URL.Path[1:], "/")
    cpu_id := urlVars[1]
    files, err := ioutil.ReadDir("json_data/" + cpu_id)
    if err != nil {
        res := Response{"error": Response{"code": 2, "message": "No data for " + cpu_id}}
        fmt.Fprint(w, res)
        fmt.Println("ERROR\t" + err.Error())
    } else {
        res := make(Response)
        for _, file := range files {
            //fmt.Println(file.Name())
            var sData SensorData
            node_id := strings.Split(file.Name(), ".")[0]
            file_bytes, read_err := ioutil.ReadFile("json_data/" + cpu_id + "/" + node_id + ".json")
            json_err := json.Unmarshal(file_bytes, &sData)
            if read_err != nil {
                fmt.Println("ERROR\t" + read_err.Error())
            }
            if json_err != nil {
                fmt.Println("ERROR\t" + json_err.Error())
            }
            res[node_id] = sData
        }
        //fmt.Println(files)
        if len(files) > 0 {
            res["error"] = Response{"code": 0, "message": "No error"}
        } else {
            res["error"] = Response{"code": 2, "message": "No data for " + cpu_id} 
        }
        fmt.Println("RESPONSE\t" + res.String())
        fmt.Fprint(w, res)
    }
}

// Handles all Javascript, images, and HTML
// file requests
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
    http.HandleFunc("/unchecked/", dataRemoveHandler)
    http.HandleFunc("/get_data/", dataGetHandler)
    fmt.Println("Running server on localhost:8080")
    http.ListenAndServe(":8080", nil)
}
