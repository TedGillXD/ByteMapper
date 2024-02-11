package TargetGenerator

const (
	CPP = "cpp"
)

type TargetGenerator struct {
	targetLanguage string
	parsedObj      interface{}
}

func NewTargetGenerator(language string) *TargetGenerator {
	return &TargetGenerator{targetLanguage: language}
}

func (t *TargetGenerator) getCppFile(header string, source string) {

}
