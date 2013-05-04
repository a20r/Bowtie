getLocation();
getSensorData();
function getSensorData() {
  if(window.DeviceOrientationEvent) {
    document.getElementById("doEvent").innerHTML="DeviceOrientation"; 
    window.addEventListener('deviceorientation', function(eventData) {
      var tiltLR=eventData.gamma;
      var tiltFB=eventData.beta;
      var dir=eventData.alpha;
      var motUD=null;
      deviceOrientationHandler(tiltLR,tiltFB,dir,"Not supported");
    }, false);
  } else if (window.OrientationEvent) {
    document.getElementById("doEvent").innerHTML="MozOrientation";
    window.addEventListener('MozOrientation', function(eventData) {
      var tiltLR=eventData.x*90;
      var tiltFB=eventData.y*-90;
      var dir=null;
      var motUD=eventData.z;
      deviceOrientationHandler(tiltLR,tiltFB,dir,motUD);
    }, false);
  } else {
      document.getElementById("doEvent").innerHTML="Not supported on your device or browser.  Sorry."
    }
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
    navigator.geolocation.watchPosition(
      devicePositionHandler, 
      function (position) { 
        alert("Position Error");
        $.ajax({
          type : "POST",
          url : "/",
          data: {sensor_data: JSON.stringify({location: "None"})}
          //dataType: 'json'
        });
      },
      { enableHighAccuracy: true, timeout: timeoutVal, maximumAge: 0 });
  } else {
    alert("Geolocation is not supported by this browser");
  }
}
function devicePositionHandler(position) {
  document.getElementById("latPos").innerHTML = position.coords.latitude;
  document.getElementById("longPos").innerHTML = position.coords.longitude;
  $.ajax({
    type: 'POST',
    url:'/',
    data: {sensor_data: JSON.stringify({location: {latitude: position.coords.latitude, longitude: position.coords.longitude}})}
  });
}