"use client";

import { AppShell } from "@/components/layout/app-shell";
import { ListSkeleton } from "@/components/ui/skeleton";
import { api } from "@/services/api";
import { useI18nStore } from "@/stores/i18n-store";
import type { Branch, DashboardSummary, Sale } from "@/types/domain";
import { BarChart3, PackageSearch, ReceiptText, TrendingUp } from "lucide-react";
import { useEffect, useMemo, useState } from "react";

const money = (value: number) => `฿${(value / 100).toLocaleString("th-TH", { minimumFractionDigits: 2 })}`;
const compactMoney = (value: number) =>
  `฿${(value / 100).toLocaleString("th-TH", { notation: "compact", maximumFractionDigits: 1 })}`;

const emptySummary: DashboardSummary = {
  daily_sales: 0,
  monthly_sales: 0,
  revenue: 0,
  profit: 0,
  low_stock: [],
  top_products: [],
  branch_comparison: []
};

export default function DashboardPage() {
  const language = useI18nStore((state) => state.language);
  const t = useI18nStore((state) => state.t);
  const [branches, setBranches] = useState<Branch[]>([]);
  const [selectedBranchId, setSelectedBranchId] = useState("");
  const [summary, setSummary] = useState<DashboardSummary>(emptySummary);
  const [sales, setSales] = useState<Sale[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    api.myBranches()
      .then((data) => {
        setBranches(Array.isArray(data) ? data : []);
        if (data.length === 1) {
          setSelectedBranchId(data[0].id);
        }
      })
      .catch((err) => setError(err instanceof Error ? err.message : "Cannot load branches"));
  }, []);

  useEffect(() => {
    setLoading(true);
    Promise.all([api.dashboardSummary(selectedBranchId || undefined), api.branchSales(selectedBranchId || undefined)])
      .then(([summaryData, salesData]) => {
        setSummary(summaryData ?? emptySummary);
        setSales(Array.isArray(salesData) ? salesData : []);
      })
      .catch((err) => setError(err instanceof Error ? err.message : "Cannot load dashboard"))
      .finally(() => setLoading(false));
  }, [selectedBranchId]);

  const selectedBranch = useMemo(
    () => branches.find((branch) => branch.id === selectedBranchId),
    [branches, selectedBranchId]
  );

  const widgets = [
    { label: t("dashboard.dailySales"), value: summary.daily_sales.toLocaleString(), icon: ReceiptText },
    { label: t("dashboard.monthlySales"), value: summary.monthly_sales.toLocaleString(), icon: BarChart3 },
    { label: t("dashboard.revenue"), value: money(summary.revenue), icon: TrendingUp },
    { label: t("dashboard.profit"), value: money(summary.profit), icon: PackageSearch }
  ];

  const salesTrend = useMemo(() => buildSalesTrend(sales, language), [sales, language]);

  return (
    <AppShell>
      <div className="mb-5 flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold">{t("dashboard.title")}</h1>
          <p className="text-sm text-slate-500">
            {selectedBranch ? `${selectedBranch.code} · ${selectedBranch.name}` : t("dashboard.description")}
          </p>
        </div>
        <select
          className="h-10 w-full rounded-md border border-line bg-white px-3 text-sm sm:w-72"
          value={selectedBranchId}
          onChange={(event) => setSelectedBranchId(event.target.value)}
        >
          <option value="">{t("dashboard.allBranches")}</option>
          {branches.map((branch) => (
            <option key={branch.id} value={branch.id}>
              {branch.code} · {branch.name}
            </option>
          ))}
        </select>
      </div>
      {error ? <div className="mb-4 rounded-md border border-red-200 bg-red-50 p-3 text-sm text-red-700">{error}</div> : null}
      <section className="grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
        {widgets.map((widget) => {
          const Icon = widget.icon;
          return (
            <article key={widget.label} className="rounded-md border border-line bg-white p-4">
              <div className="flex items-center justify-between">
                <span className="text-sm text-slate-500">{widget.label}</span>
                <Icon className="h-5 w-5 text-brand" />
              </div>
              <div className="mt-3 text-2xl font-bold">{widget.value}</div>
            </article>
          );
        })}
      </section>
      {loading ? (
        <div className="mt-5">
          <ListSkeleton rows={5} />
        </div>
      ) : (
        <section className="mt-5 grid gap-4 xl:grid-cols-[1.4fr_1fr]">
          <RevenueTrendChart data={salesTrend} />
          <RevenueDonut revenue={summary.revenue} profit={summary.profit} />
        </section>
      )}
      <section className="mt-5 grid gap-4 xl:grid-cols-[1.2fr_1fr]">
        <TopProductsChart products={summary.top_products} />
        <BranchComparisonChart branches={summary.branch_comparison} />
      </section>
      <section className="mt-5">
        <LowStockChart items={summary.low_stock} />
      </section>
      <section className="mt-5 grid gap-4 lg:grid-cols-3">
        <div className="rounded-md border border-line bg-white p-4">
          <h2 className="font-semibold">{t("dashboard.topProducts")}</h2>
          <div className="mt-3 space-y-2">
            {summary.top_products.length === 0 ? <div className="text-sm text-slate-500">{t("dashboard.noProductSales")}</div> : null}
            {summary.top_products.map((product) => (
              <div key={product.product_id} className="flex justify-between gap-3 rounded-md bg-field p-3 text-sm">
                <span className="font-medium">{product.name}</span>
                <span>{product.quantity} · {money(product.revenue)}</span>
              </div>
            ))}
          </div>
        </div>
        <div className="rounded-md border border-line bg-white p-4">
          <h2 className="font-semibold">{t("dashboard.lowStock")}</h2>
          <div className="mt-3 space-y-2">
            {summary.low_stock.length === 0 ? <div className="text-sm text-slate-500">{t("dashboard.noLowStock")}</div> : null}
            {summary.low_stock.map((item) => (
              <div key={`${item.branch_id}-${item.product_id}`} className="flex justify-between gap-3 rounded-md bg-field p-3 text-sm">
                <span className="font-medium">{item.branch_code} · {item.product_name}</span>
                <span>{item.quantity}</span>
              </div>
            ))}
          </div>
        </div>
        <div className="rounded-md border border-line bg-white p-4">
          <h2 className="font-semibold">{t("dashboard.branchComparison")}</h2>
          <div className="mt-3 space-y-2">
            {summary.branch_comparison.length === 0 ? <div className="text-sm text-slate-500">{t("dashboard.noBranchData")}</div> : null}
            {summary.branch_comparison.map((branch) => (
              <div key={branch.branch_id} className="flex justify-between gap-3 rounded-md bg-field p-3 text-sm">
                <span className="font-medium">{branch.branch_code} · {branch.branch_name}</span>
                <span>{branch.sales_count} · {money(branch.revenue)}</span>
              </div>
            ))}
          </div>
        </div>
      </section>
    </AppShell>
  );
}

