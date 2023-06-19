package schema

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	templatePrefixRegExp = regexp.MustCompile(`<%=\s*template\s+`)
	templateRegExp       = regexp.MustCompile(`<%=\s*template\s+((\w+)\s*=\s*"?([^"]*)"?\s*)*%>`)
	templateParamRegExp  = regexp.MustCompile(`(\w+)\s*=\s*"([^"]*)"`)
	dynamicParamRegExp   = regexp.MustCompile(`\A\$\{(\w+)\}\z`)
)

func ApplyTemplates(refPath, content string) (string, error) {
	tmplInd := findTemplateDeclarations(content)
	if len(tmplInd) == 0 {
		return content, nil
	}

	templates, err := renderTemplates(refPath, content, tmplInd)
	if err != nil {
		return "", err
	}

	sb := strings.Builder{}
	sb.WriteString(strings.Trim(content[:tmplInd[0][0]], " "))

	size := len(tmplInd)
	for i := 0; i < size; i++ {
		sb.WriteString(templates[i])
		pair := tmplInd[i]

		if i+1 < size {
			begin := pair[0] + pair[1]
			end := tmplInd[i+1][0]
			sb.WriteString(strings.Trim(content[begin:end], " "))
		}
	}

	index := tmplInd[size-1][0] + tmplInd[size-1][1]
	sb.WriteString(strings.Trim(content[index:], " "))

	return sb.String(), nil
}

func renderTemplates(refPath, content string, indexes [][]int) ([]string, error) {
	templates := make([]string, len(indexes))
	for i, pair := range indexes {
		begin := pair[0]
		end := pair[0] + pair[1]
		tmplDefinition := content[begin:end]

		excerpt, err := renderTemplate(refPath, tmplDefinition)
		if err != nil {
			return nil, err
		}

		templates[i] = excerpt
	}

	return templates, nil
}

func renderTemplate(refPath, tmplDefinition string) (string, error) {
	if !templateRegExp.MatchString(tmplDefinition) {
		return "", fmt.Errorf("invalid template definition `%s'", tmplDefinition)
	}

	params := templateParamRegExp.FindAllStringSubmatch(tmplDefinition, -1)
	if err := validateParams(params); err != nil {
		return "", err
	}

	pathIndex := findPathParam(params)
	if pathIndex < 0 {
		return "", fmt.Errorf("path parameter is required `%s'", tmplDefinition)
	}

	path := resolvePath(refPath, params[pathIndex][2])
	if err := validatePath(path); err != nil {
		return "", err
	}

	tmpl, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	template := string(tmpl)
	for _, param := range params {
		key := param[1]
		value := param[2]

		if dynamicParamRegExp.MatchString(value) {
			value = dynamicParamRegExp.FindStringSubmatch(value)[1]
		} else {
			key = fmt.Sprintf("${%s}", key)
		}

		template = strings.ReplaceAll(template, key, value)
	}

	return template, nil
}

func validateParams(paramGroups [][]string) error {
	for i, params := range paramGroups {
		pname := params[1]
		if params[2] == "" {
			return fmt.Errorf("template parameter '%s' is empty", pname)
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
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("%s not found", path)
	} else if info.IsDir() {
		return fmt.Errorf("%s is a directory", path)
	}

	return nil
}

func resolvePath(refPath, path string) string {
	if filepath.IsAbs(path) {
		return path
	}

	return filepath.Join(filepath.Dir(refPath), path)
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
