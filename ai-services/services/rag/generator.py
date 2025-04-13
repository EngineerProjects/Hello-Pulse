"""
Response generation service for RAG
"""

from typing import Dict, List, Any, Optional, AsyncIterator, Tuple
import asyncio

from core.logging import logger, LoggingContext
from services.llm_providers import get_llm_provider
from services.rag.retriever import retrieve_relevant_documents


async def generate_answer(
    query: str,
    organization_id: str,
    user_id: str,
    filter_dict: Optional[Dict[str, Any]] = None,
    system_message: Optional[str] = None,
    temperature: float = 0.7,
    max_tokens: Optional[int] = None,
    num_documents: int = 5,
    **kwargs: Any
) -> Tuple[str, List[Dict[str, Any]]]:
    """
    Generate an answer to a query using retrieved documents
    
    Args:
        query: User query
        organization_id: Organization ID
        user_id: User ID
        filter_dict: Filters for document retrieval
        system_message: Optional system message
        temperature: Temperature for generation
        max_tokens: Maximum tokens to generate
        num_documents: Number of documents to retrieve
        **kwargs: Additional parameters
        
    Returns:
        Tuple of (generated_answer, documents_used)
    """
    with LoggingContext(organization_id=organization_id, user_id=user_id):
        try:
            # Retrieve relevant documents
            documents = await retrieve_relevant_documents(
                query=query,
                organization_id=organization_id,
                user_id=user_id,
                filter_dict=filter_dict,
                limit=num_documents
            )
            
            if not documents:
                # No documents found, generate a response without context
                llm_provider = get_llm_provider(organization_id)
                
                default_system = system_message or (
                    "You are a helpful assistant. If you don't know the answer to a question, "
                    "just say that you don't have enough information to provide a reliable answer."
                )
                
                # Generate response
                response = await llm_provider.generate(
                    prompt=query,
                    system_message=default_system,
                    temperature=temperature,
                    max_tokens=max_tokens,
                    **kwargs
                )
                
                logger.info(
                    f"Generated answer without context (no relevant documents found)",
                    extra={
                        "query": query,
                        "organization_id": organization_id,
                        "user_id": user_id,
                    }
                )
                
                return response, []
            
            # Format documents for the prompt
            context = format_documents_for_context(documents)
            
            # Get LLM provider
            llm_provider = get_llm_provider(organization_id)
            
            # Create system message with RAG instructions
            if system_message:
                rag_system_message = system_message
            else:
                rag_system_message = (
                    "You are a helpful assistant with access to a knowledge base. "
                    "Answer the user's question based on the provided context. "
                    "If the context doesn't contain the necessary information to answer the question, "
                    "just say that you don't have enough information to provide a reliable answer. "
                    "Don't make up information that's not in the context."
                )
            
            # Create prompt with context
            rag_prompt = f"""I need information about the following question:
{query}

Here is the relevant context from our knowledge base:

{context}

Based on this context, please provide a comprehensive answer to the question."""
            
            # Generate response
            response = await llm_provider.generate(
                prompt=rag_prompt,
                system_message=rag_system_message,
                temperature=temperature,
                max_tokens=max_tokens,
                **kwargs
            )
            
            logger.info(
                f"Generated RAG answer",
                extra={
                    "query": query,
                    "organization_id": organization_id,
                    "user_id": user_id,
                    "document_count": len(documents)
                }
            )
            
            return response, documents
        
        except Exception as e:
            logger.error(
                f"Error generating answer: {str(e)}",
                extra={
                    "query": query,
                    "organization_id": organization_id,
                    "user_id": user_id
                }
            )
            raise


