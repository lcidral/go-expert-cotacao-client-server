# Sistema de Cotação USD-BRL

Sistema cliente-servidor em Go para consulta de cotações USD-BRL usando a API da AwesomeAPI.

## Pré-requisitos

- Go 1.21+
- SQLite3

## Como executar

1. Clone o projeto
```bash
git clone https://github.com/lcidral/go-expert-cotacao-client-server.git
cd go-expert-cotacao-client-server
```

2. Instale as dependências
```bash
go mod download
```

3. Execute o servidor
```bash
go run server.go
```

4. Em outro terminal, execute o cliente
```bash
go run client.go
```

## Funcionalidades

- Servidor (`/cotacao`):
    - Timeout de 200ms para API externa
    - Timeout de 10ms para banco de dados
    - Porta 8080

- Cliente:
    - Timeout de 300ms
    - Salva cotação em cotacao.txt