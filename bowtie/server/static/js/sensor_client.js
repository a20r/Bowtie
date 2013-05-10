
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
getLocation();
getSensorData(200);

// Sets the window exiting function
window.onbeforeunload = on_exit;

// Occurs when somebody clicks the alert
// close button
function warning_closed() {
  document.getElementById("alert_msg").style.display = "none";
}

// Function fires once the page is closed
function on_exit() {
  var phone_id_box = document.getElementById("phone_id");
  var cpu_id_box = document.getElementById("cpu_id");
  if (phone_id_box.value != "" && cpu_id_box.value != "") {
    $.ajax({
      type: 'POST',
      url:'/unchecked/' + document.getElementById("cpu_id").value + '/' + document.getElementById("phone_id").value
    });
  }
}

// Toggles whether the data is being shown to the user
// and whether it gets sent to the server
function toggle_readonly() {
 var phone_id_box = document.getElementById("phone_id");
 var cpu_id_box = document.getElementById("cpu_id");
  if(phone_id_box.hasAttribute('readonly')){   

      phone_id_box.removeAttribute('readonly');
      cpu_id_box.removeAttribute('readonly');
      document.getElementById("sub_button").innerHTML = "Start sensing";
      document.getElementById("sub_button").className = "btn btn-large btn-success"
      document.getElementById("sensor_table").style.display = "none";

      // deletes the data from the server
      if (phone_id_box.value != "" && cpu_id_box.value != "") {
        $.ajax({
          type: 'POST',
          url:'/unchecked/' + document.getElementById("cpu_id").value + '/' + document.getElementById("phone_id").value
        });
      }
  } else {
    if(phone_id_box.value != "" && cpu_id_box.value != "") {
      cpu_id_box.setAttribute('readonly', 'readonly');
      phone_id_box.setAttribute('readonly', 'readonly');
      document.getElementById("sub_button").innerHTML = "Stop sensing";
      document.getElementById("sub_button").className = "btn btn-large btn-primary btn-danger";
      document.getElementById("sensor_table").style.display = "block";
      document.getElementById("alert_msg").style.display = "none";
    } else {
      document.getElementById("alert_msg").style.display = "block";
    }
  }
}

// Sets the event handler for the orientation sensor.
// If the data cannot be gathered, messages will be
// shown and the checkboxes will be disabled
function getSensorData(time_interval_ms) {
  if(window.DeviceOrientationEvent) {
    window.addEventListener('deviceorientation', orientationEventHandler, false);
    window.setInterval(sendAjax, time_interval_ms);

  } else {
      document.getElementById("doTiltLR").innerHTML = "Not supported";
      document.getElementById("doTiltLRCheckbox").checked = false;
      document.getElementById("doTiltLRCheckbox").disabled = true;
      document.getElementById("doTiltFB").innerHTML = "Not supported";
      document.getElementById("doTiltFBCheckbox").checked = false;
      document.getElementById("doTiltFBCheckbox").disabled = true;
      document.getElementById("doDirection").innerHTML = "Not supported";
      document.getElementById("doDirectionCheckbox").checked = false;
      document.getElementById("doDirectionCheckbox").disabled = true;
    }
}

// Handles the orientation data
function orientationEventHandler(eventData) {
  var tiltLR=eventData.gamma;
  var tiltFB=eventData.beta;
  var dir=eventData.alpha;
  deviceOrientationHandler(tiltLR,tiltFB,dir);
  //sendAjax({code: 0, message: "No Error"})
}

// Writes the orientation data on the HTML page
function deviceOrientationHandler(tiltLR,tiltFB,dir) {
  document.getElementById("doTiltLR").innerHTML=Math.round(tiltLR);
  document.getElementById("doTiltFB").innerHTML=Math.round(tiltFB);
  document.getElementById("doDirection").innerHTML=Math.round(dir);
}

// Sets the locaction event handlers
function getLocation() {
  if (navigator.geolocation) {
    var timeoutVal = 6000;
    var extraGeoParam = {enableHighAccuracy: true, timeout: timeoutVal, maximumAge: 0};
    navigator.geolocation.watchPosition(devicePositionHandler, positionError, extraGeoParam);
  } else {
    alert("Geolocation is not supported by this browser");
  }
}

// Fired if getting the position raises an error
function positionError (position) { 
  document.getElementById("latPos").innerHTML = "Not supported";
  document.getElementById("latPosCheckbox").checked = false;
  document.getElementById("latPosCheckbox").disabled = true;
  document.getElementById("longPos").innerHTML = "Not supported";
  document.getElementById("longPosCheckbox").checked = false;
  document.getElementById("longPosCheckbox").disabled = true;
  //sendAjax({code: 1, message: "Position Error"});
}

// Writes the position of the coordinates onto the
// HTML page
function devicePositionHandler(position) {
  document.getElementById("latPos").innerHTML = position.coords.latitude;
  document.getElementById("longPos").innerHTML = position.coords.longitude;
  document.getElementById("latPosCheckbox").checked = true;
  document.getElementById("latPosCheckbox").disabled = false;
  document.getElementById("longPosCheckbox").checked = true;
  document.getElementById("longPosCheckbox").disabled = false;

  //sendAjax({code: 0, message:"No Error"});
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
  if (lat == null && lon == null && document.getElementById("latPos").innerHTML == "Not supported") {
    error_data = {code: 1, message: "Position error"}
  } else {
    error_data = {code: 0, message: "No error"}
  }
  if (document.getElementById("cpu_id").hasAttribute('readonly')) {
    $.ajax({
      type: 'POST',
      url:'/checked/' + document.getElementById("cpu_id").value + '/' + document.getElementById("phone_id").value,
      data: {sensor_data: JSON.stringify({orientation: {tilt_horizontal: tiltLR, 
        tilt_vertical: tiltFB, direction: dir}, location: {latitude: lat, longitude: lon}, error: error_data})}
    });
  }
}

// Checks if the data is valid before sending
// it out
function getIfValid(element_id) {
  if (document.getElementById(element_id).innerHTML == "Not supported" ||
      document.getElementById(element_id).innerHTML == "" || 
      !document.getElementById(element_id + "Checkbox").checked) {
    return null;
  } else {
    return parseFloat(document.getElementById(element_id).innerHTML);
  }
}