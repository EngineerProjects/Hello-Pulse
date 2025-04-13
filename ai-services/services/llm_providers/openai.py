"""
OpenAI provider implementation
"""

from typing import Dict, List, Optional, Any, AsyncIterator, Union
import asyncio

from openai import AsyncOpenAI
from openai.types.chat import ChatCompletionMessageParam
import tiktoken

from core.config import settings
from core.logging import logger
from .base import BaseLLMProvider

class OpenAIProvider(BaseLLMProvider):
    """OpenAI API provider implementation"""
    
    def __init__(self, organization_id: str, company_config: Dict[str, Any]):
        """
        Initialize OpenAI provider
        
        Args:
            organization_id: Organization ID
            company_config: Company-specific configuration
        """
        self.organization_id = organization_id
        self.model = company_config.get("openai_model", settings.OPENAI_DEFAULT_MODEL)
        self.embedding_model = company_config.get(
            "openai_embedding_model", settings.OPENAI_EMBEDDING_MODEL
        )
        
        # Get API key - try company-specific one first, then fall back to global
        api_key = company_config.get("openai_api_key", settings.OPENAI_API_KEY)
        if not api_key:
            raise ValueError("OpenAI API key not configured")
        
        # Initialize client
        self.client = AsyncOpenAI(api_key=api_key)
        
        # Get tokenizer for token counting
        try:
            self.tokenizer = tiktoken.encoding_for_model(self.model)
        except Exception:
            # Fall back to cl100k for newer models that may not be in tiktoken yet
            self.tokenizer = tiktoken.get_encoding("cl100k_base")
        
        logger.info(
            f"Initialized OpenAI provider",
            extra={
                "model": self.model,
                "embedding_model": self.embedding_model,
                "organization_id": organization_id
            }
        )
    
    async def generate(
        self,
        prompt: str,
        system_message: Optional[str] = None,
        temperature: float = 0.7,
        max_tokens: Optional[int] = None,
        stop_sequences: Optional[List[str]] = None,
        **kwargs: Any
    ) -> str:
        """Generate text with OpenAI"""
        # For OpenAI, we use the chat interface for all generation
        messages = []
        
        if system_message:
            messages.append({"role": "system", "content": system_message})
        
        messages.append({"role": "user", "content": prompt})
        
        try:
            response = await self.client.chat.completions.create(
                model=self.model,
                messages=[ChatCompletionMessageParam(**msg) for msg in messages],
                temperature=temperature,
                max_tokens=max_tokens,
                stop=stop_sequences,
                **kwargs
            )
            
            return response.choices[0].message.content or ""
        except Exception as e:
            logger.error(
                f"OpenAI generation error: {str(e)}",
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
        """Stream text generation with OpenAI"""
        messages = []
        
        if system_message:
            messages.append({"role": "system", "content": system_message})
        
        messages.append({"role": "user", "content": prompt})
        
        try:
            stream = await self.client.chat.completions.create(
                model=self.model,
                messages=[ChatCompletionMessageParam(**msg) for msg in messages],
                temperature=temperature,
                max_tokens=max_tokens,
                stop=stop_sequences,
                stream=True,
                **kwargs
            )
            
            async for chunk in stream:
                if chunk.choices and chunk.choices[0].delta.content:
                    yield chunk.choices[0].delta.content
        except Exception as e:
            logger.error(
                f"OpenAI streaming error: {str(e)}",
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
        """Generate chat response with OpenAI"""
        try:
            response = await self.client.chat.completions.create(
                model=self.model,
                messages=[ChatCompletionMessageParam(**msg) for msg in messages],
                temperature=temperature,
                max_tokens=max_tokens,
                stop=stop_sequences,
                **kwargs
            )
            
            return response.choices[0].message.content or ""
        except Exception as e:
            logger.error(
                f"OpenAI chat error: {str(e)}",
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
        """Stream chat response with OpenAI"""
        try:
            stream = await self.client.chat.completions.create(
                model=self.model,
                messages=[ChatCompletionMessageParam(**msg) for msg in messages],
                temperature=temperature,
                max_tokens=max_tokens,
                stop=stop_sequences,
                stream=True,
                **kwargs
            )
            
            async for chunk in stream:
                if chunk.choices and chunk.choices[0].delta.content:
                    yield chunk.choices[0].delta.content
        except Exception as e:
            logger.error(
                f"OpenAI chat streaming error: {str(e)}",
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
        """Generate embeddings with OpenAI"""
        try:
            # Process in batches of 100 to avoid rate limits
            batch_size = 100
            all_embeddings = []
            
            for i in range(0, len(texts), batch_size):
                batch_texts = texts[i:i+batch_size]
                
                response = await self.client.embeddings.create(
                    model=self.embedding_model,
                    input=batch_texts,
                    **kwargs
                )
                
                # Extract embeddings and ensure they're in the same order as input
                sorted_embeddings = sorted(
                    response.data, 
                    key=lambda x: x.index
                )
                
                batch_embeddings = [item.embedding for item in sorted_embeddings]
                all_embeddings.extend(batch_embeddings)
                
                # Sleep a bit to avoid rate limits if there are more batches
                if i + batch_size < len(texts):
                    await asyncio.sleep(0.5)
            
            return all_embeddings
        except Exception as e:
            logger.error(
                f"OpenAI embeddings error: {str(e)}",
                extra={
                    "model": self.embedding_model,
                    "organization_id": self.organization_id,
                    "text_count": len(texts)
                }
            )
            raise
    
    def get_provider_name(self) -> str:
        """Get provider name"""
        return "openai"
    
    def get_model_name(self) -> str:
        """Get model name"""
        return self.model