type TrendPoint = {
  label: string;
  key: string;
  revenue: number;
  sales: number;
};

function localDateKey(date: Date) {
  const year = date.getFullYear();
  const month = `${date.getMonth() + 1}`.padStart(2, "0");
  const day = `${date.getDate()}`.padStart(2, "0");
  return `${year}-${month}-${day}`;
}

function buildSalesTrend(sales: Sale[], language: "en" | "th"): TrendPoint[] {
  const formatter = new Intl.DateTimeFormat(language === "th" ? "th-TH" : "en", { weekday: "short" });
  const today = new Date();
  const days = Array.from({ length: 7 }, (_, index) => {
    const date = new Date(today);
    date.setDate(today.getDate() - (6 - index));
    date.setHours(0, 0, 0, 0);
    return {
      label: formatter.format(date),
      key: localDateKey(date),
      revenue: 0,
      sales: 0
    };
  });
  const byKey = new Map(days.map((day) => [day.key, day]));
  sales.forEach((sale) => {
    const key = localDateKey(new Date(sale.created_at));
    const point = byKey.get(key);
    if (point) {
      point.revenue += sale.total;
      point.sales += 1;
    }
  });
  return days;
}

function RevenueTrendChart({ data }: { data: TrendPoint[] }) {
  const t = useI18nStore((state) => state.t);
  const maxRevenue = Math.max(...data.map((point) => point.revenue), 1);
  const mobilePoints = data.map((point, index) => {
    const x = 24 + index * 45;
    const y = 142 - (point.revenue / maxRevenue) * 96;
    return { ...point, x, y };
  });
  const desktopPoints = data.map((point, index) => {
    const x = 32 + index * 64;
    const y = 178 - (point.revenue / maxRevenue) * 132;
    return { ...point, x, y };
  });
  const mobileLine = mobilePoints.map((point) => `${point.x},${point.y}`).join(" ");
  const mobileArea = `24,142 ${mobileLine} ${mobilePoints[mobilePoints.length - 1]?.x ?? 294},142`;
  const desktopLine = desktopPoints.map((point) => `${point.x},${point.y}`).join(" ");
  const desktopArea = `32,178 ${desktopLine} ${desktopPoints[desktopPoints.length - 1]?.x ?? 416},178`;
  const totalRevenue = data.reduce((sum, point) => sum + point.revenue, 0);
  const totalSales = data.reduce((sum, point) => sum + point.sales, 0);

  return (
    <article className="rounded-md border border-line bg-white p-4">
      <div className="flex flex-col gap-1 sm:flex-row sm:items-start sm:justify-between">
        <div>
          <h2 className="font-semibold">{t("dashboard.trendTitle")}</h2>
          <p className="text-sm text-slate-500">{t("dashboard.trendDescription")}</p>
        </div>
        <div className="text-sm text-slate-500 sm:text-right">
          <div className="font-semibold text-slate-900">{compactMoney(totalRevenue)}</div>
          <div>{totalSales} {t("dashboard.saleCount")}</div>
        </div>
      </div>
      <div className="mt-4 min-w-0">
        <svg viewBox="0 0 318 178" className="h-auto w-full overflow-visible sm:hidden" preserveAspectRatio="xMidYMid meet" role="img" aria-label="Seven day revenue trend chart">
          <defs>
            <linearGradient id="revenueGradient" x1="0" x2="0" y1="0" y2="1">
              <stop offset="0%" stopColor="#2f7f75" stopOpacity="0.28" />
              <stop offset="100%" stopColor="#2f7f75" stopOpacity="0.04" />
            </linearGradient>
          </defs>
          {[46, 78, 110, 142].map((y) => (
            <line key={y} x1="24" x2="294" y1={y} y2={y} stroke="#e2e8f0" strokeWidth="1" />
          ))}
          <polygon points={mobileArea} fill="url(#revenueGradient)" />
          <polyline points={mobileLine} fill="none" stroke="#2f7f75" strokeWidth="3" strokeLinecap="round" strokeLinejoin="round" />
          {mobilePoints.map((point) => (
            <g key={point.key}>
              <circle cx={point.x} cy={point.y} r="4" fill="#2f7f75" stroke="white" strokeWidth="2" />
              <text x={point.x} y="164" textAnchor="middle" className="fill-slate-500 text-[10px]">
                {point.label}
              </text>
              <text x={point.x} y={Math.max(point.y - 10, 16)} textAnchor="middle" className="fill-slate-700 text-[9px] font-semibold">
                {point.revenue > 0 ? compactMoney(point.revenue) : ""}
              </text>
            </g>
          ))}
        </svg>
        <svg viewBox="0 0 448 220" className="hidden h-64 w-full overflow-visible sm:block" preserveAspectRatio="xMidYMid meet" role="img" aria-label="Seven day revenue trend chart">
          <defs>
            <linearGradient id="revenueGradientDesktop" x1="0" x2="0" y1="0" y2="1">
              <stop offset="0%" stopColor="#2f7f75" stopOpacity="0.28" />
              <stop offset="100%" stopColor="#2f7f75" stopOpacity="0.04" />
            </linearGradient>
          </defs>
          {[46, 90, 134, 178].map((y) => (
            <line key={y} x1="32" x2="416" y1={y} y2={y} stroke="#e2e8f0" strokeWidth="1" />
          ))}
          <polygon points={desktopArea} fill="url(#revenueGradientDesktop)" />
          <polyline points={desktopLine} fill="none" stroke="#2f7f75" strokeWidth="4" strokeLinecap="round" strokeLinejoin="round" />
          {desktopPoints.map((point) => (
            <g key={point.key}>
              <circle cx={point.x} cy={point.y} r="5" fill="#2f7f75" stroke="white" strokeWidth="3" />
              <text x={point.x} y="204" textAnchor="middle" className="fill-slate-500 text-[12px]">
                {point.label}
              </text>
              <text x={point.x} y={Math.max(point.y - 12, 16)} textAnchor="middle" className="fill-slate-700 text-[11px] font-semibold">
                {point.revenue > 0 ? compactMoney(point.revenue) : ""}
              </text>
            </g>
          ))}
        </svg>
      </div>
    </article>
  );
}

