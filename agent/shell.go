package agent

import (
	"fmt"
	"io"
	"os/exec"
	"runtime"

	"Stowaway/utils"
)

// CreatInteractiveShell 创建交互式shell
func CreatInteractiveShell() (io.Reader, io.Writer) {
	var cmd *exec.Cmd
	sys := utils.CheckSystem()
	//判断操作系统后决定启动哪一种shell
	switch sys {
	case 0x01:
		cmd = exec.Command("c:\\windows\\system32\\cmd.exe")
	default:
		cmd = exec.Command("/bin/sh", "-i")
		if runtime.GOARCH == "386" || runtime.GOARCH == "amd64" {
			cmd = exec.Command("/bin/bash", "-i")
		}
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err)
	}

	cmd.Stderr = cmd.Stdout //将stderr重定向至stdout
	cmd.Start()

	return stdout, stdin
}

// StartShell 启动shell
func StartShell(command string, stdin io.Writer, stdout io.Reader, currentid string) {
	buf := make([]byte, 1024)
	stdin.Write([]byte(command))

	for {
		count, err := stdout.Read(buf)

		if err != nil {
			return
		}

		respShell, _ := utils.ConstructPayload(utils.AdminId, "", "DATA", "SHELLRESP", " ", string(buf[:count]), 0, currentid, AgentStatus.AESKey, false)
		AgentStuff.ProxyChan.ProxyChanToUpperNode <- respShell
	}
}
