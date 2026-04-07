import os
import time
from typing import Any

import pytest
import requests


BASE_URL = os.getenv("API_BASE_URL", "http://localhost:8081/api/v1")


@pytest.fixture(scope="session")
def api_base_url() -> str:
    return BASE_URL


@pytest.fixture(scope="session", autouse=True)
def wait_for_api(api_base_url: str) -> None:
    deadline = time.time() + 90
    last_error: Any = None

    while time.time() < deadline:
        try:
            response = requests.get(f"{api_base_url}/categories", timeout=3)
            if response.status_code in (200, 401):
                return
        except requests.RequestException as exc:
            last_error = exc
        time.sleep(2)

    raise RuntimeError(f"API is not ready at {api_base_url}. Last error: {last_error}")
