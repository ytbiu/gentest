package iostream

import (
	"bufio"
	"errors"
	"io"
	"os"
)

func File(path string) *os.File {
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		return nil
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		panic(err)
	}
	return f
}

func ReadByLineWithDo(fileName string, do func(line string)) {
	f, _ := os.OpenFile(fileName, os.O_RDONLY, 0755)
	buf := bufio.NewReader(f)
	for {
		line, _, err := buf.ReadLine()
		if err != nil && err != io.EOF {
			panic(err)
		}
		if errors.Is(err, io.EOF) {
			return
		}
		do(string(line))
	}
}
