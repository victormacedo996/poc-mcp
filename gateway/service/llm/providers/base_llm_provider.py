from abc import ABC
from abc import abstractmethod


class BaseLLMProvider(ABC):
    """
    Base class for all LLM providers.
    """
    @abstractmethod
    def call_model(self, prompt: str) -> str:
        raise NotImplementedError("call_model() must be implemented in the subclass")