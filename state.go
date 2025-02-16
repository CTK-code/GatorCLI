package main

import (
	"github.com/CTK-code/GatorCLI/internal/config"
	"github.com/CTK-code/GatorCLI/internal/database"
)

type state struct {
	db     *database.Queries
	Config *config.Config
}
