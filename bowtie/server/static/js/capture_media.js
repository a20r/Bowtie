var video_capturer = {
    ready: false,
    video: null,
    canvas: null,
    canvas_context: null,

    ws: null,
    ws_url: null,
    time_interval: null
}

var audio_capturer = {
    ready: false,
    recorder: null,

    ws: null,
    ws_url: null,
    time_interval: null
}

function streamVideo(stream) {
    window.stream = stream;
    window.URL = window.URL || window.webkitURL;
    var video = document.getElementById("live_stream")
    console.log("Streaming video");
    if (window.URL) {
        video.src = window.URL.createObjectURL(stream);
    } else {
        video.src = stream;
    }
    video.play();
    video_capturer.ready = true;
}

function transmitVideoToURL(video_capturer) {
    console.log("Transmitting video to url");

    return setInterval(
        function () {
            // condition that stops transmitting video stream to server
            if (video_capturer.ws.readyState == 3) { // 3 - socket is closed
                console.log("Stop video stream transmission!");
                clearInterval(timer);
            } else if (video_capturer.ready != false) {
                // draw video stream to canvas, obtain canvas data as jpg
                // then transmit to server
                video_capturer.canvas_context.drawImage(
                    video_capturer.video,
                    0,
                    0,
                    320,
                    240
                );
                var data = video_capturer.canvas.toDataURL(
                    'image/jpeg',
                    1.0
                );
                console.log("Transmitting video slice!");
                video_capturer.ws.send(data);
            }
        },
        video_capturer.time_interval
    );
}

function initMediaStream() {
    navigator.getUserMedia = navigator.getUserMedia
        || navigator.webkitGetUserMedia
        || navigator.mozGetUserMedia
        || navigator.msGetUserMedia;

    console.log("Initializing audio stream");
    navigator.getUserMedia(
        {
            audio : true,
            video : true
        },
        streamAudio,
        function(e) {
            console.log('Error! Failed to initialize audio stream:', e);
            alert('Error! Failed to initialize audio stream!');
        }
    );
}

function blobToBase64(blob) {
    var reader = new FileReader();
    reader.readAsDataURL(blob);

    reader.onload = function(reader_event) {
        var binary_string = reader_event.target.result;
        blob = btoa(binary_string);
    };

    reader.onloadend = function(reader_event) {
        // when finished encoding blob to base64
        data = reader_event.target.result;
        console.log("Transmitting audio slice!");
        audio_capturer.ws.send(data)
    }

    reader.onerror = function(reader_event) {
        console.log("FileReader Error!:" + reader_event);
        alert("Failed to encode audio binary to Base64!");
    }
}

function streamAudio(stream) {
    streamVideo(stream);
    window.AudioContext = window.AudioContext || window.webkitAudioContext;
    var audio_context = new AudioContext();
    var input_point = audio_context.createGain();

    console.log("Streaming audio");
    // Create an AudioNode from the stream.
    var real_audio_input = audio_context.createMediaStreamSource(stream);
    audio_input = real_audio_input;
    audio_input.connect(input_point);

    // Create audio recorder
    audio_capturer.recorder = new Recorder(input_point);
    audio_capturer.ready = true;
    console.log("ready!!");
}

function streamVideo(stream) {
    window.stream = stream;
    window.URL = window.URL || window.webkitURL;
    var video = document.getElementById("live_stream")
    console.log("Streaming video");
    if (window.URL) {
        video.src = window.URL.createObjectURL(stream);
    } else {
        video.src = stream;
    }
    video.play();
    video_capturer.ready = true;
}


function encodeAudio(blob) {
    return blobToBase64(blob);
}

function transmitAudioToURL(audio_capturer) {
    console.log("Transmitting audio to url");

    return setInterval(
        function() {
            if (audio_capturer.ws.readyState == 3) { // 3 - socket is closed
                console.log("Stop audio stream transmission!");
                clearInterval(timer);
            } else if (audio_capturer.ready == true) {
                setTimeout(
                    function() {
                        // stop recording
                        audio_capturer.recorder.stop();
                        audio_capturer.recorder.exportWAV(blobToBase64);

                        // start recording
                        audio_capturer.recorder.clear();
                        audio_capturer.recorder.record();
                    },
                    audio_capturer.time_interval
                );
            }
        },
        audio_capturer.time_interval
    );
}
