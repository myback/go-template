package tmpl_functions

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	"github.com/ghodss/yaml"
	"github.com/myback/sprig"
)

type noValueError struct {
	msg string
}

func (e *noValueError) Error() string {
	return e.msg
}

func createFuncMap() template.FuncMap {
	aliased := template.FuncMap{
		"isFile":         isFile,
		"readFile":       readFile,
		"readDir":        readDir,
		"toYaml":         toYaml,
		"fromYaml":       fromYaml,
		"setValueAtPath": setValueAtPath,
		"requiredEnv":    requiredEnv,
		"get":            get,
		"getOrNil":       getOrNil,
		"tpl":            tpl,
		"required":       required,
	}

	funcMap := sprig.TxtFuncMap()

	for name, f := range aliased {
		funcMap[name] = f
	}

	return funcMap
}

func isFile(filename string) (bool, error) {
	path, _ := filepath.Abs(filename)
	stat, err := os.Stat(path)
	if err == nil {
		return !stat.IsDir(), nil
	}

	return false, err
}

func readFile(filename string) (string, error) {
	path, _ := filepath.Abs(filename)

	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func readDir(path string) ([]string, error) {
	contextPath, _ := filepath.Abs(path)

	entries, err := os.ReadDir(contextPath)
	if err != nil {
		return nil, fmt.Errorf("ReadDir %q: %w", contextPath, err)
	}

	var filenames []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		filenames = append(filenames, filepath.Join(path, entry.Name()))
	}

	return filenames, nil
}

func tpl(text string, data interface{}) (string, error) {
	buf, err := renderTemplateToBuffer(text, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func newTemplate() *template.Template {
	return template.
		New("stringTemplate").
		Funcs(createFuncMap()).
		Option("missingkey=error")
}

func renderTemplateToBuffer(s string, data ...interface{}) (*bytes.Buffer, error) {
	var t, parseErr = newTemplate().Parse(s)
	if parseErr != nil {
		return nil, parseErr
	}

	var tplString bytes.Buffer
	var d interface{}
	if len(data) > 0 {
		d = data[0]
	}
	var execErr = t.Execute(&tplString, d)

	if execErr != nil {
		return &tplString, execErr
	}

	return &tplString, nil
}

func toYaml(v interface{}) (string, error) {
	data, err := yaml.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func fromYaml(str string) (map[string]any, error) {
	m := map[string]any{}

	if err := yaml.Unmarshal([]byte(str), &m); err != nil {
		return nil, fmt.Errorf("%s, offending yaml: %s", err, str)
	}
	return m, nil
}

func setValueAtPath(path string, value interface{}, values map[string]any) (map[string]any, error) {
	var current interface{}
	current = values
	components := strings.Split(path, ".")
	pathToMap := components[:len(components)-1]
	key := components[len(components)-1]
	for _, k := range pathToMap {
		var elem interface{}

		switch typedCurrent := current.(type) {
		case map[string]interface{}:
			v, exists := typedCurrent[k]
			if !exists {
				return nil, fmt.Errorf("failed to set value at path \"%s\": value for key \"%s\" does not exist", path, k)
			}
			elem = v
		case map[interface{}]interface{}:
			v, exists := typedCurrent[k]
			if !exists {
				return nil, fmt.Errorf("failed to set value at path \"%s\": value for key \"%s\" does not exist", path, k)
			}
			elem = v
		default:
			return nil, fmt.Errorf("failed to set value at path \"%s\": value for key \"%s\" was not a map", path, k)
		}

		switch typedElem := elem.(type) {
		case map[string]interface{}, map[interface{}]interface{}:
			current = typedElem
		default:
			return nil, fmt.Errorf("failed to set value at path \"%s\": value for key \"%s\" was not a map", path, k)
		}
	}

	switch typedCurrent := current.(type) {
	case map[string]interface{}:
		typedCurrent[key] = value
	case map[interface{}]interface{}:
		typedCurrent[key] = value
	default:
		return nil, fmt.Errorf("failed to set value at path \"%s\": value for key \"%s\" was not a map", path, key)
	}
	return values, nil
}

func requiredEnv(name string) (string, error) {
	if val, exists := os.LookupEnv(name); exists && len(val) > 0 {
		return val, nil
	}

	return "", fmt.Errorf("required env var `%s` is not set", name)
}

func required(warn string, val interface{}) (interface{}, error) {
	if val == nil {
		return nil, fmt.Errorf(warn)
	} else if _, ok := val.(string); ok {
		if val == "" {
			return nil, fmt.Errorf(warn)
		}
	}

	return val, nil
}

func get(path string, varArgs ...interface{}) (interface{}, error) {
	var defSet bool
	var def interface{}
	var obj interface{}
	switch len(varArgs) {
	case 1:
		defSet = false
		def = nil
		obj = varArgs[0]
	case 2:
		defSet = true
		def = varArgs[0]
		obj = varArgs[1]
	default:
		return nil, fmt.Errorf("unexpected number of args pased to the template function get(path, [def, ]obj): expected 1 or 2, got %d, args was %v", len(varArgs), varArgs)
	}

	if path == "" {
		return obj, nil
	}
	keys := strings.Split(path, ".")
	var v interface{}
	var ok bool
	switch typedObj := obj.(type) {
	case *map[string]interface{}:
		obj = *typedObj
	}
	switch typedObj := obj.(type) {
	case map[string]interface{}:
		v, ok = typedObj[keys[0]]
		if !ok {
			if defSet {
				return def, nil
			}
			return nil, &noValueError{fmt.Sprintf("no value exist for key \"%s\" in %v", keys[0], typedObj)}
		}
	case map[interface{}]interface{}:
		v, ok = typedObj[keys[0]]
		if !ok {
			if defSet {
				return def, nil
			}
			return nil, &noValueError{fmt.Sprintf("no value exist for key \"%s\" in %v", keys[0], typedObj)}
		}
	default:
		maybeStruct := reflect.ValueOf(typedObj)
		if maybeStruct.Kind() != reflect.Struct {
			return nil, &noValueError{fmt.Sprintf("unexpected type(%v) of value for key \"%s\": it must be either map[string]interface{} or any struct", reflect.TypeOf(obj), keys[0])}
		} else if maybeStruct.NumField() < 1 {
			return nil, &noValueError{fmt.Sprintf("no accessible struct fields for key \"%s\"", keys[0])}
		}
		f := maybeStruct.FieldByName(keys[0])
		if !f.IsValid() {
			if defSet {
				return def, nil
			}
			return nil, &noValueError{fmt.Sprintf("no field named \"%s\" exist in %v", keys[0], typedObj)}
		}
		v = f.Interface()
	}

	if defSet {
		return get(strings.Join(keys[1:], "."), def, v)
	}
	return get(strings.Join(keys[1:], "."), v)
}

func getOrNil(path string, o interface{}) (interface{}, error) {
	v, err := get(path, o)
	if err != nil {
		switch err.(type) {
		case *noValueError:
			return nil, nil
		default:
			return nil, err
		}
	}
	return v, nil
}
