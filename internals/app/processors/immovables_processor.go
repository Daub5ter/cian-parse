package processors

import (
	"cian-parse/internals/app/db"
	"cian-parse/internals/models"
	"context"
	"fmt"
)

var errIdIsEmpty = fmt.Errorf("the id of immovable should be not empty")

type ImmovablesProcessor struct {
	storage *db.ImmovablesStorage
}

func NewImmovablesProcessor(storage *db.ImmovablesStorage) *ImmovablesProcessor {
	processor := new(ImmovablesProcessor)
	processor.storage = storage
	return processor
}

func (p *ImmovablesProcessor) Create(ctx context.Context, immovable models.Immovable) (string, error) {
	if immovable.Link == "" || immovable.Title == "" || immovable.Price == 0 {
		return "", fmt.Errorf("the parts of immovable should not be emply")
	}

	return p.storage.Create(ctx, immovable)
}

func (p *ImmovablesProcessor) FindOne(ctx context.Context, id string) (models.Immovable, error) {
	if id == "" {
		return models.Immovable{}, errIdIsEmpty
	}

	return p.storage.FindOne(ctx, id)
}

func (p *ImmovablesProcessor) FindAll(ctx context.Context) ([]models.Immovable, error) {
	return p.storage.FindAll(ctx)
}

func (p *ImmovablesProcessor) Update(ctx context.Context, id string, immovable models.Immovable) error {
	if id == "" {
		return errIdIsEmpty
	}

	return p.storage.Update(ctx, id, immovable)
}

func (p *ImmovablesProcessor) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errIdIsEmpty
	}

	return p.storage.Delete(ctx, id)
}
