package service

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"tic3001-go-server/common/dto"
	"tic3001-go-server/entity"
	"time"
)

type notesService struct{}

var NotesService = new(notesService)

// dbMap is used to act as database storage
var db = make(map[string]*entity.Notes)

func init() {
	// fake data for now
	// data := generateDataSet()
	// for i := range data {
	// 	database[data[i].Id] = &data[i]
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
	fmt.Printf("keyword: %s", keyword)
	// query from database
	keys := make([]string, 0)
	for _, e := range db {
		// not enter key, collect all keys
		if keyword == "" {
			keys = append(keys, e.Id)
		} else {
			if strings.Contains(e.Name, keyword) {
				keys = append(keys, e.Id)
			}
		}
	}
	return s.getSortedNotesByCreateDate(true, keys)
}

func (s *notesService) getSortedNotesByCreateDate(desc bool, keys []string) []dto.NotesDto {
	sort.SliceStable(keys, func(i, j int) bool {
		if desc {
			return db[keys[i]].CreateDate > db[keys[j]].CreateDate
		}
		return db[keys[i]].CreateDate < db[keys[j]].CreateDate
	})

	dtos := make([]dto.NotesDto, len(keys))
	for i, k := range keys {
		e := db[k]
		notesDto := dto.NotesDto{
			Id:          e.Id,
			Name:        e.Name,
			Description: e.Description,
			CreateDate:  e.CreateDate,
		}
		dtos[i] = notesDto
	}
	return dtos
}

func (s *notesService) Create(form dto.NotesForm) {
	e := entity.Notes{
		Id:          uuid.New().String(),
		Name:        form.Name,
		Description: form.Description,
		CreateDate:  time.Now().Format("2006-01-02 15:04:05"),
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
