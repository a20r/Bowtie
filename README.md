# Bowtie
Bowtie is a smartphone sensory data collector implemented in Go. To collect sensory
data, the smartphone is required to visit a HTML5 web page served by Bowtie.
Bowtie is capable of capturing:

- Accelerometer
- Gyroscope
- GPS
- Camera
- Microphone
- And more...


## Use cases
### Robotics
Our initial use case is to use Bowtie for robotics purposes. We realize there
is a large investment in terms of time and cost associated in developing the
electronics for robots, therefore we are attempting mitigate those barriers by
using what is already widely available, smartphones.

In particular a robot can utilize the phone's wide array of sensors without
having to resort to specialized sensors that are difficult to be reused for
other robots. Additionally given a swarm of robots using Bowtie, robots can
communicate and act as a collective based on data of each other.

### Crowd Sourcing
Crowd sourcing from multiple smartphones can be combined to build a
spatiotemporal view of the phenomenon of interest and also to extract important
community statistics. Given the ubiquity of mobile phones and the high density
of people in metropolitan areas, participatory sensing can achieve an
unprecedented level of coverage in both space and time for observing events of
interest in urban spaces.


## How Bowtie Works
Bowtie implements a client-server model. To make deployment simple, the client
(smartphone) is not required to install anything, the only requirement is for
the smartphone to use a HTML5 compliant web browser to visit the web page
Bowtie is currently serving.

![Client Server](images/BowtieModel.png)


## Requirements and Dependencies

- **[Go](http://golang.org/)** language
- Go's websocket library

## Usage

### Setup

	export GOPATH=$(pwd)/bowtie/server/
    export PATH=$PATH:$GOPATH/bin
    go get code.google.com/p/go.net/websocket
    go install bowtie

### Run 

    cd bowtie/server/
    bowtie -addr=<Address to run on> -port=<Port to run on>
