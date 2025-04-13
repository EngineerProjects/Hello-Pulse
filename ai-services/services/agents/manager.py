"""
Agent management service
"""

import uuid
import json
from datetime import datetime
from typing import Dict, List, Optional, Any, AsyncIterator, Union, Tuple

from core.logging import logger, LoggingContext
from services.llm_providers import get_llm_provider
from services.vector_databases import get_vector_db
from services.rag.retriever import retrieve_relevant_documents
from services.agents.personalization import RAGAgent
from utils.filtering import get_access_filters, sanitize_metadata_for_storage


class AgentManager:
    """Agent management service"""
    
    def __init__(self):
        """Initialize agent manager"""
        # Cache for active agents
        self._agent_cache: Dict[str, RAGAgent] = {}
    
    async def create_agent(
        self,
        organization_id: str,
        user_id: str,
        name: str,
        instructions: str,
        description: Optional[str] = None,
        document_ids: Optional[List[str]] = None,
        document_filter: Optional[Dict[str, Any]] = None,
        metadata: Optional[Dict[str, Any]] = None
    ) -> Dict[str, Any]:
        """
        Create a new agent
        
        Args:
            organization_id: Organization ID
            user_id: User ID
            name: Agent name
            instructions: Agent instructions
            description: Optional agent description
            document_ids: Optional list of document IDs to associate with the agent
            document_filter: Optional filter for finding relevant documents
            metadata: Optional additional metadata
            
        Returns:
            Agent information
        """
        with LoggingContext(organization_id=organization_id, user_id=user_id):
            try:
                # Get vector database
                vector_db = get_vector_db(organization_id)
                
                # Generate agent ID
                agent_id = str(uuid.uuid4())
                
                # Create agent metadata
                agent_metadata = {
                    "agent_id": agent_id,
                    "name": name,
                    "description": description,
                    "instructions": instructions,
                    "created_by": user_id,
                    "created_at": datetime.utcnow().isoformat(),
                    "organization_id": organization_id,
                    "type": "rag_agent",
                }
                
                # Add custom metadata if provided
                if metadata:
                    for key, value in metadata.items():
                        if key not in agent_metadata:
                            agent_metadata[key] = value
                
                # Store agent metadata in the vector database
                # We store it as a special document with no embedding
                placeholder_embedding = [0.0] * 1536  # Placeholder embedding
                
                agent_doc_ids = await vector_db.add_documents(
                    texts=[json.dumps(agent_metadata)],
                    embeddings=[placeholder_embedding],
                    metadata=[{
                        "agent_id": agent_id,
                        "organization_id": organization_id,
                        "user_id": user_id,
                        "visibility": "private",
                        "document_type": "agent_metadata",
                        "timestamp": int(datetime.utcnow().timestamp())
                    }]
                )
                
                agent_doc_id = agent_doc_ids[0]
                
                # Associate documents with agent if document_ids is provided
                if document_ids:
                    # Verify each document and associate it with the agent
                    for doc_id in document_ids:
                        # Get document to verify access
                        doc = await vector_db.get_document(doc_id)
                        
                        if not doc:
                            logger.warning(
                                f"Document not found: {doc_id}",
                                extra={
                                    "agent_id": agent_id,
                                    "organization_id": organization_id,
                                    "user_id": user_id
                                }
                            )
                            continue
                        
                        # Verify organization access
                        doc_org_id = doc["metadata"].get("organization_id")
                        if doc_org_id != organization_id:
                            logger.warning(
                                f"Attempted to associate document from another organization",
                                extra={
                                    "agent_id": agent_id,
                                    "document_id": doc_id,
                                    "organization_id": organization_id,
                                    "document_org_id": doc_org_id
                                }
                            )
                            continue
                        
                        # Add agent_id to document metadata
                        updated_metadata = dict(doc["metadata"])
                        
                        # Initialize or append to associated_agents list
                        if "associated_agents" not in updated_metadata:
                            updated_metadata["associated_agents"] = [agent_id]
                        elif agent_id not in updated_metadata["associated_agents"]:
                            updated_metadata["associated_agents"].append(agent_id)
                        
                        # Update document
                        await vector_db.update_document(
                            document_id=doc_id,
                            metadata=updated_metadata
                        )
                
                # Store document filter if provided
                if document_filter:
                    agent_metadata["document_filter"] = document_filter
                    
                    # Update the agent metadata document
                    await vector_db.update_document(
                        document_id=agent_doc_id,
                        text=json.dumps(agent_metadata)
                    )
                
                logger.info(
                    f"Created agent: {name}",
                    extra={
                        "agent_id": agent_id,
                        "organization_id": organization_id,
                        "user_id": user_id,
                        "document_count": len(document_ids) if document_ids else 0
                    }
                )
                
                # Format the response
                response = {
                    "id": agent_id,
                    "name": name,
                    "description": description,
                    "instructions": instructions,
                    "created_at": agent_metadata["created_at"],
                    "metadata": metadata or {},
                    "document_count": len(document_ids) if document_ids else 0
                }
                
                return response
            
            except Exception as e:
                logger.error(
                    f"Error creating agent: {str(e)}",
                    extra={
                        "name": name,
                        "organization_id": organization_id,
                        "user_id": user_id
                    }
                )
                raise
    
    async def get_agent(
        self,
        agent_id: str,
        organization_id: str,
        user_id: str
    ) -> Optional[Dict[str, Any]]:
        """
        Get agent information
        
        Args:
            agent_id: Agent ID
            organization_id: Organization ID
            user_id: User ID
            
        Returns:
            Agent information or None if not found
        """
        with LoggingContext(organization_id=organization_id, user_id=user_id):
            try:
                # Get vector database
                vector_db = get_vector_db(organization_id)
                
                # Find agent metadata document
                filter_dict = {
                    "agent_id": agent_id,
                    "organization_id": organization_id,
                    "document_type": "agent_metadata"
                }
                
                docs = await vector_db.list_documents(
                    filter_dict=filter_dict,
                    limit=1
                )
                
                if not docs:
                    logger.warning(
                        f"Agent not found: {agent_id}",
                        extra={
                            "agent_id": agent_id,
                            "organization_id": organization_id,
                            "user_id": user_id
                        }
                    )
                    return None
                
                # Extract agent metadata
                agent_doc = docs[0]
                agent_metadata = json.loads(agent_doc["text"])
                
                # Count associated documents
                associated_docs_filter = {
                    "associated_agents": agent_id,
                    "organization_id": organization_id
                }
                
                document_count = await vector_db.count_documents(associated_docs_filter)
                
                # Format the response
                response = {
                    "id": agent_id,
                    "name": agent_metadata.get("name", "Unnamed Agent"),
                    "description": agent_metadata.get("description"),
                    "instructions": agent_metadata.get("instructions", ""),
                    "created_at": agent_metadata.get("created_at"),
                    "updated_at": agent_metadata.get("updated_at"),
                    "metadata": {k: v for k, v in agent_metadata.items() if k not in [
                        "agent_id", "name", "description", "instructions", 
                        "created_at", "updated_at", "created_by", "organization_id", "type"
                    ]},
                    "document_count": document_count
                }
                
                return response
            
            except Exception as e:
                logger.error(
                    f"Error getting agent: {str(e)}",
                    extra={
                        "agent_id": agent_id,
                        "organization_id": organization_id,
                        "user_id": user_id
                    }
                )
                raise
    
    async def list_agents(
        self,
        organization_id: str,
        user_id: str,
        filter_dict: Optional[Dict[str, Any]] = None,
        limit: int = 100,
        offset: int = 0
    ) -> Tuple[List[Dict[str, Any]], int]:
        """
        List agents
        
        Args:
            organization_id: Organization ID
            user_id: User ID
            filter_dict: Optional filter criteria
            limit: Maximum number of results
            offset: Pagination offset
            
        Returns:
            Tuple of (list of agents, total count)
        """
        with LoggingContext(organization_id=organization_id, user_id=user_id):
            try:
                # Get vector database
                vector_db = get_vector_db(organization_id)
                
                # Set up base filter
                base_filter = {
                    "organization_id": organization_id,
                    "document_type": "agent_metadata"
                }
                
                # Merge with additional filters
                if filter_dict:
                    for key, value in filter_dict.items():
                        if key not in ["organization_id", "document_type"]:
                            base_filter[key] = value
                
                # Count total matches
                total_count = await vector_db.count_documents(base_filter)
                
                # Get agent documents
                docs = await vector_db.list_documents(
                    filter_dict=base_filter,
                    limit=limit,
                    offset=offset
                )
                
                # Process each agent
                agents = []
                for doc in docs:
                    try:
                        agent_metadata = json.loads(doc["text"])
                        agent_id = agent_metadata.get("agent_id")
                        
                        # Count associated documents
                        document_count = await vector_db.count_documents({
                            "associated_agents": agent_id,
                            "organization_id": organization_id
                        })
                        
                        # Format agent data
                        agent_data = {
                            "id": agent_id,
                            "name": agent_metadata.get("name", "Unnamed Agent"),
                            "description": agent_metadata.get("description"),
                            "instructions": agent_metadata.get("instructions", ""),
                            "created_at": agent_metadata.get("created_at"),
                            "updated_at": agent_metadata.get("updated_at"),
                            "metadata": {k: v for k, v in agent_metadata.items() if k not in [
                                "agent_id", "name", "description", "instructions", 
                                "created_at", "updated_at", "created_by", "organization_id", "type"
                            ]},
                            "document_count": document_count
                        }
                        
                        agents.append(agent_data)
                    except Exception as e:
                        logger.error(
                            f"Error processing agent document: {str(e)}",
                            extra={
                                "document_id": doc["id"],
                                "organization_id": organization_id
                            }
                        )
                
                return agents, total_count
            
            except Exception as e:
                logger.error(
                    f"Error listing agents: {str(e)}",
                    extra={
                        "organization_id": organization_id,
                        "user_id": user_id
                    }
                )
                raise
    
    async def delete_agent(
        self,
        agent_id: str,
        organization_id: str,
        user_id: str
    ) -> bool:
        """
        Delete an agent
        
        Args:
            agent_id: Agent ID
            organization_id: Organization ID
            user_id: User ID
            
        Returns:
            True if deleted, False otherwise
        """
        with LoggingContext(organization_id=organization_id, user_id=user_id):
            try:
                # Get vector database
                vector_db = get_vector_db(organization_id)
                
                # Find agent metadata document
                filter_dict = {
                    "agent_id": agent_id,
                    "organization_id": organization_id,
                    "document_type": "agent_metadata"
                }
                
                docs = await vector_db.list_documents(
                    filter_dict=filter_dict,
                    limit=1
                )
                
                if not docs:
                    logger.warning(
                        f"Agent not found for deletion: {agent_id}",
                        extra={
                            "agent_id": agent_id,
                            "organization_id": organization_id,
                            "user_id": user_id
                        }
                    )
                    return False
                
                # Get agent metadata
                agent_doc = docs[0]
                agent_metadata = json.loads(agent_doc["text"])
                
                # Check if user is authorized to delete
                created_by = agent_metadata.get("created_by")
                if created_by != user_id:
                    # TODO: Check if user is admin
                    is_admin = False
                    
                    if not is_admin:
                        logger.warning(
                            f"Unauthorized attempt to delete agent: {agent_id}",
                            extra={
                                "agent_id": agent_id,
                                "organization_id": organization_id,
                                "user_id": user_id,
                                "created_by": created_by
                            }
                        )
                        return False
                
                # Delete agent metadata document
                agent_doc_id = agent_doc["id"]
                success = await vector_db.delete_documents([agent_doc_id])
                
                if not success:
                    logger.error(
                        f"Failed to delete agent metadata document",
                        extra={
                            "agent_id": agent_id,
                            "document_id": agent_doc_id,
                            "organization_id": organization_id
                        }
                    )
                    return False
                
                # Remove agent from all associated documents
                # Find documents associated with this agent
                associated_docs = await vector_db.list_documents(
                    filter_dict={
                        "associated_agents": agent_id,
                        "organization_id": organization_id
                    },
                    limit=1000  # Set a reasonable limit
                )
                
                # Update each document to remove the agent association
                for doc in associated_docs:
                    doc_id = doc["id"]
                    metadata = doc["metadata"]
                    
                    if "associated_agents" in metadata and agent_id in metadata["associated_agents"]:
                        # Remove agent from list
                        metadata["associated_agents"].remove(agent_id)
                        
                        # If list is empty, remove it altogether
                        if not metadata["associated_agents"]:
                            del metadata["associated_agents"]
                        
                        # Update document
                        await vector_db.update_document(
                            document_id=doc_id,
                            metadata=metadata
                        )
                
                # Remove from agent cache if present
                if agent_id in self._agent_cache:
                    del self._agent_cache[agent_id]
                
                logger.info(
                    f"Deleted agent: {agent_id}",
                    extra={
                        "agent_id": agent_id,
                        "organization_id": organization_id,
                        "user_id": user_id,
                        "associated_docs_updated": len(associated_docs)
                    }
                )
                
                return True
            
            except Exception as e:
                logger.error(
                    f"Error deleting agent: {str(e)}",
                    extra={
                        "agent_id": agent_id,
                        "organization_id": organization_id,
                        "user_id": user_id
                    }
                )
                return False
    
    async def update_agent(
        self,
        agent_id: str,
        organization_id: str,
        user_id: str,
        name: Optional[str] = None,
        instructions: Optional[str] = None,
        description: Optional[str] = None,
        document_ids: Optional[List[str]] = None,
        document_filter: Optional[Dict[str, Any]] = None,
        metadata: Optional[Dict[str, Any]] = None
    ) -> Optional[Dict[str, Any]]:
        """
        Update an agent
        
        Args:
            agent_id: Agent ID
            organization_id: Organization ID
            user_id: User ID
            name: Optional new name
            instructions: Optional new instructions
            description: Optional new description
            document_ids: Optional new document IDs
            document_filter: Optional new document filter
            metadata: Optional new metadata
            
        Returns:
            Updated agent information or None if failed
        """
        with LoggingContext(organization_id=organization_id, user_id=user_id):
            try:
                # Get vector database
                vector_db = get_vector_db(organization_id)
                
                # Find agent metadata document
                filter_dict = {
                    "agent_id": agent_id,
                    "organization_id": organization_id,
                    "document_type": "agent_metadata"
                }
                
                docs = await vector_db.list_documents(
                    filter_dict=filter_dict,
                    limit=1
                )
                
                if not docs:
                    logger.warning(
                        f"Agent not found for update: {agent_id}",
                        extra={
                            "agent_id": agent_id,
                            "organization_id": organization_id,
                            "user_id": user_id
                        }
                    )
                    return None
                
                # Get agent metadata
                agent_doc = docs[0]
                agent_metadata = json.loads(agent_doc["text"])
                
                # Check if user is authorized to update
                created_by = agent_metadata.get("created_by")
                if created_by != user_id:
                    # TODO: Check if user is admin
                    is_admin = False
                    
                    if not is_admin:
                        logger.warning(
                            f"Unauthorized attempt to update agent: {agent_id}",
                            extra={
                                "agent_id": agent_id,
                                "organization_id": organization_id,
                                "user_id": user_id,
                                "created_by": created_by
                            }
                        )
                        return None
                
                # Update agent metadata
                updated = False
                
                if name is not None:
                    agent_metadata["name"] = name
                    updated = True
                
                if instructions is not None:
                    agent_metadata["instructions"] = instructions
                    updated = True
                
                if description is not None:
                    agent_metadata["description"] = description
                    updated = True
                
                if document_filter is not None:
                    agent_metadata["document_filter"] = document_filter
                    updated = True
                
                if metadata is not None:
                    # Update custom metadata fields
                    for key, value in metadata.items():
                        if key not in [
                            "agent_id", "name", "description", "instructions", 
                            "created_at", "updated_at", "created_by", "organization_id", "type"
                        ]:
                            agent_metadata[key] = value
                            updated = True
                
                if updated:
                    # Update timestamp
                    agent_metadata["updated_at"] = datetime.utcnow().isoformat()
                    
                    # Update agent metadata document
                    await vector_db.update_document(
                        document_id=agent_doc["id"],
                        text=json.dumps(agent_metadata)
                    )
                
                # Handle document associations if provided
                if document_ids is not None:
                    # First, remove agent from all currently associated documents
                    current_docs = await vector_db.list_documents(
                        filter_dict={
                            "associated_agents": agent_id,
                            "organization_id": organization_id
                        },
                        limit=1000  # Set a reasonable limit
                    )
                    
                    for doc in current_docs:
                        doc_id = doc["id"]
                        metadata = doc["metadata"]
                        
                        if "associated_agents" in metadata and agent_id in metadata["associated_agents"]:
                            # Remove agent from list
                            metadata["associated_agents"].remove(agent_id)
                            
                            # If list is empty, remove it altogether
                            if not metadata["associated_agents"]:
                                del metadata["associated_agents"]
                            
                            # Update document
                            await vector_db.update_document(
                                document_id=doc_id,
                                metadata=metadata
                            )
                    
                    # Then, associate new documents
                    for doc_id in document_ids:
                        # Get document to verify access
                        doc = await vector_db.get_document(doc_id)
                        
                        if not doc:
                            logger.warning(
                                f"Document not found: {doc_id}",
                                extra={
                                    "agent_id": agent_id,
                                    "organization_id": organization_id,
                                    "user_id": user_id
                                }
                            )
                            continue
                        
                        # Verify organization access
                        doc_org_id = doc["metadata"].get("organization_id")
                        if doc_org_id != organization_id:
                            logger.warning(
                                f"Attempted to associate document from another organization",
                                extra={
                                    "agent_id": agent_id,
                                    "document_id": doc_id,
                                    "organization_id": organization_id,
                                    "document_org_id": doc_org_id
                                }
                            )
                            continue
                        
                        # Add agent_id to document metadata
                        updated_metadata = dict(doc["metadata"])
                        
                        # Initialize or append to associated_agents list
                        if "associated_agents" not in updated_metadata:
                            updated_metadata["associated_agents"] = [agent_id]
                        elif agent_id not in updated_metadata["associated_agents"]:
                            updated_metadata["associated_agents"].append(agent_id)
                        
                        # Update document
                        await vector_db.update_document(
                            document_id=doc_id,
                            metadata=updated_metadata
                        )
                
                # Remove from agent cache if present
                if agent_id in self._agent_cache:
                    del self._agent_cache[agent_id]
                
                # Get updated agent info
                return await self.get_agent(agent_id, organization_id, user_id)
            
            except Exception as e:
                logger.error(
                    f"Error updating agent: {str(e)}",
                    extra={
                        "agent_id": agent_id,
                        "organization_id": organization_id,
                        "user_id": user_id
                    }
                )
                return None
    
    async def chat_with_agent(
        self,
        agent_id: str,
        organization_id: str,
        user_id: str,
        messages: List[Dict[str, str]],
        temperature: float = 0.7,
        max_tokens: Optional[int] = None,
        stream: bool = False,
        **kwargs: Any
    ) -> Union[str, AsyncIterator[str]]:
        """
        Chat with an agent
        
        Args:
            agent_id: Agent ID
            organization_id: Organization ID
            user_id: User ID
            messages: List of message dictionaries with 'role' and 'content'
            temperature: Controls randomness (0.0 = deterministic, 1.0 = random)
            max_tokens: Maximum number of tokens to generate
            stream: Whether to stream the response
            **kwargs: Additional parameters
            
        Returns:
            Agent's response or stream
        """
        with LoggingContext(organization_id=organization_id, user_id=user_id):
            try:
                # Get agent instance
                agent = await self._get_agent_instance(agent_id, organization_id, user_id)
                
                if not agent:
                    raise ValueError(f"Agent not found: {agent_id}")
                
                # Chat with agent
                if stream:
                    return await agent.chat_stream(
                        messages=messages,
                        temperature=temperature,
                        max_tokens=max_tokens,
                        **kwargs
                    )
                else:
                    return await agent.chat(
                        messages=messages,
                        temperature=temperature,
                        max_tokens=max_tokens,
                        **kwargs
                    )
            
            except Exception as e:
                logger.error(
                    f"Error chatting with agent: {str(e)}",
                    extra={
                        "agent_id": agent_id,
                        "organization_id": organization_id,
                        "user_id": user_id
                    }
                )
                raise
    
    async def _get_agent_instance(
        self,
        agent_id: str,
        organization_id: str,
        user_id: str
    ) -> Optional[RAGAgent]:
        """
        Get or create an agent instance
        
        Args:
            agent_id: Agent ID
            organization_id: Organization ID
            user_id: User ID
            
        Returns:
            Agent instance or None if not found
        """
        # Check cache first
        cache_key = f"{organization_id}:{agent_id}"
        if cache_key in self._agent_cache:
            return self._agent_cache[cache_key]
        
        # Get agent info
        agent_info = await self.get_agent(agent_id, organization_id, user_id)
        
        if not agent_info:
            return None
        
        # Create agent instance
        agent = RAGAgent(
            agent_id=agent_id,
            agent_name=agent_info["name"],
            agent_description=agent_info["description"],
            agent_instructions=agent_info["instructions"],
            organization_id=organization_id,
            user_id=user_id,
            metadata=agent_info["metadata"]
        )
        
        # Cache the agent
        self._agent_cache[cache_key] = agent
        
        return agent

# Create a global agent manager instance
agent_manager = AgentManager()