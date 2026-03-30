"use client";

import { Header } from "@/components/Header";
import { api } from "@/lib/api";
import { authStore } from "@/lib/auth";
import { CartData } from "@/lib/types";
import { useEffect, useState } from "react";

export default function CartPage() {
  const [cart, setCart] = useState<CartData>({ items: [], amount_cost: 0 });
  const [error, setError] = useState("");

  useEffect(() => {
    const token = authStore.getToken();
    if (!token) {
      setError("Виконайте вхід для перегляду кошика");
      return;
    }

    api.cart(token)
      .then((res) => setCart(res.data))
      .catch((e) => setError(e instanceof Error ? e.message : "Помилка завантаження кошика"));
  }, []);

  return (
    <main className="container">
      <Header />
      <div className="card" style={{ padding: 16 }}>
        <h2>Кошик</h2>
        {error && <p style={{ color: "crimson" }}>{error}</p>}
        {cart.items.map((item) => (
          <div
            key={item.product_id}
            style={{ display: "grid", gridTemplateColumns: "1fr 120px 120px", borderBottom: "1px solid #e5e7eb", padding: "10px 0" }}
          >
            <div>{item.product_name}</div>
            <div>{item.quantity} шт.</div>
            <div>{item.line_total} грн</div>
          </div>
        ))}
        <h3 style={{ textAlign: "right" }}>Разом: {cart.amount_cost} грн</h3>
      </div>
    </main>
  );
}
