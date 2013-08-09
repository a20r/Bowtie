function AudioCapturer() {

}

AudioCapturer.prototype.initAudioStream = function() {
    var audio_status = true;

    navigator.getUserMedia(
        {audio: true},
        streamAudio,
        function(e) {
            console.log('Error! Failed to initialize audio stream:', e);
            audio_status = false;
        }
    );

    return audio_status;
}

AudioCapturer.prototype.streamAudio = function(stream) {
    window.AudioContext = window.AudioContext || window.webkitAudioContext;

    var audio_context = new AudioContext();
    var input_point = audio_context.createGain();

    // Create an AudioNode from the stream.
    var real_audio_input = audio_context.createMediaStreamSource(stream);
    audio_input = real_audio_input;
    audio_input.connect(input_point);

    // Create audio recorder
    audio_recorder = new Recorder(input_point);

    // Record and transmit
    var interval = 1000
    intervalRecordAndTransmit(audio_recorder, interval);
}

AudioCapturer.prototype.blobToBase64 = function(blob) {
    var reader = new FileReader();
    reader.readAsDataURL(blob);

    reader.onload = function(reader_event) {
        var binary_string = reader_event.target.result;
        blob = btoa(binary_string);
    };

    reader.onloadend = function(reader_event) {
        // when finished encoding blob to base64
        return reader_event.target.result;
    }

    reader.onerror = function(reader_event) {
        console.log("FileReader Error!:" + reader_event);
        return null;
    }
}

AudioCapturer.prototype.transmitAudioStreamToURL = function(blob) {
    setTimeout(
        function() {
            return blobToBase64(blob);
        },
        1000
    );
}

AudioCapturer.prototype.intervalRecord = function(audio_recorder, interval) {
    var ws = openWebSocket(ws_url);

    audio_recorder.clear();
    audio_recorder.record();
    console.log("Start recording!");

    setTimeout(
        function() {
            console.log("Stop recording!");
            audio_recorder.stop();
            var data = audio_recorder.exportWAV(transmitAudioStreamToURL);
            console.log("DATA : " +  data);
        },
        interval
    );
}

