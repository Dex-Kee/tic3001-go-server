package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"tic3001-go-server/common/constant"
	"tic3001-go-server/common/dto"
)

var engine = setupRouter()

func TestMain(t *testing.M) {
	t.Run()
	clearDataSet(engine)
}

func setupRouter() *gin.Engine {
	e := gin.Default()
	notes := e.Group("/api/notes")
	{
		notes.GET("/list", NotesController.List)
		notes.POST("/create", NotesController.Create)
		notes.PUT("/update", NotesController.Update)
		notes.DELETE("/delete", NotesController.Delete)
	}
	return e
}

func TestGetWithSuccess(t *testing.T) {
	recorder := doHttpRequest(engine, http.MethodGet, "/api/notes/list", []byte{})

	// check status code
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestGetWithEmptyResultSuccess(t *testing.T) {
	queryParam := "keyword=NIL-INFO"
	recorder := doHttpRequest(engine, http.MethodGet, "/api/notes/list"+"?"+queryParam, []byte{})

	// check status code
	assert.Equal(t, http.StatusOK, recorder.Code)

	resp := parserRespDto(recorder.Body.Bytes())

	// check data length
	list := resp.Data.([]interface{})
	assert.Equal(t, 0, len(list))
}

func TestPostWithSuccess(t *testing.T) {
	form := new(dto.NotesForm)
	form.Name = "tic3002 assignment"
	form.Description = "assignment description"
	data, _ := json.Marshal(form)

	recorder := doHttpRequest(engine, http.MethodPost, "/api/notes/create", data)

	// check status
	assert.Equal(t, http.StatusOK, recorder.Code)

	// query again, check if content match
	m := getData(engine)
	assert.Equal(t, "tic3002 assignment", m["name"])
	assert.Equal(t, "assignment descriptionxxx", m["description"])

	t.Cleanup(func() {
		clearDataSet(engine)
	})
}

func TestPostWithEmptyName(t *testing.T) {
	nameEmptyForm := new(dto.NotesForm)
	nameEmptyForm.Name = ""
	nameEmptyForm.Description = "deploy web server via docker"
	data, _ := json.Marshal(nameEmptyForm)

	recorder := doHttpRequest(engine, http.MethodPost, "/api/notes/create", data)

	// check status code
	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	// check response code
	resp := parserRespDto(recorder.Body.Bytes())
	assert.Equal(t, constant.RespCodeClientParamError, resp.Code)
}

func TestPostWithEmptyDescription(t *testing.T) {
	nameEmptyForm := new(dto.NotesForm)
	nameEmptyForm.Name = "tic3001 assignment"
	nameEmptyForm.Description = ""
	data1, _ := json.Marshal(nameEmptyForm)

	recorder := doHttpRequest(engine, http.MethodPost, "/api/notes/create", data1)

	// check status code
	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	// check response code
	resp := parserRespDto(recorder.Body.Bytes())
	assert.Equal(t, constant.RespCodeClientParamError, resp.Code)
}

func TestUpdateWithSuccess(t *testing.T) {
	// prefill one data
	form := new(dto.NotesForm)
	form.Name = "tic3001 assignment"
	form.Description = "assignment description"
	data, _ := json.Marshal(form)
	recorder := doHttpRequest(engine, http.MethodPost, "/api/notes/create", data)

	// get data id
	m := getData(engine)

	// build update form
	form = new(dto.NotesForm)
	form.Id = m["id"].(string)
	form.Name = "tic4001 assignment"
	form.Description = "new assignment description"
	data, _ = json.Marshal(form)
	recorder = doHttpRequest(engine, http.MethodPut, "/api/notes/update", data)

	// check status code
	assert.Equal(t, http.StatusOK, recorder.Code)

	// get data again, check if fields are updated
	m = getData(engine)
	assert.Equal(t, "tic4001 assignment", m["name"])
	assert.Equal(t, "new assignment description", m["description"])

	t.Cleanup(func() {
		clearDataSet(engine)
	})
}

func TestUpdateWithInvalidId(t *testing.T) {
	form := new(dto.NotesForm)
	form.Id = "abcd1234"
	form.Name = "tic4001 assignment"
	form.Description = "new assignment description"
	data, _ := json.Marshal(form)
	recorder := doHttpRequest(engine, http.MethodPut, "/api/notes/update", data)

	// check status code
	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	// check resp data
	respDto := parserRespDto(recorder.Body.Bytes())
	assert.Equal(t, constant.RespCodeClientParamError, respDto.Code)
}

func TestDeleteWithSuccess(t *testing.T) {
	// prefill one data
	form := new(dto.NotesForm)
	form.Name = "tic3001 assignment"
	form.Description = "assignment description"
	data, _ := json.Marshal(form)
	doHttpRequest(engine, http.MethodPost, "/api/notes/create", data)

	// get data id
	m := getData(engine)

	queryParam := "id=" + m["id"].(string)
	recorder := doHttpRequest(engine, http.MethodDelete, "/api/notes/delete?"+queryParam, []byte{})

	// check status
	assert.Equal(t, http.StatusOK, recorder.Code)

	// query again, check if the list is empty
	list := fetchDataList(engine)
	assert.Equal(t, 0, len(list))
}

func TestDeleteWithInvalidId(t *testing.T) {
	queryParam := "id=abcd1234"
	recorder := doHttpRequest(engine, http.MethodDelete, "/api/notes/delete?"+queryParam, []byte{})

	// check status
	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	respDto := parserRespDto(recorder.Body.Bytes())
	assert.Equal(t, constant.RespCodeClientParamError, respDto.Code)
}

func getData(engine *gin.Engine) map[string]interface{} {
	list := fetchDataList(engine)
	if len(list) == 0 {
		return map[string]interface{}{}
	}
	return list[0].(map[string]interface{})
}

func fetchDataList(engine *gin.Engine) []interface{} {
	recorder := doHttpRequest(engine, "GET", "/api/notes/list", []byte{})
	respDto := parserRespDto(recorder.Body.Bytes())
	return respDto.Data.([]interface{})
}

func doHttpRequest(engine *gin.Engine, method string, path string, data []byte) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, path, bytes.NewReader(data))
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, request)
	return recorder
}

func parserRespDto(data []byte) *dto.ResponseDto {
	resp := new(dto.ResponseDto)
	err := json.Unmarshal(data, resp)
	if err != nil {
		panic("error when perform http calls")
	}
	return resp
}

func clearDataSet(engine *gin.Engine) {
	list := fetchDataList(engine)
	if len(list) == 0 {
		_ = os.Remove("data.json")
		return
	}
	for _, e := range list {
		m := e.(map[string]interface{})
		doHttpRequest(engine, http.MethodDelete, fmt.Sprintf("/api/notes/delete?id=%s", m["id"]), []byte{})
	}
	_ = os.Remove("data.json")
}
