"""
Ollama provider implementation
"""

import json
import aiohttp
from typing import Dict, List, Optional, Any, AsyncIterator
import asyncio

from core.config import settings
from core.logging import logger
from .base import BaseLLMProvider

class OllamaProvider(BaseLLMProvider):
    """Ollama API provider implementation"""
    
    def __init__(self, organization_id: str, company_config: Dict[str, Any]):
        """
        Initialize Ollama provider
        
        Args:
            organization_id: Organization ID
            company_config: Company-specific configuration
        """
        self.organization_id = organization_id
        self.model = company_config.get("ollama_model", settings.OLLAMA_DEFAULT_MODEL)
        self.api_base = company_config.get("ollama_api_base", settings.OLLAMA_API_BASE)
        
        logger.info(
            f"Initialized Ollama provider",
            extra={
                "model": self.model,
                "api_base": self.api_base,
                "organization_id": organization_id
            }
        )
    
    async def _make_request(self, endpoint: str, data: Dict[str, Any]) -> Dict[str, Any]:
        """
        Make a request to the Ollama API
        
        Args:
            endpoint: API endpoint
            data: Request data
            
        Returns:
            Response data
        """
        url = f"{self.api_base}/{endpoint}"
        
        async with aiohttp.ClientSession() as session:
            async with session.post(url, json=data) as response:
                if response.status != 200:
                    error_text = await response.text()
                    logger.error(
                        f"Ollama API error: {error_text}",
                        extra={
                            "status_code": response.status,
                            "endpoint": endpoint,
                            "model": self.model,
                            "organization_id": self.organization_id
                        }
                    )
                    raise ValueError(f"Ollama API error: {error_text}")
                
                return await response.json()
    
    async def _stream_request(self, endpoint: str, data: Dict[str, Any]) -> AsyncIterator[Dict[str, Any]]:
        """
        Stream a request from the Ollama API
        
        Args:
            endpoint: API endpoint
            data: Request data
            
        Returns:
            AsyncIterator of response chunks
        """
        url = f"{self.api_base}/{endpoint}"
        data["stream"] = True
        
        async with aiohttp.ClientSession() as session:
            async with session.post(url, json=data) as response:
                if response.status != 200:
                    error_text = await response.text()
                    logger.error(
                        f"Ollama API streaming error: {error_text}",
                        extra={
                            "status_code": response.status,
                            "endpoint": endpoint,
                            "model": self.model,
                            "organization_id": self.organization_id
                        }
                    )
                    raise ValueError(f"Ollama API error: {error_text}")
                
                # Process the streaming response
                buffer = ""
                async for chunk in response.content:
                    buffer += chunk.decode("utf-8")
                    
                    # Split by newlines and process complete JSON objects
                    lines = buffer.split("\n")
                    
                    # Keep the last potentially incomplete line
                    buffer = lines[-1]
                    
                    # Process all complete lines
                    for line in lines[:-1]:
                        if line.strip():
                            try:
                                yield json.loads(line)
                            except json.JSONDecodeError:
                                logger.warning(f"Invalid JSON in Ollama response: {line}")
    
    async def generate(
        self,
        prompt: str,
        system_message: Optional[str] = None,
        temperature: float = 0.7,
        max_tokens: Optional[int] = None,
        stop_sequences: Optional[List[str]] = None,
        **kwargs: Any
    ) -> str:
        """Generate text with Ollama"""
        data = {
            "model": self.model,
            "prompt": prompt,
            "temperature": temperature,
        }
        
        if system_message:
            data["system"] = system_message
            
        if max_tokens:
            data["num_predict"] = max_tokens
            
        if stop_sequences:
            data["stop"] = stop_sequences
            
        # Add any additional parameters
        data.update({k: v for k, v in kwargs.items() if v is not None})
        
        try:
            response = await self._make_request("generate", data)
            return response.get("response", "")
        except Exception as e:
            logger.error(
                f"Ollama generation error: {str(e)}",
                extra={
                    "model": self.model,
                    "organization_id": self.organization_id
                }
            )
            raise
    
    async def generate_stream(
        self,
        prompt: str,
        system_message: Optional[str] = None,
        temperature: float = 0.7,
        max_tokens: Optional[int] = None,
        stop_sequences: Optional[List[str]] = None,
        **kwargs: Any
    ) -> AsyncIterator[str]:
        """Stream text generation with Ollama"""
        data = {
            "model": self.model,
            "prompt": prompt,
            "temperature": temperature,
        }
        
        if system_message:
            data["system"] = system_message
            
        if max_tokens:
            data["num_predict"] = max_tokens
            
        if stop_sequences:
            data["stop"] = stop_sequences
            
        # Add any additional parameters
        data.update({k: v for k, v in kwargs.items() if v is not None})
        
        try:
            async for chunk in self._stream_request("generate", data):
                if "response" in chunk:
                    yield chunk["response"]
        except Exception as e:
            logger.error(
                f"Ollama streaming error: {str(e)}",
                extra={
                    "model": self.model,
                    "organization_id": self.organization_id
                }
            )
            raise
    
    async def chat(
        self,
        messages: List[Dict[str, str]],
        temperature: float = 0.7,
        max_tokens: Optional[int] = None,
        stop_sequences: Optional[List[str]] = None,
        **kwargs: Any
    ) -> str:
        """Generate chat response with Ollama"""
        # Convert messages to Ollama format
        ollama_messages = []
        
        for message in messages:
            role = message["role"]
            content = message["content"]
            
            if role == "system":
                # System messages are handled differently in Ollama
                continue
            
            ollama_messages.append({
                "role": role,
                "content": content
            })
        
        # Find system message if it exists
        system_message = next((m["content"] for m in messages if m["role"] == "system"), None)
        
        data = {
            "model": self.model,
            "messages": ollama_messages,
            "temperature": temperature,
        }
        
        if system_message:
            data["system"] = system_message
            
        if max_tokens:
            data["num_predict"] = max_tokens
            
        if stop_sequences:
            data["stop"] = stop_sequences
            
        # Add any additional parameters
        data.update({k: v for k, v in kwargs.items() if v is not None})
        
        try:
            response = await self._make_request("chat", data)
            return response.get("message", {}).get("content", "")
        except Exception as e:
            logger.error(
                f"Ollama chat error: {str(e)}",
                extra={
                    "model": self.model,
                    "organization_id": self.organization_id
                }
            )
            raise
    
    async def chat_stream(
        self,
        messages: List[Dict[str, str]],
        temperature: float = 0.7,
        max_tokens: Optional[int] = None,
        stop_sequences: Optional[List[str]] = None,
        **kwargs: Any
    ) -> AsyncIterator[str]:
        """Stream chat response with Ollama"""
        # Convert messages to Ollama format
        ollama_messages = []
        
        for message in messages:
            role = message["role"]
            content = message["content"]
            
            if role == "system":
                # System messages are handled differently in Ollama
                continue
            
            ollama_messages.append({
                "role": role,
                "content": content
            })
        
        # Find system message if it exists
        system_message = next((m["content"] for m in messages if m["role"] == "system"), None)
        
        data = {
            "model": self.model,
            "messages": ollama_messages,
            "temperature": temperature,
        }
        
        if system_message:
            data["system"] = system_message
            
        if max_tokens:
            data["num_predict"] = max_tokens
            
        if stop_sequences:
            data["stop"] = stop_sequences
            
        # Add any additional parameters
        data.update({k: v for k, v in kwargs.items() if v is not None})
        
        try:
            async for chunk in self._stream_request("chat", data):
                if "message" in chunk and "content" in chunk["message"]:
                    yield chunk["message"]["content"]
        except Exception as e:
            logger.error(
                f"Ollama chat streaming error: {str(e)}",
                extra={
                    "model": self.model,
                    "organization_id": self.organization_id
                }
            )
            raise
    
    async def get_embeddings(
        self,
        texts: List[str],
        **kwargs: Any
    ) -> List[List[float]]:
        """Generate embeddings with Ollama"""
        # Process each text individually as Ollama doesn't support batch embedding
        embeddings = []
        
        for text in texts:
            try:
                data = {
                    "model": self.model,
                    "prompt": text,
                }
                
                # Add any additional parameters
                data.update({k: v for k, v in kwargs.items() if v is not None})
                
                response = await self._make_request("embeddings", data)
                embeddings.append(response.get("embedding", []))
                
                # Sleep a bit to avoid overwhelming the API
                await asyncio.sleep(0.1)
            except Exception as e:
                logger.error(
                    f"Ollama embeddings error for text: {str(e)}",
                    extra={
                        "model": self.model,
                        "organization_id": self.organization_id
                    }
                )
                raise
        
        return embeddings
    
    def get_provider_name(self) -> str:
        """Get provider name"""
        return "ollama"
    
    def get_model_name(self) -> str:
        """Get model name"""
        return self.model