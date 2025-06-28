#!/bin/bash
if [ "$EUID" -ne 0 ]; then
  echo "Please run this script as root (e.g., with sudo)"
  exit 1
else
    if [ -e gamepadkey ]; then
        kdocker ./gamepadkey
    else
        go build -o gamepadkey
        kdocker ./gamepadkey
    fi
fi
