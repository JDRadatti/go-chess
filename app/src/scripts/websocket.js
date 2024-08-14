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
    return CONN
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
