package wut

type (
	LoadOptions interface {
		TemplateFactory(tf TemplateFactory) LoadOptions
		AddFiles(file ...string) LoadOptions
		AsDefault() error
		Build() (LangFactory, error)
	}

	options struct {
		templateFactory TemplateFactory
		files           []string
	}
)

var defaultFactory LangFactory

func SetDefaultFactory(f LangFactory) {
	defaultFactory = f
}

func Lang(code string) LangSource {
	return defaultFactory.Lang(code)
}

func LangNoFallback(code string) LangSource {
	return defaultFactory.LangNoFallback(code)
}

func Setup() LoadOptions {
	return &options{
		templateFactory: nil,
		files:           make([]string, 0),
	}
}

func (o *options) TemplateFactory(tf TemplateFactory) LoadOptions {
	o.templateFactory = tf
	return o
}

func (o *options) AddFiles(file ...string) LoadOptions {
	o.files = append(o.files, file...)
	return o
}

func (o *options) Build() (LangFactory, error) {
	factory := o.templateFactory
	if factory == nil {
		factory = &DefaultTemplateFactory{}
	}

	parsedFiles := make([]*LangFile, 0)
	for _, file := range o.files {
		langFile, err := ReadFile(file)
		if err != nil {
			return nil, err
		}
		parsedFiles = append(parsedFiles, langFile)
	}

	return NewLangFactory(factory, parsedFiles...)
}

func (o *options) AsDefault() error {
	factory, err := o.Build()
	if err != nil {
		return err
	}
	SetDefaultFactory(factory)
	return nil
}
