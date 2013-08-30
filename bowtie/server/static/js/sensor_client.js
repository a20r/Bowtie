
////////////////////////////////////////////////
// sensor_client.js
//
// Main JavaScript code for the Bowtie website.
// In charge of retrieving the sensor data
// and sending it to the server.
//
// Improvements:
//    Need to make sure that the sensor data
//    stops sending once the page has been
//    changed. For instance when the user
//    clicks on the about page.
//
////////////////////////////////////////////////

var audioInterval = undefined;
var videoInterval = undefined;
var sendingInterval = undefined;
var waitTime = 700; // ms

//var ws_url = "ws://localhost:8000/websocket/";
var ws_url = "ws://82.196.12.41/websocket/";

// Two functions that need ro run
// in order for the code to work properly
window.onload = function() {
  getLocation();
  getSensorData();
  getUserMediaData();
}

// Sets the window exiting function
window.onbeforeunload = on_exit;

function hasGetUserMedia() {
    // Note: Opera is unprefixed.
    return !!(navigator.getUserMedia || navigator.webkitGetUserMedia ||
    navigator.mozGetUserMedia || navigator.msGetUserMedia);
}

function getUserMediaData() {
    var video = document.getElementById('live_stream');
    var canvas = document.getElementById('vid_img');

    // change the IP for different testing scenarios!!

    var time_interval = 1000;

    // Setup
    // Video Capturer
    video_capturer.video = video;
    video_capturer.canvas = canvas;
    video_capturer.canvas_context = canvas.getContext("2d");
    video_capturer.ws_url = ws_url;
    video_capturer.time_interval = time_interval;

    // Audio Capturer
    audio_capturer.ws_url = ws_url;
    audio_capturer.time_interval = time_interval;

    if (!hasGetUserMedia()) {
        alert('Error! getUserMedia() is not supported in your browser!');
    }
}

// Occurs when somebody clicks the alert
// close button
function warning_closed() {
  $("#alert_msg").css("display", "none");
}

// Function fires once the page is closed
function on_exit() {
    if (
            $("#node_id").val() != "" &&
            $("#group_id").val() != ""
    ) {
        $.ajax(
            {
                type: 'POST',
                url : (
                    '/unchecked/' +
                    $("#group_id").val() + '/' +
                    $("#node_id").val()
                )
            }
        );
  }
}

// Toggles whether the data is being shown to the user
// and whether it gets sent to the server
function toggle_readonly() {

    var node_id_box = document.getElementById("node_id");
    var group_id_box = document.getElementById("group_id");
    if(node_id_box.hasAttribute('readonly')) {

        $("#accelerometer-chart").css("display", "none");

        try {
            clearInterval(sendingInterval);
            clearInterval(audioInterval);
            clearInterval(videoInterval);
        } catch (err) {}

        node_id_box.removeAttribute('readonly');
        group_id_box.removeAttribute('readonly');

        $("#sub_button").html("Start sensing");
        $("#sub_button").attr(
            "class",
            "btn btn-large btn-success"
        );

        $("#sensor_table").css(
            "display",
            "none"
        );

    // deletes the data from the server
        if (
                node_id_box.value != "" &&
                group_id_box.value != ""
        ) {
            $.ajax(
                {
                    type : 'POST',
                    url : (
                        '/unchecked/' +
                        $("#group_id").val() + '/' +
                        $("#node_id").val()
                    )
                }
            );
        }
    } else {
        if (
                node_id_box.value != "" &&
                group_id_box.value != ""
        ) {
            sendingInterval = window.setInterval(sendAjax, waitTime);

            if (hasGetUserMedia()) {
                var ws_client = new WebSocketClient();
                var ws = ws_client.init(ws_url);
                video_capturer.ws = ws;
                audio_capturer.ws = ws;

                initMediaStream();
            }

            // Transmitting video and audio
            transmission_details.group_id = $("#group_id").val();
            transmission_details.node_id = $("#node_id").val();
            audioInterval = transmitAudioToURL(audio_capturer);
            videoInterval = transmitVideoToURL(video_capturer);

            group_id_box.setAttribute('readonly', 'readonly');
            node_id_box.setAttribute('readonly', 'readonly');
            $("#accelerometer-chart").empty();
            $("#accelerometer-chart").css("display", "block");
            realtime_demo();

            $("#sub_button").html("Stop sensing");
            $("#sub_button").attr(
                "class",
                "btn btn-large btn-primary btn-danger"
            );

            $("#sensor_table").css("display", "block");
            $("#alert_msg").css("display", "none");
        } else {
            $("#alert_msg").css("display", "block");
        }
    }
}

