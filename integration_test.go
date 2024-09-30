package wut

import (
	"testing"
)

type user struct {
	Name string
}

func TestSetup(t *testing.T) {
	err := Setup().
		AddFiles("testdata/i18n/en.toml",
			"testdata/i18n/es.toml").
		AsDefault()

	if err != nil {
		t.Fatal(err)
	}

	english := Lang("en")
	spanish := Lang("es")

	// "Hello"
	if english.Get("display.hello").Txt != "Hello" {
		t.Errorf("english.Get(\"display.hello\").Txt != \"Hello\"")
	}
	// "Hola
	if spanish.Get("display.hello").Txt != "Hola" {
		t.Errorf("spanish.Get(\"display.hello\").Txt != \"Hola\"")
	}

	msg := english.Get("display.bye")
	// ""
	if msg.Txt != "" {
		t.Errorf("english.Get(\"display.bye\").Txt != \"\"")
	}
	// "Bye"
	if msg.Or("Bye") != "Bye" {
		t.Errorf("english.Get(\"display.bye\").Txt != \"Bye\"")
	}

	u := user{Name: "Wut"}
	templatedMsg := english.Get("display.hello_name", u)

	// "Hello, Wut"
	if templatedMsg.Txt != "Hello, Wut" {
		t.Errorf("templatedMsg.Txt != \"Hello, Wut\" but was %s", templatedMsg.Txt)
	}

	templatedGetMsg := spanish.Get("display.hello_template", u)
	if templatedGetMsg.Txt != "Hola, Wut" {
		t.Errorf("templatedGetMsg.Txt != \"Hola, Wut\" but was %s", templatedGetMsg.Txt)
	}
}
