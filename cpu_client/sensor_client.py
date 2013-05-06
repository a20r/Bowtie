
import urllib2
import json
#import threading
import time
import uuid

"""
Module will be used on the CPU side to get the 
data from the website
"""

__author__ = "Alex Wallar <aw204@st-andrews.ac.uk>"

class SensorDataGetter:

	def __init__(self, http_ser, port, droid_id = None):
		"""
		For the droid_id I would like to use the MAC address of the 
		CPU client.
		"""
		if droid_id == None:
			self.droid_id = str(uuid.getnode())
		else:
			self.droid_id = droid_id
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
		except (ValueError, urllib2.URLError):
			return dict()

if __name__ == '__main__':
	gsd = SensorDataGetter('http://127.0.0.1', 5000)
	while True:
		print gsd.get_data()
		time.sleep(0.5)
