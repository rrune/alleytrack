package data

import (
	"reflect"
	"testing"

	"github.com/rrune/alleytrack/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestCheckpointsAdd_Table(t *testing.T) {
	tests := []struct {
		name string
		ch   models.Checkpoint
	}{
		{"simple", TestCheckpoints[0]},
		{"text", TestCheckpoints[1]},
		{"third", TestCheckpoints[2]},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := db.Checkpoints.Add(tc.ch)
			assert.Nil(t, err)
		})
	}

	cleanupDB(t)
}

func TestCheckpointsGetAll(t *testing.T) {
	assert.Nil(t, db.Checkpoints.Add(TestCheckpoints[0]))
	assert.Nil(t, db.Checkpoints.Add(TestCheckpoints[1]))

	chs, err := db.Checkpoints.GetAll()
	assert.Nil(t, err)

	assert.Equal(t, len(chs), 2)
	assert.True(t, reflect.DeepEqual(chs[0], TestCheckpoints[0]))
	assert.True(t, reflect.DeepEqual(chs[1], TestCheckpoints[1]))

	cleanupDB(t)
}

func TestCheckpointsGetById_Table(t *testing.T) {
	assert.Nil(t, db.Checkpoints.Add(TestCheckpoints[0]))

	tests := []struct {
		name  string
		id    int
		exist bool
		ch    models.Checkpoint
	}{
		{"simple", TestCheckpoints[0].ID, true, TestCheckpoints[0]},
		{"does not exist", 999, false, models.Checkpoint{}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ch, exists, err := db.Checkpoints.GetById(tc.id)
			assert.Nil(t, err)
			assert.Equal(t, exists, tc.exist)
			assert.True(t, reflect.DeepEqual(ch, tc.ch))
		})
	}

	cleanupDB(t)
}

func TestCheckpointsUpdateById(t *testing.T) {
	assert.Nil(t, db.Checkpoints.Add(TestCheckpoints[0]))

	changedCheckpoint := models.Checkpoint{
		ID:       TestCheckpoints[0].ID,
		Link:     "uud",
		Location: "loc",
		Info:     "infooo",
		Text:     true,
	}

	success, err := db.Checkpoints.UpdateById(TestCheckpoints[0].ID, changedCheckpoint)
	assert.Nil(t, err)
	assert.True(t, success)

	ch, exists, err := db.Checkpoints.GetById(TestCheckpoints[0].ID)
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.True(t, reflect.DeepEqual(ch, changedCheckpoint))

	success, err = db.Checkpoints.UpdateById(999, models.Checkpoint{})
	assert.Nil(t, err)
	assert.False(t, success)

	cleanupDB(t)
}

func TestCheckpointsRemoveById(t *testing.T) {
	assert.Nil(t, db.Checkpoints.Add(TestCheckpoints[0]))

	exists, err := db.Checkpoints.RemoveById(TestCheckpoints[0].ID)
	assert.Nil(t, err)
	assert.True(t, exists)

	exists, err = db.Checkpoints.RemoveById(999)
	assert.Nil(t, err)
	assert.False(t, exists)

	cleanupDB(t)
}
