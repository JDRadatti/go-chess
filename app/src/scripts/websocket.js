import { ref } from 'vue'
import { getPlayerID } from './api.js'

const color = ref("")
const gameID = ref("")
const playerID = getPlayerID()


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
                color.value = parsed["color"];
                gameID.value = parsed["gameID"];
            }
        };

        CONN.onopen = function(event) {
            // Request game and join
            const msg = {
                action: "join",
                playerID: playerID,
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
            PlayerID: playerID,
            GameID: gameID.value,
            Move: move,
            date: Date.now(),
        };
        CONN.send(JSON.stringify(msg));
    }
}
