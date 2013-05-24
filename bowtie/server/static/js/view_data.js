
////////////////////////////////////////////////
// view_data.js
//
////////////////////////////////////////////////

// Occurs when somebody clicks the alert
// close button
function warning_closed() {
  document.getElementById("alert_msg").style.display = "none";
}

// Shows a warning message in the 
// form of a alert
function show_warning(msg) {
  document.getElementById("alert_msg").innerHTML =  "<button type='button' class='close' onClick='warning_closed()'" + 
                                                    "data-dismiss='alert'>&times;</button>" + 
                                                    "<strong>Error!</strong> " + msg;
  document.getElementById("alert_msg").style.display = "block";
}

// Gets the form ready for data
// to be received from the server
function ready_to_start() {
  var cpu_id_box = document.getElementById("cpu_id");
  cpu_id_box.removeAttribute('readonly');
  document.getElementById("sub_button").innerHTML = "Grab data";
  document.getElementById("sub_button").className = "btn btn-large btn-success"
  document.getElementById("sensor_table").style.display = "none";
}

// Gets the form ready to stop 
// getting data from the server
function ready_to_stop() {
  var cpu_id_box = document.getElementById("cpu_id");
  cpu_id_box.setAttribute('readonly', 'readonly');
  document.getElementById("sub_button").innerHTML = "Stop grabbing";
  document.getElementById("sub_button").className = "btn btn-large btn-primary btn-danger";
  document.getElementById("sensor_table").style.display = "block";
  warning_closed();
}

// Toggles whether the data is being shown to the user
// and whether it gets sent to the server
var intervalVar;
function toggle_readonly() {
  var cpu_id_box = document.getElementById("cpu_id");
  if(cpu_id_box.hasAttribute('readonly')){   
    clearInterval(intervalVar);
    ready_to_start();
  } else {
    if(cpu_id_box.value != "") {
      intervalVar = setInterval(function() {$.getJSON('/get_data/' + cpu_id_box.value, visualize_data)}, 50);
    } else {
      show_warning("Please enter a CPU Id before continuing");
    }
  }
}

// Makes the string presentable (obviously)
// Capitalizes the first letter and 
// puts space where there is an underscore
function makeStringPresentable(string) {
  var spaceString = string.replace("_", " ");
  return spaceString.charAt(0).toUpperCase() + spaceString.slice(1);
}

// Creates the tables dynamically that visualize
// the data being sent to the server
function createTables(cpu_data, node_name, s_table) {
  if (node_name != "error") {
    s_table.innerHTML += "<span id = " + node_name + "_picdiv class='container' style='display:block -webkit-perspective: 400px;'>"
                       + "<b style='margin-top:2;margin-right:10;margin-bottom:4'><font size='5'>" + node_name + "</font></b>" + 
                         "<img src='img/black-bow-tie.png' width = '100' height = '40' " + 
                         "id=" + node_name + "_picture class='logo'></span>";
    s_table.innerHTML += "<table id =" + node_name + " class='table table-hover' border='0'>";
    var n_table = document.getElementById(node_name);
    for (var sensor_name in cpu_data[node_name]) {
      n_table.innerHTML += "<tr>";
      for (var sensor_component in cpu_data[node_name][sensor_name]) {
        n_table.innerHTML += "<td><b>" + makeStringPresentable(sensor_component) + "</b>" + 
                             "</td><td id = " + node_name + "_" + sensor_component + ">" + 
                             cpu_data[node_name][sensor_name][sensor_component] + "</td>";
      }
      n_table.innerHTML += "</tr>";
    }
    s_table.innerHTML += "</table>";
  }
}

// Updates the existing tables with new sensor data
function updateTables(cpu_data, node_name, s_table) {
  if (node_name != "error") {
    for (var sensor_name in cpu_data[node_name]) {
      for (var sensor_component in cpu_data[node_name][sensor_name]) {
        document.getElementById(node_name + "_" + sensor_component).innerHTML = String(cpu_data[node_name][sensor_name][sensor_component]);
      }
    }
  }
  var tiltLR = document.getElementById(node_name + "_tilt_horizontal").innerHTML;
  var tiltFB = document.getElementById(node_name + "_tilt_vertical").innerHTML;
  var dir = document.getElementById(node_name + "_direction").innerHTML;
  document.getElementById(node_name + "_picture").style.webkitTransform = "rotateX(" + (tiltFB * -1) + "deg)" + 
                                                                          " rotateY(" + tiltLR + "deg)";
}

// Takes the set difference of two arrays
Array.prototype.diff = function(a) {
    return this.filter(function(i) {return !(a.indexOf(i) > -1);});
};

function jsonToArray(jsondata) { 
  var retArray = new Array();
  var i = 0;
  for (var label in jsondata) {
    retArray[i++] = label;
  }
  return retArray;
}

// Revives the data if the table
// was created and made to display none
function reviveData(data_id) {
  document.getElementById(data_id).style.display = "block";
  document.getElementById(data_id + "_picdiv").style.display = "block";
  document.getElementById(data_id + "_header").style.display = "block";
}

// Sets the display of table to none
function killData(data_id) {
  document.getElementById(data_id).style.display = "none";
  document.getElementById(data_id + "_picdiv").style.display = "none";
  //document.getElementById(data_id + "_header").style.display = "none";
}

// Gets rid of sensor tables
// that are representing nodes
// that data is not being sent
// for
var prev_nodes = new Array();
function filterOldData(cpu_data) {
  if (prev_nodes.length > 0) {
    var rmData = prev_nodes.diff(jsonToArray(cpu_data));
    for (var i in rmData) {
      try {
        killData(rmData[i]);
      } catch (err) {}
    }
  }
}

// Visualizes the data by putting the data
// values in a table. It is dynamic so
// once a sensor node starts sending or stops
// sending data it is noticable.
function visualize_data(cpu_data) {
  //alert(JSON.stringify(cpu_data["error"]));
  if (cpu_data['error']['code'] == 2) {
    show_warning(cpu_data['error']['message']);
    clearInterval(intervalVar);
    ready_to_start();
    return;
  }
  filterOldData(cpu_data);
  prev_nodes = jsonToArray(cpu_data)
  ready_to_stop();
  var s_table = document.getElementById('sensor_table');
  for (var node_name in cpu_data) {
    if (document.getElementById(node_name)) {
      if (document.getElementById(node_name).style.display == "none") {
        reviveData(node_name);
      }
      updateTables(cpu_data, node_name, s_table);
    } else {
      createTables(cpu_data, node_name, s_table);
    }
  }
}
