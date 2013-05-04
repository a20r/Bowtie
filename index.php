<html>
<head>
  <body>
    <div class="main">
      <h2>Device Orientation</h2>
      <table>
        <tr>
          <td>Event Supported</td>
          <td id="doEvent"></td>
        </tr>
        <tr>
          <td>Tilt Left/Right [tiltLR]</td>
          <td id="doTiltLR"></td>
        </tr>
        <tr>
          <td>Tilt Front/Back [tiltFB]</td>
          <td id="doTiltFB"></td>
        </tr>
        <tr>
          <td>Direction [direction]</td>
          <td id="doDirection"></td>
        </tr>
        <tr>
          <td>Motion Up/Down</td>
          <td id="doMotionUD"></td>
        </tr>
        <tr>
          <td>Latitude</td>
          <td id="latPos"></td>
        </tr>
        <tr>
          <td>Longitude</td>
          <td id="longPos"></td>
        </tr>
      </table>
    </div>
    <script>
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
            function (position) { alert("Position Error"); },
            { enableHighAccuracy: true, timeout: timeoutVal, maximumAge: 0 }
          );
        }
        else {
          alert("Geolocation is not supported by this browser");
        }
      }
      function devicePositionHandler(position) {
        document.getElementById("latPos").innerHTML = position.coords.latitude;
        document.getElementById("longPos").innerHTML = position.coords.longitude;
      }
    </script>
  </body>
</head>
</html>
