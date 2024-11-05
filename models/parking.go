package models

import (
	"fmt"
	"math/rand"
	"parking_lot_simulator/characters"
	"time"

	"fyne.io/fyne/v2"
)

type Vehicle struct {
	ID    int
	Space int             // Espacio asignado a este vehículo
	Car   *characters.Car // Referencia a la instancia de Car
}

type ParkingLot struct {
	Spaces          []bool
	availableSpaces chan struct{} // Canal para espacios disponibles
	parkedChannel   chan Vehicle  // Notificar cuando un auto se estaciona
	exitedChannel   chan Vehicle  // Notificar cuando un auto sale
	exitCount       int           // Contador de salidas
}

func NewParkingLot(numSpaces int) *ParkingLot {
	p := &ParkingLot{
		Spaces:          make([]bool, numSpaces),
		availableSpaces: make(chan struct{}, numSpaces), // Limita a la cantidad de espacios
		parkedChannel:   make(chan Vehicle),             // Canal de notificación para parqueo
		exitedChannel:   make(chan Vehicle),             // Canal de notificación para salida
		exitCount:       0,
	}

	// Inicializar el canal de espacios disponibles con tantos espacios como capacidad
	for i := 0; i < numSpaces; i++ {
		p.availableSpaces <- struct{}{}
	}
	return p
}

// VehicleArrives maneja la llegada y el estacionamiento de un vehículo.
func (p *ParkingLot) VehicleArrives(v Vehicle, container *fyne.Container) {
	fmt.Printf("Car %d arrives.\n", v.ID)

	// Esperar a que haya un espacio disponible (bloquea si está lleno)
	<-p.availableSpaces

	// Crear la instancia gráfica del auto y asignarla al campo Car de Vehicle
	v.Car = characters.NewCar(container, false)
	for i, occupied := range p.Spaces {
		if !occupied {
			p.Spaces[i] = true
			v.Space = i // Asignar el espacio a v.Space
			fmt.Printf("Car %d parked in space %d.\n", v.ID, i)
			v.Car.Park(i)        // Llamamos a Park en la instancia de Car
			p.parkedChannel <- v // Notificar que se ha estacionado
			break
		}
	}

	// Simula tiempo de estacionamiento antes de que el vehículo salga
	time.Sleep(time.Second * time.Duration(rand.Intn(3)+3))

	// Llamar a la función de salida para el vehículo después del tiempo de estacionamiento
	p.VehicleExits(v)
}

// VehicleExits maneja la salida de un vehículo y libera el espacio específico
func (p *ParkingLot) VehicleExits(v Vehicle) {
	// Liberar el espacio específico asignado a este vehículo
	if v.Space >= 0 && v.Space < len(p.Spaces) && p.Spaces[v.Space] {
		p.Spaces[v.Space] = false
		fmt.Printf("Car %d freed space %d.\n", v.ID, v.Space)
	} else {
		fmt.Printf("Warning: Car %d attempting to free an invalid or already freed space %d.\n", v.ID, v.Space)
	}

	fmt.Printf("Car %d heads to the exit\n", v.ID)
	p.exitCount++
	v.Car.Exit(p.exitCount) // Pasar el índice `exitCount` a Exit para mover el vehículo y esconder la imagen
	time.Sleep(time.Second * time.Duration(rand.Intn(2)))
	fmt.Printf("Car %d exits the parking lot\n", v.ID)

	// Notificar que el vehículo ha salido
	p.exitedChannel <- v

	// Liberar el espacio en el estacionamiento
	p.availableSpaces <- struct{}{} // Libera un espacio en el canal de disponibilidad
}

// Getter para los canales de notificación
func (p *ParkingLot) ParkedChannel() <-chan Vehicle {
	return p.parkedChannel
}

func (p *ParkingLot) ExitedChannel() <-chan Vehicle {
	return p.exitedChannel
}
