from pydantic import BaseModel

class LLMInsertRequest(BaseModel):
    name: str
    size: str
    num_parameters: int
    description: str
    pros: str
    cons: str