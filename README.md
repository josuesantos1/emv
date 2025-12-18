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

### Passo 3: Execute o processador de transações

Em outro terminal, execute:

```bash
go run cmd/main.go
```

Você verá o prompt interativo:

```
EMV Transaction Processor
=========================
Enter TLV hex data (or 'exit' to quit)

TLV>
```

### Passo 4: Insira dados TLV

Cole o TLV hex e pressione Enter. Exemplo:

```
TLV> 5A0845395787636214865F2404251200009F340400000000

========== TRANSACTION RESULT ==========
Status: APPROVED
Message: Transaction authorized successfully
PAN: 4539578763621486
Expiry Date: 12/2025
CVM: 00000000
Timestamp: 2025-12-17 10:30:45
========================================

TLV>
```

O processador:
1. Decodifica os dados TLV do cartão
2. Valida os dados (PAN via Luhn, data de validade, CVM)
3. Envia para autorização no gateway
4. Exibe o resultado formatado
5. Registra em `transactions.json`
6. Aguarda nova entrada

Para sair, digite `exit` ou `quit`.

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
