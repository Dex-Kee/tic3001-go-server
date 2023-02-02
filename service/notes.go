package service

import (
	"github.com/google/uuid"
	"tic3001-go-server/common/dto"
	"tic3001-go-server/entity"
)

type notesService struct{}

var NotesService = new(notesService)

// dbMap is used to act as database storage
var db = make(map[string]*entity.Notes)

func init() {
	// fake data for now
	data := generateDataSet()
	for i := range data {
		db[data[i].Id] = &data[i]
	}
}

func generateDataSet() []entity.Notes {
	e1 := entity.Notes{
		Id:          uuid.New().String(),
		Name:        "do tic3001 assignment1",
		Description: "deploy docker & k8s",
	}
	e2 := entity.Notes{
		Id:          uuid.New().String(),
		Name:        "do tic3001 assignment2",
		Description: "build a simple web server",
	}
	e3 := entity.Notes{
		Id:          uuid.New().String(),
		Name:        "do tic3001 assignment3",
		Description: "implement the cache for large data set query",
	}
	return []entity.Notes{e1, e2, e3}
}

func (s *notesService) List(keyword string) []dto.NotesDto {
	dtos := make([]dto.NotesDto, 0)
	// query from db
	for _, e := range db {
		dtos = append(dtos, dto.NotesDto{
			Id:          e.Id,
			Name:        e.Name,
			Description: e.Description,
		})
	}
	return dtos
}

func (s *notesService) Create(form dto.NotesForm) {
	e := entity.Notes{
		Id:          uuid.New().String(),
		Name:        form.Name,
		Description: form.Description,
	}
	db[e.Id] = &e
}

func (s *notesService) Update(form dto.NotesForm) {
	e := db[form.Id]
	e.Name = form.Name
	e.Description = form.Description
}

func (s *notesService) Delete(id string) {
	delete(db, id)
}
