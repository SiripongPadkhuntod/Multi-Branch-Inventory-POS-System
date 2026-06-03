export type Role = "OWNER" | "MANAGER" | "EMPLOYEE";
export type PaymentMethod = "CASH" | "PROMPTPAY" | "BANK_TRANSFER" | "CREDIT_CARD";

export type ApiResponse<T> = {
  success: boolean;
  message: string;
  data: T;
};

export type User = {
  id: string;
  branch_id: string | null;
  branch_ids?: string[];
  role: Role;
  name: string;
  email: string;
  status: string;
};

export type AuditLog = {
  id: string;
  user_id: string;
  user_name: string;
  user_email: string;
  action: string;
  entity_type: string;
  entity_id: string;
  old_data: string;
  new_data: string;
  ip_address: string;
  created_at: string;
};

export type Branch = {
  id: string;
  code: string;
  name: string;
  address: string;
  phone: string;
  status: string;
  created_at: string;
};

export type Category = {
  id: string;
  name: string;
  description: string;
};

export type Product = {
  id: string;
  sku: string;
  barcode: string;
  qr_code: string;
  name: string;
  description: string;
  category_id: string;
  image_url: string;
  cost_price: number;
  sell_price: number;
  status: string;
};

export type Inventory = {
  id: string;
  branch_id: string;
  product_id: string;
  quantity: number;
  reserved_quantity: number;
  reorder_threshold: number;
  updated_at: string;
};

export type MovementType = "RECEIVE" | "SALE" | "RETURN" | "ADJUSTMENT" | "TRANSFER_IN" | "TRANSFER_OUT";

export type InventoryMovementDetail = {
  id: string;
  branch_id: string;
  branch_code: string;
  branch_name: string;
  product_id: string;
  sku: string;
  barcode: string;
  product_name: string;
  movement_type: MovementType;
  quantity: number;
  reference_id: string;
  created_by: string;
  created_by_name: string;
  created_at: string;
};

export type BranchStockDetail = {
  branch_id: string;
  branch_code: string;
  branch_name: string;
  quantity: number;
  reserved_quantity: number;
  reorder_threshold: number;
  updated_at: string;
};

export type ProductStockSummary = {
  product_id: string;
  sku: string;
  barcode: string;
  image_url: string;
  product_name: string;
  total_quantity: number;
  total_reserved: number;
  branches: BranchStockDetail[];
};

export type Sale = {
  id: string;
  receipt_number: string;
  branch_id: string;
  employee_id: string;
  subtotal: number;
  discount: number;
  tax: number;
  total: number;
  payment_status: string;
  refund_status: "NONE" | "PARTIAL_REFUND" | "REFUNDED";
  created_at: string;
};

export type SaleItemDetail = {
  id: string;
  product_id: string;
  sku: string;
  barcode: string;
  product_name: string;
  quantity: number;
  returned_quantity: number;
  original_price: number;
  final_price: number;
  discount_amount: number;
  discount_reason: string;
  line_total: number;
};

export type PaymentDetail = {
  id: string;
  payment_method: PaymentMethod;
  amount: number;
};

export type SaleDetail = Sale & {
  branch_code: string;
  branch_name: string;
  employee_name: string;
  items: SaleItemDetail[];
  payments: PaymentDetail[];
};

export type TopProduct = {
  product_id: string;
  name: string;
  quantity: number;
  revenue: number;
};

export type LowStockItem = {
  branch_id: string;
  branch_code: string;
  product_id: string;
  product_name: string;
  quantity: number;
  reorder_threshold: number;
};

export type BranchSalesSummary = {
  branch_id: string;
  branch_code: string;
  branch_name: string;
  revenue: number;
  sales_count: number;
};

export type DashboardSummary = {
  daily_sales: number;
  monthly_sales: number;
  revenue: number;
  profit: number;
  low_stock: LowStockItem[];
  top_products: TopProduct[];
  branch_comparison: BranchSalesSummary[];
};

export type EmployeeSalesSummary = {
  user_id: string;
  branch_id: string;
  branch_code: string;
  name: string;
  email: string;
  role: Role;
  status: string;
  sales_count: number;
  revenue: number;
};

export type CartItem = {
  product: Product;
  quantity: number;
  final_price: number;
  discount_amount: number;
  discount_reason?: string;
};
