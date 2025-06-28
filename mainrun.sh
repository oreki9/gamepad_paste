#!/bin/bash
if [ -e gamepadkey ]; then
    kdocker ./gamepadkey
else
    go build -o gamepadkey
    kdocker ./gamepadkey
fi