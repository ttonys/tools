package main

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/yosssi/gohtml"
	"io/ioutil"
	"os"
	"path/filepath"
)

type HttpxOutput struct {
	URL         string `json:"url"`
	StatusCode  int    `json:"status_code"`
	Title       string `json:"title"`
	WebServer   string `json:"webserver"`
	ContentType string `json:"content_type"`
	Input       string `json:"input"`
	Host        string `json:"host"`
	Path        string `json:"path"`
	RawHeader   string `json:"raw_header"`
	Request     string `json:"request"`
	Body        string `json:"body"`
}

func main() {
	outDir := "out"
	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		os.Mkdir(outDir, 0755)
	}

	indexFile, err := os.Create(filepath.Join(outDir, "index"))
	if err != nil {
		fmt.Println("Error creating index file:", err)
		return
	}
	defer indexFile.Close()

	index2File, err := os.Create(filepath.Join(outDir, "index2"))
	if err != nil {
		fmt.Println("Error creating index2 file:", err)
		return
	}
	defer index2File.Close()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var output HttpxOutput
		rawJSON := scanner.Bytes()
		err := json.Unmarshal(rawJSON, &output)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			continue
		}

		hostDir := filepath.Join(outDir, output.Input)
		if _, err := os.Stat(hostDir); os.IsNotExist(err) {
			os.Mkdir(hostDir, 0755)
		}

		hash := fmt.Sprintf("%x", sha1.Sum([]byte(output.URL)))
		urlFilePath := filepath.Join(hostDir, hash)
		jsonFilePath := urlFilePath + ".json"

		// 格式化 JSON 数据
		var formattedRawJSON bytes.Buffer
		err = json.Indent(&formattedRawJSON, rawJSON, "", "  ")
		if err != nil {
			fmt.Println("Error formatting JSON:", err)
			continue
		}

		// 保存格式化后的 JSON 内容
		err = ioutil.WriteFile(jsonFilePath, formattedRawJSON.Bytes(), 0644)
		if err != nil {
			fmt.Println("Error writing JSON file:", err)
			continue
		}

		// 对 HTML body 进行格式化
		formattedBody := gohtml.Format(output.Body)

		detailedResponse := fmt.Sprintf("%s\n\n%s\n\n%s\n\n%s\n",
			output.URL,
			output.Request,
			output.RawHeader,
			formattedBody,
		)
		err = ioutil.WriteFile(urlFilePath, []byte(detailedResponse), 0644)
		if err != nil {
			fmt.Println("Error writing response file:", err)
			continue
		}

		// 去掉当前目录部分的路径
		relativeURLFilePath := filepath.Join(output.Input, hash)
		relativeJSONFilePath := relativeURLFilePath + ".json"

		globalIndexEntry := fmt.Sprintf("%s %s (%d %s)\n", relativeURLFilePath, output.URL, output.StatusCode, output.Title)
		indexFile.WriteString(globalIndexEntry)

		globalIndex2Entry := fmt.Sprintf("%s %s (%d %s)\n", relativeJSONFilePath, output.URL, output.StatusCode, output.Title)
		index2File.WriteString(globalIndex2Entry)

		localIndexFilePath := filepath.Join(hostDir, "index")
		localIndex2FilePath := filepath.Join(hostDir, "index2")

		localIndexEntry := fmt.Sprintf("./%s %s (%d %s)\n", hash, output.URL, output.StatusCode, output.Title)
		localIndex2Entry := fmt.Sprintf("./%s.json %s (%d %s)\n", hash, output.URL, output.StatusCode, output.Title)

		localIndexFile, err := os.OpenFile(localIndexFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Error creating/opening local index file:", err)
			continue
		}
		defer localIndexFile.Close()
		localIndexFile.WriteString(localIndexEntry)

		localIndex2File, err := os.OpenFile(localIndex2FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Error creating/opening local index2 file:", err)
			continue
		}
		defer localIndex2File.Close()
		localIndex2File.WriteString(localIndex2Entry)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}
}
