package trmnl

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"strconv"
	"time"
)

var (
	fontFile = loadFont()
	LONDON   = loadLondonTimezone()
)

func loadLondonTimezone() *time.Location {
	london, err := time.LoadLocation("Europe/London")
	if err != nil {
		log.Fatal(err)
	}
	return london
}

func NewCanvas() *image.Paletted {
	const width, height = 800, 480
	img := image.NewPaletted(image.Rect(0, 0, width, height), color.Palette{color.Black, color.White})
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.White)
		}
	}
	return img
}

func RenderTempGraphs(img *image.Paletted, outdoorTemps []float64, indoorTemps []float64) {
	annotationDrawer := fontDrawerForSize(img, fontFile, 11)
	titleDrawer := fontDrawerForSize(img, fontFile, 14)
	statDrawer := fontDrawerForSize(img, fontFile, 42)

	maxOutdoor := renderTempGraph(outdoorTemps, img, 200, 10, 20, annotationDrawer)
	maxIndoor := renderTempGraph(indoorTemps, img, 200, 10, 260, annotationDrawer)

	titleDrawer.Dot = fixed.P(650, 50)
	titleDrawer.DrawString("Outdoor 48h Max")
	statDrawer.Dot = fixed.P(650, 95)
	statDrawer.DrawString(strconv.FormatFloat(maxOutdoor, 'f', 1, 64) + "°")
	titleDrawer.Dot = fixed.P(650, 150)
	titleDrawer.DrawString("Outdoor Last")
	statDrawer.Dot = fixed.P(650, 195)
	statDrawer.DrawString(strconv.FormatFloat(outdoorTemps[len(outdoorTemps)-1], 'f', 1, 64) + "°")

	titleDrawer.Dot = fixed.P(650, 290)
	titleDrawer.DrawString("Indoor 48h Max")
	statDrawer.Dot = fixed.P(650, 335)
	statDrawer.DrawString(strconv.FormatFloat(maxIndoor, 'f', 1, 64) + "°")
	titleDrawer.Dot = fixed.P(650, 390)
	titleDrawer.DrawString("Indoor Last")
	statDrawer.Dot = fixed.P(650, 435)
	statDrawer.DrawString(strconv.FormatFloat(indoorTemps[len(indoorTemps)-1], 'f', 1, 64) + "°")

	annotationDrawer.Dot = fixed.P(650, 470)
	annotationDrawer.DrawString(time.Now().In(LONDON).Format("2006-01-02 15:04"))
}

func fontDrawerForSize(img *image.Paletted, loadedFont *truetype.Font, size float64) *font.Drawer {
	return &font.Drawer{
		Dst: img,
		Src: image.Black,
		Face: truetype.NewFace(loadedFont, &truetype.Options{
			Size:    size,
			Hinting: font.HintingFull,
		}),
	}
}

func renderTempGraph(temps []float64, img *image.Paletted, graphHeight int, graphLeft int, graphTop int, drawer *font.Drawer) float64 {
	var maxTemp = -1000.0
	var minTemp = +1000.0
	for i := 0; i < len(temps); i++ {
		if temps[i] < minTemp {
			minTemp = temps[i]
		}
		if temps[i] > maxTemp {
			maxTemp = temps[i]
		}
	}
	var tempRange = maxTemp - minTemp

	var lastY = -1000.0
	for i := 0; i < len(temps); i++ {
		x := i + graphLeft
		thisY := float64(graphTop) + float64(graphHeight)*0.90 - ((temps[i]-minTemp)/tempRange)*float64(graphHeight)*0.80
		img.Set(x, int(thisY), color.Black)

		if lastY != -1000.0 {
			stepToLastY := 1
			if thisY > lastY {
				stepToLastY = -1
			}
			for y := int(thisY); y != int(lastY); y += stepToLastY {
				img.Set(x, y, color.Black)
			}
		}

		img.Set(x, graphTop, color.Black)
		img.Set(x, graphTop+graphHeight, color.Black)
		lastY = thisY
	}

	step := 1.0
	if tempRange >= 10.0 {
		step = 5.0
	}
	minTempLine := math.Round((minTemp-step)/step) * step
	maxTempLine := math.Round((maxTemp+step)/step) * step
	for tempLine := minTempLine; tempLine < maxTempLine; tempLine += step {
		if tempLine > minTemp-tempRange*.1 && tempLine < maxTemp+tempRange*.1 {
			y := float64(graphTop) + float64(graphHeight)*0.90 - ((tempLine-minTemp)/tempRange)*float64(graphHeight)*0.80
			for i := 0; i < len(temps); i += 3 {
				x := i + graphLeft
				img.Set(x, int(y), color.Black)
			}
			drawer.Dot = fixed.P(graphLeft+len(temps)+2, int(y))
			drawer.DrawString(strconv.FormatFloat(tempLine, 'f', 0, 64) + "°")
		}
	}
	return maxTemp
}

func loadFont() *truetype.Font {
	fontFile, err := os.ReadFile("/usr/share/fonts/truetype/noto/NotoSans-Bold.ttf")
	if err != nil {
		log.Fatal(err)
	}
	loadedFont, err := truetype.Parse(fontFile)
	if err != nil {
		log.Fatal(err)
	}
	return loadedFont
}
