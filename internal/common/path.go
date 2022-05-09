package common

import (
	"path/filepath"

	"github.com/adrg/xdg"
)

var (
	// PathAsunder is the directory where asunder stores its data.
	PathAsunder = filepath.Join(xdg.DataHome, "asunder")

	// PathDB is the database where encrypted entries are stored.
	PathDB = filepath.Join(PathAsunder, "asunder.db")

	// PathData is where the hashed secret is stored.
	PathData = filepath.Join(PathAsunder, "data.json")
)