function RevenueDonut({ revenue, profit }: { revenue: number; profit: number }) {
  const t = useI18nStore((state) => state.t);
  const safeRevenue = Math.max(revenue, 1);
  const profitPercent = Math.max(0, Math.min(100, Math.round((profit / safeRevenue) * 100)));
  const costPercent = 100 - profitPercent;

  return (
    <article className="rounded-md border border-line bg-white p-4">
      <h2 className="font-semibold">{t("dashboard.revenueQuality")}</h2>
      <p className="text-sm text-slate-500">{t("dashboard.revenueQualityDescription")}</p>
      <div className="mt-5 flex flex-col items-center gap-5 sm:flex-row">
        <div
          className="grid h-44 w-44 place-items-center rounded-full"
          style={{ background: `conic-gradient(#2f7f75 0 ${profitPercent}%, #162033 ${profitPercent}% 100%)` }}
        >
          <div className="grid h-28 w-28 place-items-center rounded-full bg-white text-center">
            <div>
              <div className="text-2xl font-bold">{profitPercent}%</div>
              <div className="text-xs text-slate-500">{t("dashboard.margin")}</div>
            </div>
          </div>
        </div>
        <div className="w-full space-y-3 text-sm">
          <div className="flex items-center justify-between rounded-md bg-field p-3">
            <span className="flex items-center gap-2"><span className="h-3 w-3 rounded-sm bg-brand" />{t("dashboard.profit")}</span>
            <span className="font-semibold">{money(profit)}</span>
          </div>
          <div className="flex items-center justify-between rounded-md bg-field p-3">
            <span className="flex items-center gap-2"><span className="h-3 w-3 rounded-sm bg-slate-800" />{t("dashboard.costShare")}</span>
            <span className="font-semibold">{costPercent}%</span>
          </div>
          <div className="flex items-center justify-between rounded-md bg-field p-3">
            <span>{t("dashboard.revenue")}</span>
            <span className="font-semibold">{money(revenue)}</span>
          </div>
        </div>
      </div>
    </article>
  );
}

