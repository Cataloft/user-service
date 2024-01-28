package apis

import "log/slog"

type Enricher struct {
	logger *slog.Logger
}

func New(logger *slog.Logger) *Enricher {
	return &Enricher{
		logger: logger,
	}
}
