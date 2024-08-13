import { ref } from 'vue'
import { getPlayerID } from './api.js'

const color = ref("")
const gameID = ref("")


let CONN = null;

export function useWebsocket(id) {
    if (window["WebSocket"]) {
        CONN = new WebSocket("ws://" + document.location.host + "/game/" + id);
        CONN.onclose = function(event) {
            CONN = null;
        };

        CONN.error = function(event) {
            CONN = null;
        };
        CONN.onmessage = function(event) {
            var messages = event.data.split('\n');
            for (var i = 0; i < messages.length; i++) {
                var message = messages[i];
                var parsed = JSON.parse(message)
                if (parsed.Action == "join success") {
                    localStorage.setItem("playerID", parsed["PlayerID"])
                }
                color.value = parsed["color"];
                gameID.value = parsed["gameID"];
            }
        };

        CONN.onopen = function(event) {
            // Request game and join
            const msg = {
                action: "join",
                playerID: getPlayerID(),
                date: Date.now(),
            };
            CONN.send(JSON.stringify(msg));
        }
    } else {
        console.log("Your browser does not support WebSockets.")
    }
}


export function sendMove(move) {
    if (CONN != null) {
        const msg = {
            Action: "move",
            PlayerID: getPlayerID(),
            GameID: gameID.value,
            Move: move,
            date: Date.now(),
        };
        CONN.send(JSON.stringify(msg));
    }
}