function TopProductsChart({ products }: { products: DashboardSummary["top_products"] }) {
  const t = useI18nStore((state) => state.t);
  const maxRevenue = Math.max(...products.map((product) => product.revenue), 1);
  return (
    <article className="rounded-md border border-line bg-white p-4">
      <h2 className="font-semibold">{t("dashboard.topProductsChart")}</h2>
      <p className="text-sm text-slate-500">{t("dashboard.topProductsDescription")}</p>
      <div className="mt-4 space-y-3">
        {products.length === 0 ? <div className="rounded-md bg-field p-4 text-sm text-slate-500">{t("dashboard.noProductSales")}</div> : null}
        {products.map((product) => (
          <div key={product.product_id} className="grid gap-2 text-sm">
            <div className="flex items-center justify-between gap-3">
              <span className="truncate font-medium">{product.name}</span>
              <span className="shrink-0 text-slate-500">{product.quantity} pcs · {compactMoney(product.revenue)}</span>
            </div>
            <div className="h-3 overflow-hidden rounded-full bg-field">
              <div className="h-full rounded-full bg-brand" style={{ width: `${Math.max(6, (product.revenue / maxRevenue) * 100)}%` }} />
            </div>
          </div>
        ))}
      </div>
    </article>
  );
}

function BranchComparisonChart({ branches }: { branches: DashboardSummary["branch_comparison"] }) {
  const t = useI18nStore((state) => state.t);
  const maxRevenue = Math.max(...branches.map((branch) => branch.revenue), 1);
  return (
    <article className="rounded-md border border-line bg-white p-4">
      <h2 className="font-semibold">{t("dashboard.branchRevenue")}</h2>
      <p className="text-sm text-slate-500">{t("dashboard.branchRevenueDescription")}</p>
      <div className="mt-4 flex h-64 items-end gap-3 overflow-x-auto border-b border-line pb-2 pt-7">
        {branches.length === 0 ? <div className="self-center text-sm text-slate-500">{t("dashboard.noBranchData")}</div> : null}
        {branches.map((branch) => {
          const height = Math.max(12, (branch.revenue / maxRevenue) * 150);
          return (
            <div key={branch.branch_id} className="flex min-w-20 flex-1 flex-col items-center gap-2 text-center">
              <div className="min-h-4 text-xs font-semibold text-slate-700">{compactMoney(branch.revenue)}</div>
              <div className="flex h-40 w-full items-end justify-center">
                <div className="w-10 rounded-t-md bg-slate-800" style={{ height }} />
              </div>
              <div className="w-full truncate text-xs text-slate-500">{branch.branch_code}</div>
            </div>
          );
        })}
      </div>
    </article>
  );
}

