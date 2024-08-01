<script setup>
import GameBoard from '../components/GameBoard.vue'
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import axios from 'axios'

const move = ref("")
const gameID = ref("")
const playerID = ref("")
const color = ref("")
var connected = false
const route = useRoute()

playerID.value = localStorage.getItem("playerID")
if (playerID.value == null) {
    axios.post("/token").then(response => {
        playerID.value = response.data
        localStorage.setItem("playerID", playerID.value)
    }).catch(err => {
        console.log("error", err)
    })
}

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
            color.value = parsed["Color"];
            gameID.value = parsed["GameID"];
            playerID.value = parsed["PlayerID"];
        }
    };

    conn.onopen = function (event) {
        connected = true
        const msg = {
            PlayerID: playerID.value,
            date: Date.now(),
        };
        conn.send(JSON.stringify(msg));
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
        <p>PlayerID: {{ playerID }}</p>
        <p>GameID: {{ gameID }}</p>
    </main>
</template>

<style></style>
