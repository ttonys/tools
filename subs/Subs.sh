#!/bin/bash


func() {
    echo "子域名扫描脚本(ASN/Brute/CRT)"
    echo "Usage:"
    echo "Subs.sh [-d Domain name] [-a ASN] [-w Wrod List] -f -x..."
    echo "Description:"
    echo "-d 指定域名 example: -d example.com"
    echo "-f 指定子域名文件(开启-f不执行子域名挖掘) example: -f sub.txt"
    echo "-a 指定ASN号(https://bgp.he.net/) example: -a AS714"
    echo "-w 指定域名爆破字典 example: -w ~/subdomains-top1000.txt"
    echo "-x 执行子域名fuzz example: -x"
    exit -1
}

ASN="null"
Start=flase
SubFuzz=false
DomainScan=false
WordList="/Users/sys71m/Tools/wordlist/SecLists/Discovery/DNS/subdomains-top1million-110000.txt"


function programExists() {
    local ret='0'
    command -v $1 >/dev/null 2>&1 || { local ret='1'; }

    # fail on non-zero return value
    if [[ "$ret" -ne 0 ]]; then
        Start=true
        echo -e "\033[31m[Error]命令不存在:$1 \033[0m"
        return 1
    fi
    echo -e "\033[32m[Success]命令存在:$1 \033[0m"
    return 0
}


# 获取参数
while getopts 'h:a:d:w:f:x' OPT; do
    case $OPT in
        a) ASN="$OPTARG";;
        d) Domain="$OPTARG";;
        w) WordList="$OPTARG";;
        f) SubFile="$OPTARG";;
        x) SubFuzz=true;;
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
# programExists assetfinder
programExists anew
programExists ./findomain
programExists ./puredns
echo -e "*****结束执行参数检查*****\n"

# 命令不全不执行shell
if [[ $Start == false ]]; then
    exit -1
fi

# Domain必须指定
if [[ $Domain == "" ]]; then
    echo -e "\033[31m[Fatal]必须指定Domain:-d example.com \033[0m"
    exit -1
fi

# 当制定subfile时, 不执行子域名查找
if [[ $SubFile == "" ]]; then
    SubFile="./subdomains/$Domain.subs.txt"
    DomainScan=true 
fi


# 输出参数信息
echo -e "\033[32m[ASN]        $ASN \033[0m"
echo -e "\033[32m[Domain]     $Domain \033[0m"
echo -e "\033[32m[SubFuzz]    $SubFuzz \033[0m"
echo -e "\033[32m[WordList]   $WordList \033[0m"
echo -e "\033[32m[SubFile]    $SubFile \033[0m"
echo -e "\033[32m[DomainScan] $DomainScan \033[0m"
echo -e "\n"


# 通过ASN查找域名
if [[ $ASN != "null" && $DomainScan == true ]]; then
    echo -e "*****ASN*****"
    echo -e "[Shell] whois -h whois.radb.net  -- '-i origin $ASN' | grep -Eo \"([0-9.]+){4}/[0-9]+\" | uniq | mapcidr -silent | dnsx -ptr -resp-only"
    whois -h whois.radb.net  -- "-i origin $ASN" | grep -Eo "([0-9.]+){4}/[0-9]+" | uniq | mapcidr -silent | dnsx -ptr -resp-only | tee ./result/$ASN.txt
    echo -e "*****ASN*****\n"
fi


# amass暴力破解域名
if [[ $Domain != "" && $DomainScan == true ]]; then
    echo -e "*****开始执行Amass遍历*****"
    amass enum -active -d $Domain -brute -w $WordList -o ./result/$Domain.amass.txt
    echo -e "*****结束执行Amass遍历*****\n"
fi


# subfinder查找域名
if [[ $Domain != "" && $DomainScan == true ]]; then
    echo -e "*****开始执行Subfinder*****"
    subfinder -d $Domain -o ./result/$Domain.subfinder.txt
    echo -e "*****结束执行Subfinder*****\n"
fi


