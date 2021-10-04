# Super Duper Adventure

O objetivo deste projeto é usar a seguinte stack:

- GoLang
- Solr
- RabbitMQ
- Redis

e JS

## Tabela de conteúdo

- [Primeiros passos](#primeiros-passos)
- [Fluxo de trabalho da aplicação](#fluxo-de-trabalho-da-aplicação)
- [Utilidades](#utilidades)

## Primeiros passos

Certifique-se de ter as ferramentas instaladas:

- [Docker](https://docs.docker.com/get-docker/)
- [Docker-compose](https://docs.docker.com/compose/install/)

Em seguida, rode o comando:

```shell
cp .env.example .env
```

O comando acima é responsável por adicionar um arquivo de variáveis de ambiente.

O próximo passo é executar a aplicação. Para isso, na raiz do projeto, execute:

```shell
docker-compose up -d
```

Por fim, para processar a fila execute:

```shell
docker exec gsr_go go run .
```

Então, após, as seguintes urls estarão disponíveis para você visualizar os acontecimentos:

- http://localhost:3000 - Mock da API em JS
    - [/courses](http://localhost:3000/courses)
    - [/courses?ids={ ids_separados_por_virgula }](http://localhost:3000/courses?ids=)
    - [/course/{ id }](http://localhost:3000/course/1)
    - [/course/{ id }/students](http://localhost:3000/course/1/students)
    - [/students](http://localhost:3000/students)
    - [/students?ids={ ids_separados_por_virgula }](http://localhost:3000/students?ids=)
    - [/student/{ id }](http://localhost:3000/student/1)

- [http://localhost:15672](http://localhost:15672) - RabbitMQ
- [http://localhost:8983](http://localhost:8983) - Apache Solr
- [http://localhost:6379](http://localhost:6379) - Redis

## Fluxo de trabalho de aplicação

A ideia é:

Ao executar uma request na API de JS ela ira adicionar um novo 'id' de documentos na fila do RabbitMQ. Ao receber um
novo Id na fila do RabbitMQ o Go irá consumir esse ID e fará uma nova request para pegar os filhos daquela referência.

Ao processar todos os filhos daquela referência o Go irá adicionar-los ao Solr e criará a chave de cache do solr no
redis.

Este projeto é um projeto de estudo e de teste de desempenho.

## Utilidades

### Mock api JS

Você pode executar os seguintes curl para popular a memória da API **(OBS: os registros da api ficam em memória. Ao
reiniciar a api estes registros são perdidos)**

**Seed de Cursos**

```shell
curl -X GET 'http://localhost:3000/seed/courses' -H 'Content-Type: application/json'
```

**Seed de Estudantes**

```shell
curl -X GET 'http://localhost:3000/seed/students' -H 'Content-Type: application/json'
```

**Seed de Fila do Rabbit**

```shell
curl -X GET 'http://localhost:3000/seed/queue' -H 'Content-Type: application/json'
```

### Redis

Acesse o redis `docker exec -it gsr_redis redis-cli -h localhost` e execute o comando:

```shell
// PARA PEGAR TODAS AS CHAVES DO REDIS 
KEYS course_*

// PARA PEGAR O ENDEREÇO DO SHARD DO CURSO ID 123
GET course_123
```

### Apache Solr

Acesse a [Dashboard do Solr](http://localhost:8983) e selecione o Shard/Core retornado pelo [Redis](#redis). Na barra à
esquerda selecione "Query" e no campo de query você pode usar a query: `curso:123`. Com isso você terá todos os
registros daquele curso dentro do Apache Solr

### RabbitMQ

Para visualizar o processamento da fila do RabbitMQ basta acessar a [Dashboard do RabbitMQ](http://localhost:15672) com
o usuário `rabbitmq` e a senha `rabbitmq`.

Selecionar Queues e em seguida a fila queueCourseId