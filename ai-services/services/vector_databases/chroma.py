"""
ChromaDB provider implementation
"""

from typing import Dict, List, Optional, Any, Union, Tuple
import uuid
import chromadb
from chromadb.api.models.Collection import Collection
from chromadb.config import Settings

from core.config import settings as app_settings
from core.logging import logger
from .base import BaseVectorDatabase

class ChromaDBProvider(BaseVectorDatabase):
    """ChromaDB provider implementation"""
    
    def __init__(self, organization_id: str, company_config: Dict[str, Any], shared: bool = True):
        """
        Initialize ChromaDB provider
        
        Args:
            organization_id: Organization ID
            company_config: Company-specific configuration
            shared: Whether this is a shared database
        """
        self.organization_id = organization_id
        self.shared = shared
        
        # Get connection settings
        self.host = company_config.get("chroma_host", app_settings.CHROMA_HOST)
        self.port = company_config.get("chroma_port", app_settings.CHROMA_PORT)
        
        # Set up collection name
        if shared:
            # For shared DB, we use a single collection and filter by organization
            self.collection_name = app_settings.CHROMA_COLLECTION_NAME
        else:
            # For private DB, we create a separate collection per organization
            self.collection_name = f"{app_settings.CHROMA_COLLECTION_NAME}_{organization_id}"
        
        # Initialize client
        self.client = chromadb.HttpClient(
            host=self.host,
            port=self.port
        )
        
        # Try to get/create collection
        try:
            self.collection = self.client.get_or_create_collection(
                name=self.collection_name,
                metadata={"organization_id": organization_id}
            )
            logger.info(
                f"Connected to ChromaDB collection",
                extra={
                    "collection": self.collection_name,
                    "organization_id": organization_id,
                    "shared": shared
                }
            )
        except Exception as e:
            logger.error(
                f"ChromaDB connection error: {str(e)}",
                extra={
                    "collection": self.collection_name,
                    "organization_id": organization_id
                }
            )
            raise
    
    async def add_documents(
        self,
        texts: List[str],
        embeddings: List[List[float]],
        metadata: List[Dict[str, Any]],
        **kwargs: Any
    ) -> List[str]:
        """Add documents to ChromaDB"""
        try:
            # Generate IDs if not provided
            ids = kwargs.get("ids", [str(uuid.uuid4()) for _ in range(len(texts))])
            
            # Ensure all metadata includes organization_id for shared collections
            if self.shared:
                for meta in metadata:
                    meta["organization_id"] = self.organization_id
            
            # Add documents
            self.collection.add(
                documents=texts,
                embeddings=embeddings,
                metadatas=metadata,
                ids=ids
            )
            
            return ids
        except Exception as e:
            logger.error(
                f"ChromaDB add_documents error: {str(e)}",
                extra={
                    "collection": self.collection_name,
                    "organization_id": self.organization_id,
                    "document_count": len(texts)
                }
            )
            raise
    
    async def search(
        self,
        query_embedding: List[float],
        filter_dict: Dict[str, Any],
        limit: int = 10,
        **kwargs: Any
    ) -> List[Dict[str, Any]]:
        """Search for similar documents in ChromaDB"""
        try:
            # Ensure organization_id filter for shared collections
            if self.shared and "organization_id" not in filter_dict:
                filter_dict["organization_id"] = self.organization_id
            
            # Convert filter dict to ChromaDB format
            where = {}
            for key, value in filter_dict.items():
                where[key] = value
            
            # Perform search
            results = self.collection.query(
                query_embeddings=[query_embedding],
                where=where,
                n_results=limit,
                include=["documents", "metadatas", "distances"]
            )
            
            # Format results
            formatted_results = []
            for i in range(len(results["ids"][0])):
                formatted_results.append({
                    "id": results["ids"][0][i],
                    "text": results["documents"][0][i],
                    "metadata": results["metadatas"][0][i],
                    "score": 1.0 - float(results["distances"][0][i])  # Convert distance to similarity score
                })
            
            return formatted_results
        except Exception as e:
            logger.error(
                f"ChromaDB search error: {str(e)}",
                extra={
                    "collection": self.collection_name,
                    "organization_id": self.organization_id,
                    "filter": filter_dict
                }
            )
            raise
    
    async def delete_documents(
        self,
        document_ids: List[str],
        **kwargs: Any
    ) -> bool:
        """Delete documents from ChromaDB"""
        try:
            if self.shared:
                # For shared collections, ensure we only delete documents from this organization
                # First, get the documents
                results = self.collection.get(
                    ids=document_ids,
                    include=["metadatas"]
                )
                
                # Filter out IDs that don't belong to this organization
                filtered_ids = []
                for i, metadata in enumerate(results["metadatas"]):
                    if metadata.get("organization_id") == self.organization_id:
                        filtered_ids.append(document_ids[i])
                
                if not filtered_ids:
                    return True  # No documents to delete
                
                # Delete the filtered documents
                self.collection.delete(ids=filtered_ids)
            else:
                # For private collections, we can delete directly
                self.collection.delete(ids=document_ids)
            
            return True
        except Exception as e:
            logger.error(
                f"ChromaDB delete_documents error: {str(e)}",
                extra={
                    "collection": self.collection_name,
                    "organization_id": self.organization_id,
                    "document_count": len(document_ids)
                }
            )
            return False
    
    async def update_document(
        self,
        document_id: str,
        text: Optional[str] = None,
        embedding: Optional[List[float]] = None,
        metadata: Optional[Dict[str, Any]] = None,
        **kwargs: Any
    ) -> bool:
        """Update a document in ChromaDB"""
        try:
            # Get current document
            results = self.collection.get(
                ids=[document_id],
                include=["documents", "metadatas", "embeddings"]
            )
            
            if not results["ids"]:
                logger.warning(
                    f"Document not found for update",
                    extra={
                        "collection": self.collection_name,
                        "document_id": document_id,
                        "organization_id": self.organization_id
                    }
                )
                return False
            
            # Check organization ownership for shared collections
            if self.shared:
                current_metadata = results["metadatas"][0]
                if current_metadata.get("organization_id") != self.organization_id:
                    logger.warning(
                        f"Attempted to update document from another organization",
                        extra={
                            "collection": self.collection_name,
                            "document_id": document_id,
                            "organization_id": self.organization_id,
                            "owner_org_id": current_metadata.get("organization_id")
                        }
                    )
                    return False
            
            # Prepare update data
            update_data = {}
            if text is not None:
                update_data["documents"] = text
            
            if embedding is not None:
                update_data["embeddings"] = embedding
            
            if metadata is not None:
                # Ensure organization_id is preserved for shared collections
                if self.shared and "organization_id" not in metadata:
                    metadata["organization_id"] = self.organization_id
                update_data["metadatas"] = metadata
            
            # Update the document
            if update_data:
                self.collection.update(
                    ids=document_id,
                    **update_data
                )
            
            return True
        except Exception as e:
            logger.error(
                f"ChromaDB update_document error: {str(e)}",
                extra={
                    "collection": self.collection_name,
                    "document_id": document_id,
                    "organization_id": self.organization_id
                }
            )
            return False
    
    async def get_document(
        self,
        document_id: str,
        **kwargs: Any
    ) -> Optional[Dict[str, Any]]:
        """Get a document by ID from ChromaDB"""
        try:
            results = self.collection.get(
                ids=[document_id],
                include=["documents", "metadatas", "embeddings"]
            )
            
            if not results["ids"]:
                return None
            
            # Check organization ownership for shared collections
            if self.shared:
                metadata = results["metadatas"][0]
                if metadata.get("organization_id") != self.organization_id:
                    logger.warning(
                        f"Attempted to access document from another organization",
                        extra={
                            "collection": self.collection_name,
                            "document_id": document_id,
                            "organization_id": self.organization_id,
                            "owner_org_id": metadata.get("organization_id")
                        }
                    )
                    return None
            
            return {
                "id": results["ids"][0],
                "text": results["documents"][0],
                "metadata": results["metadatas"][0],
                "embedding": results["embeddings"][0] if "embeddings" in results else None
            }
        except Exception as e:
            logger.error(
                f"ChromaDB get_document error: {str(e)}",
                extra={
                    "collection": self.collection_name,
                    "document_id": document_id,
                    "organization_id": self.organization_id
                }
            )
            return None
    
    async def list_documents(
        self,
        filter_dict: Dict[str, Any],
        limit: int = 100,
        offset: int = 0,
        **kwargs: Any
    ) -> List[Dict[str, Any]]:
        """List documents matching filters from ChromaDB"""
        try:
            # Ensure organization_id filter for shared collections
            if self.shared and "organization_id" not in filter_dict:
                filter_dict["organization_id"] = self.organization_id
            
            # Convert filter dict to ChromaDB format
            where = {}
            for key, value in filter_dict.items():
                where[key] = value
            
            # ChromaDB doesn't support offset/limit directly
            # We'll get more results and then slice
            safe_limit = limit + offset
            
            # Perform query
            results = self.collection.get(
                where=where,
                limit=safe_limit,
                include=["documents", "metadatas"]
            )
            
            # Apply offset/limit manually
            formatted_results = []
            for i in range(min(len(results["ids"]), safe_limit)):
                if i >= offset:
                    formatted_results.append({
                        "id": results["ids"][i],
                        "text": results["documents"][i],
                        "metadata": results["metadatas"][i]
                    })
            
            return formatted_results
        except Exception as e:
            logger.error(
                f"ChromaDB list_documents error: {str(e)}",
                extra={
                    "collection": self.collection_name,
                    "organization_id": self.organization_id,
                    "filter": filter_dict
                }
            )
            return []
    
    async def count_documents(
        self,
        filter_dict: Dict[str, Any],
        **kwargs: Any
    ) -> int:
        """Count documents matching filters in ChromaDB"""
        try:
            # Ensure organization_id filter for shared collections
            if self.shared and "organization_id" not in filter_dict:
                filter_dict["organization_id"] = self.organization_id
            
            # Convert filter dict to ChromaDB format
            where = {}
            for key, value in filter_dict.items():
                where[key] = value
            
            # ChromaDB doesn't have a count method
            # We'll get just the IDs to count them
            results = self.collection.get(
                where=where,
                include=[]  # Only get IDs
            )
            
            return len(results["ids"])
        except Exception as e:
            logger.error(
                f"ChromaDB count_documents error: {str(e)}",
                extra={
                    "collection": self.collection_name,
                    "organization_id": self.organization_id,
                    "filter": filter_dict
                }
            )
            return 0
    
    def get_provider_name(self) -> str:
        """Get provider name"""
        return "chroma"
    
    def get_collection_name(self) -> str:
        """Get collection name"""
        return self.collection_name
    
    async def create_collection(
        self,
        collection_name: str,
        embedding_dimension: int,
        **kwargs: Any
    ) -> bool:
        """Create a new collection in ChromaDB"""
        try:
            metadata = kwargs.get("metadata", {})
            if self.shared:
                metadata["organization_id"] = self.organization_id
            
            self.client.create_collection(
                name=collection_name,
                metadata=metadata
            )
            return True
        except Exception as e:
            logger.error(
                f"ChromaDB create_collection error: {str(e)}",
                extra={
                    "collection": collection_name,
                    "organization_id": self.organization_id
                }
            )
            return False
    
    async def delete_collection(
        self,
        collection_name: str,
        **kwargs: Any
    ) -> bool:
        """Delete a collection from ChromaDB"""
        try:
            # For safety, verify ownership in shared mode
            if self.shared and collection_name != self.collection_name:
                # Get collection metadata
                collection = self.client.get_collection(collection_name)
                metadata = collection.metadata
                
                # Check if we own this collection
                if metadata.get("organization_id") != self.organization_id:
                    logger.warning(
                        f"Attempted to delete collection from another organization",
                        extra={
                            "collection": collection_name,
                            "organization_id": self.organization_id,
                            "owner_org_id": metadata.get("organization_id")
                        }
                    )
                    return False
            
            self.client.delete_collection(collection_name)
            return True
        except Exception as e:
            logger.error(
                f"ChromaDB delete_collection error: {str(e)}",
                extra={
                    "collection": collection_name,
                    "organization_id": self.organization_id
                }
            )
            return False
    
    async def list_collections(
        self,
        **kwargs: Any
    ) -> List[str]:
        """List available collections in ChromaDB"""
        try:
            collections = self.client.list_collections()
            
            # Filter by organization for shared mode
            if self.shared:
                result = []
                for collection in collections:
                    metadata = collection.metadata
                    if metadata and metadata.get("organization_id") == self.organization_id:
                        result.append(collection.name)
                return result
            else:
                return [collection.name for collection in collections]
        except Exception as e:
            logger.error(
                f"ChromaDB list_collections error: {str(e)}",
                extra={
                    "organization_id": self.organization_id
                }
            )
            return []