function LowStockChart({ items }: { items: DashboardSummary["low_stock"] }) {
  const t = useI18nStore((state) => state.t);
  const shownItems = items.slice(0, 8);
  const maxQuantity = Math.max(...shownItems.map((item) => item.quantity), 10);
  return (
    <article className="rounded-md border border-line bg-white p-4">
      <div className="flex flex-col gap-1 sm:flex-row sm:items-start sm:justify-between">
        <div>
          <h2 className="font-semibold">{t("dashboard.lowStockPressure")}</h2>
          <p className="text-sm text-slate-500">{t("dashboard.lowStockDescription")}</p>
        </div>
        <div className="text-sm font-semibold text-slate-700">{items.length} {t("dashboard.itemCount")}</div>
      </div>
      <div className="mt-4 grid gap-3 md:grid-cols-2">
        {shownItems.length === 0 ? <div className="rounded-md bg-field p-4 text-sm text-slate-500">{t("dashboard.noLowStock")}</div> : null}
        {shownItems.map((item) => (
          <div key={`${item.branch_id}-${item.product_id}`} className="rounded-md bg-field p-3">
            <div className="flex justify-between gap-3 text-sm">
              <span className="truncate font-medium">{item.branch_code} · {item.product_name}</span>
              <span className="font-semibold">{item.quantity}</span>
            </div>
            <div className="mt-2 h-2 overflow-hidden rounded-full bg-white">
              <div className="h-full rounded-full bg-accent" style={{ width: `${Math.max(8, (item.quantity / maxQuantity) * 100)}%` }} />
            </div>
          </div>
        ))}
      </div>
    </article>
  );
}
