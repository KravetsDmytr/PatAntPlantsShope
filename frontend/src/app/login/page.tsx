"use client";

import { Header } from "@/components/Header";
import { api } from "@/lib/api";
import { authStore } from "@/lib/auth";
import { useRouter } from "next/navigation";
import { useState } from "react";

export default function LoginPage() {
  const router = useRouter();
  const [login, setLogin] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  async function submit(e: React.FormEvent) {
    e.preventDefault();
    try {
      setError("");
      const res = await api.login({ login, password });
      authStore.setToken(res.data.token);
      router.push("/");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Помилка входу");
    }
  }

  return (
    <main className="container">
      <Header />
      <form onSubmit={submit} className="card" style={{ maxWidth: 420, margin: "30px auto", padding: 16 }}>
        <h2>Вхід</h2>
        <input className="input" placeholder="Логін" value={login} onChange={(e) => setLogin(e.target.value)} />
        <div style={{ height: 8 }} />
        <input
          className="input"
          placeholder="Пароль"
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
        {error && <p style={{ color: "crimson" }}>{error}</p>}
        <button className="btn" type="submit">
          Увійти
        </button>
      </form>
    </main>
  );
}
