
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
    "errors"

    //ADTs
    "time"

    // rethinkdb
    rethink "github.com/christopherhesse/rethinkgo"

    // custom pkgs
)

// JSON response mapping
type Response map[string]interface{}

// Type definition for disambiguation. Holds the sensor data
type SensorData map[string]interface{}

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

// Removes the JSON data once the node stops
// sending sensor data
func dataRemoveHandler(w http.ResponseWriter, r *http.Request) {
    timePrinter("GET\t" + r.URL.Path)
    urlVars := strings.Split(r.URL.Path[1:], "/")
    group_id, node_id := urlVars[1], urlVars[2]
    err := os.Remove("json_data/" + group_id + "/" + node_id + ".json")

    if err != nil {
        timePrinter("ERROR\t" + err.Error())
    }
}

// Handler called when data is sent
// to the server from a node
func dataSentHandler(w http.ResponseWriter, r *http.Request) {
    timePrinter("POST\t" + r.URL.Path)

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
        timePrinter("ERROR\t" + err.Error())
        return
    }

    file.Write([]byte(r.Form["sensor_data"][0]))
    file.Close()
}

// Responds to the GET request from a client.
// Used for the visualization and for APIs
// for users to query the data.
func dataGetHandler(w http.ResponseWriter, r *http.Request) {
    timePrinter("GET\t" + r.URL.Path)
    urlVars := strings.Split(r.URL.Path[1:], "/")
    group_id := urlVars[1]
    files, err := ioutil.ReadDir("json_data/" + group_id)

    if err != nil {
        res := Response{"error": Response{"code": 2, "message": "No data for " + group_id}}
        fmt.Fprint(w, res)
        timePrinter("ERROR\t" + err.Error())

    } else {
        res := make(Response)

        for _, file := range files {
            //timePrinter(file.Name())
            var sData SensorData
            node_id := strings.Split(file.Name(), ".")[0]
            file_bytes, read_err := ioutil.ReadFile(
                "json_data/" + 
                group_id + "/" + 
                node_id + ".json",
            )
            json_err := json.Unmarshal(file_bytes, &sData)

            if read_err != nil {
                timePrinter("ERROR\t" + read_err.Error())
            }
            if json_err != nil {
                timePrinter("ERROR\t" + json_err.Error())
            }
            res[node_id] = sData
        }

        //timePrinter(files)
        if len(files) > 0 {
            res["error"] = Response{"code": 0, "message": "No error"}
        } else {
            res["error"] = Response{"code": 2, "message": "No data for " + group_id}
        }

        timePrinter("RESPONSE\t" + res.String())
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
        timePrinter("ERROR\t" + err.Error())
        return
    }

    // Decode Base64 string to binary
    data_header := strings.Split(ms.Data, ",")[0]
    data_raw := strings.Split(ms.Data, ",")[1]

    if (data_header == "data:image/jpeg;base64") {
        img_data, err := base64.StdEncoding.DecodeString(data_raw)
        if err != nil {
            timePrinter("ERROR\t" + err.Error())
            return
        }
        file.Write([]byte(img_data))

    } else {
        err := "video data format [" + data_header + "] not supported!"
        timePrinter("ERROR\t" + err)
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
        timePrinter("ERROR\t" + err.Error())
        return
    }

    // Decode Base64 string to binary
    data_header := strings.Split(ms.Data, ",")[0]
    data_raw := strings.Split(ms.Data, ",")[1]

    if (data_header == "data:audio/wav;base64") {
        audio_data, err := base64.StdEncoding.DecodeString(data_raw)
        if err != nil {
            timePrinter("ERROR:\t" + err.Error())
            return
        }

        // Write out the image binary
        file.Write([]byte(audio_data))
    } else {
        err := "audio data format [" + data_header + "] not supported!"
        timePrinter("ERROR:\t" +  err)
    }

    file.Close()
}

// Websocket Parser
func websocketMsgParser(msg string) {
    b := []byte(msg)
    var ms MediaSlice

    err := json.Unmarshal(b, &ms)
    if err != nil {
        timePrinter("ERROR:\tSocket --> " + err.Error())
        return
    }

    timePrinter("Parsing Websocket message [" + ms.Media_Type + "]")
    if (ms.Media_Type == "video") {
        videoStreamHandler(ms)
    } else if (ms.Media_Type == "audio") {
        audioStreamHandler(ms)
    }
}

