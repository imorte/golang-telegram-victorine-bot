# golang-telegram-victorine-bot


### This is a golang realization bot https://github.com/synnz/telegram-victorine-bot

#### - How to use?
1) Create a file **bot_token.go** into project root directory
2) Put into this code, and get api token here *https://telegram.me/botfather*
```GO
 package main
 
 var (
  TOKEN = "Your telegram bot api token"
 )
 
```

#### - Dependency
```
go get github.com/mattn/go-sqlite3
go get gopkg.in/telegram-bot-api.v4
```

#### - How to put in container
```
# Dockerfile
FROM golang:1.8

WORKDIR /go/src/app
COPY . .

RUN go-wrapper download 
RUN go-wrapper install 

CMD ["go-wrapper", "run"]
```

#### - Developers
[synnz](https://github.com/synnz), [LikiPiki](https://github.com/LikiPiki) and [endorphin82](https://github.com/endorphin82)
