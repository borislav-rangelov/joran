package wut

import "strings"

type Message struct {
	Key        []string
	NoFallback bool
	Context    any
}

func Msg(key string, ctx ...any) *Message {
	path := strings.Split(key, ".")
	if len(ctx) == 0 {
		return &Message{Key: path, Context: nil}
	}
	return &Message{Key: path, Context: ctx[0]}
}
