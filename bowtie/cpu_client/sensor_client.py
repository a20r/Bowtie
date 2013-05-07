
import urllib2
import json
#import threading
import time
import uuid

"""
Module will be used on the CPU side to get the 
data from the website

Improvements:
	Need to have more of an interface to see which phone_ids exist
	and to disable getting data from a certain sensor from a certain
	phone_id
"""

__author__ = "Alexander Wallar <aw204@st-andrews.ac.uk>"

class Bowtie:

	def __init__(self, http_ser, port, cpu_id = None):
		"""
		For the cpu_id I would like to use the MAC address of the 
		CPU client.
		"""
		if cpu_id == None:
			self.cpu_id = str(uuid.getnode())
		else:
			self.cpu_id = str(cpu_id)
		self.http_ser = http_ser
		self.port = port

		print 'Using identifier:', self.cpu_id

	def get_data_url(self):
		return self.http_ser + ':' + str(self.port) + '/json_data/' + self.cpu_id + '/'

	def get_data(self, phone_id = ""):
		if phone_id != "":
			phone_id = phone_id + '.json'
		try:
			ser = urllib2.urlopen(self.get_data_url() + phone_id)
			s_data = ser.read()
			ser.close()
			return json.loads(s_data)
		except urllib2.URLError:
			return {"error": {"code": 3, "message": "No connection to host"}}
		except ValueError:
			return {"error": {"code": 4, "message": "JSON data could not be parsed"}}

if __name__ == '__main__':
	gsd = Bowtie('http://127.0.0.1', 5000, 1234)
	while True:
		print gsd.get_data()
		time.sleep(0.1)
