# 📦 Consulta de CEP com Múltiplas APIs (ViaCEP e BrasilAPI)

Este projeto é uma aplicação em Go que realiza a consulta de endereço a partir de um CEP, utilizando **duas APIs em paralelo**:

- [ViaCEP](https://viacep.com.br)
- [BrasilAPI](https://brasilapi.com.br)

A aplicação usa a **resposta mais rápida** entre elas e **descarta a mais lenta**.  
Se nenhuma responder em até **1 segundo**, um erro de timeout é retornado.

---

## 🚀 Como Executar

1. Clone o repositório:

   ```bash
   git clone https://github.com/seu-usuario/seu-repositorio.git
   cd seu-repositorio
   ```

2. Execute a aplicação:

   ```bash
   go run main.go
   ```

3. Acesse no navegador ou via `curl`:

   ```
   http://localhost:8080/?cep=01001000
   ```

---

## ⚙️ Funcionamento

- O servidor escuta na porta `8080`.
- Ao receber uma requisição GET com o parâmetro `cep`, duas goroutines são iniciadas:
  - Uma consulta à **API ViaCEP**
  - Outra consulta à **BrasilAPI**
- Ambas enviam suas respostas para um canal buffered.
- Um `select` aguarda a **primeira resposta com sucesso**.
- A resposta é:
  - Impressa no terminal (origem + JSON formatado)
  - Enviada ao cliente via HTTP
- Se nenhuma API responder dentro de **1 segundo**, retorna:
  - Código HTTP `504 Gateway Timeout`
  - Mensagem: `"Timeout: nenhuma API respondeu a tempo"`

---

## ✅ Exemplo de Saída no Terminal

```
Servidor rodando em http://localhost:8080
Resposta recebida de: BrasilAPI
Dados do endereço:
{
  "cep": "01001-000",
  "state": "SP",
  "city": "São Paulo",
  "neighborhood": "Sé",
  "street": "Praça da Sé",
  "service": "correios"
}
```

---

## 🧱 Estrutura de Código

- **`main.go`**:
  - Define as structs `ViaCEP` e `BrasilApi`
  - Configura o servidor HTTP (`:8080`)
  - Implementa `BuscaViaCep` e `BuscaBrasilApi`
  - Dispara as buscas em goroutines concorrentes
  - Usa canais com buffer e `select` para responder com a mais rápida
  - Timeout implementado com `time.After(1 * time.Second)`

---

## 🔗 Endpoint

| Método | Endpoint | Parâmetros        | Exemplo                                 |
|--------|----------|-------------------|------------------------------------------|
| GET    | `/`      | `?cep=01001000`   | `http://localhost:8080/?cep=01001000`    |

---

## ❌ Tratamento de Erros

| Status | Descrição                                           |
|--------|-----------------------------------------------------|
| 400    | CEP não informado                                   |
| 500    | Erro interno na consulta do CEP                     |
| 504    | Timeout – nenhuma API respondeu dentro de 1 segundo |

---

## 📦 Dependências

- `net/http`
- `encoding/json`
- `time`
- `io`

---

## 📄 Licença

Este projeto está licenciado sob a [MIT License](LICENSE).
