package shell

import (
	"fmt"
	"os"
	"os/exec"
)

func MockGen(mockFileName, mockSrcFileName, mockPkgName string) {

	mockCmd := fmt.Sprintf(
		"mockgen -destination %s -source %s -package %s",
		mockFileName,
		mockSrcFileName,
		mockPkgName,
	)

	fmt.Println(
		fmt.Sprintf(
			`generating mock file:
		mockFileName : %s
		mockSrcFileName : %s
		mockDir : %s`,
			mockFileName,
			mockSrcFileName,
			mockPkgName,
		),
	)

	cmd := exec.Command("/bin/bash", "-c", mockCmd)
	stdout, _ := cmd.StdoutPipe()
	go echo(stdout)
	stderr, _ := cmd.StderrPipe()
	go echo(stderr)

	if err := cmd.Run(); err != nil {
		fmt.Println(mockCmd)
		fmt.Println(err)
		os.Exit(1)
	}

}
