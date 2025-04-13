"""
SerpAPI web search provider implementation
"""

import json
import aiohttp
from typing import Dict, List, Optional, Any, Union

from core.config import settings
from core.logging import logger
from .base import BaseWebSearchProvider, WebSearchResult


class SerpAPIProvider(BaseWebSearchProvider):
    """SerpAPI provider implementation"""
    
    def __init__(self, api_key: Optional[str] = None):
        """
        Initialize SerpAPI provider
        
        Args:
            api_key: SerpAPI API key (defaults to settings if not provided)
        """
        self.api_key = api_key or settings.SERPAPI_API_KEY
        
        if not self.api_key:
            raise ValueError("SerpAPI API key not configured")
        
        self.base_url = "https://serpapi.com/search"
    
    async def search(
        self,
        query: str,
        num_results: int = 5,
        **kwargs: Any
    ) -> List[WebSearchResult]:
        """
        Search the web using SerpAPI
        
        Args:
            query: Search query
            num_results: Number of results to return
            **kwargs: Additional parameters for the search
                engine: Search engine to use (google, bing, etc.)
                country: Country code (us, uk, etc.)
                language: Language code (en, fr, etc.)
                safe: Safe search (active, off)
                time_period: Time period (last_day, last_week, last_month, last_year)
            
        Returns:
            List of search results
        """
        # Set up parameters
        params = {
            "q": query,
            "api_key": self.api_key,
            "engine": kwargs.get("engine", "google"),
            "num": min(num_results * 2, 100),  # Request more results to account for filtering
            "gl": kwargs.get("country", "us"),
            "hl": kwargs.get("language", "en"),
            "safe": kwargs.get("safe", "active"),
        }
        
        # Add time period if specified
        if "time_period" in kwargs:
            params["time_period"] = kwargs["time_period"]
        
        try:
            # Make request
            async with aiohttp.ClientSession() as session:
                async with session.get(self.base_url, params=params) as response:
                    if response.status != 200:
                        error_text = await response.text()
                        logger.error(
                            f"SerpAPI error: {error_text}",
                            extra={
                                "status_code": response.status,
                                "query": query
                            }
                        )
                        raise Exception(f"SerpAPI error: {error_text}")
                    
                    data = await response.json()
            
            # Extract results
            results = []
            
            # Process organic results
            if "organic_results" in data:
                for i, result in enumerate(data["organic_results"]):
                    if i >= num_results:
                        break
                    
                    title = result.get("title", "")
                    url = result.get("link", "")
                    snippet = result.get("snippet", "")
                    
                    # Additional info
                    additional_info = {}
                    if "displayed_link" in result:
                        additional_info["displayed_link"] = result["displayed_link"]
                    if "position" in result:
                        additional_info["position"] = result["position"]
                    if "sitelinks" in result:
                        additional_info["sitelinks"] = result["sitelinks"]
                    
                    results.append(
                        WebSearchResult(
                            title=title,
                            url=url,
                            snippet=snippet,
                            source="Google",
                            position=i,
                            additional_info=additional_info
                        )
                    )
            
            # Log search results
            logger.info(
                f"SerpAPI search completed",
                extra={
                    "query": query,
                    "result_count": len(results),
                    "engine": params["engine"]
                }
            )
            
            return results[:num_results]
        
        except Exception as e:
            logger.error(
                f"SerpAPI search error: {str(e)}",
                extra={
                    "query": query,
                    "engine": params["engine"]
                }
            )
            raise
    
    def get_provider_name(self) -> str:
        """Get provider name"""
        return "serpapi"