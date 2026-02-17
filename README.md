# 🐹 Go-RAG: Extreme AI Intelligence Engine

<p align="center">
  <img src="https://media.giphy.com/media/v1.Y2lkPTc5MGI3NjExNHJueGZ4bm96ZzRneHBleGZ4eGZ4eGZ4eGZ4eGZ4eGZ4eGZ4JmVwPXYxX2ludGVybmFsX2dpZl9ieV9pZCBmcm9tPXNlYXJjaCZjdD1n/3o7TKSjP3u8R4S2YwM/giphy.gif" width="600" alt="Go-RAG Animation">
</p>

---

## ⚡ The Project Overview
This is a **High-Performance RAG (Retrieval-Augmented Generation)** system built in Golang. It doesn't just chat; it *thinks* by connecting your PDFs and live URLs to a massive vector brain.

### 🛡️ Why This Stack?
* **Golang:** For lightning-fast concurrent processing of large PDFs.
* **Astra DB (DataStax):** Serverless Vector DB for global-scale AI retrieval.
* **Supabase/PG:** The backbone for user accounts, session history, and relational data.
* **Groq:** Sub-second LLM inference (Llama-3/Mixtral).
* **Voyage AI:** State-of-the-art text embeddings.

---

## 🚀 Core Features

| Feature | Description | Status |
| :--- | :--- | :--- |
| **🔐 Full Auth** | Google OAuth2 + JWT Integration | ✅ Done |
| **📄 PDF Ingestion** | Extract, Chunk, and Embed complex PDF documents | ✅ Done |
| **🌐 URL Scraper** | Live web scraping for real-time AI context | ✅ Done |
| **🧠 Vector Search** | Semantic retrieval using Astra DB & Voyage AI | ✅ Done |
| **💬 Chat History** | Persistent sessions stored in Supabase/PostgreSQL | ✅ Done |
| **👮 Admin Panel** | Full control over users, data, and embedding indices | ✅ Done |

---

## 🛠️ System Architecture



1.  **Ingestion:** Files/Links are processed via Go Goroutines.
2.  **Vectorization:** Text chunks $\rightarrow$ Voyage AI $\rightarrow$ High-dimensional Vectors.
3.  **Storage:** Vectors live in **Astra DB**; Meta-data lives in **Supabase**.
4.  **Retrieval:** User Query $\rightarrow$ Vector Search $\rightarrow$ Context Injection $\rightarrow$ Groq Answer.

---

## 📦 Environment Setup

Create a `.env` file in your root folder:

```env
# Database & Auth
DATABASE_URL=postgresql://postgres.ryaspdobmptunhmrkslr:8uh6zM07UhuX21Vv@aws-1-eu-west-1.pooler.supabase.com:5432/postgres?sslmode=require
JWT_KEY="izzyGO"
PORT=8080
APP_ENV=production

# AI Keys (Astra & Inference)
ASTRA=your_astra_token
ASTRA_END_POINT=your_endpoint_here
VOYAGE_API_KEY=your_voyage_key
GROQ_API_KEY=your_groq_key

# Integrations
API_BERVO=your_brevo_key
EMAIL=your_admin_email
PASSWORD=your_app_password
GOOGLE_CLIENT_ID=your_google_id
```

🚦 Getting Started

1. Installation

go mod download

2. Run the Engine

go run ./cmd/main.go

<p align="center">
<img src="https://www.google.com/search?q=https://media.giphy.com/media/v1.Y2lkPTc5MGI3NjExNHJueGZ4bm96ZzRneHBleGZ4eGZ4eGZ4eGZ4eGZ4eGZ4eGZ4JmVwPXYxX2ludGVybmFsX2dpZl9ieV9pZCBmcm9tPXNlYXJjaCZjdD1n/L3X9GvptxDcuY/giphy.gif" width="300">
</p>




