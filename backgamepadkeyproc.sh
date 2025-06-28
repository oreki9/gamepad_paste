#!/bin/bash
gamepadWindow="main"
if [ -e /dev/input/event6 ]; then
    sudo evtest /dev/input/event6 | grep --line-buffered 'code 315.*value 1' | while read line; do
        echo "KEY 315 PRESSED!"
        matches=$(ps aux | grep -i "$gamepadWindow" | grep -v grep)
        checkprocess=$(echo "$matches" | wc -l)
        if [ "$checkprocess" -eq 1 ]; then #process not found
            kdocker "$gamepadWindow"
        fi
        # put your command here, e.g.:
        # notify-send "Key 315 pressed"
    done
fi
