from fastmcp import Client
from fastmcp.client.transports import SSETransport


async def main():

    client = Client(SSETransport("http://localhost:8000/sse"))
    async with client:
        result = await client.list_tools()
        print(result)


        
if __name__ == "__main__":
    import asyncio
    asyncio.run(main())