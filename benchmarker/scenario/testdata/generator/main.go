package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const (
	width  = 640
	height = 640
)

//go:generate go run main.go -dir ../icons -width 640 -height 640
//go:generate go run main.go -dir ../covers -width 1500 -height 600
//go:generate go run main.go -dir ../photos -width 1800 -height 1800

// Usage: go run generator/main.go -dir icons -width 640 -height 640
//
//	or go run testdata/generator.go -dir covers -width 1500 -height 600
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	directory := flag.String("dir", "icons", "directory to save images")
	width := flag.Int("width", 640, "width of image")
	height := flag.Int("height", 640, "height of image")
	flag.Parse()

	if _, err := os.Stat(*directory); !errors.Is(err, os.ErrNotExist) {
		if err := os.RemoveAll(*directory); err != nil {
			panic(err)
		}
	}
	if err := os.Mkdir(*directory, 0755); err != nil {
		panic(err)
	}

	generatedCodes := map[string]struct{}{}

	for n := 0; n < 100; n++ {
		palettes, err := getColorhuntPalettes(ctx, n)
		if err != nil {
			panic(err)
		}
		if len(palettes) == 0 {
			break
		}

		for _, palette := range palettes {
			if _, ok := generatedCodes[palette.Code]; ok {
				continue
			}
			generatedCodes[palette.Code] = struct{}{}

			colors := []color.NRGBA{
				parseRGB(palette.Code[0:6]),
				parseRGB(palette.Code[6:12]),
				parseRGB(palette.Code[12:18]),
				parseRGB(palette.Code[18:24]),
			}

			img := image.NewNRGBA(image.Rect(0, 0, *width, *height))

			for y := 0; y < img.Rect.Dy(); y++ {
				idx := int(math.Floor(float64(y) / float64(img.Rect.Dy()) * 4))
				color := colors[idx]
				for x := 0; x < img.Rect.Dx(); x++ {
					img.Set(x, y, color)
				}
			}

			log.Printf("Generate: %s", palette.Code)
			fd, err := os.Create(fmt.Sprintf("%s/palette-%s.png", *directory, palette.Code))
			if err != nil {
				panic(err)
			}

			if err := png.Encode(fd, img); err != nil {
				panic(err)
			}
			fd.Close()
		}
	}
}

type ColorPalette struct {
	Code  string
	Likes string
	Date  string
}

func getColorhuntPalettes(ctx context.Context, step int) ([]ColorPalette, error) {
	values := url.Values{}
	values.Add("sort", "random")
	values.Add("timeframe", "30")
	values.Add("tags", "")
	values.Add("step", strconv.Itoa(step))
	body := strings.NewReader(values.Encode())

	req, err := http.NewRequestWithContext(ctx, "POST", "https://colorhunt.co/php/feed.php", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Encoding", "deflate")
	// req.Header.Set("X-Requested-With", "XMLHttpRequest")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	log.Printf("status: %s", res.Status)

	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	colorPalettes := []ColorPalette{}
	if err := json.NewDecoder(bytes.NewReader(rawBody)).Decode(&colorPalettes); err != nil {
		return nil, err
	}

	return colorPalettes, nil
}

func parseRGB(rgb string) color.NRGBA {
	rgb = strings.TrimPrefix(rgb, "#")
	r, _ := strconv.ParseUint(rgb[0:2], 16, 8)
	g, _ := strconv.ParseUint(rgb[2:4], 16, 8)
	b, _ := strconv.ParseUint(rgb[4:6], 16, 8)
	return color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}
}
