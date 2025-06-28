sudo evtest /dev/input/event6 | grep --line-buffered 'code 315.*value 1' | while read line; do
    echo "KEY 315 PRESSED!"
    # put your command here, e.g.:
    # notify-send "Key 315 pressed"
done
