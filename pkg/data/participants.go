package data

import (
	"database/sql"

	"github.com/rrune/alleytrack/internal/util"
	"github.com/rrune/alleytrack/pkg/models"
)

type Participants interface {
	Create() (err error)

	Clear() (err error)

	Add(pt models.Participant) (err error)

	GetAll() (pt []models.Participant, err error)

	GetByNumber(number int) (pt models.Participant, exists bool, err error)

	UpdateByNumber(number int, pt models.Participant) (success bool, err error)

	RemoveByNumber(number int) (exists bool, err error)
}

type participants struct {
	DB *sql.DB
}

func (p participants) Create() (err error) {
	_, err = p.DB.Exec(`
	CREATE TABLE IF NOT EXISTS participants (
		number INT not null primary key, 
		name TEXT not null, 
		outoftown BOOLEAN, 
		flinta BOOLEAN
		);
		`)
	if util.Check(err) {
		return
	}

	_, err = p.DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_number_unique ON participants(number);")
	return
}

func (p participants) Clear() (err error) {
	_, err = p.DB.Exec("DELETE FROM participants")
	return
}

func (p participants) Add(pt models.Participant) (err error) {
	_, err = p.DB.Exec(`
	INSERT INTO participants (
		number, 
		name, 
		outoftown, 
		flinta
	) 
	VALUES (?, ?, ?, ?);`,
		pt.Number,
		pt.Name,
		pt.OutOfTown,
		pt.Flinta,
	)

	return
}

func (p participants) GetAll() (pts []models.Participant, err error) {
	rows, err := p.DB.Query("SELECT number, name, outoftown, flinta FROM participants")

	for rows.Next() {
		if util.Check(rows.Err()) {
			return nil, rows.Err()
		}
		pt := models.Participant{}
		err = rows.Scan(&pt.Number, &pt.Name, &pt.OutOfTown, &pt.Flinta)
		if util.Check(err) {
			return
		}
		pts = append(pts, pt)
	}
	return
}

func (p participants) GetByNumber(number int) (pt models.Participant, exists bool, err error) {
	row := p.DB.QueryRow("SELECT number, name, outoftown, flinta FROM participants WHERE number = ?", number)
	err = row.Scan(&pt.Number, &pt.Name, &pt.OutOfTown, &pt.Flinta)
	if err == sql.ErrNoRows {
		err = nil
		exists = false
		return
	}
	exists = true
	if util.Check(err) {
		return
	}
	return
}

func (p participants) UpdateByNumber(number int, pt models.Participant) (success bool, err error) {
	res, err := p.DB.Exec(`
	UPDATE participants SET 
		number = ?, 
		name = ?, 
		outoftown = ?, 
		flinta = ? 
	WHERE number = ?`,
		pt.Number,
		pt.Name,
		pt.OutOfTown,
		pt.Flinta,
		number,
	)
	if util.Check(err) {
		return
	}
	affected, err := res.RowsAffected()
	return affected > 0, err
}

func (p participants) RemoveByNumber(number int) (exists bool, err error) {
	res, err := p.DB.Exec("DELETE FROM participants WHERE number = ?", number)
	if util.Check(err) {
		return
	}
	num, err := res.RowsAffected()
	exists = num > 0

	return
}
