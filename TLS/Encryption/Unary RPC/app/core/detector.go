package core

import (
	"fmt"

	"github.com/abadojack/whatlanggo"
)

type LanguageResult struct {
	Name       string
	Script     string
	Confidence float64
}

func (output LanguageResult) String() string {
	return fmt.Sprintf("Language: %s\tScript: %s\tConfidence: %f", output.Name, output.Script, output.Confidence)
}

func NewLanguageResult(text string) LanguageResult {
	info := whatlanggo.Detect(text)
	return LanguageResult{
		Name:       info.Lang.String(),
		Script:     whatlanggo.Scripts[info.Script],
		Confidence: info.Confidence,
	}
}
