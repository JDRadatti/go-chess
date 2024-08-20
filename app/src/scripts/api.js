import axios from 'axios'

const defaultTime = 10
const defaultIncrement = 0

// Get the PlayerID, or retrieve a new PlayerID from the server
export function getPlayerID() {
    return localStorage.getItem("playerID")
}

export function setPlayerID(playerID) {
    localStorage.setItem("playerID", playerID)
}


export async function startGame(time, increment) {
    let playerID = getPlayerID()
    return axios.post('/play', {
        playerID: playerID,
        time: ((time) ? time : defaultTime),
        increment: ((increment) ? increment : defaultIncrement),
    }).then(response => {
        console.log("playerID", response.data.PlayerID)
        setPlayerID(response.data.PlayerID)
        return response.data
    }).catch(error => {
        console.log("ERROR: ", error)
    })
}
