package storage

import (
	"github.com/AlexanderYukhanov/90minsGOexp/data/appointments"
	"github.com/AlexanderYukhanov/90minsGOexp/data/trainers"
)

// Storage - storage interface.
type Storage interface {
	ConnectAppointmentsRepo() (appointments.Repo, error)
	ConnectTrainersRepo() (trainers.Repo, error)
}
