#!/bin/bash

start_time=$(date +%s)
# nohup ./command.sh >/dev/null 2>&1
while true; do
    matches=$(ps aux | grep -i main.go | grep -v grep)
    checkprocess=$(echo "$matches" | wc -l)
    echo "$checkprocess"
    if [ "$checkprocess" -eq 1 ]; then
        getcopystr=$(xclip -selection clipboard -o)
        sleep 2
        xdotool type "$getcopystr"
        echo "hello"
        break
    fi
    # Check if more than 60 seconds have passed
    current_time=$(date +%s)
    elapsed=$((current_time - start_time))
    if [ "$elapsed" -ge 60 ]; then
        echo "Timeout reached: 1 minute"
        break
    fi
    sleep 1  # Optional: add delay to avoid CPU overload
done
