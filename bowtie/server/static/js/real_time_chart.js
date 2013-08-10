function realtime_demo () {

  var
    x = [],
    dataA = [],
    dataB = [],
    dataC = [],
    data = [[x, dataA], [x, dataB], [x, dataC]],
    options, i, timesries;

  // Mock Data:
  function sample(i) {
    x.push(i);
    dataA.push(+$("#doTiltLR").html());
    dataB.push(+$("#doTiltFB").html());
    dataC.push(+$("#doDirection").html());

  }

  // Initial Data:
  for (i = 0; i < 100; i++) {
    sample(i);
  }

  // Envision Timeseries Options
  options = {
    container : document.getElementById("accelerometer-chart"),
    data : {
      detail : data,
      summary : data
    },
    defaults : {
      summary : {
        config : {
          handles : { show : false }
        }
      }
    }
  }

  // Render the timeseries
  timeseries = new envision.templates.TimeSeries(options);

  // Method to get new data
  // This could be part of an Ajax callback, a websocket callback,
  // or streaming / long-polling data source.
  function getNewData () {
    i++;

    // Short circuit (no need to keep going!  you get the idea)
    if (i > 1000) return;

    sample(i);
    animate(i);
  }

  // Initial request for new data
  getNewData();

  // Animate the new data
  function animate (i) {

    var
      start = (new Date()).getTime(),
      length = 500, // 500ms animation length
      max = i - 1,  // One new point comes in at a time
      min = i - 51, // Show 50 in the top
      offset = 0;   // Animation frame offset

    // Render animation frame
    (function frame () {

      var
        time = (new Date()).getTime(),
        tick = Math.min(time - start, length),
        offset = (Math.sin(Math.PI * (tick) / length - Math.PI / 2) + 1) / 2;

      // Draw the summary first
      timeseries.summary.draw(null, {
        xaxis : {
          min : 0,
          max : max + offset
        }
      });

      // Trigger the select interaction.
      // Update the select region and draw the detail graph.
      timeseries.summary.trigger('select', {
        data : {
          x : {
            min : min + offset,
            max : max + offset
          }
        }
      });

      if (tick < length) {
        setTimeout(frame, 20);
      } else {
        // Pretend new data comes in every second
        setTimeout(getNewData, 500);
      }
    })();
  }
}