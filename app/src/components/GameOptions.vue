<script setup>
import { ref, onMounted } from 'vue'

const time = ref(3)
const increment = ref(0)

const times = ref([1, 3, 5, 10])
const increments = ref([0, 1, 2, 10])
const incRefs = ref([])
const timeRefs = ref([])
const activeTime = ref(null)
const activeInc = ref(null)

const activeIndex = 0

function timeClick(event, index) {
    event.target.classList.add("active")
    activeTime.value.classList.remove("active")
    activeTime.value = event.target
    time.value = times.value[index]
}

function incClick(event, index) {
    event.target.classList.add("active")
    activeInc.value.classList.remove("active")
    activeInc.value = event.target
    increment.value = increments.value[index]
}

onMounted(() => {
    for (let i = 0; i < times.value.length; i++) {
        if (times.value[i] == time.value) { // convert to minutes
            timeRefs.value[i].classList.add("active");
            activeTime.value = timeRefs.value[i];
        }
    }

    for (let i = 0; i < increments.value.length; i++) {
        if (increments.value[i] == increment.value) {
            incRefs.value[i].classList.add("active");
            activeInc.value = incRefs.value[i];
        }
    }

})
</script>

<template>
    <div class="board-side-container">
        <div class="container">
            <h2>Time Control</h2>
            <div class="time">
                <button v-for="(time, index) in times" :key="index" @click="timeClick($event, index)"
                    data-type="secondary" ref="timeRefs">
                    {{ time }} Min.
                </button>
            </div>
            <h2> Increment</h2>
            <div class="increment">
                <button v-for="(inc, index) in increments" :key="index" @click="incClick($event, index)"
                    data-type="secondary" ref="incRefs">
                    {{ inc }} Sec.
                </button>
            </div>
            <button class="play" data-type="primary" @click="$emit('start', time * 60, increment)">Play</button>
        </div>
    </div>
</template>

<style scoped>
.container {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
}

.play {
    width: 100%;
    margin-top: 1.5rem;
}

.time>* {
    margin: 0rem 0.10rem;
}

.increment>* {
    margin: 0rem 0.10rem;
}

.container>* {
    padding-bottom: 1rem;
}
</style>
