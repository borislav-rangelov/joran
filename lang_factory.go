package wut

import (
	"fmt"
	"slices"
	"strings"
)

type (
	LangFactory interface {
		Lang(code string) LangSource
		LangNoFallback(code string) LangSource
	}

	langFactory struct {
		codeToConfig map[string]ConfigMap
		fallbackMap  map[string]string
	}
)

func NewLangFactory(tf TemplateFactory, files ...*LangFile) (LangFactory, error) {
	fileMap, fallbackMap, err := mergeLangFiles(files)

	if err != nil {
		return nil, err
	}

	configMaps := make(map[string]ConfigMap)
	for k, file := range fileMap {
		m, er := toConfigMap(file.Language, []string{k}, file.Values, tf)
		if er != nil {
			return nil, er
		}
		configMaps[k] = m
	}
	return &langFactory{
		codeToConfig: configMaps,
		fallbackMap:  fallbackMap,
	}, nil
}

func (l *langFactory) Lang(code string) LangSource {
	code = strings.ToLower(code)
	if configMap, ok := l.codeToConfig[code]; ok {
		var parent LangSource
		fallbackCode := l.fallbackMap[code]
		if fallbackCode != "" {
			parent = l.Lang(fallbackCode)
		}
		return NewLookupSource(parent, configMap)
	}
	return NewLookupSource(nil, make(ConfigMap))
}

func (l *langFactory) LangNoFallback(code string) LangSource {
	code = strings.ToLower(code)
	if configMap, ok := l.codeToConfig[code]; ok {
		return NewLookupSource(nil, configMap)
	}
	return NewLookupSource(nil, make(ConfigMap))
}

func mergeLangFiles(files []*LangFile) (map[string]*LangFile, map[string]string, error) {
	fileMap := make(map[string]*LangFile)
	fallbackMap := make(map[string]string)

	for _, f := range files {
		if fileMap[f.Language] == nil {
			fileMap[f.Language] = f
			fallbackMap[f.Language] = f.Fallback
		} else {
			err := fileMap[f.Language].Include(f)
			if err != nil {
				return nil, nil, err
			}
		}
	}
	return fileMap, fallbackMap, nil
}

func toConfigMap(lang string, path []string, values map[string]any, tf TemplateFactory) (ConfigMap, error) {
	result := make(ConfigMap)

	for k, v := range values {
		keyPath := slices.Clone(path)
		keyPath = append(keyPath, k)

		if entry, ok := v.(string); ok {
			template, err := tf.GetTemplate(lang, keyPath, entry)
			if err != nil {
				return nil, err
			}
			result[k] = &KeyConfig{
				Template: template,
			}
			continue
		}

		if cfg, ok := v.(map[string]any); ok {
			configMap, err := toConfigMap(lang, keyPath, cfg, tf)
			if err != nil {
				return nil, err
			}
			result[k] = &KeyConfig{
				Configs: configMap,
			}
			continue
		}

		return nil, fmt.Errorf("unknown value type: %T", v)
	}

	return result, nil
}
