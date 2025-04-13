"""
Base interface for LLM providers
"""

from abc import ABC, abstractmethod
from typing import Dict, List, Optional, Any, AsyncIterator, Union

class BaseLLMProvider(ABC):
    """
    Base interface for LLM providers
    All LLM integration implementations must inherit from this class
    """
    
    @abstractmethod
    async def generate(
        self,
        prompt: str,
        system_message: Optional[str] = None,
        temperature: float = 0.7,
        max_tokens: Optional[int] = None,
        stop_sequences: Optional[List[str]] = None,
        **kwargs: Any
    ) -> str:
        """
        Generate text completion from prompt
        
        Args:
            prompt: The user prompt to generate from
            system_message: Optional system message for chat models
            temperature: Controls randomness (0.0 = deterministic, 1.0 = random)
            max_tokens: Maximum number of tokens to generate
            stop_sequences: List of sequences where generation should stop
            **kwargs: Additional provider-specific parameters
            
        Returns:
            Generated text as string
        """
        pass
    
    @abstractmethod
    async def generate_stream(
        self,
        prompt: str,
        system_message: Optional[str] = None,
        temperature: float = 0.7,
        max_tokens: Optional[int] = None,
        stop_sequences: Optional[List[str]] = None,
        **kwargs: Any
    ) -> AsyncIterator[str]:
        """
        Stream text completion from prompt
        
        Args:
            prompt: The user prompt to generate from
            system_message: Optional system message for chat models
            temperature: Controls randomness (0.0 = deterministic, 1.0 = random)
            max_tokens: Maximum number of tokens to generate
            stop_sequences: List of sequences where generation should stop
            **kwargs: Additional provider-specific parameters
            
        Returns:
            AsyncIterator yielding chunks of generated text
        """
        pass
    
    @abstractmethod
    async def chat(
        self,
        messages: List[Dict[str, str]],
        temperature: float = 0.7,
        max_tokens: Optional[int] = None,
        stop_sequences: Optional[List[str]] = None,
        **kwargs: Any
    ) -> str:
        """
        Generate response from a chat history
        
        Args:
            messages: List of message dictionaries with 'role' and 'content'
                     roles can be 'system', 'user', or 'assistant'
            temperature: Controls randomness (0.0 = deterministic, 1.0 = random)
            max_tokens: Maximum number of tokens to generate
            stop_sequences: List of sequences where generation should stop
            **kwargs: Additional provider-specific parameters
            
        Returns:
            Generated response as string
        """
        pass
    
    @abstractmethod
    async def chat_stream(
        self,
        messages: List[Dict[str, str]],
        temperature: float = 0.7,
        max_tokens: Optional[int] = None,
        stop_sequences: Optional[List[str]] = None,
        **kwargs: Any
    ) -> AsyncIterator[str]:
        """
        Stream response from a chat history
        
        Args:
            messages: List of message dictionaries with 'role' and 'content'
                     roles can be 'system', 'user', or 'assistant'
            temperature: Controls randomness (0.0 = deterministic, 1.0 = random)
            max_tokens: Maximum number of tokens to generate
            stop_sequences: List of sequences where generation should stop
            **kwargs: Additional provider-specific parameters
            
        Returns:
            AsyncIterator yielding chunks of generated text
        """
        pass
    
    @abstractmethod
    async def get_embeddings(
        self,
        texts: List[str],
        **kwargs: Any
    ) -> List[List[float]]:
        """
        Generate embeddings for a list of texts
        
        Args:
            texts: List of text strings to embed
            **kwargs: Additional provider-specific parameters
            
        Returns:
            List of embedding vectors as lists of floats
        """
        pass
    
    @abstractmethod
    def get_provider_name(self) -> str:
        """
        Get the name of this LLM provider
        
        Returns:
            Provider name as string
        """
        pass
    
    @abstractmethod
    def get_model_name(self) -> str:
        """
        Get the name of the model being used
        
        Returns:
            Model name as string
        """
        pass