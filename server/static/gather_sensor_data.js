
getLocation();
getSensorData();

function getSensorData() {
  if(window.DeviceOrientationEvent) {
    document.getElementById("doEvent").innerHTML="DeviceOrientation"; 
    window.addEventListener('deviceorientation', orientationEventHandler, false);
  } else {
      document.getElementById("doEvent").innerHTML="Not supported on your device or browser.  Sorry."
    }
}

function orientationEventHandler(eventData) {
  var tiltLR=eventData.gamma;
  var tiltFB=eventData.beta;
  var dir=eventData.alpha;
  if (document.getElementById("latPos").innerHTML == "Not supported" ||
      document.getElementById("latPos").innerHTML == "") {
    lat = null;
    lon = null;
  } else {
    lat = float(document.getElementById("latPos").innerHTML);
    lon = float(document.getElementById("longPos").innerHTML);
  }
  $.ajax({
    type: 'POST',
    url:'/',
    data: {sensor_data: JSON.stringify({tilt: {tilt_horizontal: tiltLR, 
      tilt_vertical: tiltFB, direction: dir}, location: {latitude: lat, longitude: lon}})}
  });
  deviceOrientationHandler(tiltLR,tiltFB,dir);
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
  document.getElementById("longPos").innerHTML = "Not supported";
}

function devicePositionHandler(position) {

  // changes the numbers on the html page
  document.getElementById("latPos").innerHTML = position.coords.latitude;
  document.getElementById("longPos").innerHTML = position.coords.longitude;
}