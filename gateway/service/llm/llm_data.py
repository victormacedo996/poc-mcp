from repository.sqlite import get_available_llms, get_available_mcp_servers, add_mcp_server
from typing import List
from models.available_llms import AvailableLLM

def get_available_llm_models() -> List[AvailableLLM]:
    """
    Fetches available LLM models from the database.
    """
    llm_models = get_available_llms()
    
    available_llms = list()
    for model in llm_models:
        available_llms.append(AvailableLLM(
            name=model[0],
            size=model[1],
            num_parameters=model[2],
            description=model[3],
            pros=model[4],
            cons=model[5]
        ))

    return available_llms

def insert_llm_model(name: str, size: str, num_parameters: int, description: str, pros: str, cons: str) -> None:
    """
    Inserts a new LLM model into the database.
    """
    add_mcp_server(name, size, num_parameters, description, pros, cons)