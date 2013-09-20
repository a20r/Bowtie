package main

import (
    // sys pkgs
    "os"

    "code.google.com/p/go.net/websocket"

    // string pkgs
    "strings"
    "encoding/json"
    "encoding/base64"
)

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
    path := VIDEO_DIR + group_id + "/"

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
    path := AUDIO_DIR + group_id + "/"

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
            websocket.Message.Send(ws, "FAIL:" + err.Error())
            return
        }
        // timePrinter("ProcessSocket: got message", msg)
        websocketMsgParser(msg)
    }

    timePrinter("Finish handling websocket with wsHandler")
}
