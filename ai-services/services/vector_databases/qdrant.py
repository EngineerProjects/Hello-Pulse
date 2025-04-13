"""
Qdrant provider implementation
"""

from typing import Dict, List, Optional, Any, Union, Tuple
import uuid
import json
from qdrant_client import QdrantClient
from qdrant_client.http import models as rest
from qdrant_client.http.exceptions import UnexpectedResponse

from core.config import settings as app_settings
from core.logging import logger
from .base import BaseVectorDatabase

class QdrantProvider(BaseVectorDatabase):
    """Qdrant provider implementation"""
    
    def __init__(self, organization_id: str, company_config: Dict[str, Any], shared: bool = True):
        """
        Initialize Qdrant provider
        
        Args:
            organization_id: Organization ID
            company_config: Company-specific configuration
            shared: Whether this is a shared database
        """
        self.organization_id = organization_id
        self.shared = shared
        
        # Get connection settings
        self.host = company_config.get("qdrant_host", app_settings.QDRANT_HOST)
        self.port = company_config.get("qdrant_port", app_settings.QDRANT_PORT)
        
        # Set up collection name
        if shared:
            # For shared DB, we use a single collection and filter by organization
            self.collection_name = app_settings.QDRANT_COLLECTION_NAME
        else:
            # For private DB, we create a separate collection per organization
            self.collection_name = f"{app_settings.QDRANT_COLLECTION_NAME}_{organization_id}"
        
        # Initialize client
        self.client = QdrantClient(
            host=self.host,
            port=self.port
        )
        
        # Try to get/create collection
        try:
            # Check if collection exists
            collections = self.client.get_collections().collections
            collection_names = [c.name for c in collections]
            
            if self.collection_name not in collection_names:
                # Create collection with necessary fields for filtering
                self.client.create_collection(
                    collection_name=self.collection_name,
                    vectors_config=rest.VectorParams(
                        size=1536,  # Default size, will be adjusted during first add
                        distance=rest.Distance.COSINE
                    ),
                    # Define payload fields for filtering
                    metadata={
                        "organization_id": self.organization_id,
                        "shared": shared
                    }
                )
                
                # Create index for fast organization_id filtering
                self.client.create_payload_index(
                    collection_name=self.collection_name,
                    field_name="organization_id",
                    field_schema=rest.PayloadSchemaType.KEYWORD
                )
                
                # Create index for visibility filtering
                self.client.create_payload_index(
                    collection_name=self.collection_name,
                    field_name="visibility",
                    field_schema=rest.PayloadSchemaType.KEYWORD
                )
                
                # Create index for user_id filtering
                self.client.create_payload_index(
                    collection_name=self.collection_name,
                    field_name="user_id",
                    field_schema=rest.PayloadSchemaType.KEYWORD
                )
            
            logger.info(
                f"Connected to Qdrant collection",
                extra={
                    "collection": self.collection_name,
                    "organization_id": organization_id,
                    "shared": shared
                }
            )
        except Exception as e:
            logger.error(
                f"Qdrant connection error: {str(e)}",
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
        """Add documents to Qdrant"""
        try:
            # Generate IDs if not provided
            ids = kwargs.get("ids", [str(uuid.uuid4()) for _ in range(len(texts))])
            
            # Prepare points
            points = []
            for i in range(len(texts)):
                # Ensure all metadata includes organization_id for shared collections
                payload = dict(metadata[i])
                if self.shared:
                    payload["organization_id"] = self.organization_id
                
                # Add text to payload
                payload["text"] = texts[i]
                
                # Create point
                points.append(
                    rest.PointStruct(
                        id=ids[i],
                        vector=embeddings[i],
                        payload=payload
                    )
                )
            
            # Add points
            self.client.upsert(
                collection_name=self.collection_name,
                points=points
            )
            
            return ids
        except Exception as e:
            logger.error(
                f"Qdrant add_documents error: {str(e)}",
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
        """Search for similar documents in Qdrant"""
        try:
            # Ensure organization_id filter for shared collections
            if self.shared and "organization_id" not in filter_dict:
                filter_dict["organization_id"] = self.organization_id
            
            # Convert filter dict to Qdrant format
            filter_conditions = []
            for key, value in filter_dict.items():
                if isinstance(value, list):
                    # Handle list values (IN operator)
                    filter_conditions.append(
                        rest.FieldCondition(
                            key=key,
                            match=rest.MatchAny(any=value)
                        )
                    )
                else:
                    # Handle single values (match)
                    filter_conditions.append(
                        rest.FieldCondition(
                            key=key,
                            match=rest.MatchValue(value=value)
                        )
                    )
            
            # Combine conditions with AND
            if filter_conditions:
                filter_query = rest.Filter(
                    must=filter_conditions
                )
            else:
                filter_query = None
            
            # Perform search
            search_params = rest.SearchParams(hnsw_ef=128)
            results = self.client.search(
                collection_name=self.collection_name,
                query_vector=query_embedding,
                query_filter=filter_query,
                limit=limit,
                search_params=search_params,
                with_payload=True,
                with_vectors=kwargs.get("with_vectors", False)
            )
            
            # Format results
            formatted_results = []
            for result in results:
                # Extract text from payload
                text = result.payload.pop("text", "")
                
                formatted_results.append({
                    "id": str(result.id),
                    "text": text,
                    "metadata": result.payload,
                    "score": result.score
                })
                
                # Add vector if requested
                if kwargs.get("with_vectors", False) and hasattr(result, "vector"):
                    formatted_results[-1]["embedding"] = result.vector
            
            return formatted_results
        except Exception as e:
            logger.error(
                f"Qdrant search error: {str(e)}",
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
        """Delete documents from Qdrant"""
        try:
            if self.shared:
                # For shared collections, ensure we only delete documents from this organization
                # First, get the documents
                points = self.client.retrieve(
                    collection_name=self.collection_name,
                    ids=document_ids,
                    with_payload=True
                )
                
                # Filter out IDs that don't belong to this organization
                filtered_ids = []
                for point in points:
                    if point.payload.get("organization_id") == self.organization_id:
                        filtered_ids.append(point.id)
                
                if not filtered_ids:
                    return True  # No documents to delete
                
                # Delete the filtered documents
                self.client.delete(
                    collection_name=self.collection_name,
                    points_selector=rest.PointIdsList(
                        points=filtered_ids
                    )
                )
            else:
                # For private collections, we can delete directly
                self.client.delete(
                    collection_name=self.collection_name,
                    points_selector=rest.PointIdsList(
                        points=document_ids
                    )
                )
            
            return True
        except Exception as e:
            logger.error(
                f"Qdrant delete_documents error: {str(e)}",
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
        """Update a document in Qdrant"""
        try:
            # Get current document
            points = self.client.retrieve(
                collection_name=self.collection_name,
                ids=[document_id],
                with_payload=True,
                with_vectors=embedding is None  # Only get vectors if we're not updating them
            )
            
            if not points:
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
                current_payload = points[0].payload
                if current_payload.get("organization_id") != self.organization_id:
                    logger.warning(
                        f"Attempted to update document from another organization",
                        extra={
                            "collection": self.collection_name,
                            "document_id": document_id,
                            "organization_id": self.organization_id,
                            "owner_org_id": current_payload.get("organization_id")
                        }
                    )
                    return False
            
            # Prepare update payload
            payload = {}
            if text is not None:
                payload["text"] = text
            
            if metadata is not None:
                # Ensure organization_id is preserved for shared collections
                if self.shared and "organization_id" not in metadata:
                    metadata["organization_id"] = self.organization_id
                payload.update(metadata)
            
            # Update the document
            if payload:
                self.client.set_payload(
                    collection_name=self.collection_name,
                    payload=payload,
                    points=[document_id]
                )
            
            # Update the embedding if provided
            if embedding is not None:
                self.client.update_vectors(
                    collection_name=self.collection_name,
                    points=[
                        rest.PointVectors(
                            id=document_id,
                            vector=embedding
                        )
                    ]
                )
            
            return True
        except Exception as e:
            logger.error(
                f"Qdrant update_document error: {str(e)}",
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
        """Get a document by ID from Qdrant"""
        try:
            with_vectors = kwargs.get("with_vectors", False)
            
            points = self.client.retrieve(
                collection_name=self.collection_name,
                ids=[document_id],
                with_payload=True,
                with_vectors=with_vectors
            )
            
            if not points:
                return None
            
            point = points[0]
            
            # Check organization ownership for shared collections
            if self.shared:
                if point.payload.get("organization_id") != self.organization_id:
                    logger.warning(
                        f"Attempted to access document from another organization",
                        extra={
                            "collection": self.collection_name,
                            "document_id": document_id,
                            "organization_id": self.organization_id,
                            "owner_org_id": point.payload.get("organization_id")
                        }
                    )
                    return None
            
            # Extract text from payload
            text = point.payload.pop("text", "")
            
            result = {
                "id": str(point.id),
                "text": text,
                "metadata": point.payload
            }
            
            if with_vectors and hasattr(point, "vector"):
                result["embedding"] = point.vector
            
            return result
        except Exception as e:
            logger.error(
                f"Qdrant get_document error: {str(e)}",
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
        """List documents matching filters from Qdrant"""
        try:
            # Ensure organization_id filter for shared collections
            if self.shared and "organization_id" not in filter_dict:
                filter_dict["organization_id"] = self.organization_id
            
            # Convert filter dict to Qdrant format
            filter_conditions = []
            for key, value in filter_dict.items():
                if isinstance(value, list):
                    # Handle list values (IN operator)
                    filter_conditions.append(
                        rest.FieldCondition(
                            key=key,
                            match=rest.MatchAny(any=value)
                        )
                    )
                else:
                    # Handle single values (match)
                    filter_conditions.append(
                        rest.FieldCondition(
                            key=key,
                            match=rest.MatchValue(value=value)
                        )
                    )
            
            # Combine conditions with AND
            if filter_conditions:
                filter_query = rest.Filter(
                    must=filter_conditions
                )
            else:
                filter_query = None
            
            # Use scroll API for pagination
            points, next_page_offset = self.client.scroll(
                collection_name=self.collection_name,
                scroll_filter=filter_query,
                limit=limit,
                offset=offset,
                with_payload=True,
                with_vectors=kwargs.get("with_vectors", False)
            )
            
            # Format results
            formatted_results = []
            for point in points:
                # Extract text from payload
                text = point.payload.pop("text", "")
                
                result = {
                    "id": str(point.id),
                    "text": text,
                    "metadata": point.payload
                }
                
                if kwargs.get("with_vectors", False) and hasattr(point, "vector"):
                    result["embedding"] = point.vector
                
                formatted_results.append(result)
            
            return formatted_results
        except Exception as e:
            logger.error(
                f"Qdrant list_documents error: {str(e)}",
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
        """Count documents matching filters in Qdrant"""
        try:
            # Ensure organization_id filter for shared collections
            if self.shared and "organization_id" not in filter_dict:
                filter_dict["organization_id"] = self.organization_id
            
            # Convert filter dict to Qdrant format
            filter_conditions = []
            for key, value in filter_dict.items():
                if isinstance(value, list):
                    # Handle list values (IN operator)
                    filter_conditions.append(
                        rest.FieldCondition(
                            key=key,
                            match=rest.MatchAny(any=value)
                        )
                    )
                else:
                    # Handle single values (match)
                    filter_conditions.append(
                        rest.FieldCondition(
                            key=key,
                            match=rest.MatchValue(value=value)
                        )
                    )
            
            # Combine conditions with AND
            if filter_conditions:
                filter_query = rest.Filter(
                    must=filter_conditions
                )
            else:
                filter_query = None
            
            # Get count
            count_result = self.client.count(
                collection_name=self.collection_name,
                count_filter=filter_query
            )
            
            return count_result.count
        except Exception as e:
            logger.error(
                f"Qdrant count_documents error: {str(e)}",
                extra={
                    "collection": self.collection_name,
                    "organization_id": self.organization_id,
                    "filter": filter_dict
                }
            )
            return 0
    
    def get_provider_name(self) -> str:
        """Get provider name"""
        return "qdrant"
    
    def get_collection_name(self) -> str:
        """Get collection name"""
        return self.collection_name
    
    async def create_collection(
        self,
        collection_name: str,
        embedding_dimension: int,
        **kwargs: Any
    ) -> bool:
        """Create a new collection in Qdrant"""
        try:
            metadata = kwargs.get("metadata", {})
            if self.shared:
                metadata["organization_id"] = self.organization_id
            
            # Create collection
            self.client.create_collection(
                collection_name=collection_name,
                vectors_config=rest.VectorParams(
                    size=embedding_dimension,
                    distance=rest.Distance.COSINE
                ),
                metadata=metadata
            )
            
            # Create index for fast organization_id filtering
            self.client.create_payload_index(
                collection_name=collection_name,
                field_name="organization_id",
                field_schema=rest.PayloadSchemaType.KEYWORD
            )
            
            # Create index for visibility filtering
            self.client.create_payload_index(
                collection_name=collection_name,
                field_name="visibility",
                field_schema=rest.PayloadSchemaType.KEYWORD
            )
            
            # Create index for user_id filtering
            self.client.create_payload_index(
                collection_name=collection_name,
                field_name="user_id",
                field_schema=rest.PayloadSchemaType.KEYWORD
            )
            
            return True
        except Exception as e:
            logger.error(
                f"Qdrant create_collection error: {str(e)}",
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
        """Delete a collection from Qdrant"""
        try:
            # For safety, verify ownership in shared mode
            if self.shared and collection_name != self.collection_name:
                # Get collection info
                collection_info = self.client.get_collection(collection_name)
                
                # Check if we own this collection
                metadata = collection_info.metadata
                if metadata and metadata.get("organization_id") != self.organization_id:
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
                f"Qdrant delete_collection error: {str(e)}",
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
        """List available collections in Qdrant"""
        try:
            collections = self.client.get_collections().collections
            
            # Filter by organization for shared mode
            if self.shared:
                result = []
                for collection in collections:
                    collection_info = self.client.get_collection(collection.name)
                    metadata = collection_info.metadata
                    if metadata and metadata.get("organization_id") == self.organization_id:
                        result.append(collection.name)
                return result
            else:
                return [collection.name for collection in collections]
        except Exception as e:
            logger.error(
                f"Qdrant list_collections error: {str(e)}",
                extra={
                    "organization_id": self.organization_id
                }
            )
            return []