import axios from 'axios'

const defaultTime = 10
const defaultIncrement = 0

// Get the PlayerID, or retrieve a new PlayerID from the server
export function getPlayerID() {
    let playerID = localStorage.getItem("playerID")
    if (playerID == null) {
        axios.post("/token").then(response => {
            playerID = response.data
            localStorage.setItem("playerID", playerID)
        }).catch(err => {
            console.log("error", err)
        })
    }
    return playerID
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
