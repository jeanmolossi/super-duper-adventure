package services

import (
	"errors"
	"fmt"
	"github.com/jeanmolossi/super-duper-adventure/adapters/redis"
	solr2 "github.com/jeanmolossi/super-duper-adventure/adapters/solr"
	"github.com/jeanmolossi/super-duper-adventure/domain"
	"log"
	"regexp"
)

type StudentProcessor struct {
	Students        []domain.Student
	SolrConnection  *solr2.Connection
	RedisConnection *redis.Redis
}

func NewStudentProcessor() *StudentProcessor {
	return &StudentProcessor{
		SolrConnection:  solr2.NewSolr(),
		RedisConnection: redis.NewRedisConnection(),
	}
}

func (sp *StudentProcessor) Start() error {
	course := sp.Students[0].Curso
	redisClient := sp.RedisConnection
	defer redisClient.Client.Close()

	hasCacheCourse, err := redisClient.Get(course)
	if err != nil && err.Error() != "redis: nil" {
		return err
	}

	if hasCacheCourse != "" {
		e := fmt.Sprintf("Curso já existe no redis")
		return errors.New(e)
	}

	_, err = sp.SolrConnection.Update(map[string]interface{}{
		"add": sp.Students,
	}, true)

	if err != nil {
		return err
	}

	solrShard := replaceHost(sp.SolrConnection.URL)

	_, err = sp.RedisConnection.Set(course, solrShard)
	if err != nil {
		return err
	}

	log.Printf("Curso %s armazenado no shard: %s", course, solrShard)

	return nil
}

func replaceHost(host string) string {
	r := regexp.MustCompile("gsr_solr")
	replacedHost := r.ReplaceAllString(host, "localhost")

	return replacedHost
}
