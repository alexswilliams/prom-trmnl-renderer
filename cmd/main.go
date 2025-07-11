package main

import (
	"bytes"
	"image"
	"image/png"
	"log"
	"os"
	"os/signal"
	"prom-trmnl-renderer/lib/trmnl"
	"time"
)

func main() {
	drawAndUpload()

	var sigIntReceived = closeOnSigInt(make(chan bool, 1))
	ticker := time.NewTicker(20 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			drawAndUpload()
		case <-sigIntReceived:
			log.Println("Exiting...")
			return
		}
	}
}

func closeOnSigInt(channel chan bool) chan bool {
	var signals = make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	go func() { <-signals; println("Received SIGINT"); close(channel) }()
	return channel
}

func drawAndUpload() {
	outdoorTemps := trmnl.FetchLast48Hours("min(govee_temperature_celsius{alias=~\"Outside.*\"})")
	indoorTemps := trmnl.FetchLast48Hours("max(govee_temperature_celsius{alias!~\"Outside.*|Fridge|Car\"})")

	img := trmnl.NewCanvas()
	trmnl.RenderTempGraphs(img, outdoorTemps, indoorTemps)

	pngBytes := encodeToPng(img)
	trmnl.UploadToS3(pngBytes)

	//if err := os.WriteFile("/tmp/out.png", pngBytes, 0644); err != nil {
	//	log.Fatal(err)
	//}
}

func encodeToPng(img *image.Paletted) []byte {
	var encoded bytes.Buffer
	err := png.Encode(&encoded, img)
	if err != nil {
		log.Fatal(err)
	}
	return encoded.Bytes()
}
