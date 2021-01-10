package storage

import (
	"sync"
	"time"

	"github.com/AlexanderYukhanov/90minsGOexp/data/appointments"
	"github.com/AlexanderYukhanov/90minsGOexp/data/models"
	"github.com/AlexanderYukhanov/90minsGOexp/data/trainers"
	sti "github.com/AlexanderYukhanov/90minsGOexp/storage"
	"github.com/google/uuid"
)

var _ sti.Storage = &Storage{}
var _ appointments.Repo = &Storage{}
var _ trainers.Repo = &Storage{}

// Storage is an in-memory Storage impl
type Storage struct {
	mutex sync.Mutex
	appointments []*models.Appointment
	trainers     map[models.TrainerID]*models.Trainer
}

func simpleUUID(i int) string {
	u := uuid.UUID{}
	u[15] = byte(i)
	return u.String()
}

// NewStorageWithDefaultData ...
func NewStorageWithDefaultData() *Storage {
	s := Storage{
		trainers: map[models.TrainerID]*models.Trainer{},
	}
	// some test trainers with predictable uuids
	for i := 1; i < 11; i++ {
		id := models.TrainerID(simpleUUID(i))
		s.trainers[id] = &models.Trainer{
			ID: id,
		}
	}
	return &s
}

// ConnectAppointmentsRepo ...
func (s *Storage) ConnectAppointmentsRepo() (appointments.Repo, error) {
	return s, nil
}

// ConnectTrainersRepo ...
func (s *Storage) ConnectTrainersRepo() (trainers.Repo, error) {
	return s, nil
}

// InsertAppointment ...
func (s *Storage) InsertAppointment(appointment models.Appointment) (models.AppointmentID, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, app := range s.appointments {
		if app.TrainerID == appointment.TrainerID && app.StartTime == appointment.StartTime {
			return "", sti.ErrConflict
		}
	}
	appointment.ID = models.AppointmentID(uuid.New().String())
	s.appointments = append(s.appointments, &appointment)
	return appointment.ID, nil
}

// ListTrainerAppointments ...
func (s *Storage) ListTrainerAppointments(id models.TrainerID, from time.Time, to time.Time) ([]*models.Appointment, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	var result []*models.Appointment
	for _, app := range s.appointments {
		if app.TrainerID == id && !from.After(app.StartTime) && !app.StartTime.After(to) {
			result = append(result, app)
		}
	}
	return result, nil
}

// FindTrainer ...
func (s *Storage) FindTrainer(id models.TrainerID) (*models.Trainer, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	trainer, found := s.trainers[id]
	if !found {
		return nil, sti.ErrNotFound
	}
	return trainer, nil
}