// Websocket Handler
func websocketHandler(ws *websocket.Conn) {
    timePrinter("Handling websocket request with wsHandler")
    var msg string

    // Process incomming websocket messages
    for {
        err := websocket.Message.Receive(ws, &msg)
        if err != nil {
            timePrinter("ERROR:\tSocket --> " + err.Error())
            _ = websocket.Message.Send(ws, "FAIL:" + err.Error())
            return
        }
        // timePrinter("ProcessSocket: got message", msg)
        websocketMsgParser(msg)
    }

    timePrinter("Finish handling websocket with wsHandler")
}

// using a more defined type for the restful API
type NodeSensorData struct {
    Value interface{} `value` 
    Type string `type`
    Time string `time`
}

func (nsd NodeSensorData) String() string {
    bytes, err := json.Marshal(nsd)
    if err != nil {
        return ""
    } else {
        return string(bytes)
    }
}

// for easier querying
type BowtieQueries struct {
    Session *rethink.Session
    GroupId string
    NodeId string
    Sensor string
}

func (bq BowtieQueries) GroupExists() bool {
    var groupData rethink.Map
    rethink.Table("sensor_table").Get( 
        bq.GroupId,
    ).Run(bq.Session).One(&groupData)

    return groupData != nil
}

func (bq BowtieQueries) NodeExists() bool {
    var nodeExists bool
    rethink.Table("sensor_table").Get(
        bq.GroupId,
    ).Attr("nodes").Contains(
        bq.NodeId,
    ).Run(bq.Session).One(&nodeExists)

    return nodeExists
}

func (bq BowtieQueries) InsertGroupWithData(sData NodeSensorData) {
    rethink.Table("sensor_table").Insert(
        rethink.Map{
            "groupId" : bq.GroupId,
            "nodes" : rethink.Map{
                bq.NodeId : rethink.Map{
                    bq.Sensor : rethink.Map{
                        "value" : sData.Value,
                        "type" : sData.Type,
                        "time" : sData.Time,
                    },
                },
            },
        },
    ).Run(bq.Session).Exec()
}

func (bq BowtieQueries) CreateGroup() {
    rethink.Table("sensor_table").Insert(
        rethink.Map{
            "groupId" : bq.GroupId,
            "nodes" : rethink.Map{},
        },
    ).Run(bq.Session).Exec()
}

// creates node if not there, updates if it exists
func (bq BowtieQueries) UpdateSensor(sData NodeSensorData) error {

    if !bq.GroupExists() {
        bq.CreateGroup()
    }

    var nodes map[string]rethink.Map
    rethink.Table("sensor_table").Get(
        bq.GroupId,
    ).Attr("nodes").Run(bq.Session).One(&nodes)

    if len(nodes[bq.NodeId]) > 0 {
        nodes[bq.NodeId][bq.Sensor] = rethink.Map{
            "value" : sData.Value,
            "type" : sData.Type,
            "time" : sData.Time,
        }
    } else { 
        nodes[bq.NodeId] = rethink.Map{
            bq.Sensor : rethink.Map{
                "value" : sData.Value,
                "type" : sData.Type,
                "time" : sData.Time,
            },
        }
    }

    rethink.Table("sensor_table").Get(
        bq.GroupId,
    ).Update(
        rethink.Map{
            "nodes" : nodes,
        },
    ).Run(bq.Session).Exec()

    return nil
}

func (bq BowtieQueries) UpdateNode(sDataMap map[string]NodeSensorData) error {
    if !bq.GroupExists() {
        bq.CreateGroup() // check this
    }

    var nodes map[string]rethink.Map
    rethink.Table("sensor_table").Get(
        bq.GroupId,
    ).Attr("nodes").Run(bq.Session).One(&nodes)

    if len(nodes[bq.NodeId]) == 0 {
        nodes[bq.NodeId] = make(rethink.Map)
    }

    for key, sensorNode := range sDataMap {
        nodes[bq.NodeId][key] = rethink.Map{
            "value" : sensorNode.Value,
            "type" : sensorNode.Type,
            "time" : sensorNode.Time,
        }
    }

    rethink.Table("sensor_table").Get(
        bq.GroupId,
    ).Update(
        rethink.Map{
            "nodes" : nodes,
        },
    ).Run(bq.Session).Exec()

    return nil
}

