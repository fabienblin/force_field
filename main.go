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
	PARTICULE_NB       int     = 200
	PERLIN_Z_INCREMENT float64 = 0.001
)

var PARTICULE_COLOR color.Color = color.RGBA{R: 10, G: 10, B: 250, A: 255}
var forceField *perlin.Perlin
var particules []*Particule

type Particule struct {
	InitX         float64
	InitY         float64
	X             float64
	Y             float64
	IsOutOfBounds bool
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
		particule := &Particule{X: rand.Float64() * float64(IMAGE_WIDTH), Y: rand.Float64() * float64(IMAGE_HEIGHT)}
		particule.InitX = particule.X
		particule.InitY = particule.Y

		particules = append(particules, particule)
	}
}

func reInitParticules() {
	for _, particule := range particules {
		particule.X = particule.InitX
		particule.Y = particule.InitY
		particule.IsOutOfBounds = false
	}
}

func fillImageWithBlack(img *image.RGBA) {
	for y := 0; y < IMAGE_HEIGHT; y++ {
		for x := 0; x < IMAGE_WIDTH; x++ {
			img.Set(x, y, color.Black)
		}
	}
}

// return true if all particules are out of bounds
func isParticulesAllOutOfBounds() bool {
	for _, particule := range particules {
		if !particule.IsOutOfBounds {
			return false
		}
	}

	return true
}

// refreshImage modifies the underlying RGBA image buffer and updates the canvas image.
func refreshImage(canvasImage *canvas.Image) {
	rgba := canvasImage.Image.(*image.RGBA)
	var z float64

	for {
		fillImageWithBlack(rgba)
		for !isParticulesAllOutOfBounds() {
			for _, particule := range particules {
				if particule.IsOutOfBounds {
					continue
				}

				force := forceField.Noise3D(float64(particule.X)/PERLIN_ZOOM, float64(particule.Y)/PERLIN_ZOOM, z)

				angle := force * 2 * math.Pi
				particule.X += (math.Cos(angle) * PARTICULE_SPEED)
				particule.Y += (math.Sin(angle) * PARTICULE_SPEED)

				particule.IsOutOfBounds = particule.X < 0 || particule.Y < 0 || particule.X > float64(IMAGE_WIDTH) || particule.Y > float64(IMAGE_HEIGHT)

				rgba.Set(int(particule.X), int(particule.Y), PARTICULE_COLOR)
			}
		}

		reInitParticules()
		canvasImage.Refresh()
		z += PERLIN_Z_INCREMENT

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
