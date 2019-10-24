package shell

import (
	"fmt"
	"io"
	"io/ioutil"
)

func echo(stdout io.ReadCloser, prefix ...string) {

	content, _ := ioutil.ReadAll(stdout)
	if len(prefix) > 0 {
		fmt.Printf("%s : %s \n", prefix[0], string(content))
		return
	}

	fmt.Println(string(content))
}
