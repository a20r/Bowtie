all:
	echo 'Enter run_server or run_cpu_client'

dependency_check:
	@python DEPENDENCIES; if [ $$? -eq 1 ]; then exit -1; fi;

run_server: dependency_check
	cd bowtie/server; python http_server.py; cd ..

run_client: dependency_check
	cd bowtie/cpu_client; python sensor_client.py; cd ..