# https://crt.sh查找域名
if [[ $Domain != "" && $DomainScan == true ]]; then
    echo -e "*****开始执行CTRF(https://crt.sh)*****"
    python ctrf.py -d $Domain -o ./result/$Domain.ctrf.txt
    echo -e "*****结束执行CTRF(https://crt.sh)*****\n"
fi


# assetfinder查找域名
# if [[ $Domain != "" && $DomainScan == true ]]; then
#     echo -e "*****开始执行Assetfinder*****"
#     assetfinder --subs-only $Domain | tee ./result/$Domain.assetfinder.txt
#     echo -e "*****结束执行Assetfinder*****\n"
# fi


# findomain查找域名
if [[ $Domain != "" && $DomainScan == true ]]; then
    echo -e "*****开始执行Findomain*****"
    ./findomain --quiet -t $Domain -u ./result/$Domain.findomain.txt
    echo -e "*****结束执行Findomain*****\n"
fi


# 结果去重
if [[ $Domain != "" && $DomainScan == true ]]; then
    echo -e "*****开始执行去重[ASN/Amass/Subfinder/CTRF/Assetfinder/Findomain]*****"
    # cat ./result/$ASN.txt \
    # ./result/$Domain.amass.txt \
    # ./result/$Domain.subfinder.txt \
    # ./result/$Domain.ctrf.txt \
    # ./result/$Domain.assetfinder.txt \
    # ./result/$Domain.findomain.txt > ./result/$Domain.subs.unsort.txt
    cat ./result/$Domain.subfinder.txt \
    ./result/$Domain.amass.txt \
    ./result/$Domain.ctrf.txt \
    ./result/$Domain.findomain.txt | anew $SubFile
    echo -e "\033[32m[Success]执行子域名挖掘结束, 子域名保存位置: $SubFile \033[0m"
    echo -e "*****结束执行去重[ASN/Amass/Subfinder/CTRF/Assetfinder/Findomain]*****\n"
fi


# Gotator查找域名
if [[ $Domain != "" && $SubFuzz == true ]]; then
    echo -e "*****开始执行Gotator[排列组合]*****"
    ./gotator -sub $SubFile -perm permutations_list.txt -depth 1 -numbers 10 -mindup -adv -md > ./result/$Domain.gotator.txt
    echo -e "\033[32m文件位置:./result/$Domain.gotator.txt \033[0m"
    echo -e "*****结束执行Gotator[排列组合]*****\n"
fi


# Regulator查找域名
if [[ $Domain != "" && $SubFuzz == true ]]; then
    echo -e "*****开始执行Regulator[排列组合]*****"
    python generator_rules.py -t $Domain -f $SubFile -o "./result/$Domain.regulator.rules.txt"
    if [[ $? -ne 0 ]]; then
        exit -1
    fi
    cat ./result/$Domain.regulator.rules.txt | python generator_urls.py | sed -E 's/\.{2,}/./g' | sort -fu | grep -vE '(\._|_\.|\-\.|\.\-|_\-|\-_)' > ./result/$Domain.regulator.txt
    echo -e "\033[32m文件位置:./result/$Domain.regulator.txt \033[0m"
    echo -e "*****结束执行Regulaor[排列组合]*****\n"
fi


# puredns验证
if [[ $Domain != "" && $SubFuzz == true ]]; then
    echo -e "*****开始执行Puredns验证[排列组合]*****"
    cat ./result/$Domain.gotator.txt ./result/$Domain.regulator.txt > ./result/$Domain.puredens.txt
    ./puredns resolve ./result/$Domain.puredens.txt --write ./subdomains/$Domain.subs.new.txt
    cat ./subdomains/$Domain.subs.new.txt $SubFile | sort | uniq | tee ./subdomains/$Domain.final.txt
    echo -e "\033[32m[Success]执行子域名Fuzz结束, 保存位置: ./subdomains/$Domain.final.txt \033[0m"
    diff $SubFile ./subdomains/$Domain.final.txt
    echo -e "*****结束执行Puredns验证[排列组合]*****\n"
fi