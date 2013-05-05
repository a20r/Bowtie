
from flask import Flask, render_template, jsonify, Response, make_response, session, request, redirect, url_for, send_from_directory
from flask.ext.assets import Environment, Bundle
import uuid
import json
import os.path as path

app = Flask(__name__)

assets = Environment(app)
js_pre = Bundle('gather_sensor_data.js')
assets.register('js_pre',js_pre)

def _generate_user_id():
    return uuid.uuid4()

#application index
@app.route('/')
def index():
	response = make_response(render_template('index.html'))
	return response

@app.route('/', methods=['POST', 'GET'])
def get_sensor_data():
	error = None
	if request.method == 'POST':
		sensor_data = json.loads(request.form['sensor_data'])
		parse_sensor_data(sensor_data, 'json_data/sensor_data.json')
		return render_template('index.html')
	else:
		error = 'WHATTT????'
	return render_template('index.html', error=error)

@app.route('/json_data/<data_name>')
def send_sensor_data(data_name):
	file_path = 'json_data/' + data_name
	print file_path
	if not path.isfile(file_path):
		return render_template('index.html', error='Requested JSON data does not exist')
	with open(file_path, 'r+') as sensor_file:
		requested_data = sensor_file.readline()
	return Response(json.dumps(requested_data), mimetype='application/json')

def parse_sensor_data(sensor_data, file_path):
	if not path.isfile(file_path):
		open(file_path, 'a').close()
	with open(file_path, 'w+b') as sensor_file:
		prev_readings_str = sensor_file.readline()
		if prev_readings_str != '':
			prev_readings = json.loads(prev_readings)
		else:
			prev_readings = dict()

		if 'location' in sensor_data.keys():
			prev_readings['location'] = sensor_data['location']
		if 'tilt' in sensor_data.keys():
			prev_readings['tilt'] = sensor_data['tilt']

		sensor_file.write(json.dumps(prev_readings))

if __name__ == '__main__':
	app.run(debug=True)