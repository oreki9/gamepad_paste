#!/bin/bash
if [ "$1" = "copy" ]; then
    sleep 0.5
    xclip -o -selection primary | xclip -selection clipboard
elif [ "$1" = "paste" ]; then
    sleep 0.5
    getcopystr=$(xclip -selection clipboard -o)
    xdotool type "$getcopystr"
elif [ "$1" = "getpaste" ]; then
    getcopystr=$(xclip -selection clipboard -o)
    echo "$getcopystr"
fi
