package main

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

func NewLanguage(info whatlanggo.Info) LanguageResult {
	return LanguageResult{
		Name:       info.Lang.String(),
		Script:     whatlanggo.Scripts[info.Script],
		Confidence: info.Confidence,
	}
}

func main() {
	info := whatlanggo.Detect("For weeks, as the stock market regularly climbed to records, investors wondered what it would take to snap Wall Street out of its blissful state")
	fmt.Println(NewLanguage(info))

	info2 := whatlanggo.Detect("Esta variante se propaga más rápido y podría conllevar mayor riesgo de infección, por eso, los investigadores creen que esta nueva variante se propagó de Brasil a Asia.")
	fmt.Println(NewLanguage(info2))

	info3 := whatlanggo.Detect("Президент Владимир Зеленский 29 января встретился с делегацией швейцарской компании Stadler Rail AG во главе с председателем совета директоров")
	fmt.Println(NewLanguage(info3))
}
