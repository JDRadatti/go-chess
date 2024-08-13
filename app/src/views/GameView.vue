<script setup>
import GameBoard from '../components/GameBoard.vue'
import GameSide from '../components/GameSide.vue'
import { useWebsocket } from '../scripts/websocket.js'
import { useRoute } from 'vue-router'
import { onMounted } from 'vue'

const route = useRoute()

onMounted(() => {
    useWebsocket(route.params.id)
})
</script>

<template>
    <main class="game-container">
        <GameBoard />
        <div>
            <GameSide />
        </div>
    </main>
</template>

<style></style>
