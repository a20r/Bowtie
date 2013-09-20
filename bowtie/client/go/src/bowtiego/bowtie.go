
package bowtiego

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
)

const(
	SENSORS = "sensors"
	MEDIA = "media"
	NODES = "nodes"
)

type SensorData struct {
    Value interface{} `value` 
    Type string `type`
    Time string `time`
}

type NodeData map[string]SensorData

type GroupData map[string]NodeData

type BowtieServer struct {
	URL string
}

func (bs BowtieServer) GetSensor(groupId, nodeId, sensor string) (*SensorData, error) {
	resp, err := http.Get(
		bq.URL + SENSORS + "/" + 
		groupId + "/" + 
		nodeId + "/" + sensor,
	)

	if err != nil {
		return nil, err
	}

	body, readErr := ioutil.ReadAll(resp.Body)

	if readErr != nil {
		return nil, readErr
	}

	var sData SensorData

	jsonErr := json.Unmarshal(body, &sData)

	if jsonErr != nil {
		return nil, jsonErr
	}

	return sData, nil
}
