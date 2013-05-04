
from flask import Flask, render_template, jsonify, Response, make_response, session, request, redirect, url_for, send_from_directory
from flask.ext.assets import Environment, Bundle
import uuid
import json

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
		return render_template('index.html')
	else:
		error = 'WHATTT????'
	# the code below this is executed if the request method
	# was GET or the credentials were invalid
	return render_template('index.html', error=error)

if __name__ == '__main__':
	app.run(debug=True)