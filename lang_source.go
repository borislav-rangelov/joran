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
	}
)

func NewLookupSource(parent LangSource, configs ConfigMap) *LookupSource {
	return &LookupSource{
		fallback: parent,
		configs:  configs,
	}
}

func (l *LookupSource) Get(key string, ctx ...any) *Result {
	if len(ctx) > 0 {
		return l.Msg(MsgCtx(key, ctx[0]))
	}
	return l.Msg(Msg(key))
}

func (l *LookupSource) Msg(m *Message) *Result {
	config := l.configs.findConfig(m.Key)

	if config == nil || config.Template == nil {
		if l.fallback != nil {
			return l.fallback.Msg(m)
		}
		return &Result{
			Msg: m,
			Txt: "",
		}
	}

	return &Result{
		Msg: m,
		Txt: config.Template(m),
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
	if len(r.Txt) > 0 {
		return r.Txt
	}
	return other
}

func (r *Result) OrErr() (string, error) {
	if len(r.Txt) > 0 {
		return r.Txt, nil
	}
	return "", errors.Join(ErrTranslationNotFound, errors.New("Missing key: "+r.Key()))
}
