#!/bin/bash

# cat v2ray-panel.service
replace="__SERVER_EXEC_PATH__"
binpath=$(pwd)/bin/v2ray-panel-server
serviceContent=$(cat scripts/v2ray-panel.service)
res=$(echo "$serviceContent" | sed "s#${replace}#${binpath}#g")
printf "%s" "$res" > /etc/systemd/system/v2ray-panel.service