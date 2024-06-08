#!/bin/bash

func() {
    echo "端口监听脚本(Naabu)"
    echo "Usage:"
    echo "ports.sh [-f File]"
    echo "Description:"
    echo "-f 指定监听文件 example: -f /path/to/output/subs.txt"
    exit 1
}

function programExists() {
    local ret='0'
    command -v $1 >/dev/null 2>&1 || { local ret='1'; }

    # fail on non-zero return value
    if [[ "$ret" -ne 0 ]]; then
        echo -e "\033[31m[Error]命令不存在:$1 \033[0m"
        return 1
    fi
    echo -e "\033[32m[Success]命令存在:$1 \033[0m"
    return 0
}

# 检查必需的程序
programExists naabu || exit 1
programExists notify || exit 1

# 获取参数
while getopts ':hf:' OPT; do
    case $OPT in
        f) File="$OPTARG";;
        h) func;;
        \?) func;;
    esac
done

# 如果未指定文件，显示用法
if [ -z "$File" ]; then
    func
fi

# 如果文件不存在，退出
if [ ! -f "$File" ]; then
    echo -e "\033[31m[Error]监听文件不存在: $File \033[0m"
    exit 1
fi

# 保存端口状态的文件
STATE_FILE="${File}.state"

# 执行扫描的函数
perform_scan() {
    naabu -silent -list "$File" -o current_scan.txt
}

# 发送通知的函数
send_notification() {
    local message="$1"
    echo "$message" | notify -bulk -id ports
}

# 获取当前扫描结果
perform_scan

# 初始化状态文件
if [ ! -f "$STATE_FILE" ]; then
    touch "$STATE_FILE"
fi

# 比较并更新状态文件
new_ports=$(comm -23 <(sort current_scan.txt) <(sort "$STATE_FILE"))
closed_ports=$(comm -13 <(sort current_scan.txt) <(sort "$STATE_FILE"))

# 解析新开放端口并分类
echo "$new_ports" | while IFS= read -r line; do
    domain=$(echo "$line" | cut -d ':' -f 1)
    port=$(echo "$line" | cut -d ':' -f 2)
    echo "$domain 新增 $port" >> new_ports.tmp
done

# 解析关闭端口并分类
echo "$closed_ports" | while IFS= read -r line; do
    domain=$(echo "$line" | cut -d ':' -f 1)
    port=$(echo "$line" | cut -d ':' -f 2)
    echo "$domain 减少 $port" >> closed_ports.tmp
done

# 生成通知消息
message=""
for domain in $(awk '{print $1}' new_ports.tmp closed_ports.tmp | sort | uniq); do
    message+="[$domain]:\n"
    new_ports_list=$(grep "^$domain 新增" new_ports.tmp | awk '{print $3}' | paste -sd ',' -)
    closed_ports_list=$(grep "^$domain 减少" closed_ports.tmp | awk '{print $3}' | paste -sd ',' -)

    [ -n "$new_ports_list" ] && message+="新增：$new_ports_list\n"
    [ -n "$closed_ports_list" ] && message+="减少：$closed_ports_list\n"
done

if [ -n "$message" ]; then
    send_notification "$message"
fi

# 更新状态文件
if [ -n "$new_ports" ]; then
    echo "$new_ports" >> "$STATE_FILE"
fi

if [ -n "$closed_ports" ]; then
    grep -v -F -x -f <(echo "$closed_ports") "$STATE_FILE" > temp_state && mv temp_state "$STATE_FILE"
fi

# 清理临时文件
rm current_scan.txt new_ports.tmp closed_ports.tmp
