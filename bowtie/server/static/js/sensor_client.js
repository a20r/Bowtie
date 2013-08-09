
////////////////////////////////////////////////
// sensor_client.js
//
// Main JavaScript code for the Bowtie website.
// In charge of retrieving the sensor data
// and sending it to the server. 
// Improvements: 
//    Need to make sure that the sensor data
//    stops sending once the page has been 
//    changed. For instance when the user
//    clicks on the about page.
//
////////////////////////////////////////////////

// Two functions that need ro run
// in order for the code to work properly
window.onload = function() {
  getLocation();
  getSensorData(200);
}

// Sets the window exiting function
window.onbeforeunload = on_exit;

// Occurs when somebody clicks the alert
// close button
function warning_closed() {
  $("#alert_msg").css("display", "none");
}

// Function fires once the page is closed
function on_exit() {
    if (
            $("#phone_id").val() != "" && 
            $("#cpu_id").val() != ""
    ) {
        $.ajax(
            {
                type: 'POST',
                url : (
                    '/unchecked/' + 
                    $("#cpu_id").val() + '/' + 
                    $("#phone_id").val()
                )
            }
        );
  }
}

// Toggles whether the data is being shown to the user
// and whether it gets sent to the server
var sendingInterval;
function toggle_readonly() {

    var phone_id_box = document.getElementById("phone_id");
    var cpu_id_box = document.getElementById("cpu_id");
    if(phone_id_box.hasAttribute('readonly')) { 

        $("#accelerometer-chart").css("display", "none");

        try {
            window.clearInterval(sendingInterval);
        } catch (err) {}

        phone_id_box.removeAttribute('readonly');
        cpu_id_box.removeAttribute('readonly');

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
                phone_id_box.value != "" && 
                cpu_id_box.value != ""
        ) {
            $.ajax(
                {
                    type : 'POST',
                    url : (
                        '/unchecked/' + 
                        $("#cpu_id").val() + '/' + 
                        $("#phone_id").val()
                    )
                }
            );
        }
    } else {
        if (
                phone_id_box.value != "" && 
                cpu_id_box.value != ""
        ) {
            sendingInterval = window.setInterval(sendAjax, 150);

            cpu_id_box.setAttribute('readonly', 'readonly');
            phone_id_box.setAttribute('readonly', 'readonly');
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
function getSensorData(time_interval_ms) {
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
// Ajax if the cpu_id and the node_id have
// been entered
function sendAjax() {
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
            $("#cpu_id").attr('readonly') != undefined
    ) {
        $.ajax(
            {
                type : 'POST',

                url : (
                    '/checked/' + 
                    $("#cpu_id").val() + '/' + 
                    $("#phone_id").val()
                ),

                data: {
                    sensor_data: JSON.stringify(
                        {
                            orientation: {
                                tilt_horizontal: tiltLR, 
                                tilt_vertical: tiltFB, 
                                direction: dir
                            }, 

                            location: {
                                latitude: lat, 
                                longitude: lon
                            },

                            error: error_data
                        }
                    )
                }
            }
        );
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
        return parseFloat(document.getElementById(element_id).innerHTML);
    }
}