package wut

import (
	"errors"
	"strings"
)

var ErrTranslationNotFound = errors.New("translation not found")

type (
	// LangSource searches for a translation to the given key
	LangSource interface {
		Get(key string, ctx ...any) *Result
		GetFirst(key []string, ctx ...any) *Result
		Msg(*Message) *Result
	}

	ConfigMap map[string]*KeyConfig

	LookupSource struct {
		fallback LangSource
		configs  ConfigMap
	}

	KeyConfig struct {
		Configs  ConfigMap
		Template Template
	}

	Result struct {
		Msg *Message
		Txt string
		Err error
	}
)

func NewLookupSource(parent LangSource, configs ConfigMap) *LookupSource {
	return &LookupSource{
		fallback: parent,
		configs:  configs,
	}
}

func (l *LookupSource) Get(key string, ctx ...any) *Result {
	return l.Msg(Msg(key, ctx...))
}

func (l *LookupSource) GetFirst(keys []string, ctx ...any) *Result {
	if len(keys) == 0 {
		return emptyResult(Msg("<empty-key-provided>", ctx...))
	}
	for _, key := range keys {
		r := l.Get(key, ctx...)
		if r.HasTxt() {
			return r
		}
	}
	return emptyResult(Msg(keys[0], ctx...))
}

func (l *LookupSource) Msg(m *Message) *Result {
	config := l.configs.findConfig(m.Key)

	if config == nil || config.Template == nil {
		if !m.NoFallback && l.fallback != nil {
			return l.fallback.Msg(m)
		}
		return emptyResult(m)
	}

	txt, err := config.Template(m.Context)
	return &Result{
		Msg: m,
		Txt: txt,
		Err: err,
	}
}

func emptyResult(m *Message) *Result {
	return &Result{
		Msg: m,
		Txt: "",
	}
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

func (r *Result) Key() string {
	return strings.Join(r.Msg.Key, ".")
}

func (r *Result) Or(other string) string {
	if r.HasTxt() {
		return r.Txt
	}
	return other
}

func (r *Result) OrErr() (string, error) {
	if r.HasTxt() {
		return r.Txt, nil
	}
	if r.HasError() {
		return "", r.Err
	}
	return "", errors.Join(ErrTranslationNotFound, errors.New("Missing key: "+r.Key()))
}

func (r *Result) HasTxt() bool {
	return len(r.Txt) > 0
}

func (r *Result) HasError() bool {
	return r.Err != nil
}
