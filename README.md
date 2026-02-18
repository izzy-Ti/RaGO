# 🐹 Go-RAG: Extreme AI Intelligence Engine

<p align="center">
  <img src="https://media4.giphy.com/media/v1.Y2lkPTc5MGI3NjExeDFjMTljZTZ2MGVpa2N5bDN4YmhpY3hyaHJ2N2p3MzY4ODdmb2tjNCZlcD12MV9pbnRlcm5hbF9naWZfYnlfaWQmY3Q9Zw/tiOnBLtDu8iQVOvIgV/giphy.gif" width="600" alt="Go-RAG Animation">
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
DATABASE_URL=postgresql:supabase_url
JWT_KEY=any_secret_key
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
<img src="https://media2.giphy.com/media/v1.Y2lkPTc5MGI3NjExdHhsa2V6bXFjNXBzc2EyeWEwaDJ6OXJibWR2NW5lM2NrMGcxbXN1NCZlcD12MV9pbnRlcm5hbF9naWZfYnlfaWQmY3Q9Zw/YPIrsRqqO7oB2/giphy.gif" width="300">
</p>




