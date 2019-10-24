package parser

import "fmt"

var (
	pkgNameTemplate      = "{{ package.Name }}"
	importPkgsTemplate   = "{{ import.pkgs }}"
	funcContentsTemplate = "{{ func.contents }}"
)

var (
	funcContent = `func Test%s(t *testing.T) {
	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mock data here
	%s
}`

	content = fmt.Sprintf(`package %s

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	%s
)

%s
`,
		pkgNameTemplate,
		importPkgsTemplate,
		funcContentsTemplate,
	)
)
