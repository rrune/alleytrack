package data

import (
	"database/sql"

	"github.com/rrune/alleytrack/internal/util"
	_ "modernc.org/sqlite"
)

type Data struct {
	Participants            Participants
	Checkpoints             Checkpoints
	ParticipantsCheckpoints ParticipantsCheckpoints
	CheckpointDependencies  CheckpointDependencies
	DB                      *sql.DB
}

func Init(path string) (d Data, err error) {
	db, err := sql.Open("sqlite", path)
	if util.Check(err) {
		return
	}
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if util.Check(err) {
		return
	}

	d = Data{
		Participants:            participants{DB: db},
		Checkpoints:             checkpoints{DB: db},
		ParticipantsCheckpoints: participantsCheckpoints{DB: db},
		CheckpointDependencies:  checkpointDependencies{DB: db},
		DB:                      db,
	}

	d.Participants.Create()
	if util.Check(err) {
		return
	}
	d.Checkpoints.Create()
	if util.Check(err) {
		return
	}
	d.ParticipantsCheckpoints.Create()
	if util.Check(err) {
		return
	}
	d.CheckpointDependencies.Create()

	return
}