func (bq BowtieQueries) GetSensorData() (*NodeSensorData, error) {
    node, err := bq.GetNode()
    if err != nil {
        return nil, err
    }

    if node[bq.Sensor] == nil {
        return nil, errors.New("Sensor does not exist")
    }

    sensor := node[bq.Sensor].(map[string]interface{})
    return &NodeSensorData{
        sensor["value"],
        sensor["type"].(string),
        sensor["time"].(string),
    }, nil
}

func (bq BowtieQueries) GetNode() (rethink.Map, error) {
    if !bq.GroupExists() {
        return nil, errors.New("Group does not yet exist")
    }
    var nodes map[string]rethink.Map
    rethink.Table("sensor_table").Get(
        bq.GroupId,
    ).Attr("nodes").Run(bq.Session).One(&nodes)

    if len(nodes[bq.NodeId]) == 0 {
        return nil, errors.New("Node does not exist")
    }

    return nodes[bq.NodeId], nil
}

func (bq BowtieQueries) GetGroup() (rethink.Map, error) {
    if !bq.GroupExists() {
        return nil, errors.New("Group does not yet exist")
    }
    var group rethink.Map
    rethink.Table("sensor_table").Get(
        bq.GroupId,
    ).Run(bq.Session).One(&group)

    return group, nil
}

func (bq BowtieQueries) DeleteGroup() error {
    rethink.Table("sensor_table").Get(
        bq.GroupId,
    ).Delete().Run(bq.Session).Exec()
    return nil
}

func (bq BowtieQueries) DeleteNode() error {
    group, err := bq.GetGroup()

    if err != nil {
        timePrinter("ERROR\t" + err.Error())
        return err
    }

    nodesMap := group["nodes"].(map[string]interface{})
    fmt.Println(nodesMap)

    delete(
        nodesMap, 
        bq.NodeId,
    )

    rethink.Table("sensor_table").Get(
        bq.GroupId,
    ).Update(
        rethink.Map{
            "nodes" : nodesMap,
        },
    ).Run(bq.Session).Exec()

    return nil
}

func (bq BowtieQueries) DeleteSensor() error {
    group, err := bq.GetGroup()

    if err != nil {
        timePrinter("ERROR\t" + err.Error())
        return err
    }
    
    delete(
        group["nodes"].(map[string]interface{})[bq.NodeId].(map[string]interface{}),
        bq.Sensor,
    )

    rethink.Table("sensor_table").Get(
        bq.GroupId,
    ).Update(
        rethink.Map{
            "nodes" : group["nodes"],
        },
    ).Run(bq.Session).Exec()

    return nil
}

func makeBowtieQueriesWithPath(
    URLStr string, 
    rethinkSession *rethink.Session,
) *BowtieQueries {
    groupId, nodeId, sensor := parseRestfulURL(URLStr)
    bq := BowtieQueries{
        rethinkSession,
        groupId,
        nodeId,
        sensor,
    }
    return &bq
}

func parseRestfulURL(
    // params
    URLStr string,
) (
    // return values
    groupId string, 
    nodeId string, 
    sensor string,
) {
    var splitURL = strings.Split(URLStr[1:], "/")

    if len(splitURL) >= 2 {
        groupId = splitURL[1]
        if len(splitURL) >= 3 {
            nodeId = splitURL[2]
            if len(splitURL) >= 4 {
                sensor = splitURL[3]
            }
        }
    }
    return
}

func restfulHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
        case "GET":
            restfulGet(w, r)
        case "POST":
            restfulPost(w, r)
        // when all sensor data is sent at one time
        case "PUT":
            restfulPut(w, r)
        case "DELETE":
            restfulDelete(w, r)
        default:
            timePrinter("ERROR:\tUnknown request method")
    }
}

func restfulDelete(w http.ResponseWriter, r *http.Request) {
    timePrinter("DELETE\t" + r.URL.Path)
    bq := makeBowtieQueriesWithPath(r.URL.Path, session)
    if bq.NodeId == "" && bq.Sensor == "" {
        bq.DeleteGroup()
    } else if bq.Sensor == "" {
        bq.DeleteNode()
    } else {
        bq.DeleteSensor()
    }
}

