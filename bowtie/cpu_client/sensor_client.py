
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
		Creates a Bowtie instance for a certain CPU id
		"""
		if cpu_id == None:
			self.cpu_id = str(uuid.getnode())
		else:
			self.cpu_id = str(cpu_id)
		self.http_ser = http_ser
		self.port = port

		print 'Using identifier:', self.cpu_id

	def get_data_url(self):
		"""
		Gets the URL for the cpu_id
		"""
		return self.http_ser + ':' + str(self.port) + '/' + self.cpu_id + '/'

	def get_data(self, phone_id = ""):
		"""
		Gets the node data. If the phone_id is specified,
		only data for that node is sent otherwise, the data
		for all of the nodes is sent
		"""
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
	"""
	Currently using that weird IP address for local 
	testing 
	"""
	gsd = Bowtie('http://192.168.1.95', 5000, 1234)
	while True:
		print gsd.get_data()
		time.sleep(0.1)
