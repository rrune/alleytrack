package config

import (
	"os"

	"github.com/rrune/alleytrack/internal/util"
	"github.com/rrune/alleytrack/pkg/data"
	"github.com/rrune/alleytrack/pkg/models"
	"gopkg.in/yaml.v2"
)

func ReadManifestIntoDatabase(path string, db *data.Data) (err error) {
	chs := []models.Checkpoint{}

	manifest, err := os.ReadFile(path)
	if util.Check(err) {
		return
	}
	err = yaml.Unmarshal(manifest, &chs)
	if util.Check(err) {
		return
	}

	err = db.Checkpoints.Clear()
	if util.Check(err) {
		return
	}
	err = db.CheckpointDependencies.Clear()
	if util.Check(err) {
		return
	}

	for _, ch := range chs {
		err = db.Checkpoints.Add(ch)
		if util.Check(err) {
			return
		}

		for _, d := range ch.Requirements {
			err = db.CheckpointDependencies.Add(ch.ID, d)
			if util.Check(err) {
				return
			}
		}

	}

	return
}
