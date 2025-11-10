package data

import (
	"database/sql"

	"github.com/rrune/alleytrack/internal/util"
	"github.com/rrune/alleytrack/pkg/models"
)

type Checkpoints interface {
	Create() (err error)

	Clear() (err error)

	Add(ch models.Checkpoint) (err error)

	GetAll() (ch []models.Checkpoint, err error)

	GetById(id int) (ch models.Checkpoint, exists bool, err error)

	GetByLink(link string) (ch models.Checkpoint, exists bool, err error)

	UpdateById(id int, ch models.Checkpoint) (success bool, err error)

	RemoveById(id int) (exists bool, err error)
}

type checkpoints struct {
	DB *sql.DB
}

func (c checkpoints) Create() (err error) {
	_, err = c.DB.Exec(`
	CREATE TABLE IF NOT EXISTS checkpoints (
		id INTEGER not null primary key, 
		link TEXT not null, 
		location TEXT not null, 
		info TEXT not null, 
		text BOOLEAN
		);
		`)
	if util.Check(err) {
		return
	}

	_, err = c.DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_number_unique ON checkpoints(id);")
	return
}

func (c checkpoints) Clear() (err error) {
	_, err = c.DB.Exec("DELETE FROM checkpoints")
	return
}

func (c checkpoints) Add(ch models.Checkpoint) (err error) {
	_, err = c.DB.Exec(`
	INSERT INTO checkpoints (
		id, 
		link, 
		location, 
		info,
		text
	) 
	VALUES (?, ?, ?, ?, ?);`,
		ch.ID,
		ch.Link,
		ch.Location,
		ch.Info,
		ch.Text,
	)

	return
}

func (c checkpoints) GetAll() (chs []models.Checkpoint, err error) {
	rows, err := c.DB.Query("SELECT id, link, location, info, text FROM checkpoints")

	for rows.Next() {
		if util.Check(rows.Err()) {
			return nil, rows.Err()
		}
		ch := models.Checkpoint{}
		err = rows.Scan(&ch.ID, &ch.Link, &ch.Location, &ch.Info, &ch.Text)
		if util.Check(err) {
			return
		}
		chs = append(chs, ch)
	}
	return
}

func (c checkpoints) GetById(id int) (ch models.Checkpoint, exists bool, err error) {
	row := c.DB.QueryRow("SELECT id, link, location, info, text FROM checkpoints WHERE id = ?", id)
	err = row.Scan(&ch.ID, &ch.Link, &ch.Location, &ch.Info, &ch.Text)
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

func (c checkpoints) GetByLink(link string) (ch models.Checkpoint, exists bool, err error) {
	row := c.DB.QueryRow("SELECT id, link, location, info, text FROM checkpoints WHERE link = ?", link)
	err = row.Scan(&ch.ID, &ch.Link, &ch.Location, &ch.Info, &ch.Text)
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

func (c checkpoints) UpdateById(id int, ch models.Checkpoint) (success bool, err error) {
	res, err := c.DB.Exec(`
	UPDATE checkpoints SET 
		link = ?, 
		location = ?, 
		info = ?,
		text = ? 
	WHERE id = ?`,
		ch.Link,
		ch.Location,
		ch.Info,
		ch.Text,
		id,
	)
	if util.Check(err) {
		return
	}
	affected, err := res.RowsAffected()
	return affected > 0, err
}

func (c checkpoints) RemoveById(id int) (exists bool, err error) {
	res, err := c.DB.Exec("DELETE FROM checkpoints WHERE id = ?", id)
	if util.Check(err) {
		return
	}
	num, err := res.RowsAffected()
	exists = num > 0

	return
}
