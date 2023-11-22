package handlers

import (
	"cian-parse/internals/app/processors"
	"cian-parse/internals/models"
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

type ImmovablesHandler struct {
	processor *processors.ImmovablesProcessor
	ctx       context.Context
}

func NewImmovablesHandler(ctx context.Context, processor *processors.ImmovablesProcessor) *ImmovablesHandler {
	handler := new(ImmovablesHandler)
	handler.processor = processor
	handler.ctx = ctx
	return handler
}

func (h *ImmovablesHandler) Create(w http.ResponseWriter, r *http.Request) {
	var newImmovable models.Immovable

	err := json.NewDecoder(r.Body).Decode(&newImmovable)
	if err != nil {
		WrapError(w, err)
		return
	}

	id, err := h.processor.Create(h.ctx, newImmovable)
	if err != nil {
		WrapError(w, err)
		return
	}

	var m = map[string]interface{}{
		"result": "OK",
		"data":   id,
	}

	WrapOK(w, m)
}

func (h *ImmovablesHandler) FindOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		WrapError(w, errors.New("missing id"))
		return
	}

	immovable, err := h.processor.FindOne(h.ctx, id)
	if err != nil {
		WrapError(w, err)
		return
	}

	var m = map[string]interface{}{
		"result": "OK",
		"data":   immovable,
	}

	WrapOK(w, m)
}

func (h *ImmovablesHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	immovables, err := h.processor.FindAll(h.ctx)
	if err != nil {
		WrapError(w, err)
	}

	var m = map[string]interface{}{
		"result": "OK",
		"data":   immovables,
	}

	WrapOK(w, m)
}

func (h *ImmovablesHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		WrapError(w, errors.New("missing id"))
		return
	}

	immovable, err := h.processor.FindOne(h.ctx, id)
	if err != nil {
		WrapError(w, err)
		return
	}

	err = h.processor.Update(h.ctx, id, immovable)
	if err != nil {
		WrapError(w, err)
		return
	}

	var m = map[string]interface{}{
		"result": "OK",
		"data":   "",
	}

	WrapOK(w, m)
}

func (h *ImmovablesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		WrapError(w, errors.New("missing id"))
		return
	}

	err := h.processor.Delete(h.ctx, id)
	if err != nil {
		WrapError(w, err)
		return
	}

	var m = map[string]interface{}{
		"result": "OK",
		"data":   "",
	}

	WrapOK(w, m)
}
