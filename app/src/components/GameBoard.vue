<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'

const PieceToIcon = {
    "B": "bb",
    "b": "bw",
    "R": "rb",
    "r": "rw",
    "K": "kb",
    "k": "kw",
    "Q": "qb",
    "q": "qw",
    "N": "nb",
    "n": "nw",
    "P": "pb",
    "p": "pw",
}

function updateBoard(fen) {
    const rows = fen.split("/")
    var index = 0
    for (let i = 0; i < rows.length; i++) { // must be 8 rows
        for (let j = 0; j < rows[i].length; j++) {
            if (!isNaN(rows[i][j])) {
                for (let k = Number(rows[i][j]); k > 0; k--) {
                    updateSquare(index, "")
                    index++
                }
            } else {
                updateSquare(index, rows[i][j])
                index++
            }
        }
    }
}

// Update the square at index to piece
function updateSquare(index, piece) {
    if (piece in PieceToIcon) {
        squares.value[index] = PieceToIcon[piece]
    }
}

const start = ref(-1);
const dest = ref(-1);
let dragged = null;
function dragstartHandler(ev) {
    if (dragged != null) {
        return
    }

    dragged = ev.target
    for (let i = 0; i < dragged.classList.length; i++) {
        if (dragged.classList[i].startsWith('square-')) {
            start.value = Number(dragged.classList[i].slice(7))
        }
    }

    ev.target.classList.add("hide");
    return
}

function dragenterHandler(ev) {

    if (dragged == null) {
        return
    }

    ev.target.classList.add("drag-hover")
    ev.dataTransfer.dropEffect = "move";
}
function dragleaveHandler(ev) {

    if (dragged == null) {
        return
    }
    ev.preventDefault();
    ev.dataTransfer.dropEffect = "move";
    ev.target.classList.remove("drag-hover")
}

function dropHandler(ev, n) {
    if (dragged == null) {
        return
    }
    ev.preventDefault();
    ev.dataTransfer.dropEffect = "move";
    ev.target.classList.remove("drag-hover")
    dest.value = n - 1
    if (dragged != null) {
        dragged.classList.remove("square-" + start.value);
        dragged.classList.add("square-" + dest.value);
        dragged.classList.remove("hide");
    }
    dragged = null;
}

onMounted(() => {
    //updateBoard("RNBQKBNR/PPPPPPPP/8/8/8/8/pppppppp/rnbqkbnr")
    const squareElements = document.querySelectorAll('.square')
    var rank = 0
    for (let i = 0; i < squareElements.length; i++) {
        if (i % 8 == 0 && i > 0) {
            rank++
        }
        if (i % 2 == (rank % 2)) {
            squareElements[i].classList.add("dark")
        } else {
            squareElements[i].classList.add("light")
        }
    }

});

</script>

