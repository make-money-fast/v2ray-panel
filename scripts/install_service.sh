#!/bin/bash

# cat v2ray-panel.service
replace="__SERVER_EXEC_PATH__"
binpath=$(pwd)/bin/v2ray-panel-server
serviceContent=$(cat scripts/v2ray-panel.service)
res=$(echo "$serviceContent" | sed "s#${replace}#${binpath}#g")
printf "%s" "$res" > /etc/systemd/system/v2ray-panel.service

echo "查看状态: system ctl status v2ray-panel"
echo "开始: system ctl start v2ray-panel"
echo "重启: system ctl reload v2ray-panel"
echo "停止: system ctl stop v2ray-panel"
echo "开机启动: system ctl enable v2ray-panel"