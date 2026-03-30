"use client";

import Link from "next/link";
import { useEffect, useState } from "react";
import { authStore } from "@/lib/auth";

export function Header() {
  const [isAuth, setIsAuth] = useState(false);

  useEffect(() => {
    setIsAuth(Boolean(authStore.getToken()));
  }, []);

  return (
    <header className="card" style={{ padding: 12, marginBottom: 16 }}>
      <div style={{ display: "flex", gap: 16, alignItems: "center" }}>
        <Link href="/" style={{ fontWeight: 700, color: "#5f8f1f" }}>
          Каталог
        </Link>
        <div style={{ marginLeft: "auto", display: "flex", gap: 12 }}>
          <Link href="/cart">Кошик</Link>
          {isAuth ? (
            <button
              className="btn secondary"
              onClick={() => {
                authStore.clear();
                setIsAuth(false);
              }}
            >
              Вийти
            </button>
          ) : (
            <>
              <Link href="/login">Вхід</Link>
              <Link href="/register">Реєстрація</Link>
            </>
          )}
        </div>
      </div>
    </header>
  );
}
