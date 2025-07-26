#!/bin/bash
cd "/home/oreki/Documents/Github/gamepad_paste"
if [ -e maingamepadkey ]; then
     nohup ./backgamepadkeyproc.sh
else
    go build -o maingamepadkey
    nohup ./backgamepadkeyproc.sh
fi
