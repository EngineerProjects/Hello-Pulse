"""
Filtering utilities for secure access control
"""

from typing import Dict, List, Any, Optional
from enum import Enum

from core.logging import logger

class Visibility(str, Enum):
    """Document visibility levels"""
    PRIVATE = "private"  # Only accessible to the uploader
    SHARED = "shared"    # Accessible to all users in the same organization
    PUBLIC = "public"    # Accessible to anyone (across organizations)


def get_access_filters(
    organization_id: str,
    user_id: str,
    visibility: Optional[List[str]] = None,
    **additional_filters: Any
) -> Dict[str, Any]:
    """
    Create secure access filters for vector database queries
    
    Args:
        organization_id: Organization ID
        user_id: User ID
        visibility: Optional list of visibility levels to include
                   (defaults to all levels with appropriate filtering)
        **additional_filters: Additional filters to include
    
    Returns:
        Dict of filter conditions that ensure the user only accesses
        documents they're allowed to see
    """
    # Start with the base filter - always filter by organization
    filters = {"organization_id": organization_id}
    
    # Add visibility filters
    if visibility is None:
        # Default visibility filtering:
        # - User can see their own documents (private)
        # - User can see shared documents from their organization
        # - User can see public documents
        visibility_filter = [
            Visibility.SHARED.value,  # Shared within organization
            Visibility.PUBLIC.value   # Public to everyone
        ]
        
        # User can always see their own documents regardless of visibility
        filters["$or"] = [
            {"user_id": user_id},                    # User's own documents
            {"visibility": {"$in": visibility_filter}}  # Shared/public docs
        ]
    else:
        # Custom visibility filtering
        if Visibility.PRIVATE.value in visibility:
            # For private, only show user's own documents
            if len(visibility) == 1:
                # Only private, just filter by user_id
                filters["user_id"] = user_id
            else:
                # Mix of private and other visibilities
                private_filter = {"user_id": user_id}
                other_filter = {"visibility": {"$in": [v for v in visibility if v != Visibility.PRIVATE.value]}}
                filters["$or"] = [private_filter, other_filter]
        else:
            # No private documents, just filter by visibility
            filters["visibility"] = {"$in": visibility}
    
    # Add any additional filters
    for key, value in additional_filters.items():
        filters[key] = value
    
    logger.debug(
        f"Created access filters",
        extra={
            "organization_id": organization_id,
            "user_id": user_id,
            "filters": filters
        }
    )
    
    return filters


def document_access_check(
    document_metadata: Dict[str, Any],
    user_id: str,
    organization_id: str
) -> bool:
    """
    Check if a user has access to a specific document
    
    Args:
        document_metadata: Document metadata
        user_id: User ID
        organization_id: Organization ID
        
    Returns:
        Boolean indicating whether the user has access
    """
    # Check organization first
    doc_org_id = document_metadata.get("organization_id")
    if doc_org_id != organization_id:
        # Some documents might be public across organizations
        if document_metadata.get("visibility") == Visibility.PUBLIC.value:
            return True
        return False
    
    # Check user and visibility
    doc_user_id = document_metadata.get("user_id")
    doc_visibility = document_metadata.get("visibility", Visibility.PRIVATE.value)
    
    # User can always access their own documents
    if doc_user_id == user_id:
        return True
    
    # Otherwise, check visibility
    if doc_visibility == Visibility.PRIVATE.value:
        return False  # Private documents are only for the owner
    
    if doc_visibility == Visibility.SHARED.value:
        return True  # Shared documents are accessible to all users in the organization
    
    if doc_visibility == Visibility.PUBLIC.value:
        return True  # Public documents are accessible to everyone
    
    # Default deny
    return False


def sanitize_metadata_for_storage(
    metadata: Dict[str, Any],
    organization_id: str,
    user_id: str
) -> Dict[str, Any]:
    """
    Sanitize and prepare metadata for storage
    - Ensures required security fields are present
    - Validates visibility levels
    
    Args:
        metadata: Raw metadata
        organization_id: Organization ID
        user_id: User ID
        
    Returns:
        Sanitized metadata ready for storage
    """
    # Create a copy to avoid modifying the original
    sanitized = dict(metadata)
    
    # Ensure organization_id and user_id are set
    sanitized["organization_id"] = organization_id
    sanitized["user_id"] = user_id
    
    # Validate visibility
    if "visibility" not in sanitized:
        sanitized["visibility"] = Visibility.PRIVATE.value
    elif sanitized["visibility"] not in [v.value for v in Visibility]:
        logger.warning(
            f"Invalid visibility value, defaulting to private",
            extra={
                "organization_id": organization_id,
                "user_id": user_id,
                "invalid_visibility": sanitized["visibility"]
            }
        )
        sanitized["visibility"] = Visibility.PRIVATE.value
    
    # Add timestamp if not present
    if "timestamp" not in sanitized:
        import time
        sanitized["timestamp"] = int(time.time())
    
    return sanitized