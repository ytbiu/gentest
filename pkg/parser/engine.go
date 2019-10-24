package parser

import (
	"fmt"
	"gentest/pkg/parser/iostream"
	"gentest/pkg/parser/shell"
	"io/ioutil"
	"math"
	"os"
	"strings"
)

const (
	pkgGoMock  = "github.com/golang/mock/gomock"
	pkgTestify = "github.com/stretchr/testify"
	pkgMockgen = "github.com/golang/mock/mockgen"
)

type Engine interface {
	Run()
}

type engine struct {
	fileName        string
	mockSrcFileName string
	mockPath        string
	methodName      string
	doGoGet         bool

	content string
}

func NewEngine(setArgs ...func(e *engine)) Engine {

	e := &engine{}

	for _, setArg := range setArgs {
		setArg(e)
	}

	e.content = content
	e.check()

	return e
}

func (e *engine) check() {

	if e.fileName == "" {
		fmt.Println("fileName is empty")
		os.Exit(1)
	}

	dirPath, _ := os.Getwd()
	dir, err := ioutil.ReadDir(dirPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	matchFileName := func(f os.FileInfo) bool {
		if !f.IsDir() {
			if f.Name() == e.fileName {
				return true

			}
		}
		return false
	}

	for _, f := range dir {
		if matchFileName(f) {
			return
		}
	}
	fmt.Println("file not found in current dir")
	os.Exit(1)
}

func (e *engine) Run() {
	shell.GoFmt()
	defer shell.GoFmt()

	if e.doGoGet {
		shell.GoGet(pkgGoMock, pkgTestify, pkgMockgen)
	}

	e.execMockCmd()
	e.writeTestFile()

	fmt.Println("done")
}

func (e *engine) testFileName() string {

	prefix := strings.Split(e.fileName, ".")[0]
	return fmt.Sprintf("%s_test.go", prefix)
}

func (e *engine) mockSrcName() string {
	if e.mockSrcFileName != "" {
		return e.mockSrcFileName
	}

	return e.fileName
}

func (e *engine) testMethodName() string {

	return "Test" + strings.Title(e.methodName)
}

func (e *engine) testPackageName() string {

	dir, _ := os.Getwd()
	dirGroup := strings.Split(dir, "/")
	return dirGroup[len(dirGroup)-1]
}

func (e *engine) mockFileName() string {

	prefix := strings.Split(e.fileName, ".")[0]
	return e.mockPath + "/" + fmt.Sprintf("%s_mock.go", prefix)
}

func (e *engine) mockPkgName() string {

	if !strings.Contains(e.mockPath, "/") {
		return e.mockPath
	}
	return strings.Split(e.mockPath, "/")[len(e.mockPath)-1]
}

func (e *engine) mockInfo(infoTable map[string][]string) (string, string) {
	return getMockDatas(infoTable, e.mockPkgName()), getImportPath(e.mockPath)
}

func getImportPath(mockPath string) string {
	var importsContent strings.Builder

	goPath := shell.GoPath()
	if strings.HasPrefix(mockPath, goPath) {
		// 绝对路径
		importPath := strings.TrimPrefix(mockPath, goPath+"/src/")
		importsContent.WriteString(importPath)
	} else {
		// 相对路径
		pwdDir, _ := os.Getwd()
		if strings.HasPrefix(pwdDir, goPath) {
			importPath := fmt.Sprintf(`"%s/%s"`,
				strings.TrimPrefix(pwdDir, goPath+"/src/"),
				mockPath,
			)
			importsContent.WriteString(importPath)
		} else {
			fmt.Println("warning : project should be set in GOPATH")
		}
	}
	return importsContent.String()

}

func getMockDatas(infoTable map[string][]string, mockPkgName string) string {
	var newMockDatas strings.Builder
	for name, methods := range infoTable {

		mockName := strings.ToLower(name)
		mockNameSuffix := strings.Title(strings.TrimPrefix(mockName, "mock"))
		newMockDatas.WriteString(fmt.Sprintf("mock%s := %s.New%s(ctrl) \n", mockNameSuffix, mockPkgName, name))
		for _, method := range methods {
			newMockDatas.WriteString(fmt.Sprintf("mock%s.EXPECT().%s.Return() \n", mockNameSuffix, method))
			newMockDatas.WriteString("\n")
			newMockDatas.WriteString(`// call func here after mock`)
			newMockDatas.WriteString("\n \n")
			newMockDatas.WriteString("a.Equal(nil, nil)")
			newMockDatas.WriteString("\n")

		}
		newMockDatas.WriteString("\n")
	}
	return newMockDatas.String()
}

func (e *engine) writeTestFile() {

	testFileName := e.testFileName()
	testFile := iostream.File(testFileName)
	if testFile == nil {
		fmt.Printf("the file : %s is exist \n", testFileName)
		os.Exit(1)
	}
	defer testFile.Close()

	infoTable := e.interfaceFromMock()
	mockDatas, importsPkgs := e.mockInfo(infoTable)

	methodNames := e.getMethondNames()

	var funcContents strings.Builder
	for _, methodName := range methodNames {
		funcContent := fmt.Sprintf(funcContent, strings.Title(methodName), mockDatas)
		funcContents.WriteString(funcContent)
		funcContents.WriteString("\n")
	}

	replacer := strings.NewReplacer(
		pkgNameTemplate, e.testPackageName(),
		importPkgsTemplate, importsPkgs,
		funcContentsTemplate, funcContents.String(),
	)
	content = replacer.Replace(content)

	testFile.Truncate(0)
	if _, err := testFile.WriteString(content); err != nil {
		fmt.Println(err)
		return
	}
}

func (e *engine) execMockCmd() {
	shell.MockGen(e.mockFileName(), e.mockSrcName(), e.mockPkgName())
}

func (e *engine) interfaceFromMock() map[string][]string {
	infoTable := make(map[string][]string)

	iostream.ReadByLineWithDo(e.mockFileName(), func(line string) {
		if strings.Contains(line, "type") &&
			strings.Contains(line, "struct ") &&
			!strings.Contains(line, "Recorder") &&
			!strings.Contains(line, "//") {

			tmp := strings.Split(strings.TrimSpace(strings.Split(line, "struct")[0]), " ")
			infName := tmp[len(tmp)-1]

			infoTable[infName] = []string{}
			return
		}

		for name := range infoTable {
			methodFlag := fmt.Sprintf("*%s) ", name)
			if strings.Contains(line, methodFlag) {
				tmp := strings.Split(line, methodFlag)
				methodSign := strings.Split(tmp[len(tmp)-1], ")")[0] + ")"

				if methodSign == "EXPECT()" {
					continue
				}

				indexLeft := strings.Index(methodSign, "(")
				indexRight := strings.Index(methodSign, ")")
				params := strings.TrimSpace(methodSign[indexLeft+1 : indexRight])

				if len(params) > 0 {
					paramLen := len(strings.Split(params, " "))
					if paramNum := math.Ceil(float64(paramLen) / 2.0); paramNum > 0 {
						var mockParams strings.Builder
						for i := 0; i < int(paramNum); i++ {
							mockParams.WriteString("nil,")
						}
						methodSign = methodSign[:indexLeft+1] + mockParams.String() + methodSign[indexRight:]
					}
				}

				infoTable[name] = append(infoTable[name], methodSign)

			}
		}
	})

	return infoTable
}

func (e *engine) getMethondNames() []string {
	var methodNames []string
	if e.methodName != "" {
		return append(methodNames, e.methodName)
	}

	iostream.ReadByLineWithDo(e.fileName, func(line string) {
		if strings.Contains(line, "func ") &&
			strings.Contains(line, "(") &&
			strings.Contains(line, ")") &&
			strings.Contains(line, "{") &&
			!strings.Contains(line, "//") {

			splited := strings.Split(line, " ")
			var methodName string
			if strings.Contains(line, "func (") {
				methodName = splited[3]
			} else {
				methodName = splited[1]
			}
			i := strings.Index(methodName, "(")
			methodNames = append(methodNames, methodName[:i])
		}
	})

	return methodNames
}
