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
const status = ref("")

onMounted(() => {
    let CONN = useWebsocket(route.params.id)
    if (CONN == null) {
        console.log("failed to connect to websocket")
    }
    CONN.onmessage = function (event) {
        console.log(event)
        var messages = event.data.split('\n');
        for (var i = 0; i < messages.length; i++) {
            var message = messages[i];
            var parsed = JSON.parse(message)
            if (parsed.Action == "join success") {
                console.log("JOI*N SUCCESS")
                localStorage.setItem("playerID", parsed["PlayerID"])
                color.value = parsed["Color"];
                gameID.value = parsed["GameID"];
            } else if (parsed.Action == "game start") {
                started.value = true
            } else if (parsed.Action == "move") {
                console.log("handle move")
            } else if (parsed.Action == "join fail") {
                alert("game full... redirecting")
                router.push('/play')
            }
            status.value = parsed.Action
        }
    };
})

</script>


<template>
    <main class="game-container">
        <GameBoard :start="started" :color="color" :status="status" />
        <div>
            <GameSide />
        </div>
    </main>
</template>

<style></style>
