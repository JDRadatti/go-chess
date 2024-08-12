<script setup>
import GameBoard from '../components/GameBoard.vue'
import GameOptions from '../components/GameOptions.vue'
import { startGame } from '../scripts/api.js'
import { useRouter } from 'vue-router'

const router = useRouter();

function clickStart() {
    startGame().then((response) => {
        if (response["GameID"]) {
        console.log("RESPONSE", response);
            router.push('/game/' + response["GameID"]);
        }
    }).catch((error) => {
        console.log(error);
    })
}
</script>

<template>
    <main>
        <h1>Game</h1>
        <GameBoard />
        <GameOptions />
        <button @click="clickStart"> Start Game </button>
    </main>
</template>

<style></style>
