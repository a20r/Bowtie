
from flask import Flask, render_template, jsonify, Response, make_response, session, request, redirect, url_for, send_from_directory
#from flask.ext.assets import Environment, Bundle
import uuid
import json
import os.path as path
import os

"""
Main server module
"""

__author__ = "Alex Wallar <aw204@st-andrews.ac.uk>"

app = Flask(__name__)

#assets = Environment(app)
#js_pre = Bundle('gather_sensor_data.js')
#assets.register('js_pre',js_pre)

#application index
@app.route('/')
def index():
	"""
	First page
	"""
	response = make_response(render_template('index.html'))
	return response

@app.route('/<cpu_id>', methods=['POST'])
def get_sensor_data(cpu_id):
	"""
	Gets data from the JavaScript
	"""
	sensor_data = json.loads(request.form['sensor_data'])
	parse_sensor_data(sensor_data, 'json_data/%s.json' % cpu_id)
	return render_template('index.html')

@app.route('/', methods=['POST', 'GET'])
def cpu_id_not_specified():
	return render_template('index.html', error="CPU identifier not specified")

@app.route('/unchecked_<cpu_id>', methods=['POST'])
def cpu_id_unchecked(cpu_id):
	try:
		os.remove('json_data/' + cpu_id + '.json')
	except IOError:
		pass
	return render_template('index.html')

@app.route('/json_data/<data_name>', methods=['GET'])
def send_sensor_data(data_name):
	"""
	Sends data to a CPU client
	"""
	file_path = 'json_data/' + data_name
	print file_path
	if not path.isfile(file_path):
		requested_data = {"error": {2: "No data for " + data_name.split('.')[0]}}
		return Response(json.dumps(requested_data), mimetype = 'application/json')
	with open(file_path, 'r+') as sensor_file:
		requested_data = sensor_file.readline()
	return Response(requested_data, mimetype='application/json')

def parse_sensor_data(sensor_data, file_path):
	"""
	Parses and saves the sensor data
	"""
	if not path.isfile(file_path):
		open(file_path, 'a').close()
	with open(file_path, 'w') as sensor_file:
		sensor_file.write(json.dumps(sensor_data))

if __name__ == '__main__':
	app.run(debug=True)