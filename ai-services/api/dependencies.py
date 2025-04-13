"""
Dependencies for API routes
"""

import uuid
from typing import Annotated, Dict, Any, Optional

from fastapi import Depends, HTTPException, Header, status, Request

from core.security import get_current_user, TokenPayload, validate_organization_access
from core.logging import LoggingContext, logger
from services.llm_providers.base import BaseLLMProvider
from services.llm_providers import get_llm_provider
from services.vector_databases.base import BaseVectorDatabase
from services.vector_databases import get_vector_db
from utils.filtering import get_access_filters

# Create request ID dependency
async def get_request_id(x_request_id: Optional[str] = Header(None)) -> str:
    """Get or generate request ID"""
    if x_request_id is None:
        return str(uuid.uuid4())
    return x_request_id

# Enhanced current user dependency with additional validation
async def get_validated_user(
    request: Request,
    request_id: str = Depends(get_request_id),
    token_data: TokenPayload = Depends(get_current_user),
) -> TokenPayload:
    """Get current user with additional validation and logging context"""
    # Set up logging context with user and organization info
    with LoggingContext(
        organization_id=token_data.organization_id,
        user_id=token_data.user_id,
        request_id=request_id
    ):
        logger.info(
            f"Authenticated request",
            extra={
                "path": request.url.path,
                "method": request.method,
                "organization_id": token_data.organization_id,
                "user_id": token_data.user_id,
            }
        )
        
        return token_data

# Get organization-specific LLM provider
async def get_llm_for_request(
    token_data: TokenPayload = Depends(get_validated_user)
) -> BaseLLMProvider:
    """Get the appropriate LLM provider for the current organization"""
    organization_id = token_data.organization_id
    try:
        # This will use the organization's configured provider
        return get_llm_provider(organization_id)
    except Exception as e:
        logger.error(f"Failed to get LLM provider: {str(e)}")
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail="Failed to initialize LLM provider"
        )

# Get organization-specific vector database
async def get_vector_db_for_request(
    token_data: TokenPayload = Depends(get_validated_user)
) -> BaseVectorDatabase:
    """Get the appropriate vector database for the current organization"""
    organization_id = token_data.organization_id
    try:
        # This will use the organization's configured vector database
        return get_vector_db(organization_id)
    except Exception as e:
        logger.error(f"Failed to get vector database: {str(e)}")
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail="Failed to initialize vector database"
        )

# Get access filters for the current user and organization
async def get_user_access_filters(
    token_data: TokenPayload = Depends(get_validated_user)
) -> Dict[str, Any]:
    """Get access filters for the current user"""
    return get_access_filters(
        organization_id=token_data.organization_id,
        user_id=token_data.user_id
    )

# Validate organization ID in request
async def validate_organization(
    organization_id: str,
    token_data: TokenPayload = Depends(get_validated_user)
) -> str:
    """Validate that the user has access to the specified organization"""
    if not validate_organization_access(token_data, organization_id):
        logger.warning(
            f"Organization access denied",
            extra={
                "user_org_id": token_data.organization_id,
                "requested_org_id": organization_id,
                "user_id": token_data.user_id,
            }
        )
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Access to this organization is forbidden"
        )
    return organization_id