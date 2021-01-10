package service

import (
	"github.com/AlexanderYukhanov/90minsGOexp/data/appointments"
	"github.com/AlexanderYukhanov/90minsGOexp/data/trainers"
	sm "github.com/AlexanderYukhanov/90minsGOexp/server/models"
	st "github.com/AlexanderYukhanov/90minsGOexp/server/restapi/operations/trainers"
	su "github.com/AlexanderYukhanov/90minsGOexp/server/restapi/operations/users"
	"github.com/AlexanderYukhanov/90minsGOexp/storage"
	"github.com/go-openapi/runtime/middleware"
)

// TODO: split the routes
// Service - all in one routes support implementation
type Service struct {
	appointments appointments.Repo
	trainers     trainers.Repo
}

// New ...
func New(storage storage.Storage) (*Service, error) {
	s := &Service{}
	var err error
	if s.appointments, err = storage.ConnectAppointmentsRepo(); err != nil {
		return nil, err
	}
	if s.trainers, err = storage.ConnectTrainersRepo(); err != nil {
		return nil, err
	}
	return s, nil
}

// ListAvailableTimesForTrainer ...
func (s* Service) ListAvailableTimesForTrainer(params su.ListAvailableTimesForTrainerParams) middleware.Responder {
	var payload sm.AvailableSlots
	return su.NewListAvailableTimesForTrainerOK().WithPayload(payload)
}

// ListTrainerAppointments ...
func (s *Service) ListTrainerAppointments(params st.ListTrainerAppointmentsParams) middleware.Responder {
	var payload sm.TrainerAppointments
	return st.NewListTrainerAppointmentsOK().WithPayload(payload)
}

// CreateAppointment ...
func (s *Service) CreateAppointment(params su.CreateAppointmentParams) middleware.Responder {
	return su.NewCreateAppointmentCreated().WithPayload(params.Appointment)
}