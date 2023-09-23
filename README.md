# rrt-bot
A discord bot with some random capabilities.

## Running
```sh
git clone https://github.com/johannesHHM/rrt-bot
cd rrt-bot
go run main.go
```
Note: rrt-bot/ must be in your GOPATH

## Config
The bot is configured using environment variables.
To easily keep track of them, create a `.env` file using this template:
```sh
export BOT_TOKEN="[discord bot token]"
export BOT_MAP_DIR="[map directory]"
export BOT_PREFIX="!"
```

## Commands
| Command              | Description                                                   |
|----------------------|---------------------------------------------------------------|
| !upmap [map files]   | uploads all given maps to BOT_MAP_DIR                         |
| !ping                | pongs                                                         |
