"use client";

import { AppShell } from "@/components/layout/app-shell";
import { Button } from "@/components/ui/button";
import { ConfirmModal } from "@/components/ui/confirm-modal";
import { Input } from "@/components/ui/input";
import { ProductImage } from "@/components/ui/product-image";
import { ListSkeleton } from "@/components/ui/skeleton";
import { api } from "@/services/api";
import { useToastStore } from "@/stores/toast-store";
import type { Category, Product } from "@/types/domain";
import { Edit3, PackagePlus, Plus, Save, Search, Trash2, X } from "lucide-react";
import { useEffect, useMemo, useState } from "react";

const money = (value: number) => `฿${(value / 100).toLocaleString("th-TH", { minimumFractionDigits: 2 })}`;

type ProductForm = {
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

const emptyForm: ProductForm = {
  id: "",
  sku: "",
  barcode: "",
  qr_code: "",
  name: "",
  description: "",
  category_id: "",
  image_url: "",
  cost_price: 0,
  sell_price: 0,
  status: "ACTIVE"
};

export default function ProductsPage() {
  const [products, setProducts] = useState<Product[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [query, setQuery] = useState("");
  const [form, setForm] = useState<ProductForm>(emptyForm);
  const [formOpen, setFormOpen] = useState(false);
  const [confirmUpdateOpen, setConfirmUpdateOpen] = useState(false);
  const [loading, setLoading] = useState(true);
  const [notice, setNotice] = useState("");
  const [error, setError] = useState("");
  const toast = useToastStore((state) => state.show);
  const categoryMap = useMemo(() => new Map(categories.map((category) => [category.id, category])), [categories]);

  async function load(search = query) {
    setError("");
    try {
      setLoading(true);
      const [productData, categoryData] = await Promise.all([api.products(search), api.categories()]);
      setProducts(Array.isArray(productData) ? productData : []);
      const list = Array.isArray(categoryData) ? categoryData : [];
      setCategories(list);
      setForm((value) => ({ ...value, category_id: value.category_id || list[0]?.id || "" }));
    } catch (err) {
      setError(err instanceof Error ? err.message : "Cannot load products");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    load("");
  }, []);

  function edit(product: Product) {
    setNotice("");
    setError("");
    setForm({
      id: product.id,
      sku: product.sku,
      barcode: product.barcode,
      qr_code: product.qr_code,
      name: product.name,
      description: product.description,
      category_id: product.category_id,
      image_url: product.image_url,
      cost_price: product.cost_price / 100,
      sell_price: product.sell_price / 100,
      status: product.status
    });
    setFormOpen(true);
  }

  function clearForm() {
    setForm({ ...emptyForm, category_id: categories[0]?.id ?? "" });
  }

  function resetForm() {
    setNotice("");
    setError("");
    clearForm();
    setFormOpen(true);
  }

  function closeForm() {
    setFormOpen(false);
    setConfirmUpdateOpen(false);
    clearForm();
  }

  function requestSave() {
    if (form.id) {
      setConfirmUpdateOpen(true);
      return;
    }
    save();
  }

  async function save() {
    setConfirmUpdateOpen(false);
    setNotice("");
    setError("");
    if (!form.sku || !form.barcode || !form.name || !form.category_id) {
      setError("SKU, barcode, name, and category are required");
      return;
    }
    const payload = {
      sku: form.sku,
      barcode: form.barcode,
      qr_code: form.qr_code,
      name: form.name,
      description: form.description,
      category_id: form.category_id,
      image_url: form.image_url,
      cost_price: Math.round(Number(form.cost_price) * 100),
      sell_price: Math.round(Number(form.sell_price) * 100),
      status: form.status
    };
    try {
      if (form.id) {
        await api.updateProduct(form.id, payload);
        setNotice("Product updated");
        toast({ type: "success", title: "Product updated", message: form.name });
      } else {
        await api.createProduct(payload);
        setNotice("Product created");
        toast({ type: "success", title: "Product created", message: form.name });
      }
      closeForm();
      await load();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Save failed");
      toast({ type: "error", title: "Save failed", message: err instanceof Error ? err.message : "Please try again" });
    }
  }

  async function remove() {
    if (!form.id) return;
    setNotice("");
    setError("");
    try {
      await api.deleteProduct(form.id);
      setNotice("Product deleted");
      toast({ type: "success", title: "Product deleted" });
      closeForm();
      await load();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Delete failed");
      toast({ type: "error", title: "Delete failed", message: err instanceof Error ? err.message : "Please try again" });
    }
  }

  return (
    <AppShell>
      <div className="mb-5 flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold">Products</h1>
          <p className="text-sm text-slate-500">Manage SKU, barcode, pricing, category, and active status.</p>
        </div>
        <Button onClick={resetForm}>
          <Plus className="h-4 w-4" />
          New Product
        </Button>
      </div>
      {notice ? <div className="mb-4 rounded-md border border-emerald-200 bg-emerald-50 p-3 text-sm text-emerald-800">{notice}</div> : null}
      {error ? <div className="mb-4 rounded-md border border-red-200 bg-red-50 p-3 text-sm text-red-700">{error}</div> : null}

      <div>
        <section className="space-y-4">
          <div className="rounded-md border border-line bg-white p-4">
            <div className="flex flex-col gap-2 sm:flex-row">
              <Input placeholder="Search by name, SKU, or barcode" value={query} onChange={(event) => setQuery(event.target.value)} onKeyDown={(event) => event.key === "Enter" && load(query)} />
              <Button onClick={() => load(query)}>
                <Search className="h-4 w-4" />
                Search
              </Button>
            </div>
          </div>
          {loading ? <ListSkeleton rows={4} /> : (
          <div className="overflow-hidden rounded-md border border-line bg-white shadow-sm">
            {products.length === 0 ? <div className="p-5 text-sm text-slate-500">No products found.</div> : null}
            {products.map((product) => (
              <div key={product.id} className="grid gap-3 border-b border-line p-4 last:border-b-0 md:grid-cols-[1fr_120px_120px_92px] md:items-center">
                <div className="flex min-w-0 items-center gap-3">
                  <ProductImage src={product.image_url} name={product.name} />
                  <div className="min-w-0">
                    <div className="truncate font-semibold">{product.name}</div>
                    <div className="text-xs text-slate-500">{product.sku} · {product.barcode}</div>
                    <div className="mt-1 text-xs text-slate-500">{categoryMap.get(product.category_id)?.name ?? "Uncategorized"} · {product.status}</div>
                  </div>
                </div>
                <div className="text-sm">
                  <div className="text-xs text-slate-500">Cost</div>
                  <div className="font-semibold">{money(product.cost_price)}</div>
                </div>
                <div className="text-sm md:text-right">
                  <div className="text-xs text-slate-500">Sell</div>
                  <div className="font-bold">{money(product.sell_price)}</div>
                </div>
                <Button className="bg-slate-800 hover:bg-slate-700" onClick={() => edit(product)}>
                  <Edit3 className="h-4 w-4" />
                  Edit
                </Button>
              </div>
            ))}
          </div>
          )}
        </section>
      </div>
      {formOpen ? (
        <div className="fixed inset-0 z-50 grid place-items-center bg-slate-950/50 px-4 py-6 backdrop-blur-sm">
          <div className="max-h-[calc(100vh-48px)] w-full max-w-3xl overflow-y-auto rounded-md border border-line bg-white p-5 shadow-2xl">
            <div className="mb-4 flex items-start justify-between gap-3">
              <div>
                <h2 className="flex items-center gap-2 text-lg font-bold">
                  <PackagePlus className="h-5 w-5 text-brand" />
                  {form.id ? "Edit Product" : "New Product"}
                </h2>
                <p className="mt-1 text-sm text-slate-500">
                  {form.id ? "Update product master data used by POS and inventory." : "Create a product with SKU, barcode, pricing, and image."}
                </p>
              </div>
              <button className="grid h-9 w-9 shrink-0 place-items-center rounded-md border border-line text-slate-600 hover:bg-field" onClick={closeForm} aria-label="Close product form">
                <X className="h-4 w-4" />
              </button>
            </div>
            <div className="grid gap-3 sm:grid-cols-2">
              <label className="block text-sm font-semibold sm:col-span-2">Name<Input className="mt-1" value={form.name} onChange={(event) => setForm({ ...form, name: event.target.value })} /></label>
              <label className="block text-sm font-semibold">SKU<Input className="mt-1" value={form.sku} onChange={(event) => setForm({ ...form, sku: event.target.value })} /></label>
              <label className="block text-sm font-semibold">Barcode<Input className="mt-1" value={form.barcode} onChange={(event) => setForm({ ...form, barcode: event.target.value })} /></label>
              <label className="block text-sm font-semibold">QR Code<Input className="mt-1" value={form.qr_code} onChange={(event) => setForm({ ...form, qr_code: event.target.value })} /></label>
              <label className="block text-sm font-semibold">Category
                <select className="mt-1 h-10 w-full rounded-md border border-line bg-white px-3 text-sm" value={form.category_id} onChange={(event) => setForm({ ...form, category_id: event.target.value })}>
                  {categories.map((category) => <option key={category.id} value={category.id}>{category.name}</option>)}
                </select>
              </label>
              <label className="block text-sm font-semibold">Cost Price<Input className="mt-1" type="number" min={0} value={form.cost_price} onChange={(event) => setForm({ ...form, cost_price: Number(event.target.value) })} /></label>
              <label className="block text-sm font-semibold">Sell Price<Input className="mt-1" type="number" min={0} value={form.sell_price} onChange={(event) => setForm({ ...form, sell_price: Number(event.target.value) })} /></label>
              <label className="block text-sm font-semibold sm:col-span-2">Image URL<Input className="mt-1" value={form.image_url} onChange={(event) => setForm({ ...form, image_url: event.target.value })} /></label>
              <label className="block text-sm font-semibold sm:col-span-2">Description<Input className="mt-1" value={form.description} onChange={(event) => setForm({ ...form, description: event.target.value })} /></label>
              <label className="block text-sm font-semibold sm:col-span-2">Status
                <select className="mt-1 h-10 w-full rounded-md border border-line bg-white px-3 text-sm" value={form.status} onChange={(event) => setForm({ ...form, status: event.target.value })}>
                  <option value="ACTIVE">Active</option>
                  <option value="INACTIVE">Inactive</option>
                </select>
              </label>
            </div>
            <div className="mt-5 flex flex-col-reverse gap-2 sm:flex-row sm:justify-between">
              <div>
                {form.id ? (
                  <Button className="w-full bg-red-600 hover:bg-red-700 sm:w-auto" onClick={remove}>
                    <Trash2 className="h-4 w-4" />
                    Delete Product
                  </Button>
                ) : null}
              </div>
              <div className="flex flex-col-reverse gap-2 sm:flex-row">
                <Button className="!bg-white !text-slate-700 ring-1 ring-line hover:!bg-field" onClick={closeForm}>Cancel</Button>
                <Button onClick={requestSave}><Save className="h-4 w-4" />Save</Button>
              </div>
            </div>
          </div>
        </div>
      ) : null}
      <ConfirmModal
        open={confirmUpdateOpen}
        title="Confirm Product Update"
        description="This will update the product master data used by POS and inventory."
        confirmLabel="Update Product"
        onCancel={() => setConfirmUpdateOpen(false)}
        onConfirm={save}
      >
        <div className="font-semibold">{form.name || "Selected product"}</div>
        <div className="text-xs text-slate-500">{form.sku} · {form.barcode}</div>
      </ConfirmModal>
    </AppShell>
  );
}
