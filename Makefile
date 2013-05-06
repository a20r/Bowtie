all: 
	echo 'Enter run_server or run_cpu_client'

run_server:
	cd bowtie/server; python sensorserver.py; cd ..

run_client:
	cd bowtie/cpu_client; python sensor_client.py; cd ..
