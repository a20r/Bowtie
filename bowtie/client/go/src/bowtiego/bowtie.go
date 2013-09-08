
package bowtiego

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

}
