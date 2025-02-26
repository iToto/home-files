package onair

import "on-air/internal/entities"

type SVC interface {
	SetOnAirStatus(ctx context.Context, wl wlog.Logger, entities.OnAirStatus) error
	GetOnAirStatus(ctx context.Context, wl wlog.Logger) (entities.OnAirStatus, error)

}
