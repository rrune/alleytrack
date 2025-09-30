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

	// init the table if it doesn't exist
	exists, err := d.tableExists("participants")
	if util.Check(err) {
		return
	}
	if !exists {
		d.NewDB()
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

func (d Data) GetAllParticipants() (ps []models.Participant, err error) {
	rows, err := d.DB.Query("SELECT number, name, outoftown, flinta, checkpoints FROM participants")

	for rows.Next() {
		if util.Check(rows.Err()) {
			return nil, rows.Err()
		}
		p := models.Participant{}
		var checkpointsString string
		err = rows.Scan(&p.Number, &p.Name, &p.OutOfTown, &p.Flinta, &checkpointsString)
		if util.Check(err) {
			return
		}
		// unmarshal json into usable format
		err = json.Unmarshal([]byte(checkpointsString), &p.Checkpoints)
		if util.Check(err) {
			return
		}
		ps = append(ps, p)
	}
	return
}

func (d Data) UpdateCheckpoints(p models.Participant) (err error) {
	// marshal the checkpoint list into json to store them
	checkpoints, err := json.Marshal(p.Checkpoints)
	if util.Check(err) {
		return
	}
	_, err = d.DB.Exec("UPDATE participants SET checkpoints = ? WHERE number = ?", checkpoints, p.Number)
	return
}

func (d Data) UpdateParticipant(p models.Participant) (err error) {
	_, err = d.DB.Exec("UPDATE participants SET name = ?, outoftown = ?, flinta = ? WHERE number = ?", p.Name, p.OutOfTown, p.Flinta, p.Number)
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

func (d Data) RemoveParticipantByNumber(number string) (exists bool, err error) {
	res, err := d.DB.Exec("DELETE FROM participants WHERE number = ?", number)
	if util.Check(err) {
		return
	}
	num, err := res.RowsAffected()
	exists = num > 0

	return
}

/* func (d Data) RemoveCheckpointFromParticipant(number string, checkpoint string) (success bool, err error) {
	num, err := strconv.Atoi(number)
	if util.Check(err) {
		return
	}
	ch, err := strconv.Atoi(checkpoint)
	if util.Check(err) {
		return
	}

	p, exists, err := d.GetParicipantFromNumber(num)
	if !exists {
		return false, err
	}

	delete(p.Checkpoints, ch)
	checkpointsString, err := json.Marshal(p.Checkpoints)
	if util.Check(err) {
		return
	}

	_, err = d.DB.Exec("UPDATE participants SET checkpoints = ? WHERE number = ?", checkpointsString, number)
	if util.Check(err) {
		return
	}
	return true, err
} */

func (d Data) tableExists(tableName string) (bool, error) {
	var count int
	err := d.DB.QueryRow("SELECT COUNT (name) FROM sqlite_master WHERE type = 'table' AND name =? ;", tableName).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
