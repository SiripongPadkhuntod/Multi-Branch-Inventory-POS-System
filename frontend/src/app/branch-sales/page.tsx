"use client";

import { AppShell } from "@/components/layout/app-shell";
import { Button } from "@/components/ui/button";
import { api } from "@/services/api";
import type { Branch, Sale, SaleDetail } from "@/types/domain";
import { ChevronDown, ChevronUp, CreditCard, Package, ReceiptText, Users } from "lucide-react";
import { useEffect, useMemo, useState } from "react";

const money = (value: number) => `฿${(value / 100).toLocaleString("th-TH", { minimumFractionDigits: 2 })}`;

export default function BranchSalesPage() {
  const [branches, setBranches] = useState<Branch[]>([]);
  const [selectedBranchId, setSelectedBranchId] = useState("");
  const [sales, setSales] = useState<Sale[]>([]);
  const [selectedSaleId, setSelectedSaleId] = useState("");
  const [detail, setDetail] = useState<SaleDetail | null>(null);
  const [loadingDetail, setLoadingDetail] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    api.myBranches()
      .then((data) => {
        const list = Array.isArray(data) ? data : [];
        setBranches(list);
        if (list.length === 1) {
          setSelectedBranchId(list[0].id);
        }
      })
      .catch((err) => setError(err instanceof Error ? err.message : "Cannot load branches"));
  }, []);

  useEffect(() => {
    setError("");
    setSelectedSaleId("");
    setDetail(null);
    api.branchSales(selectedBranchId || undefined)
      .then((data) => setSales(Array.isArray(data) ? data : []))
      .catch((err) => {
        setSales([]);
        setError(err instanceof Error ? err.message : "Cannot load branch sales");
      });
  }, [selectedBranchId]);

  const selectedBranch = useMemo(
    () => branches.find((branch) => branch.id === selectedBranchId),
    [branches, selectedBranchId]
  );

  async function toggleDetail(sale: Sale) {
    if (selectedSaleId === sale.id) {
      setSelectedSaleId("");
      setDetail(null);
      return;
    }
    setSelectedSaleId(sale.id);
    setDetail(null);
    setError("");
    setLoadingDetail(true);
    try {
      setDetail(await api.saleDetail(sale.id));
    } catch (err) {
      setError(err instanceof Error ? err.message : "Cannot load receipt detail");
    } finally {
      setLoadingDetail(false);
    }
  }

  return (
    <AppShell>
      <div className="mb-5 flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold">Branch Sales</h1>
          <p className="text-sm text-slate-500">
            {selectedBranch ? `${selectedBranch.code} · ${selectedBranch.name}` : "Receipts from branches you manage."}
          </p>
        </div>
        <select
          className="h-10 w-full rounded-md border border-line bg-white px-3 text-sm sm:w-72"
          value={selectedBranchId}
          onChange={(event) => setSelectedBranchId(event.target.value)}
        >
          <option value="">All managed branches</option>
          {branches.map((branch) => (
            <option key={branch.id} value={branch.id}>
              {branch.code} · {branch.name}
            </option>
          ))}
        </select>
      </div>
      {error ? <div className="mb-4 rounded-md border border-red-200 bg-red-50 p-3 text-sm text-red-700">{error}</div> : null}
      <div className="overflow-hidden rounded-md border border-line bg-white">
        {sales.length === 0 ? (
          <div className="p-5 text-sm text-slate-500">No branch sales yet.</div>
        ) : (
          sales.map((sale) => (
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
                  <div>
                    <div className="font-bold">{money(sale.total)}</div>
                    <div className="text-xs text-emerald-700">{sale.payment_status}</div>
                  </div>
                  {selectedSaleId === sale.id ? <ChevronUp className="h-4 w-4" /> : <ChevronDown className="h-4 w-4" />}
                </div>
              </button>
              {selectedSaleId === sale.id ? (
                <div className="border-t border-line bg-slate-50 p-4">
                  {loadingDetail ? <div className="text-sm text-slate-500">Loading receipt detail...</div> : null}
                  {detail ? (
                    <div className="space-y-4">
                      <div className="grid gap-3 text-sm md:grid-cols-3">
                        <div className="rounded-md border border-line bg-white p-3">
                          <div className="text-xs text-slate-500">Branch</div>
                          <div className="font-semibold">{detail.branch_code} · {detail.branch_name}</div>
                        </div>
                        <div className="rounded-md border border-line bg-white p-3">
                          <div className="flex items-center gap-2 text-xs text-slate-500">
                            <Users className="h-3.5 w-3.5" />
                            Employee
                          </div>
                          <div className="font-semibold">{detail.employee_name}</div>
                        </div>
                        <div className="rounded-md border border-line bg-white p-3">
                          <div className="text-xs text-slate-500">Receipt Time</div>
                          <div className="font-semibold">{new Date(detail.created_at).toLocaleString()}</div>
                        </div>
                      </div>
                      <div className="overflow-hidden rounded-md border border-line bg-white">
                        {detail.items.map((item) => (
                          <div key={item.id} className="grid gap-3 border-b border-line p-3 text-sm last:border-b-0 md:grid-cols-[1fr_90px_120px_120px] md:items-center">
                            <div className="flex items-start gap-3">
                              <Package className="mt-0.5 h-4 w-4 shrink-0 text-brand" />
                              <div>
                                <div className="font-semibold">{item.product_name}</div>
                                <div className="text-xs text-slate-500">{item.sku} · {item.barcode}</div>
                                {item.discount_reason ? <div className="mt-1 text-xs text-accent">{item.discount_reason}</div> : null}
                              </div>
                            </div>
                            <div>
                              <div className="text-xs text-slate-500 md:hidden">Qty</div>
                              <div className="font-semibold">{item.quantity}</div>
                            </div>
                            <div>
                              <div className="text-xs text-slate-500 md:hidden">Unit Price</div>
                              <div>{money(item.final_price)}</div>
                              {item.discount_amount > 0 ? <div className="text-xs text-slate-500">Discount {money(item.discount_amount)}</div> : null}
                            </div>
                            <div className="font-bold md:text-right">{money(item.line_total)}</div>
                          </div>
                        ))}
                      </div>
                      <div className="grid gap-4 md:grid-cols-[1fr_280px]">
                        <div className="rounded-md border border-line bg-white p-3">
                          <div className="mb-2 flex items-center gap-2 font-semibold">
                            <CreditCard className="h-4 w-4 text-brand" />
                            Payments
                          </div>
                          {detail.payments.map((payment) => (
                            <div key={payment.id} className="flex justify-between py-1 text-sm">
                              <span>{payment.payment_method}</span>
                              <span className="font-semibold">{money(payment.amount)}</span>
                            </div>
                          ))}
                        </div>
                        <div className="rounded-md border border-line bg-white p-3 text-sm">
                          <div className="flex justify-between py-1"><span>Subtotal</span><span>{money(detail.subtotal)}</span></div>
                          <div className="flex justify-between py-1"><span>Discount</span><span>{money(detail.discount)}</span></div>
                          <div className="flex justify-between py-1"><span>Tax</span><span>{money(detail.tax)}</span></div>
                          <div className="mt-2 flex justify-between border-t border-line pt-2 text-lg font-bold"><span>Total</span><span>{money(detail.total)}</span></div>
                        </div>
                      </div>
                      <Button className="bg-slate-800 hover:bg-slate-700" onClick={() => window.print()}>
                        Print Receipt
                      </Button>
                    </div>
                  ) : null}
                </div>
              ) : null}
            </div>
          ))
        )}
      </div>
    </AppShell>
  );
}
