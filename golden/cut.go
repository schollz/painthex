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
	for i, inPath := range jpgFiles {
		fmt.Println(i, inPath)
		colorName := strings.Split(inPath, "_")[1]
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

var html = `<html>
<head>
	<style>
	.col-row {
display: grid;
grid-auto-flow: row;
grid-template-columns: repeat(12, 1fr);
grid-column-gap: 1em;
grid-row-gap: 1em;
margin: auto;
max-width: 30em;
}


/*
COLUMN-CONTROLLED ALIGNMENT
*/
.col-top,
.xs\=col-top {
align-self: start;
}
.col-center,
.xs\=col-center {
align-self: center;
}
.col-bottom,
.xs\=col-bottom {
align-self: end;
}
.col-stretch,
.xs\=col-stretch {
align-self: stretch;
}


/*
EXTRA SMALL
*/
.col-1,
.xs\=col-1 {
grid-column: span 1;
}
.col-2,
.xs\=col-2 {
grid-column: span 2;
}
.col-3,
.xs\=col-3 {
grid-column: span 3;
}
.col-4,
.xs\=col-4 {
grid-column: span 4;
}
.col-5,
.xs\=col-5 {
grid-column: span 5;
}
.col-6,
.xs\=col-6 {
grid-column: span 6;
}
.col-7,
.xs\=col-7 {
grid-column: span 7;
}
.col-8,
.xs\=col-8 {
grid-column: span 8;
}
.col-9,
.xs\=col-9 {
grid-column: span 9;
}
.col-10,
.xs\=col-10 {
grid-column: span 10;
}
.col-11,
.xs\=col-11 {
grid-column: span 11;
}
.col-12,
.xs\=col-12 {
grid-column: span 12;
}

@media only screen and (min-width: 48em) {
.col-row {
max-width: 48em;
}

/*
SMALL
*/
.sm\=col-1 {
grid-column: span 1;
}
.sm\=col-2 {
grid-column: span 2;
}
.sm\=col-3 {
grid-column: span 3;
}
.sm\=col-4 {
grid-column: span 4;
}
.sm\=col-5 {
grid-column: span 5;
}
.sm\=col-6 {
grid-column: span 6;
}
.sm\=col-7 {
grid-column: span 7;
}
.sm\=col-8 {
grid-column: span 8;
}
.sm\=col-9 {
grid-column: span 9;
}
.sm\=col-10 {
grid-column: span 10;
}
.sm\=col-11 {
grid-column: span 11;
}
.sm\=col-12 {
grid-column: span 12;
}

.sm\=col-top {
align-self: start;
}
.sm\=col-center {
align-self: center;
}
.sm\=col-bottom {
align-self: end;
}
.sm\=col-stretch {
align-self: stretch;
}

}

@media only screen and (min-width: 64em) {
.col-row {
max-width: 64em;
}

/*
MEDIUM
*/
.md\=col-1 {
grid-column: span 1;
}
.md\=col-2 {
grid-column: span 2;
}
.md\=col-3 {
grid-column: span 3;
}
.md\=col-4 {
grid-column: span 4;
}
.md\=col-5 {
grid-column: span 5;
}
.md\=col-6 {
grid-column: span 6;
}
.md\=col-7 {
grid-column: span 7;
}
.md\=col-8 {
grid-column: span 8;
}
.md\=col-9 {
grid-column: span 9;
}
.md\=col-10 {
grid-column: span 10;
}
.md\=col-11 {
grid-column: span 11;
}
.md\=col-12 {
grid-column: span 12;
}

.md\=col-top {
align-self: start;
}
.md\=col-center {
align-self: center;
}
.md\=col-bottom {
align-self: end;
}
.md\=col-stretch {
align-self: stretch;
}
}

@media only screen and (min-width: 75em) {
.col-row {
max-width: 75em;
}

/*
LARGE
*/
.lg\=col-1 {
grid-column: span 1;
}
.lg\=col-2 {
grid-column: span 2;
}
.lg\=col-3 {
grid-column: span 3;
}
.lg\=col-4 {
grid-column: span 4;
}
.lg\=col-5 {
grid-column: span 5;
}
.lg\=col-6 {
grid-column: span 6;
}
.lg\=col-7 {
grid-column: span 7;
}
.lg\=col-8 {
grid-column: span 8;
}
.lg\=col-9 {
grid-column: span 9;
}
.lg\=col-10 {
grid-column: span 10;
}
.lg\=col-11 {
grid-column: span 11;
}
.lg\=col-12 {
grid-column: span 12;
}

.lg\=col-top {
align-self: start;
}
.lg\=col-center {
align-self: center;
}
.lg\=col-bottom {
align-self: end;
}
.lg\=col-stretch {
align-self: stretch;
}

}

.col-fluid {
max-width: 100vw;
}
</style>
</head>
<body>
		{{ range . }}
		<div class='col-row'>
			<div class='lg=col-12'><h2>{{.Name}}</h2></div>
		</div>
		<div class='col-row'>
			<div class='lg=col-12'><img src="{{.Path}}" width=360></div>
		</div>
		<div class='col-row'>
				<div class='lg=col-4'><img src="{{.Path}}.1.png" width=120></div>
				<div class='lg=col-4'><img src="{{.Path}}.2.png" width=120></div>
				<div class='lg=col-4'><img src="{{.Path}}.3.png" width=120></div>
		</div>
		<div class='col-row'>
			{{ range .Hexes}}
			<div class='lg=col-4'><div style="float:left;width:120px;height:120px;background:{{.}};"></div></div>
			{{ end}}
		</div>
		{{ end }}
</body>
</html>`
