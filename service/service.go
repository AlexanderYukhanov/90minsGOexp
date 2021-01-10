package service

import (
	"fmt"
	"github.com/AlexanderYukhanov/90minsGOexp/data/appointments"
	"github.com/AlexanderYukhanov/90minsGOexp/data/models"
	"github.com/AlexanderYukhanov/90minsGOexp/data/trainers"
	sm "github.com/AlexanderYukhanov/90minsGOexp/server/models"
	st "github.com/AlexanderYukhanov/90minsGOexp/server/restapi/operations/trainers"
	su "github.com/AlexanderYukhanov/90minsGOexp/server/restapi/operations/users"
	"github.com/AlexanderYukhanov/90minsGOexp/storage"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	CodeNone = iota
	BadTrainerID = iota + 1
	StartTimeInPast
	TimeIntervalTooLarge
	InvalidAppointmentStartTime
	InvalidAppointmentLength
	OutsideWorkingHours
	OutsideWorkingDays
	Holiday

	// TODO: should come from a config
	MaxTimeIntervalDays = 31
	MaxTimeInterval = time.Hour * 24 * MaxTimeIntervalDays
	AppointmentLengthMin = 30
)

var (
	// TODO: should come from a config
	timeZone  = "America/Los_Angeles"
	startTime = parseKitchenTimeOrDie("08:00AM")
	endTime   = parseKitchenTimeOrDie("4:30PM")
	weekends  = []time.Weekday {time.Saturday, time.Sunday}
)

// TODO: split the routes
// Service - all in one routes support implementation
type Service struct {
	appointments appointments.Repo
	trainers     trainers.Repo
	allowedSessionStarts []time.Time
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
	s.populateAllowedSessionsStart()
	return s, nil
}

func (s *Service) populateAllowedSessionsStart() {
	start := startTime
	for !start.After(endTime) {
		s.allowedSessionStarts = append(s.allowedSessionStarts, start)
		start = start.Add(time.Minute * AppointmentLengthMin)
	}
}

// ListAvailableTimesForTrainer ...
func (s* Service) ListAvailableTimesForTrainer(params su.ListAvailableTimesForTrainerParams) middleware.Responder {
	var payload sm.AvailableSlots
	return su.NewListAvailableTimesForTrainerOK().WithPayload(payload)
}

// ListTrainerAppointments ...
func (s *Service) ListTrainerAppointments(params st.ListTrainerAppointmentsParams) middleware.Responder {
	if time.Time(params.EndsAt).Sub(time.Time(params.EndsAt)) > MaxTimeInterval {
		return st.NewListTrainerAppointmentsBadRequest().WithPayload(&sm.Error{
			Code: codeAsString(StartTimeInPast),
			Devmessage: "Start time in the past.",
			Attributes: []string{strconv.FormatInt(MaxTimeIntervalDays, 10)},
		})
	}
	apps, err := s.appointments.ListTrainerAppointments(
		models.TrainerID(params.Trainerid), time.Time(params.StartsAt), time.Time(params.EndsAt))
	if err != nil {
		log.Println(err)
		return middleware.Error(http.StatusInternalServerError, nil)
	}
	var payload sm.TrainerAppointments
	for _, app := range apps {
		payload = append(payload, &sm.TrainerAppointment{
			ID:       strfmt.UUID(app.ID),
			UserID:   strfmt.UUID(app.UserID),
			StartsAt: strfmt.DateTime(app.StartTime),
			EndsAt:   strfmt.DateTime(app.StartTime.Add(30 * time.Minute)),
		})
	}
	return st.NewListTrainerAppointmentsOK().WithPayload(payload)
}

// CreateAppointment ...
func (s *Service) CreateAppointment(params su.CreateAppointmentParams) middleware.Responder {
	if errPayload := s.checkAppointmentParameters(params.Appointment); errPayload != nil {
		return su.NewCreateAppointmentBadRequest().WithPayload(errPayload)
	}
	id, err := s.appointments.InsertAppointment(models.Appointment{
		UserID:    models.UserID(params.Userid),
		TrainerID: models.TrainerID(params.Appointment.TrainerID),
		StartTime: time.Time(params.Appointment.StartsAt),
	})
	if err != nil {
		log.Println(err)
		return middleware.Error(http.StatusInternalServerError, nil)
	}
	appointment := params.Appointment
	appointment.ID = strfmt.UUID(id)
	return su.NewCreateAppointmentCreated().WithPayload(appointment)
}

func (s *Service) checkAppointmentParameters(app *sm.UserAppointment) *sm.Error {
	start := time.Time(app.StartsAt)
	end := time.Time(app.EndsAt)
	y, m, d := start.Date()
	opening := startTime.AddDate(y, int(m) - 1, d - 1)
	last := endTime.AddDate(y, int(m) - 1, d - 1)

	// check if the appointment time is in future
	if time.Now().After(start) {
		return &sm.Error{Code: codeAsString(StartTimeInPast), Devmessage: "Start time in the past."}
	}

	// check for length
	if end.Sub(start) != time.Minute * AppointmentLengthMin {
		return &sm.Error{
			Code: codeAsString(InvalidAppointmentLength),
			Devmessage: "Unsupported appointment length",
			Attributes: []string{strconv.FormatInt(AppointmentLengthMin, 10)},
		}
	}

	// check for weekends
	day := start.Weekday()
	for _, we := range weekends {
		if day == we {
			return &sm.Error{
				Code: codeAsString(OutsideWorkingDays),
				Devmessage: "Appointment is outside of working days",
				Attributes: []string{
					opening.String(),
					last.Add(time.Minute*AppointmentLengthMin).String(),
				},
			}
		}
	}

	// check for working hours
	if start.Before(opening) || start.After(last) {
		return &sm.Error{
			Code: codeAsString(OutsideWorkingHours),
			Devmessage: "Appointment is outside of working hours",
			Attributes: []string{
				opening.String(),
				last.Add(time.Minute*AppointmentLengthMin).String(),
			},
		}
	}

	// check is appointment time is aligned with appointment length
	if !s.isAllowedAppointmentStartAligned(start) {
		return &sm.Error{
			Code:       codeAsString(InvalidAppointmentStartTime),
			Devmessage: "Start time must be aligned on appointment length boundary",
			Attributes: []string{strconv.FormatInt(AppointmentLengthMin, 10)},
		}
	}

	if _, err := s.trainers.FindTrainer(models.TrainerID(app.TrainerID)); err != nil {
		return &sm.Error{
			Code: strconv.FormatInt(BadTrainerID, 16),
			Devmessage: "Bad trainer identifier."}
	}

	return nil
}

func (s *Service) isAllowedAppointmentStartAligned(t time.Time) bool {
	y, m, d := t.Date()
	for _, al := range s.allowedSessionStarts {
		if t.Equal(al.AddDate(y, int(m) - 1, d - 1)) {
			return true
		}
	}
	return false
}

func codeAsString(code int64) string {
	return strconv.FormatInt(code, 16)
}

func parseKitchenTimeOrDie(s string) time.Time {
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		panic(fmt.Sprintf("Failed to load location %v: %v", timeZone, err))
	}
	t, err := time.ParseInLocation(time.Kitchen, s, loc)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse %v as kitchen time: %v", s, err))
	}
	return t
}