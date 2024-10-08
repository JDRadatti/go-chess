<script setup>
import { ref, onMounted, watch, computed } from 'vue'
import CopyLink from '../components/CopyLink.vue'
import { sendAbort, sendResign, acceptDraw, denyDraw, sendDrawRequest } from '../scripts/websocket.js'

// time and increment should be in seconds
const props = defineProps(['whiteTurn', 'whiteTime', 'blackTime', 'increment', 'start', 'color', 'over', 'status', 'move'])

const lastMove = ref("")
const moves = ref([])
const anchorRef = ref(null)

const started = ref(false)
const gameOver = ref(false)
const drawStatus = ref("Request Draw")
const abortClassList = ref(["abort", "hide"])
const resignClassList = ref(["resign", "hide"])
const drawClassList = ref(["draw", "hide"])
const clocksClassList = ref(["clocks", ""])
const clockWhiteList = ref(["clock", ""])
const clockBlackList = ref(["clock", ""])
const whiteTime = ref(0) // seconds left in white's time clock
const blackTime = ref(0)
const whiteTimeFormatted = computed(() => formatSeconds(props.whiteTime))
const blackTimeFormatted = computed(() => formatSeconds(props.blackTime))

function flip() {
    clocksClassList.value[1] = "flipped"
}

function activateWhiteClock() {
    clockWhiteList.value[1] = "active"
    clockBlackList.value[1] = ""
}

function activateBlackClock() {
    clockWhiteList.value[1] = ""
    clockBlackList.value[1] = "active"
}

function formatSeconds(seconds) {
    return new Date(seconds * 1000).toISOString().slice(14, 19)
}

function showButtons() {
    abortClassList.value[1] = ""
    drawClassList.value[1] = ""
    resignClassList.value[1] = ""
}

function hideButtons() {
    abortClassList.value[1] = "hide"
    drawClassList.value[1] = "hide"
    resignClassList.value[1] = "hide"
}

function hideAbort() {
    abortClassList.value[1] = "hide"
}

function drawRequest() {
    drawStatus.value = "Waiting..."
    sendDrawRequest()
}

function addMove() {
    if (moves.value.length != 0 && moves.value[moves.value.length - 1].length == 1) {
        moves.value[moves.value.length - 1].push(lastMove.value)
    } else {
        moves.value.push([lastMove.value])
    }
}

watch(started, () => {
    if (props.color == 1) {
        flip()
    }
    showButtons()
})

watch(props, (props) => {
    if (props.start) {
        started.value = true
    }
    if (props.over) {
        hideButtons()
        gameOver.value = true
    }
    if (props.whiteTurn) {
        activateWhiteClock()
    } else {
        activateBlackClock()
    }
    if (props.status == "draw_deny1" && props.color == 0) {
        drawStatus.value = "Draw Denied"
    } else if (props.status == "draw_deny0" && props.color == 1) {
        drawStatus.value = "Draw Denied"
    }
    if (props.move != lastMove.value) {
        hideAbort()
        lastMove.value = props.move
        addMove()
        anchorRef.value.scrollIntoView()
    }
})

</script>

<template>
    <div class="board-side-container">
        <div :class="clocksClassList">
            <div :class="clockBlackList">
                <p>{{ blackTimeFormatted }}</p>
            </div>
            <div class="middle-container">
                <ol class="moves-container">
                    <li v-for="(row, index) in moves" class="moveRow" :key="index" :id="index">
                        {{ index + 1 }}.
                        <div v-for="(move, index) in row" class="move" :key="index" :id="index">
                            {{ move }}
                        </div>
                    </li>
                    <div id="anchor" ref="anchorRef"></div>
                </ol>
                <CopyLink :show="start"></CopyLink>
                <div class="buttons-container">
                    <button :class="abortClassList" data-type="secondary" @click="sendAbort">Abort</button>
                    <button :class="resignClassList" data-type="secondary" @click="drawRequest">{{ drawStatus
                        }}</button>
                    <button :class="drawClassList" data-type="secondary" @click="sendResign">Resign</button>
                </div>
            </div>
            <div :class="clockWhiteList">
                <p>{{ whiteTimeFormatted }}</p>
            </div>
        </div>
    </div>
</template>

<style scoped>
.hide {
    display: none;
}

.middle-container {
    height: 100%;
    margin-bottom: 1rem;
}

.moves-container {
    height: 50dvh;
    width: 100%;
    padding-left: 0;
    overflow-y: scroll;
}

.moves-container * {
    overflow-anchor: none;
}

#anchor {
    overflow-anchor: auto;
    height: 1px;
}

.move {
    padding: 0.5rem;
}

.moveRow {
    display: flex;
    flex-direction: row;
    justify-content: flex-start;
    align-items: center;
}

.buttons-container {
    display: flex;
    flex-direction: row;
    justify-content: right;
}

p {
    color: var(--light-square);
    width: 10ch;
    text-align: center;
}

.flipped {
    flex-direction: column-reverse !important;
}

.clocks {
    height: var(--board-container-size);
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    align-items: flex-start;
}


.clock {
    padding: 0.5rem;
    background-color: var(--dark-square);
}

.active {
    border: solid var(--light-square);
}
</style>
