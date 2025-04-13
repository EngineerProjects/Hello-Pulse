# 📄 Prompt: Help Me Build a Scalable Modular Python AI Service (with strict security filtering)

---

I am building a collaborative brainstorming and project management platform, and I need help **designing and starting the development** of the **Python AI Microservice** part.

Here’s the global architecture of the system:

```
+--------------------------------------------------------------+
|                       Frontend (Web / Mobile)                |
|        (User Interface: Projects, Brainstorming, Agents)     |
+--------------------------------------------------------------+
                               |
                               ▼
+--------------------------------------------------------------+
|                        Backend (Golang)                     |
|  - User Authentication                                       |
|  - Company and Project Management                            |
|  - Permissions and Access Control                            |
|  - Session and Calendar Management                           |
|  - Document Management (metadata)                            |
|  - Orchestration of AI Requests (secure API calls to Python) |
+--------------------------------------------------------------+
                               |
                      (Secure API Calls)
                               |
                               ▼
+--------------------------------------------------------------+
|                    Python AI Microservice                   |
|  - Communicate with LLMs (Ollama, OpenAI)                    |
|  - Perform RAG (Embedding, Search, Generate)                 |
|  - Perform Web Search and Enrich Responses                   |
|  - Manage AI Agents (Creation, Personalization)              |
|  - Serve AI API Endpoints for Golang                         |
+--------------------------------------------------------------+
    |                           |                           |
    ▼                           ▼                           ▼
+-----------+            +-----------------+         +----------------+
| Ollama API|            | OpenAI API       |         | Web Search API |
| (Local/Cloud)          | (External Cloud) |         | (SerpAPI, etc.)|
+-----------+            +-----------------+         +----------------+
                               |
                               ▼
+--------------------------------------------------------------+
|                 Vector Database (Qdrant / ChromaDB)         |
|  - Store document embeddings                                |
|  - Metadata: organization_id, user_id, project_id, visibility|
|  - Enforce security through strict query filtering          |
+--------------------------------------------------------------+
```

---

## 🛠️ Requirements for the Python AI Microservice:

- **FastAPI** (or another performant async framework) for APIs.
- **Modular Design**:
  - Add or remove LLM providers easily (Ollama, OpenAI, HuggingFace, etc.).
  - Add or switch vector databases easily (ChromaDB, Qdrant, others).
  - Each provider/vector database should have its own **independent integration layer**.
- **Secure Authentication**:
  - Use JWTs or API Keys to authenticate internal communication between Golang and Python services.
- **Multi-company Ready**:
  - Handle multi-tenancy properly.
  - Allow companies to either:
    - Use a shared vector database (with strict access filtering),
    - Or request their **own isolated private vector database** for maximum security.

- **Strict Filtering on All Search and Data Access**:
  - When performing vector searches or RAG, **always apply strict filters** on:
    - `organization_id` (company),
    - `user_id` (who uploaded the document),
    - `visibility` (whether the document is public, shared in company, or private).
  - No document should be accessible unless it matches both the user’s access rights and company boundaries.

- **Pluggable Search and RAG**:
  - Modular document retrieval and augmentation.
  - Each company could eventually have custom RAG pipelines.

- **Web Search Module**:
  - Modular web search integration (e.g., SerpAPI) that can be swapped later if needed.

---

## ⚡ Additional Development Goals:

- APIs to implement first:
  - `/generate`: Simple LLM text generation.
  - `/rag_generate`: RAG-based answer generation with strict access filtering.
  - `/search_and_generate`: Web Search + LLM enrichment.
  - `/create_agent`: Create a new AI agent with specific documents and instructions.
  - `/chat_with_agent`: Long conversation with a personalized AI agent.

- **Scalable and Clean Project Structure**:
  - Clear separation between API layer, LLM providers, vector storage modules, search modules, agent management modules.

- **Configuration Management**:
  - Dynamically manage per-company configurations:
    - LLM provider selection,
    - Vector database selection (shared or private),
    - Document visibility policies.

- **Logging and Monitoring**:
  - Good error handling, request logging, and operational metrics.

---

# ❓ Additional Questions:

1. **Is it possible** to fully design the microservice this way so I can easily add/remove new LLM providers or vector databases in the future without major refactoring?
2. **What is the best modular project structure** you recommend to organize this Python code (production-ready)?
3. **Can you help generate a clean project starter template** with folders like:
   - `/api`
   - `/services/llm_providers`
   - `/services/vector_databases`
   - `/services/agents`
   - `/services/websearch`
   - `/config`
   - `/schemas`
   - `/utils`
4. **How should I manage secure dynamic connections** when companies want separate vector databases or have different LLM providers?
5. **What is the best way to enforce strict filtering** across all operations in the service?

---

# ✅ Final Goal:

**Build a modular, scalable, secure, multi-tenant Python AI Microservice** that can evolve easily as I add new AI models, new vector stores, and onboard many different companies, while **strictly protecting data access by organization, user, and visibility rules**.

---

✅ Please use this architecture and requirements to guide me and help me start building the first clean version.

---