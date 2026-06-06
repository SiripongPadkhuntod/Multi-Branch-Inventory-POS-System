import type { Sale, SaleDetail } from "@/types/domain";
import { useI18nStore } from "@/stores/i18n-store";

type SaleStatusProps = {
  sale: Pick<Sale, "payment_status" | "refund_status">;
};

export function SaleStatus({ sale }: SaleStatusProps) {
  const t = useI18nStore((state) => state.t);
  if (sale.refund_status === "REFUNDED") {
    return <span className="text-xs font-semibold text-red-600">{t("status.refunded")}</span>;
  }
  if (sale.refund_status === "PARTIAL_REFUND") {
    return <span className="text-xs font-semibold text-amber-700">{t("status.partialRefund")}</span>;
  }
  return <span className="text-xs font-semibold text-emerald-700">{sale.payment_status}</span>;
}

export function canRefundSale(detail: SaleDetail | null) {
  if (!detail || detail.refund_status === "REFUNDED") {
    return false;
  }
  return detail.items.some((item) => item.quantity - (item.returned_quantity ?? 0) > 0);
}
