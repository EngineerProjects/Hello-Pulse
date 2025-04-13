"""
API routes for RAG (Retrieval-Augmented Generation)
"""

from typing import Dict, Any, Optional, List
from fastapi import APIRouter, Depends, HTTPException, Query, Request, Body
from fastapi.responses import StreamingResponse

from api.dependencies import (
    get_validated_user,
    get_llm_for_request,
    get_vector_db_for_request,
    get_request_id
)
from core.security import TokenPayload
from schemas.requests import RAGRequest, AddDocumentsRequest, SearchRequest
from schemas.responses import (
    RAGResponse,
    DocumentResponse,
    DocumentListResponse,
    SearchResponse,
    ErrorResponse
)
from services.llm_providers.base import BaseLLMProvider
from services.vector_databases.base import BaseVectorDatabase
from services.rag.generator import generate_answer, generate_answer_stream
from services.rag.retriever import (
    retrieve_relevant_documents,
    store_document,
    delete_document
)
from utils.embeddings import get_embeddings, create_document_chunks_with_embeddings
from utils.filtering import sanitize_metadata_for_storage
from core.logging import logger, LoggingContext

router = APIRouter(tags=["RAG"])

@router.post(
    "/rag/generate",
    response_model=RAGResponse,
    responses={
        400: {"model": ErrorResponse},
        401: {"model": ErrorResponse},
        500: {"model": ErrorResponse}
    },
    summary="Generate an answer using RAG",
    description="Generate an answer based on retrieved documents"
)
async def rag_generate(
    request: Request,
    data: RAGRequest,
    request_id: str = Depends(get_request_id),
    token_data: TokenPayload = Depends(get_validated_user),
    llm_provider: BaseLLMProvider = Depends(get_llm_for_request)
):
    """
    Generate an answer using RAG
    
    Args:
        request: Request object
        data: RAG request data
        request_id: Request ID
        token_data: Token payload with user and organization info
        llm_provider: LLM provider
        
    Returns:
        Generated answer and retrieved documents
    """
    # Extract request parameters
    organization_id = token_data.organization_id
    user_id = token_data.user_id
    
    with LoggingContext(
        organization_id=organization_id,
        user_id=user_id,
        request_id=request_id
    ):
        try:
            # Check if streaming is requested
            if data.stream:
                # Create streaming response
                async def rag_stream():
                    try:
                        # Generate streaming answer
                        stream, docs = await generate_answer_stream(
                            query=data.query,
                            organization_id=organization_id,
                            user_id=user_id,
                            filter_dict=data.filter,
                            system_message=data.system_message,
                            temperature=data.temperature,
                            max_tokens=data.max_tokens,
                            num_documents=data.num_documents
                        )
                        
                        # Stream the answer
                        async for chunk in stream:
                            # Yield chunk as server-sent event
                            yield f"data: {chunk}\n\n"
                        
                        # End of stream
                        yield "data: [DONE]\n\n"
                    except Exception as e:
                        logger.error(
                            f"RAG streaming error: {str(e)}",
                            extra={
                                "organization_id": organization_id,
                                "user_id": user_id,
                                "request_id": request_id
                            }
                        )
                        # Send error as event
                        yield f"data: [ERROR] {str(e)}\n\n"
                
                return StreamingResponse(
                    rag_stream(),
                    media_type="text/event-stream"
                )
            
            # Non-streaming request
            answer, docs = await generate_answer(
                query=data.query,
                organization_id=organization_id,
                user_id=user_id,
                filter_dict=data.filter,
                system_message=data.system_message,
                temperature=data.temperature,
                max_tokens=data.max_tokens,
                num_documents=data.num_documents
            )
            
            logger.info(
                "RAG generation completed",
                extra={
                    "query_length": len(data.query),
                    "answer_length": len(answer),
                    "document_count": len(docs),
                    "organization_id": organization_id,
                    "user_id": user_id,
                    "request_id": request_id,
                    "model": llm_provider.get_model_name(),
                    "provider": llm_provider.get_provider_name()
                }
            )
            
            # Format documents for response
            doc_responses = []
            for doc in docs:
                doc_responses.append(
                    DocumentResponse(
                        id=doc["id"],
                        text=doc["text"],
                        metadata=doc["metadata"],
                        score=doc.get("score")
                    )
                )
            
            return RAGResponse(
                answer=answer,
                documents=doc_responses,
                model=llm_provider.get_model_name(),
                provider=llm_provider.get_provider_name(),
                finish_reason="stop",  # Placeholder
                usage=None  # Placeholder
            )
        
        except Exception as e:
            logger.error(
                f"RAG generation error: {str(e)}",
                extra={
                    "organization_id": organization_id,
                    "user_id": user_id,
                    "request_id": request_id
                }
            )
            raise HTTPException(
                status_code=500,
                detail=f"RAG generation failed: {str(e)}"
            )


