#!/bin/bash


func() {
    echo "目前仅支持GNU/Linux平台"
    echo "Usage:"
    echo "regulator.sh [-f SubDomain File] [-d Domain name] [-i]"
    echo "Description:"
    echo "-f 指定子域名文件 example: -f sub.txt"
    echo "-d 指定扫描的域名 example: -d baidu"
    echo "-i 自动安装依赖"
    exit -1
}

Install="false"
Start="true"


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



function installProgram() {
    # echo -e "\033[36m下载https://github.com/cramppet/regulator \033[0m"
    # git clone https://github.com/cramppet/regulator


    echo -e "\033[36m安装regulator依赖 \033[0m"
    pip install -r requirements.txt
    if [[ $? -ne 0 ]]; then
        exit -1
    fi


    # echo -e "\033[36m下载resolvers.txt解析文件 \033[0m"
    # wget https://raw.githubusercontent.com/blechschmidt/massdns/master/lists/resolvers.txt
    # if [[ $? -ne 0 ]]; then
    #     exit -1
    # fi


    # echo -e "\033[36m下载https://github.com/d3mondev/puredns \033[0m"
    # git clone https://github.com/d3mondev/puredns


    # echo -e "\033[36m修改puredns, issues:https://github.com/d3mondev/puredns/issues/17 \033[0m"
    #  ====================


    # echo -e "\033[36m编译go build puredns \033[0m"
    # cd ./puredns && go build
    # if [[ $? -ne 0 ]]; then
    #     exit -1
    # fi

    echo -e "\033[36m安装massdns \033[0m"
    git clone https://github.com/blechschmidt/massdns.git
    cd ./massdns && make && make install
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
while getopts 'h:f:d:i' OPT; do
    case $OPT in
        f) SubFile="$OPTARG";;
        d) Domain="$OPTARG";;
        i) Install="true";;
        h) func;;
        ?) func;;
    esac
done

# 检查命令是否存在
echo -e "*****开始执行参数检查*****"
programExists go
programExists git
programExists python
programExists ./puredns
programExists massdns
echo -e "*****结束执行参数检查*****\n"
if [[ "$Start" == "false" && "$Install" == "false" ]]; then
    exit -1
fi


if [[ "$Install" == "true" ]]; then
    echo -e "*****开始安装*****"
    installProgram
    echo -e "*****结束安装*****\n"
    exit -1
fi


if [[ "$SubFile" != "" && $Domain != "" ]]; then

    echo $Domain
    echo $NewSubFile
    echo -e "*****开始执行Regulator*****"
    execRegulator $Domain $SubFile
    echo -e "*****结束执行Regulator*****\n"
fi


