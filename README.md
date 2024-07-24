# leilaoGolang
Desafio Concorrência com Golang - Leilão - Go Express

## Índice
1. [Objetivo](#objetivo)
2. [Descrição](#descrição)
3. [Pré-requisitos](#pré-requisitos)
4. [Configurações](#configurações)
5. [Executando o projeto](#executando-o-projeto)


## Objetivo
O objetivo deste projeto é utilizar go routines para adicionar a rotina de fechamento automático do leilão a partir de um tempo.

## Descrição
Uma função irá calcular o tempo do leilão, baseado em parâmetros previamente definidos em variáveis de ambiente permitindo que uma nova go routine valide a existência de um leilão (auction) vencido (em que o tempo já se esgotou) e realizar o update fechando o mesmo. Tambémm é feito um teste para validar se o fechamento está acontecendo de forma automatizada;

## Pré-requisitos
Assegure-se de ter as seguintes ferramentas instaladas:
- [Golang](https://go.dev/doc/install)
- [Docker](https://docs.docker.com/compose/install/)
- A instalação do [MongoDB Compass](https://www.mongodb.com/docs/compass/current/install/) é opcional

## Configurações
O projeto é configurado no arquivo `./cmd/auction/.env`.

As configurações são: 
- BATCH_INSERT_INTERVAL= Tempo de espera da aplicação para inserir um lote de lances no bando de dados;
- MAX_BATCH_SIZE= Quantidade de lances para fechar um lote e inserir no banco de dados;
- AUCTION_INTERVAL= Tempo que um leilão permanece aberto (recebendo lances) após a sua criação;

## Executando o projeto

O projeto possui um Makefile com comandos utilitários para execução do projeto listados abaixo:

- Constrói as imagens e sobe os containers do projeto na porta 8080:
```
$ make build-run
```

- Sobe os containers do projeto sem build na porta 8080:
```
$ make run
```

- Encerra a execução dos containers do projeto:
```
$ make stop
```
- Sobe o Mongo-DB na porta 27017:
```
$ make run-mongo
```
### Execução por API
Dentro da pasta ./docs existe uma collection do Postman com 3 chamadas HTTP necessárias para realizar o teste de integração.
- **POST /auctions:** Chamada HTTP para criar um leilão. É possível trocar a chamada do Postman pelo seguinte CURL:
```
curl --location 'localhost:8080/auctions' \
--header 'Content-Type: application/json' \
--data '{
    "product_name": "Ferrari",
    "category": "carros",
    "description": "Leilão de carros"
}'
```
- **GET /auctions:** Chamada HTTP para buscar leilões por status. É possível trocar a chamada do Postman pelo seguinte CURL:
```
curl --location 'localhost:8080/auctions?status=0'
```
- **POST /bid:** Chamada HTTP para dar lances no leilão. Para fazer essa chamada é necessário passar a ID to leilão criado pela chamada POST /auctions, através da chamda GET /auctions acima ou por uma IDE de acesso ao MongoDB (sugerido o MongoDB Compass). É possível trocar a chamada do Postman pelo seguinte CURL:
```
curl --location 'localhost:8080/bid' \
--header 'Content-Type: application/json' \
--data '{
    "user_id": "18ab66ba-6c2e-4917-b9fa-5ae47018b1d3",
    "auction_id": "[INSIRA AQUI O ID DO LEILÃO CRIADO]",
    "amount": 1000.00
}'
```

### Execução por teste automatizado
Após fazer o clone deste projeto, abra o terminal no diretório raiz do projeto e execute os comandos abaixo:
```
$ make run-mongo
$ go mod tidy
$ go test ./... -count=1
```
O ultimo comando executará um teste que foi codificado dentro do diretório `./internal/infra/database/auction/create_auction_test.go`, que executará um teste de integração em um database de teste. Para que funcione adequadamente, o container do mongodb precisa estar no ar com as configurações padrão definidas no arquivo `cmd/auction/.env`.
