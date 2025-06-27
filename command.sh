#!/bin/bash

start_time=$(date +%s)
# nohup ./command.sh >/dev/null 2>&1
while true; do
    matches=$(ps aux | grep -i main.go | grep -v grep)
    checkprocess=$(echo "$matches" | wc -l)
    touch "$checkprocess" -eq 1
    if [ "$checkprocess" -eq 1 ]; then
        xdotool key --clearmodifiers ctrl+v
        # echo "hello"
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
