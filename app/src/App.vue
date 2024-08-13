<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>

<template>
    <div>
        <nav class="main-nav">
            <RouterLink to="/">Home</RouterLink>
            <RouterLink to="/play/">Play</RouterLink>
        </nav>
    </div>
    <RouterView />
</template>

<style scoped>
.main-nav {
    display: flex;
    max-width: var(--navbar-width);
    flex-direction: column;
    position: absolute;
    top: 2rem;
    left: 2rem;
}
</style>
