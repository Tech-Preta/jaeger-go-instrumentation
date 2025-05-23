# Go Example Service with Jaeger Instrumentation

Este é um exemplo de aplicação Go que demonstra a instrumentação com Jaeger para observabilidade.

## Requisitos

- Go 1.22 ou superior
- Docker
- Kubernetes (opcional, para deploy com Helm)
- Jaeger (para visualização dos traces)

## Boas Práticas

### Configuração de Endpoints

É uma boa prática **nunca** hardcodear endpoints de serviços externos (como o Jaeger) no código ou em arquivos de configuração. Isso porque:

1. **Portabilidade**: O código deve ser portável entre diferentes ambientes (desenvolvimento, homologação, produção)
2. **Segurança**: Evita exposição desnecessária de endpoints internos
3. **Manutenibilidade**: Facilita mudanças de configuração sem necessidade de recompilar o código
4. **Flexibilidade**: Permite diferentes configurações em diferentes ambientes

Nesta aplicação, seguimos estas práticas:
- O endpoint do Jaeger é configurado via variável de ambiente.
- No Kubernetes, usamos o DNS interno do cluster.
- Localmente, usamos o endpoint do ingress.
- Nenhum endpoint está hardcoded no código.
## Configuração

A aplicação utiliza as seguintes variáveis de ambiente:

- `PORT`: Porta em que a aplicação irá rodar (padrão: 8080)
- `JAEGER_ENDPOINT`: Endpoint do Jaeger Collector (obrigatório)

## Executando Localmente

1. Instale as dependências:
```bash
cd src
go mod download
```

2. Execute a aplicação:
```bash
go run main.go
```

## Build e Execução com Docker

1. Build da imagem:
```bash
docker build -t nataliagranato/jaeger-go-instrumentation:0.1.1 .
```

2. Execute o container:
```bash
# Para execução local (fora do cluster Kubernetes)
docker run -p 8080:8080 -e JAEGER_ENDPOINT="https://jaeger-http-collector.nataliagranato.xyz/api/traces" nataliagranato/jaeger-go-instrumentation:0.1.1
```

## Deploy com Helm

1. Instale o chart:
```bash
helm install jaeger-go-instrumentation ./charts
```

2. Para atualizar o endpoint do Jaeger, edite o arquivo `charts/values.yaml` e atualize o valor de `JAEGER_ENDPOINT`.

## Configuração do Jaeger

### No Kubernetes (dentro do cluster)
Quando a aplicação está rodando dentro do cluster Kubernetes, usamos o DNS interno do Kubernetes:

```yaml
env:
  - name: JAEGER_ENDPOINT
    value: "http://jaeger-collector.jaeger.svc.cluster.local:14268/api/traces"
```

Onde:
- `jaeger-collector`: nome do serviço do Jaeger
- `jaeger`: namespace onde o Jaeger está instalado
- `14268`: porta HTTP do collector

### Localmente (fora do cluster)
Quando a aplicação está rodando localmente (fora do cluster), usamos o endpoint do ingress do Jaeger:

```bash
JAEGER_ENDPOINT="https://jaeger-http-collector.nataliagranato.xyz/api/traces"
```

## Testando

Após iniciar a aplicação, você pode testar enviando uma requisição para o endpoint:

```bash
curl http://localhost:8080/
```

Os traces serão enviados para o Jaeger e podem ser visualizados na interface do Jaeger UI.

### Visualizando os Traces

1. Acesse a interface do Jaeger UI:
```bash
kubectl port-forward -n jaeger svc/jaeger-query 16686:16686
```

2. Abra o navegador em `http://localhost:16686`

3. Na interface do Jaeger:
   - Selecione o serviço "go-example-service"
   - Clique em "Find Traces"
   - Você verá os traces gerados pela aplicação

Cada trace mostrará:
- Nome da operação (helloHandler)
- Duração da operação
- Atributos configurados (event e message)

## Estrutura do Projeto

```
.
├── Dockerfile
├── README.md
├── charts/
│   ├── Chart.yaml
│   ├── templates/
│   │   ├── deployment.yaml
│   │   └── service.yaml
│   └── values.yaml
└── src/
    ├── go.mod
    ├── go.sum
    └── main.go
```
