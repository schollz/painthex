package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/oliamb/cutter"
)

func main() {
	jpgFiles, err := filepath.Glob("*_*_*.jpg")
	if err != nil {
		log.Fatal(err)
	}
	for _, inPath := range jpgFiles {
		err = crop(inPath, inPath+".1.png", image.Point{60, 300})
		if err != nil {
			log.Fatal(err)
		}
		err = crop(inPath, inPath+".2.png", image.Point{500, 300})
		if err != nil {
			log.Fatal(err)
		}
		err = crop(inPath, inPath+".3.png", image.Point{900, 700})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func crop(inPath string, outPath string, tlPoint image.Point) (err error) {
	fi, err := os.Open(inPath)
	if err != nil {
		return
	}
	defer fi.Close()
	fmt.Println(fi)
	img, err := jpeg.Decode(fi)
	if err != nil {
		return
	}

	cImg, err := cutter.Crop(img, cutter.Config{
		Height: 50,             // height in pixel or Y ratio(see Ratio Option below)
		Width:  50,             // width in pixel or X ratio
		Mode:   cutter.TopLeft, // Accepted Mode: TopLeft, Centered
		Anchor: tlPoint,        // Position of the top left point
		// Anchor: image.Point{500, 300}, // Position of the top left point
		// Anchor:  image.Point{900, 700}, // Position of the top left point
		Options: 0, // Accepted Option: Ratio
	})
	if err != nil {
		return
	}

	fo, err := os.Create(outPath)
	if err != nil {
		return
	}
	defer fo.Close()

	err = png.Encode(fo, cImg)
	return
}
