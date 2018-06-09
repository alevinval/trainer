package trainer

// Enricher interface allows custom logic to enrich activities data.
type Enricher interface {
	Enrich(a *Activity) (err error)
}

// EnrichActivities applies enrichers to a list of activities.
func EnrichActivities(activities ActivityList, enrichers ...Enricher) (err error) {
	for _, activity := range activities {
		for _, enricher := range enrichers {
			err2 := enricher.Enrich(activity)
			if err2 != nil {
				err = err2
			}
		}
	}
	return
}
