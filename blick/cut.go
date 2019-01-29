package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/go-playground/colors"
	"github.com/oliamb/cutter"
)

func main() {
	jpgFiles, err := filepath.Glob("*_*_*.jpg")
	if err != nil {
		log.Fatal(err)
	}
	for i, inPath := range jpgFiles {
		fmt.Println(inPath)
		if i > 0 {
			break
		}
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
	cnrgba, err := AverageImageColor(outPath)
	rgb, err := colors.RGB(cnrgba.R, cnrgba.G, cnrgba.B)
	if err != nil {
		return
	}
	fmt.Println(rgb.ToHEX().String())
	return
}

func AverageImageColor(inPath string) (cnrgba color.NRGBA, err error) {
	fi, err := os.Open(inPath)
	if err != nil {
		return
	}
	defer fi.Close()
	im, err := png.Decode(fi)
	if err != nil {
		return
	}
	rgba := imageToRGBA(im)
	size := rgba.Bounds().Size()
	w, h := size.X, size.Y
	var r, g, b int
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := rgba.RGBAAt(x, y)
			r += int(c.R)
			g += int(c.G)
			b += int(c.B)
		}
	}
	r /= w * h
	g /= w * h
	b /= w * h
	cnrgba = color.NRGBA{uint8(r), uint8(g), uint8(b), 255}
	return
}

func imageToRGBA(src image.Image) *image.RGBA {
	dst := image.NewRGBA(src.Bounds())
	draw.Draw(dst, dst.Rect, src, image.ZP, draw.Src)
	return dst
}
