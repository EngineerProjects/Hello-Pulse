"""
Document retrieval service for RAG
"""

from typing import Dict, List, Any, Optional, Tuple
import asyncio

from core.logging import logger, LoggingContext
from services.llm_providers import get_llm_provider
from services.vector_databases import get_vector_db
from utils.embeddings import get_embeddings
from utils.filtering import get_access_filters


async def retrieve_relevant_documents(
    query: str,
    organization_id: str,
    user_id: str,
    filter_dict: Dict[str, Any] = None,
    limit: int = 5,
    rerank: bool = True
) -> List[Dict[str, Any]]:
    """
    Retrieve relevant documents for a query
    
    Args:
        query: Query text
        organization_id: Organization ID
        user_id: User ID
        filter_dict: Additional filter criteria
        limit: Maximum number of documents to retrieve
        rerank: Whether to rerank results for better relevance
        
    Returns:
        List of relevant documents with metadata and scores
    """
    with LoggingContext(organization_id=organization_id, user_id=user_id):
        try:
            # Get LLM provider for embeddings
            llm_provider = get_llm_provider(organization_id)
            
            # Get vector database
            vector_db = get_vector_db(organization_id)
            
            # Generate query embedding
            query_embeddings = await llm_provider.get_embeddings([query])
            query_embedding = query_embeddings[0]
            
            # Get access filters
            access_filters = get_access_filters(
                organization_id=organization_id,
                user_id=user_id
            )
            
            # Merge with additional filters if provided
            if filter_dict:
                # We need to be careful how we merge filters to preserve security
                for key, value in filter_dict.items():
                    # Don't allow overriding security filters
                    if key not in ["organization_id", "user_id", "$or"]:
                        access_filters[key] = value
            
            # Retrieve documents
            documents = await vector_db.search(
                query_embedding=query_embedding,
                filter_dict=access_filters,
                limit=limit * 2 if rerank else limit  # Get more docs if reranking
            )
            
            # Log retrieval results
            logger.info(
                f"Retrieved {len(documents)} documents for query",
                extra={
                    "query": query,
                    "organization_id": organization_id,
                    "user_id": user_id,
                    "document_count": len(documents)
                }
            )
            
            # Rerank if requested and we have a semantic search provider
            if rerank and len(documents) > limit:
                # For now, just use the scores from the vector search
                # In a real implementation, you could use a more sophisticated reranking model
                documents = documents[:limit]
            
            return documents
        
        except Exception as e:
            logger.error(
                f"Error retrieving documents: {str(e)}",
                extra={
                    "query": query,
                    "organization_id": organization_id,
                    "user_id": user_id
                }
            )
            raise


async def store_document(
    text: str,
    metadata: Dict[str, Any],
    organization_id: str,
    user_id: str
) -> str:
    """
    Store a document in the vector database
    
    Args:
        text: Document text
        metadata: Document metadata
        organization_id: Organization ID
        user_id: User ID
        
    Returns:
        Document ID
    """
    with LoggingContext(organization_id=organization_id, user_id=user_id):
        try:
            # Get LLM provider for embeddings
            llm_provider = get_llm_provider(organization_id)
            
            # Get vector database
            vector_db = get_vector_db(organization_id)
            
            # Generate embedding
            embeddings = await llm_provider.get_embeddings([text])
            embedding = embeddings[0]
            
            # Ensure metadata includes required fields
            metadata["organization_id"] = organization_id
            metadata["user_id"] = user_id
            
            # Store document
            doc_ids = await vector_db.add_documents(
                texts=[text],
                embeddings=[embedding],
                metadata=[metadata]
            )
            
            logger.info(
                f"Stored document in vector database",
                extra={
                    "document_id": doc_ids[0],
                    "organization_id": organization_id,
                    "user_id": user_id,
                    "text_length": len(text)
                }
            )
            
            return doc_ids[0]
        
        except Exception as e:
            logger.error(
                f"Error storing document: {str(e)}",
                extra={
                    "organization_id": organization_id,
                    "user_id": user_id,
                    "text_length": len(text)
                }
            )
            raise


async def delete_document(
    document_id: str,
    organization_id: str,
    user_id: str
) -> bool:
    """
    Delete a document from the vector database
    
    Args:
        document_id: Document ID
        organization_id: Organization ID
        user_id: User ID
        
    Returns:
        True if successful, False otherwise
    """
    with LoggingContext(organization_id=organization_id, user_id=user_id):
        try:
            # Get vector database
            vector_db = get_vector_db(organization_id)
            
            # Get document to verify ownership
            document = await vector_db.get_document(document_id)
            
            if not document:
                logger.warning(
                    f"Document not found for deletion",
                    extra={
                        "document_id": document_id,
                        "organization_id": organization_id,
                        "user_id": user_id
                    }
                )
                return False
            
            # Verify document ownership or admin status
            doc_org_id = document["metadata"].get("organization_id")
            doc_user_id = document["metadata"].get("user_id")
            
            if doc_org_id != organization_id:
                logger.warning(
                    f"Attempted to delete document from another organization",
                    extra={
                        "document_id": document_id,
                        "organization_id": organization_id,
                        "user_id": user_id,
                        "document_org_id": doc_org_id
                    }
                )
                return False
            
            if doc_user_id != user_id:
                # Check if user is an admin - we'd implement proper admin checks in a real system
                # For now, just check if the user is allowed to delete others' documents
                is_admin = False  # Example placeholder - would be actually checked
                
                if not is_admin:
                    logger.warning(
                        f"Attempted to delete another user's document without admin rights",
                        extra={
                            "document_id": document_id,
                            "organization_id": organization_id,
                            "user_id": user_id,
                            "document_user_id": doc_user_id
                        }
                    )
                    return False
            
            # Delete document
            success = await vector_db.delete_documents([document_id])
            
            if success:
                logger.info(
                    f"Deleted document from vector database",
                    extra={
                        "document_id": document_id,
                        "organization_id": organization_id,
                        "user_id": user_id
                    }
                )
            
            return success
        
        except Exception as e:
            logger.error(
                f"Error deleting document: {str(e)}",
                extra={
                    "document_id": document_id,
                    "organization_id": organization_id,
                    "user_id": user_id
                }
            )
            return False