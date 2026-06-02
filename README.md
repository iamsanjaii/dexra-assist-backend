# Dexra Assist

AI-powered customer support assistant with Retrieval-Augmented Generation (RAG), vector search, analytics, and secure admin knowledge management.

---

# Overview

Dexra Assist is a full-stack AI customer support platform that enables organizations to create intelligent chat assistants powered by their own documents and custom knowledge sources.

The system uses a Retrieval-Augmented Generation (RAG) pipeline with semantic vector search to provide grounded, context-aware responses while minimizing hallucinations.

The platform includes:

* AI-powered chatbot
* Admin dashboard
* Knowledge base management
* Custom Q&A management
* AI analytics dashboard
* Token and latency monitoring
* Prompt injection protection
* Session-based conversational memory

---

# Features

## AI Chatbot

* Context-aware conversational assistant
* Retrieval-Augmented Generation (RAG)
* Semantic vector search using ChromaDB
* Session memory support
* Hallucination prevention
* Source attribution for responses
* Prompt injection filtering

---

## Admin Dashboard

* Upload and manage documents
* Add/Edit/Delete custom Q&A pairs
* AI analytics dashboard
* Token usage monitoring
* Estimated AI cost tracking
* Request and latency analytics

---

## Security Features

* JWT authentication
* Protected admin APIs
* Prompt injection defense
* Retrieval grounding
* Security middleware
* Input sanitization
* Rate limiting

---

# Tech Stack

## Frontend

* Next.js
* Tailwind CSS
* shadcn/ui
* Recharts
* Zustand

## Backend

* Golang
* Gin Framework
* MongoDB
* ChromaDB

## AI Stack

* OpenRouter API
* RAG Architecture
* Semantic Search
* Embedding-based Retrieval

---


<img width="1230" height="935" alt="diagram-export" src="https://github.com/user-attachments/assets/2dfdf56e-00f8-45fd-b9c9-f24501f967be" />


<img width="1671" height="1189" alt="diagram-export-2-6-2026-1_58_42-pm" src="https://github.com/user-attachments/assets/87369088-d3d2-4484-a2a5-1a6297f097f3" />















# System Architecture

```text
Frontend (Next.js)
        ↓
Gin Backend API
        ↓
-------------------------
| MongoDB              |
| - sessions           |
| - chat history       |
| - documents          |
| - analytics          |
-------------------------
        ↓
-------------------------
| ChromaDB             |
| - embeddings         |
| - semantic retrieval |
-------------------------
        ↓
OpenRouter LLM API
```

---

# RAG Workflow

1. User submits query
2. Backend generates embeddings
3. ChromaDB retrieves relevant chunks
4. Context + memory are injected into prompt
5. LLM generates grounded response
6. Sources and analytics are returned

---

# AI Analytics

Dexra Assist includes observability and monitoring features:

* Token usage tracking
* Request analytics
* Average latency monitoring
* Estimated AI cost tracking
* Usage visualization dashboards

---

# Hallucination Prevention

The assistant is designed to minimize hallucinations by:

* enforcing retrieval grounding
* applying similarity thresholds
* refusing unsupported queries
* restricting prompt injection attempts

If sufficient context is unavailable, the assistant responds safely instead of generating misleading answers.

---

# Screenshots

## Admin Dashboard

* Knowledge base management
* AI analytics dashboard
* Q&A management

## Client Chat Interface

* Conversational AI assistant
* Source citations
* Session memory

---

# Setup Instructions

## Frontend

```bash
cd frontend
npm install
npm run dev
```

---

## Backend

```bash
cd backend
go mod tidy
go run cmd/main.go
```

---

# Environment Variables

## Backend

```env
MONGO_URI=
OPENROUTER_API_KEY=
JWT_SECRET=
CHROMADB_URL=
```

---

# Deployment

## Frontend

* Vercel

## Backend

* VPS + Nginx Reverse Proxy

## Vector Database

* Self-hosted ChromaDB

---

# Future Improvements

* Streaming AI responses
* Multi-user RBAC
* Redis caching
* WebSocket support
* Multi-tenant architecture
* File preview system

---

# Demo

## Live Demo

https://dexra.sanjai.app

## Admin Panel

https://admin.sanjai.app



## Demo Video

(Add video link)

---

# Author

Built by Sanjai Kumar.
