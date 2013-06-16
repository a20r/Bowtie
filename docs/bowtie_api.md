# Bowtie RESTFUL API
## Sensors
Below are methods to obtain sensory data from nodes:
### Get Sensory Data

    `GET sensors/temperature/:node_id`
    `GET sensors/geo_location/:node_id`
    `GET sensors/gyroscope/:node_id`
    `GET sensors/accelerometer/:node_id`
    `GET sensors/device_orientation/:node_id`
    `GET sensors/camera/:node_id`
    `GET sensors/microphone/:node_id`

#### Parameters
- **node_id** (Required): Node id of the node you wish to obtain data from.

### Post Sensory Data

    `POST sensors/temperature/`
    `POST sensors/geo_location/`
    `POST sensors/gyroscope/`
    `POST sensors/accelerometer/`
    `POST sensors/device_orientation/`
    `POST sensors/camera/`
    `POST sensors/microphone/`

#### Parameters
No parameters are needed since the server will be using the IP for identifying
the device.


## OAuth
OAuth is an open standard for authorization. OAuth provides a method for
clients to access server resources on behalf of a resource owner (such as a
different client or an end-user). It also provides a process for end-users to
authorize third-party access to their server resources without sharing their
credentials (typically, a username and password pair), using user-agent
redirections. Below are methods that provide those features:

### Authorize 

    `GET oauth/authorize`

Allows a node to use an OAuth Request Token to request user authorization.

### Access Token

    `POST oauth/access_token`

Allows a node to exchange the OAuth Request Token for an OAuth Access Token.

### Request Token

    `POST oauth/request_token`

Allows a node to obtain an OAuth Request Token to request user authorization.

### Invalidate Token

    `POST oauth/invalidate_token`

Allows a registered application to revoke an issued OAuth token by presenting
its client credentials. Once a token has been invalidated, new creation
attempts will yield a different Bearer Token and usage of the invalidated token
will no longer be allowed.

#### Parameters
No parameters are needed since the server will be using the IP for identifying
the unique device.
