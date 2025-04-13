"""
Personalized agent implementations
"""

from typing import Dict, List, Optional, Any, AsyncIterator, Union
from datetime import datetime

from core.logging import logger, LoggingContext
from services.llm_providers import get_llm_provider
from services.vector_databases import get_vector_db
from services.rag.retriever import retrieve_relevant_documents
from utils.filtering import get_access_filters
from .base import BaseAgent


class RAGAgent(BaseAgent):
    """
    Retrieval-Augmented Generation Agent
    This agent uses RAG to enhance its responses with knowledge from documents
    """
    
    def __init__(
        self,
        agent_id: str,
        agent_name: str,
        agent_instructions: str,
        organization_id: str,
        user_id: str,
        agent_description: Optional[str] = None,
        metadata: Optional[Dict[str, Any]] = None
    ):
        """
        Initialize RAG agent
        
        Args:
            agent_id: Agent ID
            agent_name: Agent name
            agent_instructions: Agent instructions
            organization_id: Organization ID
            user_id: User ID of the user accessing the agent
            agent_description: Optional agent description
            metadata: Optional additional metadata
        """
        self._agent_id = agent_id
        self._agent_name = agent_name
        self._agent_description = agent_description
        self._agent_instructions = agent_instructions
        self._organization_id = organization_id
        self._user_id = user_id
        self._metadata = metadata or {}
    
    @property
    def agent_id(self) -> str:
        """Get agent ID"""
        return self._agent_id
    
    @property
    def agent_name(self) -> str:
        """Get agent name"""
        return self._agent_name
    
    @property
    def agent_description(self) -> Optional[str]:
        """Get agent description"""
        return self._agent_description
    
    @property
    def agent_instructions(self) -> str:
        """Get agent instructions"""
        return self._agent_instructions
    
    @property
    def agent_metadata(self) -> Dict[str, Any]:
        """Get agent metadata"""
        return self._metadata
    
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
        # Start with standard access filters
        filters = get_access_filters(
            organization_id=self._organization_id,
            user_id=self._user_id
        )
        
        # Add agent filter - only retrieve documents associated with this agent
        agent_filter = {
            "associated_agents": self._agent_id
        }
        
        # If there's a document_filter in metadata, merge it
        if "document_filter" in self._metadata:
            doc_filter = self._metadata["document_filter"]
            for key, value in doc_filter.items():
                if key not in filters and key != "associated_agents":
                    agent_filter[key] = value
        
        # Merge with filters
        if "$and" not in filters:
            filters["$and"] = []
        filters["$and"].append(agent_filter)
        
        # Retrieve documents
        docs = await retrieve_relevant_documents(
            query=query,
            organization_id=self._organization_id,
            user_id=self._user_id,
            filter_dict=filters,
            limit=limit
        )
        
        return docs
    
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
        with LoggingContext(
            organization_id=self._organization_id,
            user_id=self._user_id,
            agent_id=self._agent_id
        ):
            try:
                # Get the last user message
                user_messages = [msg for msg in messages if msg["role"] == "user"]
                if not user_messages:
                    return "I need a question to assist you."
                
                last_user_message = user_messages[-1]["content"]
                
                # Retrieve relevant documents
                docs = await self.retrieve_relevant_docs(last_user_message)
                
                # Get LLM provider
                llm_provider = get_llm_provider(self._organization_id)
                
                # Create system message with agent instructions and context
                system_message = self._create_system_message(docs)
                
                # Create chat messages including system message
                chat_messages = [
                    {"role": "system", "content": system_message}
                ]
                
                # Add user conversation history
                chat_messages.extend(messages)
                
                # Generate response
                response = await llm_provider.chat(
                    messages=chat_messages,
                    temperature=temperature,
                    max_tokens=max_tokens,
                    **kwargs
                )
                
                logger.info(
                    f"Agent chat response generated",
                    extra={
                        "agent_id": self._agent_id,
                        "organization_id": self._organization_id,
                        "user_id": self._user_id,
                        "document_count": len(docs),
                        "temperature": temperature
                    }
                )
                
                return response
            
            except Exception as e:
                logger.error(
                    f"Error in agent chat: {str(e)}",
                    extra={
                        "agent_id": self._agent_id,
                        "organization_id": self._organization_id,
                        "user_id": self._user_id
                    }
                )
                raise
    
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
        with LoggingContext(
            organization_id=self._organization_id,
            user_id=self._user_id,
            agent_id=self._agent_id
        ):
            try:
                # Get the last user message
                user_messages = [msg for msg in messages if msg["role"] == "user"]
                if not user_messages:
                    yield "I need a question to assist you."
                    return
                
                last_user_message = user_messages[-1]["content"]
                
                # Retrieve relevant documents
                docs = await self.retrieve_relevant_docs(last_user_message)
                
                # Get LLM provider
                llm_provider = get_llm_provider(self._organization_id)
                
                # Create system message with agent instructions and context
                system_message = self._create_system_message(docs)
                
                # Create chat messages including system message
                chat_messages = [
                    {"role": "system", "content": system_message}
                ]
                
                # Add user conversation history
                chat_messages.extend(messages)
                
                # Generate streaming response
                stream = llm_provider.chat_stream(
                    messages=chat_messages,
                    temperature=temperature,
                    max_tokens=max_tokens,
                    **kwargs
                )
                
                logger.info(
                    f"Agent chat stream started",
                    extra={
                        "agent_id": self._agent_id,
                        "organization_id": self._organization_id,
                        "user_id": self._user_id,
                        "document_count": len(docs),
                        "temperature": temperature
                    }
                )
                
                async for chunk in stream:
                    yield chunk
            
            except Exception as e:
                logger.error(
                    f"Error in agent chat stream: {str(e)}",
                    extra={
                        "agent_id": self._agent_id,
                        "organization_id": self._organization_id,
                        "user_id": self._user_id
                    }
                )
                raise
    
    def _create_system_message(self, docs: List[Dict[str, Any]]) -> str:
        """
        Create system message with agent instructions and context
        
        Args:
            docs: Retrieved documents
            
        Returns:
            System message
        """
        # Format documents for the prompt
        context = self._format_documents_for_context(docs)
        
        # Create system message
        if context:
            system_message = f"""You are {self._agent_name}"""
            
            if self._agent_description:
                system_message += f", {self._agent_description}"
            
            system_message += f""". Your instructions are:

{self._agent_instructions}

You have access to the following knowledge:

{context}

Use this knowledge to help answer the user's questions. If the context doesn't contain the necessary information, use your general knowledge but prioritize the context information when available.

Always answer questions in a clear and helpful way, and maintain the personality described in your instructions.
"""
        else:
            # No documents available
            system_message = f"""You are {self._agent_name}"""
            
            if self._agent_description:
                system_message += f", {self._agent_description}"
            
            system_message += f""". Your instructions are:

{self._agent_instructions}

You don't have any specific knowledge available for this question, so please use your general knowledge to help the user.

Always answer questions in a clear and helpful way, and maintain the personality described in your instructions.
"""
        
        return system_message
    
    def _format_documents_for_context(self, documents: List[Dict[str, Any]]) -> str:
        """
        Format retrieved documents into a single context string
        
        Args:
            documents: List of document dictionaries with text and metadata
            
        Returns:
            Formatted context string
        """
        if not documents:
            return ""
        
        context_parts = []
        
        for i, doc in enumerate(documents):
            # Extract document information
            text = doc["text"]
            metadata = doc["metadata"]
            
            # Add document metadata if available
            source_info = []
            if "title" in metadata:
                source_info.append(f"Title: {metadata['title']}")
            if "source" in metadata:
                source_info.append(f"Source: {metadata['source']}")
            if "url" in metadata:
                source_info.append(f"URL: {metadata['url']}")
            if "date" in metadata:
                source_info.append(f"Date: {metadata['date']}")
            if "author" in metadata:
                source_info.append(f"Author: {metadata['author']}")
            
            # Format document with metadata
            if source_info:
                source_str = " | ".join(source_info)
                context_parts.append(f"Document {i+1} ({source_str}):\n{text}\n")
            else:
                context_parts.append(f"Document {i+1}:\n{text}\n")
        
        return "\n".join(context_parts)
    
    def to_dict(self) -> Dict[str, Any]:
        """
        Convert agent to dictionary representation
        
        Returns:
            Dictionary representation of the agent
        """
        return {
            "id": self._agent_id,
            "name": self._agent_name,
            "description": self._agent_description,
            "instructions": self._agent_instructions,
            "organization_id": self._organization_id,
            "metadata": self._metadata
        }