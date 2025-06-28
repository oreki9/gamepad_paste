#!/bin/bash
if [ "$EUID" -ne 0 ]; then
  echo "Please run this script as root (e.g., with sudo)"
  exit 1
else
    if [ -e gamepadkey ]; then
        sudo kdocker ./gamepadkey
    else
        go build -o gamepadkey
        sudo kdocker ./gamepadkey
    fi
fi
