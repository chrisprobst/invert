package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

func main() {
	runtime.GOMAXPROCS(4)
	var group sync.WaitGroup
	for _, arg := range os.Args[1:] {
		group.Add(1)
		go func() {
			defer group.Done()
			path, err := filepath.Abs(arg)
			if err != nil {
				log.Fatalln(err)
			}
			log.Println(path)

			f, err := os.Open(path)
			if err != nil {
				log.Fatalln(err)
			}
			log.Println("Opened file")

			img, err := jpeg.Decode(f)
			if err != nil {
				log.Fatalln(err)
			}
			log.Println("Decoded image file")

			newImg := image.NewRGBA(img.Bounds())
			for x := 0; x < img.Bounds().Dx(); x++ {
				for y := 0; y < img.Bounds().Dy(); y++ {
					r, g, b, _ := img.At(x, y).RGBA()
					r = 255 - r
					g = 255 - g
					b = 255 - b

					newImg.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), 255})
				}
			}

			fout, err := os.Create(path + "_inverted.jpg")
			if err != nil {
				log.Fatalln(err)
			}
			if err := jpeg.Encode(fout, newImg, nil); err != nil {
				log.Fatalln(err)
			}
		}()
	}
	group.Wait()
}
