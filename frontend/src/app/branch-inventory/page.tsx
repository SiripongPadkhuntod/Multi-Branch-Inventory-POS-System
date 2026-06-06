"use client";

import { AppShell } from "@/components/layout/app-shell";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { ProductImage } from "@/components/ui/product-image";
import { ListSkeleton } from "@/components/ui/skeleton";
import { api } from "@/services/api";
import { useI18nStore } from "@/stores/i18n-store";
import type { Branch, Category, Inventory, Product } from "@/types/domain";
import { Search, X } from "lucide-react";
import { useEffect, useMemo, useState } from "react";

export default function BranchInventoryPage() {
  const [branches, setBranches] = useState<Branch[]>([]);
  const [selectedBranchId, setSelectedBranchId] = useState("");
  const [inventories, setInventories] = useState<Inventory[]>([]);
  const [products, setProducts] = useState<Product[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [query, setQuery] = useState("");
  const [categoryId, setCategoryId] = useState("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const t = useI18nStore((state) => state.t);
  const productMap = useMemo(() => new Map(products.map((product) => [product.id, product])), [products]);
  const selectedBranch = useMemo(
    () => branches.find((branch) => branch.id === selectedBranchId),
    [branches, selectedBranchId]
  );

  useEffect(() => {
    Promise.all([api.myBranches(), api.products(), api.categories()])
      .then(([branchData, productData, categoryData]) => {
        const list = Array.isArray(branchData) ? branchData : [];
        setBranches(list);
        setProducts(productData);
        setCategories(Array.isArray(categoryData) ? categoryData : []);
        if (list.length > 0) {
          setSelectedBranchId(list[0].id);
        }
      })
      .catch((err) => setError(err instanceof Error ? err.message : t("inventory.loadFailed")))
      .finally(() => setLoading(false));
  }, []);

  useEffect(() => {
    if (!selectedBranchId) {
      return;
    }
    loadInventory();
  }, [selectedBranchId]);

  async function loadInventory() {
    if (!selectedBranchId) return;
    setLoading(true);
    setError("");
    try {
      const data = await api.inventories(query, selectedBranchId, categoryId || undefined);
      setInventories(Array.isArray(data) ? data : []);
    } catch (err) {
      setError(err instanceof Error ? err.message : t("inventory.loadFailed"));
    } finally {
      setLoading(false);
    }
  }

  function clearFilters() {
    setQuery("");
    setCategoryId("");
    if (selectedBranchId) {
      setLoading(true);
      api.inventories("", selectedBranchId)
      .then((data) => setInventories(Array.isArray(data) ? data : []))
        .catch((err) => setError(err instanceof Error ? err.message : t("inventory.loadFailed")))
        .finally(() => setLoading(false));
    }
  }

  return (
    <AppShell>
      <div className="mb-5 flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold">{t("branchInventory.title")}</h1>
          <p className="text-sm text-slate-500">
            {selectedBranch ? `${selectedBranch.code} · ${selectedBranch.name}` : t("inventory.description")}
          </p>
        </div>
        <select
          className="h-10 w-full rounded-md border border-line bg-white px-3 text-sm sm:w-72"
          value={selectedBranchId}
          onChange={(event) => setSelectedBranchId(event.target.value)}
        >
          {branches.map((branch) => (
            <option key={branch.id} value={branch.id}>
              {branch.code} · {branch.name}
            </option>
          ))}
        </select>
      </div>
      {error ? <div className="rounded-md border border-red-200 bg-red-50 p-3 text-sm text-red-700">{error}</div> : null}
      <section className="mb-4 rounded-md border border-line bg-white p-4 shadow-sm">
        <div className="grid gap-3 lg:grid-cols-[1fr_240px_auto_auto]">
          <Input
            placeholder={t("branchInventory.searchPlaceholder")}
            value={query}
            onChange={(event) => setQuery(event.target.value)}
            onKeyDown={(event) => event.key === "Enter" && loadInventory()}
          />
          <select className="h-10 w-full rounded-md border border-line bg-white px-3 text-sm" value={categoryId} onChange={(event) => setCategoryId(event.target.value)}>
            <option value="">{t("branchInventory.allCategories")}</option>
            {categories.map((category) => (
              <option key={category.id} value={category.id}>{category.name}</option>
            ))}
          </select>
          <Button onClick={loadInventory}>
            <Search className="h-4 w-4" />
            {t("common.search")}
          </Button>
          <Button className="!bg-white !text-slate-700 ring-1 ring-line hover:!bg-field" onClick={clearFilters}>
            <X className="h-4 w-4" />
            {t("common.clear")}
          </Button>
        </div>
      </section>
      {loading ? <ListSkeleton rows={3} /> : (
      <div className="overflow-hidden rounded-md border border-line bg-white shadow-sm">
        {inventories.length === 0 ? <div className="p-5 text-sm text-slate-500">{t("branchInventory.empty")}</div> : null}
        {inventories.map((item) => {
          const product = productMap.get(item.product_id);
          return (
            <div key={item.id} className="flex flex-wrap items-center justify-between gap-3 border-b border-line p-4 last:border-b-0">
              <div className="flex min-w-0 items-center gap-3">
                <ProductImage src={product?.image_url} name={product?.name ?? item.product_id} />
                <div>
                  <div className="font-semibold">{product?.name ?? item.product_id}</div>
                  <div className="text-xs text-slate-500">{product ? `${product.sku} · ${product.barcode}` : t("inventory.detailsUnavailable")}</div>
                </div>
              </div>
              <div className="text-right">
                <div className="text-xl font-bold">{item.quantity}</div>
                <div className="text-xs text-slate-500">{t("inventory.reserved")} {item.reserved_quantity}</div>
              </div>
            </div>
          );
        })}
      </div>
      )}
    </AppShell>
  );
}
