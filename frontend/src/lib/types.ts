export type ApiResponse<T> = {
  message: string;
  data: T;
  error: string | null;
};

export type Category = {
  id: number;
  name: string;
};

export type Product = {
  id: number;
  name: string;
  producer: string;
  type: string;
  description: string;
  title: string | null;
  cost: number;
  weight: number | null;
  unit: string | null;
  guarantee: string | null;
  image_url: string;
  category_id: number;
  created_at: string;
  updated_at: string;
};

export type CartItem = {
  product_id: number;
  product_name: string;
  cost: number;
  quantity: number;
  line_total: number;
};

export type CartData = {
  items: CartItem[];
  amount_cost: number;
};
