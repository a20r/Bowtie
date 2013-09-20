package main

import (

    // io pkgs
    "fmt"

    // network pkgs
    "net/http"

    // string pkgs
    "encoding/json"
)

func restfulSensorsHandler(w http.ResponseWriter, r *http.Request) {
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
    bq := MakeQueriesWithPath(r.URL.Path)
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
    bq := MakeQueriesWithPath(r.URL.Path)
    var bytes []byte
    var err error
    var marshalData interface{}
    errNum := 1

    if bq.Sensor == "" && bq.NodeId == "" {
        marshalData, err, errNum = bq.GetGroup()
    } else if bq.Sensor == ""{
        marshalData, err, errNum = bq.GetNode()
    }  else {
        marshalData, err, errNum = bq.GetSensor()
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
            Response{"Error" : JSON_ERROR, "Message": err.Error()},
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
    bq := MakeQueriesWithPath(r.URL.Path)

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
                    "Error": JSON_ERROR, 
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
                    "Error": JSON_ERROR, 
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

    bq := MakeQueriesWithPath(r.URL.Path)
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
                "Error" : JSON_ERROR, 
                "Message" : err.Error(),
            },
        )
    } else {
        w.Write(bytes) 
    }
}

func restfulMediaHandler(w http.ResponseWriter, r *http.Request) {
    timePrinter("GET\t" + r.URL.Path)

    bq := MakeQueriesWithPath(r.URL.Path)
    media, timeStamp, err := bq.GetMedia()

    if err != nil {
        fmt.Fprint(
            w,
            Response{"Error" : 1, "Message" : err.Error()},
        )
        return
    }

    var mediaType string

    switch bq.Sensor {
        case "audio":
            mediaType = "audio/base64"
        case "video":
            mediaType = "video/base64"
        default:
            fmt.Fprint(
                w,
                Response{
                    "Error" : 1, 
                    "Message" : "I have no idea what happened here",
                },
            )
            return    
    }

    fmt.Fprint(
        w,
        Response{
            "Value" : media,
            "Type" : mediaType,
            "Time" : *timeStamp,
        },
    )
}
