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
        - maps from human readable to index
    - piece implements function PossibleMoves() that has a list of all possible moves, even if they are not legal
4. 
  
- build backend first 
- play around with scaling the servers and testing
- try to build out frontend

