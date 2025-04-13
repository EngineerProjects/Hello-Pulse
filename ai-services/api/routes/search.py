"""
API routes for web search
"""

from typing import Dict, Any, Optional, List
from fastapi import APIRouter, Depends, HTTPException, Query, Request
from fastapi.responses import StreamingResponse

from api.dependencies import (
    get_validated_user,
    get_llm_for_request,
    get_request_id
)
from core.security import TokenPayload
from schemas.requests import WebSearchRequest, WebSearchAndGenerateRequest
from schemas.responses import (
    WebSearchResponse,
    WebSearchAndGenerateResponse,
    WebSearchResult,
    ErrorResponse
)
from services.llm_providers.base import BaseLLMProvider
from services.websearch import get_search_provider, search_and_generate
from core.logging import logger, LoggingContext

router = APIRouter(tags=["Web Search"])

@router.post(
    "/search",
    response_model=WebSearchResponse,
    responses={
        400: {"model": ErrorResponse},
        401: {"model": ErrorResponse},
        500: {"model": ErrorResponse}
    },
    summary="Search the web",
    description="Search the web using the configured search provider"
)
async def web_search(
    request: Request,
    data: WebSearchRequest,
    request_id: str = Depends(get_request_id),
    token_data: TokenPayload = Depends(get_validated_user)
):
    """
    Search the web
    
    Args:
        request: Request object
        data: Web search request data
        request_id: Request ID
        token_data: Token payload with user and organization info
        
    Returns:
        Web search results
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
            # Get search provider
            search_provider = get_search_provider()
            
            # Perform search
            results = await search_provider.search(
                query=data.query,
                num_results=data.num_results
            )
            
            logger.info(
                "Web search completed",
                extra={
                    "query": data.query,
                    "result_count": len(results),
                    "organization_id": organization_id,
                    "user_id": user_id,
                    "request_id": request_id,
                    "provider": search_provider.get_provider_name()
                }
            )
            
            # Format results for response
            search_results = []
            for result in results:
                search_results.append(
                    WebSearchResult(
                        title=result.title,
                        url=result.url,
                        snippet=result.snippet if data.include_snippets else None,
                        source=result.source
                    )
                )
            
            return WebSearchResponse(
                results=search_results,
                query=data.query,
                count=len(results)
            )
        
        except Exception as e:
            logger.error(
                f"Web search error: {str(e)}",
                extra={
                    "query": data.query,
                    "organization_id": organization_id,
                    "user_id": user_id,
                    "request_id": request_id
                }
            )
            raise HTTPException(
                status_code=500,
                detail=f"Web search failed: {str(e)}"
            )


@router.post(
    "/search_and_generate",
    response_model=WebSearchAndGenerateResponse,
    responses={
        400: {"model": ErrorResponse},
        401: {"model": ErrorResponse},
        500: {"model": ErrorResponse}
    },
    summary="Search the web and generate an answer",
    description="Search the web and generate an enhanced answer based on the results"
)
async def web_search_and_generate(
    request: Request,
    data: WebSearchAndGenerateRequest,
    request_id: str = Depends(get_request_id),
    token_data: TokenPayload = Depends(get_validated_user),
    llm_provider: BaseLLMProvider = Depends(get_llm_for_request)
):
    """
    Search the web and generate an answer
    
    Args:
        request: Request object
        data: Web search and generate request data
        request_id: Request ID
        token_data: Token payload with user and organization info
        llm_provider: LLM provider
        
    Returns:
        Generated answer and search results
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
                async def search_stream():
                    try:
                        # Perform search first
                        search_provider = get_search_provider()
                        search_results = await search_provider.search(
                            query=data.query,
                            num_results=data.num_results
                        )
                        
                        # Format search results for the prompt
                        from services.websearch import format_search_results
                        search_context = format_search_results(search_results)
                        
                        # Create system message
                        enhanced_system_message = data.system_message or (
                            "You are a helpful assistant with access to web search results. "
                            "Answer the user's question based on the provided search results. "
                            "If the search results don't contain the necessary information to answer the question, "
                            "just summarize what information is available. "
                            "Include relevant URLs in your answer where appropriate, but do not list all URLs."
                        )
                        
                        # Create prompt with search results
                        search_prompt = f"""I need information about the following query:
{data.query}

Here are the top search results:

{search_context}

Based on these search results, please provide a comprehensive answer to the query."""
                        
                        # Stream the response
                        async for chunk in llm_provider.generate_stream(
                            prompt=search_prompt,
                            system_message=enhanced_system_message,
                            temperature=data.temperature,
                            max_tokens=data.max_tokens
                        ):
                            # Yield chunk as server-sent event
                            yield f"data: {chunk}\n\n"
                        
                        # End of stream
                        yield "data: [DONE]\n\n"
                    except Exception as e:
                        logger.error(
                            f"Search and generate streaming error: {str(e)}",
                            extra={
                                "query": data.query,
                                "organization_id": organization_id,
                                "user_id": user_id,
                                "request_id": request_id
                            }
                        )
                        # Send error as event
                        yield f"data: [ERROR] {str(e)}\n\n"
                
                return StreamingResponse(
                    search_stream(),
                    media_type="text/event-stream"
                )
            
            # Non-streaming request
            result = await search_and_generate(
                query=data.query,
                organization_id=organization_id,
                user_id=user_id,
                num_results=data.num_results,
                system_message=data.system_message,
                temperature=data.temperature,
                max_tokens=data.max_tokens
            )
            
            logger.info(
                "Web search and generate completed",
                extra={
                    "query": data.query,
                    "answer_length": len(result["answer"]),
                    "result_count": len(result["search_results"]),
                    "organization_id": organization_id,
                    "user_id": user_id,
                    "request_id": request_id,
                    "model": llm_provider.get_model_name(),
                    "provider": llm_provider.get_provider_name()
                }
            )
            
            # Format search results for response
            search_results = []
            for result_dict in result["search_results"]:
                search_results.append(
                    WebSearchResult(
                        title=result_dict["title"],
                        url=result_dict["url"],
                        snippet=result_dict.get("snippet"),
                        source=result_dict["source"]
                    )
                )
            
            return WebSearchAndGenerateResponse(
                answer=result["answer"],
                search_results=search_results,
                model=llm_provider.get_model_name(),
                provider=llm_provider.get_provider_name(),
                finish_reason="stop",  # Placeholder
                usage=None  # Placeholder
            )
        
        except Exception as e:
            logger.error(
                f"Search and generate error: {str(e)}",
                extra={
                    "query": data.query,
                    "organization_id": organization_id,
                    "user_id": user_id,
                    "request_id": request_id
                }
            )
            raise HTTPException(
                status_code=500,
                detail=f"Search and generate failed: {str(e)}"
            )