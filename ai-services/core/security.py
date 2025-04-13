"""
Security utilities for the Hello Pulse AI Microservice
Handles JWT validation, access control, and authentication
"""

from datetime import datetime, timedelta
from typing import Dict, Optional, Any, Union

import jwt
from fastapi import Depends, HTTPException, status
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials
from pydantic import BaseModel

from core.config import settings
from core.logging import logger

# Security scheme for API authentication
security = HTTPBearer()

class TokenPayload(BaseModel):
    """JWT token payload model"""
    organization_id: str
    user_id: str
    exp: Optional[datetime] = None
    role: Optional[str] = None

def create_access_token(
    organization_id: str,
    user_id: str,
    role: Optional[str] = None,
    expires_delta: Optional[timedelta] = None,
) -> str:
    """
    Create a new JWT access token
    """
    if expires_delta:
        expire = datetime.utcnow() + expires_delta
    else:
        expire = datetime.utcnow() + timedelta(
            minutes=settings.AUTH_TOKEN_EXPIRE_MINUTES
        )
    
    to_encode = {
        "organization_id": organization_id,
        "user_id": user_id,
        "exp": expire,
    }
    
    if role:
        to_encode["role"] = role
    
    encoded_jwt = jwt.encode(
        to_encode, 
        settings.SECRET_KEY, 
        algorithm=settings.AUTH_ALGORITHM
    )
    
    return encoded_jwt

def decode_access_token(token: str) -> TokenPayload:
    """
    Decode and validate a JWT access token
    """
    try:
        payload = jwt.decode(
            token,
            settings.SECRET_KEY,
            algorithms=[settings.AUTH_ALGORITHM]
        )
        
        return TokenPayload(**payload)
    except jwt.PyJWTError as e:
        logger.error(f"Token validation error: {str(e)}")
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Could not validate credentials",
            headers={"WWW-Authenticate": "Bearer"},
        )

async def get_current_user(
    credentials: HTTPAuthorizationCredentials = Depends(security)
) -> TokenPayload:
    """
    Dependency to get the current authenticated user from the JWT token
    """
    try:
        token = credentials.credentials
        token_data = decode_access_token(token)
        
        # Check token expiration
        if token_data.exp and datetime.utcnow() > token_data.exp:
            raise HTTPException(
                status_code=status.HTTP_401_UNAUTHORIZED,
                detail="Token expired",
                headers={"WWW-Authenticate": "Bearer"},
            )
        
        return token_data
    except Exception as e:
        logger.error(f"Authentication error: {str(e)}")
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Could not validate credentials",
            headers={"WWW-Authenticate": "Bearer"},
        )

def validate_organization_access(token_data: TokenPayload, target_org_id: str) -> bool:
    """
    Validate that the user has access to the target organization
    """
    # Simple check - user must belong to the organization they're trying to access
    # In a real system, this might be more complex with roles and permissions
    return token_data.organization_id == target_org_id