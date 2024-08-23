<script setup>
import GameBoard from '../components/GameBoard.vue'
import GameSide from '../components/GameSide.vue'
import { useWebsocket } from '../scripts/websocket.js'
import { getPlayerID } from '../scripts/api.js'
import { useRoute, useRouter } from 'vue-router'
import { onMounted, ref } from 'vue'

const route = useRoute()
const router = useRouter()

const color = ref(-1)
var whiteTime = ref(0)
var blackTime = ref(0)
var increment = ref(0)
const gameID = ref("")
const status = ref("")
const started = ref(false)
const waiting = ref(false)
const move = ref("")
const fen = ref("")
const gameOver = ref(false)
const whiteTurn = ref(true)
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
            if (parsed.Action == "join_success") {
                localStorage.setItem("playerID", parsed["PlayerID"])
                color.value = parsed.Player
                gameID.value = parsed.GameID;
                whiteTime.value = parsed.WhiteTime;
                blackTime.value = parsed.BlackTime;
                increment.value = parsed.Increment;
                waiting.value = true
            } else if (parsed.Action == "game_start") {
                started.value = true
                waiting.value = false
                move.value = parsed.Move
                fen.value = parsed.FEN
            } else if (parsed.Action == "move_success") {
                move.value = parsed.Move
                fen.value = parsed.FEN
            } else if (parsed.Action == "join_fail") {
                alert("game full... redirecting")
                router.push('/play')
            } else if (parsed.Action == "game_kill") {
                alert("Could not find an opponenet... redirecting")
                router.push('/play')
            } else if (parsed.Action == "game_end_time") {
                gameOver.value = true
                if (color.value == 1 && parsed.PlayerID == getPlayerID()) {
                    status.value = "WHITE WON"
                } else if (color.value == 0 && parsed.PlayerID != getPlayerID()) {
                    status.value = "WHITE WON"
                } else {
                    status.value = "BLACK WON"
                }
            } else if (parsed.Action == "game_end") {
                gameOver.value = true
                fen.value = parsed.FEN
                move.value = parsed.Move
                if (move.value == "1-0") {
                    status.value = "WHITE WON"
                } else if (move.value == "0-1") {
                    status.value = "BLACK WON"
                } else if (move.value == "1/2-1/2") {
                    status.value = "DRAW"
                }
            } else if (parsed.Action == "draw_request") {
                status.value = "draw_request"
            } else if (parsed.Action == "draw_deny") {
                status.value = "draw_deny" + parsed.Player
            } else if (parsed.Action == "draw") {
                status.value = "draw"
                gameOver.value = true
            } else if (parsed.Action == "abort") {
                status.value = "aborted"
                gameOver.value = true
            } else if (parsed.Action == "resign") {
                status.value = "resigned"
                gameOver.value = true
            } else if (parsed.Action == "time_update") {
                whiteTime.value = parsed.WhiteTime
                blackTime.value = parsed.BlackTime
            }
            if (parsed.Turn == 0) {
                whiteTurn.value = true
            } else {
                whiteTurn.value = false
            }
            messageCount.value++
        }
    };
})

</script>


<template>
    <main class="game-container">
        <GameBoard :start="started" :color="color" :waiting="waiting" :fen="fen" :count="messageCount" :over="gameOver"
            :status="status" />
        <div>
            <GameSide :start="started" :whiteTurn="whiteTurn" :blackTime="blackTime" :whiteTime="whiteTime"
                :color="color" :over="gameOver" :status="status" :move="move" />
        </div>
    </main>
</template>

<style></style>
