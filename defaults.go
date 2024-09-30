package wut

type (
	LoadOptions interface {

		// TemplateFactory sets which factory will provide the templates
		// that will be returned in the end. Default is DefaultTemplateFactory
		TemplateFactory(tf TemplateFactory) LoadOptions

		// AddFiles accepts a list of files to be parsed.
		// They need to at least have the language property defined in them.
		// If adding multiple files for the same language, their fallback language MUST match.
		// The 'language' and 'fallback' settings could be moved to a separate file in the future.
		AddFiles(file ...string) LoadOptions

		// AsDefault sets the configured factory from Build as the default factory to be used by the global methods available
		AsDefault() error

		// Build configures and returns the lang factory
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
