
import urllib2
import json
#import threading
import time

"""
Module will be used on the CPU side to get the 
data from the website
"""

__author__ = "Alex Wallar <aw204@st-andrews.ac.uk>"

class SensorDataGetter:

	def __init__(self, http_ser, port, droid_id):
		"""
		For the droid_id I would like to use the MAC address of the 
		CPU client.
		"""
		self.droid_id = droid_id
		self.http_ser = http_ser
		self.port = port

	def get_data_url(self):
		return self.http_ser + ':' + str(self.port) + '/json_data/' + self.droid_id + '.json'

	def get_data(self):
		ser = urllib2.urlopen(self.get_data_url())
		s_data = ser.read()
		ser.close()
		try:
			return json.loads(s_data)
		except ValueError:
			return dict()

if __name__ == '__main__':
	gsd = SensorDataGetter('http://127.0.0.1', 5000, 'sensor_data')
	while True:
		print gsd.get_data()
		time.sleep(0.5)
