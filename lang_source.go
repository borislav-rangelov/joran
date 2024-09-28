package joran

import (
	"errors"
	"fmt"
)

var ErrTranslationNotFound = errors.New("translation not found")

type (
	LangSource interface {
		Translate(*Message) string
		TryTranslate(*Message) (string, error)
	}

	ConfigMap map[string]*KeyConfig

	LookupSource struct {
		parent  LangSource
		configs ConfigMap
	}

	KeyConfig struct {
		Configs  ConfigMap
		Template Template
	}
)

func NewLookupSource(parent LangSource, configs ConfigMap) *LookupSource {
	return &LookupSource{
		parent:  parent,
		configs: configs,
	}
}

func (l *LookupSource) Translate(m *Message) string {
	val, err := l.TryTranslate(m)
	if err != nil {
		panic(fmt.Sprintf("Failed to find translation for key: %s", m.Key))
	}
	return val
}

func (l *LookupSource) TryTranslate(m *Message) (string, error) {
	config := l.configs.findConfig(m.Key)

	if config == nil || config.Template == nil {
		if l.parent != nil {
			return l.parent.TryTranslate(m)
		}
		return "", ErrTranslationNotFound
	}

	return config.Template(m), nil
}

func (c ConfigMap) findConfig(keys []string) *KeyConfig {
	key, rest := keys[0], keys[1:]
	if config, ok := c[key]; ok {
		if len(rest) == 0 {
			return config
		}
		return config.Configs.findConfig(rest)
	}
	return nil
}
