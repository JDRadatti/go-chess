<script setup>
import { ref, onMounted, watch, computed } from 'vue'

// time and increment should be in seconds
const props = defineProps(['whiteTurn', 'whiteTime', 'blackTime', 'increment', 'start', 'color', 'over'])

const started = ref(false)
const gameOver = ref(false)
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

watch(started, () => {
    if (props.color == 1) {
        flip()
    }
})

watch(props, (props) => {
    if (props.start) {
        started.value = true
    }
    if (props.over) {
        gameOver.value = true
    }
    if (props.whiteTurn) {
        activateWhiteClock()
    } else {
        activateBlackClock()
    }
})
</script>

<template>
    <div class="board-side-container">
        <div :class="clocksClassList">
            <div :class="clockBlackList">
                <p>{{ blackTimeFormatted }}</p>
            </div>
            <div :class="clockWhiteList">
                <p>{{ whiteTimeFormatted }}</p>
            </div>
        </div>
    </div>
</template>

<style scoped>
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
    align-items: center;
}


.clock {
    padding: 0.5rem;
    background-color: var(--dark-square);
}

.active {
    border: solid var(--light-square);
}
</style>
