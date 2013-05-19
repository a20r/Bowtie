GoBowtie
========

##Purpose:
Bowtie is a sensor integration server implemented using Go, HTML5 and Javascript.
Sensory data from mobile devies are collected via the HTML5's API and pushed to a server. 
The data is stored in a unique folder defined by the CPU Id and inside the folder are JSON 
files that are defined by the Node Id. This makes it easy in any programming language on almost 
any operating system to gather the sensor data. This means that sensor data can be gathered without 
any native program running.

###Inspiration
Many small scale robotics projects suffer due to of a lack of sensors or the difficulty to
integrate sensors in the project. Many sensors come with limited documentation and require different
port specifications, different voltages, some need a pulse to be sent to them before data can 
be read, and so on. By using devices that support an internet connection and HTML5 (i.e. a 
smart phone) as sensor nodes, a layer of abstraction is removed, which leads to better overall
integration. A sample case is a small mobile robot that needs an orientation and GPS sensor.
By using Bowtie, any Android phone can be used for these sensors. The phone 
would just need to load Bowtie's specified website in the browser to gather 
the data. To integrate the data, the robot's processing unit would need to use Bowtie's
client side Python module, which would extract the sensor data.
###But why get the data through a server?
Using a main server allows for data to be easily integrated for swarm robotics. Imagine a swarm 
of robots that are dispatched from one main robot, each with their own sensor node. All of the
data could be gathered by a central processing unit on a main robot to make conclusions about the
environment. Also passing the data through a server allows for a more robust, less limited system
to grow.

##To run:
* The server: make run_server
* The CPU client: make run_client

##Requirements:
* sudo apt-get install golang

##Improvements Needed
##Using a database
* Use a database instead of storing everything in files so further analysis of the data will be available
	* Using sqlite3 in Go is easy

##Administrator website
* Make admin website that can see all of the different phone nodes and their respective CPU Ids
	* Make bootsrap password and username field
	* Have drop down views using bootsrap

##User Login 
* Allowing users to login for a certain CPU Id to see graphical data from the sensors and geolocation on maps
	* Look at Google API for maps
	* Find a way to make graphs in HTML5 so that they can be updated dynamically

###Phone Identifiers --> done
* Need to have a way to distinguish phone nodes for different robots
	* Setup a folder for each CPU identifier
	* Have a text field to enter the phone id
	* Once the CPU checkbox is unclicked, delete the JSON in the CPU id folder for the respective phone id
	* Nothing should be sending until both the phone id and the CPU id have been entered

###Launching on a server
* Just use Go!
