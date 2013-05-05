
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
  $.ajax({
    type: 'POST',
    url:'/',
    data: {sensor_data: JSON.stringify({location: {tilt_horizontal: tiltLR, 
      tilt_vertical: tiltFB, direction: dir}})}
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
  // sends the data to the server
  $.ajax({
    type : "POST",
    url : "/",
    data: {sensor_data: JSON.stringify({location: null})}
  });
}

function devicePositionHandler(position) {

  // changes the numbers on the html page
  document.getElementById("latPos").innerHTML = position.coords.latitude;
  document.getElementById("longPos").innerHTML = position.coords.longitude;

  // sends the data to the server
  $.ajax({
    type: 'POST',
    url:'/',
    data: {sensor_data: JSON.stringify({location: {latitude: position.coords.latitude, 
      longitude: position.coords.longitude}})}
  });
}