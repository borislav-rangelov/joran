package wut

import (
	"errors"
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"os"
	"path"
)

type (
	Values map[string]any

	LangFile struct {
		Language string
		Fallback string
		Values   Values
	}
)

// ReadFile opens the file and reads the config from it.
// Supported types: toml
func ReadFile(filename string) (*LangFile, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	values, err := readFile(filename, bytes)

	if err != nil {
		return nil, err
	}

	return ValuesToLangFile(values), nil
}

func ValuesToLangFile(values Values) *LangFile {
	lang := getLanguage(&values)
	fallback := getFallback(&values)

	delete(values, "language")
	delete(values, "fallback")

	cleanUp(values)

	return &LangFile{
		Language: lang,
		Fallback: fallback,
		Values:   values,
	}
}

func readFile(filename string, bytes []byte) (Values, error) {
	switch path.Ext(filename) {
	case ".toml":
		return readToml(bytes)
	default:
		return nil, errors.New(fmt.Sprintf("file type not supported: %s", filename))
	}
}

func readToml(bytes []byte) (Values, error) {
	var values Values
	if err := toml.Unmarshal(bytes, &values); err != nil {
		return nil, err
	}
	return values, nil
}

func getLanguage(values *Values) string {
	lan, ok := (*values)["language"]
	if !ok {
		return ""
	}
	s, ok := lan.(string)
	if !ok {
		return ""
	}
	return s
}

func getFallback(values *Values) string {
	val, ok := (*values)["fallback"]
	if !ok {
		return ""
	}
	s, ok := val.(string)
	if !ok {
		return ""
	}
	return s
}

func (l *LangFile) Include(other *LangFile) error {
	if l == nil || other == nil {
		return errors.New("nil LangFile")
	}
	if l.Language != other.Language ||
		l.Fallback != other.Fallback {
		return errors.New("mismatched language and/or fallback")
	}
	return mergeIntoLeft(l.Values, other.Values)
}

func cleanUp(m map[string]any) {
	for k, v := range m {
		if isEmpty(m, k) {
			delete(m, k)
			continue
		}
		if isStruct(v) {
			cleanUp(v.(map[string]any))
		}
	}
}

func isStruct(val any) bool {
	_, ok := val.(map[string]any)
	return ok
}

func isEmpty(m map[string]any, key string) bool {
	if m[key] == nil {
		return true
	}
	if val, ok := m[key].(string); ok {
		return len(val) == 0
	}
	if val, ok := m[key].(map[string]any); ok {
		return len(val) == 0
	}
	return false
}

func mergeIntoLeft(left, right map[string]any) error {
	for k, v := range right {

		if isEmpty(right, k) {
			continue
		}

		if isEmpty(left, k) {
			left[k] = v
			continue
		}

		err := tryMerge(left, right, k)
		if err != nil {
			return err
		}
	}
	return nil
}

func tryMerge(left, right map[string]any, k string) error {
	leftVal := left[k]
	rightVal := right[k]

	if _, ok := leftVal.(string); ok {
		return cannotMergeError(k, leftVal, rightVal)
	}

	if leftMap, ok := leftVal.(map[string]any); ok {
		if rightMap, ok := rightVal.(map[string]any); ok {
			return mergeIntoLeft(leftMap, rightMap)
		}
	}

	return cannotMergeError(k, leftVal, rightVal)
}

func cannotMergeError(k string, leftVal any, rightVal any) error {
	return errors.New(fmt.Sprintf("key %s: cannot merge %v into %v", k, leftVal, rightVal))
}
