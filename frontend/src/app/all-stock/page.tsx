"use client";

import { AppShell } from "@/components/layout/app-shell";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { ProductImage } from "@/components/ui/product-image";
import { ListSkeleton } from "@/components/ui/skeleton";
import { api } from "@/services/api";
import { useI18nStore } from "@/stores/i18n-store";
import type { ProductStockSummary } from "@/types/domain";
import { ChevronDown, ChevronUp, Search } from "lucide-react";
import { useEffect, useMemo, useState } from "react";

export default function AllStockPage() {
  const [items, setItems] = useState<ProductStockSummary[]>([]);
  const [query, setQuery] = useState("");
  const [expandedProductId, setExpandedProductId] = useState("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const t = useI18nStore((state) => state.t);

  const totalQuantity = useMemo(() => items.reduce((sum, item) => sum + item.total_quantity, 0), [items]);
  const totalReserved = useMemo(() => items.reduce((sum, item) => sum + item.total_reserved, 0), [items]);

  async function load(search = query) {
    setError("");
    try {
      setLoading(true);
      const data = await api.allStock(search);
      setItems(Array.isArray(data) ? data : []);
    } catch (err) {
      setError(err instanceof Error ? err.message : t("inventory.loadFailed"));
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    load("");
  }, []);

  return (
    <AppShell>
      <div className="mb-5">
        <h1 className="text-2xl font-bold">{t("allStock.title")}</h1>
        <p className="text-sm text-slate-500">{t("allStock.description")}</p>
      </div>
      {error ? <div className="mb-4 rounded-md border border-red-200 bg-red-50 p-3 text-sm text-red-700">{error}</div> : null}

      <section className="mb-4 grid gap-3 sm:grid-cols-3">
        <article className="rounded-md border border-line bg-white p-4">
          <div className="text-sm text-slate-500">{t("allStock.productCount")}</div>
          <div className="mt-2 text-2xl font-bold">{items.length.toLocaleString()}</div>
        </article>
        <article className="rounded-md border border-line bg-white p-4">
          <div className="text-sm text-slate-500">{t("allStock.totalOnHand")}</div>
          <div className="mt-2 text-2xl font-bold">{totalQuantity.toLocaleString()}</div>
        </article>
        <article className="rounded-md border border-line bg-white p-4">
          <div className="text-sm text-slate-500">{t("allStock.totalReserved")}</div>
          <div className="mt-2 text-2xl font-bold">{totalReserved.toLocaleString()}</div>
        </article>
      </section>

      <section className="mb-4 rounded-md border border-line bg-white p-4">
        <div className="flex flex-col gap-2 sm:flex-row">
          <Input placeholder={t("allStock.searchPlaceholder")} value={query} onChange={(event) => setQuery(event.target.value)} onKeyDown={(event) => event.key === "Enter" && load(query)} />
          <Button onClick={() => load(query)}>
            <Search className="h-4 w-4" />
            {t("common.search")}
          </Button>
        </div>
      </section>

      {loading ? <ListSkeleton rows={5} /> : (
      <section className="overflow-hidden rounded-md border border-line bg-white shadow-sm">
        {items.length === 0 ? <div className="p-5 text-sm text-slate-500">{t("allStock.empty")}</div> : null}
        {items.map((item) => {
          const expanded = expandedProductId === item.product_id;
          return (
            <div key={item.product_id} className="border-b border-line last:border-b-0">
              <button className="grid w-full gap-3 p-4 text-left hover:bg-field md:grid-cols-[1fr_130px_130px_32px] md:items-center" onClick={() => setExpandedProductId(expanded ? "" : item.product_id)}>
                <div className="flex min-w-0 items-center gap-3">
                  <ProductImage src={item.image_url} name={item.product_name} />
                  <div className="min-w-0">
                    <div className="truncate font-semibold">{item.product_name}</div>
                    <div className="text-xs text-slate-500">{item.sku} · {item.barcode}</div>
                  </div>
                </div>
                <div>
                  <div className="text-xs text-slate-500">{t("inventory.onHand")}</div>
                  <div className="text-xl font-bold">{item.total_quantity.toLocaleString()}</div>
                </div>
                <div>
                  <div className="text-xs text-slate-500">{t("inventory.reserved")}</div>
                  <div className="font-semibold">{item.total_reserved.toLocaleString()}</div>
                </div>
                {expanded ? <ChevronUp className="h-4 w-4" /> : <ChevronDown className="h-4 w-4" />}
              </button>
              {expanded ? (
                <div className="border-t border-line bg-slate-50 p-4">
                  <div className="grid gap-2 md:grid-cols-2 xl:grid-cols-3">
                    {item.branches.map((branch) => (
                      <div key={branch.branch_id} className="rounded-md border border-line bg-white p-3 text-sm">
                        <div className="font-semibold">{branch.branch_code} · {branch.branch_name}</div>
                        <div className="mt-2 flex justify-between"><span>{t("inventory.onHand")}</span><span className="font-bold">{branch.quantity.toLocaleString()}</span></div>
                        <div className="flex justify-between text-slate-500"><span>{t("inventory.reserved")}</span><span>{branch.reserved_quantity.toLocaleString()}</span></div>
                        <div className="flex justify-between text-slate-500"><span>{t("inventory.reorderAt")}</span><span>{branch.reorder_threshold.toLocaleString()}</span></div>
                      </div>
                    ))}
                  </div>
                </div>
              ) : null}
            </div>
          );
        })}
      </section>
      )}
    </AppShell>
  );
}
