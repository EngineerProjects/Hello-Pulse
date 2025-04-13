"""
Pydantic models for API request validation
"""

from typing import Dict, List, Optional, Any, Union
from pydantic import BaseModel, Field, validator

class GenerateRequest(BaseModel):
    """Request model for text generation"""
    prompt: str = Field(..., description="The prompt to generate from")
    system_message: Optional[str] = Field(None, description="Optional system message for chat models")
    temperature: float = Field(0.7, description="Controls randomness (0.0 = deterministic, 1.0 = random)")
    max_tokens: Optional[int] = Field(None, description="Maximum number of tokens to generate")
    stop_sequences: Optional[List[str]] = Field(None, description="List of sequences where generation should stop")
    stream: bool = Field(False, description="Whether to stream the response")
    model_params: Optional[Dict[str, Any]] = Field(None, description="Additional model-specific parameters")


class ChatMessage(BaseModel):
    """Chat message model"""
    role: str = Field(..., description="Message role (system, user, or assistant)")
    content: str = Field(..., description="Message content")


class ChatRequest(BaseModel):
    """Request model for chat API"""
    messages: List[ChatMessage] = Field(..., description="List of messages in the conversation")
    temperature: float = Field(0.7, description="Controls randomness (0.0 = deterministic, 1.0 = random)")
    max_tokens: Optional[int] = Field(None, description="Maximum number of tokens to generate")
    stop_sequences: Optional[List[str]] = Field(None, description="List of sequences where generation should stop")
    stream: bool = Field(False, description="Whether to stream the response")
    model_params: Optional[Dict[str, Any]] = Field(None, description="Additional model-specific parameters")


class DocumentMetadata(BaseModel):
    """Document metadata model"""
    visibility: Optional[str] = Field("private", description="Document visibility (private, shared, public)")
    project_id: Optional[str] = Field(None, description="Project ID if applicable")
    tags: Optional[List[str]] = Field(None, description="List of tags")
    source: Optional[str] = Field(None, description="Document source")
    
    # Allow additional properties
    class Config:
        extra = "allow"


class Document(BaseModel):
    """Document model"""
    text: str = Field(..., description="Document text")
    metadata: DocumentMetadata = Field(DocumentMetadata(), description="Document metadata")


class AddDocumentsRequest(BaseModel):
    """Request model for adding documents"""
    documents: List[Document] = Field(..., description="List of documents to add")


class SearchRequest(BaseModel):
    """Request model for vector search"""
    query: str = Field(..., description="Query text")
    filter: Optional[Dict[str, Any]] = Field({}, description="Filter criteria")
    limit: int = Field(10, description="Maximum number of results")
    include_embeddings: bool = Field(False, description="Whether to include embeddings in results")


class RAGRequest(BaseModel):
    """Request model for RAG generation"""
    query: str = Field(..., description="Query text")
    filter: Optional[Dict[str, Any]] = Field({}, description="Filter criteria for document retrieval")
    system_message: Optional[str] = Field(None, description="Optional system message for chat models")
    temperature: float = Field(0.7, description="Controls randomness (0.0 = deterministic, 1.0 = random)")
    max_tokens: Optional[int] = Field(None, description="Maximum number of tokens to generate")
    num_documents: int = Field(5, description="Number of documents to retrieve")
    stream: bool = Field(False, description="Whether to stream the response")


class WebSearchRequest(BaseModel):
    """Request model for web search"""
    query: str = Field(..., description="Search query")
    num_results: int = Field(5, description="Number of search results to return")
    include_snippets: bool = Field(True, description="Whether to include result snippets")


class WebSearchAndGenerateRequest(BaseModel):
    """Request model for web search + generation"""
    query: str = Field(..., description="Search query")
    num_results: int = Field(5, description="Number of search results to return")
    system_message: Optional[str] = Field(None, description="Optional system message for chat models")
    temperature: float = Field(0.7, description="Controls randomness (0.0 = deterministic, 1.0 = random)")
    max_tokens: Optional[int] = Field(None, description="Maximum number of tokens to generate")
    stream: bool = Field(False, description="Whether to stream the response")


class AgentCreateRequest(BaseModel):
    """Request model for creating an agent"""
    name: str = Field(..., description="Agent name")
    description: Optional[str] = Field(None, description="Agent description")
    instructions: str = Field(..., description="Agent instructions")
    document_ids: Optional[List[str]] = Field(None, description="Document IDs to associate with the agent")
    document_filter: Optional[Dict[str, Any]] = Field(None, description="Filter criteria for documents")
    metadata: Optional[Dict[str, Any]] = Field(None, description="Additional agent metadata")


class AgentChatRequest(BaseModel):
    """Request model for chatting with an agent"""
    agent_id: str = Field(..., description="Agent ID")
    messages: List[ChatMessage] = Field(..., description="List of messages in the conversation")
    temperature: float = Field(0.7, description="Controls randomness (0.0 = deterministic, 1.0 = random)")
    max_tokens: Optional[int] = Field(None, description="Maximum number of tokens to generate")
    stream: bool = Field(False, description="Whether to stream the response")
    model_params: Optional[Dict[str, Any]] = Field(None, description="Additional model-specific parameters")


class AgentListRequest(BaseModel):
    """Request model for listing agents"""
    filter: Optional[Dict[str, Any]] = Field({}, description="Filter criteria")
    limit: int = Field(100, description="Maximum number of results")
    offset: int = Field(0, description="Pagination offset")