#!/bin/bash
if [ -e maingamepadkey ]; then
    kdocker ./maingamepadkey
else
    go build -o maingamepadkey
    kdocker ./maingamepadkey
fi