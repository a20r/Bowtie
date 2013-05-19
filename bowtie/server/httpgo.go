package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "strings"
    "encoding/json"
)

type Response map[string]interface{}
type SensorData Response

type Page struct {
	Title string
	Body []byte
}

func (r Response) String() (s string) {
    b, err := json.Marshal(r)
    if err != nil {
            s = ""
            return
    }
    s = string(b)
    return
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

func dataRemoveHandler(w http.ResponseWriter, r *http.Request) {
    urlVars := strings.Split(r.URL.Path[1:], "/")
    cpu_id, node_id := urlVars[1], urlVars[2]
    err := os.Remove("json_data/" + cpu_id + "/" + node_id + ".json")
    if err != nil {
        fmt.Println(err)
    }
}

func dataSentHandler(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    urlVars := strings.Split(r.URL.Path[1:], "/")
    cpu_id, node_id := urlVars[1], urlVars[2]
    path := "./json_data/" + cpu_id + "/"
    os.Mkdir(path, os.ModePerm | os.ModeType)
    file, err := os.Create(path + node_id + ".json")
    if err != nil {
        fmt.Println(err)
        return
    }
    file.Write([]byte(r.Form["sensor_data"][0]))
    file.Close()
}

func dataGetHandler(w http.ResponseWriter, r *http.Request) {
    urlVars := strings.Split(r.URL.Path[1:], "/")
    cpu_id := urlVars[1]
    files, err := ioutil.ReadDir("json_data/" + cpu_id)
    if err != nil {
        res := Response{"error": Response{"code": 2, "message": "No data for " + cpu_id}}
        fmt.Fprint(w, res)
        fmt.Println(err)
    } else {
        res := make(Response)
        for _, file := range files {
            //fmt.Println(file.Name())
            node_id := strings.Split(file.Name(), ".")[0]
            var sData SensorData
            file_bytes, read_err := ioutil.ReadFile("json_data/" + cpu_id + "/" + node_id + ".json")
            json_err := json.Unmarshal(file_bytes, &sData)
            if read_err != nil {
                fmt.Println(read_err)
            }
            if json_err != nil {
                fmt.Println(json_err)
            }
            res[node_id] = sData
        }
        res["error"] = Response{"code": 0, "message": "No error"}
        fmt.Println(res)
        fmt.Fprint(w, res)
    }
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
    http.HandleFunc("/unchecked/", dataRemoveHandler)
    http.HandleFunc("/get_data/", dataGetHandler)
    http.ListenAndServe(":8080", nil)
}
