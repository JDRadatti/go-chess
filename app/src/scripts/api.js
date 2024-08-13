import axios from 'axios'
import { ref } from 'vue'

const defaultTime = 10
const defaultIncrement = 0
const playerID = ref("")

// Get the PlayerID, or retrieve a new PlayerID from the server
export function getPlayerID() {
    if (playerID.value != "") {
        return playerID.value
    }
    if (localStorage.getItem("playerID") != null) {
        playerID.value = localStorage.getItem("playerID")
        return playerID.value
    }

    return playerID.value
}

export function setPlayerID() {
    axios.post("/token").then(response => {
        localStorage.setItem("playerID", response.data)
    }).catch(err => {
        console.log("error", err)
    })
}


export async function startGame(time, increment) {
    let playerID = getPlayerID()
    return axios.post('/play', {
        playerID: playerID,
        time: ((time) ? time : defaultTime),
        increment: ((increment) ? increment : defaultIncrement),
    }).then(response => {
        return response.data
    }).catch(error => {
        console.log("ERROR: ", error)
    })
}
