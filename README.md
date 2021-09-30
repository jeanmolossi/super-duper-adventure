# Super Duper Adventure

O objetivo deste projeto é usar a seguinte stack:

- GoLang
- Solr
- RabbitMQ

e JS

## Primeiros passos

Certifique-se de ter as ferramentas instaladas:
- [Docker]()
- [Docker-compose]()

Em seguida, rode o comando:

`cp .env.example .env`

O comando acima é responsável por adicionar um arquivo de variáveis de ambiente legível pelo [docker-compose.yml](./docker-compose.yml)

O próximo passo é executar a aplicação. Para isso, na raiz do projeto, execute:

`docker-compose up -d`

Então, após, as seguintes urls estarão disponíveis para você visualizar os acontecimentos:

- http://localhost:3000 - Mock da API em JS
  - [/cursos](http://localhost:3000/cursos)
- [http://localhost:15672](http://localhost:15672) - RabbitMQ
- [http://localhost:8983](http://localhost:8983) - Apache Solr

## Fluxo de trabalho de aplicação

A ideia é:

Ao executar uma request na API de JS ela ira adicionar um 
novo 'id' de documentos na fila do RabbitMQ. 
Ao receber um novo Id na fila do RabbitMQ o Go irá consumir 
esse ID e fará uma nova request para pegar os filhos daquela 
referência.

Ao processar todos os filhos daquela referência o Go irá adicionar-los 
ao Solr.

Este projeto é um projeto de estudo e de teste de desempenho.
