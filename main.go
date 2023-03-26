package main

import (
	"flag"
	"image"
	"image/color"
	"image/png"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readImage(path string) image.Image {
	f, err := os.Open(path)
	check(err)
	defer f.Close()
	i, _, err := image.Decode(f)
	check(err)
	return i
}

func main() {
	var in, mask, out string
	flags := flag.NewFlagSet("demasker", flag.ExitOnError)
	flags.StringVar(&in, "in", "", "input file")
	flags.StringVar(&mask, "mask", "", "mask file")
	flags.StringVar(&out, "out", "", "output file")
	err := flags.Parse(os.Args[1:])
	if in == "" || mask == "" || out == "" || err != nil {
		flags.Usage()
		os.Exit(1)
	}

	var (
		inImg   = readImage(in)
		maskImg = readImage(mask)
		outImg  = image.NewRGBA(inImg.Bounds())
	)

	if inImg.Bounds() != maskImg.Bounds() {
		panic("image bounds do not match")
	}

	// read every pixel in the mask and set the corresponding pixel in the
	// input image to transparent based on darkness
	for x := 0; x < inImg.Bounds().Dx(); x++ {
		for y := 0; y < inImg.Bounds().Dy(); y++ {

			var (
				maskColor = color.RGBAModel.Convert(maskImg.At(x, y)).(color.RGBA)
				inColor   = color.RGBAModel.Convert(inImg.At(x, y)).(color.RGBA)
				a         = 255 - maskColor.R // assuming grayscale mask
				r         = uint8(uint16(inColor.R) * uint16(a) / 255)
				g         = uint8(uint16(inColor.G) * uint16(a) / 255)
				b         = uint8(uint16(inColor.B) * uint16(a) / 255)
			)

			outImg.Set(x, y, color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: a})
		}
	}

	// write the result to the output file
	f, err := os.Create(out)
	check(err)
	defer f.Close()

	err = png.Encode(f, outImg)
	check(err)
}
