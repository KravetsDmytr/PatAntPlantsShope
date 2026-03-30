"use client";

import { Product } from "@/lib/types";

type Props = {
  product: Product;
  onAdd: (productId: number) => void;
};

export function ProductCard({ product, onAdd }: Props) {
  return (
    <div className="card" style={{ padding: 12 }}>
      <img
        src={product.image_url}
        alt={product.name}
        style={{ width: "100%", height: 150, objectFit: "cover", borderRadius: 8 }}
      />
      <h4 style={{ margin: "10px 0 6px" }}>{product.name}</h4>
      <p style={{ margin: 0, fontSize: 13, color: "#6b7280", minHeight: 40 }}>{product.description}</p>
      <div style={{ marginTop: 10, display: "flex", alignItems: "center", justifyContent: "space-between" }}>
        <strong>{product.cost} грн</strong>
        <button className="btn" onClick={() => onAdd(product.id)}>
          До кошика
        </button>
      </div>
    </div>
  );
}
