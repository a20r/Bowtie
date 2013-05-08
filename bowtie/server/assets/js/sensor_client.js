getLocation();
getSensorData();

window.onbeforeunload = on_exit;

function warning_closed() {
  document.getElementById("alert_msg").style.display = "none";
}

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

function toggle_readonly() {
 var phone_id_box = document.getElementById("phone_id");
 var cpu_id_box = document.getElementById("cpu_id");
  if(phone_id_box.hasAttribute('readonly')){   

      phone_id_box.removeAttribute('readonly');
      cpu_id_box.removeAttribute('readonly');
      document.getElementById("sub_button").innerHTML = "Start sensing";
      document.getElementById("sub_button").className = "btn btn-large btn-success"

      document.getElementById("sensor_table").style.display = "none";
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

function getSensorData() {
  if(window.DeviceOrientationEvent) {
    window.addEventListener('deviceorientation', orientationEventHandler, false);
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

function orientationEventHandler(eventData) {
  var tiltLR=eventData.gamma;
  var tiltFB=eventData.beta;
  var dir=eventData.alpha;
  deviceOrientationHandler(tiltLR,tiltFB,dir);
  sendAjax({code: 0, message: "No Error"})
}

function deviceOrientationHandler(tiltLR,tiltFB,dir) {
  document.getElementById("doTiltLR").innerHTML=Math.round(tiltLR);
  document.getElementById("doTiltFB").innerHTML=Math.round(tiltFB);
  document.getElementById("doDirection").innerHTML=Math.round(dir);
}

function getLocation() {
  if (navigator.geolocation) {
    var timeoutVal = 6000;
    var extraGeoParam = {enableHighAccuracy: true, timeout: timeoutVal, maximumAge: 0};
    navigator.geolocation.watchPosition(devicePositionHandler, positionError, extraGeoParam);
  } else {
    alert("Geolocation is not supported by this browser");
  }
}

function positionError (position) { 
  document.getElementById("latPos").innerHTML = "Not supported";
  document.getElementById("latPosCheckbox").checked = false;
  document.getElementById("latPosCheckbox").disabled = true;
  document.getElementById("longPos").innerHTML = "Not supported";
  document.getElementById("longPosCheckbox").checked = false;
  document.getElementById("longPosCheckbox").disabled = true;
  sendAjax({code: 1, message: "Position Error"});
}
function devicePositionHandler(position) {
  document.getElementById("latPos").innerHTML = position.coords.latitude;
  document.getElementById("longPos").innerHTML = position.coords.longitude;
  document.getElementById("latPosCheckbox").checked = true;
  document.getElementById("latPosCheckbox").disabled = false;
  document.getElementById("longPosCheckbox").checked = true;
  document.getElementById("longPosCheckbox").disabled = false;

  sendAjax({code: 0, message:"No Error"});
}
function sendAjax(error_data) {
  tiltLR = getIfValid("doTiltLR");
  tiltFB = getIfValid("doTiltFB");
  dir = getIfValid("doDirection");
  lat = getIfValid("latPos");
  lon = getIfValid("longPos");
  if (document.getElementById("cpu_id").hasAttribute('readonly')) {
    $.ajax({
      type: 'POST',
      url:'/checked/' + document.getElementById("cpu_id").value + '/' + document.getElementById("phone_id").value,
      data: {sensor_data: JSON.stringify({orientation: {tilt_horizontal: tiltLR, 
        tilt_vertical: tiltFB, direction: dir}, location: {latitude: lat, longitude: lon}, error: error_data})}
    });
  }
}

function getIfValid(element_id) {
  if (document.getElementById(element_id).innerHTML == "Not supported" ||
      document.getElementById(element_id).innerHTML == "" || 
      !document.getElementById(element_id + "Checkbox").checked) {
    return null;
  } else {
    return parseFloat(document.getElementById(element_id).innerHTML);
  }
}