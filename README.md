# reptile

This branch is meant to experiment with servers in GoLang

Project goals:
- Chess.com clone 
- Testable and scallable api
- Server stats

Tech stack: 
- Backend: go
- Frontend: Vue.js
- Possibly dockerize
- Want to try and deploy with aws ec2. Maybe amplify. Try GCP?
- If I get reaallllly far. User accounts?
- integrate different bots? custom and stockfish
- MAYBE s3 to send png's based on theme
- Load balancer and rate limiter?

How to reach goal: 
1. Get UDP server up and running
- will there be one server that is constantly running and then each game creates
a new "game" and the server keeps track of multople games?
- Different thread for each game?
- go I need user accounts or can I just have random users?
- server needs to hold game logic
- two clients per game
- clients need to send moves to the server. "kf2"
- potential problem: how to make the client keep track of things like moving. 
Do I need to implements the game on both the client and the server?
A user can move an element to a square and then if the server decides it's valid 
it will move it on the server and reflect online?
server needs to send time contro?

How fast of games can we play?
2. Handle match making. at first it can be a queue or with a code?
3. Implement chess on top of server. 
    - allow for different types of users (person, bot, etc...)
    - game holds info about the users, board, time control. whos turn it is. etc. STATE
    - board keeps track of the pieces and their positions as well as legal moves
        - has functions like Move(startSquare, endSquare) and legalMoves(square)
        - maps from human readable to index (a1 to 0)
    - piece implements function PossibleMoves() that has a list of all possible moves, even if they are not legal
4. 
  
- build backend first 
- play around with scaling the servers and testing
- try to build out frontend

Starting point. Connect to the serfver with one client and display the time remaining in 
a game. maybe make a button that says move and see if it adds time. game 
ends when time ran out

How do I get a custom url with a unique key for each game?

Lobby:
    lobby []Game (make FIFO queue?)
Game:
    id: 
    address: // address to the game server
    white Player // address to client 1 -- will play with the white pieces
    black Player // address to client 2 -- will play with thye black pieces
    moves: "d4 d5 ... "
    timeStarted: 
    timeEnded:
    status: "waiting", "ongoing", "draw", "white", "black"
    clock increment: time.Time // total time left in the game
    clock initial
Player: 
    addr
    name? string "Anonymous"
    currentGame Game
Lobby: 
    list of games currently being played


Want the games to be exportable via PGN

do I need to add a token as a cookie and then use that to auth a Game. 
i.e. you're only allowed in the url if you have the token 


How do I make custom urls for each gamge
How do I auth each game url
How do I connect to the UDP server from the website?

Steps: 
- press button
- backend creates a game and adds it to the lobby
- creates a new url and traverses there
- new url contains the game state that it gets from the udp server
- for now, just make the clock go
    

Day 1 notes: 
- add basic http and udp servers
- learn about vue and add sample app
- how to use the udp server for the game?

Day 2 notes: 
- realised I needed a websocket. i.e. upgrade the http server to allow
a constant stream of data
- 


- What I need: 
- websockt: creates client, registers to game or creates a game, registers game to lobby if created
    
User - connects server to game
- game pointer
- buffered channel of outbount moves (maybe. premoves)
- connection pointer

Game - connects User to Lobby
White *User
Black *User
Moves channel for inbound moves from client
init game and add game to games pointer. Start go routine to play game.
remove game: close player's move channels, cleans game, and removes pointer from games

Use Reddis to store game state?


Lobby - controls all games (creation, deletion, moves)
- []Game pointers
- max number of games

Lobby should be separate 

Chess watch? 
Build out clans. Live stream chess games. 
Chess Clash
chessclash.com


import { uuid } from 'vue-uuid';



go routine for each game that runs the game. to start, it can just be a 
clock and stream of moves
