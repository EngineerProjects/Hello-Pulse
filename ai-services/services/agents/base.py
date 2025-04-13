"""
Base interface for AI agents
"""

from abc import ABC, abstractmethod
from typing import Dict, List, Optional, Any, AsyncIterator, Union


class BaseAgent(ABC):
    """
    Base interface for AI agents
    All agent implementations must inherit from this class
    """
    
    @property
    @abstractmethod
    def agent_id(self) -> str:
        """Get agent ID"""
        pass
    
    @property
    @abstractmethod
    def agent_name(self) -> str:
        """Get agent name"""
        pass
    
    @property
    @abstractmethod
    def agent_description(self) -> Optional[str]:
        """Get agent description"""
        pass
    
    @property
    @abstractmethod
    def agent_instructions(self) -> str:
        """Get agent instructions"""
        pass
    
    @property
    @abstractmethod
    def agent_metadata(self) -> Dict[str, Any]:
        """Get agent metadata"""
        pass
    
    @abstractmethod
    async def chat(
        self,
        messages: List[Dict[str, str]],
        temperature: float = 0.7,
        max_tokens: Optional[int] = None,
        **kwargs: Any
    ) -> str:
        """
        Chat with the agent
        
        Args:
            messages: List of message dictionaries with 'role' and 'content'
            temperature: Controls randomness (0.0 = deterministic, 1.0 = random)
            max_tokens: Maximum number of tokens to generate
            **kwargs: Additional parameters
            
        Returns:
            Agent's response
        """
        pass
    
    @abstractmethod
    async def chat_stream(
        self,
        messages: List[Dict[str, str]],
        temperature: float = 0.7,
        max_tokens: Optional[int] = None,
        **kwargs: Any
    ) -> AsyncIterator[str]:
        """
        Stream a chat with the agent
        
        Args:
            messages: List of message dictionaries with 'role' and 'content'
            temperature: Controls randomness (0.0 = deterministic, 1.0 = random)
            max_tokens: Maximum number of tokens to generate
            **kwargs: Additional parameters
            
        Returns:
            AsyncIterator of response chunks
        """
        pass
    
    @abstractmethod
    async def retrieve_relevant_docs(
        self,
        query: str,
        limit: int = 5
    ) -> List[Dict[str, Any]]:
        """
        Retrieve relevant documents for this agent based on a query
        
        Args:
            query: Query text
            limit: Maximum number of documents to retrieve
            
        Returns:
            List of relevant documents
        """
        pass
    
    @abstractmethod
    def to_dict(self) -> Dict[str, Any]:
        """
        Convert agent to dictionary representation
        
        Returns:
            Dictionary representation of the agent
        """
        pass