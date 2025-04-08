<div align="center">
  <img src="docs/statics/logo.png" alt="Hello Pulse Logo" width="200">
  <h1>Hello Pulse</h1>
  <p>A collaborative brainstorming platform powered by artificial intelligence</p>
</div>

## 📋 Overview

Hello Pulse is a modern collaborative platform that combines team workflows with AI-powered assistance. It provides organizations with powerful tools for project management, event scheduling, and real-time collaboration, all enhanced by artificial intelligence.

## 🚀 Features

- **User & Organization Management**
  - User authentication and session management
  - Organization creation and membership management
  - Secure invitation system

- **Project Management**
  - Hierarchical project structure
  - Team collaboration on projects
  - Task assignment and tracking

- **Event Scheduling**
  - Team event coordination
  - Participant management
  - Importance-based prioritization

- **File Management**
  - Secure file storage
  - Access control for shared resources
  - Version tracking

- **AI Integration** (Coming soon)
  - Natural language processing
  - Retrieval-augmented generation
  - Computer vision capabilities

## 🛠️ Tech Stack

Hello Pulse uses a modern, scalable architecture:

- **Backend**: Golang with Gin framework
- **API**: RESTful with JWT authentication
- **Database**: PostgreSQL with pgvector extension
- **Storage**: MinIO for object storage
- **Infrastructure**: Docker and Docker Compose

## 🏗️ Architecture

The application follows a clean, domain-driven design with well-separated concerns:

- **Models**: Core business entities
- **Repositories**: Data access layer
- **Services**: Business logic
- **API Handlers**: Request handling
- **Middleware**: Cross-cutting concerns

## 📊 Project Structure

```
hello-pulse/
├── backend/           # Go-based backend services
│   ├── cmd/           # Application entry points
│   ├── internal/      # Internal packages
│   │   ├── api/       # API handlers and routes
│   │   ├── models/    # Domain models
│   │   ├── repositories/ # Data access
│   │   └── services/  # Business logic
│   └── pkg/           # Shared utilities
├── ai-services/       # Python-based AI services (coming soon)
├── docs/              # Documentation
├── proxy/             # API gateway configuration
├── scripts/           # Utility scripts
├── volumes/           # Data volumes for containers
├── .github/           # GitHub workflows
├── docker-compose.dev.yml # Development environment
├── docker-compose.yml # Production environment
├── Makefile           # Build and development tasks
└── up.sh, down.sh     # Convenience scripts
```

## 🚦 Getting Started

### Prerequisites

- Docker and Docker Compose
- Go 1.24+ (for development)
- Make (optional, for using the Makefile)

### Quick Start

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/hello-pulse.git
   cd hello-pulse
   ```

2. **Set up environment**
   ```bash
   make setup
   ```
   This will create the necessary directories and configuration files.

3. **Generate a secure JWT secret**
   ```bash
   make gen-jwt
   ```

4. **Start the development environment**
   ```bash
   make dev
   # or
   ./up.sh
   ```

5. **Initialize the database**
   ```bash
   make db-setup
   ```

6. **Access the services**
   - Backend API: http://localhost:5000
   - API via Nginx: http://localhost:7000
   - MinIO Console: http://localhost:8900
   - PGAdmin: http://localhost:8888

### Stopping the Services

```bash
make down
# or
./down.sh
```

## 🧹 Maintenance

Hello Pulse provides several commands for maintaining your development environment:

```bash
# View logs from all services
make logs

# Remove containers and volumes (preserves images)
make clean

# Prune Docker resources
make prune-containers  # Remove only containers
make prune-images      # Remove only images
make prune-volumes     # Remove only volumes
make prune-networks    # Remove only networks
make prune-builders    # Clean builder cache
make prune-all         # Prune everything
```

## 🧪 Testing

```bash
# Run backend tests
make backend-test
```

## 📚 Documentation

For more detailed documentation, see the [docs](./docs) directory.

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 👥 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📧 Contact

For questions or support, please contact [your-email@example.com](mailto:your-email@example.com).