package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var (
	remoteClient         string
	bucket               string
	domain               string
	subdomain            string
	listDomainsFlag      bool
	listSubdomainsFlag   bool
	listPortsFlag        bool
	listHttpxResultsFlag bool
)

func init() {
	flag.StringVar(&remoteClient, "r", "local", "Specify the MinIO remote client")
	flag.StringVar(&bucket, "s3", "airflow-subs", "Specify the MinIO bucket")
	flag.StringVar(&domain, "d", "", "Specify the domain")
	flag.StringVar(&subdomain, "sd", "", "Specify the subdomain")
	flag.BoolVar(&listDomainsFlag, "dl", false, "List all domains")
	flag.BoolVar(&listSubdomainsFlag, "sl", false, "List all subdomains")
	flag.BoolVar(&listPortsFlag, "pl", false, "List all ports")
	flag.BoolVar(&listHttpxResultsFlag, "hl", false, "List httpx results")
}

func executeCommand(cmdStr string) string {
	cmdArgs := strings.Split(cmdStr, " ")
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error executing command: %s\n", err)
		os.Exit(1)
	}
	return string(output)
}

func listDomains() {
	cmdStr := fmt.Sprintf("mc ls %s/%s", remoteClient, bucket)
	output := executeCommand(cmdStr)
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if len(line) > 0 {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				domain := strings.TrimSuffix(parts[len(parts)-1], "/")
				fmt.Println(domain)
			}
		}
	}
}

func listSubdomains() {
	if domain == "" {
		fmt.Println("Domain must be specified with -d")
		os.Exit(1)
	}
	cmdStr := fmt.Sprintf("mc cat %s/%s/%s/subs.txt", remoteClient, bucket, domain)
	output := executeCommand(cmdStr)
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if len(line) > 0 {
			fmt.Println(line)
		}
	}
}

func listPorts() {
	if domain == "" {
		fmt.Println("Domain must be specified with -d")
		os.Exit(1)
	}
	cmdStr := fmt.Sprintf("mc cat %s/%s/%s/ports.txt", remoteClient, bucket, domain)
	output := executeCommand(cmdStr)
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if len(line) > 0 {
			fmt.Println(line)
		}
	}
}

func listHttpxResults() {
	if domain == "" {
		fmt.Println("Domain must be specified with -d")
		os.Exit(1)
	}
	if subdomain != "" {
		cmdStr := fmt.Sprintf("mc cat %s/%s/%s/%s/httpx.json", remoteClient, bucket, domain, subdomain)
		output := executeCommand(cmdStr)
		fmt.Println(output)
		return
	}
	cmdStr := fmt.Sprintf("mc ls %s/%s/%s", remoteClient, bucket, domain)
	output := executeCommand(cmdStr)
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if len(line) > 0 {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				subDir := parts[len(parts)-1]
				if strings.HasSuffix(subDir, "/") {
					subDir = strings.TrimSuffix(subDir, "/")
					cmdStr := fmt.Sprintf("mc cat %s/%s/%s/%s/httpx.json", remoteClient, bucket, domain, subDir)
					output := executeCommand(cmdStr)
					fmt.Println(output)
				}
			}
		}
	}
}

func main() {
	flag.Parse()

	if listDomainsFlag {
		listDomains()
	} else if listSubdomainsFlag {
		listSubdomains()
	} else if listPortsFlag {
		listPorts()
	} else if listHttpxResultsFlag {
		listHttpxResults()
	} else {
		fmt.Println("Usage: mcc [options]")
		fmt.Println("Options:")
		fmt.Println("  -dl            List all domains")
		fmt.Println("  -sl            List all subdomains for a domain")
		fmt.Println("  -pl            List all ports for a domain")
		fmt.Println("  -hl            List httpx results for a domain or subdomain")
		flag.PrintDefaults()
	}
}
