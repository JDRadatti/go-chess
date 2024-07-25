<script setup>
import GameBoard from '../components/GameBoard.vue'
import GameOptions from '../components/GameOptions.vue'
import { useRouter } from 'vue-router'
import axios from 'axios'

const router = useRouter()
var time = 10
var increment = 0

function startGame() {
    axios.post('/play', {
        Time: time,
        Increment: increment,
    }).then(response => {
        const gameID = response.data
        router.push(`/game/${gameID}`)
    }).catch(error => {
        console.log(error)
    })
}
</script>

<template>
    <main>
        <h1>Game</h1>
        <GameBoard />
        <GameOptions />
        <button @click="startGame"> Start Game </button>
    </main>
</template>

<style></style>
