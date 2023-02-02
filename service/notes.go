package service

import (
	"encoding/json"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"tic3001-go-server/common/dto"
	"tic3001-go-server/entity"
)

type notesService struct{}

var NotesService = new(notesService)

// dbMap is used to act as database storage
var db = make(map[string]*entity.Notes)

func init() {
	// fake data for now
	// data := generateDataSet()
	// for i := range data {
	// 	db[data[i].Id] = &data[i]
	// }
	initDataSet()
}

func initDataSet() {
	file := "data.json"
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		log.Info("data.json file is not found ...")
		return
	}

	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Error("err when read from file: ", err)
		return
	}

	if len(content) == 0 {
		log.Info("empty content ...")
		return
	}

	list := make([]entity.Notes, 0)
	err = json.Unmarshal(content, &list)
	if err != nil {
		log.Error("err when parse file to json: ", err)
		return
	}

	for i := range list {
		db[list[i].Id] = &list[i]
	}
}

func generateRawDataSet() []entity.Notes {
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
	go s.persistData()
}

func (s *notesService) Update(form dto.NotesForm) {
	e := db[form.Id]
	e.Name = form.Name
	e.Description = form.Description
	go s.persistData()
}

func (s *notesService) Delete(id string) {
	delete(db, id)
	go s.persistData()
}

func (s *notesService) GetById(id string) *entity.Notes {
	notes, ok := db[id]
	if !ok {
		return nil
	}
	return notes
}

func (s *notesService) persistData() {
	list := make([]entity.Notes, 0)
	for _, e := range db {
		list = append(list, *e)
	}
	data, err := json.Marshal(list)
	if err != nil {
		log.Error("err: ", err.Error())
	}
	_ = ioutil.WriteFile("data.json", data, 0755)
}
