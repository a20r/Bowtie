all:
	@echo 'options:'
	@echo "\t run_server"
	@echo "\t run_cpu_client"
	@echo ''

dependency_check:
	@python DEPENDENCIES; if [ $$? -eq 1 ]; then exit -1; fi;

run_server:
	@echo "Running Bowtie server ..."
	cd bowtie/server; go run bowtie_server.go; cd ../..

run_client: dependency_check
	@echo "Running Bowtie client ..."
	cd bowtie/cpu_client; python sensor_client.py; cd ../..
