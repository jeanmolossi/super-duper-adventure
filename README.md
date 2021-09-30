# Super Duper Adventure

O objetivo deste projeto é usar a seguinte stack:

- GoLang
- Solr
- RabbitMQ

e JS

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
