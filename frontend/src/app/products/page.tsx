"use client";

import { AppShell } from "@/components/layout/app-shell";
import { Button } from "@/components/ui/button";
import { ConfirmModal } from "@/components/ui/confirm-modal";
import { Input } from "@/components/ui/input";
import { ProductImage } from "@/components/ui/product-image";
import { ListSkeleton } from "@/components/ui/skeleton";
import { api } from "@/services/api";
import { useI18nStore } from "@/stores/i18n-store";
import { useToastStore } from "@/stores/toast-store";
import type { Category, Product } from "@/types/domain";
import { Edit3, ImagePlus, PackagePlus, Plus, Save, Search, Trash2, Upload, X } from "lucide-react";
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
  const [uploadingImage, setUploadingImage] = useState(false);
  const [notice, setNotice] = useState("");
  const [error, setError] = useState("");
  const t = useI18nStore((state) => state.t);
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
      setError(err instanceof Error ? err.message : t("products.loadFailed"));
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
      setError(t("products.required"));
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
        setNotice(t("products.updated"));
        toast({ type: "success", title: t("products.updated"), message: form.name });
      } else {
        await api.createProduct(payload);
        setNotice(t("products.created"));
        toast({ type: "success", title: t("products.created"), message: form.name });
      }
      closeForm();
      await load();
    } catch (err) {
      setError(err instanceof Error ? err.message : t("products.saveFailed"));
      toast({ type: "error", title: t("products.saveFailed"), message: err instanceof Error ? err.message : t("common.tryAgain") });
    }
  }

  async function uploadImage(file?: File) {
    if (!file) {
      return;
    }
    setError("");
    if (!file.type.startsWith("image/")) {
      setError(t("products.imageRequired"));
      return;
    }
    if (file.size > 5 * 1024 * 1024) {
      setError(t("products.imageTooLarge"));
      return;
    }
    try {
      setUploadingImage(true);
      const result = await api.uploadProductImage(file);
      setForm((value) => ({ ...value, image_url: result.image_url }));
      toast({ type: "success", title: t("products.imageUploaded"), message: file.name });
    } catch (err) {
      setError(err instanceof Error ? err.message : t("products.imageUploadFailed"));
      toast({ type: "error", title: t("products.imageUploadFailed"), message: err instanceof Error ? err.message : t("common.tryAgain") });
    } finally {
      setUploadingImage(false);
    }
  }

  async function remove() {
    if (!form.id) return;
    setNotice("");
    setError("");
    try {
      await api.deleteProduct(form.id);
      setNotice(t("products.deleted"));
      toast({ type: "success", title: t("products.deleted") });
      closeForm();
      await load();
    } catch (err) {
      setError(err instanceof Error ? err.message : t("products.deleteFailed"));
      toast({ type: "error", title: t("products.deleteFailed"), message: err instanceof Error ? err.message : t("common.tryAgain") });
    }
  }

  return (
    <AppShell>
      <div className="mb-5 flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold">{t("products.title")}</h1>
          <p className="text-sm text-slate-500">{t("products.description")}</p>
        </div>
        <Button onClick={resetForm}>
          <Plus className="h-4 w-4" />
          {t("products.new")}
        </Button>
      </div>
      {notice ? <div className="mb-4 rounded-md border border-emerald-200 bg-emerald-50 p-3 text-sm text-emerald-800">{notice}</div> : null}
      {error ? <div className="mb-4 rounded-md border border-red-200 bg-red-50 p-3 text-sm text-red-700">{error}</div> : null}

      <div>
        <section className="space-y-4">
          <div className="rounded-md border border-line bg-white p-4">
            <div className="flex flex-col gap-2 sm:flex-row">
              <Input placeholder={t("products.searchPlaceholder")} value={query} onChange={(event) => setQuery(event.target.value)} onKeyDown={(event) => event.key === "Enter" && load(query)} />
              <Button onClick={() => load(query)}>
                <Search className="h-4 w-4" />
                {t("common.search")}
              </Button>
            </div>
          </div>
          {loading ? <ListSkeleton rows={4} /> : (
          <div className="overflow-hidden rounded-md border border-line bg-white shadow-sm">
            {products.length === 0 ? <div className="p-5 text-sm text-slate-500">{t("products.empty")}</div> : null}
            {products.map((product) => (
              <div key={product.id} className="grid gap-3 border-b border-line p-4 last:border-b-0 md:grid-cols-[1fr_120px_120px_92px] md:items-center">
                <div className="flex min-w-0 items-center gap-3">
                  <ProductImage src={product.image_url} name={product.name} />
                  <div className="min-w-0">
                    <div className="truncate font-semibold">{product.name}</div>
                    <div className="text-xs text-slate-500">{product.sku} · {product.barcode}</div>
                    <div className="mt-1 text-xs text-slate-500">{categoryMap.get(product.category_id)?.name ?? t("common.uncategorized")} · {product.status}</div>
                  </div>
                </div>
                <div className="text-sm">
                  <div className="text-xs text-slate-500">{t("field.cost")}</div>
                  <div className="font-semibold">{money(product.cost_price)}</div>
                </div>
                <div className="text-sm md:text-right">
                  <div className="text-xs text-slate-500">{t("field.sell")}</div>
                  <div className="font-bold">{money(product.sell_price)}</div>
                </div>
                <Button className="bg-slate-800 hover:bg-slate-700" onClick={() => edit(product)}>
                  <Edit3 className="h-4 w-4" />
                  {t("common.edit")}
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
                  {form.id ? t("products.edit") : t("products.new")}
                </h2>
                <p className="mt-1 text-sm text-slate-500">
                  {form.id ? t("products.editDescription") : t("products.createDescription")}
                </p>
              </div>
              <button className="grid h-9 w-9 shrink-0 place-items-center rounded-md border border-line text-slate-600 hover:bg-field" onClick={closeForm} aria-label={t("common.close")}>
                <X className="h-4 w-4" />
              </button>
            </div>
            <div className="grid gap-3 sm:grid-cols-2">
              <label className="block text-sm font-semibold sm:col-span-2">{t("field.name")}<Input className="mt-1" value={form.name} onChange={(event) => setForm({ ...form, name: event.target.value })} /></label>
              <label className="block text-sm font-semibold">{t("field.sku")}<Input className="mt-1" value={form.sku} onChange={(event) => setForm({ ...form, sku: event.target.value })} /></label>
              <label className="block text-sm font-semibold">{t("field.barcode")}<Input className="mt-1" value={form.barcode} onChange={(event) => setForm({ ...form, barcode: event.target.value })} /></label>
              <label className="block text-sm font-semibold">{t("field.qrCode")}<Input className="mt-1" value={form.qr_code} onChange={(event) => setForm({ ...form, qr_code: event.target.value })} /></label>
              <label className="block text-sm font-semibold">{t("field.category")}
                <select className="mt-1 h-10 w-full rounded-md border border-line bg-white px-3 text-sm" value={form.category_id} onChange={(event) => setForm({ ...form, category_id: event.target.value })}>
                  {categories.map((category) => <option key={category.id} value={category.id}>{category.name}</option>)}
                </select>
              </label>
              <label className="block text-sm font-semibold">{t("field.costPrice")}<Input className="mt-1" type="number" min={0} value={form.cost_price} onChange={(event) => setForm({ ...form, cost_price: Number(event.target.value) })} /></label>
              <label className="block text-sm font-semibold">{t("field.sellPrice")}<Input className="mt-1" type="number" min={0} value={form.sell_price} onChange={(event) => setForm({ ...form, sell_price: Number(event.target.value) })} /></label>
              <div className="grid gap-3 rounded-md border border-line bg-field p-3 sm:col-span-2 sm:grid-cols-[96px_1fr]">
                <ProductImage src={form.image_url} name={form.name || t("products.image")} className="h-24 w-24" />
                <div className="space-y-3">
                  <div>
                    <div className="flex items-center gap-2 text-sm font-semibold">
                      <ImagePlus className="h-4 w-4 text-brand" />
                      {t("products.image")}
                    </div>
                    <p className="mt-1 text-xs text-slate-500">{t("products.imageHelp")}</p>
                  </div>
                  <div className="flex flex-col gap-2 sm:flex-row">
                    <label className="inline-flex h-10 cursor-pointer items-center justify-center gap-2 rounded-md bg-slate-800 px-4 text-sm font-semibold text-white hover:bg-slate-700">
                      <Upload className="h-4 w-4" />
                      {uploadingImage ? t("products.uploading") : t("products.upload")}
                      <input
                        className="sr-only"
                        type="file"
                        accept="image/png,image/jpeg,image/webp,image/gif"
                        disabled={uploadingImage}
                        onChange={(event) => uploadImage(event.target.files?.[0])}
                      />
                    </label>
                    {form.image_url ? (
                      <Button className="!bg-white !text-slate-700 ring-1 ring-line hover:!bg-white" onClick={() => setForm({ ...form, image_url: "" })}>
                        {t("products.removeImage")}
                      </Button>
                    ) : null}
                  </div>
                  <label className="block text-xs font-semibold text-slate-500">
                    {t("field.imagePathUrl")}
                    <Input className="mt-1" value={form.image_url} onChange={(event) => setForm({ ...form, image_url: event.target.value })} />
                  </label>
                </div>
              </div>
              <label className="block text-sm font-semibold sm:col-span-2">{t("field.description")}<Input className="mt-1" value={form.description} onChange={(event) => setForm({ ...form, description: event.target.value })} /></label>
              <label className="block text-sm font-semibold sm:col-span-2">{t("field.status")}
                <select className="mt-1 h-10 w-full rounded-md border border-line bg-white px-3 text-sm" value={form.status} onChange={(event) => setForm({ ...form, status: event.target.value })}>
                  <option value="ACTIVE">{t("status.active")}</option>
                  <option value="INACTIVE">{t("status.inactive")}</option>
                </select>
              </label>
            </div>
            <div className="mt-5 flex flex-col-reverse gap-2 sm:flex-row sm:justify-between">
              <div>
                {form.id ? (
                  <Button className="w-full bg-red-600 hover:bg-red-700 sm:w-auto" onClick={remove}>
                    <Trash2 className="h-4 w-4" />
                    {t("products.delete")}
                  </Button>
                ) : null}
              </div>
              <div className="flex flex-col-reverse gap-2 sm:flex-row">
                <Button className="!bg-white !text-slate-700 ring-1 ring-line hover:!bg-field" onClick={closeForm}>{t("common.cancel")}</Button>
                <Button onClick={requestSave}><Save className="h-4 w-4" />{t("common.save")}</Button>
              </div>
            </div>
          </div>
        </div>
      ) : null}
      <ConfirmModal
        open={confirmUpdateOpen}
        title={t("products.confirmTitle")}
        description={t("products.confirmDescription")}
        confirmLabel={t("products.confirmButton")}
        onCancel={() => setConfirmUpdateOpen(false)}
        onConfirm={save}
      >
        <div className="font-semibold">{form.name || t("products.selected")}</div>
        <div className="text-xs text-slate-500">{form.sku} · {form.barcode}</div>
      </ConfirmModal>
    </AppShell>
  );
}
