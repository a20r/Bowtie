all: 
	echo 'Enter run_server or run_cpu_client'

run_server:
	cd server; python sensorserver.py; cd ..

run_client:
	cd cpu_client; python sensor_client.py; cd ..
