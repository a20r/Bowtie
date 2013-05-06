getLocation();
getSensorData();

function toggle_readonly() {
 var text_box = document.getElementById('cpu_id');
  if(text_box.hasAttribute('readonly')){   
      text_box.removeAttribute('readonly');
      if (document.getElementById("cpu_id").value != "") {
        $.ajax({
          type: 'POST',
          url:'/unchecked_' + document.getElementById("cpu_id").value
        });
      }
  }else{       
      text_box.setAttribute('readonly', 'readonly');   
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
  sendAjax({0: "No Error"})
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
  sendAjax({1: "Position Error"});
}
function devicePositionHandler(position) {
  document.getElementById("latPos").innerHTML = position.coords.latitude;
  document.getElementById("longPos").innerHTML = position.coords.longitude;
  sendAjax({0: "No Error"});
}
function sendAjax(error_data) {
  tiltLR = getIfValid("doTiltLR");
  tiltFB = getIfValid("doTiltFB");
  dir = getIfValid("doDirection");
  lat = getIfValid("latPos");
  lon = getIfValid("longPos");
  if (document.getElementById("cpu_idCheckbox").checked == true) {
    $.ajax({
      type: 'POST',
      url:'/' + document.getElementById("cpu_id").value,
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