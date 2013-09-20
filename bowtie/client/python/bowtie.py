
import json
import urllib2
import time

class BowtieClient:

	def __init__(self, url="http://www.bowtie.mobi/"):

		self.url = url

		if self.url[:4] != "http":
			self.url = "http://" + self.url

		if self.url[-1] != "/":
			self.url += "/"

	def getGroup(self, groupId):
		jsonStr = urllib2.urlopen(
			self.url + "sensors/" + 
			groupId
		).read()
		return json.loads(jsonStr)

	def getNode(self, groupId, nodeId):
		jsonStr = urllib2.urlopen(
			self.url + "sensors/" + 
			groupId + "/" + 
			nodeId
		).read()
		return json.loads(jsonStr)

	def getSensor(self, groupId, nodeId, sensor):
		jsonStr = urllib2.urlopen(
			self.url + "sensors/" + 
			groupId + "/" + 
			nodeId + "/" + 
			sensor
		).read()
		return json.loads(jsonStr)
