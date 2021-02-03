package main

import (
	"fmt"

	"github.com/anthonynsimon/bild/effect"
	"github.com/anthonynsimon/bild/imgio"
)

func main() {
	img, err := imgio.Open("assets/grimmjow.jpg")
	if err != nil {
		fmt.Println(err)
		return
	}

	result := effect.EdgeDetection(img, 2.0)

	if err := imgio.Save("assets/result2.jpg", result, imgio.JPEGEncoder(100)); err != nil {
		fmt.Println(err)
		return
	}
}