@router.post(
    "/rag/search",
    response_model=SearchResponse,
    responses={
        400: {"model": ErrorResponse},
        401: {"model": ErrorResponse},
        500: {"model": ErrorResponse}
    },
    summary="Search for documents",
    description="Search for documents using vector similarity"
)
async def search_documents(
    request: Request,
    data: SearchRequest,
    request_id: str = Depends(get_request_id),
    token_data: TokenPayload = Depends(get_validated_user)
):
    """
    Search for documents
    
    Args:
        request: Request object
        data: Search request data
        request_id: Request ID
        token_data: Token payload with user and organization info
        
    Returns:
        Search results
    """
    # Extract request parameters
    organization_id = token_data.organization_id
    user_id = token_data.user_id
    
    with LoggingContext(
        organization_id=organization_id,
        user_id=user_id,
        request_id=request_id
    ):
        try:
            # Get query embedding
            embeddings = await get_embeddings(
                texts=[data.query],
                organization_id=organization_id
            )
            
            # Execute search
            results = await retrieve_relevant_documents(
                query=data.query,
                organization_id=organization_id,
                user_id=user_id,
                filter_dict=data.filter,
                limit=data.limit
            )
            
            logger.info(
                "Document search completed",
                extra={
                    "query_length": len(data.query),
                    "result_count": len(results),
                    "organization_id": organization_id,
                    "user_id": user_id,
                    "request_id": request_id
                }
            )
            
            # Format results for response
            doc_responses = []
            for doc in results:
                doc_response = DocumentResponse(
                    id=doc["id"],
                    text=doc["text"],
                    metadata=doc["metadata"],
                    score=doc.get("score")
                )
                
                # Include embedding if requested
                if data.include_embeddings and "embedding" in doc:
                    doc_response.embedding = doc["embedding"]
                
                doc_responses.append(doc_response)
            
            return SearchResponse(
                results=doc_responses,
                query=data.query,
                count=len(results),
                total=None  # We don't have a total count for now
            )
        
        except Exception as e:
            logger.error(
                f"Document search error: {str(e)}",
                extra={
                    "organization_id": organization_id,
                    "user_id": user_id,
                    "request_id": request_id
                }
            )
            raise HTTPException(
                status_code=500,
                detail=f"Document search failed: {str(e)}"
            )


