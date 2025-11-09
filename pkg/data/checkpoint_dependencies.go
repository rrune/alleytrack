package data

import (
	"database/sql"

	"github.com/rrune/alleytrack/internal/util"
	"github.com/rrune/alleytrack/pkg/models"
)

type CheckpointDependencies interface {
	Create() (err error)

	Clear() (err error)

	Add(chID int, chDependID int) (err error)

	GetDependencies(chID int) (depends []int, err error)

	GetAvailableByNumber(number int) (checkpoints []models.Checkpoint, err error)

	Remove(chID int, chDependID int) (exists bool, err error)
}

type checkpointDependencies struct {
	DB *sql.DB
}

func (c checkpointDependencies) Create() (err error) {
	_, err = c.DB.Exec(`
	CREATE TABLE checkpoint_dependencies (
    	checkpoint_id INTEGER NOT NULL REFERENCES checkpoints(id) ON DELETE CASCADE,
    	required_checkpoint_id INTEGER NOT NULL REFERENCES checkpoints(id) ON DELETE CASCADE,
    	PRIMARY KEY (checkpoint_id, required_checkpoint_id),
    	CONSTRAINT no_self_dependency CHECK (checkpoint_id <> required_checkpoint_id)
	);
	`)
	return
}

func (c checkpointDependencies) Clear() (err error) {
	_, err = c.DB.Exec("DELETE FROM checkpoint_dependencies")
	return
}

func (c checkpointDependencies) Add(chID int, chDependID int) (err error) {
	_, err = c.DB.Exec(`
	INSERT INTO checkpoint_dependencies (
		checkpoint_id, 
		required_checkpoint_id
	) 
	VALUES (?, ?);`,
		chID,
		chDependID,
	)

	return
}

func (c checkpointDependencies) GetDependencies(chID int) (depends []int, err error) {
	rows, err := c.DB.Query("SELECT required_checkpoint_id FROM checkpoint_dependencies WHERE checkpoint_id = ?", chID)
	if util.Check(err) {
		return
	}

	depends = []int{}

	for rows.Next() {
		if util.Check(rows.Err()) {
			err = rows.Err()
			return
		}
		var checkpoint int

		err = rows.Scan(&checkpoint)
		if util.Check(err) {
			return
		}
		depends = append(depends, checkpoint)
	}

	return
}

func (c checkpointDependencies) GetAvailableByNumber(number int) (checkpoints []models.Checkpoint, err error) {
	rows, err := c.DB.Query(`
	SELECT c.id, c.link, c.location, c.info, c.text
	FROM checkpoints c
	LEFT JOIN checkpoint_dependencies d 
  		ON c.id = d.checkpoint_id
	LEFT JOIN participant_checkpoints pc 
  		ON pc.checkpoint_id = d.required_checkpoint_id
  	AND pc.participant_number = ?
	WHERE c.id NOT IN (
    	SELECT checkpoint_id
    	FROM participant_checkpoints
    	WHERE participant_number = ?
	)
	GROUP BY c.id, c.link, c.location, c.info, c.text
	HAVING COUNT(d.required_checkpoint_id) = COUNT(pc.checkpoint_id);
	`,
		number,
		number,
	)
	if util.Check(err) {
		return
	}

	for rows.Next() {
		if util.Check(rows.Err()) {
			err = rows.Err()
			return
		}
		ch := models.Checkpoint{}
		err = rows.Scan(&ch.ID, &ch.Link, &ch.Location, &ch.Info, &ch.Text)
		if util.Check(err) {
			return
		}
		checkpoints = append(checkpoints, ch)
	}

	return
}

func (c checkpointDependencies) Remove(chID int, chDependID int) (exists bool, err error) {
	res, err := c.DB.Exec("DELETE FROM checkpoint_dependencies WHERE (checkpoint_id = ? AND required_checkpoint_id = ?)", chID, chDependID)
	if util.Check(err) {
		return
	}
	num, err := res.RowsAffected()
	exists = num > 0

	return
}
