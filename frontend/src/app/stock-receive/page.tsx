"use client";

import { AppShell } from "@/components/layout/app-shell";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { ProductImage } from "@/components/ui/product-image";
import { api } from "@/services/api";
import { useToastStore } from "@/stores/toast-store";
import type { Branch, Product } from "@/types/domain";
import { ClipboardPlus, Search } from "lucide-react";
import { useEffect, useState } from "react";

export default function StockReceivePage() {
  const [branches, setBranches] = useState<Branch[]>([]);
  const [selectedBranchId, setSelectedBranchId] = useState("");
  const [products, setProducts] = useState<Product[]>([]);
  const [selectedProductId, setSelectedProductId] = useState("");
  const [query, setQuery] = useState("");
  const [quantity, setQuantity] = useState(1);
  const [reason, setReason] = useState("Manager stock receive");
  const [notice, setNotice] = useState("");
  const [error, setError] = useState("");
  const toast = useToastStore((state) => state.show);

  useEffect(() => {
    Promise.all([api.myBranches(), api.products()])
      .then(([branchData, productData]) => {
        const list = Array.isArray(branchData) ? branchData : [];
        setBranches(list);
        if (list.length > 0) {
          setSelectedBranchId(list[0].id);
        }
        setProducts(productData);
      })
      .catch((err) => setError(err instanceof Error ? err.message : "Cannot load products"));
  }, []);

  async function searchProducts() {
    setError("");
    try {
      setProducts(await api.products(query));
    } catch (err) {
      setError(err instanceof Error ? err.message : "Search failed");
    }
  }

  async function receive() {
    setError("");
    setNotice("");
    if (!selectedProductId) {
      setError("Select a product first");
      return;
    }
    if (!selectedBranchId) {
      setError("Select a branch first");
      return;
    }
    if (quantity < 1) {
      setError("Quantity must be at least 1");
      return;
    }
    try {
      await api.receiveStock(selectedProductId, quantity, reason, selectedBranchId);
      const product = products.find((item) => item.id === selectedProductId);
      setNotice(`Received ${quantity} item(s) into stock${product ? `: ${product.name}` : ""}`);
      toast({ type: "success", title: "Stock received", message: product ? product.name : `${quantity} item(s)` });
      setQuantity(1);
      setSelectedProductId("");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Receive stock failed");
      toast({ type: "error", title: "Receive stock failed", message: err instanceof Error ? err.message : "Please try again" });
    }
  }

  return (
    <AppShell>
      <div className="mb-5">
        <h1 className="text-2xl font-bold">Receive Stock</h1>
        <p className="text-sm text-slate-500">Add incoming stock to your assigned branch.</p>
      </div>
      {notice ? <div className="mb-4 rounded-md border border-emerald-200 bg-emerald-50 p-3 text-sm font-medium text-emerald-800">{notice}</div> : null}
      {error ? <div className="mb-4 rounded-md border border-red-200 bg-red-50 p-3 text-sm font-medium text-red-700">{error}</div> : null}
      <div className="grid gap-4 lg:grid-cols-[1fr_360px]">
        <section className="rounded-md border border-line bg-white p-4">
          <label className="text-sm font-semibold">Search Product</label>
          <div className="mt-2 flex gap-2">
            <Input value={query} onChange={(event) => setQuery(event.target.value)} onKeyDown={(event) => event.key === "Enter" && searchProducts()} />
            <Button onClick={searchProducts}>
              <Search className="h-4 w-4" />
              Search
            </Button>
          </div>
          <div className="mt-4 grid gap-2 md:grid-cols-2">
            {products.map((product) => (
              <button
                key={product.id}
                onClick={() => setSelectedProductId(product.id)}
                className={`rounded-md border p-3 text-left ${selectedProductId === product.id ? "border-brand bg-emerald-50" : "border-line hover:border-brand"}`}
              >
                <div className="flex gap-3">
                  <ProductImage src={product.image_url} name={product.name} />
                  <div className="min-w-0">
                    <div className="truncate font-semibold">{product.name}</div>
                    <div className="text-xs text-slate-500">{product.sku} · {product.barcode}</div>
                  </div>
                </div>
              </button>
            ))}
          </div>
        </section>
        <aside className="rounded-md border border-line bg-white p-4">
          <h2 className="flex items-center gap-2 text-lg font-bold">
            <ClipboardPlus className="h-5 w-5 text-brand" />
            Receive
          </h2>
          <div className="mt-4 space-y-4">
            <label className="block text-sm font-semibold">
              Branch
              <select
                className="mt-1 h-10 w-full rounded-md border border-line bg-white px-3 text-sm"
                value={selectedBranchId}
                onChange={(event) => setSelectedBranchId(event.target.value)}
              >
                {branches.map((branch) => (
                  <option key={branch.id} value={branch.id}>
                    {branch.code} · {branch.name}
                  </option>
                ))}
              </select>
            </label>
            <label className="block text-sm font-semibold">
              Quantity
              <Input className="mt-1" type="number" min={1} value={quantity} onChange={(event) => setQuantity(Number(event.target.value))} />
            </label>
            <label className="block text-sm font-semibold">
              Reason
              <Input className="mt-1" value={reason} onChange={(event) => setReason(event.target.value)} />
            </label>
            <div className="rounded-md bg-field p-3 text-xs text-slate-600">
              Stock will be added only to branches you are allowed to manage.
            </div>
            <Button className="w-full" onClick={receive} disabled={!selectedProductId || !selectedBranchId || quantity < 1}>
              Receive Stock
            </Button>
          </div>
        </aside>
      </div>
    </AppShell>
  );
}
