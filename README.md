

## tools

> 子域名枚举脚本

#### 安装
- python（版本python3）
- pip install -r requirements.txt

#### 使用
- 仅进行子域名扫描: `cd subs && ./Subs -d example.com`
- 进行子域名扫描并fuzz: `cd subs && ./Subs -d example.com -x`
- 仅进行fuzz: `cd subs && ./Subs -d example.com -f subs.txt -x`