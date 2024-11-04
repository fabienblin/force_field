package main

import (
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"github.com/aquilax/go-perlin"
)

const (
	IMAGE_WIDTH        int     = 500
	IMAGE_HEIGHT       int     = 500
	PERLIN_ALPHA       float64 = 10
	PERLIN_BETA        float64 = 3
	PERLIN_N           int32   = 3
	PERLIN_ZOOM        float64 = 100
	PARTICULE_SPEED    float64 = 1.0
	PARTICULE_NB       int     = 100
	// PERLIN_Z_INCREMENT float64 = 0.01
)

var PARTICULE_COLOR color.Color = color.RGBA{R: 10, G: 10, B: 250, A: 255}
var forceField *perlin.Perlin
var particules []*Particule

type Particule struct {
	X float64
	Y float64
}

func init() {
	seed := rand.Int63()
	rand.NewSource(seed)
	log.Printf("seed : %d", seed)

	forceField = perlin.NewPerlin(PERLIN_ALPHA, PERLIN_BETA, PERLIN_N, seed)
	initParticules()
}

func initParticules() {
	for i := 0; i < PARTICULE_NB; i++ {
		particules = append(particules, &Particule{X: rand.Float64() * float64(IMAGE_WIDTH), Y: rand.Float64() * float64(IMAGE_HEIGHT)})
	}
}

func fillImageWithBlack(img *image.RGBA) {
	for y := 0; y < IMAGE_HEIGHT; y++ {
		for x := 0; x < IMAGE_WIDTH; x++ {
			img.Set(x, y, color.Black)
		}
	}
}

// refreshImage modifies the underlying RGBA image buffer and updates the canvas image.
func refreshImage(canvasImage *canvas.Image) {
	rgba := canvasImage.Image.(*image.RGBA)
	// var z float64

	for {
		// fillImageWithBlack(rgba)
		// z += PERLIN_Z_INCREMENT

		for _, particule := range particules {
			// force := forceField.Noise3D(float64(particule.X)/PERLIN_ZOOM, float64(particule.Y)/PERLIN_ZOOM, z)
			force := forceField.Noise2D(float64(particule.X)/PERLIN_ZOOM, float64(particule.Y)/PERLIN_ZOOM)

			angle := force * 2 *math.Pi
			particule.X += (math.Cos(angle) * PARTICULE_SPEED)
			particule.Y += (math.Sin(angle) * PARTICULE_SPEED)

			rgba.Set(int(particule.X), int(particule.Y), PARTICULE_COLOR)
		}

		// Refresh the canvas image to apply the changes
		canvasImage.Refresh()

		// Add a delay for smoother transitions
		time.Sleep(50 * time.Millisecond)
	}
}

func main() {
	// Create a new Fyne application
	app := app.New()
	window := app.NewWindow("Color Changing Image")

	// Create an RGBA image that can be modified
	r := image.Rect(0, 0, IMAGE_WIDTH, IMAGE_HEIGHT)
	rgba := image.NewRGBA(r)

	// Create a Fyne canvas image from the RGBA image
	canvasImage := canvas.NewImageFromImage(rgba)
	canvasImage.FillMode = canvas.ImageFillOriginal

	// Use a goroutine to continuously modify the image over time
	go refreshImage(canvasImage)

	// Set the canvas image as the window content
	window.SetContent(canvasImage)
	window.ShowAndRun()
}
