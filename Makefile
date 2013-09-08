
clean: 
	rm -r bowtie/server/audio_data
	rm -r bowtie/server/video_data
	rm -r bowtie/server/json_data

setup:
	mkdir bowtie/server/audio_data
	mkdir bowtie/server/video_data
	mkdir bowtie/server/json_data
