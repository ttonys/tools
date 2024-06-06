package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

// 定义一个结构体来解析crt.sh的响应
type crtshEntry struct {
	CommonName string `json:"common_name"`
	NameValue  string `json:"name_value"`
}

// 去重函数
func unique(slice []string) []string {
	uniqueMap := make(map[string]bool)
	var uniqueList []string
	for _, item := range slice {
		if _, exists := uniqueMap[item]; !exists {
			uniqueMap[item] = true
			uniqueList = append(uniqueList, item)
		}
	}
	return uniqueList
}

// 过滤子域名函数
func filterSubdomains(subdomains []string, domain string) []string {
	var filtered []string
	re := regexp.MustCompile(fmt.Sprintf(`^[a-zA-Z0-9._-]+\.%s$`, regexp.QuoteMeta(domain)))
	for _, subdomain := range subdomains {
		if re.MatchString(subdomain) {
			filtered = append(filtered, subdomain)
		}
	}
	return filtered
}

func main() {
	// 定义命令行参数
	domain := flag.String("d", "", "The domain to search for")
	flag.Parse()

	// 检查是否提供了域名参数
	if *domain == "" {
		fmt.Println("Please specify a domain using -d")
		return
	}

	// 构建查询URL
	url := fmt.Sprintf("https://crt.sh/?q=%%25.%s&output=json", *domain)

	// 发送HTTP请求
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// 解析JSON数据
	var entries []crtshEntry
	err = json.Unmarshal(body, &entries)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// 提取子域名并去重
	var subdomains []string
	for _, entry := range entries {
		subdomains = append(subdomains, strings.Split(entry.NameValue, "\n")...)
	}
	subdomains = unique(subdomains)

	// 过滤子域名
	filteredSubdomains := filterSubdomains(subdomains, *domain)

	// 打印子域名
	for _, subdomain := range filteredSubdomains {
		fmt.Println(subdomain)
	}
}
