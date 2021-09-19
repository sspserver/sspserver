package httpserver

import "geniusrabbit.dev/sspserver/internal/models"

type dummyTargetAccessor struct{}

func (tacc dummyTargetAccessor) TargetByID(id uint64) models.Target { return nil }
