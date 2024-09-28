package joran

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadLangFile(t *testing.T) {
	filename, _, err := createTestFile(t, "test.toml", `
language = "de"
fallback = "en"

[validation]
email = 'Email error'

[validation.Class]
email = 'Class email error'
`)
	if err != nil {
		t.Fatal(err)
	}

	file, err := ReadLangFile(filename)
	if err != nil {
		t.Fatal(err)
	}

	if file.Language != "de" {
		t.Errorf("got %q, want 'de_de'", file.Language)
	}

	if file.Fallback != "en" {
		t.Errorf("got %q, want 'en'", file.Fallback)
	}

	if file.Values["validation"].(map[string]any)["email"] != "Email error" {
		t.Errorf("vlidation.email missing")
	}

	if file.Values["validation"].(map[string]any)["Class"].(map[string]any)["email"] != "Class email error" {
		t.Errorf("vlidation.Class.email missing")
	}
}

func TestLangFile_Include(t *testing.T) {
	filename, _, err := createTestFile(t, "test.toml", `
language = "de"
fallback = "en"

[validation]
email = 'Email error'

[validation.Class]
email = 'Class email error'
`)
	if err != nil {
		t.Fatal(err)
	}
	filename2, _, err := createTestFile(t, "test2.toml", `
language = "de"
fallback = "en"

[validation]
email2 = 'Email error'

[other]
email = 'other error'
`)
	if err != nil {
		t.Fatal(err)
	}

	file, err := ReadLangFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	file2, err := ReadLangFile(filename2)
	if err != nil {
		t.Fatal(err)
	}

	err = file.Include(file2)
	if err != nil {
		t.Fatal(err)
	}
}

func createTestFile(t *testing.T, name, content string) (string, *os.File, error) {
	dir := t.TempDir()
	path := filepath.Join(dir, name)
	file, err := os.Create(path)
	if err != nil {
		return "", nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)
	_, err = file.Write([]byte(content))
	if err != nil {
		return "", nil, err
	}
	return path, file, nil
}
