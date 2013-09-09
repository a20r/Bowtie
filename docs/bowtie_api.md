# Bowtie RESTful API
## Sensors
Below are methods to obtain sensory data from nodes:
### Get Sensory Data

    `GET sensors/:group_id/:node_id/latitude/`
    `GET sensors/:group_id/:node_id/longitude/`
    `GET sensors/:group_id/:node_id/orientation/`
    `GET sensors/:group_id/:node_id/tilt_vertical/`
    `GET sensors/:group_id/:node_id/tilt_horizontal/`
    `GET sensors/:group_id/:node_id/:custom_sensor`

#### Return Structure

    {
        Value : `JSON Object`,
        Type : `String`,
        Time : `String`
    }

### Get Node Data

    `GET sensors/:group_id/:node_id/`

#### Return Structure

    {
        `Sensor Name` : {
            Value : `JSON Object`,
            Type : `String`,
            Time : `String`
        }
    } 

### Get Group Data

    `GET sensors/:group_id/`

### Return Structure

    {
        `Node Id` : {
            `Sensor Name` : {
                Value : `JSON Object`,
                Type : `String`,
                Time : `String`
            }
        }
    }

### Get Media Data

    `GET media/:group_id/:node_id/:media_type`

#### Return Structure

    {
        Value : `Base 64 String`,
        Type : `:media_type/base64`,
        Time : `String`
    }

### Post Sensory Data

    `POST sensors/:group_id/:node_id/latitude/`
    `POST sensors/:group_id/:node_id/longitude/`
    `POST sensors/:group_id/:node_id/orientation/`
    `POST sensors/:group_id/:node_id/tilt_vertical/`
    `POST sensors/:group_id/:node_id/tilt_horizontal/`
    `POST sensors/:group_id/:node_id/:custom_sensor`

#### Post Structure

    {
        Value : `JSON Object`,
        Type : `String`,
        Time : `String`
    }


### Post Node Data

    `POST sensors/:group_id/:node_id/`

#### Post Structure

    {
        `Sensor Name` : {
            Value : `JSON Object`,
            Type : `String`,
            Time : `String`
        }
    } 

### Delete Sensory Data

    `DELETE sensors/:group_id/:node_id/latitude/`
    `DELETE sensors/:group_id/:node_id/longitude/`
    `DELETE sensors/:group_id/:node_id/orientation/`
    `DELETE sensors/:group_id/:node_id/tilt_vertical/`
    `DELETE sensors/:group_id/:node_id/tilt_horizontal/`
    `DELETE sensors/:group_id/:node_id/:custom_sensor`

### Delete Node Data

    `DELETE sensors/:group_id/:node_id/`

### Delete Group Data

    `DELETE sensors/:group_id/`

### Get the List of Nodes for a Group

    `GET nodes/:group_id`

#### Return Structure

    [ String ]

#### Parameters
- **group_id**: Group id of the group you wish to act upon.
- **node_id**: Node id of the node you wish to act upon
- **media_type**: Either `audio` or `video`. Represents the type of media you would like to obtain from the server.
- **custom_sensor**: This is a custom sensor name. The other sensors are provided by our web application. You can create a custom sensor by posting to a custom sensor's name. You can operate on this type of sensors as you would any of the sensors provided by the web app.
