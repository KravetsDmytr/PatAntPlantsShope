"use client";

import { Header } from "@/components/Header";
import { api } from "@/lib/api";
import { authStore } from "@/lib/auth";
import { Product } from "@/lib/types";
import Link from "next/link";
import { useParams } from "next/navigation";
import { useEffect, useState } from "react";

export default function ProductDetailsPage() {
  const params = useParams<{ id: string }>();
  const id = Number(params.id);

  const [product, setProduct] = useState<Product | null>(null);
  const [error, setError] = useState("");

  useEffect(() => {
    if (!id || Number.isNaN(id)) {
      setError("Невірний id товару");
      return;
    }

    api
      .productById(id)
      .then((res) => setProduct(res.data))
      .catch((e) => setError(e instanceof Error ? e.message : "Помилка завантаження товару"));
  }, [id]);

  async function addToCart() {
    if (!product) return;
    const token = authStore.getToken();
    if (!token) {
      alert("Спочатку виконайте вхід");
      return;
    }

    try {
      await api.addToCart(token, { product_id: product.id, quantity: 1 });
      alert("Товар додано до кошика");
    } catch (e) {
      alert(e instanceof Error ? e.message : "Помилка додавання");
    }
  }

  return (
    <main className="container">
      <Header />
      <Link href="/" style={{ color: "#6b7280" }}>
        ← Назад до каталогу
      </Link>

      {error && <p style={{ color: "crimson" }}>{error}</p>}

      {product && (
        <section className="card" style={{ marginTop: 12, padding: 16, display: "grid", gridTemplateColumns: "360px 1fr", gap: 18 }}>
          <img
            src={product.image_url}
            alt={product.name}
            style={{ width: "100%", height: 320, objectFit: "cover", borderRadius: 12 }}
          />

          <div>
            <h1 style={{ marginTop: 0 }}>{product.name}</h1>
            {product.title && <p style={{ color: "#6b7280", marginTop: -6 }}>{product.title}</p>}
            <p>{product.description}</p>

            <div style={{ display: "grid", gridTemplateColumns: "170px 1fr", rowGap: 8, columnGap: 10 }}>
              <strong>Вартість:</strong>
              <span>{product.cost} грн</span>

              <strong>Виробник:</strong>
              <span>{product.producer || "-"}</span>

              <strong>Тип:</strong>
              <span>{product.type || "-"}</span>

              <strong>Вага:</strong>
              <span>{product.weight ? `${product.weight} ${product.unit ?? ""}` : "-"}</span>

              <strong>Гарантія:</strong>
              <span>{product.guarantee ? new Date(product.guarantee).toLocaleDateString("uk-UA") : "-"}</span>

              <strong>Створено:</strong>
              <span>{new Date(product.created_at).toLocaleString("uk-UA")}</span>
            </div>

            <div style={{ marginTop: 16, display: "flex", gap: 10 }}>
              <button className="btn" onClick={addToCart}>
                Додати до кошика
              </button>
              <Link className="btn secondary" href="/cart">
                До кошика
              </Link>
            </div>
          </div>
        </section>
      )}
    </main>
  );
}
