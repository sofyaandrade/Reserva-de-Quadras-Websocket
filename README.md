# 🏐 Reserva de Quadras Esportivas

Projeto desenvolvido na disciplina de **Programação Orientada a Eventos** do curso de **Sistemas de Informação** na **UNIDAVI**.

---

## 📖 Descrição

Este projeto simula um **sistema de reserva de quadras esportivas** com **comunicação em tempo real** entre o **frontend (React)** e o ** (Go)**, utilizando **WebSockets**.

O sistema permite que múltiplos usuários (ou abas) vejam as ações em tempo real — como criação, cancelamento e finalização de reservas — de forma sincronizada entre todos os clientes conectados.

---

## 🧩 Estrutura do Projeto

A aplicação é dividida em duas partes:

- **Backend (Go):** Gerencia os eventos, reservas, estado do sistema e quadras disponíveis, além de manter a comunicação WebSocket com todos os clientes conectados.
- **Frontend (React + Vite):** Interface gráfica para administração e interação com o sistema em tempo real.

---

## ⚙️ Requisitos

- [Go (1.22 ou superior)](https://go.dev/dl/)
- [Node.js (18 ou superior)](https://nodejs.org/en/download/)
- npm (instalado junto com o Node.js)

---

## 🧭 Configurando o Ambiente Go

1️⃣ **Verifique se o Go está instalado**
```bash
go version
```
Se aparecer algo como `go version go1.22.0 linux/amd64`, está tudo certo ✅

2️⃣ **Configure o GOPATH (se ainda não existir)**  
Por padrão, o Go cria uma pasta para seus projetos em:
- **Windows:** `C:\Users\SeuUsuario\go`
- **Linux/macOS:** `~/go`

Verifique com:
```bash
go env GOPATH
```

3️⃣ **Adicione o Go ao PATH (caso necessário)**  
No Linux/macOS:
```bash
export PATH=$PATH:$(go env GOPATH)/bin
```
---

## 🚀 Execução do Projeto

### 1️⃣ Clone o repositório
```bash
git clone https://github.com/sofyaandrade/Reserva-de-Quadras-Websocket.git
cd Reserva-de-Quadras-Websocke
```

---

### 2️⃣ Execute o Backend (Go)
```bash
cd websocket
go mod tidy
go run main.go
```

O servidor será iniciado em:
```
http://localhost:8080
```

---

### 3️⃣ Execute o Frontend (React + Vite)
Abra outro terminal:
```bash
cd frontend
npm install
npm run dev
```

Acesse no navegador:
```
http://localhost:5173
```

---

## 💻 Funcionalidades Principais

✅ **Administração**
- Adicionar novas quadras com nome e capacidade.
- Iniciar e parar o sistema de reservas.

✅ **Reservas em tempo real**
- Criação, confirmação, cancelamento e finalização automática.
- Atualizações instantâneas em todas as abas via WebSocket.

✅ **Validação de conflitos**
- Impede reservas em horários sobrepostos.

✅ **Temporizadores**
- Mostra o tempo para confirmação e o tempo restante de uso da quadra.

✅ **Feedbacks visuais**
- Logs em tempo real.
- Linhas de “reserva negada” exibidas na tabela com o motivo.

---

## 📦 Dependências

### Backend (Go)
| Biblioteca | Função |
|-------------|--------|
| `net/http` | Servidor HTTP nativo |
| `github.com/gorilla/websocket` | Comunicação em tempo real |
| `sync` | Controle de concorrência com Mutex |

### Frontend (React + Vite)
| Pacote | Função |
|--------|--------|
| `react` | Biblioteca para construção da interface |
| `vite` | Build e servidor de desenvolvimento rápido |
| `eslint` | Boas práticas e linting do código |

---

## 🔁 Comunicação em Tempo Real

A comunicação entre o **frontend** e o **** ocorre por **eventos WebSocket**.

| Evento | Origem | Descrição |
|--------|---------|-----------|
| `sistema.iniciado` | Admin → Todos | Libera reservas |
| `sistema.parado` | Admin → Todos | Bloqueia reservas |
| `quadra.adicionada` | Admin → Todos | Adiciona nova quadra |
| `horario.reservado` | Usuário → Todos | Cria nova reserva |
| `horario.confirmado` | Servidor → Todos | Confirma reserva após tempo limite |
| `reserva.cancelada` | Usuário → Todos | Cancela reserva |
| `reserva.negada` | Servidor → Usuário | Reserva negada (horário em conflito) |
| `jogo.finalizado` | Servidor → Todos | Finaliza o uso da quadra |

---

## 🧠 Como Testar

1. Abra o frontend no navegador (`http://localhost:5173`).
2. Abra mais de uma aba (simulando usuários diferentes).
3. Clique em **“Iniciar sistema”** em uma delas.
4. Crie uma quadra e tente fazer reservas de diferentes abas.
5. Observe as atualizações em tempo real em todas as janelas.  
   Se uma reserva conflitar com outra, o sistema exibirá uma linha “Negada”.

---

## 🤝 Desenvolvido por
| [**Sofya Andrade**](https://github.com/sofyaandrade) | [**Matheus Ferrari Dos Santos**](https://github.com/matheusferrarimf) |
Disciplina: *Lingaugem de Programação e Paradigmas*  
Professor: *Ademar Perfoll Junior*  
Curso de *Sistemas de Informação — UNIDAVI*  

---
