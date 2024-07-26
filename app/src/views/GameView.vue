<script setup>
import GameBoard from '../components/GameBoard.vue'
import { ref } from 'vue'
import { useRoute } from 'vue-router'

const move = ref("")
const gameID = ref("")
const playerID = ref("")
var connected = false
const route = useRoute()

if (window["WebSocket"]) {
    var conn = new WebSocket("ws://" + document.location.host + "/game/" + route.params.id);
    conn.onclose = function (event) {
        var item = document.createElement("div");
        item.innerHTML = "<b>Connection closed.</b>";
        connected = false
    };
    conn.onmessage = function (event) {
        var messages = event.data.split('\n');
        for (var i = 0; i < messages.length; i++) {
            var message = messages[i];
            var parsed = JSON.parse(message)
            playerID.value = parsed["PlayerID"];
            gameID.value = parsed["GameID"];
        }
    };

    conn.onopen = function (event) {
        connected = true
    }


} else {
    var item = document.createElement("div");
    item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
    console.log(item);
}

function sendMove() {
    if (connected) {
        const msg = {
            PlayerID: playerID.value,
            GameID: gameID.value,
            Move: move.value,
            date: Date.now(),
        };
        conn.send(JSON.stringify(msg));
        console.log("message", msg)
    }
}
</script>

<template>
    <main>
        <h1>Game</h1>
        <GameBoard />
        <form>
            <input v-model="move" placeholder="Make Move"></input>
        </form>
        <button @click="sendMove">Send</button>
        <p>Previous move: {{ move }}</p>
        <p>PlayerID: {{ gameID }}</p>
        <p>GameID: {{ playerID }}</p>
    </main>
</template>

<style></style>
