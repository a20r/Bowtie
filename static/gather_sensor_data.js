
getLocation();
getSensorData();

function getSensorData() {
  if(window.DeviceOrientationEvent) {
    document.getElementById("doEvent").innerHTML="DeviceOrientation"; 
    window.addEventListener('deviceorientation', orientationEventHandler, false);
  } else if (window.OrientationEvent) {
    document.getElementById("doEvent").innerHTML="MozOrientation";
    window.addEventListener('MozOrientation', mozEventHandler, false);
  } else {
      document.getElementById("doEvent").innerHTML="Not supported on your device or browser.  Sorry."
    }
}

function mozEventHandler(eventData) {
  var tiltLR=eventData.x*90;
  var tiltFB=eventData.y*-90;
  var dir=null;
  var motUD=eventData.z;
  deviceOrientationHandler(tiltLR,tiltFB,dir,motUD);
}

function orientationEventHandler(eventData) {
  var tiltLR=eventData.gamma;
  var tiltFB=eventData.beta;
  var dir=eventData.alpha;
  var motUD=null;
  deviceOrientationHandler(tiltLR,tiltFB,dir,"Not supported");
}

function deviceOrientationHandler(tiltLR,tiltFB,dir,motionUD) {
  document.getElementById("doTiltLR").innerHTML=Math.round(tiltLR);
  document.getElementById("doTiltFB").innerHTML=Math.round(tiltFB);
  document.getElementById("doDirection").innerHTML=Math.round(dir);
  document.getElementById("doMotionUD").innerHTML=motionUD;
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
  alert("Position Error");
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