"use client";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useI18nStore } from "@/stores/i18n-store";
import type { SaleDetail } from "@/types/domain";
import { RotateCcw, X } from "lucide-react";
import { useEffect, useMemo, useState } from "react";

const money = (value: number) => `฿${(value / 100).toLocaleString("th-TH", { minimumFractionDigits: 2 })}`;

type RefundModalProps = {
  open: boolean;
  detail: SaleDetail | null;
  loading?: boolean;
  onClose: () => void;
  onConfirm: (items: { product_id: string; quantity: number }[]) => void;
};

export function RefundModal({ open, detail, loading = false, onClose, onConfirm }: RefundModalProps) {
  const [quantities, setQuantities] = useState<Record<string, number>>({});
  const t = useI18nStore((state) => state.t);

  useEffect(() => {
    if (!open || !detail) {
      setQuantities({});
    }
  }, [open, detail]);

  const selectedItems = useMemo(() => {
    if (!detail) {
      return [];
    }
    return detail.items
      .map((item) => ({
        product_id: item.product_id,
        quantity: Math.max(0, Math.min(refundableQuantity(item), Math.floor(quantities[item.id] ?? 0)))
      }))
      .filter((item) => item.quantity > 0);
  }, [detail, quantities]);

  const refundTotal = useMemo(() => {
    if (!detail) {
      return 0;
    }
    return detail.items.reduce((sum, item) => {
      const quantity = Math.max(0, Math.min(refundableQuantity(item), Math.floor(quantities[item.id] ?? 0)));
      return sum + quantity * item.final_price;
    }, 0);
  }, [detail, quantities]);

  if (!open || !detail) {
    return null;
  }

  return (
    <div className="fixed inset-0 z-50 grid place-items-center bg-slate-950/50 px-4 py-6 backdrop-blur-sm">
      <div className="max-h-[calc(100vh-48px)] w-full max-w-2xl overflow-y-auto rounded-md border border-line bg-white p-5 shadow-2xl">
        <div className="mb-4 flex items-start justify-between gap-3">
          <div>
            <h2 className="flex items-center gap-2 text-lg font-bold">
              <RotateCcw className="h-5 w-5 text-brand" />
              Refund / Return
            </h2>
            <p className="mt-1 text-sm text-slate-500">{detail.receipt_number} · {t("refund.selectedItems")}</p>
          </div>
          <button className="grid h-9 w-9 shrink-0 place-items-center rounded-md border border-line text-slate-600 hover:bg-field" onClick={onClose} aria-label={t("common.close")}>
            <X className="h-4 w-4" />
          </button>
        </div>

        <div className="overflow-hidden rounded-md border border-line">
          {detail.items.map((item) => {
            const quantity = quantities[item.id] ?? 0;
            const remaining = refundableQuantity(item);
            return (
              <div key={item.id} className="grid gap-3 border-b border-line p-3 last:border-b-0 md:grid-cols-[1fr_120px_120px] md:items-center">
                <div className="min-w-0">
                  <div className="truncate font-semibold">{item.product_name}</div>
                  <div className="text-xs text-slate-500">{item.sku} · sold {item.quantity} · returned {item.returned_quantity} · left {remaining}</div>
                </div>
                <label className="block text-sm font-semibold">
                  {t("sales.qty")}
                  <Input
                    className="mt-1"
                    type="number"
                    min={0}
                    max={remaining}
                    value={quantity}
                    disabled={remaining === 0}
                    onChange={(event) => setQuantities({ ...quantities, [item.id]: Number(event.target.value) })}
                  />
                </label>
                <div className="text-sm md:text-right">
                  <div className="text-xs text-slate-500">{t("refund.estimate")}</div>
                  <div className="font-bold">{money(Math.max(0, Math.min(remaining, quantity)) * item.final_price)}</div>
                </div>
              </div>
            );
          })}
        </div>

        <div className="mt-4 rounded-md bg-field p-3 text-sm">
          <div className="flex justify-between">
            <span>{t("refund.selectedItems")}</span>
            <span className="font-semibold">{selectedItems.reduce((sum, item) => sum + item.quantity, 0)}</span>
          </div>
          <div className="mt-1 flex justify-between">
            <span>{t("refund.estimatedRefund")}</span>
            <span className="font-bold">{money(refundTotal)}</span>
          </div>
        </div>

        <div className="mt-5 flex flex-col-reverse gap-2 sm:flex-row sm:justify-end">
          <Button className="!bg-white !text-slate-700 ring-1 ring-line hover:!bg-field" onClick={onClose}>{t("common.cancel")}</Button>
          <Button disabled={loading || selectedItems.length === 0} onClick={() => onConfirm(selectedItems)}>
            <RotateCcw className="h-4 w-4" />
            {loading ? t("refund.refunding") : t("refund.confirm")}
          </Button>
        </div>
      </div>
    </div>
  );
}

function refundableQuantity(item: SaleDetail["items"][number]) {
  return Math.max(0, item.quantity - (item.returned_quantity ?? 0));
}
