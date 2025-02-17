// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package replay

import (
	"context"
	"fmt"
	"time"

	"github.com/DataDog/datadog-agent/cmd/agent/api/response"
	pb "github.com/DataDog/datadog-agent/pkg/proto/pbgo"
	pbutils "github.com/DataDog/datadog-agent/pkg/proto/utils"
	"github.com/DataDog/datadog-agent/pkg/status/health"
	"github.com/DataDog/datadog-agent/pkg/tagger/collectors"
	"github.com/DataDog/datadog-agent/pkg/tagger/telemetry"
	"github.com/DataDog/datadog-agent/pkg/tagger/types"
	"github.com/DataDog/datadog-agent/pkg/util"
	"github.com/DataDog/datadog-agent/pkg/util/log"
)

// Tagger stores tags to entity as stored in a replay state.
type Tagger struct {
	store *tagStore

	ctx    context.Context
	cancel context.CancelFunc

	health          *health.Handle
	telemetryTicker *time.Ticker
}

// NewTagger returns an allocated tagger. You still have to run Init()
// once the config package is ready.
func NewTagger() *Tagger {
	return &Tagger{
		store: newTagStore(),
	}
}

// Init initializes the connection to the replay tagger and starts watching for
// events.
func (t *Tagger) Init() error {
	t.health = health.RegisterLiveness("tagger")
	t.telemetryTicker = time.NewTicker(1 * time.Minute)

	t.ctx, t.cancel = context.WithCancel(context.Background())

	return nil
}

// Stop closes the connection to the replay tagger and stops event collection.
func (t *Tagger) Stop() error {
	t.cancel()

	t.telemetryTicker.Stop()
	err := t.health.Deregister()
	if err != nil {
		return err
	}

	log.Info("replay tagger stopped successfully")

	return nil
}

// Tag returns tags for a given entity at the desired cardinality.
func (t *Tagger) Tag(entityID string, cardinality collectors.TagCardinality) ([]string, error) {
	entity, ok := t.store.getEntity(entityID)

	if !ok {
		telemetry.Queries.Inc(collectors.TagCardinalityToString(cardinality), telemetry.QueryEmptyTags)
		return []string{}, fmt.Errorf("Entity not found")
	}

	telemetry.Queries.Inc(collectors.TagCardinalityToString(cardinality), telemetry.QuerySuccess)
	return entity.GetTags(cardinality), nil
}

// TagBuilder returns tags for a given entity at the desired cardinality.
func (t *Tagger) TagBuilder(entityID string, cardinality collectors.TagCardinality, tb *util.TagsBuilder) error {
	tags, err := t.Tag(entityID, cardinality)
	if err == nil {
		tb.Append(tags...)

	}
	return err
}

// Standard returns the standard tags for a given entity.
func (t *Tagger) Standard(entityID string) ([]string, error) {
	entity, ok := t.store.getEntity(entityID)
	if !ok {
		return []string{}, fmt.Errorf("Entity not found")
	}

	return entity.StandardTags, nil
}

// List returns all the entities currently stored by the tagger.
func (t *Tagger) List(cardinality collectors.TagCardinality) response.TaggerListResponse {
	entities := t.store.listEntities()
	resp := response.TaggerListResponse{
		Entities: make(map[string]response.TaggerListEntity),
	}

	for _, e := range entities {
		resp.Entities[e.ID] = response.TaggerListEntity{
			Tags: map[string][]string{
				replaySource: e.GetTags(collectors.HighCardinality),
			},
		}
	}

	return resp
}

// Subscribe does nothing in the replay tagger this tagger does not respond to events.
func (t *Tagger) Subscribe(cardinality collectors.TagCardinality) chan []types.EntityEvent {
	// NOP
	return nil
}

// Unsubscribe does nothing in the replay tagger this tagger does not respond to events.
func (t *Tagger) Unsubscribe(ch chan []types.EntityEvent) {
	// NOP
}

// LoadState loads the state for the tagger from the supplied map.
func (t *Tagger) LoadState(state map[string]*pb.Entity) {

	if state == nil {
		return
	}

	// better stores these as the native type
	for id, entity := range state {
		entityID, err := pbutils.Pb2TaggerEntityID(entity.Id)
		if err != nil {
			log.Errorf("Error getting identity ID for %v: %v", id, err)
			continue
		}

		e := types.Entity{
			ID:                          entityID,
			HighCardinalityTags:         entity.HighCardinalityTags,
			OrchestratorCardinalityTags: entity.OrchestratorCardinalityTags,
			LowCardinalityTags:          entity.LowCardinalityTags,
			StandardTags:                entity.StandardTags,
		}

		err = t.store.addEntity(id, e)
		if err != nil {
			log.Errorf("Error storing identity with ID %v in store: %v", id, err)
		}
	}

	log.Debugf("Loaded %v elements into tag store", len(t.store.store))
}

// GetEntity returns the Entity for the supplied entity id..
func (t *Tagger) GetEntity(entityID string) (*types.Entity, error) {

	entity, ok := t.store.getEntity(entityID)
	if !ok {
		return nil, fmt.Errorf("No entity found for supplied id :%v", entityID)
	}

	return entity, nil
}
