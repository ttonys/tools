#!/bin/bash


func() {
    echo "子域名扫描脚本(ASN/Brute/CRT)"
    echo "Usage:"
    echo "Subs.sh [-d Domain name] [-a ASN] [-w Wrod List]"
    echo "Description:"
    echo "-d 指定域名 example: -d example.com"
    echo "-a 指定ASN号(https://bgp.he.net/) example: -a AS714"
    echo "-w 指定域名爆破字典 example: -w ~/subdomains-top1000.txt"
    exit -1
}

Start="true"
WordList="/usr/share/seclists/Discovery/DNS/subdomains-top1million-110000.txt"


function programExists() {
    local ret='0'
    command -v $1 >/dev/null 2>&1 || { local ret='1'; }

    # fail on non-zero return value
    if [[ "$ret" -ne 0 ]]; then
        Start="false"
        echo -e "\033[31m[Error]命令不存在:$1 \033[0m"
        return 1
    fi
    echo -e "\033[32m[Success]命令存在:$1 \033[0m"
    return 0
}


function execRegulator() {
    echo -e "\033[36m生成规则文件 \033[0m"
    python3 main.py -t $1 -f $2 -o "./result/"$1".rules"
    if [[ $? -ne 0 ]]; then
        exit -1
    fi


    echo -e "\033[36m生成子域名爆破文件 \033[0m"
    ./make_brute_list.sh "./result/"$1".rules" "./result/"$1".brute"


    echo -e "\033[36mDNS验证子域名 \033[0m"
    ./puredns resolve "./result/"$1".brute" --write "./result/"$1".valid"
}


# 获取参数
while getopts 'h:a:d:w:' OPT; do
    case $OPT in
        a) ASN="$OPTARG";;
        d) Domain="$OPTARG";;
        w) WordList="$OPTARG";;
        h) func;;
        ?) func;;
    esac
done

# 检查命令是否存在
echo -e "*****开始执行参数检查*****"
programExists whois
programExists mapcidr
programExists dnsx
programExists amass
programExists massdns
programExists httpx
programExists assetfinder
programExists ./findoamin
echo -e "*****结束执行参数检查*****\n"

# 命令不全不执行shell
if [[ "$Start" == "false" ]]; then
    exit -1
fi


# 输出参数信息
echo -e "\033[32m[ASN]        $ASN \033[0m"
echo -e "\033[32m[Domain]     $Domain \033[0m"
echo -e "\033[32m[WordList]   $WordList \033[0m"
echo -e "\n"


# 通过ASN查找域名
if [[ "$ASN" != "" ]]; then
    echo -e "*****ASN*****"
    echo -e "[Shell] whois -h whois.radb.net  -- '-i origin $ASN' | grep -Eo \"([0-9.]+){4}/[0-9]+\" | uniq | mapcidr -silent | dnsx -ptr -resp-only"
    whois -h whois.radb.net  -- "-i origin $ASN" | grep -Eo "([0-9.]+){4}/[0-9]+" | uniq | mapcidr -silent | dnsx -ptr -resp-only | tee -a ./result/$ASN.txt
    echo -e "*****ASN*****\n"
fi


# amass暴力破解域名
if [[ "$WordList" != "" && "$Domain" != "" ]]; then
    echo -e "*****开始执行Amass遍历*****"
    amass enum -active -d $Domain -brute -w $WordList -o ./result/$Domain.amass.txt
    echo -e "*****结束执行Amass遍历*****\n"
fi


# subfinder查找域名
if [[ "$Domain" != "" ]]; then
    echo -e "*****开始执行Subfinder*****"
    subfinder -d $Domain -o ./result/$Domain.subfinder.txt
    echo -e "*****结束执行Subfinder*****\n"
fi


# https://crt.sh查找域名
if [[ "$Domain" != "" ]]; then
    echo -e "*****开始执行CTRF(https://crt.sh)*****"
    python ctrf.py -d $Domain -o ./result/$Domain.ctrf.txt
    echo -e "*****结束执行CTRF(https://crt.sh)*****\n"
fi


# assetfinder查找域名
if [[ "$Domain" != "" ]]; then
    echo -e "*****开始执行Assetfinder*****"
    assetfinder --subs-only $Domain | tee -a ./result/$Domain.assetfinder.txt
    echo -e "*****结束执行Assetfinder*****\n"
fi


# findoamin查找域名
if [[ "$Domain" != "" ]]; then
    echo -e "*****开始执行Findoamin*****"
    ./findoamin --quiet -t $Domain -u ./result/$Domain.findoamin.txt
    echo -e "*****结束执行Findoamin*****\n"
fi





# 结果去重
# cat ./result/$Domain.subfinder.txt ./result/$Domain.amass.txt > ./result/$Domain.subs.unsort.txt
# sort -n ./result/$Domain.subs.unsort.txt | uniq | tee -a ./result/$Domain.subs.txt