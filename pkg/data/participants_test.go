package data

import (
	"reflect"
	"testing"

	"github.com/rrune/alleytrack/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestParticipantsAdd_Table(t *testing.T) {
	tests := []struct {
		name string
		pt   models.Participant
	}{
		{"simple", TestParticipants[0]},
		{"OoT", TestParticipants[1]},
		{"flinta", TestParticipants[2]},
		{"both", TestParticipants[3]},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := db.Participants.Add(tc.pt)
			assert.Nil(t, err)
		})
	}

	cleanupDB(t)
}

func TestParticipantsGetAll(t *testing.T) {
	assert.Nil(t, db.Participants.Add(TestParticipants[0]))
	assert.Nil(t, db.Participants.Add(TestParticipants[1]))

	pts, err := db.Participants.GetAll()
	assert.Nil(t, err)

	assert.Equal(t, len(pts), 2)
	assert.True(t, reflect.DeepEqual(pts[0], TestParticipants[0]))
	assert.True(t, reflect.DeepEqual(pts[1], TestParticipants[1]))

	cleanupDB(t)
}

func TestParticipantsGetByNumber_Table(t *testing.T) {
	assert.Nil(t, db.Participants.Add(TestParticipants[0]))

	tests := []struct {
		name   string
		number int
		exist  bool
		pt     models.Participant
	}{
		{"simple", TestParticipants[0].Number, true, TestParticipants[0]},
		{"does not exist", 999, false, models.Participant{}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pt, exists, err := db.Participants.GetByNumber(tc.number)
			assert.Nil(t, err)
			assert.Equal(t, exists, tc.exist)
			assert.True(t, reflect.DeepEqual(pt, tc.pt))
		})
	}

	cleanupDB(t)
}

func TestParticipantsUpdateByNumber(t *testing.T) {
	assert.Nil(t, db.Participants.Add(TestParticipants[0]))

	changedParticipant := models.Participant{
		Number:    TestParticipants[0].Number,
		Name:      "test",
		OutOfTown: true,
		Flinta:    true,
	}

	success, err := db.Participants.UpdateByNumber(TestParticipants[0].Number, changedParticipant)
	assert.Nil(t, err)
	assert.True(t, success)

	pt, exists, err := db.Participants.GetByNumber(TestParticipants[0].Number)
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.True(t, reflect.DeepEqual(pt, changedParticipant))

	success, err = db.Participants.UpdateByNumber(999, models.Participant{})
	assert.Nil(t, err)
	assert.False(t, success)

	cleanupDB(t)
}

func TestParticipantsRemoveByNumber(t *testing.T) {
	assert.Nil(t, db.Participants.Add(TestParticipants[0]))

	exists, err := db.Participants.RemoveByNumber(TestParticipants[0].Number)
	assert.Nil(t, err)
	assert.True(t, exists)

	exists, err = db.Participants.RemoveByNumber(999)
	assert.Nil(t, err)
	assert.False(t, exists)

	cleanupDB(t)
}
