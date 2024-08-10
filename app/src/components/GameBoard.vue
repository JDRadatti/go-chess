<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'

const PieceToIcon = {
    "B": "/pieces/bb.svg",
    "b": "/pieces/bw.svg",
    "R": "/pieces/rb.svg",
    "r": "/pieces/rw.svg",
    "K": "/pieces/kb.svg",
    "k": "/pieces/kw.svg",
    "Q": "/pieces/qb.svg",
    "q": "/pieces/qw.svg",
    "N": "/pieces/nb.svg",
    "n": "/pieces/nw.svg",
    "P": "/pieces/pb.svg",
    "p": "/pieces/pw.svg",
    "": "/pieces/empty.svg",
}

const squares = ref([
    '', '', '', '', '', '', '', '',
    '', '', '', '', '', '', '', '',
    '', '', '', '', '', '', '', '',
    '', '', '', '', '', '', '', '',
    '', '', '', '', '', '', '', '',
    '', '', '', '', '', '', '', '',
    '', '', '', '', '', '', '', '',
    '', '', '', '', '', '', '', '',
])
squares.value[0] = PieceToIcon["P"]

function updateBoard(fen) {
    const rows = fen.split("/")
    var index = 0
    for (let i = 0; i < rows.length; i++) { // must be 8 rows
        for (let j = 0; j < rows[i].length; j++) {
            console.log("PIECE: ", rows[i][j])
            console.log("NUMBER", Number(rows[i][j]))
            if (!isNaN(rows[i][j])) {
                for (let k = Number(rows[i][j]); k > 0; k--) {
                    updateSquare(index, "")
                    index++
                }
            } else {
                console.log("HERE", index, rows[i][j])
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

onMounted(() => {
    updateBoard("RNBQKBNR/PPPPPPPP/8/8/8/8/pppppppp/rnbqkbnr")
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
    console.log(squareElements)
});

</script>

<template>
    <h1 class="green">This is your game board {{ $route.params.id }}</h1>
    <div class="board">
        <div class="square" v-for="piece in squares">
            <img class="piece" :src="piece" />
        </div>
    </div>
</template>

<style scoped>
.board {
    display: grid;
    grid-template-columns: repeat(8, 1fr);
    gap: 0;
}

.square {
    background-color: yellow;
    height: 5vw;
    width: 5vw;
    display: flex;
    align-items: center;
    justify-content: center;
    border: none;
}

.square.light {
    background-color: #f0f1f0;
}

.square.dark {
    background-color: #8476ba;
}

.piece {
    cursor: grab;
    width: 4.5vw;
    height: 4.5vw;
    border: none;
}
</style>
