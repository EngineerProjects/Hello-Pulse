"""
LLM providers package
"""

from typing import Dict, Optional, Type
import importlib

from core.config import settings
from core.logging import logger
from .base import BaseLLMProvider
# Import specific providers for registration
from .openai import OpenAIProvider
from .ollama import OllamaProvider

# LLM Provider registry
_PROVIDER_REGISTRY: Dict[str, Type[BaseLLMProvider]] = {
    "openai": OpenAIProvider,
    "ollama": OllamaProvider,
}

# LLM Provider instance cache
_PROVIDER_INSTANCES: Dict[str, BaseLLMProvider] = {}

def register_provider(name: str, provider_class: Type[BaseLLMProvider]) -> None:
    """
    Register a new LLM provider
    
    Args:
        name: Provider name
        provider_class: Provider class
    """
    _PROVIDER_REGISTRY[name] = provider_class
    logger.info(f"Registered LLM provider: {name}")

def get_llm_provider(organization_id: str) -> BaseLLMProvider:
    """
    Get LLM provider for an organization
    
    Args:
        organization_id: Organization ID
        
    Returns:
        LLM provider instance
    
    Raises:
        ValueError: If provider not found
    """
    # Get provider name from company configuration
    provider_name = settings.get_llm_provider_for_company(organization_id)
    
    # Create cache key (composite of org ID and provider name)
    cache_key = f"{organization_id}:{provider_name}"
    
    # Return from cache if exists
    if cache_key in _PROVIDER_INSTANCES:
        return _PROVIDER_INSTANCES[cache_key]
    
    # Get provider class
    if provider_name not in _PROVIDER_REGISTRY:
        logger.error(f"LLM provider not found: {provider_name}")
        # Fall back to default provider
        provider_name = settings.DEFAULT_LLM_PROVIDER
    
    provider_class = _PROVIDER_REGISTRY[provider_name]
    
    # Get company-specific configuration
    company_config = settings.get_company_config(organization_id)
    
    # Create provider instance
    provider_instance = provider_class(organization_id, company_config)
    
    # Cache the instance
    _PROVIDER_INSTANCES[cache_key] = provider_instance
    
    logger.info(
        f"Created LLM provider instance",
        extra={
            "provider": provider_name,
            "organization_id": organization_id,
            "model": provider_instance.get_model_name()
        }
    )
    
    return provider_instance

def clear_provider_cache() -> None:
    """Clear provider cache"""
    _PROVIDER_INSTANCES.clear()
    logger.info("Cleared LLM provider cache")

# Function to dynamically load and register additional providers
def load_provider_module(module_path: str) -> None:
    """
    Dynamically load a provider module
    
    Args:
        module_path: Import path to the module
    """
    try:
        module = importlib.import_module(module_path)
        for attr_name in dir(module):
            attr = getattr(module, attr_name)
            # Check if it's a class that inherits from BaseLLMProvider and isn't BaseLLMProvider itself
            if (isinstance(attr, type) and 
                issubclass(attr, BaseLLMProvider) and 
                attr is not BaseLLMProvider):
                # Get provider name from the class's method
                provider_name = attr().get_provider_name()
                register_provider(provider_name, attr)
    except Exception as e:
        logger.error(f"Failed to load provider module {module_path}: {str(e)}")