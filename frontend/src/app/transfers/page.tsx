"use client";

import { AppShell } from "@/components/layout/app-shell";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { api } from "@/services/api";
import { useI18nStore } from "@/stores/i18n-store";
import type { Branch, Inventory, Product } from "@/types/domain";
import { ArrowRightLeft, Boxes, Save } from "lucide-react";
import { useEffect, useMemo, useState } from "react";

export default function TransfersPage() {
  const [branches, setBranches] = useState<Branch[]>([]);
  const [products, setProducts] = useState<Product[]>([]);
  const [sourceStock, setSourceStock] = useState<Inventory[]>([]);
  const [fromBranchId, setFromBranchId] = useState("");
  const [toBranchId, setToBranchId] = useState("");
  const [productId, setProductId] = useState("");
  const [quantity, setQuantity] = useState(1);
  const [notice, setNotice] = useState("");
  const [error, setError] = useState("");
  const t = useI18nStore((state) => state.t);

  const productMap = useMemo(() => new Map(products.map((product) => [product.id, product])), [products]);

  async function load() {
    setError("");
    try {
      const [branchData, productData] = await Promise.all([api.myBranches(), api.products()]);
      const branchList = Array.isArray(branchData) ? branchData : [];
      const productList = Array.isArray(productData) ? productData : [];
      setBranches(branchList);
      setProducts(productList);
      setFromBranchId((value) => value || branchList[0]?.id || "");
      setToBranchId((value) => value || branchList[1]?.id || branchList[0]?.id || "");
      setProductId((value) => value || productList[0]?.id || "");
    } catch (err) {
      setError(err instanceof Error ? err.message : t("transfers.failed"));
    }
  }

  useEffect(() => {
    load();
  }, []);

  useEffect(() => {
    if (fromBranchId) {
      api.inventories("", fromBranchId)
        .then((data) => setSourceStock(Array.isArray(data) ? data : []))
        .catch((err) => setError(err instanceof Error ? err.message : t("inventory.loadFailed")));
    }
  }, [fromBranchId]);

  async function transfer() {
    setNotice("");
    setError("");
    if (!fromBranchId || !toBranchId || !productId || quantity < 1) {
      setError(t("inventory.required"));
      return;
    }
    if (fromBranchId === toBranchId) {
      setError(t("transfers.failed"));
      return;
    }
    try {
      await api.transferStock(fromBranchId, toBranchId, productId, quantity);
      setNotice(t("transfers.transferred"));
      const data = await api.inventories("", fromBranchId);
      setSourceStock(Array.isArray(data) ? data : []);
    } catch (err) {
      setError(err instanceof Error ? err.message : t("transfers.failed"));
    }
  }

  return (
    <AppShell>
      <div className="mb-5">
        <h1 className="text-2xl font-bold">{t("transfers.title")}</h1>
        <p className="text-sm text-slate-500">{t("transfers.description")}</p>
      </div>
      {notice ? <div className="mb-4 rounded-md border border-emerald-200 bg-emerald-50 p-3 text-sm text-emerald-800">{notice}</div> : null}
      {error ? <div className="mb-4 rounded-md border border-red-200 bg-red-50 p-3 text-sm text-red-700">{error}</div> : null}
      <div className="grid gap-4 xl:grid-cols-[1fr_380px]">
        <section className="overflow-hidden rounded-md border border-line bg-white">
          <div className="border-b border-line p-4 font-semibold">{t("transfers.sourceStock")}</div>
          {sourceStock.length === 0 ? <div className="p-5 text-sm text-slate-500">{t("transfers.emptySource")}</div> : null}
          {sourceStock.map((item) => {
            const product = productMap.get(item.product_id);
            return (
              <button key={item.id} className="flex w-full items-center justify-between gap-3 border-b border-line p-4 text-left last:border-b-0 hover:bg-field" onClick={() => setProductId(item.product_id)}>
                <div className="flex min-w-0 items-center gap-3">
                  <Boxes className="h-5 w-5 shrink-0 text-brand" />
                  <div className="min-w-0">
                    <div className="truncate font-semibold">{product?.name ?? item.product_id}</div>
                    <div className="text-xs text-slate-500">{product ? `${product.sku} · ${product.barcode}` : t("inventory.detailsUnavailable")}</div>
                  </div>
                </div>
                <div className="text-right">
                  <div className="text-xl font-bold">{item.quantity}</div>
                  <div className="text-xs text-slate-500">{t("inventory.onHand")}</div>
                </div>
              </button>
            );
          })}
        </section>
        <aside className="rounded-md border border-line bg-white p-4">
          <h2 className="mb-4 flex items-center gap-2 text-lg font-bold"><ArrowRightLeft className="h-5 w-5 text-brand" />{t("transfers.new")}</h2>
          <div className="space-y-3">
            <label className="block text-sm font-semibold">From {t("field.branch")}
              <select className="mt-1 h-10 w-full rounded-md border border-line bg-white px-3 text-sm" value={fromBranchId} onChange={(event) => setFromBranchId(event.target.value)}>
                {branches.map((branch) => <option key={branch.id} value={branch.id}>{branch.code} · {branch.name}</option>)}
              </select>
            </label>
            <label className="block text-sm font-semibold">To {t("field.branch")}
              <select className="mt-1 h-10 w-full rounded-md border border-line bg-white px-3 text-sm" value={toBranchId} onChange={(event) => setToBranchId(event.target.value)}>
                {branches.map((branch) => <option key={branch.id} value={branch.id}>{branch.code} · {branch.name}</option>)}
              </select>
            </label>
            <label className="block text-sm font-semibold">{t("field.product")}
              <select className="mt-1 h-10 w-full rounded-md border border-line bg-white px-3 text-sm" value={productId} onChange={(event) => setProductId(event.target.value)}>
                {products.map((product) => <option key={product.id} value={product.id}>{product.sku} · {product.name}</option>)}
              </select>
            </label>
            <label className="block text-sm font-semibold">{t("field.quantity")}<Input className="mt-1" type="number" min={1} value={quantity} onChange={(event) => setQuantity(Number(event.target.value))} /></label>
            <Button className="w-full" onClick={transfer}><Save className="h-4 w-4" />{t("transfers.transferStock")}</Button>
          </div>
        </aside>
      </div>
    </AppShell>
  );
}