@router.post(
    "/rag/documents",
    response_model=DocumentListResponse,
    responses={
        400: {"model": ErrorResponse},
        401: {"model": ErrorResponse},
        500: {"model": ErrorResponse}
    },
    summary="Add documents",
    description="Add documents to the vector database"
)
async def add_documents(
    request: Request,
    data: AddDocumentsRequest,
    request_id: str = Depends(get_request_id),
    token_data: TokenPayload = Depends(get_validated_user)
):
    """
    Add documents
    
    Args:
        request: Request object
        data: Add documents request data
        request_id: Request ID
        token_data: Token payload with user and organization info
        
    Returns:
        Added documents
    """
    # Extract request parameters
    organization_id = token_data.organization_id
    user_id = token_data.user_id
    
    with LoggingContext(
        organization_id=organization_id,
        user_id=user_id,
        request_id=request_id
    ):
        try:
            # Get vector database
            vector_db = get_vector_db_for_request(token_data)
            
            # Process each document
            document_texts = []
            document_metadata = []
            
            for doc in data.documents:
                # Sanitize metadata
                metadata = sanitize_metadata_for_storage(
                    metadata=doc.metadata.dict(),
                    organization_id=organization_id,
                    user_id=user_id
                )
                
                document_texts.append(doc.text)
                document_metadata.append(metadata)
            
            # Generate embeddings
            embeddings = await get_embeddings(
                texts=document_texts,
                organization_id=organization_id
            )
            
            # Store documents
            doc_ids = await vector_db.add_documents(
                texts=document_texts,
                embeddings=embeddings,
                metadata=document_metadata
            )
            
            logger.info(
                "Documents added",
                extra={
                    "document_count": len(document_texts),
                    "organization_id": organization_id,
                    "user_id": user_id,
                    "request_id": request_id
                }
            )
            
            # Fetch the added documents for response
            added_docs = []
            for doc_id in doc_ids:
                doc = await vector_db.get_document(doc_id)
                if doc:
                    added_docs.append(
                        DocumentResponse(
                            id=doc["id"],
                            text=doc["text"],
                            metadata=doc["metadata"]
                        )
                    )
            
            return DocumentListResponse(
                documents=added_docs,
                count=len(added_docs),
                limit=len(doc_ids),
                offset=0
            )
        
        except Exception as e:
            logger.error(
                f"Add documents error: {str(e)}",
                extra={
                    "organization_id": organization_id,
                    "user_id": user_id,
                    "request_id": request_id
                }
            )
            raise HTTPException(
                status_code=500,
                detail=f"Adding documents failed: {str(e)}"
            )


@router.get(
    "/rag/documents/{document_id}",
    response_model=DocumentResponse,
    responses={
        400: {"model": ErrorResponse},
        401: {"model": ErrorResponse},
        404: {"model": ErrorResponse},
        500: {"model": ErrorResponse}
    },
    summary="Get document",
    description="Get a document by ID"
)
async def get_document(
    request: Request,
    document_id: str,
    request_id: str = Depends(get_request_id),
    token_data: TokenPayload = Depends(get_validated_user),
    vector_db: BaseVectorDatabase = Depends(get_vector_db_for_request)
):
    """
    Get document
    
    Args:
        request: Request object
        document_id: Document ID
        request_id: Request ID
        token_data: Token payload with user and organization info
        vector_db: Vector database
        
    Returns:
        Document
    """
    # Extract request parameters
    organization_id = token_data.organization_id
    user_id = token_data.user_id
    
    with LoggingContext(
        organization_id=organization_id,
        user_id=user_id,
        request_id=request_id
    ):
        try:
            # Get document
            doc = await vector_db.get_document(document_id)
            
            if not doc:
                raise HTTPException(
                    status_code=404,
                    detail="Document not found"
                )
            
            # Check access permissions
            doc_org_id = doc["metadata"].get("organization_id")
            if doc_org_id != organization_id:
                logger.warning(
                    f"Attempted to access document from another organization",
                    extra={
                        "document_id": document_id,
                        "organization_id": organization_id,
                        "document_org_id": doc_org_id,
                        "user_id": user_id,
                        "request_id": request_id
                    }
                )
                raise HTTPException(
                    status_code=404,
                    detail="Document not found"
                )
            
            return DocumentResponse(
                id=doc["id"],
                text=doc["text"],
                metadata=doc["metadata"]
            )
        
        except HTTPException:
            raise
        
        except Exception as e:
            logger.error(
                f"Get document error: {str(e)}",
                extra={
                    "document_id": document_id,
                    "organization_id": organization_id,
                    "user_id": user_id,
                    "request_id": request_id
                }
            )
            raise HTTPException(
                status_code=500,
                detail=f"Error retrieving document: {str(e)}"
            )


