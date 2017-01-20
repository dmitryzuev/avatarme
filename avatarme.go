package main

import (
	"crypto/sha1"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/lucasb-eyer/go-colorful"
	"net/http"
	"log"
	"image"
	// "image/png"
	"image/color"
	"image/draw"
	"bytes"
	"image/jpeg"
	"strconv"
)

type Person struct{
	email string
	ip4 string
	publicKey string
}

func (p Person) String() string {
	return p.email + p.ip4 + p.publicKey
}

func (p Person) Hash() [sha1.Size]byte {
	return sha1.Sum([]byte(p.String()))
}

func main() {
	router := httprouter.New()
	router.GET("/avatar", Avatar)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func Avatar(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  p := Person{
		email: r.FormValue("email"),
		ip4: r.FormValue("ip4"),
		publicKey: r.FormValue("publicKey"),
	}
	hash := p.Hash()
	hash_str := fmt.Sprintf("%x", hash)
	hue_int, _ := strconv.ParseInt(hash_str[33:], 16, 64)
	// fmt.Fprintf(w, "%v", float64(hue_int / 0xffffff))

	bg_color := color.RGBA{240, 240, 240, 255}
	fg_color := colorful.Hsl(float64(hue_int)/0xfffff, 0.5, 0.7)

	m := image.NewRGBA(image.Rect(0, 0, 250, 250))
	for i := 0; i < 15; i++ {
		var rect_color color.Color
		switch hash_char, _ := strconv.ParseInt(string(hash_str[i]), 16, 64); hash_char % 2 {
		case 1:
			rect_color = bg_color
		case 0:
			rect_color = fg_color
		}

		if i < 5 {
			rectangle(2 * 50, i * 50, 50, 50, rect_color, m)
		} else if i < 10 {
			rectangle(1 * 50, (i - 5) * 50, 50, 50, rect_color, m);
      rectangle(3 * 50, (i - 5) * 50, 50, 50, rect_color, m);
		} else {
			rectangle(0 * 50, (i - 10) * 50, 50, 50, rect_color, m);
      rectangle(4 * 50, (i - 10) * 50, 50, 50, rect_color, m);
		}
	}

	// var img image.Image = m

	buffer := new(bytes.Buffer)
  if err := jpeg.Encode(buffer, m, nil); err != nil {
      log.Println("unable to encode image.")
  }

  w.Header().Set("Content-Type", "image/jpeg")
  w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
  if _, err := w.Write(buffer.Bytes()); err != nil {
      log.Println("unable to write image.")
  }
}

func rectangle(x int, y int, w int, h int, color color.Color, img *image.RGBA) {
	sub := img.SubImage(image.Rectangle{image.Point{x, y}, image.Point{x+w, y+h}}).(*image.RGBA)
	draw.Draw(sub, sub.Bounds(), &image.Uniform{color}, image.ZP, draw.Src)
	// var i, j;
  // for (i = x; i < x + w; i++) {
  //     for (j = y; j < y + h; j++) {
  //         image.buffer[image.index(i, j)] = color;
  //     }
  // }
}
