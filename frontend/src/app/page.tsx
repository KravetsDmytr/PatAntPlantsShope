"use client";

import { Header } from "@/components/Header";
import { ProductCard } from "@/components/ProductCard";
import { api } from "@/lib/api";
import { authStore } from "@/lib/auth";
import { Category, Product } from "@/lib/types";
import { useEffect, useState } from "react";

export default function HomePage() {
  const [categories, setCategories] = useState<Category[]>([]);
  const [products, setProducts] = useState<Product[]>([]);
  const [categoryId, setCategoryId] = useState<number | undefined>(undefined);
  const [q, setQ] = useState("");
  const [minCost, setMinCost] = useState("");
  const [maxCost, setMaxCost] = useState("");
  const [error, setError] = useState("");

  async function loadData() {
    try {
      setError("");
      const [c, p] = await Promise.all([
        api.categories(),
        api.products({
          category_id: categoryId,
          q: q || undefined,
          min_cost: minCost ? Number(minCost) : undefined,
          max_cost: maxCost ? Number(maxCost) : undefined
        })
      ]);
      setCategories(c.data || []);
      setProducts(p.data || []);
    } catch (e) {
      setError(e instanceof Error ? e.message : "Помилка завантаження");
    }
  }

  useEffect(() => {
    void loadData();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [categoryId]);

  async function addToCart(productId: number) {
    const token = authStore.getToken();
    if (!token) {
      alert("Спочатку виконайте вхід");
      return;
    }
    try {
      await api.addToCart(token, { product_id: productId, quantity: 1 });
      alert("Товар додано до кошика");
    } catch (e) {
      alert(e instanceof Error ? e.message : "Помилка додавання");
    }
  }

  return (
    <main className="container">
      <Header />

      <section style={{ display: "grid", gridTemplateColumns: "220px 1fr", gap: 16 }}>
        <aside className="card" style={{ padding: 12 }}>
          <h3 style={{ marginTop: 0 }}>Категорії</h3>
          <button className="btn secondary" style={{ width: "100%", marginBottom: 8 }} onClick={() => setCategoryId(undefined)}>
            Усі товари
          </button>
          {categories.map((c) => (
            <button
              key={c.id}
              className="btn secondary"
              style={{ width: "100%", marginBottom: 8, border: categoryId === c.id ? "2px solid #92c746" : undefined }}
              onClick={() => setCategoryId(c.id)}
            >
              {c.name}
            </button>
          ))}
        </aside>

        <div>
          <div className="card" style={{ padding: 12, marginBottom: 12, display: "grid", gap: 8, gridTemplateColumns: "1fr 140px 140px auto" }}>
            <input className="input" placeholder="Я шукаю..." value={q} onChange={(e) => setQ(e.target.value)} />
            <input className="input" placeholder="min" value={minCost} onChange={(e) => setMinCost(e.target.value)} />
            <input className="input" placeholder="max" value={maxCost} onChange={(e) => setMaxCost(e.target.value)} />
            <button className="btn" onClick={() => void loadData()}>
              Фільтр
            </button>
          </div>

          {error && <p style={{ color: "crimson" }}>{error}</p>}

          <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fill, minmax(220px, 1fr))", gap: 12 }}>
            {products.map((p) => (
              <ProductCard key={p.id} product={p} onAdd={addToCart} />
            ))}
          </div>
        </div>
      </section>
    </main>
  );
}
