#!/bin/bash
mainrun(){
    gamepadWindow="maingamepadkey"
    mainproc=$(ps aux | grep -i "/bin/bash backgamepadkeyproc.sh" | grep -v grep)
    checkmainprco=$(echo "$mainproc" | wc -l)
    if [ "$checkmainprco" -eq 1 ]; then
        if [ -e /dev/input/event6 ]; then
            evtest /dev/input/event6 | grep --line-buffered 'code 315.*value 1' | while read line; do
                echo "KEY 315 PRESSED!"
                matches=$(ps aux | grep -i "./$gamepadWindow" | grep -v grep)
#                 checkprocess=$(echo "$matches" | wc -l)
#                 echo "$checkprocess"
                if [ -z "$matches" ]; then #process not found
                    nohup "./$gamepadWindow"
                fi
                # put your command here, e.g.:
                # notify-send "Key 315 pressed"
            done
        else
            echo "device is not found"
        fi
    else
        if [ "$checkmainprco" -gt 1 ]; then
            themostold=$(echo "$mainproc" | awk 'NR > 1 { if ($10 > max) { max = $10; line = $0 } } END { print line }' | awk '{print $2}')
            countoldproc=$(echo "$themostold" | wc -l)
            if [ "$countoldproc" -eq 1 ]; then
                kill "$themostold"
                mainrun
            fi
        fi
    fi
}
mainrun
