<script setup>
import GameBoard from '../components/GameBoard.vue'
import GameSide from '../components/GameSide.vue'
import { useWebsocket } from '../scripts/websocket.js'
import { useRoute, useRouter } from 'vue-router'
import { onMounted, ref } from 'vue'

const route = useRoute()
const router = useRouter()

const color = ref(-1)
const gameID = ref("")
const started = ref(false)
const waiting = ref(false)
const move = ref("")
const fen = ref("")
const messageCount = ref(0)

onMounted(() => {
    let CONN = useWebsocket(route.params.id)
    if (CONN == null) {
        alert("failed to connect to websocket")
        router.push('/play')
    }
    CONN.onmessage = function (event) {
        var messages = event.data.split('\n');
        for (var i = 0; i < messages.length; i++) {
            var message = messages[i];
            var parsed = JSON.parse(message)
            if (parsed.Action == "join success") {
                localStorage.setItem("playerID", parsed["PlayerID"])
                color.value = parsed["Color"];
                gameID.value = parsed["GameID"];
                waiting.value = true
            } else if (parsed.Action == "game start") {
                started.value = true
                waiting.value = false
            } else if (parsed.Action == "move") {
                move.value = parsed.Move
                fen.value = parsed.FEN
            } else if (parsed.Action == "join fail") {
                alert("game full... redirecting")
                router.push('/play')
            }
            messageCount.value++
        }
    };
})

</script>


<template>
    <main class="game-container">
        <GameBoard :start="started" :color="color" :waiting="waiting" :fen="fen" :count="messageCount" />
        <div>
            <GameSide />
        </div>
    </main>
</template>

<style></style>
