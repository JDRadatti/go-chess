<script setup>
import GameBoard from '../components/GameBoard.vue'
import { useWebsocket } from '../scripts/websocket.js'
import { useRoute } from 'vue-router'
const route = useRoute()
useWebsocket(route.params.id)
</script>

<template>
    <main>
        <GameBoard />
    </main>
</template>

<style></style>
