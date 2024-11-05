package characters

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type Car struct {
	Image *canvas.Image
}

func NewCar(container *fyne.Container, isMoving bool) *Car {
	imagePath := "assets/car.png"
	if isMoving {
		imagePath = "assets/car_move.png"
	}

	carImage := canvas.NewImageFromFile(imagePath)
	carImage.Resize(fyne.NewSize(40, 70))
	carImage.Move(fyne.NewPos(20, 230))

	// Solo agregar la imagen al contenedor si no es `nil`
	if container != nil {
		container.Add(carImage)
		container.Refresh()
	}

	return &Car{Image: carImage}
}

func (c *Car) Park(space int) {
	x, y := calculateParkingPosition(space)

	animateMove(c.Image, x, y)
}

func (c *Car) Exit(exitIndex int) {
	x := 20 * exitIndex
	animateMove(c.Image, float32(x), 200)
	c.Image.Hide()
}

func calculateParkingPosition(space int) (float32, float32) {
	x := float32(130 + 50*space)
	y := float32(130)
	if space > 10 {
		x = float32(90 + 55*(space-10))
		y = float32(300)
	}
	return x, y
}

func animateMove(image *canvas.Image, targetX, targetY float32) {
	startX, startY := image.Position().X, image.Position().Y
	duration := time.Millisecond * 500
	steps := 20
	stepDuration := duration / time.Duration(steps)

	for i := 0; i <= steps; i++ {
		progress := float32(i) / float32(steps)
		newX := startX + progress*(targetX-startX)
		newY := startY + progress*(targetY-startY)

		time.Sleep(stepDuration)
		image.Move(fyne.NewPos(newX, newY))
		image.Refresh()
	}
}
