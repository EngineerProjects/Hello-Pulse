"""
Vector database package
"""

from typing import Dict, Optional, Type
import importlib

from core.config import settings
from core.logging import logger
from .base import BaseVectorDatabase
# Import specific vector database implementations
from .chroma import ChromaDBProvider
from .qdrant import QdrantProvider

# Vector database registry
_VECTOR_DB_REGISTRY: Dict[str, Type[BaseVectorDatabase]] = {
    "chroma": ChromaDBProvider,
    "qdrant": QdrantProvider,
}

# Vector database instance cache
_VECTOR_DB_INSTANCES: Dict[str, BaseVectorDatabase] = {}

def register_vector_db(name: str, vector_db_class: Type[BaseVectorDatabase]) -> None:
    """
    Register a new vector database
    
    Args:
        name: Vector database name
        vector_db_class: Vector database class
    """
    _VECTOR_DB_REGISTRY[name] = vector_db_class
    logger.info(f"Registered vector database: {name}")

def get_vector_db(organization_id: str) -> BaseVectorDatabase:
    """
    Get vector database for an organization
    
    Args:
        organization_id: Organization ID
        
    Returns:
        Vector database instance
    
    Raises:
        ValueError: If vector database not found
    """
    # Get vector database name from company configuration
    company_config = settings.get_company_config(organization_id)
    vector_db_name = company_config.get("vector_db", settings.DEFAULT_VECTOR_DB)
    
    # Check if we should use a shared vector database
    shared_vector_db = company_config.get("shared_vector_db", settings.SHARED_VECTOR_DB)
    
    # Create cache key based on organization and whether it's shared
    if shared_vector_db:
        cache_key = f"shared:{vector_db_name}"
    else:
        cache_key = f"{organization_id}:{vector_db_name}"
    
    # Return from cache if exists
    if cache_key in _VECTOR_DB_INSTANCES:
        return _VECTOR_DB_INSTANCES[cache_key]
    
    # Get vector database class
    if vector_db_name not in _VECTOR_DB_REGISTRY:
        logger.error(f"Vector database not found: {vector_db_name}")
        # Fall back to default
        vector_db_name = settings.DEFAULT_VECTOR_DB
    
    vector_db_class = _VECTOR_DB_REGISTRY[vector_db_name]
    
    # Create vector database instance
    vector_db_instance = vector_db_class(organization_id, company_config, shared_vector_db)
    
    # Cache the instance
    _VECTOR_DB_INSTANCES[cache_key] = vector_db_instance
    
    logger.info(
        f"Created vector database instance",
        extra={
            "provider": vector_db_name,
            "organization_id": organization_id,
            "shared": shared_vector_db,
            "collection": vector_db_instance.get_collection_name()
        }
    )
    
    return vector_db_instance

def clear_vector_db_cache() -> None:
    """Clear vector database cache"""
    _VECTOR_DB_INSTANCES.clear()
    logger.info("Cleared vector database cache")

# Function to dynamically load and register additional vector databases
def load_vector_db_module(module_path: str) -> None:
    """
    Dynamically load a vector database module
    
    Args:
        module_path: Import path to the module
    """
    try:
        module = importlib.import_module(module_path)
        for attr_name in dir(module):
            attr = getattr(module, attr_name)
            # Check if it's a class that inherits from BaseVectorDatabase and isn't BaseVectorDatabase itself
            if (isinstance(attr, type) and 
                issubclass(attr, BaseVectorDatabase) and 
                attr is not BaseVectorDatabase):
                # Get provider name from the class's method
                provider_name = attr(None, {}, True).get_provider_name()
                register_vector_db(provider_name, attr)
    except Exception as e:
        logger.error(f"Failed to load vector database module {module_path}: {str(e)}")