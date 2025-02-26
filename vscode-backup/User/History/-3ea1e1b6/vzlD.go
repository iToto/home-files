package onair

import "on-air/internal/entities"

type SVC interface {
	SetOnAirStatus(ctx context.Context, wl wlog.Logger, entities.OnAirStatus) error
	GetOnAirStatus(ctx context.Context, wl wlog.Logger) (entities.OnAirStatus, error)
	ToggleOnAirStatus(ctx context.Context, wl wlog.Logger) error
}

type onAirService struct {
	// add any dependencies here (DB, Client, etc.)
}

func New(
// pass in any dependencies
) (SVC, error) {
	return &onAirService{}, nil
}
