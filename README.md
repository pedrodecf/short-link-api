# Short Link API Go

Este é um projeto simples de encurtador de URLs desenvolvido em GoLang com o framework Gin. Foi inspirado por um tutorial em Node.js disponível [aqui](https://www.youtube.com/watch?v=az7NpD02RM4).

## Funcionalidades

- Encurta URLs longas.
- Redireciona para URLs originais.
- Mantém um ranking das URLs mais acessadas utilizando Redis.

## Tecnologias Utilizadas

- PostgreSQL: para armazenamento dos dados das URLs encurtadas.
- Redis: para armazenamento do ranking.

## Como Usar

1. Clone o repositório.
2. Instale as dependências.
3. Configure o PostgreSQL e o Redis utilizando Docker Compose:

```bash
docker-compose up -d
```

4. Inicie o servidor.
5. Use os endpoints para encurtar URLs e acessar estatísticas.
