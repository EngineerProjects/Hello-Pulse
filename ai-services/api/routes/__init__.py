"""
API routes package
"""

# Import all route modules to make them available
from . import generate, rag, search, agents

all_routes = [
    generate.router,
    rag.router,
    search.router,
    agents.router
]