PyBowtie
========

##Purpose:
* The goal is to create a seamless integration between phone sensor data and robotics applications. 
* Small scale robotics has suffered because of a lack of sensor data due to cost and CPU pin compatability
* Bowtie offers a cheaper solution where the developer only needs a smart phone and an internet connection on the robot

##To run:
* The server: make run_server
* The CPU client: make run_client

##Requirements:
* Pip:
	* easy_install pip
* Flask:
	* pip install flask

##Improvements Needed
##Using a database
* Use a database instead of storing everything in files so further analysis of the data will be available
	* Using sqlite3 in Python is easy 

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
* Need to put the server onto an actual server for testing
	* Google app engine?
	* School server using a virtual env?
	* Nathan's server?