<template>
    <h1 class="green">{{ start }} {{ dest }}</h1>
    <div class="board" draggable="false">
        <div class="square unselectable" v-for="n in 64" draggable="false" @dragenter="dragenterHandler($event)"
            @dragleave="dragleaveHandler($event)" @drop="dropHandler($event, n)" @dragover.prevent @dragenter.prevent>
        </div>
    </div>


    <div class="pieces">
        <div class="piece kb square-4" draggable="true" @dragstart="dragstartHandler($event)"> </div>
        <div class="piece qb square-3" draggable="true" @dragstart="dragstartHandler($event)"> </div>
        <div class="piece rb square-7" draggable="true" @dragstart="dragstartHandler($event)"> </div>
        <div class="piece rb square-0" draggable="true" @dragstart="dragstartHandler($event)"> </div>
        <div class="piece bb square-5" draggable="true" @dragstart="dragstartHandler($event)"> </div>
        <div class="piece bb square-2" draggable="true" @dragstart="dragstartHandler($event)"> </div>
        <div class="piece nb square-6" draggable="true" @dragstart="dragstartHandler($event)"> </div>
        <div class="piece nb square-1" draggable="true" @dragstart="dragstartHandler($event)"> </div>
        <div class="piece pb square-8" draggable="true" @dragstart="dragstartHandler($event)"> </div>
        <div class="piece pb square-9" draggable="true" @dragstart="dragstartHandler($event)"> </div>
        <div class="piece pb square-10" draggable="true" @dragstart="dragstartHandler($event)"> </div>
        <div class="piece pb square-11" draggable="true" @dragstart="dragstartHandler($event)"> </div>
        <div class="piece pb square-12" draggable="true" @dragstart="dragstartHandler($event)"> </div>
        <div class="piece pb square-13" draggable="true" @dragstart="dragstartHandler($event)"> </div>
        <div class="piece pb square-14" draggable="true" @dragstart="dragstartHandler($event)"> </div>
        <div class="piece pb square-15" draggable="true" @dragstart="dragstartHandler($event)"> </div>
        <div class="piece kw square-60" draggable="true" @dragstart="dragstartHandler($event)"> </div>
        <div class="piece qw square-59" draggable="true" @dragstart="dragstartHandler($event)"> </div>
        <div class="piece rw square-56" draggable="true" @dragstart="dragstartHandler($event)"> </div>
        <div class="piece rw square-63" draggable="true" @dragstart="dragstartHandler($event)"> </div>
        <div class="piece bw square-58" draggable="true" @dragstart="dragstartHandler($event)"> </div>
        <div class="piece bw square-61" draggable="true" @dragstart="dragstartHandler($event)"> </div>
        <div class="piece nw square-62" draggable="true" @dragstart="dragstartHandler($event)"> </div>
        <div class="piece nw square-57" draggable="true" @dragstart="dragstartHandler($event)"> </div>
        <div class="piece pw square-48" draggable="true" @dragstart="dragstartHandler($event)"> </div>
        <div class="piece pw square-49" draggable="true" @dragstart="dragstartHandler($event)"> </div>
        <div class="piece pw square-50" draggable="true" @dragstart="dragstartHandler($event)"> </div>
        <div class="piece pw square-51" draggable="true" @dragstart="dragstartHandler($event)"> </div>
        <div class="piece pw square-52" draggable="true" @dragstart="dragstartHandler($event)"> </div>
        <div class="piece pw square-53" draggable="true" @dragstart="dragstartHandler($event)"> </div>
        <div class="piece pw square-54" draggable="true" @dragstart="dragstartHandler($event)"> </div>
        <div class="piece pw square-55" draggable="true" @dragstart="dragstartHandler($event)"> </div>
    </div>
</template>

<style scoped>
.board {
    display: grid;
    grid-template-columns: repeat(8, 1fr);
    gap: 0;
    position: relative;
}

.drag-hover {
    border: var(--square-hover-border-width) var(--square-hover-border-type) var(--square-hover-border-color);
}

.hide {
    display: none;
    transform: translate(1000%, 1000%);
}

.square {
    background-color: var(--color-background);
    height: var(--square-size);
    width: var(--square-size);
    display: flex;
    align-items: center;
    justify-content: center;
}

.square.light {
    background-color: var(--light-square);
}

.square.dark {
    background-color: var(--dark-square);
}

.square-0 {
    transform: translate(0%, -800%)
}

.square-1 {
    transform: translate(100%, -800%)
}

.square-2 {
    transform: translate(200%, -800%)
}

.square-3 {
    transform: translate(300%, -800%)
}

.square-4 {
    transform: translate(400%, -800%)
}

.square-5 {
    transform: translate(500%, -800%)
}

.square-6 {
    transform: translate(600%, -800%)
}


.square-7 {
    transform: translate(700%, -800%)
}

.square-8 {
    transform: translate(0%, -700%)
}

.square-9 {
    transform: translate(100%, -700%)
}

.square-10 {
    transform: translate(200%, -700%)
}

.square-11 {
    transform: translate(300%, -700%)
}

.square-12 {
    transform: translate(400%, -700%)
}

.square-13 {
    transform: translate(500%, -700%)
}

.square-14 {
    transform: translate(600%, -700%)
}

.square-15 {
    transform: translate(700%, -700%)
}