async def generate_answer_stream(
    query: str,
    organization_id: str,
    user_id: str,
    filter_dict: Optional[Dict[str, Any]] = None,
    system_message: Optional[str] = None,
    temperature: float = 0.7,
    max_tokens: Optional[int] = None,
    num_documents: int = 5,
    **kwargs: Any
) -> Tuple[AsyncIterator[str], List[Dict[str, Any]]]:
    """
    Stream an answer to a query using retrieved documents
    
    Args:
        query: User query
        organization_id: Organization ID
        user_id: User ID
        filter_dict: Filters for document retrieval
        system_message: Optional system message
        temperature: Temperature for generation
        max_tokens: Maximum tokens to generate
        num_documents: Number of documents to retrieve
        **kwargs: Additional parameters
        
    Returns:
        Tuple of (stream_iterator, documents_used)
    """
    with LoggingContext(organization_id=organization_id, user_id=user_id):
        try:
            # Retrieve relevant documents
            documents = await retrieve_relevant_documents(
                query=query,
                organization_id=organization_id,
                user_id=user_id,
                filter_dict=filter_dict,
                limit=num_documents
            )
            
            # Get LLM provider
            llm_provider = get_llm_provider(organization_id)
            
            if not documents:
                # No documents found, generate a response without context
                default_system = system_message or (
                    "You are a helpful assistant. If you don't know the answer to a question, "
                    "just say that you don't have enough information to provide a reliable answer."
                )
                
                # Create streaming response
                stream = llm_provider.generate_stream(
                    prompt=query,
                    system_message=default_system,
                    temperature=temperature,
                    max_tokens=max_tokens,
                    **kwargs
                )
                
                logger.info(
                    f"Streaming answer without context (no relevant documents found)",
                    extra={
                        "query": query,
                        "organization_id": organization_id,
                        "user_id": user_id,
                    }
                )
                
                return stream, []
            
            # Format documents for the prompt
            context = format_documents_for_context(documents)
            
            # Create system message with RAG instructions
            if system_message:
                rag_system_message = system_message
            else:
                rag_system_message = (
                    "You are a helpful assistant with access to a knowledge base. "
                    "Answer the user's question based on the provided context. "
                    "If the context doesn't contain the necessary information to answer the question, "
                    "just say that you don't have enough information to provide a reliable answer. "
                    "Don't make up information that's not in the context."
                )
            
            # Create prompt with context
            rag_prompt = f"""I need information about the following question:
{query}

Here is the relevant context from our knowledge base:

{context}

Based on this context, please provide a comprehensive answer to the question."""
            
            # Create streaming response
            stream = llm_provider.generate_stream(
                prompt=rag_prompt,
                system_message=rag_system_message,
                temperature=temperature,
                max_tokens=max_tokens,
                **kwargs
            )
            
            logger.info(
                f"Streaming RAG answer",
                extra={
                    "query": query,
                    "organization_id": organization_id,
                    "user_id": user_id,
                    "document_count": len(documents)
                }
            )
            
            return stream, documents
        
        except Exception as e:
            logger.error(
                f"Error streaming answer: {str(e)}",
                extra={
                    "query": query,
                    "organization_id": organization_id,
                    "user_id": user_id
                }
            )
            raise


def format_documents_for_context(documents: List[Dict[str, Any]]) -> str:
    """
    Format retrieved documents into a single context string
    
    Args:
        documents: List of document dictionaries with text and metadata
        
    Returns:
        Formatted context string
    """
    context_parts = []
    
    for i, doc in enumerate(documents):
        # Extract document information
        text = doc["text"]
        metadata = doc["metadata"]
        
        # Add document metadata if available
        source_info = []
        if "title" in metadata:
            source_info.append(f"Title: {metadata['title']}")
        if "source" in metadata:
            source_info.append(f"Source: {metadata['source']}")
        if "url" in metadata:
            source_info.append(f"URL: {metadata['url']}")
        if "date" in metadata:
            source_info.append(f"Date: {metadata['date']}")
        if "author" in metadata:
            source_info.append(f"Author: {metadata['author']}")
        
        # Format document with metadata
        if source_info:
            source_str = " | ".join(source_info)
            context_parts.append(f"Document {i+1} ({source_str}):\n{text}\n")
        else:
            context_parts.append(f"Document {i+1}:\n{text}\n")
    
    return "\n".join(context_parts)