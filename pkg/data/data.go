package data

import (
	"database/sql"
	"encoding/json"

	"github.com/rrune/alleytrack/internal/util"
	"github.com/rrune/alleytrack/pkg/models"
	_ "modernc.org/sqlite"
)

type Data struct {
	DB *sql.DB
}

func Init() (d Data, err error) {
	db, err := sql.Open("sqlite", "./data/db.sqlite")
	if util.Check(err) {
		return
	}

	d = Data{
		DB: db,
	}
	return
}

func (d Data) NewDB() (err error) {
	_, err = d.DB.Exec("DROP TABLE IF EXISTS participants;")
	if util.Check(err) {
		return
	}

	_, err = d.DB.Exec("CREATE TABLE participants (number INT not null primary key, name TEXT not null, outoftown BOOLEAN, flinta BOOLEAN, checkpoints TEXT);")
	if util.Check(err) {
		return
	}

	_, err = d.DB.Exec("CREATE UNIQUE INDEX idx_number_unique ON participants(number);")
	if util.Check(err) {
		return
	}
	return
}

func (d Data) NewParticipant(p models.Participant) (err error) {
	// marshal the checkpoint list into json to store them
	checkpoints, err := json.Marshal(p.Checkpoints)
	if util.Check(err) {
		return
	}

	_, err = d.DB.Query("INSERT INTO participants (number, name, outoftown, flinta, checkpoints) VALUES (?, ?, ?, ?, ?);", p.Number, p.Name, p.OutOfTown, p.Flinta, string(checkpoints))
	return
}

func (d Data) GetParicipantFromNumber(number int) (p models.Participant, exists bool, err error) {
	var checkpointsString string

	row := d.DB.QueryRow("SELECT number, name, outoftown, flinta, checkpoints FROM participants WHERE number = ?", number)
	err = row.Scan(&p.Number, &p.Name, &p.OutOfTown, &p.Flinta, &checkpointsString)
	if err == sql.ErrNoRows {
		exists = false
		return
	}
	exists = true
	if util.Check(err) {
		return
	}

	// unmarshal json into usable format
	err = json.Unmarshal([]byte(checkpointsString), &p.Checkpoints)
	return
}

func (d Data) UpdateCheckpoints(p models.Participant) (err error) {
	// marshal the checkpoint list into json to store them
	checkpoints, err := json.Marshal(p.Checkpoints)
	if util.Check(err) {
		return
	}
	_, err = d.DB.Query("UPDATE participants SET checkpoints = ? WHERE number = ?", checkpoints, p.Number)
	return
}

func (d Data) IsNumberTaken(number string) (b bool, err error) {
	if err = d.DB.QueryRow("SELECT number FROM participants WHERE number = ?", number).Scan(&number); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return true, err
	}
	return true, nil
}
