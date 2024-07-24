<script setup>
import GameBoard from '../components/GameBoard.vue'
import { ref } from 'vue'
import { useRoute } from 'vue-router'

const move = ref("")
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
            var item = document.createElement("div");
            item.innerText = messages[i];
            console.log("MESSAGE", messages[i]);
        }
    };

    conn.onopen = function (event) {
        connected = true
        conn.send("I have joined the game");
    }


} else {
    var item = document.createElement("div");
    item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
    console.log(item);
}

function sendMove() {
    if (connected) {
        //const msg = {
        //    move: move.value,
        //    date: Date.now(),
        //};

        // Send the msg object as a JSON-formatted string.
        //conn.send(JSON.stringify(msg));
        conn.send(move.value);
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
    </main>
</template>

<style></style>
