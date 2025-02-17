package telemetry

import "github.com/DataDog/datadog-agent/pkg/telemetry"

const (
	// QueryEmptyEntityID refers to a query made with an empty entity id
	QueryEmptyEntityID = "empty_entity_id"
	// QueryEmptyTags refers to a query that returned no tags
	QueryEmptyTags = "empty_tags"
	// QuerySuccess refers to a successful query
	QuerySuccess = "success"

	// FetchNotFound refers to a tagger fetch that did not find an entity
	FetchNotFound = "not_found"
	// FetchError refers to a tagger fetch that returned an error
	FetchError = "error"
	// FetchSuccess refers to a tagger fetch that was successful
	FetchSuccess = "success"
)

var (
	// StoredEntities tracks how many entities are stored in the tagger.
	StoredEntities = telemetry.NewGaugeWithOpts("tagger", "stored_entities",
		[]string{"source", "prefix"}, "Number of entities in the store.",
		telemetry.Options{NoDoubleUnderscoreSep: true})

	// UpdatedEntities tracks the number of updates to tagger entities.
	UpdatedEntities = telemetry.NewCounterWithOpts("tagger", "updated_entities",
		[]string{}, "Number of updates made to entities.",
		telemetry.Options{NoDoubleUnderscoreSep: true})

	// PrunedEntities tracks the number of pruned tagger entities.
	PrunedEntities = telemetry.NewGaugeWithOpts("tagger", "pruned_entities",
		[]string{}, "Number of pruned tagger entities.",
		telemetry.Options{NoDoubleUnderscoreSep: true})

	// Queries tracks the number of queries made against the tagger.
	Queries = telemetry.NewCounterWithOpts("tagger", "queries",
		[]string{"cardinality", "status"}, "Queries made against the tagger.",
		telemetry.Options{NoDoubleUnderscoreSep: true})

	// Fetches tracks the number of fetches from the underlying collectors.
	Fetches = telemetry.NewCounterWithOpts("tagger", "fetches",
		[]string{"collector", "status"}, "Fetches from collectors.",
		telemetry.Options{NoDoubleUnderscoreSep: true})

	// ClientStreamErrors tracks how many errors were received when streaming
	// tagger events.
	ClientStreamErrors = telemetry.NewCounterWithOpts("tagger", "client_stream_errors",
		[]string{}, "Errors received when streaming tagger events",
		telemetry.Options{NoDoubleUnderscoreSep: true})

	// ServerStreamErrors tracks how many errors happened when streaming
	// out tagger events.
	ServerStreamErrors = telemetry.NewCounterWithOpts("tagger", "server_stream_errors",
		[]string{}, "Errors when streaming out tagger events",
		telemetry.Options{NoDoubleUnderscoreSep: true})

	// Subscribers tracks how many subscribers the tagger has.
	Subscribers = telemetry.NewGaugeWithOpts("tagger", "subscribers",
		[]string{}, "Number of channels subscribing to tagger events",
		telemetry.Options{NoDoubleUnderscoreSep: true})

	// Events tracks the number of tagger events being sent out.
	Events = telemetry.NewCounterWithOpts("tagger", "events",
		[]string{"cardinality"}, "Number of tagger events being sent out",
		telemetry.Options{NoDoubleUnderscoreSep: true})

	// Sends tracks the number of times the tagger has sent a
	// notification with a group of events.
	Sends = telemetry.NewCounterWithOpts("tagger", "sends",
		[]string{}, "Number of of times the tagger has sent a notification with a group of events",
		telemetry.Options{NoDoubleUnderscoreSep: true})

	// Receives tracks the number of times the tagger has received a
	// notification with a group of events.
	Receives = telemetry.NewCounterWithOpts("tagger", "receives",
		[]string{}, "Number of of times the tagger has received a notification with a group of events",
		telemetry.Options{NoDoubleUnderscoreSep: true})
)
