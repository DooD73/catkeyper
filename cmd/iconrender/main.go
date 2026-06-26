package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"catkeyper/internal/iconrender"
)

func main() {
	in := flag.String("in", "assets/openmoji-cat-face.svg", "source SVG")
	out := flag.String("out", "build/icon.iconset", "destination .iconset directory")
	flag.Parse()

	if err := os.MkdirAll(*out, 0o755); err != nil {
		log.Fatalf("create iconset: %v", err)
	}

	for name, size := range iconrender.IconSizes {
		img, err := iconrender.RenderIcon(*in, size)
		if err != nil {
			log.Fatalf("render %s: %v", name, err)
		}
		if err := iconrender.WritePNG(filepath.Join(*out, name), img); err != nil {
			log.Fatalf("write %s: %v", name, err)
		}
	}
}
