# JSON API experimental implementation


Test application
```
go test ./...
```


Run application
```
go run cmd/server/main.go
```

Request Balance
```
curl -H "Content-Type: application/json" \
--request POST \
--data '{"jsonrpc":"2.0","method":"getBalance","params":{"callerId":1,"playerName":"player1","currency":"EUR","gameId":"riot"},"id":0}' \
http://localhost:8080
```

Withdraw and deposit
```
curl -H "Content-Type: application/json" \
--request POST \
--data '{"jsonrpc":"2.0","method":"withdrawAndDeposit","params":{"callerId":1,"playerName":"player1","withdraw":0,"deposit":200,"currency":"EUR","transactionRef":"1:UOwGgNHPgq3OkqRE","gameRoundRef":"1wawxl:39","gameId":"riot","reason":"GAME_PLAY_FINAL","sessionId":"qx9sgvvpihtrlug","spinDetails":{"betType":"spin","winType":"standart"}},"id":0}' \
http://localhost:8080
```


Rollback 
```
curl -H "Content-Type: application/json" \
--request POST \
--data '{"jsonrpc":"2.0","method":"rollbackTransaction","params":{"callerId":1,"playerName":"player1","transactionRef":"1:UOwGgNHPgq3OkqRE","gameId":"riot","sessionId":"qx9sgvvpihtrlug","gameRoundRef":"1wawxl:39"},"id":0}' \
http://localhost:8080
```


