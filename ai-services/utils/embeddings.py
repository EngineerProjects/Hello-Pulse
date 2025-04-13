"""
Utilities for handling text embeddings
"""

import re
from typing import List, Dict, Any, Optional, Tuple
import numpy as np

from core.logging import logger
from services.llm_providers import get_llm_provider


async def get_embeddings(
    texts: List[str],
    organization_id: str
) -> List[List[float]]:
    """
    Get embeddings for a list of texts
    
    Args:
        texts: List of text strings to embed
        organization_id: Organization ID for LLM provider selection
    
    Returns:
        List of embedding vectors
    """
    if not texts:
        return []
    
    # Get the LLM provider for this organization
    llm_provider = get_llm_provider(organization_id)
    
    try:
        # Generate embeddings
        embeddings = await llm_provider.get_embeddings(texts)
        return embeddings
    except Exception as e:
        logger.error(
            f"Error generating embeddings: {str(e)}",
            extra={
                "organization_id": organization_id,
                "text_count": len(texts),
                "provider": llm_provider.get_provider_name()
            }
        )
        raise


def chunk_text(
    text: str,
    chunk_size: int = 1000,
    chunk_overlap: int = 200
) -> List[str]:
    """
    Split text into overlapping chunks
    
    Args:
        text: Input text
        chunk_size: Maximum chunk size in characters
        chunk_overlap: Overlap between chunks in characters
    
    Returns:
        List of text chunks
    """
    if len(text) <= chunk_size:
        return [text]
    
    # Simple paragraph-based chunking
    paragraphs = re.split(r'\n\s*\n', text)
    
    chunks = []
    current_chunk = ""
    
    for paragraph in paragraphs:
        # If adding this paragraph would exceed chunk size, start a new chunk
        if len(current_chunk) + len(paragraph) > chunk_size:
            chunks.append(current_chunk.strip())
            # Keep some overlap from the previous chunk
            if chunk_overlap > 0 and len(current_chunk) > chunk_overlap:
                # Try to find a paragraph break for overlap
                overlap_text = current_chunk[-chunk_overlap:]
                # Find last paragraph break in the overlap text
                last_break = overlap_text.rfind("\n\n")
                if last_break != -1:
                    current_chunk = current_chunk[-(chunk_overlap - last_break):]
                else:
                    current_chunk = current_chunk[-chunk_overlap:]
            else:
                current_chunk = ""
        
        # Add paragraph to current chunk
        if current_chunk and not current_chunk.endswith("\n"):
            current_chunk += "\n\n"
        current_chunk += paragraph
    
    # Add the last chunk if it's not empty
    if current_chunk.strip():
        chunks.append(current_chunk.strip())
    
    return chunks


def compute_similarity(
    query_embedding: List[float],
    document_embeddings: List[List[float]]
) -> List[float]:
    """
    Compute cosine similarity between query and document embeddings
    
    Args:
        query_embedding: Query embedding vector
        document_embeddings: List of document embedding vectors
    
    Returns:
        List of similarity scores (0-1 range, higher is more similar)
    """
    # Convert to numpy arrays for efficient computation
    query_array = np.array(query_embedding)
    doc_array = np.array(document_embeddings)
    
    # Normalize vectors
    query_norm = query_array / np.linalg.norm(query_array)
    doc_norms = doc_array / np.linalg.norm(doc_array, axis=1, keepdims=True)
    
    # Compute cosine similarity
    similarities = np.dot(doc_norms, query_norm)
    
    return similarities.tolist()


async def create_document_chunks_with_embeddings(
    text: str,
    metadata: Dict[str, Any],
    organization_id: str,
    chunk_size: int = 1000,
    chunk_overlap: int = 200
) -> List[Tuple[str, List[float], Dict[str, Any]]]:
    """
    Create document chunks with embeddings
    
    Args:
        text: Document text
        metadata: Document metadata
        organization_id: Organization ID for LLM provider selection
        chunk_size: Maximum chunk size in characters
        chunk_overlap: Overlap between chunks in characters
    
    Returns:
        List of tuples containing (text_chunk, embedding, chunk_metadata)
    """
    # Chunk the text
    chunks = chunk_text(text, chunk_size, chunk_overlap)
    
    # Generate embeddings for all chunks
    embeddings = await get_embeddings(chunks, organization_id)
    
    # Create result tuples
    results = []
    for i, (chunk, embedding) in enumerate(zip(chunks, embeddings)):
        # Copy metadata for each chunk
        chunk_metadata = dict(metadata)
        
        # Add chunk-specific metadata
        chunk_metadata["chunk_index"] = i
        chunk_metadata["chunk_count"] = len(chunks)
        if len(chunks) > 1:
            chunk_metadata["is_chunk"] = True
            chunk_metadata["parent_id"] = metadata.get("parent_id", metadata.get("id", None))
        
        results.append((chunk, embedding, chunk_metadata))
    
    return results