import { ApiResponse, CartData, Category, Product } from "./types";

const BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080/api/v1";

async function request<T>(path: string, init?: RequestInit): Promise<ApiResponse<T>> {
  const res = await fetch(`${BASE_URL}${path}`, {
    ...init,
    headers: {
      "Content-Type": "application/json",
      ...(init?.headers ?? {})
    },
    cache: "no-store"
  });

  const data = (await res.json()) as ApiResponse<T>;
  if (!res.ok) {
    throw new Error(data?.error || data?.message || "API error");
  }
  return data;
}

export const api = {
  register: (payload: {
    login: string;
    first_name: string;
    last_name: string;
    email: string;
    password: string;
  }) => request<null>("/auth/register", { method: "POST", body: JSON.stringify(payload) }),

  login: (payload: { login: string; password: string }) =>
    request<{ token: string }>("/auth/login", { method: "POST", body: JSON.stringify(payload) }),

  categories: () => request<Category[]>("/categories"),

  products: (filters: { category_id?: number; min_cost?: number; max_cost?: number; min_price?: number; max_price?: number; q?: string }) => {
    const p = new URLSearchParams();
    if (filters.category_id) p.set("category_id", String(filters.category_id));
    if (filters.min_cost) p.set("min_cost", String(filters.min_cost));
    else if (filters.min_price) p.set("min_cost", String(filters.min_price));

    if (filters.max_cost) p.set("max_cost", String(filters.max_cost));
    else if (filters.max_price) p.set("max_cost", String(filters.max_price));

    if (filters.q) p.set("q", filters.q);
    return request<Product[]>(`/products${p.toString() ? `?${p.toString()}` : ""}`);
  },

  productById: (id: number) => request<Product>(`/products/${id}`),

  addToCart: (token: string, payload: { product_id: number; quantity: number }) =>
    request<null>("/cart/items", {
      method: "POST",
      body: JSON.stringify(payload),
      headers: { Authorization: `Bearer ${token}` }
    }),

  cart: (token: string) =>
    request<CartData>("/cart", {
      headers: { Authorization: `Bearer ${token}` }
    })
};
