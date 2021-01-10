package models

import "time"

// Appointment models an appointment.
type Appointment struct {
	ID        AppointmentID
	UserID    UserID
	TrainerID TrainerID
	StartTime time.Time
}
