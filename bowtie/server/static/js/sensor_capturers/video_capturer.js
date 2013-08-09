function VideoCapturer(video, canvas, ws_url, time_interval) {
    VideoCapturer.ready_to_transmit = false;

    VideoCapturer.video = video;
    VideoCapturer.canvas = canvas;
    VideoCapturer.canvas_context = canvas.getContext('2d');

    this.ws_url = ws_url;
    this.time_interval = time_interval;
}

VideoCapturer.prototype.initVideoStream = function() {
    var video_status = true;
    navigator.getUserMedia = navigator.getUserMedia
        || navigator.webkitGetUserMedia
        || navigator.mozGetUserMedia
        || navigator.msGetUserMedia;

    console.log("Initializing video stream");
    navigator.getUserMedia(
        {video: true},
        this.streamVideo,
        function(e) {
            console.log('Error! Failed to initialize video stream:', e);
            video_status = false;
        }
    );

    return video_status;
}

VideoCapturer.prototype.streamVideo = function(stream) {
    window.stream = stream;
    window.URL = window.URL || window.webkitURL;

    console.log("Streaming video");
    if (window.URL) {
        VideoCapturer.video.src = window.URL.createObjectURL(stream);
    } else {
        VideoCapturer.video.src = stream;
    }
    VideoCapturer.video.play();
    VideoCapturer.ready_to_transmit = true;
}

VideoCapturer.prototype.transmitVideoToURL = function() {
    console.log("Transmitting video to url");
    var ws_client = new WebSocketClient();
    var ws = ws_client.init(this.ws_url);

    try {
    var timer = setInterval(
        function () {
            // condition that stops transmitting video stream to server
            if (ws.readyState == 3) { // 3 - websocket is closed
                console.log("Stop video stream transmission!");
                clearInterval(timer);
            } else if (VideoCapturer.ready_to_transmit != false) {
                // draw video stream to canvas, obtain canvas data as jpg
                // then transmit to server
                VideoCapturer.canvas_context.drawImage(VideoCapturer.video, 0, 0, 320, 240);
                var data = this.canvas.toDataURL('image/jpeg', 1.0);
                console.log("Transmitting: [" + data + "]");
                ws.send(data);
            }
        },
        this.time_interval
    );
    } catch (error) {
        console.log(error);
    }

}
