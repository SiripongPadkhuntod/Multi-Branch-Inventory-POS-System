import type { ApiResponse, Branch, CartItem, Category, DashboardSummary, EmployeeSalesSummary, Inventory, InventoryMovementDetail, PaymentMethod, Product, ProductStockSummary, Role, Sale, SaleDetail, User } from "@/types/domain";

const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080/api/v1";

export class ApiError extends Error {
  constructor(message: string, public status: number) {
    super(message);
  }
}

let redirectingToLogin = false;

function clearBrowserSession() {
  if (typeof window === "undefined") {
    return;
  }
  localStorage.removeItem("access_token");
  localStorage.removeItem("refresh_token");
  localStorage.removeItem("user");
}

function redirectToLogin() {
  if (typeof window === "undefined" || redirectingToLogin || window.location.pathname === "/login") {
    return;
  }
  redirectingToLogin = true;
  clearBrowserSession();
  fetch(`${API_URL}/auth/logout`, {
    method: "POST",
    credentials: "include",
    headers: { "Content-Type": "application/json" }
  }).finally(() => {
    window.location.assign("/login?session=expired");
  });
}

async function request<T>(path: string, init: RequestInit = {}): Promise<T> {
  const token = typeof window !== "undefined" ? localStorage.getItem("access_token") : null;
  const response = await fetch(`${API_URL}${path}`, {
    ...init,
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...init.headers
    }
  });
  const payload = (await response.json()) as ApiResponse<T>;
  if (!response.ok || !payload.success) {
    if (response.status === 401 && !path.startsWith("/auth/login") && !path.startsWith("/auth/logout")) {
      redirectToLogin();
    }
    throw new ApiError(payload.message || "Request failed", response.status);
  }
  return payload.data;
}

export const api = {
  login: (email: string, password: string) =>
    request<{ access_token: string; refresh_token: string; user: User }>("/auth/login", {
      method: "POST",
      body: JSON.stringify({ email, password })
    }),
  logout: () => request<null>("/auth/logout", { method: "POST" }),
  me: () => request<User>("/auth/me"),
  myBranches: () => request<Branch[]>("/branches/my"),
  createBranch: (input: { code: string; name: string; address: string; phone: string; status: string }) =>
    request<Branch>("/branches", { method: "POST", body: JSON.stringify(input) }),
  updateBranch: (id: string, input: { code: string; name: string; address: string; phone: string; status: string }) =>
    request<Branch>(`/branches/${id}`, { method: "PUT", body: JSON.stringify(input) }),
  categories: () => request<Category[]>("/categories"),
  createCategory: (input: { name: string; description: string }) =>
    request<Category>("/categories", { method: "POST", body: JSON.stringify(input) }),
  updateCategory: (id: string, input: { name: string; description: string }) =>
    request<Category>(`/categories/${id}`, { method: "PUT", body: JSON.stringify(input) }),
  deleteCategory: (id: string) => request<null>(`/categories/${id}`, { method: "DELETE" }),
  dashboardSummary: (branchId?: string) =>
    request<DashboardSummary>(`/dashboard/summary${branchId ? `?branch_id=${encodeURIComponent(branchId)}` : ""}`),
  users: () => request<User[]>("/users"),
  employeeSalesSummary: () => request<EmployeeSalesSummary[]>("/users/sales-summary"),
  createUser: (input: { name: string; email: string; password: string; role: Role; branch_id: string; status: string }) =>
    request<User>("/users", { method: "POST", body: JSON.stringify(input) }),
  updateUser: (id: string, input: { name: string; email: string; role: Role; branch_id: string; status: string }) =>
    request<User>(`/users/${id}`, { method: "PUT", body: JSON.stringify(input) }),
  products: (q = "") => request<Product[]>(`/products?q=${encodeURIComponent(q)}`),
  productByBarcode: (barcode: string) => request<Product>(`/products/barcode/${encodeURIComponent(barcode)}`),
  createProduct: (input: Omit<Product, "id">) => request<Product>("/products", { method: "POST", body: JSON.stringify(input) }),
  updateProduct: (id: string, input: Omit<Product, "id">) => request<Product>(`/products/${id}`, { method: "PUT", body: JSON.stringify(input) }),
  deleteProduct: (id: string) => request<null>(`/products/${id}`, { method: "DELETE" }),
  inventories: (q = "", branchId?: string, categoryId?: string) =>
    request<Inventory[]>(
      `/inventories?q=${encodeURIComponent(q)}${branchId ? `&branch_id=${encodeURIComponent(branchId)}` : ""}${categoryId ? `&category_id=${encodeURIComponent(categoryId)}` : ""}`
    ),
  inventoryMovements: (q = "", branchId?: string) =>
    request<InventoryMovementDetail[]>(`/inventories/movements?q=${encodeURIComponent(q)}${branchId ? `&branch_id=${encodeURIComponent(branchId)}` : ""}`),
  allStock: (q = "") => request<ProductStockSummary[]>(`/inventories/all-stock?q=${encodeURIComponent(q)}`),
  adjustStock: (branchId: string, productId: string, quantityDelta: number, reason: string) =>
    request<null>("/inventories/adjust", {
      method: "POST",
      body: JSON.stringify({
        branch_id: branchId,
        product_id: productId,
        quantity_delta: quantityDelta,
        reason
      })
    }),
  receiveStock: (productId: string, quantity: number, reason: string, branchId?: string) =>
    request<null>("/inventories/receive", {
      method: "POST",
      body: JSON.stringify({
        branch_id: branchId,
        product_id: productId,
        quantity_delta: quantity,
        reason
      })
    }),
  transferStock: (fromBranchId: string, toBranchId: string, productId: string, quantity: number) =>
    request<null>("/inventories/transfer", {
      method: "POST",
      body: JSON.stringify({
        from_branch_id: fromBranchId,
        to_branch_id: toBranchId,
        product_id: productId,
        quantity
      })
    }),
  sales: (dateFrom?: string, dateTo?: string) =>
    request<Sale[]>(`/sales${dateFrom || dateTo ? `?${new URLSearchParams({ ...(dateFrom ? { date_from: dateFrom } : {}), ...(dateTo ? { date_to: dateTo } : {}) }).toString()}` : ""}`),
  branchSales: (branchId?: string) => request<Sale[]>(`/sales/branch${branchId ? `?branch_id=${encodeURIComponent(branchId)}` : ""}`),
  saleDetail: (id: string) => request<SaleDetail>(`/sales/${id}`),
  createSale: (branchId: string, items: CartItem[], method: PaymentMethod, discount: number, tax: number) => {
    const subtotal = items.reduce((sum, item) => sum + item.final_price * item.quantity, 0);
    const total = Math.max(0, subtotal - discount + tax);
    return request<Sale>("/sales", {
      method: "POST",
      body: JSON.stringify({
        branch_id: branchId,
        discount,
        tax,
        items: items.map((item) => ({
          product_id: item.product.id,
          quantity: item.quantity,
          final_price: item.final_price,
          discount_amount: item.discount_amount,
          discount_reason: item.discount_reason ?? ""
        })),
        payments: [{ payment_method: method, amount: total }]
      })
    });
  }
};
