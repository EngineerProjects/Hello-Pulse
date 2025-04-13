"""
Configuration management for the Hello Pulse AI Microservice
"""

import os
from typing import List, Optional, Dict, Any, Union
from pydantic import Field
from pydantic_settings import BaseSettings

class Settings(BaseSettings):
    # Project info
    PROJECT_NAME: str = "Hello Pulse AI Microservice"
    PROJECT_DESCRIPTION: str = "AI capabilities for Hello Pulse platform"
    PROJECT_VERSION: str = "0.1.0"
    
    # Server settings
    HOST: str = "0.0.0.0"
    PORT: int = 8080
    DEBUG: bool = False
    
    # Security settings
    SECRET_KEY: str = Field(..., env="SECRET_KEY")
    AUTH_ALGORITHM: str = "HS256"
    AUTH_TOKEN_EXPIRE_MINUTES: int = 60
    
    # CORS settings
    CORS_ORIGINS: List[str] = ["*"]
    
    # Default LLM provider
    DEFAULT_LLM_PROVIDER: str = "openai"
    
    # Default vector database
    DEFAULT_VECTOR_DB: str = "chroma"
    SHARED_VECTOR_DB: bool = True  # Whether to use a shared vector DB or per-company
    
    # OpenAI settings
    OPENAI_API_KEY: Optional[str] = Field(None, env="OPENAI_API_KEY")
    OPENAI_EMBEDDING_MODEL: str = "text-embedding-3-small"
    OPENAI_DEFAULT_MODEL: str = "gpt-3.5-turbo"
    
    # Ollama settings
    OLLAMA_API_BASE: str = "http://localhost:11434/api"
    OLLAMA_DEFAULT_MODEL: str = "llama3"
    
    # ChromaDB settings
    CHROMA_HOST: str = "localhost"
    CHROMA_PORT: int = 8000
    CHROMA_COLLECTION_NAME: str = "hello_pulse"
    
    # Qdrant settings
    QDRANT_HOST: str = "localhost"
    QDRANT_PORT: int = 6333
    QDRANT_COLLECTION_NAME: str = "hello_pulse"
    
    # Web search settings
    SERPAPI_API_KEY: Optional[str] = Field(None, env="SERPAPI_API_KEY")
    
    # Logging
    LOG_LEVEL: str = "INFO"
    
    # Company configurations cache (would be replaced with DB or external config in production)
    # This is a placeholder for dynamic configuration per company
    COMPANY_CONFIGS: Dict[str, Dict[str, Any]] = {}

    class Config:
        env_file = ".env"
        env_file_encoding = "utf-8"
        case_sensitive = True

    def get_company_config(self, organization_id: str) -> Dict[str, Any]:
        """Get company-specific configuration"""
        if organization_id not in self.COMPANY_CONFIGS:
            # Return default configuration if company specific one doesn't exist
            return {
                "llm_provider": self.DEFAULT_LLM_PROVIDER,
                "vector_db": self.DEFAULT_VECTOR_DB,
                "shared_vector_db": self.SHARED_VECTOR_DB,
                "openai_model": self.OPENAI_DEFAULT_MODEL,
                "ollama_model": self.OLLAMA_DEFAULT_MODEL,
            }
        return self.COMPANY_CONFIGS[organization_id]
    
    def update_company_config(self, organization_id: str, config: Dict[str, Any]) -> None:
        """Update company-specific configuration"""
        if organization_id not in self.COMPANY_CONFIGS:
            self.COMPANY_CONFIGS[organization_id] = {}
        self.COMPANY_CONFIGS[organization_id].update(config)

    def get_llm_provider_for_company(self, organization_id: str) -> str:
        """Get LLM provider for specific company"""
        company_config = self.get_company_config(organization_id)
        return company_config.get("llm_provider", self.DEFAULT_LLM_PROVIDER)
    
    def get_vector_db_for_company(self, organization_id: str) -> str:
        """Get vector database for specific company"""
        company_config = self.get_company_config(organization_id)
        return company_config.get("vector_db", self.DEFAULT_VECTOR_DB)

# Create settings instance
settings = Settings()