package config

import (
	"path/filepath"

	"github.com/adrg/xdg"
)

var (
	PathAsunder = filepath.Join(xdg.DataHome, "asunder")
	PathDB      = filepath.Join(PathAsunder, "asunder.db")
	PathData    = filepath.Join(PathAsunder, "data.json")
)
