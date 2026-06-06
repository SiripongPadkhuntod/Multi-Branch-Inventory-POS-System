"use client";

import { AppShell } from "@/components/layout/app-shell";
import { ReceiptModal } from "@/components/sales/receipt-modal";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { ProductImage } from "@/components/ui/product-image";
import { Skeleton } from "@/components/ui/skeleton";
import { api } from "@/services/api";
import { useAuthStore } from "@/stores/auth-store";
import { useCartStore } from "@/stores/cart-store";
import { useI18nStore } from "@/stores/i18n-store";
import { useToastStore } from "@/stores/toast-store";
import type { PaymentMethod, Product, SaleDetail } from "@/types/domain";
import { CreditCard, Minus, Plus, Printer, ScanBarcode, Trash2 } from "lucide-react";
import { useMemo, useState } from "react";

const money = (value: number) => `฿${(value / 100).toLocaleString("th-TH", { minimumFractionDigits: 2 })}`;

export default function POSPage() {
  const user = useAuthStore((state) => state.user);
  const { items, discount, addProduct, updateQuantity, setDiscount, clear } = useCartStore();
  const [barcode, setBarcode] = useState("");
  const [search, setSearch] = useState("");
  const [products, setProducts] = useState<Product[]>([]);
  const [searching, setSearching] = useState(false);
  const [paymentMethod, setPaymentMethod] = useState<PaymentMethod>("CASH");
  const [receiptDetail, setReceiptDetail] = useState<SaleDetail | null>(null);
  const [receiptOpen, setReceiptOpen] = useState(false);
  const [notice, setNotice] = useState("");
  const [error, setError] = useState("");
  const t = useI18nStore((state) => state.t);
  const subtotal = useMemo(() => items.reduce((sum, item) => sum + item.final_price * item.quantity, 0), [items]);
  const tax = 0;
  const total = Math.max(0, subtotal - discount + tax);
  const toast = useToastStore((state) => state.show);

  async function scan() {
    if (!barcode.trim()) return;
    setError("");
    setNotice("");
    try {
      const product = await api.productByBarcode(barcode.trim());
      addProduct(product);
      toast({ type: "success", title: t("pos.addedToCart"), message: product.name });
      setBarcode("");
    } catch (err) {
      setError(err instanceof Error ? err.message : t("pos.productNotFound"));
      toast({ type: "error", title: t("pos.productNotFound"), message: barcode.trim() });
    }
  }

  async function runSearch() {
    setError("");
    setSearching(true);
    try {
      const data = await api.products(search);
      setProducts(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : t("pos.searchFailed"));
    } finally {
      setSearching(false);
    }
  }

  async function checkout() {
    if (!user?.branch_id || items.length === 0) return;
    setError("");
    setNotice("");
    try {
      const sale = await api.createSale(user.branch_id, items, paymentMethod, discount, tax);
      clear();
      setNotice(t("pos.receiptPaid").replace("{{receipt}}", sale.receipt_number).replace("{{amount}}", money(sale.total)));
      toast({ type: "success", title: t("pos.checkoutComplete"), message: sale.receipt_number });
      try {
        const detail = await api.saleDetail(sale.id);
        setReceiptDetail(detail);
        setReceiptOpen(true);
      } catch {
        setReceiptDetail(null);
        toast({ type: "error", title: t("pos.receiptPreviewUnavailable"), message: t("pos.openLaterMySales") });
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : t("pos.checkoutFailed"));
      toast({ type: "error", title: t("pos.checkoutFailed"), message: err instanceof Error ? err.message : t("common.tryAgain") });
    }
  }

  return (
    <AppShell>
      <div className="mb-4 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold">{t("pos.title")}</h1>
          <p className="text-sm text-slate-500">{t("pos.description")}</p>
        </div>
        <Button onClick={checkout} disabled={!user?.branch_id || items.length === 0}>
          <CreditCard className="h-4 w-4" />
          {t("pos.checkout")}
        </Button>
      </div>
      {notice ? <div className="mb-4 rounded-md border border-emerald-200 bg-emerald-50 p-3 text-sm font-medium text-emerald-800">{notice}</div> : null}
      {receiptDetail ? (
        <div className="mb-4 flex flex-col gap-2 rounded-md border border-emerald-200 bg-emerald-50 p-3 text-sm text-emerald-900 sm:flex-row sm:items-center sm:justify-between">
          <span className="font-medium">{t("pos.receiptReady").replace("{{receipt}}", receiptDetail.receipt_number)}</span>
          <Button className="bg-slate-800 hover:bg-slate-700" onClick={() => setReceiptOpen(true)}>
            <Printer className="h-4 w-4" />
            {t("pos.viewReceipt")}
          </Button>
        </div>
      ) : null}
      {error ? <div className="mb-4 rounded-md border border-red-200 bg-red-50 p-3 text-sm font-medium text-red-700">{error}</div> : null}
      <div className="grid gap-4 xl:grid-cols-[1fr_420px]">
        <section className="space-y-4">
          <div className="rounded-md border border-line bg-white p-4">
            <label className="text-sm font-semibold">{t("field.barcode")}</label>
            <div className="mt-2 flex gap-2">
              <Input value={barcode} onChange={(event) => setBarcode(event.target.value)} onKeyDown={(event) => event.key === "Enter" && scan()} autoFocus />
              <Button onClick={scan}>
                <ScanBarcode className="h-4 w-4" />
                {t("pos.add")}
              </Button>
            </div>
          </div>
          <div className="rounded-md border border-line bg-white p-4">
            <label className="text-sm font-semibold">{t("pos.productSearch")}</label>
            <div className="mt-2 flex gap-2">
              <Input value={search} onChange={(event) => setSearch(event.target.value)} onKeyDown={(event) => event.key === "Enter" && runSearch()} />
              <Button onClick={runSearch}>{t("common.search")}</Button>
            </div>
            <div className="mt-4 grid gap-2 md:grid-cols-2">
              {searching ? Array.from({ length: 4 }).map((_, index) => (
                <div key={index} className="rounded-md border border-line p-3">
                  <div className="flex gap-3">
                    <Skeleton className="h-14 w-14" />
                    <div className="flex-1 space-y-2">
                      <Skeleton className="h-4 w-3/4" />
                      <Skeleton className="h-3 w-1/2" />
                    </div>
                  </div>
                </div>
              )) : null}
              {products.map((product) => (
                <button key={product.id} onClick={() => { addProduct(product); toast({ type: "success", title: t("pos.addedToCart"), message: product.name }); }} className="rounded-md border border-line bg-white p-3 text-left shadow-sm transition hover:border-brand hover:bg-brandSoft/40">
                  <div className="flex gap-3">
                    <ProductImage src={product.image_url} name={product.name} />
                    <div className="min-w-0">
                      <div className="truncate font-semibold">{product.name}</div>
                      <div className="text-xs text-slate-500">{product.sku} · {product.barcode}</div>
                      <div className="mt-2 text-sm font-bold">{money(product.sell_price)}</div>
                    </div>
                  </div>
                </button>
              ))}
            </div>
          </div>
        </section>
        <aside className="rounded-md border border-line bg-white p-4">
          <h2 className="text-lg font-bold">{t("pos.cart")}</h2>
          <div className="mt-4 space-y-3">
            {items.map((item) => (
              <div key={item.product.id} className="rounded-md border border-line p-3">
                <div className="flex items-start justify-between gap-3">
                  <div className="flex min-w-0 gap-3">
                    <ProductImage src={item.product.image_url} name={item.product.name} className="h-12 w-12" />
                    <div className="min-w-0">
                    <div className="font-semibold">{item.product.name}</div>
                    <div className="text-xs text-slate-500">{money(item.final_price)} {t("pos.each")}</div>
                    </div>
                  </div>
                  <button onClick={() => updateQuantity(item.product.id, 0)} title={t("pos.remove")}>
                    <Trash2 className="h-4 w-4 text-red-600" />
                  </button>
                </div>
                <div className="mt-3 flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <button className="rounded-md border border-line p-2" onClick={() => updateQuantity(item.product.id, item.quantity - 1)} title={t("pos.decrease")}>
                      <Minus className="h-4 w-4" />
                    </button>
                    <span className="w-8 text-center font-semibold">{item.quantity}</span>
                    <button className="rounded-md border border-line p-2" onClick={() => updateQuantity(item.product.id, item.quantity + 1)} title={t("pos.increase")}>
                      <Plus className="h-4 w-4" />
                    </button>
                  </div>
                  <div className="font-bold">{money(item.final_price * item.quantity)}</div>
                </div>
              </div>
            ))}
          </div>
          <div className="mt-5 space-y-3 border-t border-line pt-4">
            <label className="block text-sm font-semibold">
              {t("field.discount")}
              <Input className="mt-1" type="number" value={discount / 100} onChange={(event) => setDiscount(Math.round(Number(event.target.value) * 100))} />
            </label>
            <select className="h-10 w-full rounded-md border border-line px-3" value={paymentMethod} onChange={(event) => setPaymentMethod(event.target.value as PaymentMethod)}>
              <option value="CASH">{t("payment.cash")}</option>
              <option value="PROMPTPAY">{t("payment.promptpay")}</option>
              <option value="BANK_TRANSFER">{t("payment.bankTransfer")}</option>
              <option value="CREDIT_CARD">{t("payment.creditCard")}</option>
            </select>
            <div className="flex justify-between text-sm"><span>{t("field.subtotal")}</span><span>{money(subtotal)}</span></div>
            <div className="flex justify-between text-sm"><span>{t("field.discount")}</span><span>{money(discount)}</span></div>
            <div className="flex justify-between text-xl font-bold"><span>{t("field.total")}</span><span>{money(total)}</span></div>
          </div>
        </aside>
      </div>
      <ReceiptModal open={receiptOpen} detail={receiptDetail} onClose={() => setReceiptOpen(false)} />
    </AppShell>
  );
}
