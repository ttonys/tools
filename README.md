## 子域名枚举脚本(subs)

> `subs.sh` 是一个用于子域名扫描的脚本

#### 使用
```
./subs.sh -d example.com -w ~/subdomains.txt -o /tmp/res
```



## httpx处理工具(ttt)

`ttt` 是一个用于处理 `httpx` 输出的 Go 工具。该工具会将扫描结果归类并生成索引文件。

#### 安装

要安装该工具，请运行以下命令：

```sh
go install github.com/ttonys/tools/ttt@latest
```

#### 使用方法

假设你有一个包含 URL 列表的文件 urls.txt，可以通过以下命令运行 httpx 并使用 ttt 处理输出：

```sh
cat urls.txt | ~/go/1.22.1/bin/httpx -j -irr -include-chain | ttt
```
该命令会读取 urls.txt 中的 URL 列表，使用 httpx 进行扫描，并将结果通过管道传递给 ttt 工具进行处理和归类。

#### 功能

- 处理每行的 JSON 输出并进行归类和索引。
- 输出结果到 `out` 文件夹，没有的话会新建。
- 每个 URL 会生成一个单独的文件夹，包含详细的响应内容和格式化的 JSON 文件。
- 生成全局和局部索引文件。

#### 示例

运行上述命令后，你会在 `out` 目录中看到类似以下结构：

```
out/
├── index
├── index2
├── example.com/
│   ├── index
│   ├── index2
│   ├── a813f9c4a94fb48076757e8526ddf28409ab2059
│   └── a813f9c4a94fb48076757e8526ddf28409ab2059.json
└── example.net/
    ├── index
    ├── index2
    ├── 9173b353c3423508e738212159c14c0c17a2d2f3
    └── 9173b353c3423508e738212159c14c0c17a2d2f3.json
```

每个 URL 文件夹中的 `index` 文件示例如下：

```
./a813f9c4a94fb48076757e8526ddf28409ab2059 https://www.example.com (404 Not Found)
```

`index2` 文件示例如下：

```
./a813f9c4a94fb48076757e8526ddf28409ab2059.json https://www.example.com (404 Not Found)
```



## crt.sh获取子域名(crt)

> 从crt.sh获取子域名

#### 使用

```
go install github.com/ttonys/tools/crt@latest
```

## glark飞书通知

> 发送通知到飞书机器人

#### 使用

```
go install github.com/ttonys/tools/glark@latest
```

