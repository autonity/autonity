"""Simple JSON-RPC client."""

from typing import List, Any

import requests


class RPCError(RuntimeError):
    """The JSON-RPC response has an 'error' field instead of a 'result' field.
    """


class Client:
    """Simple Web3 client for sending JSON-RPC requests."""
    nonce: int = 0
    url: str
    version: str = "2.0"

    def __init__(self, port: int):
        self.url = f"http://127.0.0.1:{port}"

    def request(self, method: str, params: List[str]) -> Any:
        """Sends a JSON-RPC request and returns the 'result' field of the
        response."""
        nonce = self.nonce
        response = requests.post(self.url, json={
            "jsonrpc": self.version,
            "method": method,
            "params": params,
            "id": nonce,
        })
        response.raise_for_status()
        self.nonce += 1
        data = response.json()
        assert data["id"] == nonce
        assert data["jsonrpc"] == self.version
        if "error" in data:
            raise RPCError(data["error"])
        return data["result"]
