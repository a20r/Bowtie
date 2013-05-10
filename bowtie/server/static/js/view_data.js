
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
function toggle_readonly() {
  var cpu_id_box = document.getElementById("cpu_id");
  if(cpu_id_box.hasAttribute('readonly')){   
    ready_to_start();
  } else {
    if(cpu_id_box.value != "") {
      $.getJSON('/' + cpu_id_box.value, visualize_data);
    } else {
      show_warning("Please enter a CPU Id and Node Id before continuing");
    }
  }
}

// Need to use tail recursion
function visualize_data(cpu_data) {
  //alert(JSON.stringify(node_data));
  if (cpu_data['error']['code'] == 2) {
    show_warning(cpu_data['error']['message']);
    return;
  }
  ready_to_stop();
  var s_table = document.getElementById('sensor_table');
  for (var node_data in cpu_data) {
    //alert(JSON.stringify(cpu_data));
    s_table.innerHTML = JSON.stringify(cpu_data);
  }
  // Need to put this in an interval
  $.getJSON('/' + document.getElementById('cpu_id').value, visualize_data);
}
