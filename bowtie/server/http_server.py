
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

"""
assets = Environment(app)
js_pre = Bundle('assets/js/bootstrap.js', 'assets/js/bootstrap.min.js', 'assets/js/html5shiv.js', 'assets/js/sensor_client.js')
css_pre = Bundle('assests/css/boostrap-responsive.min.css', 'assests/css/bootstrap.min.css')
assets.register('css_pre', css_pre)
assets.register('js_pre',js_pre)
"""

#application index
@app.route('/')
def index():
	"""
	First page
	"""
	response = make_response(render_template('index.html'))
	return response

@app.route('/about.html')
def about():
	return make_response(render_template('about.html'))

@app.route('/favicon.ico')
def favicon():
    return send_from_directory(os.path.join(app.root_path, 'static'), 'favicon.ico', mimetype='image/vnd.microsoft.icon')

@app.route('/assets/js/<js_file>.js')
def get_js(js_file):
	return send_from_directory(os.path.join(app.root_path, 'assets/js'), js_file + '.js', mimetype='text/javascript')

@app.route('/assets/css/<css_file>.css')
def get_css(css_file):
	return send_from_directory(os.path.join(app.root_path, 'assets/css'), css_file + '.css', mimetype='text/css')

@app.route('/checked/<cpu_id>/<phone_id>', methods=['POST'])
def get_sensor_data(cpu_id, phone_id):
	"""
	Gets data from the JavaScript
	"""
	sensor_data = json.loads(request.form['sensor_data'])
	parse_sensor_data(sensor_data, 'json_data/%s/%s.json' % (cpu_id, phone_id))
	return render_template('index.html')

@app.route('/', methods=['POST', 'GET'])
def cpu_id_not_specified():
	return render_template('index.html', error="CPU identifier not specified")

@app.route('/unchecked/<cpu_id>/<phone_id>', methods=['POST'])
def cpu_id_unchecked(cpu_id, phone_id):
	try:
		os.remove('json_data/%s/%s.json' % (cpu_id, phone_id))
	except OSError:
		pass
	return render_template('index.html')

@app.route('/<cpu_id>/<data_name>', methods=['GET'])
def send_single_sensor_data(cpu_id, data_name):
	"""
	Sends data to a CPU client
	"""
	file_path = 'json_data/%s/%s.json' % (cpu_id, data_name)
	if not path.isfile(file_path):
		requested_data = {"error": {2: "No data for " + data_name.split('.')[0]}}
		return Response(json.dumps(requested_data), mimetype = 'application/json')
	with open(file_path, 'r+') as sensor_file:
		requested_data = sensor_file.readline()
	return Response(requested_data, mimetype='application/json')

@app.route('/<cpu_id>/', methods=['GET'])
def send_sensor_data(cpu_id):
	"""
	Sends data to a CPU client
	"""
	file_path = 'json_data/%s/' % cpu_id
	if not path.isdir(file_path):
		requested_data = {"error": {"code": 2, "message": "No data for " + cpu_id}}
		return Response(json.dumps(requested_data), mimetype = 'application/json')
	full_data = dict()
	for it in os.walk(file_path):
		if len(it[2]) == 0:
			requested_data = {"error": {"code": 2, "message": "No data for " + cpu_id}}
			return Response(json.dumps(requested_data), mimetype = 'application/json')
		for json_file in it[2]:	
			with open(file_path + json_file, 'r+') as sensor_file:
				full_data[json_file.split('.')[0]] = json.loads(sensor_file.readline())
	return Response(json.dumps(full_data), mimetype='application/json')

def parse_sensor_data(sensor_data, file_path):
	"""
	Parses and saves the sensor data
	"""
	file_dir = '/'.join(file_path.split('/')[:-1])
	if not path.exists(file_dir):
		os.makedirs(file_dir)
	if not path.isfile(file_path):
		open(file_path, 'a').close()
	with open(file_path, 'w') as sensor_file:
		sensor_file.write(json.dumps(sensor_data))

if __name__ == '__main__':
	app.run(debug=True, host="192.168.1.95")