package onair

import "on-air/internal/entities"

type SVC interface {
	SetOnAirStatus(ctx context.Context, wl wlog.Logger, entities.OnAirStatus) error
	GetOnAirStatus(ctx context.Context, wl wlog.Logger) (entities.OnAirStatus, error)
	ToggleOnAirStatus(ctx context.Context, wl wlog.Logger) error
}

type onAirService struct {
	// add any dependencies here (DB, Client, etc.)
	db          *sqlx.DB
}

func New(
// pass in any dependencies
	sor *sqlx.DB,
) (SVC, error) {
	return &onAirService{db:sor}, nil
}
