package schema

import (
	"regexp"
	"strings"
)

var templatePrefixRegExp = regexp.MustCompile(`<%=\s*template\s+`)

func ApplyTemplates(content string) (string, error) {
	templates := findTemplateDeclarations(content)
	if len(templates) == 0 {
		return content, nil
	}

	// fmt.Println(len(content), content[templates[1][0]:templates[1][0]+templates[1][1]])

	return "", nil
}

func findTemplateDeclarations(content string) [][]int {
	indexes := templatePrefixRegExp.FindAllStringIndex(content, -1)
	if len(indexes) == 0 {
		return nil
	}

	size := len(indexes)
	declarations := make([][]int, 0, size)
	const endSize = 2

	for i := 0; i < size; i++ {
		var template string
		begin := indexes[i][0]

		if i+1 < size {
			end := indexes[i+1][0]
			template = content[begin:end]
		} else {
			template = content[begin:]
		}

		end := strings.LastIndex(template, "%>")
		if end >= 0 {
			end += endSize
			declarations = append(declarations, []int{begin, end})
		}
	}

	return declarations
}
