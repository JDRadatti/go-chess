<script setup>
import GameBoard from '../components/GameBoard.vue'
import GameOptions from '../components/GameOptions.vue'
import { startGame } from '../scripts/api.js'
import { useRouter } from 'vue-router'

const router = useRouter();

// time should be in minutes, increment in seconds
function clickStart(time, increment) {
    startGame(time, increment).then((response) => {
        if (response["GameID"]) {
            router.push('/game/' + response["GameID"]);
        }
    }).catch((error) => {
        console.log(error);
    })
}
</script>

<template>
    <main class="game-container">
        <GameBoard />
        <div>
            <GameOptions @start='clickStart' />
        </div>
    </main>
</template>

<style></style>
