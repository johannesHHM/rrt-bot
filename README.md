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
export BOT_URL_HTTP_MASTER="https://master1.ddnet.org/ddnet/15/servers.json"
export BOT_PREFIX="!"
```
Source `.env` after making changes to it.

## Commands
| Command                  | Description                                                   |
|--------------------------|---------------------------------------------------------------|
| ``!help``                | lists all commands                                            |
| ``!upmap [map files]``   | uploads all given maps to BOT_MAP_DIR                         |
| ``!lsmap``               | lists all maps in BOT_MAP_DIR                                 |
| ``!rmmap [map names]``   | removes given maps from BOT_MAP_DIR                           |
| ``!online``              | lists all online RRT members                                  |
| ``!ping``                | pongs                                                         |
