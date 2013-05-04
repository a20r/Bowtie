
from flask import Flask, render_template, jsonify, Response, make_response, session, request, redirect, url_for, send_from_directory
import uuid
import json
app = Flask(__name__)

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
		print request.form['sensor_data']
		return request.form['sensor_data']
	else:
		error = 'WHATTT????'
	# the code below this is executed if the request method
	# was GET or the credentials were invalid
	return render_template('index.html', error=error)

if __name__ == '__main__':
	app.run(debug=True)