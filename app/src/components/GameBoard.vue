<script setup>
import { ref, onMounted, watch } from 'vue'
import { sendMove } from '../scripts/websocket.js'
import { VueSpinnerBox } from 'vue3-spinners';

const props = defineProps(['start', 'color', 'waiting', 'fen', 'count'])

const waiting = ref(false)
const dragPiece = ref(null);
const boardRef = ref(null);
const pieceClassList = ref([
    ["piece", "qb", "square-3", ""],
    ["piece", "kb", "square-4", ""],
    ["piece", "rb", "square-7", ""],
    ["piece", "rb", "square-0", ""],
    ["piece", "bb", "square-5", ""],
    ["piece", "bb", "square-2", ""],
    ["piece", "nb", "square-6", ""],
    ["piece", "nb", "square-1", ""],
    ["piece", "pb", "square-8", ""],
    ["piece", "pb", "square-9", ""],
    ["piece", "pb", "square-10", ""],
    ["piece", "pb", "square-11", ""],
    ["piece", "pb", "square-12", ""],
    ["piece", "pb", "square-13", ""],
    ["piece", "pb", "square-14", ""],
    ["piece", "pb", "square-15", ""],
    ["piece", "kw", "square-60", ""],
    ["piece", "qw", "square-59", ""],
    ["piece", "rw", "square-56", ""],
    ["piece", "rw", "square-63", ""],
    ["piece", "bw", "square-58", ""],
    ["piece", "bw", "square-61", ""],
    ["piece", "nw", "square-62", ""],
    ["piece", "nw", "square-57", ""],
    ["piece", "pw", "square-48", ""],
    ["piece", "pw", "square-49", ""],
    ["piece", "pw", "square-50", ""],
    ["piece", "pw", "square-51", ""],
    ["piece", "pw", "square-52", ""],
    ["piece", "pw", "square-53", ""],
    ["piece", "pw", "square-54", ""],
    ["piece", "pw", "square-55", ""],
])

const pieceTypeIndex = 1
const pieceSquareIndex = 2
const pieceHideIndex = 3
const start = ref(-1);
const dest = ref(-1);
let dragged = null;

