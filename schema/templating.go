package schema

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	templatePrefixRegExp = regexp.MustCompile(`<%=\s*template\s+`)
	templateRegExp       = regexp.MustCompile(`<%=\s*template\s+((\w+)\s*=\s*"?([^"]*)"?\s*)*%>`)
	templateParamRegExp  = regexp.MustCompile(`(\w+)\s*=\s*"([^"]*)"`)
)

func ApplyTemplates(content string) (string, error) {
	tmplInd := findTemplateDeclarations(content)
	if len(tmplInd) == 0 {
		return content, nil
	}

	templates, err := renderTemplates(content, tmplInd)
	if err != nil {
		return "", err
	}

	fmt.Println(templates)

	return "", nil
}

func renderTemplates(content string, indexes [][]int) (string, error) {
	for _, pair := range indexes {
		begin := pair[0]
		end := pair[0] + pair[1]
		template := content[begin:end]

		if !templateRegExp.MatchString(template) {
			return "", fmt.Errorf("invalid template definition `%s'", template)
		}

		params := templateParamRegExp.FindAllStringSubmatch(template, -1)
		if err := validateParams(params); err != nil {
			return "", err
		}

		pathIndex := findPathParam(params)
		if pathIndex < 0 {
			return "", fmt.Errorf("path parameters is required `%s'", template)
		}

		path := params[pathIndex][2]
		if err := validatePath(path); err != nil {
			return "", err
		}
	}

	return "", nil
}

func validateParams(paramGroups [][]string) error {
	for i, params := range paramGroups {
		pname := params[1]
		if params[2] == "" {
			return fmt.Errorf("template parameter %s is empty", pname)
		}

		for j, params2 := range paramGroups {
			if i != j && pname == params2[1] {
				return fmt.Errorf("repeated parameter `%s'", pname)
			}
		}
	}

	return nil
}

func findPathParam(paramGroups [][]string) int {
	for i, params := range paramGroups {
		if params[1] == "path" {
			return i
		}
	}

	return -1
}

func validatePath(path string) error {
	if path == "" {
		return fmt.Errorf("path is required")
	}

	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("%s not found", path)
	} else if info.IsDir() {
		return fmt.Errorf("%s is a directory", path)
	}

	return nil
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
