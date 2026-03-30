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
  description: string;
  cost: number;
  image_url: string;
  category_id: number;
  created_at: string;
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
