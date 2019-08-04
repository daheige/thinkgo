package common

import "os/exec"

//运行shell脚本
func RunShell(exeStr string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", exeStr)
	bytes, err := cmd.CombinedOutput()
	return string(bytes), err
}
