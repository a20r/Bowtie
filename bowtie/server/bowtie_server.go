package main

import (
    // sys pkgs
    "os"
    "flag"

    // io pkgs
    "fmt"
    "io/ioutil"

    // network pkgs
    "net/http"
    "code.google.com/p/go.net/websocket"

    // string pkgs
    "strings"
    "encoding/json"
    "encoding/base64"

    // custom pkgs
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
        var p *Page
        var err error

        fmt.Println("GET\t" + r.URL.Path)

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
    group_id, node_id := urlVars[1], urlVars[2]
    err := os.Remove("json_data/" + group_id + "/" + node_id + ".json")

    if err != nil {
        fmt.Println("ERROR\t" + err.Error())
    }
}

// Handler called when data is sent
// to the server from a node
func dataSentHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("POST\t" + r.URL.Path)

    // Parse form and extract data details
    r.ParseForm()
    urlVars := strings.Split(r.URL.Path[1:], "/")
    group_id := urlVars[1]
    node_id := urlVars[2]
    path := "./json_data/" + group_id + "/"

    // Make and log data to a file
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
    group_id := urlVars[1]
    files, err := ioutil.ReadDir("json_data/" + group_id)

    if err != nil {
        res := Response{"error": Response{"code": 2, "message": "No data for " + group_id}}
        fmt.Fprint(w, res)
        fmt.Println("ERROR\t" + err.Error())

    } else {
        res := make(Response)

        for _, file := range files {
            //fmt.Println(file.Name())
            var sData SensorData
            node_id := strings.Split(file.Name(), ".")[0]
            file_bytes, read_err := ioutil.ReadFile("json_data/" + group_id + "/" + node_id + ".json")
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
            res["error"] = Response{"code": 2, "message": "No data for " + group_id}
        }

        fmt.Println("RESPONSE\t" + res.String())
        fmt.Fprint(w, res)
    }
}

// Video stream handler
// Obtains data as a string encoded in Base64 and outputs the video
// stream a single image
func videoStreamHandler(data string) {
    group_id := "testing"
    node_id := "testing"
    path := "./video_data/" + group_id + "/"

    // Make and log data to a file
    os.Mkdir(path, os.ModePerm | os.ModeType)
    file, err := os.Create(path + node_id + ".jpg")
    if err != nil {
        fmt.Println("ERROR\t" + err.Error())
        return
    }

    // Decode Base64 string to binary
    img_data, err := base64.StdEncoding.DecodeString(data)
    if err != nil {
        fmt.Println("error:", err)
        return
    }

    // Write out the image binary
    file.Write([]byte(img_data))
    file.Close()
}

// Audio stream handler
// Obtains data as a string encoded in Base64 and outputs the audio
// stream as a single wav file
func audioStreamHandler(data string) {
    group_id := "testing"
    node_id := "testing"
    path := "./audio_data/" + group_id + "/"


    // Make and log data to a file
    os.Mkdir(path, os.ModePerm | os.ModeType)
    file, err := os.Create(path + node_id + ".wav")
    if err != nil {
        fmt.Println("ERROR\t" + err.Error())
        return
    }

    // Decode Base64 string to binary
    audio_data, err := base64.StdEncoding.DecodeString(data)
    if err != nil {
        fmt.Println("error:", err)
        return
    }

    // Write out the image binary
    file.Write([]byte(audio_data))
    file.Close()
}

// Websocket Parser
func websocketMsgParser(msg string) {
    // Parse header and data
    msg_header := strings.Split(msg, ",")[0]
    msg_data := strings.Split(msg, ",")[1]

    fmt.Println("Parsing Websocket message [" + msg_header + "]")
    if (msg_header == "data:image/jpeg;base64") {
        videoStreamHandler(msg_data)
    } else if (msg_header == "data:audio/wav;base64") {
        audioStreamHandler(msg_data)
    }
}

// Websocket Handler
func websocketHandler(ws *websocket.Conn) {
    fmt.Println("Handling websocket request with wsHandler")
    var msg string

    // Process incomming websocket messages
    for {
        err := websocket.Message.Receive(ws, &msg)
        if err != nil {
            fmt.Println("ProcessSocket: got error", err)
            _ = websocket.Message.Send(ws, "FAIL:" + err.Error())
            return
        }
        // fmt.Println("ProcessSocket: got message", msg)
        websocketMsgParser(msg)
    }

    fmt.Println("Finish handling websocket with wsHandler")
}

// Handles all incomming http requests
func requestHandler() {
    staticHandler := fileResponseCreator("static")
    http.HandleFunc("/", fileResponseCreator("templates"))
    http.HandleFunc("/css/", staticHandler)
    http.HandleFunc("/js/", staticHandler)
    http.HandleFunc("/img/", staticHandler)
    http.HandleFunc("/favicon.ico", fileResponseCreator("static/img"))

    // Handle webcam stream requests
    http.Handle("/websocket/", websocket.Handler(websocketHandler))
}

// MAIN EXECUTION FLOW
func main() {
    requestHandler()

    http.HandleFunc("/checked/", dataSentHandler)
    http.HandleFunc("/unchecked/", dataRemoveHandler)
    http.HandleFunc("/get_data/", dataGetHandler)

    var addr_flag = flag.String("addr", "localhost", "Address the http server binds to")
    var port_flag = flag.String("port", "8080", "Port used for http server")

    flag.Parse()

    //fmt.Println("Running server on " + *addr_flag + ":" + *port_flag)
    http.ListenAndServe(*addr_flag + ":" + *port_flag, nil)
}


