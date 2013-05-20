
////////////////////////////////////////////////
// view_data.js
//
////////////////////////////////////////////////

// Occurs when somebody clicks the alert
// close button
function warning_closed() {
  document.getElementById("alert_msg").style.display = "none";
}

function show_warning(msg) {
  document.getElementById("alert_msg").innerHTML =  "<button type='button' class='close' onClick='warning_closed()'" + 
                                                    "data-dismiss='alert'>&times;</button>" + 
                                                    "<strong>Error!</strong> " + msg;
  document.getElementById("alert_msg").style.display = "block";
}

function ready_to_start() {
  var cpu_id_box = document.getElementById("cpu_id");
  cpu_id_box.removeAttribute('readonly');
  document.getElementById("sub_button").innerHTML = "Grab data";
  document.getElementById("sub_button").className = "btn btn-large btn-success"
  document.getElementById("sensor_table").style.display = "none";
}

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

function println(msg) {
  document.getElementById('sensor_table').innerHTML += (msg + "<br>");
}

var prev_nodes = new Array();
function visualize_data(cpu_data) {
  //alert(JSON.stringify(cpu_data["error"]));
  if (cpu_data['error']['code'] == 2) {
    show_warning(cpu_data['error']['message']);
    clearInterval(intervalVar);
    ready_to_start();
    return;
  }
  ready_to_stop();
  var s_table = document.getElementById('sensor_table');
  //s_table.innerHTML = "";

  for (var node_name in cpu_data) {
    if (document.getElementById(node_name)) {
      if (node_name != "error") {
        for (var sensor_name in cpu_data[node_name]) {
          for (var sensor_component in cpu_data[node_name][sensor_name]) {
            document.getElementById(node_name + "_" + sensor_component).innerHTML = String(cpu_data[node_name][sensor_name][sensor_component]);
          }
        }
      }
    } else {
      if (node_name != "error") {
        s_table.innerHTML += "<tr id = " + node_name + ">" + node_name;
        for (var sensor_name in cpu_data[node_name]) {
          s_table.innerHTML += "<tr>";
          for (var sensor_component in cpu_data[node_name][sensor_name]) {
            s_table.innerHTML += "<td><b>" + sensor_component + "</b> </td><td id = " + node_name + "_" + sensor_component + 
                                  ">" + cpu_data[node_name][sensor_name][sensor_component] + "</td>";
          }
          s_table.innerHTML += "</tr>";
        }
        s_table.innerHTML += "</tr>";
      }
    }
  }
}
