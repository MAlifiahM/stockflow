package xlogger

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/MAlifiahM/stockflow/backend/internal/config"
)

func TestSetup_Development(t *testing.T) {
	cfg := config.Config{IsDevelopment: true}
	Setup(cfg)
	assert.NotNil(t, Logger)
}

func TestSetup_Production(t *testing.T) {
	cfg := config.Config{}
	Setup(cfg)
	assert.NotNil(t, Logger)
}
