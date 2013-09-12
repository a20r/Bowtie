
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
			self.url + "sensors/" + groupId
		).read()
		return json.loads(jsonStr)

def testLatency(iterations, groupId):
	bc = BowtieClient()
	avg = 0
	for i in range(iterations):
		timeSent = time.time()
		a = bc.getGroup(groupId)
		timeRec = time.time()
		avg += (timeRec - timeSent)
		print str(i) + "\t: " + str(timeRec - timeSent)
	return avg / iterations
