#!/bin/zsh

echo $1
# var definitions
SSH_ADDR="root@$1"

# first download the existing .env from server root
scp $SSH_ADDR:~/.env .env.downloaded

# copy .env to server root
scp .env $SSH_ADDR:~/
