
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

    // custom pkgs
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

// Media slice
type MediaSlice struct {
    Media_Type string
    Group_ID string
    Node_ID string
    Data string
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
type SensorData struct {
    Value interface{} `value` 
    Type string `type`
    Time string `time`
}

type NodeData map[string]SensorData
type GroupData map[string]NodeData
 
/*
func (nsd SensorData) String() string {
    bytes := nsd.ToBytes()
    return string(bytes)
}

func (nsd SensorData) ToBytes() []byte {
    bytes, err := json.Marshal(nsd)
    if err != nil {
        return make([]byte, 0)
    } else {
        return bytes
    }
}
*/

// for easier querying
type BowtieQueries struct {
    GroupId string
    NodeId string
    Sensor string
}

func (bq BowtieQueries) GroupExists() bool {
    files, err := ioutil.ReadDir(
        "./json_data/" + bq.GroupId,
    )

    return err == nil && len(files) > 0
}

func (bq BowtieQueries) NodeExists() bool {
    _, readErr := ioutil.ReadFile(
        "json_data/" + 
        bq.GroupId + "/" + 
        bq.NodeId + ".json",
    )

    return readErr == nil
}

func (bq BowtieQueries) CreateGroup() error {
    path := "./json_data/" + bq.GroupId + "/"

    err := os.Mkdir(path, os.ModePerm | os.ModeType)

    if err != nil {
        return errors.New(
            "Unable to create group --> " +
            err.Error(),
        )
    } else {
        return nil
    }
}

// creates node if not there, updates if it exists
func (bq BowtieQueries) UpdateSensor(sData SensorData) error {

    if !bq.GroupExists() {
        bq.CreateGroup()
    }

    sDataMap, err := bq.GetNode()

    if err != nil {
        sDataMap = make(NodeData)
    }

    sDataMap[bq.Sensor] = sData

    path := "./json_data/" + bq.GroupId + "/"
    file, err := os.Create(path + bq.NodeId + ".json")

    if err != nil {
        return errors.New(
            "Unable to update sensor --> " +
            err.Error(),
        )
    }

    defer file.Close()

    bytes, jsonErr := json.Marshal(sDataMap)

    if jsonErr != nil {
        return errors.New(
            "Unable to update sensor --> " +
            jsonErr.Error(),
        )
    }

    file.Write(bytes)

    return nil
}

func (bq BowtieQueries) UpdateNode(sDataMap NodeData) error {

    if !bq.GroupExists() {
        bq.CreateGroup() // check this
    }

    path := "./json_data/" + bq.GroupId + "/"
    file, err := os.Create(path + bq.NodeId + ".json")

    if err != nil {
        return errors.New(
            "Unable to update node --> " +
            err.Error(),
        )
    }

    defer file.Close()

    bytes, jsonErr := json.Marshal(sDataMap)

    if jsonErr != nil {
        return errors.New(
            "Unable to update node --> " +
            jsonErr.Error(),
        )
    }

    file.Write(bytes)

    return nil
}

func (bq BowtieQueries) GetSensor() (*SensorData, error) {
    node, err := bq.GetNode()
    if err != nil {
        return nil, err
    }

    if _, exists := node[bq.Sensor]; !exists {
        return nil, errors.New("Sensor does not exist")
    }

    return &SensorData{
        node[bq.Sensor].Value,
        node[bq.Sensor].Type,
        node[bq.Sensor].Time,
    }, nil
}

func (bq BowtieQueries) GetNode() (NodeData, error) {

    var sDataMap NodeData

    fileBytes, readErr := ioutil.ReadFile(
        "json_data/" + 
        bq.GroupId + "/" + 
        bq.NodeId + ".json",
    )

    if readErr != nil {
        return nil, errors.New(
            "Unable to get node --> " + 
            readErr.Error(),
        )
    }

    jsonErr := json.Unmarshal(fileBytes, &sDataMap)

    if jsonErr != nil {
        return nil, errors.New(
            "Unable to get node --> " + 
            jsonErr.Error(),
        )
    }

    return sDataMap, nil
}

func (bq BowtieQueries) GetGroup() (GroupData, error, int) {
    if !bq.GroupExists() {
        return nil, errors.New("Group does not yet exist"), 1
    }

    files, err := ioutil.ReadDir("json_data/" + bq.GroupId)

    if err != nil {
        return nil, errors.New(
            "Unable to get group --> " +
            err.Error(),
        ), 4
    }

    group := make(GroupData)

    for _, file := range files {
        if file.Name() == ".DS_Store" {
            continue
        }
        var sDataMap NodeData
        nodeId := strings.Split(file.Name(), ".")[0]
        fileBytes, readErr := ioutil.ReadFile(
            "json_data/" + 
            bq.GroupId + "/" + 
            nodeId + ".json",
        )

        if readErr != nil {
            return nil, errors.New(
                "Unable to get group --> " + 
                readErr.Error(),
            ), 3
        }

        jsonErr := json.Unmarshal(fileBytes, &sDataMap)

        if jsonErr != nil {
            return nil, errors.New(
                "Unable to get group --> " + 
                jsonErr.Error(),
            ), 2
        }

        group[nodeId] = sDataMap
    }
    return group, nil, 0
}

func (bq BowtieQueries) GetNodeArray() ([]string, error) {

    group, err, _ := bq.GetGroup()

    if err != nil {
        return nil, err
    }

    i := 0
    retArray := make([]string, len(group))
    for key, _ := range group {
        retArray[i] = key
        i = i + 1
    }

    return retArray, nil
}

func (bq BowtieQueries) DeleteGroup() error {
    err := os.Remove(
        "json_data/" + 
        bq.GroupId,
    )

    if err != nil {
        return errors.New(
            "Unable to delete group --> " + 
            err.Error(),
        )
    }

    return nil
}

func (bq BowtieQueries) DeleteNode() error {

    err := os.Remove(
        "json_data/" + 
        bq.GroupId + "/" +
        bq.NodeId + ".json",
    )

    if err != nil {
        return errors.New(
            "Unable to delete node --> " + 
            err.Error(),
        )
    }

    return nil
}

func (bq BowtieQueries) DeleteSensor() error {

    node, err := bq.GetNode()

    if err != nil {
        return errors.New(
            "Unable to delete sensor --> " + 
            err.Error(),
        )
    }

    delete(
        node,
        bq.Sensor,
    )

    return nil
}

func makeBowtieQueriesWithPath(
    URLStr string, 
) *BowtieQueries {
    groupId, nodeId, sensor := parseRestfulURL(URLStr)
    bq := BowtieQueries{
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
        case "DELETE":
            restfulDelete(w, r)
        default:
            timePrinter("ERROR:\tUnknown request method")
    }
}

func restfulDelete(w http.ResponseWriter, r *http.Request) {
    timePrinter("DELETE\t" + r.URL.Path)
    bq := makeBowtieQueriesWithPath(r.URL.Path)
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
    bq := makeBowtieQueriesWithPath(r.URL.Path)
    var bytes []byte
    var err error
    var marshalData interface{}
    errNum := 1

    if bq.Sensor == "" && bq.NodeId == "" {
        marshalData, err, errNum = bq.GetGroup()
    } else if bq.Sensor == ""{
        marshalData, err = bq.GetNode()
    }  else {
        marshalData, err = bq.GetSensor()
    }

    if err != nil {
        fmt.Fprint(
            w,
            Response{
                "Error" : errNum, 
                "Message": err.Error(),
            },
        )
        return
    }

    bytes, err = json.Marshal(marshalData)

    if err != nil {
       fmt.Fprint(
            w,
            Response{"Error" : 1, "Message": err.Error()},
        )
        return 
    }

    w.Write(bytes)
}

/*
    Handles the sensor data posting. The actual sensor
    data is stored in the JSON form.
*/
func restfulPost(w http.ResponseWriter, r *http.Request) {
    timePrinter("POST\t" + r.URL.Path)

    r.ParseForm()
    bq := makeBowtieQueriesWithPath(r.URL.Path)

    var err error

    if bq.NodeId == "" && bq.Sensor == "" {
        fmt.Fprint(
            w, 
            Response{
                "Error": 1, 
                "Message": "Not able to post a full group",
            },
        );
        return
    } else if bq.Sensor == "" {
        var sData NodeData
        err = json.Unmarshal(
            []byte (r.Form["sensorData"][0]), 
            &sData,
        )

        if err != nil {
            fmt.Fprint(
                w, 
                Response{
                    "Error": 1, 
                    "Message": err.Error(),
                },
            );
            return
        }
        err = bq.UpdateNode(sData)
    } else {
        var sData SensorData
        err = json.Unmarshal(
            []byte (r.Form["sensorData"][0]), 
            &sData,
        )

        if err != nil {
            fmt.Fprint(
                w, 
                Response{
                    "Error": 1, 
                    "Message": err.Error(),
                },
            );
            return
        }
        err = bq.UpdateSensor(sData)
    }

    if err != nil {
        fmt.Fprint(
            w, 
            Response{"Error": 1, "Message": err.Error()},
        );
        return
    }

    fmt.Fprint(w, Response{"Error": 0, "Message": "All went well"});
}

func restfulNodesHandler(w http.ResponseWriter, r *http.Request) {
    timePrinter("GET\t" + r.URL.Path)

    bq := makeBowtieQueriesWithPath(r.URL.Path)
    nodeIds, bqErr := bq.GetNodeArray()

    if bqErr != nil {
        fmt.Fprint(
            w, 
            Response{
                "Error" : 1, 
                "Message" : bqErr.Error(),
            },
        )
        return
    }

    bytes, err := json.Marshal(nodeIds)
    if err != nil {
        fmt.Fprint(
            w, 
            Response{
                "Error" : 1, 
                "Message" : err.Error(),
            },
        )
    } else {
        w.Write(bytes) 
    }
}

// Handles all in coming http requests
func UIHandler() {
    staticHandler := fileResponseCreator("static")
    http.HandleFunc("/", fileResponseCreator("templates"))
    http.HandleFunc("/css/", staticHandler)
    http.HandleFunc("/js/", staticHandler)
    http.HandleFunc("/img/", staticHandler)
    http.HandleFunc("/favicon.ico", fileResponseCreator("static/img"))
}

// MAIN EXECUTION FLOW
func main() {

    UIHandler()

    http.HandleFunc("/sensors/", restfulHandler)
    http.HandleFunc("/nodes/", restfulNodesHandler)

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

    //timePrinter("Running server on " + *addr_flag + ":" + *port_flag)
    http.ListenAndServe(*addr_flag + ":" + *port_flag, nil)
}


