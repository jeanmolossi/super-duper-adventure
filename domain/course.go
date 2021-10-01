package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

type Course struct {
	ID     string `json:"id"`
	Titulo string `json:"titulo"`
	Aulas  int    `json:"aulas"`
	Alunos string `json:"alunos"`
	Link   string `json:"link"`
}

func NewCourse() *Course {
	return &Course{}
}

func (c *Course) GetCourse(id string) (*Course, error) {
	request := NewRequest()
	request.Url = fmt.Sprintf("http://gsr_mock_api:3000/curso/%s/ver?populador", id)

	req, err := request.Request()
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}

	response, err := request.GetResponse(req)
	if err != nil || response.StatusCode != 200 {
		if err == nil {
			err = errors.New("api nao respondeu corretamente")
		}
		log.Printf(err.Error())
		return nil, err
	}

	err = json.Unmarshal(response.Body, c)
	if err != nil {
		log.Printf("Problema com a resposta da API: %s", err.Error())
		return nil, err
	}

	return c, nil
}
