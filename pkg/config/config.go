package config

import (
	"os"

	"github.com/rrune/alleytrack/internal/util"
	"github.com/rrune/alleytrack/pkg/models"
	"gopkg.in/yaml.v2"
)

func SwitchEnabledInConfig(alleycat *models.Alleycat) (err error) {
	alleycat.Config.Enabled = !alleycat.Config.Enabled

	confByte, err := yaml.Marshal(alleycat.Config)
	if util.Check(err) {
		return
	}
	err = os.WriteFile("./config/config.yml", confByte, 0644)
	return
}