func restfulGet(w http.ResponseWriter, r *http.Request) {
    timePrinter("GET\t" + r.URL.Path)
    bq := makeBowtieQueriesWithPath(r.URL.Path, session)

    if bq.Sensor != "" && bq.NodeId != "" {
        var sensorData *NodeSensorData
        sensorData, err := bq.GetSensorData()
        if err != nil {
            fmt.Fprint(
                w, 
                Response{
                    "Error" : err.Error(),
                },
            )
        } else {
            fmt.Fprint(w, sensorData)
        }
    } else {

        var bytes []byte
        var err error
        var marshalData interface{}
        var group rethink.Map

        group, err = bq.GetGroup()
        if err != nil {
            fmt.Fprint(
                w,
                Response{
                    "error" : 1, 
                    "message": err.Error(),
                },
            )
            return
        }

        if bq.Sensor == "" && bq.NodeId == "" {
            marshalData = group["nodes"]
        } else if bq.Sensor == ""{
            marshalData = group["nodes"].(map[string]interface{})[bq.NodeId]
            if marshalData == nil {
               fmt.Fprint(
                    w,
                    Response{
                        "error" : 1, 
                        "message": "Node does not exist",
                    },
                )
                return 
            }
        } 

        bytes, err = json.Marshal(marshalData)

        if err != nil {
           fmt.Fprint(
                w,
                Response{"error" : 1, "message": err.Error()},
            )
            return 
        }

        w.Write(bytes)

    }
}

/*
    Handles the sensor data posting. The actual sensor
    data is stored in the JSON form.

    The current structure of this JSON is as follows:

    sensorData : JSON.stringify(
        {
            value : `value of the sensor being sent`
            type : `the data type of value`
            time : `time stamp from when it was sent`
            //token : `authentication token`
        }
    )
*/
func restfulPost(w http.ResponseWriter, r *http.Request) {
    timePrinter("POST\t" + r.URL.Path)

    // decodes the JSON data to be sent to the database
    var sData NodeSensorData
    r.ParseForm()

    err := json.Unmarshal(
        []byte (r.Form["sensorData"][0]), 
        &sData,
    )

    if err != nil {
        fmt.Fprint(
            w, 
            Response{"error": 1, "message": err.Error()},
        );
        return
    }

    bq := makeBowtieQueriesWithPath(r.URL.Path, session)

    bq.UpdateSensor(sData)

    fmt.Fprint(w, Response{"error": 0, "message": "All went well"});
}

func restfulPut(w http.ResponseWriter, r *http.Request) {
    timePrinter("PUT\t" + r.URL.Path)
    groupId := strings.Split(r.URL.Path[1:], "/")[1]
    nodeId := strings.Split(r.URL.Path[1:], "/")[2]

    var sDataMap map[string]NodeSensorData
    r.ParseForm()

    err := json.Unmarshal(
        []byte (r.Form["sensorData"][0]),
        &sDataMap,
    )

    if err != nil {
        fmt.Fprint(
            w, 
            Response{"error": 1, "message": err.Error()},
        );
        return
    }

    bq := BowtieQueries{session, groupId, nodeId, ""}
    err = bq.UpdateNode(sDataMap)
    if err != nil {
        fmt.Fprint(
            w, 
            Response{"error": 1, "message": err.Error()},
        );
        return
    }

    fmt.Fprint(
        w, 
        Response{"error": 0, "message": "Everything is great"},
    );
}

func restfulNodesHandler(w http.ResponseWriter, r *http.Request) {
    timePrinter("GET\t" + r.URL.Path)
    groupId := strings.Split(r.URL.Path[1:], "/")[1]
    var nodes map[string]rethink.Map
    rethink.Table("sensor_table").Get(
        groupId,
    ).Attr("nodes").Run(session).One(&nodes)

    nodeIds := make([]string, len(nodes))
    i := 0
    for key := range nodes {
        nodeIds[i] = key
        i = i + 1
    }

    bytes, err := json.Marshal(nodeIds)
    if err != nil {
        fmt.Fprint(w, Response{"Error" : "Not able to marshal JSON data"})
    } else {
        w.Write(bytes) 
    }
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

    if dbErr != nil {
        timePrinter(dbErr.Error())
        return
    } 
    requestHandler()

    http.HandleFunc("/checked/", dataSentHandler)
    http.HandleFunc("/unchecked/", dataRemoveHandler)
    http.HandleFunc("/get_data/", dataGetHandler)

    http.HandleFunc("/sensors/", restfulHandler)
    http.HandleFunc("/nodes/", restfulNodesHandler)

    var addr_flag = flag.String("addr", "localhost", "Address the http server binds to")
    var port_flag = flag.String("port", "8080", "Port used for http server")

    flag.Parse()

    //timePrinter("Running server on " + *addr_flag + ":" + *port_flag)
    http.ListenAndServe(*addr_flag + ":" + *port_flag, nil)
}


