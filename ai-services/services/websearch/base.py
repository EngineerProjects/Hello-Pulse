"""
Base interface for web search providers
"""

from abc import ABC, abstractmethod
from typing import Dict, List, Optional, Any


class WebSearchResult:
    """Web search result model"""
    
    def __init__(
        self,
        title: str,
        url: str,
        snippet: Optional[str] = None,
        source: str = "",
        position: int = 0,
        additional_info: Optional[Dict[str, Any]] = None
    ):
        self.title = title
        self.url = url
        self.snippet = snippet
        self.source = source
        self.position = position
        self.additional_info = additional_info or {}
    
    def to_dict(self) -> Dict[str, Any]:
        """Convert to dictionary"""
        return {
            "title": self.title,
            "url": self.url,
            "snippet": self.snippet,
            "source": self.source,
            "position": self.position,
            **{f"info_{k}": v for k, v in self.additional_info.items()}
        }


class BaseWebSearchProvider(ABC):
    """
    Base interface for web search providers
    All web search implementations must inherit from this class
    """
    
    @abstractmethod
    async def search(
        self,
        query: str,
        num_results: int = 5,
        **kwargs: Any
    ) -> List[WebSearchResult]:
        """
        Search the web
        
        Args:
            query: Search query
            num_results: Number of results to return
            **kwargs: Additional provider-specific parameters
            
        Returns:
            List of search results
        """
        pass
    
    @abstractmethod
    def get_provider_name(self) -> str:
        """
        Get the name of this web search provider
        
        Returns:
            Provider name as string
        """
        pass