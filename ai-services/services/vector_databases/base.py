"""
Base interface for vector databases
"""

from abc import ABC, abstractmethod
from typing import Dict, List, Optional, Any, Tuple, Union

class BaseVectorDatabase(ABC):
    """
    Base interface for vector databases
    All vector database implementations must inherit from this class
    """
    
    @abstractmethod
    async def add_documents(
        self,
        texts: List[str],
        embeddings: List[List[float]],
        metadata: List[Dict[str, Any]],
        **kwargs: Any
    ) -> List[str]:
        """
        Add documents to vector database
        
        Args:
            texts: List of document texts
            embeddings: List of embedding vectors
            metadata: List of metadata dictionaries
            **kwargs: Additional provider-specific parameters
            
        Returns:
            List of document IDs
        """
        pass
    
    @abstractmethod
    async def search(
        self,
        query_embedding: List[float],
        filter_dict: Dict[str, Any],
        limit: int = 10,
        **kwargs: Any
    ) -> List[Dict[str, Any]]:
        """
        Search for similar documents
        
        Args:
            query_embedding: Query embedding vector
            filter_dict: Filtering criteria (must include security filters)
            limit: Maximum number of results
            **kwargs: Additional provider-specific parameters
            
        Returns:
            List of document dictionaries with text, metadata, and score
        """
        pass
    
    @abstractmethod
    async def delete_documents(
        self,
        document_ids: List[str],
        **kwargs: Any
    ) -> bool:
        """
        Delete documents from the vector database
        
        Args:
            document_ids: List of document IDs to delete
            **kwargs: Additional provider-specific parameters
            
        Returns:
            True if successful, False otherwise
        """
        pass
    
    @abstractmethod
    async def update_document(
        self,
        document_id: str,
        text: Optional[str] = None,
        embedding: Optional[List[float]] = None,
        metadata: Optional[Dict[str, Any]] = None,
        **kwargs: Any
    ) -> bool:
        """
        Update a document in the vector database
        
        Args:
            document_id: Document ID to update
            text: New text (if None, keep existing)
            embedding: New embedding (if None, keep existing)
            metadata: New metadata (if None, keep existing)
            **kwargs: Additional provider-specific parameters
            
        Returns:
            True if successful, False otherwise
        """
        pass
    
    @abstractmethod
    async def get_document(
        self,
        document_id: str,
        **kwargs: Any
    ) -> Optional[Dict[str, Any]]:
        """
        Get a document by ID
        
        Args:
            document_id: Document ID
            **kwargs: Additional provider-specific parameters
            
        Returns:
            Document dictionary or None if not found
        """
        pass
    
    @abstractmethod
    async def list_documents(
        self,
        filter_dict: Dict[str, Any],
        limit: int = 100,
        offset: int = 0,
        **kwargs: Any
    ) -> List[Dict[str, Any]]:
        """
        List documents matching filters
        
        Args:
            filter_dict: Filtering criteria (must include security filters)
            limit: Maximum number of results
            offset: Pagination offset
            **kwargs: Additional provider-specific parameters
            
        Returns:
            List of document dictionaries
        """
        pass
    
    @abstractmethod
    async def count_documents(
        self,
        filter_dict: Dict[str, Any],
        **kwargs: Any
    ) -> int:
        """
        Count documents matching filters
        
        Args:
            filter_dict: Filtering criteria (must include security filters)
            **kwargs: Additional provider-specific parameters
            
        Returns:
            Document count
        """
        pass
    
    @abstractmethod
    def get_provider_name(self) -> str:
        """
        Get the name of this vector database provider
        
        Returns:
            Provider name as string
        """
        pass
    
    @abstractmethod
    def get_collection_name(self) -> str:
        """
        Get the name of the collection being used
        
        Returns:
            Collection name as string
        """
        pass
    
    @abstractmethod
    async def create_collection(
        self,
        collection_name: str,
        embedding_dimension: int,
        **kwargs: Any
    ) -> bool:
        """
        Create a new collection
        
        Args:
            collection_name: Collection name
            embedding_dimension: Dimension of embeddings
            **kwargs: Additional provider-specific parameters
            
        Returns:
            True if successful, False otherwise
        """
        pass
    
    @abstractmethod
    async def delete_collection(
        self,
        collection_name: str,
        **kwargs: Any
    ) -> bool:
        """
        Delete a collection
        
        Args:
            collection_name: Collection name
            **kwargs: Additional provider-specific parameters
            
        Returns:
            True if successful, False otherwise
        """
        pass
    
    @abstractmethod
    async def list_collections(
        self,
        **kwargs: Any
    ) -> List[str]:
        """
        List available collections
        
        Args:
            **kwargs: Additional provider-specific parameters
            
        Returns:
            List of collection names
        """
        pass