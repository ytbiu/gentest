package parser

func FileName(fileName string) func(e *engine) {
	return func(e *engine) {
		e.fileName = fileName
	}
}

func MockPath(mockPath string) func(e *engine) {
	return func(e *engine) {
		e.mockPath = mockPath
	}
}

func MockSrcFileName(srcFileName string) func(e *engine) {
	return func(e *engine) {
		e.mockSrcFileName = srcFileName
	}
}

func MethodName(methodName string) func(e *engine) {
	return func(e *engine) {
		e.methodName = methodName
	}
}

func DoGoGet(do bool) func(e *engine) {
	return func(e *engine) {
		e.doGoGet = do
	}
}
