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
	outputJsonFlag       bool
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
	flag.BoolVar(&outputJsonFlag, "json", false, "Output in JSON format")
}

func executeCommand(cmdStr string) {
	var cmd *exec.Cmd
	if strings.Contains(cmdStr, "|") {
		cmdParts := strings.Split(cmdStr, "|")
		cmd1Args := strings.Split(strings.TrimSpace(cmdParts[0]), " ")
		cmd2Args := strings.Split(strings.TrimSpace(cmdParts[1]), " ")
		cmd1 := exec.Command(cmd1Args[0], cmd1Args[1:]...)
		cmd2 := exec.Command(cmd2Args[0], cmd2Args[1:]...)
		pipe, err := cmd1.StdoutPipe()
		if err != nil {
			fmt.Printf("Error creating pipe: %s\n", err)
			os.Exit(1)
		}
		cmd2.Stdin = pipe
		cmd2.Stdout = os.Stdout
		cmd2.Stderr = os.Stderr
		err = cmd1.Start()
		if err != nil {
			fmt.Printf("Error starting command 1: %s\n", err)
			os.Exit(1)
		}
		err = cmd2.Start()
		if err != nil {
			fmt.Printf("Error starting command 2: %s\n", err)
			os.Exit(1)
		}
		err = cmd1.Wait()
		if err != nil {
			fmt.Printf("Error waiting for command 1: %s\n", err)
			os.Exit(1)
		}
		err = cmd2.Wait()
		if err != nil {
			fmt.Printf("Error waiting for command 2: %s\n", err)
			os.Exit(1)
		}
	} else {
		cmdArgs := strings.Split(cmdStr, " ")
		cmd = exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Error executing command: %s\n", err)
			fmt.Printf("Command: %s\n", cmdStr)
			os.Exit(1)
		}
	}
}

func listDomains() {
	cmdStr := fmt.Sprintf("mc ls %s/%s", remoteClient, bucket)
	executeCommand(cmdStr)
}

func listSubdomains() {
	if domain == "" {
		fmt.Println("Domain must be specified with -d")
		os.Exit(1)
	}
	cmdStr := fmt.Sprintf("mc cat %s/%s/%s/subs.txt", remoteClient, bucket, domain)
	executeCommand(cmdStr)
}

func listPorts() {
	if domain == "" {
		fmt.Println("Domain must be specified with -d")
		os.Exit(1)
	}
	cmdStr := fmt.Sprintf("mc cat %s/%s/%s/ports.txt", remoteClient, bucket, domain)
	executeCommand(cmdStr)
}

func listHttpxResults() {
	if domain == "" {
		fmt.Println("Domain must be specified with -d")
		os.Exit(1)
	}
	cmdFormat := "mc cat %s/%s/%s/%s/httpx.json"
	if !outputJsonFlag {
		cmdFormat += " | fastgron"
	}
	if subdomain != "" {
		cmdStr := fmt.Sprintf(cmdFormat, remoteClient, bucket, domain, subdomain)
		executeCommand(cmdStr)
		return
	}
	cmdStr := fmt.Sprintf("mc ls %s/%s/%s", remoteClient, bucket, domain)
	output := executeCommandWithOutput(cmdStr)
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if len(line) > 0 {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				subDir := parts[len(parts)-1]
				if strings.HasSuffix(subDir, "/") {
					subDir = strings.TrimSuffix(subDir, "/")
					cmdStr := fmt.Sprintf(cmdFormat, remoteClient, bucket, domain, subDir)
					executeCommand(cmdStr)
				}
			}
		}
	}
}

func executeCommandWithOutput(cmdStr string) string {
	cmdArgs := strings.Split(cmdStr, " ")
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error executing command: %s\n", err)
		fmt.Printf("Command: %s\n", cmdStr)
		fmt.Printf("Output: %s\n", string(output))
		os.Exit(1)
	}
	return string(output)
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
		fmt.Println("  -json          Output in JSON format")
		flag.PrintDefaults()
	}
}
