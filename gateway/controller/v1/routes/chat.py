from fastapi import APIRouter
from controller.v1.dto.request.chat_request import ChatRequest
from service.mcp_servers.handle_mcp_servers import get_available_tools_from_mcp_server, call_mcp_tool
from service.llm.providers.ollama import OllamaModelProvider
from string import Template
from fastapi.responses import StreamingResponse
import json


router = APIRouter()

TOOL_PROMPT = """
SYSTEM:  
You are an AI assistant with access to the following MCP tools (each tool follows MCP’s JSON-RPC 2.0 schema):

TOOLS:
$tools

When a user asks a question, follow these steps:

1. **Analyze** the user’s request.  
2. **Decide**:
- If the question requires external data or actions beyond your training (e.g. live weather, database lookup, computation), choose exactly one tool from the above list.
- Otherwise, plan to answer directly without calling any tool.

3. **Output** **strictly** one ARRAY containing JSON OBJECTS (no extra text) in one of these forms:

– **Tool call**  
[
  { "name": <tool_name>, "arguments": { … } }
]

User’s question begins below.

USER: $user_question

"""

FINAL_PROMPT = """
You are an expert assistant.  
You have been given:  

1. The original user question:  
   $user_question

2. The raw MCP tool outputs (already fetched via MCP’s tools/call), in JSON form:  
   $mcp_tool_outputs_json

Your task is to consume those MCP results and craft a coherent, self-contained answer to the user.  

Assistant (plan):
1. Parse the user question.
2. Examine the MCP_TOOL_OUTPUTS_JSON; identify relevant entries.
3. Summarize or compare the key fields from those entries.
4. Structure the final response: include context, findings, and recommendations.

Assistant (analysis):
– Refer to each MCP result by its index (e.g. “Result #0 shows …”).  
– Extract the most salient fields (e.g. title, summary, metrics).  
– Note any discrepancies or highlights across results.

Assistant (answer):
Provide a clear, concise answer to “$user_question” that weaves in the MCP results. Cite each result by index.  
"""



@router.post("/achat")
async def chat(chat_request: ChatRequest) -> StreamingResponse:
    """
    Chat endpoint
    """
    tools = await get_available_tools_from_mcp_server("http://localhost:8000/sse")
    model = OllamaModelProvider()
    tools_prompt = Template(TOOL_PROMPT).substitute(tools=tools, user_question=chat_request.prompt)
    async def event_stream():
        model_response = model.call_model(tools_prompt, stream=False)
        tools_2_call = json.loads(model_response)
        tools_results = list()
        for tool in tools_2_call:
            yield f"Calling Tool: {tool['name']} with arguments {tool['arguments']}\n"
            tool_call_result = await call_mcp_tool("http://localhost:8000/sse", tool['name'], **tool['arguments'])
            tools_results.append(tool_call_result[0].model_dump())
        
        yield f"Calling expert model\n"
        str_tools_results = "\n".join([str(tool) for tool in tools_results])
        print(str_tools_results)
        for chunk in model.call_model(Template(FINAL_PROMPT).substitute(user_question=chat_request.prompt, mcp_tool_outputs_json=str_tools_results), stream=True):
            yield chunk

    return StreamingResponse(event_stream(), media_type="text/event-stream")  


@router.post("/chat")
async def chat(chat_request: ChatRequest) -> dict:
    """
    Chat endpoint
    """
    tools = await get_available_tools_from_mcp_server("http://localhost:8000/sse")
    model = OllamaModelProvider()
    tools_prompt = Template(TOOL_PROMPT).substitute(tools=tools, user_question=chat_request.prompt)
    model_response = model.call_model(tools_prompt, stream=False)
    tools_2_call = json.loads(model_response)
    tools_results = list()
    for tool in tools_2_call:
        tool_call_result = await call_mcp_tool("http://localhost:8000/sse", tool['name'], **tool['arguments'])
        tools_results.append(tool_call_result[0].model_dump())

    str_tools_results = "\n".join([str(tool) for tool in tools_results])
    final_response = model.call_model(Template(FINAL_PROMPT).substitute(user_question=chat_request.prompt, mcp_tool_outputs_json=str_tools_results), stream=False)

    

    return {
        "final_response": final_response,
        "tools_2_call": tools_2_call, 
    }