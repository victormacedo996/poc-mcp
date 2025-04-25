from pydantic import BaseModel

class AvailableLLM(BaseModel):
    name: str
    size: str
    num_parameters: str
    description: str
    pros: str
    cons: str