// Sets the event handler for the orientation sensor.
// If the data cannot be gathered, messages will be
// shown and the checkboxes will be disabled
function getSensorData() {
    if(window.DeviceOrientationEvent) {
        window.addEventListener(
            'deviceorientation',
            orientationEventHandler,
            false
        );
    } else {
        $("#doTiltLR").html("Not supported");
        $("#doTiltLRCheckbox").prop("checked", false);
        $("#doTiltLRCheckbox").prop("disabled", true);
        $("#doTiltFB").html("Not supported");
        $("#doTiltFBCheckbox").prop("checked", false);
        $("#doTiltFBCheckbox").prop("disabled", true);
        $("#doDirection").html("Not supported");
        $("#doDirectionCheckbox").prop("checked", false);
        $("#doDirectionCheckbox").prop("disabled", true);
    }
}

// Handles the orientation data
function orientationEventHandler(eventData) {
    var tiltLR = eventData.gamma;
    var tiltFB = eventData.beta;
    var dir = eventData.alpha;

    $("#doTiltLR").html(
        Math.round(tiltLR)
    );

    $("#doTiltFB").html(
        Math.round(tiltFB)
    );

    $("#doDirection").html(
        Math.round(dir)
    );
}

// Sets the locaction event handlers
function getLocation() {
    if (navigator.geolocation) {
        var timeoutVal = 6000;

        var extraGeoParam = {
            enableHighAccuracy: true,
            timeout: timeoutVal,
            maximumAge: 0
        };

        navigator.geolocation.watchPosition(
            devicePositionHandler,
            positionError,
            extraGeoParam
        );
    } else {
        alert("Geolocation is not supported by this browser");
    }
}

// Fired if getting the position raises an error
function positionError (position) {
    $("#latPos").html("Not supported");
    $("#latPosCheckbox").prop("checked", false);
    $("#latPosCheckbox").prop("disabled", true);
    $("#longPos").html("Not supported");
    $("#longPosCheckbox").prop("checked", false);
    $("#longPosCheckbox").prop("disabled", true);
}

// Writes the position of the coordinates onto the
// HTML page
function devicePositionHandler(position) {
    $("#latPos").html(position.coords.latitude);
    $("#longPos").html(position.coords.longitude);
    $("#latPosCheckbox").prop("disabled", false);
    $("#longPosCheckbox").prop("disabled", false);
}

// Sends the sensory data to the server via
// Ajax if the group_id and the node_id have
// been entered
var readyToSend = true;
function sendAjax() {
    if (!readyToSend) {
        return
    }

    tiltLR = getIfValid("doTiltLR");
    tiltFB = getIfValid("doTiltFB");
    dir = getIfValid("doDirection");

    lat = getIfValid("latPos");
    lon = getIfValid("longPos");
    if (
            lat == null &&
            lon == null &&
            $("#latPos").html() == "Not supported"
    ) {
        error_data = {code: 1, message: "Position error"}
    } else {
        error_data = {code: 0, message: "No error"}
    }
    if (
            $("#group_id").attr('readonly') != undefined
    ) {
        readyToSend = false;
        // RESTful POST
        $.ajax(
            {
                type : 'PUT',

                url : (
                    '/sensors/' +
                    $("#group_id").val() + '/' +
                    $("#node_id").val()
                ),

                data : {
                    sensorData : JSON.stringify(
                        {
                            tilt_horizontal : {
                                value : tiltLR,
                                type : "integer",
                                time : new Date().toJSON()
                            },

                            tilt_vertical : {
                                value : tiltFB,
                                type : "integer",
                                time : new Date().toJSON()
                            },

                            orientation : {
                                value : dir,
                                type : "integer",
                                time : new Date().toJSON()
                            },

                            latitude : {
                                value : lat,
                                type : "float",
                                time : new Date().toJSON()
                            },

                            longitude : {
                                value : lon,
                                type : "float",
                                time : new Date().toJSON()
                            },
                        }
                    )
                },

                success : function () {
                    console.log("posted");
                    readyToSend = true;
                }
            }
        );

        // // Old post
        // $.ajax(
        //     {
        //         type : 'POST',

        //         url : (
        //             '/checked/' +
        //             $("#group_id").val() + '/' +
        //             $("#node_id").val()
        //         ),

        //         data: {
        //             sensor_data: JSON.stringify(
        //                 {
        //                     orientation: {
        //                         tilt_horizontal: tiltLR,
        //                         tilt_vertical: tiltFB,
        //                         direction: dir
        //                     },

        //                     location: {
        //                         latitude: lat,
        //                         longitude: lon
        //                     },

        //                     error: error_data
        //                 }
        //             )
        //         }
        //     }
        // );
    }
}

// Checks if the data is valid before sending
// it out
function getIfValid(element_id) {
    if (
            $("#" + element_id).html() == "Not supported" ||
            $("#" + element_id).html() == "" ||
            !$("#" + element_id + "Checkbox").prop("checked")
    ) {
        return null;
    } else {
        return parseFloat(
            $("#" + element_id).html()
        );
    }
}
