# 🏦 EMV Transaction Processor

Um processador de transações EMV em Go que simula a comunicação entre um terminal de pagamento (POS) e um gateway de autorização.

## 📋 Índice

- [Sobre o Projeto](#sobre-o-projeto)
- [Funcionalidades](#funcionalidades)
- [Como Executar](#como-executar)
- [Testes](#testes)

## 🎯 Sobre o Projeto

Este projeto implementa um módulo básico de processamento de transações EMV conforme especificação EMV, incluindo:

- Parser de dados TLV (Tag-Length-Value)
- Validação de dados do cartão (PAN, data de validade, CVM)
- Comunicação com gateway de autorização (mock)

## ✨ Funcionalidades

### 1. Parser TLV EMV

Decodificação completa de estruturas TLV conforme EMV:  
- Extração dos seguintes campos:
  - `5A` - PAN (Primary Account Number)
  - `5F24` - Data de validade
  - `9F34` - CVM (Cardholder Verification Method)

### 2. Validações

- **PAN**:
  - Comprimento entre 13 e 19 dígitos
  - Validação via Algoritmo de Luhn
- **Data de Validade**:
  - Não pode ser anterior à data atual
- **CVM**:
  - Validação de métodos suportados (bits 1, 2 e 3)

### 3. Autorização

- Gateway HTTP para comunicação com servidor acquirer (mock)
- Servidor mock de autorização (70% de aprovação)


## 🚀 Como Executar

### Pré-requisitos

- Go 1.25.0 ou superior

### Passo 1: Clone o repositório

```bash
git clone https://github.com/josuesantos1/emv.git
cd emv
```

### Passo 2: Inicie o servidor mock de autorização

Em um terminal, execute:

```bash
go run cmd/acquirer/main.go
```

O servidor iniciará na porta 8080:
```
Mock server Acquirer running on port :8080
```

### Passo 3: Execute a aplicação principal

Em outro terminal, execute:

```bash
go run cmd/main.go
```

A aplicação irá:
1. Decodificar os dados TLV do cartão
2. Validar os dados (PAN, data de validade, CVM)
3. Enviar para autorização no gateway
4. Registrar o resultado em `transactions.json`

## 🧪 Testes

Execute todos os testes:

```bash
go test ./... -v
```

Execute testes de um pacote específico:

```bash
# Testes do parser
go test ./pkg/tlv -v

# Testes do domínio
go test ./internal/domain -v
```

### Cobertura de Testes

- **Parser TLV**: Testes de Parse, ParseTag, ParseLength
- **Validações**: Testes de PAN (Luhn), Data de Validade, CVM
- **Populate**: Testes de extração e conversão de dados TLV
