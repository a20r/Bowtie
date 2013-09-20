package main

import (
    // sys pkgs
    "os"

    // io pkgs
    //"fmt"
    "io/ioutil"

    // string pkgs
    "strings"
    "encoding/json"
    "encoding/base64"
    "errors"

    //ADTs
    "time"
)

const (
    NO_ERROR = 0
    EXISTS_ERROR = 1
    JSON_ERROR = 2
    READ_ERROR = 3
    DIR_ERROR = 4
)

const (
    JSON_DIR = "./json_data/"
    AUDIO_DIR = "./audio_data/"
    VIDEO_DIR = "./video_data/"
)

// using a more defined type for the restful API
type SensorData struct {
    Value interface{} `value` 
    Type string `type`
    Time string `time`
}

type NodeData map[string]SensorData

type GroupData map[string]NodeData

// for easier querying
type BowtieQueries struct {
    GroupId string
    NodeId string
    Sensor string
}

func (bq BowtieQueries) GroupExists() bool {

    if bq.GroupId == "" {
        return false
    }

    files, err := ioutil.ReadDir(
        "./json_data/" + bq.GroupId,
    )

    return err == nil && len(files) > 0
}

func (bq BowtieQueries) NodeExists() bool {

    if bq.NodeId == "" {
        return false
    }

    _, readErr := ioutil.ReadFile(
        JSON_DIR + 
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

    sDataMap, err, _ := bq.GetNode()

    if err != nil {
        sDataMap = make(NodeData)
    }

    sDataMap[bq.Sensor] = sData

    path := JSON_DIR + bq.GroupId + "/"
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

    path := JSON_DIR + bq.GroupId + "/"
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

func (bq BowtieQueries) GetSensor() (*SensorData, error, int) {
    node, err, errNum := bq.GetNode()
    if err != nil {
        return nil, err, errNum
    }

    if _, exists := node[bq.Sensor]; !exists {
        return nil, errors.New("Sensor does not exist"), EXISTS_ERROR
    }

    return &SensorData{
        node[bq.Sensor].Value,
        node[bq.Sensor].Type,
        node[bq.Sensor].Time,
    }, nil, NO_ERROR
}

func (bq BowtieQueries) GetNode() (NodeData, error, int) {

    var sDataMap NodeData

    fileBytes, readErr := ioutil.ReadFile(
        JSON_DIR + 
        bq.GroupId + "/" + 
        bq.NodeId + ".json",
    )

    if readErr != nil {
        return nil, errors.New(
            "Unable to get node --> " + 
            readErr.Error(),
        ), READ_ERROR
    }

    jsonErr := json.Unmarshal(fileBytes, &sDataMap)

    if jsonErr != nil {
        return nil, errors.New(
            "Unable to get node --> " + 
            jsonErr.Error(),
        ), JSON_ERROR
    }

    return sDataMap, nil, NO_ERROR
}

func (bq BowtieQueries) GetGroup() (GroupData, error, int) {
    if !bq.GroupExists() {
        return nil, errors.New("Group does not yet exist"), EXISTS_ERROR
    }

    files, err := ioutil.ReadDir(JSON_DIR + bq.GroupId)

    if err != nil {
        return nil, errors.New(
            "Unable to get group --> " +
            err.Error(),
        ), DIR_ERROR
    }

    group := make(GroupData)

    for _, file := range files {
        if file.Name() == ".DS_Store" {
            continue
        }
        var sDataMap NodeData
        nodeId := strings.Split(file.Name(), ".")[0]
        fileBytes, readErr := ioutil.ReadFile(
            JSON_DIR + 
            bq.GroupId + "/" + 
            nodeId + ".json",
        )

        if readErr != nil {
            return nil, errors.New(
                "Unable to get group --> " + 
                readErr.Error(),
            ), READ_ERROR
        }

        jsonErr := json.Unmarshal(fileBytes, &sDataMap)

        if jsonErr != nil {
            return nil, errors.New(
                "Unable to get group --> " + 
                jsonErr.Error(),
            ), JSON_ERROR
        }

        group[nodeId] = sDataMap
    }
    return group, nil, NO_ERROR
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

func (bq BowtieQueries) GetMedia() ([]byte, *time.Time, error) {
    var path string
    var extension string

    if !bq.GroupExists() {
        return nil, nil, errors.New(
            "Group does not exist",
        )
    }

    if !bq.NodeExists() {
        return nil, nil, errors.New(
            "Node does not exist",
        )
    }

    switch bq.Sensor {
        case "audio":
            path = AUDIO_DIR
            extension = ".wav"
        case "video":
            path = VIDEO_DIR
            extension = ".jpg"
        default:
            return nil, nil, errors.New(
                "Unsupported media type --> " + bq.Sensor,
            )
    }

    media, readErr := ioutil.ReadFile(
        path + 
        bq.GroupId + "/" + 
        bq.NodeId + extension,
    )

    file, openErr := os.Open(
        path + 
        bq.GroupId + "/" + 
        bq.NodeId + extension,
    )

    if readErr != nil {
        return nil, nil, readErr
    }

    if openErr != nil {
        return nil, nil, openErr
    }

    stat, statErr := file.Stat()

    if statErr != nil {
        return nil, nil, statErr
    }

    mediaEncoded := base64.StdEncoding.EncodeToString(
        media,
    )

    mTime := stat.ModTime()

    return []byte(mediaEncoded), &mTime, nil
}

func (bq BowtieQueries) DeleteGroup() error {
    err := os.Remove(
        JSON_DIR + 
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
        JSON_DIR + 
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

    node, err, _ := bq.GetNode()

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

func MakeQueriesWithPath(
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