from pydantic import BaseModel

class ChatRequest(BaseModel):
    """
    Request model for chat endpoint.
    """
    prompt: str
    temperature: float = 0.7
    stream: bool = False
    