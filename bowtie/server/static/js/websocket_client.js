// WEBSOCKET CONSTANTS
var CONNECTED_WS = "CONNECTED WS";
var DISCONNECTED_WS = "DISCONNECTED WS";
var RETRY_WS = "RETRY WS";
var DISCONNECT_WS = "DISCONNECT WS";

function WebSocketClient() { }

WebSocketClient.prototype.init = function(url) {
    var ws = null;

    try {
        ws = new WebSocket(url);

        ws.onopen = function () {
            console.log("Established websocket connection!");
        };

        ws.onmessage = function(message) {
            console.log("Recieved msg from websocket: [", message.data, "]");
            if (message.data === DISCONNECT_WS ) {
                ws.close();
            }
        }

        ws.onerror = function(error) {
            if (error.reason != undefined) {
                console.log("Websocket Error!: " + error.reason);
            } else {
                console.log("Websocket Error Detected!");
            }

            ws.close();
            ws = null;
        };

        ws.onclose = function() {
            console.log("Closed websocket connection!");
            ws.close();
            ws = null;
        }

    } catch (error) {
        console.log("WebScoket Error!: " + error);
    }

    return ws;
}
