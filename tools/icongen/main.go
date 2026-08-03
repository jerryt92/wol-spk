package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
	"path/filepath"
)

var sizes = []int{16, 24, 32, 48, 64, 72, 128, 256}

var (
	brandBlue = color.RGBA{25, 118, 210, 255}
	white     = color.RGBA{255, 255, 255, 255}
)

func main() {
	dirs := []string{
		filepath.Join("synology", "ui", "images"),
		filepath.Join("cmd", "wolmanager", "images"),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			panic(err)
		}
	}

	for _, size := range sizes {
		img := renderIcon(size)
		for _, dir := range dirs {
			writePNG(filepath.Join(dir, "icon_"+itoa(size)+".png"), img)
		}
		switch size {
		case 64:
			writePNG(filepath.Join("synology", "PACKAGE_ICON.PNG"), img)
		case 256:
			writePNG(filepath.Join("synology", "PACKAGE_ICON_256.PNG"), img)
		}
	}
}

func renderIcon(size int) *image.RGBA {
	scale := 4
	if size <= 32 {
		scale = 8
	}
	large := image.NewRGBA(image.Rect(0, 0, size*scale, size*scale))
	canvas := float64(size * scale)
	drawRoundedRect(large, 0, 0, canvas, canvas, canvas*0.235, brandBlue)

	drawPower(large, canvas)
	drawNetwork(large, canvas, size)

	return downsample(large, size, scale)
}

func drawPower(img *image.RGBA, canvas float64) {
	cx := canvas * 0.5
	cy := canvas * 0.39
	r := canvas * 0.13
	stroke := canvas * 0.045
	lineTop := canvas * 0.205
	lineBottom := canvas * 0.355

	drawRingGap(img, cx, cy, r, stroke, 228, 312, white)
	drawCapsule(img, cx, lineTop, cx, lineBottom, stroke, white)
}

func drawNetwork(img *image.RGBA, canvas float64, finalSize int) {
	cx := canvas * 0.5
	y := canvas * 0.705
	outer := canvas * 0.235
	wing := canvas * 0.068
	stroke := canvas * 0.034
	dotR := canvas * 0.021

	if finalSize <= 16 {
		return
	}
	if finalSize <= 24 {
		y = canvas * 0.715
		dotR = canvas * 0.027
		for _, dx := range []float64{-0.09, 0, 0.09} {
			drawDisc(img, cx+canvas*dx, y, dotR, white)
		}
		return
	}
	if finalSize <= 32 {
		outer = canvas * 0.225
		wing = canvas * 0.06
		stroke = canvas * 0.042
		dotR = canvas * 0.023
	}

	drawCapsule(img, cx-outer, y, cx-outer+wing, y-wing, stroke, white)
	drawCapsule(img, cx-outer, y, cx-outer+wing, y+wing, stroke, white)
	drawCapsule(img, cx+outer, y, cx+outer-wing, y-wing, stroke, white)
	drawCapsule(img, cx+outer, y, cx+outer-wing, y+wing, stroke, white)

	for _, dx := range []float64{-0.078, 0, 0.078} {
		drawDisc(img, cx+canvas*dx, y, dotR, white)
	}
}

func drawRoundedRect(img *image.RGBA, x, y, w, h, r float64, c color.RGBA) {
	b := img.Bounds()
	for py := b.Min.Y; py < b.Max.Y; py++ {
		for px := b.Min.X; px < b.Max.X; px++ {
			fx := float64(px) + 0.5
			fy := float64(py) + 0.5
			dx := math.Max(math.Max(x+r-fx, 0), fx-(x+w-r))
			dy := math.Max(math.Max(y+r-fy, 0), fy-(y+h-r))
			dist := math.Hypot(dx, dy) - r
			coverage := clamp01(0.5 - dist)
			blendPixel(img, px, py, c, coverage)
		}
	}
}

func drawRingGap(img *image.RGBA, cx, cy, r, stroke, gapStart, gapEnd float64, c color.RGBA) {
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			fx := float64(x) + 0.5
			fy := float64(y) + 0.5
			angle := math.Atan2(fy-cy, fx-cx) * 180 / math.Pi
			if angle < 0 {
				angle += 360
			}
			if angle >= gapStart && angle <= gapEnd {
				continue
			}
			dist := math.Abs(math.Hypot(fx-cx, fy-cy)-r) - stroke/2
			coverage := clamp01(0.5 - dist)
			blendPixel(img, x, y, c, coverage)
		}
	}
}

func drawCapsule(img *image.RGBA, x1, y1, x2, y2, width float64, c color.RGBA) {
	b := img.Bounds()
	vx := x2 - x1
	vy := y2 - y1
	lenSq := vx*vx + vy*vy
	for py := b.Min.Y; py < b.Max.Y; py++ {
		for px := b.Min.X; px < b.Max.X; px++ {
			fx := float64(px) + 0.5
			fy := float64(py) + 0.5
			t := 0.0
			if lenSq > 0 {
				t = ((fx-x1)*vx + (fy-y1)*vy) / lenSq
				t = clamp01(t)
			}
			nx := x1 + vx*t
			ny := y1 + vy*t
			dist := math.Hypot(fx-nx, fy-ny) - width/2
			coverage := clamp01(0.5 - dist)
			blendPixel(img, px, py, c, coverage)
		}
	}
}

func drawDisc(img *image.RGBA, cx, cy, r float64, c color.RGBA) {
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			dist := math.Hypot(float64(x)+0.5-cx, float64(y)+0.5-cy) - r
			coverage := clamp01(0.5 - dist)
			blendPixel(img, x, y, c, coverage)
		}
	}
}

func downsample(src *image.RGBA, size, scale int) *image.RGBA {
	dst := image.NewRGBA(image.Rect(0, 0, size, size))
	area := float64(scale * scale)
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			var r, g, b, a uint32
			for sy := 0; sy < scale; sy++ {
				for sx := 0; sx < scale; sx++ {
					cr, cg, cb, ca := src.At(x*scale+sx, y*scale+sy).RGBA()
					r += cr
					g += cg
					b += cb
					a += ca
				}
			}
			dst.SetRGBA(x, y, color.RGBA{
				R: uint8(float64(r) / area / 257),
				G: uint8(float64(g) / area / 257),
				B: uint8(float64(b) / area / 257),
				A: uint8(float64(a) / area / 257),
			})
		}
	}
	return dst
}

func blendPixel(img draw.Image, x, y int, c color.RGBA, alpha float64) {
	if alpha <= 0 {
		return
	}
	alpha = clamp01(alpha) * float64(c.A) / 255
	old := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
	inv := 1 - alpha
	img.Set(x, y, color.RGBA{
		R: uint8(float64(c.R)*alpha + float64(old.R)*inv + 0.5),
		G: uint8(float64(c.G)*alpha + float64(old.G)*inv + 0.5),
		B: uint8(float64(c.B)*alpha + float64(old.B)*inv + 0.5),
		A: uint8(255*alpha + float64(old.A)*inv + 0.5),
	})
}

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

func writePNG(path string, img image.Image) {
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	if err := png.Encode(file, img); err != nil {
		_ = file.Close()
		panic(err)
	}
	if err := file.Close(); err != nil {
		panic(err)
	}
}

func itoa(value int) string {
	if value == 0 {
		return "0"
	}
	buf := make([]byte, 0, 8)
	for value > 0 {
		buf = append([]byte{byte('0' + value%10)}, buf...)
		value /= 10
	}
	return string(buf)
}
