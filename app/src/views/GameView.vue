<script setup>
import GameBoard from '../components/GameBoard.vue'

if (window["WebSocket"]) {
    var conn = new WebSocket("ws://" + document.location.host + "/game/123");
    console.log(conn)
    conn.onclose = function (evt) {
        var item = document.createElement("div");
        item.innerHTML = "<b>Connection closed.</b>";
        console.log(item);
    };
    conn.onmessage = function (evt) {
        var messages = evt.data.split('\n');
        for (var i = 0; i < messages.length; i++) {
            var item = document.createElement("div");
            item.innerText = messages[i];
            console.log(item);
        }
    };
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
