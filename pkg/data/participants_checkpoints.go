package data

import (
	"database/sql"
	"time"

	"github.com/rrune/alleytrack/internal/util"
	"github.com/rrune/alleytrack/pkg/models"
)

type ParticipantsCheckpoints interface {
	Create() (err error)

	Clear() (err error)

	Add(participantNumber int, checkpcheckpointID int, content string) (err error)

	GetAll(participantNumber int) (checkpoints map[int]bool, err error)

	GetCompleted(participantNumber int) (checkpoints []models.Checkpoint, err error)

	Get(participantNumber int, checkpcheckpointID int) (completed bool, content string, err error)

	Remove(participantNumber int, checkpcheckpointID int) (exists bool, err error)
}

type participantsCheckpoints struct {
	DB *sql.DB
}

func (p participantsCheckpoints) Create() (err error) {
	_, err = p.DB.Exec(`
	CREATE TABLE IF NOT EXISTS participant_checkpoints (
    	participant_number INT REFERENCES participants(number) ON DELETE CASCADE,
    	checkpoint_id INT REFERENCES checkpoints(id),
    	completed_at TIMESTAMP DEFAULT current_timestamp,
		content TEXT DEFAULT '',
    	PRIMARY KEY (participant_number, checkpoint_id)
		);
	`)
	return
}

func (p participantsCheckpoints) Clear() (err error) {
	_, err = p.DB.Exec("DELETE FROM participant_checkpoints")
	return
}

func (p participantsCheckpoints) Add(participantNumber int, checkpcheckpointID int, content string) (err error) {
	_, err = p.DB.Exec(`
	INSERT INTO participant_checkpoints (
		participant_number, 
		checkpoint_id,
		content,
		completed_at
	) 
	VALUES (?, ?, ?, ?);`,
		participantNumber,
		checkpcheckpointID,
		content,
		time.Now(),
	)

	return
}

func (p participantsCheckpoints) GetAll(participantNumber int) (checkpoints map[int]bool, err error) {
	rows, err := p.DB.Query(`
	SELECT 
  		c.id AS checkpoint_id,
  	CASE WHEN pc.checkpoint_id IS NOT NULL THEN 1 ELSE 0 END AS completed
	FROM checkpoints AS c
	LEFT JOIN participant_checkpoints AS pc
  		ON c.id = pc.checkpoint_id
  		AND pc.participant_number = ?;
	`,
		participantNumber,
	)

	checkpoints = map[int]bool{}

	for rows.Next() {
		if util.Check(rows.Err()) {
			err = rows.Err()
			return
		}
		var checkpoint int
		var completed bool

		err = rows.Scan(&checkpoint, &completed)
		if util.Check(err) {
			return
		}
		checkpoints[checkpoint] = completed
	}

	return
}

func (p participantsCheckpoints) GetCompleted(participantNumber int) (checkpoints []models.Checkpoint, err error) {
	rows, err := p.DB.Query(`
	SELECT c.*, pc.content
	FROM checkpoints AS c
	JOIN participant_checkpoints AS pc
  	ON c.id = pc.checkpoint_id
	WHERE pc.participant_number = ?;
	`,
		participantNumber,
	)
	if util.Check(err) {
		return
	}

	for rows.Next() {
		if util.Check(rows.Err()) {
			return nil, rows.Err()
		}
		ch := models.Checkpoint{}
		err = rows.Scan(&ch.ID, &ch.Link, &ch.Location, &ch.Info, &ch.Text, &ch.Content)
		if util.Check(err) {
			return
		}
		checkpoints = append(checkpoints, ch)
	}

	return
}

func (p participantsCheckpoints) Get(participantNumber int, checkpcheckpointID int) (completed bool, content string, err error) {
	rows, err := p.DB.Query("SELECT content FROM participant_checkpoints WHERE (participant_number = ? AND checkpoint_id = ?)", participantNumber, checkpcheckpointID)
	for rows.Next() {
		err = rows.Scan(&content)
		return true, content, err
	}

	return false, content, err
}

func (p participantsCheckpoints) Remove(participantNumber int, checkpcheckpointID int) (exists bool, err error) {
	res, err := p.DB.Exec("DELETE FROM participant_checkpoints WHERE (participant_number = ? AND checkpoint_id = ?)", participantNumber, checkpcheckpointID)
	if util.Check(err) {
		return
	}
	num, err := res.RowsAffected()
	exists = num > 0

	return
}
