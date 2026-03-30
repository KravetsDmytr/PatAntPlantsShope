"use client";

import { Header } from "@/components/Header";
import { api } from "@/lib/api";
import { useRouter } from "next/navigation";
import { useState } from "react";

export default function RegisterPage() {
  const router = useRouter();
  const [form, setForm] = useState({
    login: "",
    first_name: "",
    last_name: "",
    email: "",
    password: ""
  });
  const [error, setError] = useState("");

  async function submit(e: React.FormEvent) {
    e.preventDefault();
    try {
      setError("");
      await api.register(form);
      router.push("/login");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Помилка реєстрації");
    }
  }

  return (
    <main className="container">
      <Header />
      <form onSubmit={submit} className="card" style={{ maxWidth: 480, margin: "30px auto", padding: 16, display: "grid", gap: 8 }}>
        <h2>Реєстрація</h2>
        <input className="input" placeholder="Логін" value={form.login} onChange={(e) => setForm({ ...form, login: e.target.value })} />
        <input
          className="input"
          placeholder="Ім'я"
          value={form.first_name}
          onChange={(e) => setForm({ ...form, first_name: e.target.value })}
        />
        <input
          className="input"
          placeholder="Прізвище"
          value={form.last_name}
          onChange={(e) => setForm({ ...form, last_name: e.target.value })}
        />
        <input className="input" placeholder="Email" value={form.email} onChange={(e) => setForm({ ...form, email: e.target.value })} />
        <input
          className="input"
          type="password"
          placeholder="Пароль"
          value={form.password}
          onChange={(e) => setForm({ ...form, password: e.target.value })}
        />
        {error && <p style={{ color: "crimson" }}>{error}</p>}
        <button className="btn" type="submit">
          Зареєструватись
        </button>
      </form>
    </main>
  );
}
