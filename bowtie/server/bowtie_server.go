
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

    // rethinkdb
    rethink "github.com/christopherhesse/rethinkgo"

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

// Media slice
type MediaSlice struct {
    Media_Type string
    Group_ID string
    Node_ID string
    Data string
}

// database session
var session, dbErr = rethink.Connect("localhost:28015", "bowtie_db")

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
            file_bytes, read_err := ioutil.ReadFile(
                "json_data/" + 
                group_id + "/" + 
                node_id + ".json",
            )
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
func videoStreamHandler(ms MediaSlice) {
    group_id := ms.Group_ID
    node_id := ms.Node_ID
    path := "./video_data/" + group_id + "/"

    // Make and log data to a file
    os.Mkdir(path, os.ModePerm | os.ModeType)
    file, err := os.Create(path + node_id + ".jpg")
    if err != nil {
        fmt.Println("ERROR\t" + err.Error())
        return
    }

    // Decode Base64 string to binary
    data_header := strings.Split(ms.Data, ",")[0]
    data_raw := strings.Split(ms.Data, ",")[1]

    if (data_header == "data:image/jpeg;base64") {
        img_data, err := base64.StdEncoding.DecodeString(data_raw)
        if err != nil {
            fmt.Println("error:", err)
            return
        }
        file.Write([]byte(img_data))

    } else {
        err := "video data format [" + data_header + "] not supported!"
        fmt.Println("error:", err)
    }

    // Write out the image binary
    file.Close()
}

// Audio stream handler
// Obtains data as a string encoded in Base64 and outputs the audio
// stream as a single wav file
func audioStreamHandler(ms MediaSlice) {
    group_id := ms.Group_ID
    node_id := ms.Node_ID
    path := "./audio_data/" + group_id + "/"

    // Make and log data to a file
    os.Mkdir(path, os.ModePerm | os.ModeType)
    file, err := os.Create(path + node_id + ".wav")
    if err != nil {
        fmt.Println("ERROR\t" + err.Error())
        return
    }

    // Decode Base64 string to binary
    data_header := strings.Split(ms.Data, ",")[0]
    data_raw := strings.Split(ms.Data, ",")[1]

    if (data_header == "data:audio/wav;base64") {
        audio_data, err := base64.StdEncoding.DecodeString(data_raw)
        if err != nil {
            fmt.Println("error:", err)
            return
        }

        // Write out the image binary
        file.Write([]byte(audio_data))
    } else {
        err := "audio data format [" + data_header + "] not supported!"
        fmt.Println("error:", err)
    }

    file.Close()
}

// Websocket Parser
func websocketMsgParser(msg string) {
    b := []byte(msg)
    var ms MediaSlice

    err := json.Unmarshal(b, &ms)
    if err != nil {
        fmt.Println("ProcessSocket:\tgot error", err)
        return
    }

    fmt.Println("Parsing Websocket message [" + ms.Media_Type + "]")
    if (ms.Media_Type == "video") {
        videoStreamHandler(ms)
    } else if (ms.Media_Type == "audio") {
        audioStreamHandler(ms)
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
            fmt.Println("ProcessSocket:\tgot error", err)
            _ = websocket.Message.Send(ws, "FAIL:" + err.Error())
            return
        }
        // fmt.Println("ProcessSocket: got message", msg)
        websocketMsgParser(msg)
    }

    fmt.Println("Finish handling websocket with wsHandler")
}

// func restfulHandler(w http.ResponseWriter, r *http.Request) {
//     switch r.Method {
//         case "GET":
//             restfulGet(w, r)
//         case "PUT", "POST":
//             restfulPost(w, r)
//         default:
//             fmt.Println("ERROR:\tUnknown request method")
//     }
// }

// func restfulGet(w http.ResponseWriter, r *http.Request) {

// }

// /*
//     Handles the sensor data posting. The actual sensor
//     data is stored in the JSON form.

//     The current structure of this JSON is as follows:

//     sensorData : JSON.stringify(
//         {
//             value : `value of the sensor being sent`
//             type : `the data type of value`
//             time : `time stamp from when it was sent`
//             //token : `authentication token`
//         }
//     )
// */
// func restfulPost(w http.ResponseWriter, r *http.Request) {
//     fmt.Println("POST\t" + r.URL.Path)

//     // decodes the JSON data to be sent to the database
//     var sData SensorData
//     r.ParseForm()
//     json.Unmarshal([]byte (r.Form["sensorData"][0]), &sData)


//     groupId, nodeId, sensor := parseRestfulURL(r.URL.Path)
//     // checks if the entry is already in the database
//     var groupData []interface{}
//     rethink.Table("sensor_table").GetAll(
//         "groupId", 
//         groupId,
//     ).Run(session).All(&groupData)

//     entryExists := len(groupData) > 0

//     // FIX THIS SHIT BRO!
//     if entryExists {

//         var nodeExists bool
//         rethink.Table("sensor_table").GetAll(
//             "groupId",
//             groupId,
//         ).Nth(0).Attr("nodes").Contains(nodeId).Run(session).One(&nodeExists)

//         if nodeExists {
//             var mergedNode interface{}
//             rethink.Table("sensor_table").GetAll(
//                 "groupId",
//                 groupId,
//             ).Nth(0).Attr("nodes").Attr(nodeId).Merge(
//                 rethink.Map{
//                     sensor : rethink.Map{
//                         "value" : sData["value"],
//                         "type" : sData["type"],
//                         "time" : sData["time"],
//                     },
//                 },
//             ).Run(session).One(&mergedNode)

//             rethink.Table("sensor_table").GetAll(
//                 "groupId",
//                 groupId,
//             ).Nth(0).Attr("nodes")
//         }
//     } else {
//         rethink.Table("sensor_table").Insert(
//             rethink.Map{
//                 "groupId" : groupId,
//                 "nodes" : rethink.Map{
//                     nodeId : rethink.Map{
//                         sensor : rethink.Map{
//                             "value" : sData["value"],
//                             "type" : sData["type"],
//                             "time" : sData["time"],
//                         },
//                     },
//                 },
//             },
//         ).Run(session).Exec()
//     }
// }

// func parseRestfulURL(
//     // params
//     URLStr string,
// ) (
//     // return values
//     groupId string, 
//     nodeId string, 
//     sensor string,
// ) {
//     var splitURL = strings.Split(URLStr[1:], "/")

//     if len(splitURL) >= 4 {
//         groupId = splitURL[1]
//         nodeId = splitURL[2]
//         sensor = splitURL[3]
//     }

//     return
// }

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

    if dbErr != nil {
        fmt.Println(dbErr)
        return
    } 

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


