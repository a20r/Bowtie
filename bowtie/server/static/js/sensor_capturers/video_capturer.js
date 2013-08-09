var video_capturer = {
    ready: false,
    video: null,
    canvas: null,
    canvas_context: null,

    ws_url: null,
    time_interval: null
}

function initVideoStream() {
    navigator.getUserMedia = navigator.getUserMedia
        || navigator.webkitGetUserMedia
        || navigator.mozGetUserMedia
        || navigator.msGetUserMedia;

    console.log("Initializing video stream");
    navigator.getUserMedia(
        {video: true},
        streamVideo,
        function(e) {
            console.log('Error! Failed to initialize video stream:', e);
            alert('Error! Failed to initialize video stream!');
        }
    );
}

function streamVideo(stream) {
    window.stream = stream;
    window.URL = window.URL || window.webkitURL;

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
    var ws_client = new WebSocketClient();
    var ws = ws_client.init(video_capturer.ws_url);

    try {
        var timer = setInterval(
            function () {
                // condition that stops transmitting video stream to server
                if (ws.readyState == 3) { // 3 - websocket is closed
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
                    console.log("Transmitting: [" + data + "]");
                    ws.send(data);
                }
            },
            video_capturer.time_interval
        );
    } catch (error) {
        console.log(error);
    }

}
