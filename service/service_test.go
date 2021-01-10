package service

import (
	"fmt"
	"github.com/AlexanderYukhanov/90minsGOexp/data/models"
	storage "github.com/AlexanderYukhanov/90minsGOexp/memstorage"
	sm "github.com/AlexanderYukhanov/90minsGOexp/server/models"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const testTrainerID = "00000000-0000-0000-0000-000000000001"

func getTimeOnDay(t *testing.T, kitchenTime string, weekday time.Weekday) time.Time {
	now := time.Now()
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		t.Fatal(err)
	}
	now = now.In(loc).AddDate(0, 0, 1) // to be always in future
	for now.Weekday() != weekday {
		now = now.AddDate(0, 0, 1)
	}
	target, err := time.ParseInLocation(time.Kitchen, kitchenTime, loc)
	if err != nil {
		t.Fatal(err)
	}
	y, m, d := now.Date()
	return target.AddDate(y, int(m) - 1, d - 1)
}

func TestCheckAppointmentParameters(t *testing.T)  {
	s, err := New(storage.NewStorageWithDefaultData())
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct{
		name string
		trainerID models.TrainerID
		startTime time.Time
		endTime time.Time
		code int64
	} {
		{
			name:      "success path",
			trainerID: testTrainerID,
			startTime: getTimeOnDay(t, "10:00AM", time.Monday),
			endTime:   getTimeOnDay(t, "10:00AM", time.Monday).Add(time.Minute * AppointmentLengthMin),
			code:      CodeNone,
		},
		{
			name:      "start time in past",
			trainerID: testTrainerID,
			startTime: getTimeOnDay(t, "10:00AM", time.Monday).AddDate(0, 0,-7),
			endTime:   getTimeOnDay(t, "10:00AM", time.Monday).AddDate(0, 0,-7).Add(time.Minute * AppointmentLengthMin),
			code:      StartTimeInPast,
		},
		{
			name:      "weekend - sunday",
			trainerID: testTrainerID,
			startTime: getTimeOnDay(t, "10:00AM", time.Sunday),
			endTime:   getTimeOnDay(t, "10:00AM", time.Sunday).Add(time.Minute * AppointmentLengthMin),
			code:      OutsideWorkingDays,
		},
		{
			name:      "weekend - saturday",
			trainerID: testTrainerID,
			startTime: getTimeOnDay(t, "10:00AM", time.Saturday),
			endTime:   getTimeOnDay(t, "10:00AM", time.Saturday).Add(time.Minute * AppointmentLengthMin),
			code:      OutsideWorkingDays,
		},
		{
			name:      "outside working hours - too late",
			trainerID: testTrainerID,
			startTime: getTimeOnDay(t, "10:00PM", time.Monday),
			endTime:   getTimeOnDay(t, "10:00PM", time.Monday).Add(time.Minute * AppointmentLengthMin),
			code:      OutsideWorkingHours,
		},
		{
			name:      "outside working hours - too early",
			trainerID: testTrainerID,
			startTime: getTimeOnDay(t, "07:00AM", time.Monday),
			endTime:   getTimeOnDay(t, "07:00AM", time.Monday).Add(time.Minute * AppointmentLengthMin),
			code:      OutsideWorkingHours,
		},
		{
			name:      "opening boundary",
			trainerID: testTrainerID,
			startTime: getTimeOnDay(t, "08:00AM", time.Monday),
			endTime:   getTimeOnDay(t, "08:00AM", time.Monday).Add(time.Minute * AppointmentLengthMin),
			code:      CodeNone,
		},
		{
			name:      "closing boundary",
			trainerID: testTrainerID,
			startTime: getTimeOnDay(t, "04:30PM", time.Monday),
			endTime:   getTimeOnDay(t, "04:30PM", time.Monday).Add(time.Minute * AppointmentLengthMin),
			code:      CodeNone,
		},
		{
			name:      "not aligned",
			trainerID: testTrainerID,
			startTime: getTimeOnDay(t, "10:15AM", time.Monday),
			endTime:   getTimeOnDay(t, "10:15AM", time.Monday).Add(time.Minute * AppointmentLengthMin),
			code:      InvalidAppointmentStartTime,
		},
	}

	for _, cs := range cases {
		err := s.checkAppointmentParameters(&sm.UserAppointment{
			EndsAt:    strfmt.DateTime(cs.endTime),
			StartsAt:  strfmt.DateTime(cs.startTime),
			TrainerID: strfmt.UUID(cs.trainerID),
		})
		assert.True(t, (cs.code == 0) == (err == nil), fmt.Sprintf("%v: %v", cs.name, err))
		if cs.code != 0 {
			assert.Equal(t, codeAsString(cs.code), err.Code, fmt.Sprintf("%v: %v", cs.name, err))
		}
	}
}
