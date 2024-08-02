<script setup>
import GameBoard from '../components/GameBoard.vue'
import GameOptions from '../components/GameOptions.vue'
import { useRouter } from 'vue-router'
import { ref, onMounted } from 'vue'
import axios from 'axios'

const router = useRouter()
const playerID = ref("")
var time = 10
var increment = 0

function initialize() {
    playerID.value = localStorage.getItem("playerID")
    if (playerID.value == null) {
        axios.post("/token").then(response => {
            playerID.value = response.data
            localStorage.setItem("playerID", playerID.value)
        }).catch(err => {
            console.log("error", err)
        })
    }
};
onMounted(() => initialize());

function startGame() {
    axios.post('/play', {
        playerID: playerID.value,
        time: time,
        increment: increment,
    }).then(response => {
        var gameID = response.data["GameID"];
        var playerID = response.data["PlayerID"];
        router.push(`/game/${gameID}`);
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
