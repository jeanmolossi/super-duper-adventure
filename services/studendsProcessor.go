package services

import (
	"encoding/json"
	"fmt"
	"github.com/jeanmolossi/super-duper-adventure/domain"
	"github.com/jeanmolossi/super-duper-adventure/solr"
	"log"
)

type GetStudentsResponse struct {
	Curso  int       `json:"curso"`
	Total  int       `json:"total"`
	Alunos []Student `json:"users"`
}

type Student struct {
	ID    string `json:"id"`
	Nome  string `json:"nome"`
	Plano int    `json:"plano"`
	Email string `json:"email"`
	Curso string `json:"curso"`
}

type StudentProcessor struct {
	CourseResultChannel chan ProcessedCourseResult
	Students            []Student
	SolrConnection      *solr.Connection
}

func NewStudentProcessor(courseResultChannel chan ProcessedCourseResult) *StudentProcessor {
	return &StudentProcessor{
		CourseResultChannel: courseResultChannel,
		SolrConnection:      solr.NewSolr(),
	}
}

func (sp *StudentProcessor) Start(endsChannel chan string) {
	for courseResult := range sp.CourseResultChannel {
		students, err := sp.GetStudents(courseResult.Course)
		if err != nil {
			log.Fatalf("Falhou ao buscar estudantes do Curso %s", courseResult.Course.ID)
			return
		}
		_, err = sp.SolrConnection.Update(map[string]interface{}{
			"add": students,
		}, true)
		if err != nil {
			courseResult.Message.Reject(false)
			log.Fatalf("Erro ao adicionar ao solr")
			return
		}

		courseResult.Message.Ack(false)

		log.Printf("Curso %s processado. (Total %d alunos)", courseResult.Course.ID, len(students))

		endsChannel <- courseResult.Course.ID
	}
	endsChannel <- "all-done"
}

func (sp *StudentProcessor) GetStudents(course domain.Course) ([]Student, error) {
	request := domain.NewRequest()
	request.Url = fmt.Sprintf(course.Alunos)

	req, err := request.Request()
	if err != nil {
		log.Fatalf("Falha ao executar busca: %s", err.Error())
		return nil, err
	}

	response, err := request.GetResponse(req)
	if err != nil {
		log.Fatalf("Falha ao obter resposta: %s", err.Error())
		return nil, err
	}

	var result GetStudentsResponse
	err = json.Unmarshal(response.Body, &result)
	if err != nil {
		log.Fatalf("Erro ao decodificar entidade: %s", err.Error())
		return nil, err
	}

	return result.Alunos, nil
}
