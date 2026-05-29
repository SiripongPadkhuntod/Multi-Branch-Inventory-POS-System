"use client";

import { AppShell } from "@/components/layout/app-shell";
import { api } from "@/services/api";
import type { Branch, DashboardSummary, EmployeeSalesSummary, Sale, SaleDetail } from "@/types/domain";
import { BarChart3, ChevronDown, ChevronUp, Package, ReceiptText, Users } from "lucide-react";
import { useEffect, useMemo, useState } from "react";

const money = (value: number) => `฿${(value / 100).toLocaleString("th-TH", { minimumFractionDigits: 2 })}`;

const emptySummary: DashboardSummary = {
  daily_sales: 0,
  monthly_sales: 0,
  revenue: 0,
  profit: 0,
  low_stock: [],
  top_products: [],
  branch_comparison: []
};

export default function ReportsPage() {
  const [branches, setBranches] = useState<Branch[]>([]);
  const [branchId, setBranchId] = useState("");
  const [summary, setSummary] = useState<DashboardSummary>(emptySummary);
  const [sales, setSales] = useState<Sale[]>([]);
  const [employeeSales, setEmployeeSales] = useState<EmployeeSalesSummary[]>([]);
  const [selectedSaleId, setSelectedSaleId] = useState("");
  const [detail, setDetail] = useState<SaleDetail | null>(null);
  const [error, setError] = useState("");

  const selectedBranch = useMemo(() => branches.find((branch) => branch.id === branchId), [branches, branchId]);

  useEffect(() => {
    api.myBranches()
      .then((data) => setBranches(Array.isArray(data) ? data : []))
      .catch((err) => setError(err instanceof Error ? err.message : "Cannot load branches"));
  }, []);

  useEffect(() => {
    setError("");
    setSelectedSaleId("");
    setDetail(null);
    Promise.all([api.dashboardSummary(branchId || undefined), api.branchSales(branchId || undefined), api.employeeSalesSummary()])
      .then(([summaryData, salesData, employeeData]) => {
        setSummary(summaryData ?? emptySummary);
        setSales(Array.isArray(salesData) ? salesData : []);
        setEmployeeSales(Array.isArray(employeeData) ? employeeData : []);
      })
      .catch((err) => setError(err instanceof Error ? err.message : "Cannot load reports"));
  }, [branchId]);

  async function toggleDetail(sale: Sale) {
    if (selectedSaleId === sale.id) {
      setSelectedSaleId("");
      setDetail(null);
      return;
    }
    setSelectedSaleId(sale.id);
    setDetail(null);
    try {
      setDetail(await api.saleDetail(sale.id));
    } catch (err) {
      setError(err instanceof Error ? err.message : "Cannot load receipt detail");
    }
  }

  const filteredEmployeeSales = branchId ? employeeSales.filter((item) => item.branch_id === branchId) : employeeSales;

  return (
    <AppShell>
      <div className="mb-5 flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold">Reports</h1>
          <p className="text-sm text-slate-500">{selectedBranch ? `${selectedBranch.code} · ${selectedBranch.name}` : "Sales, products, employees, and branch performance."}</p>
        </div>
        <select className="h-10 w-full rounded-md border border-line bg-white px-3 text-sm sm:w-72" value={branchId} onChange={(event) => setBranchId(event.target.value)}>
          <option value="">All branches</option>
          {branches.map((branch) => <option key={branch.id} value={branch.id}>{branch.code} · {branch.name}</option>)}
        </select>
      </div>
      {error ? <div className="mb-4 rounded-md border border-red-200 bg-red-50 p-3 text-sm text-red-700">{error}</div> : null}

      <section className="grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
        {[
          ["Daily Sales", summary.daily_sales.toLocaleString()],
          ["Monthly Sales", summary.monthly_sales.toLocaleString()],
          ["Revenue", money(summary.revenue)],
          ["Profit", money(summary.profit)]
        ].map(([label, value]) => (
          <article key={label} className="rounded-md border border-line bg-white p-4">
            <div className="flex items-center justify-between text-sm text-slate-500"><span>{label}</span><BarChart3 className="h-5 w-5 text-brand" /></div>
            <div className="mt-3 text-2xl font-bold">{value}</div>
          </article>
        ))}
      </section>

      <section className="mt-5 grid gap-4 xl:grid-cols-3">
        <div className="rounded-md border border-line bg-white p-4">
          <h2 className="font-semibold">Top Products</h2>
          <div className="mt-3 space-y-2">
            {summary.top_products.length === 0 ? <div className="text-sm text-slate-500">No product sales yet.</div> : null}
            {summary.top_products.map((product) => (
              <div key={product.product_id} className="flex justify-between gap-3 rounded-md bg-field p-3 text-sm">
                <span className="font-medium">{product.name}</span>
                <span>{product.quantity} · {money(product.revenue)}</span>
              </div>
            ))}
          </div>
        </div>
        <div className="rounded-md border border-line bg-white p-4">
          <h2 className="font-semibold">Low Stock</h2>
          <div className="mt-3 space-y-2">
            {summary.low_stock.length === 0 ? <div className="text-sm text-slate-500">No low stock items.</div> : null}
            {summary.low_stock.map((item) => (
              <div key={`${item.branch_id}-${item.product_id}`} className="flex justify-between gap-3 rounded-md bg-field p-3 text-sm">
                <span className="font-medium">{item.branch_code} · {item.product_name}</span>
                <span>{item.quantity}</span>
              </div>
            ))}
          </div>
        </div>
        <div className="rounded-md border border-line bg-white p-4">
          <h2 className="font-semibold">Sales By Employee</h2>
          <div className="mt-3 space-y-2">
            {filteredEmployeeSales.length === 0 ? <div className="text-sm text-slate-500">No employee sales yet.</div> : null}
            {filteredEmployeeSales.map((item) => (
              <div key={item.user_id} className="rounded-md bg-field p-3 text-sm">
                <div className="flex items-center gap-2 font-semibold"><Users className="h-4 w-4 text-brand" />{item.name}</div>
                <div className="mt-1 flex justify-between text-xs text-slate-500"><span>{item.branch_code} · {item.email}</span><span>{item.sales_count} · {money(item.revenue)}</span></div>
              </div>
            ))}
          </div>
        </div>
      </section>

      <section className="mt-5 overflow-hidden rounded-md border border-line bg-white">
        <div className="border-b border-line p-4 font-semibold">Receipts</div>
        {sales.length === 0 ? <div className="p-5 text-sm text-slate-500">No sales yet.</div> : null}
        {sales.map((sale) => (
          <div key={sale.id} className="border-b border-line last:border-b-0">
            <button className="flex w-full flex-wrap items-center justify-between gap-3 p-4 text-left hover:bg-field" onClick={() => toggleDetail(sale)}>
              <div className="flex min-w-0 items-center gap-3">
                <ReceiptText className="h-5 w-5 shrink-0 text-brand" />
                <div className="min-w-0">
                  <div className="truncate font-semibold">{sale.receipt_number}</div>
                  <div className="text-xs text-slate-500">{new Date(sale.created_at).toLocaleString()}</div>
                </div>
              </div>
              <div className="flex items-center gap-4 text-right">
                <div><div className="font-bold">{money(sale.total)}</div><div className="text-xs text-emerald-700">{sale.payment_status}</div></div>
                {selectedSaleId === sale.id ? <ChevronUp className="h-4 w-4" /> : <ChevronDown className="h-4 w-4" />}
              </div>
            </button>
            {selectedSaleId === sale.id && detail ? (
              <div className="border-t border-line bg-slate-50 p-4">
                <div className="mb-3 grid gap-3 text-sm md:grid-cols-3">
                  <div className="rounded-md border border-line bg-white p-3"><div className="text-xs text-slate-500">Branch</div><div className="font-semibold">{detail.branch_code} · {detail.branch_name}</div></div>
                  <div className="rounded-md border border-line bg-white p-3"><div className="text-xs text-slate-500">Employee</div><div className="font-semibold">{detail.employee_name}</div></div>
                  <div className="rounded-md border border-line bg-white p-3"><div className="text-xs text-slate-500">Total</div><div className="font-semibold">{money(detail.total)}</div></div>
                </div>
                <div className="overflow-hidden rounded-md border border-line bg-white">
                  {detail.items.map((item) => (
                    <div key={item.id} className="grid gap-3 border-b border-line p-3 text-sm last:border-b-0 md:grid-cols-[1fr_80px_120px] md:items-center">
                      <div className="flex items-start gap-3"><Package className="mt-0.5 h-4 w-4 text-brand" /><div><div className="font-semibold">{item.product_name}</div><div className="text-xs text-slate-500">{item.sku} · {item.barcode}</div></div></div>
                      <div>{item.quantity} pcs</div>
                      <div className="font-bold md:text-right">{money(item.line_total)}</div>
                    </div>
                  ))}
                </div>
              </div>
            ) : null}
          </div>
        ))}
      </section>
    </AppShell>
  );
}
