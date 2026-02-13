# Exercícios — Docker para Go

## Exercício 17.1 — Multi-stage Build
Pegue a API do módulo 14 (api-pura) e crie um Dockerfile multi-stage.
- Use `golang:1.22-alpine` como builder
- Use `alpine:3.19` como runtime
- Injete a versão via `--build-arg`
- Rode e teste com `curl`
- Compare o tamanho da imagem com e sem multi-stage (`docker images`)

## Exercício 17.2 — Docker Compose
Crie um `docker-compose.yml` que sobe:
- Sua API Go (porta 8080)
- PostgreSQL (porta 5432)
- Redis (porta 6379)
- A API deve esperar o PostgreSQL estar saudável antes de iniciar
- Use `volumes` para persistir dados do banco
- Use `environment` para configurar a API

Teste:
```bash
docker compose up --build
curl http://localhost:8080/health
docker compose down
docker compose up  # dados devem persistir
```

## Exercício 17.3 — Imagem Scratch
Refaça o Dockerfile usando `FROM scratch`:
- Compile com `CGO_ENABLED=0`
- Copie os certificados SSL
- Compare o tamanho com Alpine (~10MB vs ~5MB)
- Tente rodar `docker exec -it container sh` e veja que não funciona (não tem shell!)

## Exercício 17.4 — CI/CD Simulado
Crie um `Makefile` com os seguintes targets:
```makefile
build:      # compila o binário
test:       # roda testes com coverage
lint:       # roda golangci-lint
docker:     # builda a imagem Docker
run:        # roda com docker compose
stop:       # para os containers
clean:      # remove binários e imagens
```
Execute: `make test lint docker run` e verifique que tudo funciona em sequência.
