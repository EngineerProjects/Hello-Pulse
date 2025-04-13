"""
API routes for AI agents
"""

from typing import Dict, Any, Optional, List
from fastapi import APIRouter, Depends, HTTPException, Query, Request, Body
from fastapi.responses import StreamingResponse
import json

from api.dependencies import (
    get_validated_user,
    get_llm_for_request,
    get_request_id
)
from core.security import TokenPayload
from schemas.requests import (
    AgentCreateRequest,
    AgentChatRequest,
    AgentListRequest
)
from schemas.responses import (
    AgentResponse,
    AgentListResponse,
    AgentChatResponse,
    ErrorResponse
)
from services.llm_providers.base import BaseLLMProvider
from services.agents.manager import agent_manager
from core.logging import logger, LoggingContext

router = APIRouter(tags=["AI Agents"])

@router.post(
    "/agents",
    response_model=AgentResponse,
    responses={
        400: {"model": ErrorResponse},
        401: {"model": ErrorResponse},
        500: {"model": ErrorResponse}
    },
    summary="Create an agent",
    description="Create a new AI agent"
)
async def create_agent(
    request: Request,
    data: AgentCreateRequest,
    request_id: str = Depends(get_request_id),
    token_data: TokenPayload = Depends(get_validated_user)
):
    """
    Create an agent
    
    Args:
        request: Request object
        data: Agent creation request data
        request_id: Request ID
        token_data: Token payload with user and organization info
        
    Returns:
        Created agent
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
            # Create agent
            agent_data = await agent_manager.create_agent(
                organization_id=organization_id,
                user_id=user_id,
                name=data.name,
                instructions=data.instructions,
                description=data.description,
                document_ids=data.document_ids,
                document_filter=data.document_filter,
                metadata=data.metadata
            )
            
            logger.info(
                "Agent created",
                extra={
                    "agent_id": agent_data["id"],
                    "agent_name": data.name,
                    "organization_id": organization_id,
                    "user_id": user_id,
                    "request_id": request_id
                }
            )
            
            # Convert to response model
            from datetime import datetime
            return AgentResponse(
                id=agent_data["id"],
                name=agent_data["name"],
                description=agent_data["description"],
                instructions=agent_data["instructions"],
                created_at=datetime.fromisoformat(agent_data["created_at"]),
                metadata=agent_data.get("metadata", {}),
                document_count=agent_data.get("document_count", 0)
            )
        
        except Exception as e:
            logger.error(
                f"Agent creation error: {str(e)}",
                extra={
                    "agent_name": data.name,
                    "organization_id": organization_id,
                    "user_id": user_id,
                    "request_id": request_id
                }
            )
            raise HTTPException(
                status_code=500,
                detail=f"Agent creation failed: {str(e)}"
            )


@router.get(
    "/agents/{agent_id}",
    response_model=AgentResponse,
    responses={
        400: {"model": ErrorResponse},
        401: {"model": ErrorResponse},
        404: {"model": ErrorResponse},
        500: {"model": ErrorResponse}
    },
    summary="Get agent",
    description="Get an agent by ID"
)
async def get_agent(
    request: Request,
    agent_id: str,
    request_id: str = Depends(get_request_id),
    token_data: TokenPayload = Depends(get_validated_user)
):
    """
    Get agent
    
    Args:
        request: Request object
        agent_id: Agent ID
        request_id: Request ID
        token_data: Token payload with user and organization info
        
    Returns:
        Agent
    """
    # Extract request parameters
    organization_id = token_data.organization_id
    user_id = token_data.user_id
    
    with LoggingContext(
        organization_id=organization_id,
        user_id=user_id,
        request_id=request_id,
        agent_id=agent_id
    ):
        try:
            # Get agent
            agent_data = await agent_manager.get_agent(
                agent_id=agent_id,
                organization_id=organization_id,
                user_id=user_id
            )
            
            if not agent_data:
                raise HTTPException(
                    status_code=404,
                    detail="Agent not found"
                )
            
            # Convert to response model
            from datetime import datetime
            return AgentResponse(
                id=agent_data["id"],
                name=agent_data["name"],
                description=agent_data["description"],
                instructions=agent_data["instructions"],
                created_at=datetime.fromisoformat(agent_data["created_at"]) if isinstance(agent_data["created_at"], str) else agent_data["created_at"],
                updated_at=datetime.fromisoformat(agent_data["updated_at"]) if agent_data.get("updated_at") and isinstance(agent_data["updated_at"], str) else agent_data.get("updated_at"),
                metadata=agent_data.get("metadata", {}),
                document_count=agent_data.get("document_count", 0)
            )
        
        except HTTPException:
            raise
        
        except Exception as e:
            logger.error(
                f"Get agent error: {str(e)}",
                extra={
                    "agent_id": agent_id,
                    "organization_id": organization_id,
                    "user_id": user_id,
                    "request_id": request_id
                }
            )
            raise HTTPException(
                status_code=500,
                detail=f"Error retrieving agent: {str(e)}"
            )


@router.get(
    "/agents",
    response_model=AgentListResponse,
    responses={
        400: {"model": ErrorResponse},
        401: {"model": ErrorResponse},
        500: {"model": ErrorResponse}
    },
    summary="List agents",
    description="List agents with optional filtering"
)
async def list_agents(
    request: Request,
    request_id: str = Depends(get_request_id),
    token_data: TokenPayload = Depends(get_validated_user),
    limit: int = Query(100, ge=1, le=1000),
    offset: int = Query(0, ge=0),
    filter: Optional[str] = Query(None)
):
    """
    List agents
    
    Args:
        request: Request object
        request_id: Request ID
        token_data: Token payload with user and organization info
        limit: Maximum number of results
        offset: Pagination offset
        filter: Optional JSON filter string
        
    Returns:
        List of agents
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
            
            # List agents
            agents, total_count = await agent_manager.list_agents(
                organization_id=organization_id,
                user_id=user_id,
                filter_dict=filter_dict,
                limit=limit,
                offset=offset
            )
            
            logger.info(
                "Agents listed",
                extra={
                    "agent_count": len(agents),
                    "total_count": total_count,
                    "organization_id": organization_id,
                    "user_id": user_id,
                    "request_id": request_id
                }
            )
            
            # Convert to response model
            from datetime import datetime
            agent_responses = []
            for agent_data in agents:
                agent_responses.append(
                    AgentResponse(
                        id=agent_data["id"],
                        name=agent_data["name"],
                        description=agent_data["description"],
                        instructions=agent_data["instructions"],
                        created_at=datetime.fromisoformat(agent_data["created_at"]) if isinstance(agent_data["created_at"], str) else agent_data["created_at"],
                        updated_at=datetime.fromisoformat(agent_data["updated_at"]) if agent_data.get("updated_at") and isinstance(agent_data["updated_at"], str) else agent_data.get("updated_at"),
                        metadata=agent_data.get("metadata", {}),
                        document_count=agent_data.get("document_count", 0)
                    )
                )
            
            return AgentListResponse(
                agents=agent_responses,
                count=total_count,
                limit=limit,
                offset=offset
            )
        
        except HTTPException:
            raise
        
        except Exception as e:
            logger.error(
                f"List agents error: {str(e)}",
                extra={
                    "organization_id": organization_id,
                    "user_id": user_id,
                    "request_id": request_id
                }
            )
            raise HTTPException(
                status_code=500,
                detail=f"Error listing agents: {str(e)}"
            )


