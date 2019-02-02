# painthex

See [painthex.schollz.com](https://painthex.schollz.com).

A collection of scripts to turn a set of physical paints into hexadecimal color codes.

![Example](https://user-images.githubusercontent.com/6550035/52167707-4f7c9e00-26d4-11e9-9361-b7dcd8f0e23a.png)

Once the colors are gathered go into `blick` or `golden` and then run `go run cut.go`. Afterwards, reduce size of images with

```bash
mogrify -fuzz 10% -transparent none -quality 70 -depth 8 -colors 256 *png
mogrify -quality 20% *jpg
```