.square-16 {
    transform: translate(0%, -600%)
}

.square-17 {
    transform: translate(100%, -600%)
}

.square-18 {
    transform: translate(200%, -600%)
}

.square-19 {
    transform: translate(300%, -600%)
}

.square-20 {
    transform: translate(400%, -600%)
}

.square-21 {
    transform: translate(500%, -600%)
}

.square-22 {
    transform: translate(600%, -600%)
}

.square-23 {
    transform: translate(700%, -600%)
}

.square-24 {
    transform: translate(0%, -500%)
}

.square-25 {
    transform: translate(100%, -500%)
}

.square-26 {
    transform: translate(200%, -500%)
}

.square-27 {
    transform: translate(300%, -500%)
}

.square-28 {
    transform: translate(400%, -500%)
}

.square-29 {
    transform: translate(500%, -500%)
}

.square-30 {
    transform: translate(600%, -500%)
}

.square-31 {
    transform: translate(700%, -500%)
}

.square-32 {
    transform: translate(0%, -400%)
}

.square-33 {
    transform: translate(100%, -400%)
}

.square-34 {
    transform: translate(200%, -400%)
}

.square-35 {
    transform: translate(300%, -400%)
}

.square-36 {
    transform: translate(400%, -400%)
}

.square-37 {
    transform: translate(500%, -400%)
}

.square-38 {
    transform: translate(600%, -400%)
}

.square-39 {
    transform: translate(700%, -400%)
}

.square-40 {
    transform: translate(0%, -300%)
}

.square-41 {
    transform: translate(100%, -300%)
}

.square-42 {
    transform: translate(200%, -300%)
}

.square-43 {
    transform: translate(300%, -300%)
}

.square-44 {
    transform: translate(300%, -300%)
}

.square-45 {
    transform: translate(500%, -300%)
}

.square-46 {
    transform: translate(600%, -300%)
}

.square-47 {
    transform: translate(700%, -300%)
}

.square-48 {
    transform: translate(0%, -200%)
}

.square-49 {
    transform: translate(100%, -200%)
}

.square-50 {
    transform: translate(200%, -200%)
}

.square-51 {
    transform: translate(300%, -200%)
}

.square-52 {
    transform: translate(400%, -200%)
}

.square-53 {
    transform: translate(500%, -200%)
}

.square-54 {
    transform: translate(600%, -200%)
}

.square-55 {
    transform: translate(700%, -200%)
}

.square-56 {
    transform: translate(0%, -100%)
}

.square-57 {
    transform: translate(100%, -100%)
}

.square-58 {
    transform: translate(200%, -100%)
}

.square-59 {
    transform: translate(300%, -100%)
}

.square-60 {
    transform: translate(400%, -100%)
}

.square-61 {
    transform: translate(500%, -100%)
}

.square-62 {
    transform: translate(600%, -100%)
}

.square-63 {
    transform: translate(700%, -100%)
}

.piece {
    cursor: grab;
    width: var(--square-size);
    height: var(--square-size);
    border: none;
    background-repeat: no-repeat;
    background-size: cover;
    position: absolute;
}

.kb {
    background-image: url("/pieces/kb.svg");
}

.kw {
    background-image: url("/pieces/kw.svg");
}


.rb {
    background-image: url("/pieces/rb.svg");
}

.rw {
    background-image: url("/pieces/rw.svg");
}

.qb {
    background-image: url("/pieces/qb.svg");
}

.qw {
    background-image: url("/pieces/qw.svg");
}

.nb {
    background-image: url("/pieces/nb.svg");
}

.nw {
    background-image: url("/pieces/nw.svg");
}

.bw {
    background-image: url("/pieces/bw.svg");
}

.bb {
    background-image: url("/pieces/bb.svg");
}

.pw {
    background-image: url("/pieces/pw.svg");
}

.pb {
    background-image: url("/pieces/pb.svg");
}

.unselectable {
    user-select: none;
    -moz-user-select: none;
    -webkit-user-drag: none;
    -webkit-user-select: none;
    -ms-user-select: none;
}
</style>
