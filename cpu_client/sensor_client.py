
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

	def __init__(self, http_ser, droid_id):
		self.droid_id = droid_id
		self.http_ser = http_ser

	def get_data(self):
		ser = urllib2.urlopen(self.http_ser + '/' + self.droid_id + '.json')
		s_data = ser.read()
		ser.close()
		return json.loads(s_data)

if __name__ == '__main__':
	gsd = SensorDataGetter('http://127.0.0.1:5000/json_data', 'sensor_data')
	while True:
		print gsd.get_data()
		time.sleep(0.5)
