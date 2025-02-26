// Package tradelogger will handle logging trades to BQ
package DataLogger

// this must match the schema in BQ
type record struct {

}
// Archive writes the given event to BigQuery.
func (a *Archiver) Archive(ctx context.Context, wl wlog.Logger, t *issuingeventssvc.Transaction) error {
	rng := rand.New(rand.NewSource(time.Now().UnixNano())) // nolint: gosec
	id, err := ulid.New(ulid.Now(), rng)
	if err != nil {
		return fmt.Errorf("error generating id: %s", err)
	}

	r := &record{
	}

	return a.inserter.Put(ctx, r)
