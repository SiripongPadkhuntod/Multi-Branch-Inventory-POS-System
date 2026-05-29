"use client";

import type { CartItem, Product } from "@/types/domain";
import { create } from "zustand";

type CartState = {
  items: CartItem[];
  discount: number;
  addProduct: (product: Product) => void;
  updateQuantity: (productId: string, quantity: number) => void;
  overridePrice: (productId: string, price: number, reason: string) => void;
  setDiscount: (discount: number) => void;
  clear: () => void;
};

export const useCartStore = create<CartState>((set) => ({
  items: [],
  discount: 0,
  addProduct: (product) =>
    set((state) => {
      const existing = state.items.find((item) => item.product.id === product.id);
      if (existing) {
        return {
          items: state.items.map((item) =>
            item.product.id === product.id ? { ...item, quantity: item.quantity + 1 } : item
          )
        };
      }
      return {
        items: [...state.items, { product, quantity: 1, final_price: product.sell_price, discount_amount: 0 }]
      };
    }),
  updateQuantity: (productId, quantity) =>
    set((state) => ({
      items: quantity <= 0 ? state.items.filter((item) => item.product.id !== productId) : state.items.map((item) =>
        item.product.id === productId ? { ...item, quantity } : item
      )
    })),
  overridePrice: (productId, price, reason) =>
    set((state) => ({
      items: state.items.map((item) =>
        item.product.id === productId
          ? {
              ...item,
              final_price: price,
              discount_amount: Math.max(0, item.product.sell_price - price),
              discount_reason: reason
            }
          : item
      )
    })),
  setDiscount: (discount) => set({ discount }),
  clear: () => set({ items: [], discount: 0 })
}));
