# Hello Pulse AI Microservice

A modular, scalable, secure Python AI Microservice for the Hello Pulse collaborative platform.

## 🚀 Features

- **Text Generation**: Simple LLM-based text generation with both OpenAI and Ollama support
- **RAG Generation**: Retrieval-Augmented Generation with strict access filtering
- **Web Search**: Search and enhanced generation with web search results
- **AI Agents**: Create and chat with personalized AI agents
- **Multi-tenant Architecture**: Secure isolation between organizations
- **Modular Design**: Easily swap or add new LLM providers and vector databases

## 🔧 Supported Integrations

### LLM Providers
- OpenAI (GPT-3.5, GPT-4, etc.)
- Ollama (for local and open-source LLMs)

### Vector Databases
- ChromaDB
- Qdrant

### Web Search
- SerpAPI

## 🛠️ Project Structure

```
ai-service/
├── README.md
├── requirements.txt
├── Dockerfile
├── docker-compose.yml
├── .env.example
├── main.py                  # FastAPI application entry point
├── api/
│   ├── routes/
│   │   ├── generate.py      # LLM generation endpoints
│   │   ├── rag.py           # RAG-based endpoints
│   │   ├── search.py        # Web search endpoints
│   │   └── agents.py        # AI agent endpoints
│   ├── dependencies.py      # API dependencies
│   └── middleware.py        # API middleware
├── core/
│   ├── config.py            # Configuration management
│   ├── security.py          # Authentication and authorization
│   └── logging.py           # Logging setup
├── schemas/
│   ├── requests.py          # Request models
│   └── responses.py         # Response models
├── services/
│   ├── llm_providers/
│   │   ├── base.py          # Base LLM provider interface
│   │   ├── ollama.py        # Ollama integration
│   │   └── openai.py        # OpenAI integration
│   ├── vector_databases/
│   │   ├── base.py          # Base vector DB interface
│   │   ├── chroma.py        # ChromaDB integration
│   │   └── qdrant.py        # Qdrant integration
│   ├── agents/
│   │   ├── base.py          # Base agent interface
│   │   ├── manager.py       # Agent management
│   │   └── personalization.py # Agent personalization
│   ├── websearch/
│   │   ├── base.py          # Base search interface
│   │   └── serpapi.py       # SerpAPI integration
│   └── rag/
│       ├── retriever.py     # Document retrieval
│       └── generator.py     # Response generation
└── utils/
    ├── filtering.py         # Access control filters
    ├── embeddings.py        # Embedding utilities
    └── streaming.py         # Streaming response utilities
```

## 🚦 Getting Started

### Prerequisites

- Python 3.9+
- Docker and Docker Compose (for containerized deployment)
- API keys for OpenAI and SerpAPI (optional)

### Installation

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd hello-pulse-ai-service
   ```

2. **Create and configure environment variables:**
   ```bash
   cp .env.example .env
   # Edit .env with your API keys and configuration
   ```

3. **Run with Docker Compose:**
   ```bash
   docker-compose up -d
   ```

4. **Or install and run locally:**
   ```bash
   pip install -r requirements.txt
   uvicorn main:app --reload
   ```

5. **Access the API documentation:**
   - Swagger UI: http://localhost:8080/docs
   - ReDoc: http://localhost:8080/redoc

## 🔒 Security

This microservice is designed with security in mind:

- **JWT Authentication**: Secure communication with the Go backend
- **Strict Access Filtering**: Documents are filtered by organization, user, and visibility
- **Multi-tenancy**: Organizations can have shared or isolated vector databases

## 📚 API Endpoints

### Text Generation
- `POST /api/generate`: Generate text from a prompt
- `POST /api/chat`: Chat with the AI

### RAG
- `POST /api/rag/generate`: Generate an answer using RAG
- `POST /api/rag/search`: Search for documents
- `POST /api/rag/documents`: Add documents
- `GET /api/rag/documents/{document_id}`: Get document
- `DELETE /api/rag/documents/{document_id}`: Delete document

### Web Search
- `POST /api/search`: Search the web
- `POST /api/search_and_generate`: Search and generate an answer

### Agents
- `POST /api/agents`: Create an agent
- `GET /api/agents/{agent_id}`: Get agent
- `PUT /api/agents/{agent_id}`: Update agent
- `DELETE /api/agents/{agent_id}`: Delete agent
- `POST /api/agents/{agent_id}/chat`: Chat with agent

## 🔧 Configuration

The service is highly configurable through environment variables. See `.env.example` for all options.

## 📖 Adding New Integrations

The modular design makes it easy to add new integrations:

### Adding a New LLM Provider
1. Create a new file in `services/llm_providers/`
2. Implement the `BaseLLMProvider` interface
3. Register the provider in `services/llm_providers/__init__.py`

### Adding a New Vector Database
1. Create a new file in `services/vector_databases/`
2. Implement the `BaseVectorDatabase` interface
3. Register the database in `services/vector_databases/__init__.py`

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.