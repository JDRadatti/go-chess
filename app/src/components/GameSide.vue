<script setup>
import { ref, onMounted, watch, computed } from 'vue'

// time and increment should be in seconds
const props = defineProps(['whiteTurn', 'time', 'increment', 'start', 'color', 'over'])

const started = ref(false)
const clocksClassList = ref(["clocks", ""])
const clockWhiteList = ref(["clock", "active"])
const clockBlackList = ref(["clock", ""])
const whiteTime = ref(0) // seconds left in white's time clock
const blackTime = ref(0)
const whiteTimeFormatted = computed(() => formatSeconds(whiteTime.value))
const blackTimeFormatted = computed(() => formatSeconds(blackTime.value))

setInterval(countdown, 1000)

function countdown() {
    if (!props.start) {
        return
    }
    if (props.whiteTurn) {
        whiteTime.value = whiteTime.value - 1 + props.increment
    } else {
        blackTime.value = blackTime.value - 1 + props.increment
    }
}

function flip() {
    if (clocksClassList.value[1] == "flipped") {
        clocksClassList.value[1] = ""
    } else {
        clocksClassList.value[1] = "flipped"
    }
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
    whiteTime.value = props.time
    blackTime.value = props.time
    console.log("props.color", props.color);
    if (props.color == 1) {
        flip()
    }
})

watch(props, (props) => {
    if (props.start) {
        started.value = true
    } else if (props.over) {
        started = false
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
            <div :class="clockWhiteList">
                <p>{{ whiteTimeFormatted }}</p>
            </div>
            <div :class="clockBlackList">
                <p>{{ blackTimeFormatted }}</p>
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
