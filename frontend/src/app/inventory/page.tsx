"use client";

import { AppShell } from "@/components/layout/app-shell";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { ProductImage } from "@/components/ui/product-image";
import { ListSkeleton } from "@/components/ui/skeleton";
import { api } from "@/services/api";
import { useI18nStore } from "@/stores/i18n-store";
import { useToastStore } from "@/stores/toast-store";
import type { Branch, Inventory, Product } from "@/types/domain";
import { ClipboardPlus, Save, Search, X } from "lucide-react";
import { useEffect, useMemo, useState } from "react";

type StockForm = {
  branch_id: string;
  product_id: string;
  quantity_delta: number;
  reorder_threshold: number;
  reason: string;
  mode: "receive" | "adjust";
};

export default function InventoryPage() {
  const [branches, setBranches] = useState<Branch[]>([]);
  const [products, setProducts] = useState<Product[]>([]);
  const [inventories, setInventories] = useState<Inventory[]>([]);
  const [selectedBranchId, setSelectedBranchId] = useState("");
  const [query, setQuery] = useState("");
  const [form, setForm] = useState<StockForm>({ branch_id: "", product_id: "", quantity_delta: 1, reorder_threshold: 10, reason: "Owner stock update", mode: "receive" });
  const [formOpen, setFormOpen] = useState(false);
  const [loading, setLoading] = useState(true);
  const [notice, setNotice] = useState("");
  const [error, setError] = useState("");
  const t = useI18nStore((state) => state.t);

  const productMap = useMemo(() => new Map(products.map((product) => [product.id, product])), [products]);
  const branchMap = useMemo(() => new Map(branches.map((branch) => [branch.id, branch])), [branches]);
  const toast = useToastStore((state) => state.show);

  async function load(branchId = selectedBranchId, search = query) {
    setError("");
    try {
      setLoading(true);
      const [branchData, productData, inventoryData] = await Promise.all([
        api.myBranches(),
        api.products(search),
        api.inventories(search, branchId || undefined)
      ]);
      const branchList = Array.isArray(branchData) ? branchData : [];
      const productList = Array.isArray(productData) ? productData : [];
      setBranches(branchList);
      setProducts(productList);
      setInventories(Array.isArray(inventoryData) ? inventoryData : []);
      const nextBranch = branchId || branchList[0]?.id || "";
      setSelectedBranchId(nextBranch);
      setForm((value) => ({
        ...value,
        branch_id: value.branch_id || nextBranch,
        product_id: value.product_id || productList[0]?.id || ""
      }));
    } catch (err) {
      setError(err instanceof Error ? err.message : t("inventory.loadFailed"));
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    load("", "");
  }, []);

  useEffect(() => {
    if (selectedBranchId) {
      api.inventories(query, selectedBranchId)
        .then((data) => setInventories(Array.isArray(data) ? data : []))
        .catch((err) => setError(err instanceof Error ? err.message : t("inventory.loadFailed")));
    }
  }, [selectedBranchId]);

  function selectInventory(item: Inventory) {
    setForm({
      branch_id: item.branch_id,
      product_id: item.product_id,
      quantity_delta: 1,
      reorder_threshold: item.reorder_threshold,
      reason: t("inventory.ownerDefaultReason"),
      mode: "receive"
    });
    setFormOpen(true);
  }

  function openStockAction() {
    setNotice("");
    setError("");
    setForm({
      branch_id: selectedBranchId || branches[0]?.id || "",
      product_id: products[0]?.id || "",
      quantity_delta: 1,
      reorder_threshold: 10,
      reason: t("inventory.ownerDefaultReason"),
      mode: "receive"
    });
    setFormOpen(true);
  }

  function closeStockAction() {
    setFormOpen(false);
  }

  async function saveStock() {
    setNotice("");
    setError("");
    if (!form.branch_id || !form.product_id || !form.quantity_delta) {
      setError(t("inventory.required"));
      return;
    }
    const delta = form.mode === "receive" ? Math.abs(form.quantity_delta) : form.quantity_delta;
    try {
      if (form.mode === "receive") {
        await api.receiveStock(form.product_id, Math.abs(form.quantity_delta), form.reason, form.branch_id);
      } else {
        await api.adjustStock(form.branch_id, form.product_id, delta, form.reason);
      }
      await api.setReorderThreshold(form.branch_id, form.product_id, Math.max(0, Math.round(Number(form.reorder_threshold))));
      setNotice(t("inventory.updated"));
      toast({ type: "success", title: t("inventory.updated") });
      closeStockAction();
      await load(form.branch_id, query);
    } catch (err) {
      setError(err instanceof Error ? err.message : t("inventory.updateFailed"));
      toast({ type: "error", title: t("inventory.updateFailed"), message: err instanceof Error ? err.message : t("common.tryAgain") });
    }
  }

  async function saveThresholdOnly() {
    setNotice("");
    setError("");
    if (!form.branch_id || !form.product_id) {
      setError(t("inventory.thresholdRequired"));
      return;
    }
    try {
      await api.setReorderThreshold(form.branch_id, form.product_id, Math.max(0, Math.round(Number(form.reorder_threshold))));
      setNotice(t("inventory.reorderThresholdUpdated"));
      toast({ type: "success", title: t("inventory.thresholdUpdated") });
      closeStockAction();
      await load(form.branch_id, query);
    } catch (err) {
      setError(err instanceof Error ? err.message : t("inventory.thresholdFailed"));
      toast({ type: "error", title: t("inventory.thresholdFailed"), message: err instanceof Error ? err.message : t("common.tryAgain") });
    }
  }

  return (
    <AppShell>
      <div className="mb-5 flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold">{t("inventory.title")}</h1>
          <p className="text-sm text-slate-500">{t("inventory.description")}</p>
        </div>
        <div className="flex flex-col gap-2 sm:flex-row">
          <select className="h-10 w-full rounded-md border border-line bg-white px-3 text-sm sm:w-72" value={selectedBranchId} onChange={(event) => setSelectedBranchId(event.target.value)}>
            {branches.map((branch) => <option key={branch.id} value={branch.id}>{branch.code} · {branch.name}</option>)}
          </select>
          <Button onClick={openStockAction}>
            <ClipboardPlus className="h-4 w-4" />
            {t("inventory.stockAction")}
          </Button>
        </div>
      </div>
      {notice ? <div className="mb-4 rounded-md border border-emerald-200 bg-emerald-50 p-3 text-sm text-emerald-800">{notice}</div> : null}
      {error ? <div className="mb-4 rounded-md border border-red-200 bg-red-50 p-3 text-sm text-red-700">{error}</div> : null}

      <div>
        <section className="space-y-4">
          <div className="rounded-md border border-line bg-white p-4">
            <div className="flex flex-col gap-2 sm:flex-row">
              <Input placeholder={t("inventory.searchPlaceholder")} value={query} onChange={(event) => setQuery(event.target.value)} onKeyDown={(event) => event.key === "Enter" && load(selectedBranchId, query)} />
              <Button onClick={() => load(selectedBranchId, query)}><Search className="h-4 w-4" />{t("common.search")}</Button>
            </div>
          </div>
          {loading ? <ListSkeleton rows={4} /> : (
          <div className="overflow-hidden rounded-md border border-line bg-white shadow-sm">
            {inventories.length === 0 ? <div className="p-5 text-sm text-slate-500">{t("inventory.empty")}</div> : null}
            {inventories.map((item) => {
              const product = productMap.get(item.product_id);
              const branch = branchMap.get(item.branch_id);
              const lowStock = item.quantity <= item.reorder_threshold;
              return (
                <button key={item.id} className="grid w-full gap-3 border-b border-line p-4 text-left last:border-b-0 hover:bg-field md:grid-cols-[1fr_120px_120px_120px] md:items-center" onClick={() => selectInventory(item)}>
                  <div className="flex min-w-0 items-center gap-3">
                    <ProductImage src={product?.image_url} name={product?.name ?? item.product_id} />
                    <div className="min-w-0">
                      <div className="truncate font-semibold">{product?.name ?? item.product_id}</div>
                      <div className="text-xs text-slate-500">{branch?.code ?? t("settings.branchFallback")} · {product ? `${product.sku} · ${product.barcode}` : t("inventory.detailsUnavailable")}</div>
                      {lowStock ? <div className="mt-1 text-xs font-semibold text-red-600">{t("inventory.lowStock")}</div> : null}
                    </div>
                  </div>
                  <div>
                    <div className="text-xs text-slate-500">{t("inventory.onHand")}</div>
                    <div className={lowStock ? "text-xl font-bold text-red-600" : "text-xl font-bold"}>{item.quantity}</div>
                  </div>
                  <div className="md:text-right">
                    <div className="text-xs text-slate-500">{t("inventory.reserved")}</div>
                    <div className="font-semibold">{item.reserved_quantity}</div>
                  </div>
                  <div className="md:text-right">
                    <div className="text-xs text-slate-500">{t("inventory.reorderAt")}</div>
                    <div className="font-semibold">{item.reorder_threshold}</div>
                  </div>
                </button>
              );
            })}
          </div>
          )}
        </section>
      </div>
      {formOpen ? (
        <div className="fixed inset-0 z-50 grid place-items-center bg-slate-950/50 px-4 py-6 backdrop-blur-sm">
          <div className="max-h-[calc(100vh-48px)] w-full max-w-xl overflow-y-auto rounded-md border border-line bg-white p-5 shadow-2xl">
            <div className="mb-4 flex items-start justify-between gap-3">
              <div>
                <h2 className="flex items-center gap-2 text-lg font-bold"><ClipboardPlus className="h-5 w-5 text-brand" />{t("inventory.stockAction")}</h2>
                <p className="mt-1 text-sm text-slate-500">{t("inventory.actionDescription")}</p>
              </div>
              <button className="grid h-9 w-9 shrink-0 place-items-center rounded-md border border-line text-slate-600 hover:bg-field" onClick={closeStockAction} aria-label={t("common.close")}>
                <X className="h-4 w-4" />
              </button>
            </div>
            <div className="space-y-3">
              <label className="block text-sm font-semibold">{t("field.branch")}
                <select className="mt-1 h-10 w-full rounded-md border border-line bg-white px-3 text-sm" value={form.branch_id} onChange={(event) => setForm({ ...form, branch_id: event.target.value })}>
                  {branches.map((branch) => <option key={branch.id} value={branch.id}>{branch.code} · {branch.name}</option>)}
                </select>
              </label>
              <label className="block text-sm font-semibold">{t("field.product")}
                <select className="mt-1 h-10 w-full rounded-md border border-line bg-white px-3 text-sm" value={form.product_id} onChange={(event) => setForm({ ...form, product_id: event.target.value })}>
                  {products.map((product) => <option key={product.id} value={product.id}>{product.sku} · {product.name}</option>)}
                </select>
              </label>
              <label className="block text-sm font-semibold">{t("field.action")}
                <select className="mt-1 h-10 w-full rounded-md border border-line bg-white px-3 text-sm" value={form.mode} onChange={(event) => setForm({ ...form, mode: event.target.value as StockForm["mode"] })}>
                  <option value="receive">{t("inventory.receiveStock")}</option>
                  <option value="adjust">{t("inventory.adjustStock")}</option>
                </select>
              </label>
              <label className="block text-sm font-semibold">{t("field.quantityDelta")}<Input className="mt-1" type="number" value={form.quantity_delta} onChange={(event) => setForm({ ...form, quantity_delta: Number(event.target.value) })} /></label>
              <label className="block text-sm font-semibold">{t("field.reorderThreshold")}<Input className="mt-1" type="number" min={0} value={form.reorder_threshold} onChange={(event) => setForm({ ...form, reorder_threshold: Number(event.target.value) })} /></label>
              <label className="block text-sm font-semibold">{t("field.reason")}<Input className="mt-1" value={form.reason} onChange={(event) => setForm({ ...form, reason: event.target.value })} /></label>
            </div>
            <div className="mt-5 flex flex-col-reverse gap-2 sm:flex-row sm:justify-end">
              <Button className="!bg-white !text-slate-700 ring-1 ring-line hover:!bg-field" onClick={closeStockAction}>{t("common.cancel")}</Button>
              <Button className="!bg-white !text-slate-700 ring-1 ring-line hover:!bg-field" onClick={saveThresholdOnly}>{t("inventory.saveThresholdOnly")}</Button>
              <Button onClick={saveStock}><Save className="h-4 w-4" />{t("inventory.applyStockUpdate")}</Button>
            </div>
          </div>
        </div>
      ) : null}
    </AppShell>
  );
}
