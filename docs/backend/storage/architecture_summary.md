# Hello Pulse - Storage Architecture

## Overview

The storage architecture of Hello Pulse is designed to be modular, provider-agnostic, and easily extensible. This approach allows for flexibility in deployment environments and future scaling.

## Components

### Core Abstractions

- **Storage Provider Interface** (`pkg/storage/storage.go`)
  - Defines the contract that all storage providers must implement
  - Includes operations for uploading, downloading, and managing files

- **Storage Factory** (`pkg/storage/factory.go`)
  - Creates the appropriate storage provider based on configuration
  - Single point for adding new provider implementations

### Implementations

- **MinIO Provider** (`pkg/storage/minio_provider.go`)
  - For development and smaller deployments
  - Full implementation of the Provider interface for MinIO

- **S3 Provider Example** (`pkg/storage/s3_provider.go.example`)
  - Reference implementation for AWS S3
  - Can be used as a template for production deployment

- **Provider Placeholders**
  - Sketches for other cloud providers (Azure, GCS)
  - Define the structure for future implementations

### Integration

- **File Service** (`internal/services/file/file_service.go`)
  - Business logic for file management
  - Uses the storage provider interface for all storage operations
  - Maintains database records for files

- **Configuration** (`pkg/config/storage_config.go`)
  - Loads provider-specific configuration from environment variables
  - Centralizes configuration for all storage providers

## Design Principles

1. **Separation of Concerns**
   - Business logic is independent of storage implementation
   - Storage operations are abstracted through interfaces

2. **Provider Agnosticism**
   - Code never depends on specific provider details
   - All provider-specific logic is contained in dedicated implementations

3. **Configuration Over Code**
   - Switching providers requires only configuration changes
   - No code changes needed to use different providers

4. **Graceful Degradation**
   - System handles storage provider failures gracefully
   - Clear error messages when storage is unavailable

## Usage Flow

1. Configuration is loaded from environment variables
2. The appropriate storage provider is created via the factory
3. The storage provider is initialized and tested
4. The file service uses the provider for all storage operations
5. Business logic operates on the file service without knowing the underlying provider

## Extension Points

To add a new storage provider:

1. Create a new implementation file (`pkg/storage/new_provider.go`)
2. Implement the Provider interface
3. Add the provider to the factory function
4. Update the configuration to support the new provider

No changes are required to business logic or other parts of the system.