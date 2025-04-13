"""
API routes for text generation
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
from schemas.requests import GenerateRequest, ChatRequest
from schemas.responses import GenerateResponse, ChatResponse, StreamChunk, ErrorResponse
from services.llm_providers.base import BaseLLMProvider
from core.logging import logger, LoggingContext

router = APIRouter(tags=["Text Generation"])

@router.post(
    "/generate",
    response_model=GenerateResponse,
    responses={
        400: {"model": ErrorResponse},
        401: {"model": ErrorResponse},
        500: {"model": ErrorResponse}
    },
    summary="Generate text from a prompt",
    description="Generate text from a prompt using the configured LLM provider"
)
async def generate_text(
    request: Request,
    data: GenerateRequest,
    request_id: str = Depends(get_request_id),
    token_data: TokenPayload = Depends(get_validated_user),
    llm_provider: BaseLLMProvider = Depends(get_llm_for_request)
):
    """
    Generate text from a prompt
    
    Args:
        request: Request object
        data: Generate request data
        request_id: Request ID
        token_data: Token payload with user and organization info
        llm_provider: LLM provider
        
    Returns:
        Generated text
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
                async def generate_stream():
                    try:
                        async for chunk in llm_provider.generate_stream(
                            prompt=data.prompt,
                            system_message=data.system_message,
                            temperature=data.temperature,
                            max_tokens=data.max_tokens,
                            stop_sequences=data.stop_sequences,
                            **(data.model_params or {})
                        ):
                            # Yield chunk as server-sent event
                            yield f"data: {chunk}\n\n"
                        
                        # End of stream
                        yield "data: [DONE]\n\n"
                    except Exception as e:
                        logger.error(
                            f"Streaming error: {str(e)}",
                            extra={
                                "organization_id": organization_id,
                                "user_id": user_id,
                                "request_id": request_id
                            }
                        )
                        # Send error as event
                        yield f"data: [ERROR] {str(e)}\n\n"
                
                return StreamingResponse(
                    generate_stream(),
                    media_type="text/event-stream"
                )
            
            # Non-streaming request
            generated_text = await llm_provider.generate(
                prompt=data.prompt,
                system_message=data.system_message,
                temperature=data.temperature,
                max_tokens=data.max_tokens,
                stop_sequences=data.stop_sequences,
                **(data.model_params or {})
            )
            
            logger.info(
                "Text generation completed",
                extra={
                    "prompt_length": len(data.prompt),
                    "response_length": len(generated_text),
                    "organization_id": organization_id,
                    "user_id": user_id,
                    "request_id": request_id,
                    "model": llm_provider.get_model_name(),
                    "provider": llm_provider.get_provider_name()
                }
            )
            
            return GenerateResponse(
                text=generated_text,
                model=llm_provider.get_model_name(),
                provider=llm_provider.get_provider_name(),
                finish_reason="stop",  # Placeholder - actual finish reason would come from provider
                usage=None  # Placeholder - token usage would come from provider
            )
        
        except Exception as e:
            logger.error(
                f"Generation error: {str(e)}",
                extra={
                    "organization_id": organization_id,
                    "user_id": user_id,
                    "request_id": request_id
                }
            )
            raise HTTPException(
                status_code=500,
                detail=f"Generation failed: {str(e)}"
            )


@router.post(
    "/chat",
    response_model=ChatResponse,
    responses={
        400: {"model": ErrorResponse},
        401: {"model": ErrorResponse},
        500: {"model": ErrorResponse}
    },
    summary="Chat with the AI",
    description="Generate a response based on a chat history"
)
async def chat(
    request: Request,
    data: ChatRequest,
    request_id: str = Depends(get_request_id),
    token_data: TokenPayload = Depends(get_validated_user),
    llm_provider: BaseLLMProvider = Depends(get_llm_for_request)
):
    """
    Chat with the AI
    
    Args:
        request: Request object
        data: Chat request data
        request_id: Request ID
        token_data: Token payload with user and organization info
        llm_provider: LLM provider
        
    Returns:
        Chat response
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
                async def chat_stream():
                    try:
                        messages = [msg.dict() for msg in data.messages]
                        
                        async for chunk in llm_provider.chat_stream(
                            messages=messages,
                            temperature=data.temperature,
                            max_tokens=data.max_tokens,
                            stop_sequences=data.stop_sequences,
                            **(data.model_params or {})
                        ):
                            # Yield chunk as server-sent event
                            yield f"data: {chunk}\n\n"
                        
                        # End of stream
                        yield "data: [DONE]\n\n"
                    except Exception as e:
                        logger.error(
                            f"Chat streaming error: {str(e)}",
                            extra={
                                "organization_id": organization_id,
                                "user_id": user_id,
                                "request_id": request_id
                            }
                        )
                        # Send error as event
                        yield f"data: [ERROR] {str(e)}\n\n"
                
                return StreamingResponse(
                    chat_stream(),
                    media_type="text/event-stream"
                )
            
            # Non-streaming request
            messages = [msg.dict() for msg in data.messages]
            
            response = await llm_provider.chat(
                messages=messages,
                temperature=data.temperature,
                max_tokens=data.max_tokens,
                stop_sequences=data.stop_sequences,
                **(data.model_params or {})
            )
            
            logger.info(
                "Chat completed",
                extra={
                    "message_count": len(messages),
                    "response_length": len(response),
                    "organization_id": organization_id,
                    "user_id": user_id,
                    "request_id": request_id,
                    "model": llm_provider.get_model_name(),
                    "provider": llm_provider.get_provider_name()
                }
            )
            
            return ChatResponse(
                message=response,
                model=llm_provider.get_model_name(),
                provider=llm_provider.get_provider_name(),
                finish_reason="stop",  # Placeholder - actual finish reason would come from provider
                usage=None  # Placeholder - token usage would come from provider
            )
        
        except Exception as e:
            logger.error(
                f"Chat error: {str(e)}",
                extra={
                    "organization_id": organization_id,
                    "user_id": user_id,
                    "request_id": request_id
                }
            )
            raise HTTPException(
                status_code=500,
                detail=f"Chat failed: {str(e)}"
            )