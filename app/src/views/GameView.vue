<script setup>
import GameBoard from '../components/GameBoard.vue'

if (window["WebSocket"]) {
    var conn = new WebSocket("ws://" + document.location.host + "/game/123");
    conn.onclose = function (event) {
        var item = document.createElement("div");
        item.innerHTML = "<b>Connection closed.</b>";
        console.log("CONNECTIKON CLOSED");
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
        conn.send("Here's some text that the server is urgently awaiting!");
    }


} else {
    var item = document.createElement("div");
    item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
    console.log(item);
}
</script>

<template>
    <main>
        <h1>Game</h1>
        <GameBoard />
    </main>
</template>

<style></style>
