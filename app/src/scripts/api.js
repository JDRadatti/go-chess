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
    if (localStorage.getItem("playerID", playerID) != null) {
        playerID.value = localStorage.getItem("playerID", playerID.value)
        return playerID.value
    }


    axios.post("/token").then(response => {
        playerID.value = response.data
        localStorage.setItem("playerID", playerID.value)
    }).catch(err => {
        console.log("error", err)
    })

    return playerID.value
};


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
