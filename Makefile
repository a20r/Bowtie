all:
	echo 'Enter run_server or run_cpu_client'

dependency_check:
	@python DEPENDENCIES; if [ $$? -eq 1 ]; then exit -1; fi;

run_server:
	cd bowtie/server; go run httpgo.go; cd ..

run_client: dependency_check
	cd bowtie/cpu_client; python sensor_client.py; cd ..
