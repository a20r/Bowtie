
import urllib2
import json
#import threading
import time
import uuid

"""
Module will be used on the CPU side to get the 
data from the website
"""

__author__ = "Santa>"

class Bowtie:

	def __init__(self, http_ser, port, droid_id = None):
		"""
		For the droid_id I would like to use the MAC address of the 
		CPU client.
		"""
		if droid_id == None:
			self.droid_id = str(uuid.getnode())
		else:
			self.droid_id = str(droid_id)
		self.http_ser = http_ser
		self.port = port

		print 'Using identifier:', self.droid_id

	def get_data_url(self):
		return self.http_ser + ':' + str(self.port) + '/json_data/' + self.droid_id + '.json'

	def get_data(self):
		try:
			ser = urllib2.urlopen(self.get_data_url())
			s_data = ser.read()
			ser.close()
			return json.loads(s_data)
		except urllib2.URLError:
			return {"error": {3: "No connection to host"}}
		except ValueError:
			return {"error": {4: "JSON data could not be parsed"}}

if __name__ == '__main__':
	gsd = Bowtie('http://127.0.0.1', 5000, 1234)
	while True:
		print gsd.get_data()
		time.sleep(0.1)
