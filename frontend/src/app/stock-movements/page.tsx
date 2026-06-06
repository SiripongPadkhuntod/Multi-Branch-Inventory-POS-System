"use client";

import { AppShell } from "@/components/layout/app-shell";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { api } from "@/services/api";
import { useI18nStore } from "@/stores/i18n-store";
import type { Branch, InventoryMovementDetail, MovementType } from "@/types/domain";
import { ArrowDownLeft, ArrowRightLeft, ArrowUpRight, RotateCcw, Search, ShoppingCart } from "lucide-react";
import { useEffect, useMemo, useState } from "react";

const movementLabelKey: Record<MovementType, string> = {
  RECEIVE: "stockMovements.receive",
  SALE: "stockMovements.sale",
  RETURN: "stockMovements.return",
  ADJUSTMENT: "stockMovements.adjustment",
  TRANSFER_IN: "stockMovements.transferIn",
  TRANSFER_OUT: "stockMovements.transferOut"
};

function movementIcon(type: MovementType) {
  if (type === "RECEIVE" || type === "TRANSFER_IN") return ArrowDownLeft;
  if (type === "SALE" || type === "TRANSFER_OUT") return ArrowUpRight;
  if (type === "RETURN") return RotateCcw;
  return ArrowRightLeft;
}

export default function StockMovementsPage() {
  const [branches, setBranches] = useState<Branch[]>([]);
  const [movements, setMovements] = useState<InventoryMovementDetail[]>([]);
  const [branchId, setBranchId] = useState("");
  const [query, setQuery] = useState("");
  const [error, setError] = useState("");
  const t = useI18nStore((state) => state.t);

  const transferPairs = useMemo(() => {
    const pairs = new Map<string, { from?: InventoryMovementDetail; to?: InventoryMovementDetail }>();
    for (const movement of movements) {
      if (!movement.reference_id || (movement.movement_type !== "TRANSFER_IN" && movement.movement_type !== "TRANSFER_OUT")) continue;
      const current = pairs.get(movement.reference_id) ?? {};
      if (movement.movement_type === "TRANSFER_OUT") current.from = movement;
      if (movement.movement_type === "TRANSFER_IN") current.to = movement;
      pairs.set(movement.reference_id, current);
    }
    return pairs;
  }, [movements]);

  async function load(search = query, nextBranchId = branchId) {
    setError("");
    try {
      const [branchData, movementData] = await Promise.all([api.myBranches(), api.inventoryMovements(search, nextBranchId || undefined)]);
      setBranches(Array.isArray(branchData) ? branchData : []);
      setMovements(Array.isArray(movementData) ? movementData : []);
    } catch (err) {
      setError(err instanceof Error ? err.message : t("inventory.loadFailed"));
    }
  }

  useEffect(() => {
    load("", "");
  }, []);

  useEffect(() => {
    load(query, branchId);
  }, [branchId]);

  function transferText(movement: InventoryMovementDetail) {
    const pair = transferPairs.get(movement.reference_id);
    if (!pair?.from || !pair?.to) return "";
    return `${pair.from.branch_code} -> ${pair.to.branch_code}`;
  }

  return (
    <AppShell>
      <div className="mb-5 flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold">{t("stockMovements.title")}</h1>
          <p className="text-sm text-slate-500">{t("stockMovements.description")}</p>
        </div>
        <select className="h-10 w-full rounded-md border border-line bg-white px-3 text-sm sm:w-72" value={branchId} onChange={(event) => setBranchId(event.target.value)}>
          <option value="">{t("stockMovements.allBranches")}</option>
          {branches.map((branch) => <option key={branch.id} value={branch.id}>{branch.code} · {branch.name}</option>)}
        </select>
      </div>
      {error ? <div className="mb-4 rounded-md border border-red-200 bg-red-50 p-3 text-sm text-red-700">{error}</div> : null}

      <section className="mb-4 rounded-md border border-line bg-white p-4">
        <div className="flex flex-col gap-2 sm:flex-row">
          <Input placeholder={t("stockMovements.searchPlaceholder")} value={query} onChange={(event) => setQuery(event.target.value)} onKeyDown={(event) => event.key === "Enter" && load(query, branchId)} />
          <Button onClick={() => load(query, branchId)}>
            <Search className="h-4 w-4" />
            {t("common.search")}
          </Button>
        </div>
      </section>

      <section className="overflow-hidden rounded-md border border-line bg-white">
        {movements.length === 0 ? <div className="p-5 text-sm text-slate-500">{t("stockMovements.empty")}</div> : null}
        {movements.map((movement) => {
          const Icon = movementIcon(movement.movement_type);
          const transfer = transferText(movement);
          const positive = movement.quantity > 0;
          return (
            <div key={movement.id} className="grid gap-3 border-b border-line p-4 last:border-b-0 lg:grid-cols-[1fr_150px_150px_190px] lg:items-center">
              <div className="flex min-w-0 items-start gap-3">
                <div className={positive ? "rounded-md bg-emerald-50 p-2 text-emerald-700" : "rounded-md bg-red-50 p-2 text-red-700"}>
                  <Icon className="h-5 w-5" />
                </div>
                <div className="min-w-0">
                  <div className="font-semibold">{movement.product_name}</div>
                  <div className="text-xs text-slate-500">{movement.sku} · {movement.barcode}</div>
                  {transfer ? <div className="mt-1 text-xs font-semibold text-brand">{transfer}</div> : null}
                  {movement.reference_id && !transfer ? <div className="mt-1 text-xs text-slate-500">Ref {movement.reference_id}</div> : null}
                </div>
              </div>
              <div>
                <div className="text-xs text-slate-500">{t("stockMovements.type")}</div>
                <div className="flex items-center gap-2 font-semibold">
                  {movement.movement_type === "SALE" ? <ShoppingCart className="h-4 w-4 text-brand" /> : null}
                  {t(movementLabelKey[movement.movement_type])}
                </div>
              </div>
              <div>
                <div className="text-xs text-slate-500">{t("field.branch")}</div>
                <div className="font-semibold">{movement.branch_code} · {movement.branch_name}</div>
              </div>
              <div className="lg:text-right">
                <div className={positive ? "text-xl font-bold text-emerald-700" : "text-xl font-bold text-red-700"}>
                  {positive ? "+" : ""}{movement.quantity.toLocaleString()}
                </div>
                <div className="text-xs text-slate-500">{movement.created_by_name} · {new Date(movement.created_at).toLocaleString()}</div>
              </div>
            </div>
          );
        })}
      </section>
    </AppShell>
  );
}
