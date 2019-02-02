# painthex

See [painthex.schollz.com](https://painthex.schollz.com).

A collection of scripts to turn a set of physical paints into hexadecimal color codes.

Once the colors are gathered go into `blick` or `golden` and then run `go run cut.go`. Afterwards, reduce size of images with

```bash
mogrify -fuzz 10% -transparent none -quality 70 -depth 8 -colors 256 *png
mogrify -quality 20% *jpg
```


