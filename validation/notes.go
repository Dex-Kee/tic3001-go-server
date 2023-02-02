package validation

import (
	"errors"
	"tic3001-go-server/common/dto"
	"tic3001-go-server/service"
)

type notesValidationService struct{}

var NotesValidationService = new(notesValidationService)

func (validation *notesValidationService) EntityExistedChecker(id string) error {
	e := service.NotesService.GetById(id)
	if e == nil {
		return errors.New("error: cannot found respective data for given id")
	}
	return nil
}

func (validation *notesValidationService) FormChecker(form dto.NotesForm) error {
	if form.Name == "" {
		return errors.New("error: name of notes cannot be blank")
	}
	if form.Description == "" {
		return errors.New("error: description of notes cannot be blank")
	}
	return nil
}
