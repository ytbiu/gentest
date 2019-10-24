package shell

import (
	"fmt"
	"os/exec"
	"strings"
)

func GoGet(goPkgGroup ...string) {

	// fmt.Println(fmt.Sprintf(`getting go package ...
	// 	%s,
	// 	%s,
	// 	%s,`,
	// 	pkgGoMock,
	// 	pkgTestify,
	// 	pkgMockgen),
	// )

	for _, goPkg := range goPkgGroup {
		cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("go get %s", goPkg))
		stdout, _ := cmd.StdoutPipe()
		go echo(stdout, fmt.Sprintf("exec go get %s", goPkg))
		stderr, _ := cmd.StderrPipe()
		go echo(stderr, fmt.Sprintf("exec go get %s err :", goPkg))

		if err := cmd.Run(); err != nil {
			fmt.Println(err)
		}
	}
}

func GoPath() string {
	stdout, _ := exec.Command("/bin/bash", "-c", "echo $GOPATH").Output()
	return strings.TrimSpace(string(stdout))
}

func GoFmt() {
	exec.Command("/bin/bash", "-c", "go fmt ./...").Run()
}
