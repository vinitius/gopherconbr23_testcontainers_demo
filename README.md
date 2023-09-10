# Gophercon BR 2023 - Testcontainers: Elevando o nível dos seus testes de integração

---------------------

## Sobre

Uma simples demo de um serviço em Go que se conecta a um tópico [NATS](https://nats.io/) e persiste eventos recebidos em uma base de dados [Postgres](https://www.postgresql.org/).

O foco é demonstrar o poder e praticidade do [Testcontainers](https://golang.testcontainers.org/) para criar e executar testes de integracão como se fossem unitários.

>Apresentação realizada na [GopherCon BR 2023 - Test Containers: Elevando o nível dos seus testes de integração](https://gopherconbr.org/) .

Por:  [Andreia Silva](@andreiac-silva) & [Vinítius Salomão](@vinitius)


## Dependências
 
 - Go (>=1.19)
 - Docker
 - NATS
 - Postgres
 - Make (opcional)

## Bibliotecas

 - [Testcontainers](https://golang.testcontainers.org/)
 - [Bun](https://bun.uptrace.dev/)
 - [Nats](https://pkg.go.dev/github.com/nats-io/nats.go)
 - [Goleak](https://github.com/uber-go/goleak)

## Instruções

Para executar a suíte de teste:
```makefile
make test
```

Ou simplesmente:
```go
go test ./...
```

Ou ainda:
>Execute os testes diretamente da sua IDE favorita através do seu atalho favorito :tada:

## Link da talk
@todo
