<script setup>
import { ref } from 'vue'

const props = defineProps(['show'])
const link = ref(document.location)
const copiedRef = ref(null);

function copyLink() {
    navigator.clipboard.writeText(link.value);
    copiedRef.value.classList.remove("hide")
}
</script>

<template>
    <div v-if="!show">
        <button class="color-primary copy" @click="copyLink">
            <p>Copy link and Share with opponent or wait for matchmaking:</p> {{ link }}
        </button>
        <p class="hide" ref="copiedRef">
            Copied. First person to join will be your opponenet.
        </p>
    </div>
</template>

<style scoped>
.hide {
    display: none;
}

.copy {
    background-color: var(--background-color);
    border: solid var(--color-border);
    padding: 1rem 4rem;
    color: var(--color-border-hover);
}

.copy:hover {
    border: solid var(--color-border-hover);
    color: var(--dark-square);
}

.copy p {
    font-size: 1.2em;
    color: var(--light-square)
}
</style>
