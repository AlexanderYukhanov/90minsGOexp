package appointments

import (
	"time"

	"github.com/AlexanderYukhanov/90minsGOexp/data/models"
)

// Repo - appointments repository
type Repo interface {
	InsertAppointment(appointment models.Appointment) (models.AppointmentID, error)
	ListTrainerAppointments(trainer models.TrainerID, from time.Time, to time.Time) ([]*models.Appointment, error)
}
