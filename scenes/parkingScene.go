package scenes

import (
	"fmt"
	"math/rand"
	"parking_lot_simulator/models"
	"parking_lot_simulator/views"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func StartArrival(numVehicles int) {
	mainApp := app.New()
	rand.Seed(time.Now().UnixNano())

	park := models.NewParkingLot(20)
	window, container := views.SetupWindow(mainApp) // Llama a SetupWindow desde el paquete views

	done := make(chan struct{})

	// Iniciar la simulación de vehículos llegando
	go runParkingSimulation(park, numVehicles, container)

	// Procesar vehículos que entran y salen usando los canales de notificación
	go func() {
		processedVehicles := 0 // Contador para los vehículos que han salido

		for {
			select {
			case vehicle := <-park.ParkedChannel():
				// Vehículo estacionado (solo log, no se requiere procesamiento adicional)
				fmt.Printf("Vehicle %d parked.\n", vehicle.ID)

			case <-park.ExitedChannel(): // Usamos `_` para ignorar `vehicle`
				processedVehicles++
				if processedVehicles == numVehicles {
					done <- struct{}{}
				}
			case <-done:
				window.Close()
				return
			}
		}
	}()

	window.ShowAndRun()
}

func runParkingSimulation(park *models.ParkingLot, numVehicles int, container *fyne.Container) {
	for i := 1; i <= numVehicles; i++ {
		vehicle := models.Vehicle{ID: i}
		go park.VehicleArrives(vehicle, container)
		time.Sleep(time.Duration(rand.ExpFloat64() * float64(time.Second)))
	}
}
