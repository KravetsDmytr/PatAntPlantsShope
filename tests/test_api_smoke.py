import time

import requests


def _uniq_login(prefix: str = "jenkins_user") -> str:
    return f"{prefix}_{int(time.time() * 1000)}"


def test_categories_positive(api_base_url: str) -> None:
    response = requests.get(f"{api_base_url}/categories", timeout=10)
    assert response.status_code == 200

    payload = response.json()
    assert payload["error"] is None
    assert isinstance(payload["data"], list)
    assert len(payload["data"]) > 0


def test_login_negative_invalid_password(api_base_url: str) -> None:
    login = _uniq_login()
    register_payload = {
        "login": login,
        "first_name": "Jenkins",
        "last_name": "Bot",
        "email": f"{login}@mail.test",
        "password": "correct_password_123",
    }
    reg = requests.post(f"{api_base_url}/auth/register", json=register_payload, timeout=10)
    assert reg.status_code in (200, 201), reg.text

    bad_login = requests.post(
        f"{api_base_url}/auth/login",
        json={"login": login, "password": "wrong_password"},
        timeout=10,
    )
    assert bad_login.status_code == 401
    payload = bad_login.json()
    assert payload["data"] is None
    assert payload["error"] is not None


def test_product_details_positive(api_base_url: str) -> None:
    products = requests.get(f"{api_base_url}/products", timeout=10)
    assert products.status_code == 200
    products_payload = products.json()
    assert isinstance(products_payload["data"], list)
    assert len(products_payload["data"]) > 0

    product_id = products_payload["data"][0]["id"]
    detail = requests.get(f"{api_base_url}/products/{product_id}", timeout=10)
    assert detail.status_code == 200
    detail_payload = detail.json()
    assert detail_payload["error"] is None
    assert detail_payload["data"]["id"] == product_id
