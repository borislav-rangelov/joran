package joran

import (
	"errors"
	"testing"
)

type testCase struct {
	lookupSource   *LookupSource
	message        *Message
	translation    string
	expectNotFound bool
}

func TestTranslations(t *testing.T) {
	tests := testCases()

	for _, test := range tests {
		tryTranslateValidation(t, test)
		translateValidation(t, test)
	}
}

func tryTranslateValidation(t *testing.T, test testCase) {
	translate, err := test.lookupSource.TryTranslate(test.message)

	if test.expectNotFound && !errors.Is(err, ErrTranslationNotFound) {
		t.Errorf("TryTranslate(%#v) should have emitted not found error: %v", test.message, err)
	} else if !test.expectNotFound && translate != test.translation {
		t.Errorf("TryTranslate(%#v) should have returned '%s': %s", test.message, test.translation, translate)
	}
}

func translateValidation(t *testing.T, test testCase) {
	defer func() {
		if r := recover(); r != nil && !test.expectNotFound {
			t.Errorf("Translate(%#v) should have returned '%s' instead of panicing: %v", test.message, test.translation, r)
		}
	}()
	translate := test.lookupSource.Translate(test.message)

	if test.translation != translate {
		t.Errorf("Translate(%#v) should have returned '%s': %s", test.message, test.translation, translate)
	}
}

func testCases() []testCase {
	return []testCase{
		{ // empty source
			lookupSource:   NewLookupSource(nil, make(ConfigMap)),
			message:        Msg("a"),
			expectNotFound: true,
		},
		{ // missing key
			lookupSource: NewLookupSource(nil, ConfigMap{
				"a": &KeyConfig{Configs: ConfigMap{
					"b": &KeyConfig{Template: returnT("translation-a-b")},
				}},
			}),
			message:        Msg("a.c"),
			expectNotFound: true,
		},
		{ // single key
			lookupSource: NewLookupSource(nil, ConfigMap{
				"a": &KeyConfig{Template: returnT("b")},
			}),
			message:     Msg("a"),
			translation: "b",
		},
		{ // multiple configs
			lookupSource: NewLookupSource(nil, ConfigMap{
				"a": &KeyConfig{
					Template: returnT("translation-a"),
				},
				"b": &KeyConfig{
					Template: returnT("translation-b"),
				},
			}),
			message:     Msg("a"),
			translation: "translation-a",
		},
		{ // multiple configs
			lookupSource: NewLookupSource(nil, ConfigMap{
				"a": &KeyConfig{Configs: ConfigMap{
					"b": &KeyConfig{Template: returnT("translation-a-b")}},
				},
				"b": &KeyConfig{Template: returnT("translation-b")},
			}),
			message:     Msg("a.b"),
			translation: "translation-a-b",
		},
		{ // goes to parent
			lookupSource: NewLookupSource(
				NewLookupSource(nil, ConfigMap{ // parent
					"a": &KeyConfig{Configs: ConfigMap{
						"b": &KeyConfig{Template: returnT("parent-a-b")}},
					},
					"b": &KeyConfig{Template: returnT("parent-b")},
				}),
				ConfigMap{ // current
					"b": &KeyConfig{Template: returnT("translation-b")},
				}),
			message:     Msg("a.b"),
			translation: "parent-a-b",
		},
	}
}

func returnT(msg string) Template {
	return func(ctx any) string {
		return msg
	}
}
