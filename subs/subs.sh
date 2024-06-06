#!/bin/bash

func() {
    echo "子域名扫描脚本(Subfinder/Dnsx/Alterx/...)"
    echo "Usage:"
    echo "Subs.sh [-d Domain name] [-w Word List] [-o Output Directory]"
    echo "Description:"
    echo "-d 指定域名 example: -d example.com"
    echo "-w 指定域名爆破字典 example: -w ~/subdomains-top1000.txt"
    echo "-o 指定文件保存路径 example: -o /path/to/output"
    exit 1
}

#WordList="/Users/sys71m/Tools/wordlist/SecLists/Discovery/DNS/subdomains-top1million-110000.txt"
WordList="/Users/sys71m/Tools/wordlist/SecLists/Discovery/DNS/subdomains-top1million-5000.txt"

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

# 获取参数
while getopts 'h:a:d:w:f:o:x' OPT; do
    case $OPT in
        d) Domain="$OPTARG";;
        w) WordList="$OPTARG";;
        o) OutputDir="$OPTARG";;
        h) func;;
        ?) func;;
    esac
done

# 检查命令是否存在
echo -e "*****开始执行参数检查*****"
programExists subfinder
programExists dnsx
programExists dnsgen
programExists anew
programExists alterx
programExists chaos
echo -e "*****结束执行参数检查*****\n"

# Domain必须指定
if [[ $Domain == "" ]]; then
    echo -e "\033[31m[Fatal]必须指定Domain:-d example.com \033[0m"
    exit 1
fi

# OutputDir必须指定
if [[ $OutputDir == "" ]]; then
    echo -e "\033[31m[Fatal]必须指定保存路径:-o /tmp/res \033[0m"
    exit 1
fi

# 创建输出目录
mkdir -p $OutputDir

# 输出参数信息
echo -e "\033[32m[Domain]     $Domain \033[0m"
echo -e "\033[32m[WordList]   $WordList \033[0m"
echo -e "\033[32m[OutputDir]  $OutputDir \033[0m"
echo -e "\n"


# subfinder查找域名
# https://github.com/projectdiscovery/subfinder
# subfinder -d $Domain -o $OutputDir/$Domain.subfinder.txt
echo -e "*****开始执行Subfinder*****"
subfinder -d $Domain | anew $OutputDir/subs.txt
echo -e "*****结束执行Subfinder*****\n"


# https://crt.sh查找域名
# go install github.com/ttonys/tools/crt@latest
echo -e "*****开始执行crt(https://crt.sh)*****"
crt -d $Domain | anew $OutputDir/subs.txt
echo -e "*****结束执行crt(https://crt.sh)*****\n"


# chaos
# https://github.com/projectdiscovery/chaos-client
#chaos -d $Domain -o $OutputDir/$Domain.chaos.txt
echo -e "*****开始执行Chaos*****"
chaos -d $Domain | anew $OutputDir/subs.txt
echo -e "*****结束执行Chaos*****\n"

# dnsx enum
# https://github.com/projectdiscovery/dnsx
# dnsx -silent -d $Domain -w $WordList -o $OutputDir/$Domain.dnsx.txt
echo -e "*****开始执行Dnsx*****"
dnsx -t 10 -rl 10 -stats -silent -d $Domain -w $WordList | anew $OutputDir/subs.txt
echo -e "*****结束执行Dnsx*****\n"


# dnsgen + alterx -> enum
# https://github.com/AlephNullSK/dnsgen
# https://github.com/projectdiscovery/alterx
echo -e "*****开始执行(dnsgen + alterx)*****"
cat $OutputDir/subs.txt | dnsgen - | anew -q $OutputDir/subs.enum.txt
cat $OutputDir/subs.txt | alterx | anew -q $OutputDir/subs.enum.txt
cat $OutputDir/subs.enum.txt | dnsx -silent -t 10 -rl 10 -stats | anew $OutputDir/subs.txt
echo -e "*****结束执行(dnsgen + alterx)*****\n"


echo -e "\033[32m[Success]执行子域名挖掘结束, 子域名保存位置: $OutputDir/subs.txt \033[0m"

echo -e "\033[32m执行Slack通知 \033[0m"
cat $OutputDir/subs.txt | anew $OutputDir/subs.notify.txt | notify -bulk