const PieceToType = {
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

const Squares = [
    "a8", "b8", "c8", "d8", "e8", "f8", "g8", "h8",
    "a7", "b7", "c7", "d7", "e7", "f7", "g7", "h7",
    "a6", "b6", "c6", "d6", "e6", "f6", "g6", "h6",
    "a5", "b5", "c5", "d5", "e5", "f5", "g5", "h5",
    "a4", "b4", "c4", "d4", "e4", "f4", "g4", "h4",
    "a3", "b3", "c3", "d3", "e3", "f3", "g3", "h3",
    "a2", "b2", "c2", "d2", "e2", "f2", "g2", "h2",
    "a1", "b1", "c1", "d1", "e1", "f1", "g1", "h1",
]


function updateBoard(fen) {
    const rows = fen.split("/")
    var squareIndex = 0
    var pieceIndex = 0
    for (let i = 0; i < rows.length; i++) { // must be 8 rows
        for (let j = 0; j < rows[i].length; j++) {
            if (!isNaN(rows[i][j])) {
                squareIndex += Number(rows[i][j])
            } else {
                updatePieceType(pieceIndex, rows[i][j])
                movePiece(pieceIndex, squareIndex)
                squareIndex++
                pieceIndex++
            }
        }
    }
    for (let i = pieceIndex; i < pieceClassList.value.length; i++) {
        hidePiece(i)
    }
}

function showAllPieces() {
    for (let i = 0; i < pieceClassList.value.length; i++) {
        showPiece(i)
    }
}

function hideAllPieces() {
    for (let i = 0; i < pieceClassList.value.length; i++) {
        hidePiece(i)
    }
}

function updatePieceType(pieceID, pieceIcon) {
    pieceClassList.value[pieceID][pieceTypeIndex] = PieceToType[pieceIcon]
}

function showPiece(pieceID) {
    pieceClassList.value[pieceID][pieceHideIndex] = ""
}

function hidePiece(pieceID) {
    pieceClassList.value[pieceID][pieceHideIndex] = "hide"
}

function movePiece(pieceID, dest) {
    pieceClassList.value[pieceID][pieceSquareIndex] = "square-" + dest
}

function getPieceSquareIndex(pieceID) {
    return Number(pieceClassList.value[pieceID][pieceSquareIndex].slice(7))
}

function dragstartHandler(ev) {
    if (dragged != null) {
        return
    }

    dragged = ev.target;
    hidePiece(dragged.id);

    start.value = getPieceSquareIndex(dragged.id);

    // Move dragPiece to the cursor position. 
    // This is because default drag creates an undesired opacity 
    if (dragPiece.value.classList[dragPiece.value.classList.length - 1] == "hide") {
        dragPiece.value.classList.add(dragged.classList[1]);
        dragPiece.value.classList.remove(dragPiece.value.classList[dragPiece.value.classList.length - 2]);
        dragPiece.value.style.setProperty('--cursor-y', ev.pageY - dragPiece.value.offsetWidth / 2 + "px");
        dragPiece.value.style.setProperty('--cursor-x', ev.pageX - dragPiece.value.offsetWidth / 2 + "px");
    }
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
    movePiece(dragged.id, dest.value)
    showPiece(dragged.id)
    dragged = null;

    // ASK THE SERVER
    sendMove(Squares[start.value] + Squares[dest.value]);

}

function dragendHandler(ev) {
    if (dragged != null) {
        showPiece(dragged.id)
        dragged = null;
    }

    dragPiece.value.classList.remove(dragPiece.value.classList[dragPiece.value.classList.length - 1]);
    dragPiece.value.classList.add("hide");
}

function dragoverHandler(ev) {
    // Note: the piece should always be the last element in the dragPiece's classList
    dragPiece.value.style.setProperty('--cursor-y', ev.pageY - dragPiece.value.offsetWidth / 2 + "px");
    dragPiece.value.style.setProperty('--cursor-x', ev.pageX - dragPiece.value.offsetWidth / 2 + "px");
}

function captureHandler(ev) {

    ev.preventDefault();
    dest.value = getPieceSquareIndex(ev.target.id)

    // ASK THE SERVER
    sendMove(Squares[start.value] + Squares[dest.value]);

    showPiece(dragged.id)
    dragged = null
}

onMounted(() => {
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
    hideAllPieces()
});

watch(props, (props) => {

    if (props.fen) {
        updateBoard(props.fen)
    }

    if (props.start && waiting.value == true) {
        showAllPieces()
        waiting.value = false // stop spinner
    }

    if (props.waiting) {
        waiting.value = true // start spinner
    }
})
</script>

<template>
    <div inert class="drag unselectable hide" draggable="false" ref="dragPiece"></div>
    <div class="board-container" @dragover="dragoverHandler($event)" @dragenter.prevent @dragover.prevent>
        <div class="spinner-container" v-if="waiting">
            <VueSpinnerBox size="100" color="rgba(132, 118, 186, 1)" />
            <p> Waiting for Opponenet... </p>
        </div>
        <div class="board" draggable="false">
            <div class="square unselectable" v-for="n in 64" draggable="false" @dragenter="dragenterHandler($event)"
                @dragleave="dragleaveHandler($event)" @drop="dropHandler($event, n)" @dragover.prevent
                @dragenter.prevent>
            </div>
        </div>
        <div v-for="(classList, index) in pieceClassList" :class="classList" :key="index" :id="index" draggable="true"
            @dragstart="dragstartHandler($event)" @dragend="dragendHandler($event)" @drop="captureHandler($event)"
            @dragover.prevent @dragenter.prevent>
        </div>
    </div>
</template>

<style scoped>
.spinner-container {
    display: flex;
    justify-content: center;
    align-items: center;
    flex-direction: column;
    padding: 2rem;
    background-color: var(--color-background);
    border: var(--color-background);
    position: absolute;
    width: calc(var(--square-size) * 4);
    height: calc(var(--square-size) * 4);
    top: 25%;
    left: 28%;
}

.spinner-container p {
    color: var(--light-square);
    text-align: center;
}

.board-container {
    position: relative;
    padding: 0 2rem;
}

.board {
    display: grid;
    grid-template-columns: repeat(8, 1fr);
    gap: 0;
    width: var(--board-conainter-size);
    height: var(--board-container-size);
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
    transform: translate(400%, -300%)
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


.drag {
    width: var(--square-size);
    height: var(--square-size);
    border: none;
    background-repeat: no-repeat;
    background-size: cover;
    position: absolute;
    z-index: 1;
    left: var(--cursor-x);
    top: var(--cursor-y);
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
