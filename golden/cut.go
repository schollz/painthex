package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-playground/colors"
	colorful "github.com/lucasb-eyer/go-colorful"
	"github.com/oliamb/cutter"
)

type Color struct {
	Hexes []string
	Path  template.HTMLAttr
	Name  string
}

func main() {
	jpgFiles, err := filepath.Glob("images/*_*_*.jpg")
	if err != nil {
		log.Fatal(err)
	}
	colorMap := make(map[string]string)
	colors := []Color{}
	finishedColors := make(map[string]struct{})
	for i, inPath := range jpgFiles {
		colorName := strings.Split(inPath, "_")[1]
		if _, ok := finishedColors[colorName]; ok {
			continue
		}
		fmt.Println(i, inPath)
		finishedColors[colorName] = struct{}{}
		color := Color{Name: colorName, Path: template.HTMLAttr(filepath.ToSlash(inPath)), Hexes: []string{"", "", "", "", ""}}
		var hexString string
		hexString, err = crop(inPath, inPath+".1.png", image.Point{150, 50})
		if err != nil {
			log.Fatal(err)
		}
		colorMap[colorName+"0"] = hexString
		color.Hexes[0] = hexString
		hexString, err = crop(inPath, inPath+".2.png", image.Point{450, 50})
		if err != nil {
			log.Fatal(err)
		}
		colorMap[colorName+"1"] = hexString
		color.Hexes[1] = hexString
		hexString, err = crop(inPath, inPath+".3.png", image.Point{800, 50})
		if err != nil {
			log.Fatal(err)
		}
		colorMap[colorName+"2"] = hexString
		color.Hexes[2] = hexString
		c1, _ := colorful.Hex("#ffffff")
		c2, _ := colorful.Hex(hexString)
		c3 := c1.BlendRgb(c2, 0.5)
		color.Hexes[3] = c3.Hex()

		c3 = c1.BlendRgb(c3, 0.5)
		color.Hexes[4] = c3.Hex()
		colors = append(colors, color)
	}

	colorMapBytes, err := json.MarshalIndent(colorMap, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("blick.json", colorMapBytes, 0644)
	if err != nil {
		log.Fatal(err)
	}

	var tpl bytes.Buffer
	t := template.Must(template.New("main").Parse(html))
	err = t.Execute(&tpl, colors)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("index.html", tpl.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}

	palette := ""
	for _, color := range colors {
		for i, hex := range color.Hexes {
			name := ""
			if i == 0 {
				name = "Heavy"
			} else if i >= 2 {
				name = fmt.Sprintf("Light%d", i)
			}
			palette += fmt.Sprintf("%s%s\t\t%s\n", color.Name, name, hex)
		}
	}
	palette = strings.TrimSpace(palette)
	err = ioutil.WriteFile("palette.txt", []byte(palette), 0644)
	if err != nil {
		log.Fatal(err)
	}

}

func crop(inPath string, outPath string, tlPoint image.Point) (hexString string, err error) {
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
		Height: 200,            // height in pixel or Y ratio(see Ratio Option below)
		Width:  200,            // width in pixel or X ratio
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
	hexString = rgb.ToHEX().String()
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

type ColorPalette struct {
	colorsNameLookupMap map[string]string
	colorsHexString     []string
	colorsRGB           []color.NRGBA
}

func FromNameMap(nameToHex map[string]string) (cp *ColorPalette, err error) {
	cp = new(ColorPalette)
	cp.colorsNameLookupMap = make(map[string]string)
	for name := range nameToHex {
		cp.colorsNameLookupMap[nameToHex[name]] = name
	}
	return
}

func ClosestColor(r, g, b uint8) (rc, gc, bc uint) {
	return
}

var html = `
{{ range .}}

<div class="card p1" style="margin-top: 2em;">
	<div class="card-body">
		<div class="row">
			<div class="col">
				<h2  class="display-4">{{.Name}}</h2>
			</div>
		</div>
		<div class='row'>
			<div class='col-sm-12 col-md-6'><img class="lazy" data-src="{{.Path}}" width="100%"></div>
		</div>
		<div class='row p1'>
			<div class='col p0'><img class="lazy" data-src="{{.Path}}.1.png" width=100% height=120px></div>
			<div class='col p0'><img class="lazy" data-src="{{.Path}}.2.png" width=100% height=120px></div>
			<div class='col p0'><img class="lazy" data-src="{{.Path}}.3.png" width=100% height=120px></div>
			<div class='col p0'>
				<div style="float:left;width:100%;height:120px;background:#ffffff;"></div>
			</div>
			<div class='col p0'>
				<div style="float:left;width:100%;height:120px;background:#ffffff;"></div>
			</div>

		</div>
		<div class='row p1'>
			{{ range .Hexes}}
			<div class='col p0'>
				<div style="height:120px;background:{{.}};"></div>
			</div>
			{{ end}}
		</div>
		<div class='row p1'>
		{{ range .Hexes}}
			<div class='col text-center'><code style="color:#000000">{{.}}</code></div>
			{{ end}}
		</div>

	</div>
</div>

{{ end }}`
