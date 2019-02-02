# paint2hex
A collection of scripts to turn a set of physical paints into hexadecimal color codes

Reduce size of images

```
mogrify -fuzz 10% -transparent none -quality 70 -depth 8 -colors 256 *png
mogrify -quality 20% *jpg
```