@router.delete(
    "/rag/documents/{document_id}",
    responses={
        400: {"model": ErrorResponse},
        401: {"model": ErrorResponse},
        404: {"model": ErrorResponse},
        500: {"model": ErrorResponse}
    },
    summary="Delete document",
    description="Delete a document by ID"
)
async def delete_document_endpoint(
    request: Request,
    document_id: str,
    request_id: str = Depends(get_request_id),
    token_data: TokenPayload = Depends(get_validated_user)
):
    """
    Delete document
    
    Args:
        request: Request object
        document_id: Document ID
        request_id: Request ID
        token_data: Token payload with user and organization info
        
    Returns:
        Success message
    """
    # Extract request parameters
    organization_id = token_data.organization_id
    user_id = token_data.user_id
    
    with LoggingContext(
        organization_id=organization_id,
        user_id=user_id,
        request_id=request_id
    ):
        try:
            # Delete document
            success = await delete_document(
                document_id=document_id,
                organization_id=organization_id,
                user_id=user_id
            )
            
            if not success:
                raise HTTPException(
                    status_code=404,
                    detail="Document not found or you don't have permission to delete it"
                )
            
            logger.info(
                "Document deleted",
                extra={
                    "document_id": document_id,
                    "organization_id": organization_id,
                    "user_id": user_id,
                    "request_id": request_id
                }
            )
            
            return {"success": True, "message": "Document deleted successfully"}
        
        except HTTPException:
            raise
        
        except Exception as e:
            logger.error(
                f"Delete document error: {str(e)}",
                extra={
                    "document_id": document_id,
                    "organization_id": organization_id,
                    "user_id": user_id,
                    "request_id": request_id
                }
            )
            raise HTTPException(
                status_code=500,
                detail=f"Error deleting document: {str(e)}"
            )


@router.get(
    "/rag/documents",
    response_model=DocumentListResponse,
    responses={
        400: {"model": ErrorResponse},
        401: {"model": ErrorResponse},
        500: {"model": ErrorResponse}
    },
    summary="List documents",
    description="List documents with optional filtering"
)
async def list_documents(
    request: Request,
    request_id: str = Depends(get_request_id),
    token_data: TokenPayload = Depends(get_validated_user),
    vector_db: BaseVectorDatabase = Depends(get_vector_db_for_request),
    limit: int = Query(100, ge=1, le=1000),
    offset: int = Query(0, ge=0),
    filter: Optional[str] = Query(None)
):
    """
    List documents
    
    Args:
        request: Request object
        request_id: Request ID
        token_data: Token payload with user and organization info
        vector_db: Vector database
        limit: Maximum number of results
        offset: Pagination offset
        filter: Optional JSON filter string
        
    Returns:
        List of documents
    """
    # Extract request parameters
    organization_id = token_data.organization_id
    user_id = token_data.user_id
    
    with LoggingContext(
        organization_id=organization_id,
        user_id=user_id,
        request_id=request_id
    ):
        try:
            import json
            
            # Parse filter if provided
            filter_dict = {}
            if filter:
                try:
                    filter_dict = json.loads(filter)
                except json.JSONDecodeError:
                    raise HTTPException(
                        status_code=400,
                        detail="Invalid filter JSON"
                    )
            
            # Apply organization filter for security
            if "organization_id" not in filter_dict:
                filter_dict["organization_id"] = organization_id
            elif filter_dict["organization_id"] != organization_id:
                # Don't allow filtering by other organizations
                filter_dict["organization_id"] = organization_id
            
            # Get document count
            total_count = await vector_db.count_documents(filter_dict)
            
            # Get documents
            docs = await vector_db.list_documents(
                filter_dict=filter_dict,
                limit=limit,
                offset=offset
            )
            
            logger.info(
                "Documents listed",
                extra={
                    "document_count": len(docs),
                    "total_count": total_count,
                    "organization_id": organization_id,
                    "user_id": user_id,
                    "request_id": request_id
                }
            )
            
            # Format for response
            doc_responses = []
            for doc in docs:
                doc_responses.append(
                    DocumentResponse(
                        id=doc["id"],
                        text=doc["text"],
                        metadata=doc["metadata"]
                    )
                )
            
            return DocumentListResponse(
                documents=doc_responses,
                count=total_count,
                limit=limit,
                offset=offset
            )
        
        except HTTPException:
            raise
        
        except Exception as e:
            logger.error(
                f"List documents error: {str(e)}",
                extra={
                    "organization_id": organization_id,
                    "user_id": user_id,
                    "request_id": request_id
                }
            )
            raise HTTPException(
                status_code=500,
                detail=f"Error listing documents: {str(e)}"
            )