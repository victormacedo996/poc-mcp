from service.llm.providers.base_llm_provider import BaseLLMProvider
import requests
from typing import Generator, Union
import json


class OllamaModelProvider(BaseLLMProvider):
    """
    Ollama LLM provider.
    """

    def call_model(self, prompt: str, stream: bool) -> Union[str, Generator[str, None, None]]:
        """
        Call the model with the given prompt.
        Returns a string if stream is False, otherwise a generator yielding streamed response chunks.
        """
        response = requests.post(
            url="http://localhost:11434/api/generate",
            json={
                "model": "llama3.2:3b",
                "prompt": prompt,
                "stream": stream
            },
            stream=stream
        )

        response.raise_for_status()

        if stream:
            def stream_generator():
                for line in response.iter_lines(decode_unicode=True):
                    if line:
                        yield json.loads(line)["response"]
            return stream_generator()
        else:
            return response.json()["response"]
        



