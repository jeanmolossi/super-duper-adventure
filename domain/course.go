package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

type Course struct {
	ID      string `json:"id"`
	Titulo  string `json:"titulo"`
	Aulas   int    `json:"aulas"`
	Alunos  string `json:"alunos"`
	HotSite string `json:"hot_site"`
}

func NewCourse() *Course {
	return &Course{}
}

func (c *Course) GetCourse(id string) error {
	request := NewRequest()
	request.Url = fmt.Sprintf("http://gsr_mock_api:3000/course/%s", id)

	req, err := request.Request()
	if err != nil {
		log.Printf(err.Error())
		return err
	}

	response, err := request.GetResponse(req)
	if err != nil || response.StatusCode != 200 {
		if err == nil {
			err = errors.New("api nao respondeu corretamente")
		}
		log.Printf(err.Error())
		return err
	}

	err = json.Unmarshal(response.Body, c)
	if err != nil {
		log.Printf("Problema com a resposta da API: %s", err.Error())
		return err
	}

	return nil
}

func (c *Course) GetStudents() ([]Student, error) {
	request := NewRequest()
	request.Url = fmt.Sprintf("http://gsr_mock_api:3000/course/%s/students", c.ID)

	req, err := request.Request()
	if err != nil {
		return nil, err
	}

	response, err := request.GetResponse(req)
	if err != nil {
		return nil, err
	}

	var students GetStudentsResponse
	err = json.Unmarshal(response.Body, &students)
	if err != nil {
		return nil, err
	}

	return students.Alunos, nil
}
