"""
Pydantic models for API response validation
"""

from typing import Dict, List, Optional, Any, Union
from pydantic import BaseModel, Field
from datetime import datetime


class GenerateResponse(BaseModel):
    """Response model for text generation"""
    text: str = Field(..., description="Generated text")
    model: str = Field(..., description="Model used for generation")
    provider: str = Field(..., description="Provider used for generation")
    finish_reason: Optional[str] = Field(None, description="Reason for finishing generation")
    usage: Optional[Dict[str, int]] = Field(None, description="Token usage statistics if available")


class StreamChunk(BaseModel):
    """Streaming chunk model"""
    text: str = Field(..., description="Chunk of generated text")
    finish_reason: Optional[str] = Field(None, description="Reason for finishing generation")


class ChatResponse(BaseModel):
    """Response model for chat API"""
    message: str = Field(..., description="Generated message")
    model: str = Field(..., description="Model used for generation")
    provider: str = Field(..., description="Provider used for generation")
    finish_reason: Optional[str] = Field(None, description="Reason for finishing generation")
    usage: Optional[Dict[str, int]] = Field(None, description="Token usage statistics if available")


class DocumentResponse(BaseModel):
    """Response model for document operations"""
    id: str = Field(..., description="Document ID")
    text: str = Field(..., description="Document text")
    metadata: Dict[str, Any] = Field(..., description="Document metadata")
    embedding: Optional[List[float]] = Field(None, description="Document embedding if requested")
    score: Optional[float] = Field(None, description="Similarity score if from search results")


class DocumentListResponse(BaseModel):
    """Response model for document listing"""
    documents: List[DocumentResponse] = Field(..., description="List of documents")
    count: int = Field(..., description="Total number of documents matching the query")
    limit: int = Field(..., description="Limit used for the query")
    offset: int = Field(..., description="Offset used for the query")


class SearchResponse(BaseModel):
    """Response model for vector search"""
    results: List[DocumentResponse] = Field(..., description="Search results")
    query: str = Field(..., description="Original query")
    count: int = Field(..., description="Number of results returned")
    total: Optional[int] = Field(None, description="Total number of matching documents if available")


class RAGResponse(BaseModel):
    """Response model for RAG generation"""
    answer: str = Field(..., description="Generated answer")
    documents: List[DocumentResponse] = Field(..., description="Retrieved documents")
    model: str = Field(..., description="Model used for generation")
    provider: str = Field(..., description="Provider used for generation")
    finish_reason: Optional[str] = Field(None, description="Reason for finishing generation")
    usage: Optional[Dict[str, int]] = Field(None, description="Token usage statistics if available")


class WebSearchResult(BaseModel):
    """Web search result model"""
    title: str = Field(..., description="Result title")
    url: str = Field(..., description="Result URL")
    snippet: Optional[str] = Field(None, description="Result snippet")
    source: str = Field(..., description="Source name (e.g., 'Google', 'Bing')")


class WebSearchResponse(BaseModel):
    """Response model for web search"""
    results: List[WebSearchResult] = Field(..., description="Search results")
    query: str = Field(..., description="Original query")
    count: int = Field(..., description="Number of results returned")


class WebSearchAndGenerateResponse(BaseModel):
    """Response model for web search + generation"""
    answer: str = Field(..., description="Generated answer based on search results")
    search_results: List[WebSearchResult] = Field(..., description="Search results used for generation")
    model: str = Field(..., description="Model used for generation")
    provider: str = Field(..., description="Provider used for generation")
    finish_reason: Optional[str] = Field(None, description="Reason for finishing generation")
    usage: Optional[Dict[str, int]] = Field(None, description="Token usage statistics if available")


class AgentResponse(BaseModel):
    """Response model for agent info"""
    id: str = Field(..., description="Agent ID")
    name: str = Field(..., description="Agent name")
    description: Optional[str] = Field(None, description="Agent description")
    instructions: str = Field(..., description="Agent instructions")
    created_at: datetime = Field(..., description="Creation timestamp")
    updated_at: Optional[datetime] = Field(None, description="Last update timestamp")
    metadata: Optional[Dict[str, Any]] = Field(None, description="Additional agent metadata")
    document_count: Optional[int] = Field(None, description="Number of documents associated with the agent")


class AgentListResponse(BaseModel):
    """Response model for agent listing"""
    agents: List[AgentResponse] = Field(..., description="List of agents")
    count: int = Field(..., description="Total number of agents matching the query")
    limit: int = Field(..., description="Limit used for the query")
    offset: int = Field(..., description="Offset used for the query")


class AgentChatResponse(BaseModel):
    """Response model for chatting with an agent"""
    message: str = Field(..., description="Generated message")
    agent_id: str = Field(..., description="Agent ID")
    agent_name: str = Field(..., description="Agent name")
    model: str = Field(..., description="Model used for generation")
    provider: str = Field(..., description="Provider used for generation")
    document_ids: Optional[List[str]] = Field(None, description="IDs of documents used for this response")
    finish_reason: Optional[str] = Field(None, description="Reason for finishing generation")
    usage: Optional[Dict[str, int]] = Field(None, description="Token usage statistics if available")


class ErrorResponse(BaseModel):
    """Error response model"""
    error: str = Field(..., description="Error message")
    details: Optional[Dict[str, Any]] = Field(None, description="Additional error details")
    request_id: Optional[str] = Field(None, description="Request ID for troubleshooting")