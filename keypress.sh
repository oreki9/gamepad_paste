#!/bin/bash
matches=$(ps aux | grep -i backgamepadkeyproc.sh | grep -v grep)
themostold=$(echo "$matches" | awk 'NR > 1 { if ($10 > max) { max = $10; line = $0 } } END { print line }' | awk '{print $2}')
countoldproc=$(echo "$themostold" | wc -l)
if [ "$countoldproc" -eq 1 ]; then
    isWriteSomeThing=$(lsof -p "$themostold" | grep /dev)
    echo "$isWriteSomeThing"
fi