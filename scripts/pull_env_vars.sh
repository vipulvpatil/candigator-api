#!/bin/zsh

# var definitions
SSH_ADDR="root@$1"

# copy .env from server root
scp $SSH_ADDR:~/.env .env.downloaded
