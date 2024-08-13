<script setup>
import GameBoard from '../components/GameBoard.vue'
import GameOptions from '../components/GameOptions.vue'
import { startGame, getPlayerID } from '../scripts/api.js'
import { useRouter } from 'vue-router'
import {onMounted} from 'vue'

const router = useRouter();

function clickStart() {
    startGame().then((response) => {
        if (response["GameID"]) {
            router.push('/game/' + response["GameID"]);
        }
    }).catch((error) => {
        console.log(error);
    })
}

onMounted(() => {
    getPlayerID()
})
</script>

<template>
    <main class="game-container">
        <GameBoard />
        <div>
            <GameOptions />
            <button @click="clickStart"> Start Game </button>
        </div>
    </main>
</template>

<style>
</style>
