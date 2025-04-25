from repository.sqlite import get_available_mcp_servers as repo_get_available_mcp_servers
from repository.sqlite import add_mcp_server
from fastmcp import Client
from fastmcp.client.transports import SSETransport

def get_available_mcp_servers():
    """
    Fetches available MCP servers from the database.
    """
    return repo_get_available_mcp_servers()

def add_mcp_server(url, environment):
    """
    Adds a new MCP server to the database.
    """
    add_mcp_server(url, environment)

async def get_available_tools_from_mcp_server(mcp_server_url: list[str]):
    """
    Fetches available tools from a specific MCP server.
    """     
    async with Client(SSETransport(mcp_server_url)) as client:
        result = await client.list_tools()
        return result

async def call_mcp_tool(mcp_server_url: str, tool_name: str, **kwargs):
    """
    Calls a specific tool on the MCP server with the provided input.
    """
    async with Client(SSETransport(mcp_server_url)) as client:
        result = await client.call_tool(tool_name, arguments=kwargs)
        return result
