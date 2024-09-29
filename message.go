package wut

import "strings"

type Message struct {
	Key     []string
	Context any
}

func Msg(key string) *Message {
	return &Message{
		Key:     strings.Split(key, "."),
		Context: nil,
	}
}

func MsgCtx(key string, ctx any) *Message {
	return &Message{
		Key:     strings.Split(key, "."),
		Context: ctx,
	}
}
