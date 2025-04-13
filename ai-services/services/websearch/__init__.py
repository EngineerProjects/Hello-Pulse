"""
Web search package
"""

from typing import Dict, Optional, Type, List, Any
import importlib

from core.config import settings
from core.logging import logger
from .base import BaseWebSearchProvider, WebSearchResult
# Import specific web search implementations
from .serpapi import SerpAPIProvider

# Web search provider registry
_SEARCH_PROVIDER_REGISTRY: Dict[str, Type[BaseWebSearchProvider]] = {
    "serpapi": SerpAPIProvider,
}

# Web search provider instance cache
_SEARCH_PROVIDER_INSTANCES: Dict[str, BaseWebSearchProvider] = {}

def register_search_provider(name: str, provider_class: Type[BaseWebSearchProvider]) -> None:
    """
    Register a new web search provider
    
    Args:
        name: Provider name
        provider_class: Provider class
    """
    _SEARCH_PROVIDER_REGISTRY[name] = provider_class
    logger.info(f"Registered web search provider: {name}")

def get_search_provider(provider_name: Optional[str] = None) -> BaseWebSearchProvider:
    """
    Get web search provider
    
    Args:
        provider_name: Provider name (defaults to settings.DEFAULT_SEARCH_PROVIDER)
        
    Returns:
        Web search provider instance
    
    Raises:
        ValueError: If provider not found
    """
    # Use default provider if not specified
    if provider_name is None:
        provider_name = "serpapi"  # Default
    
    # Return from cache if exists
    if provider_name in _SEARCH_PROVIDER_INSTANCES:
        return _SEARCH_PROVIDER_INSTANCES[provider_name]
    
    # Get provider class
    if provider_name not in _SEARCH_PROVIDER_REGISTRY:
        logger.error(f"Web search provider not found: {provider_name}")
        raise ValueError(f"Web search provider not found: {provider_name}")
    
    provider_class = _SEARCH_PROVIDER_REGISTRY[provider_name]
    
    # Create provider instance
    provider_instance = provider_class()
    
    # Cache the instance
    _SEARCH_PROVIDER_INSTANCES[provider_name] = provider_instance
    
    logger.info(f"Created web search provider instance: {provider_name}")
    
    return provider_instance

def clear_search_provider_cache() -> None:
    """Clear search provider cache"""
    _SEARCH_PROVIDER_INSTANCES.clear()
    logger.info("Cleared web search provider cache")

# Function to dynamically load and register additional search providers
def load_search_provider_module(module_path: str) -> None:
    """
    Dynamically load a search provider module
    
    Args:
        module_path: Import path to the module
    """
    try:
        module = importlib.import_module(module_path)
        for attr_name in dir(module):
            attr = getattr(module, attr_name)
            # Check if it's a class that inherits from BaseWebSearchProvider and isn't BaseWebSearchProvider itself
            if (isinstance(attr, type) and 
                issubclass(attr, BaseWebSearchProvider) and 
                attr is not BaseWebSearchProvider):
                # Get provider name from the class's method
                provider_name = attr().get_provider_name()
                register_search_provider(provider_name, attr)
    except Exception as e:
        logger.error(f"Failed to load search provider module {module_path}: {str(e)}")


async def search_and_generate(
    query: str,
    organization_id: str,
    user_id: str,
    num_results: int = 5,
    system_message: Optional[str] = None,
    temperature: float = 0.7,
    max_tokens: Optional[int] = None,
    search_provider_name: Optional[str] = None,
    **kwargs: Any
) -> Dict[str, Any]:
    """
    Search the web and generate an enhanced response
    
    Args:
        query: Search query
        organization_id: Organization ID
        user_id: User ID
        num_results: Number of search results to use
        system_message: Optional system message
        temperature: Temperature for generation
        max_tokens: Maximum tokens to generate
        search_provider_name: Optional search provider name
        **kwargs: Additional parameters
        
    Returns:
        Dictionary with generated answer and search results
    """
    from services.llm_providers import get_llm_provider
    
    try:
        # Get search provider
        search_provider = get_search_provider(search_provider_name)
        
        # Perform web search
        search_results = await search_provider.search(
            query=query,
            num_results=num_results,
            **kwargs
        )
        
        if not search_results:
            # No search results found
            return {
                "answer": "I couldn't find any relevant information for your query.",
                "search_results": []
            }
        
        # Format search results for the prompt
        search_context = format_search_results(search_results)
        
        # Get LLM provider
        llm_provider = get_llm_provider(organization_id)
        
        # Create system message
        if system_message:
            enhanced_system_message = system_message
        else:
            enhanced_system_message = (
                "You are a helpful assistant with access to web search results. "
                "Answer the user's question based on the provided search results. "
                "If the search results don't contain the necessary information to answer the question, "
                "just summarize what information is available. "
                "Include relevant URLs in your answer where appropriate, but do not list all URLs."
            )
        
        # Create prompt with search results
        search_prompt = f"""I need information about the following query:
{query}

Here are the top search results:

{search_context}

Based on these search results, please provide a comprehensive answer to the query."""
        
        # Generate response
        answer = await llm_provider.generate(
            prompt=search_prompt,
            system_message=enhanced_system_message,
            temperature=temperature,
            max_tokens=max_tokens
        )
        
        logger.info(
            f"Generated web search enhanced answer",
            extra={
                "query": query,
                "organization_id": organization_id,
                "user_id": user_id,
                "result_count": len(search_results)
            }
        )
        
        return {
            "answer": answer,
            "search_results": [r.to_dict() for r in search_results]
        }
    
    except Exception as e:
        logger.error(
            f"Error in search_and_generate: {str(e)}",
            extra={
                "query": query,
                "organization_id": organization_id,
                "user_id": user_id
            }
        )
        raise


def format_search_results(results: List[WebSearchResult]) -> str:
    """
    Format search results into a context string
    
    Args:
        results: List of search results
        
    Returns:
        Formatted context string
    """
    context_parts = []
    
    for i, result in enumerate(results):
        title = result.title
        url = result.url
        snippet = result.snippet or ""
        
        context_parts.append(f"Result {i+1}:\nTitle: {title}\nURL: {url}\nSnippet: {snippet}\n")
    
    return "\n".join(context_parts)