#!/usr/bin/env sh
set -x
set -e

echo "$BUTLER_API_KEY" > ~/.config/itch/butler_creds
wc ~/.config/itch/butler_creds
