package main

import (
	"bytes"
	"embed"
	"flag"
	"image"
	"image/draw"
	"image/png"
	"math/rand"
	"os"

	"github.com/dolmen-go/kittyimg"
)

//go:embed birds/*.png
var birds embed.FS

var (
	width   = flag.Int("width", 1920, "width of the image")
	height  = flag.Int("height", 1080, "height of the image")
	density = flag.Int("density", 100, "density of the birds")
	output  = flag.String("output", "image.png", "output file name")
	display = flag.Bool("display", false, "display the image using kitty image protocol after generation")
)

func main() {
	flag.Parse()

	img := image.NewRGBA(image.Rect(0, 0, *width, *height))
	draw.Draw(img, img.Bounds(), image.Black, image.Point{}, draw.Src)

	used := make(map[image.Rectangle]bool)
	for i := 0; i < *density; i++ {
		bird := randomBird()
	outer:
		for {
			x, y := rand.Intn(*width), rand.Intn(*height)
			rect := image.Rect(x, y, x+bird.Bounds().Dx(), y+bird.Bounds().Dy())
			inflated := rect.Inset(-30)
			for r := range used {
				if inflated.Overlaps(r) {
					continue outer
				}
			}
			used[rect] = true
			draw.Draw(img, rect, bird, image.Point{}, draw.Over)
			break
		}
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		panic(err)
	}
	if err := os.WriteFile(*output, buf.Bytes(), 0o644); err != nil {
		panic(err)
	}

	if *display {
		kittyimg.Fprintln(os.Stdout, img)
	} else {
		println("done!")
	}
}

func randomBird() image.Image {
	files, _ := birds.ReadDir("birds")
	file := files[rand.Intn(len(files))]
	f, _ := birds.Open("birds/" + file.Name())
	img, _ := png.Decode(f)
	return img
}