@router.delete(
    "/agents/{agent_id}",
    responses={
        400: {"model": ErrorResponse},
        401: {"model": ErrorResponse},
        404: {"model": ErrorResponse},
        500: {"model": ErrorResponse}
    },
    summary="Delete agent",
    description="Delete an agent by ID"
)
async def delete_agent(
    request: Request,
    agent_id: str,
    request_id: str = Depends(get_request_id),
    token_data: TokenPayload = Depends(get_validated_user)
):
    """
    Delete agent
    
    Args:
        request: Request object
        agent_id: Agent ID
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
        request_id=request_id,
        agent_id=agent_id
    ):
        try:
            # Delete agent
            success = await agent_manager.delete_agent(
                agent_id=agent_id,
                organization_id=organization_id,
                user_id=user_id
            )
            
            if not success:
                raise HTTPException(
                    status_code=404,
                    detail="Agent not found or you don't have permission to delete it"
                )
            
            logger.info(
                "Agent deleted",
                extra={
                    "agent_id": agent_id,
                    "organization_id": organization_id,
                    "user_id": user_id,
                    "request_id": request_id
                }
            )
            
            return {"success": True, "message": "Agent deleted successfully"}
        
        except HTTPException:
            raise
        
        except Exception as e:
            logger.error(
                f"Delete agent error: {str(e)}",
                extra={
                    "agent_id": agent_id,
                    "organization_id": organization_id,
                    "user_id": user_id,
                    "request_id": request_id
                }
            )
            raise HTTPException(
                status_code=500,
                detail=f"Error deleting agent: {str(e)}"
            )


@router.put(
    "/agents/{agent_id}",
    response_model=AgentResponse,
    responses={
        400: {"model": ErrorResponse},
        401: {"model": ErrorResponse},
        404: {"model": ErrorResponse},
        500: {"model": ErrorResponse}
    },
    summary="Update agent",
    description="Update an agent by ID"
)
async def update_agent(
    request: Request,
    agent_id: str,
    data: AgentCreateRequest,
    request_id: str = Depends(get_request_id),
    token_data: TokenPayload = Depends(get_validated_user)
):
    """
    Update agent
    
    Args:
        request: Request object
        agent_id: Agent ID
        data: Agent update data
        request_id: Request ID
        token_data: Token payload with user and organization info
        
    Returns:
        Updated agent
    """
    # Extract request parameters
    organization_id = token_data.organization_id
    user_id = token_data.user_id
    
    with LoggingContext(
        organization_id=organization_id,
        user_id=user_id,
        request_id=request_id,
        agent_id=agent_id
    ):
        try:
            # Update agent
            updated_agent = await agent_manager.update_agent(
                agent_id=agent_id,
                organization_id=organization_id,
                user_id=user_id,
                name=data.name,
                instructions=data.instructions,
                description=data.description,
                document_ids=data.document_ids,
                document_filter=data.document_filter,
                metadata=data.metadata
            )
            
            if not updated_agent:
                raise HTTPException(
                    status_code=404,
                    detail="Agent not found or you don't have permission to update it"
                )
            
            logger.info(
                "Agent updated",
                extra={
                    "agent_id": agent_id,
                    "agent_name": data.name,
                    "organization_id": organization_id,
                    "user_id": user_id,
                    "request_id": request_id
                }
            )
            
            # Convert to response model
            from datetime import datetime
            return AgentResponse(
                id=updated_agent["id"],
                name=updated_agent["name"],
                description=updated_agent["description"],
                instructions=updated_agent["instructions"],
                created_at=datetime.fromisoformat(updated_agent["created_at"]) if isinstance(updated_agent["created_at"], str) else updated_agent["created_at"],
                updated_at=datetime.fromisoformat(updated_agent["updated_at"]) if updated_agent.get("updated_at") and isinstance(updated_agent["updated_at"], str) else updated_agent.get("updated_at"),
                metadata=updated_agent.get("metadata", {}),
                document_count=updated_agent.get("document_count", 0)
            )
        
        except HTTPException:
            raise
        
        except Exception as e:
            logger.error(
                f"Update agent error: {str(e)}",
                extra={
                    "agent_id": agent_id,
                    "organization_id": organization_id,
                    "user_id": user_id,
                    "request_id": request_id
                }
            )
            raise HTTPException(
                status_code=500,
                detail=f"Error updating agent: {str(e)}"
            )


@router.post(
    "/agents/{agent_id}/chat",
    response_model=AgentChatResponse,
    responses={
        400: {"model": ErrorResponse},
        401: {"model": ErrorResponse},
        404: {"model": ErrorResponse},
        500: {"model": ErrorResponse}
    },
    summary="Chat with agent",
    description="Chat with an AI agent"
)
async def chat_with_agent(
    request: Request,
    agent_id: str,
    data: AgentChatRequest,
    request_id: str = Depends(get_request_id),
    token_data: TokenPayload = Depends(get_validated_user),
    llm_provider: BaseLLMProvider = Depends(get_llm_for_request)
):
    """
    Chat with agent
    
    Args:
        request: Request object
        agent_id: Agent ID
        data: Agent chat request data
        request_id: Request ID
        token_data: Token payload with user and organization info
        llm_provider: LLM provider
        
    Returns:
        Agent chat response
    """
    # Extract request parameters
    organization_id = token_data.organization_id
    user_id = token_data.user_id
    
    with LoggingContext(
        organization_id=organization_id,
        user_id=user_id,
        request_id=request_id,
        agent_id=agent_id
    ):
        try:
            # Check if agent exists
            agent_info = await agent_manager.get_agent(
                agent_id=agent_id,
                organization_id=organization_id,
                user_id=user_id
            )
            
            if not agent_info:
                raise HTTPException(
                    status_code=404,
                    detail="Agent not found"
                )
            
            # Check if streaming is requested
            if data.stream:
                # Create streaming response
                async def chat_stream():
                    try:
                        # Convert messages to dict format
                        messages = [msg.dict() for msg in data.messages]
                        
                        # Chat with agent (streaming)
                        stream = await agent_manager.chat_with_agent(
                            agent_id=agent_id,
                            organization_id=organization_id,
                            user_id=user_id,
                            messages=messages,
                            temperature=data.temperature,
                            max_tokens=data.max_tokens,
                            stream=True,
                            **(data.model_params or {})
                        )
                        
                        # Stream the response
                        async for chunk in stream:
                            # Yield chunk as server-sent event
                            yield f"data: {chunk}\n\n"
                        
                        # End of stream
                        yield "data: [DONE]\n\n"
                    except Exception as e:
                        logger.error(
                            f"Agent chat streaming error: {str(e)}",
                            extra={
                                "agent_id": agent_id,
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
            # Convert messages to dict format
            messages = [msg.dict() for msg in data.messages]
            
            # Chat with agent
            response = await agent_manager.chat_with_agent(
                agent_id=agent_id,
                organization_id=organization_id,
                user_id=user_id,
                messages=messages,
                temperature=data.temperature,
                max_tokens=data.max_tokens,
                stream=False,
                **(data.model_params or {})
            )
            
            logger.info(
                "Agent chat completed",
                extra={
                    "agent_id": agent_id,
                    "message_count": len(messages),
                    "response_length": len(response),
                    "organization_id": organization_id,
                    "user_id": user_id,
                    "request_id": request_id,
                    "model": llm_provider.get_model_name(),
                    "provider": llm_provider.get_provider_name()
                }
            )
            
            return AgentChatResponse(
                message=response,
                agent_id=agent_id,
                agent_name=agent_info["name"],
                model=llm_provider.get_model_name(),
                provider=llm_provider.get_provider_name(),
                finish_reason="stop",  # Placeholder
                usage=None  # Placeholder
            )
        
        except HTTPException:
            raise
        
        except Exception as e:
            logger.error(
                f"Agent chat error: {str(e)}",
                extra={
                    "agent_id": agent_id,
                    "organization_id": organization_id,
                    "user_id": user_id,
                    "request_id": request_id
                }
            )
            raise HTTPException(
                status_code=500,
                detail=f"Error chatting with agent: {str(e)}"
            )