package data

import (
	"os"
	"testing"
	"time"

	"github.com/rrune/alleytrack/pkg/models"
)

var db Data
var TestParticipants []models.Participant
var TestCheckpoints []models.Checkpoint

func TestMain(m *testing.M) {
	// create temp file
	_, err := os.Create("./test.sqlite")
	if err != nil {
		panic(err)
	}

	// init the Database struct
	db, err = Init("./test.sqlite")
	if err != nil {
		panic(err)
	}

	// create tables
	err = db.Participants.Create()
	if err != nil {
		panic(err)
	}
	err = db.Checkpoints.Create()
	if err != nil {
		panic(err)
	}
	err = db.ParticipantsCheckpoints.Create()
	if err != nil {
		panic(err)
	}
	err = db.CheckpointDependencies.Create()
	if err != nil {
		panic(err)
	}

	// create test participants
	TestParticipants = []models.Participant{
		{141, "rune", false, false},
		{69, "oot", true, false},
		{8, "flint", false, true},
		{68, "both", true, true},
	}

	TestCheckpoints = []models.Checkpoint{
		{1, "abc", "downtown", "info downtown", false, nil, "", time.Time{}},
		{2, "xyz", "uptown", "info uptown", true, nil, "", time.Time{}},
		{3, "pqd", "center", "info center", false, nil, "", time.Time{}},
	}

	// run the tests
	code := m.Run()

	// remove the temp file
	os.Remove("./test.sqlite")
	os.Exit(code)
}

func cleanupDB(t *testing.T) {
	_, err := db.DB.Exec("DELETE FROM participants")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.DB.Exec("DELETE FROM checkpoints")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.DB.Exec("DELETE FROM participant_checkpoints")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.DB.Exec("DELETE FROM checkpoint_dependencies")
	if err != nil {
		t.Fatal(err)
	}
}
