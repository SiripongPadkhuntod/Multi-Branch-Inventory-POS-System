"use client";

import { Button } from "@/components/ui/button";
import type { SaleDetail } from "@/types/domain";
import { Printer, X } from "lucide-react";

const money = (value: number) => `฿${(value / 100).toLocaleString("th-TH", { minimumFractionDigits: 2 })}`;

type ReceiptModalProps = {
  open: boolean;
  detail: SaleDetail | null;
  onClose: () => void;
};

export function ReceiptModal({ open, detail, onClose }: ReceiptModalProps) {
  if (!open || !detail) {
    return null;
  }

  return (
    <div className="fixed inset-0 z-50 grid place-items-center bg-slate-950/50 px-4 py-6 backdrop-blur-sm">
      <div className="max-h-[calc(100vh-48px)] w-full max-w-md overflow-y-auto rounded-md border border-line bg-white shadow-2xl">
        <div className="flex items-center justify-between border-b border-line p-4">
          <div>
            <h2 className="text-lg font-bold">Receipt</h2>
            <p className="text-xs text-slate-500">{detail.receipt_number}</p>
          </div>
          <button className="grid h-9 w-9 place-items-center rounded-md border border-line text-slate-600 hover:bg-field" onClick={onClose} aria-label="Close receipt">
            <X className="h-4 w-4" />
          </button>
        </div>
        <ReceiptPaper detail={detail} />
        <div className="flex flex-col-reverse gap-2 border-t border-line p-4 sm:flex-row sm:justify-end">
          <Button className="!bg-white !text-slate-700 ring-1 ring-line hover:!bg-field" onClick={onClose}>Close</Button>
          <Button onClick={() => printReceipt(detail)}>
            <Printer className="h-4 w-4" />
            Print
          </Button>
        </div>
      </div>
    </div>
  );
}

export function ReceiptPaper({ detail }: { detail: SaleDetail }) {
  return (
    <div className="bg-white p-5 font-mono text-sm text-slate-950">
      <div className="text-center">
        <div className="text-base font-bold">Multi-Branch POS</div>
        <div>{detail.branch_code} · {detail.branch_name}</div>
        <div className="mt-2 text-xs">{new Date(detail.created_at).toLocaleString()}</div>
        <div className="text-xs">Receipt: {detail.receipt_number}</div>
        <div className="text-xs">Cashier: {detail.employee_name}</div>
      </div>
      <div className="my-4 border-t border-dashed border-slate-400" />
      <div className="space-y-3">
        {detail.items.map((item) => (
          <div key={item.id}>
            <div className="font-semibold">{item.product_name}</div>
            <div className="flex justify-between gap-3 text-xs">
              <span>{item.quantity} x {money(item.final_price)}</span>
              <span>{money(item.line_total)}</span>
            </div>
            {item.discount_amount > 0 ? (
              <div className="flex justify-between gap-3 text-xs text-slate-500">
                <span>Discount {item.discount_reason ? `(${item.discount_reason})` : ""}</span>
                <span>-{money(item.discount_amount)}</span>
              </div>
            ) : null}
          </div>
        ))}
      </div>
      <div className="my-4 border-t border-dashed border-slate-400" />
      <div className="space-y-1">
        <ReceiptRow label="Subtotal" value={money(detail.subtotal)} />
        <ReceiptRow label="Discount" value={money(detail.discount)} />
        <ReceiptRow label="Tax" value={money(detail.tax)} />
        <div className="mt-2 flex justify-between border-t border-slate-300 pt-2 text-base font-bold">
          <span>Total</span>
          <span>{money(detail.total)}</span>
        </div>
      </div>
      <div className="my-4 border-t border-dashed border-slate-400" />
      <div className="space-y-1">
        {detail.payments.map((payment) => (
          <ReceiptRow key={payment.id} label={payment.payment_method} value={money(payment.amount)} />
        ))}
      </div>
      <div className="mt-5 text-center text-xs">Thank you</div>
    </div>
  );
}

function ReceiptRow({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex justify-between gap-3">
      <span>{label}</span>
      <span>{value}</span>
    </div>
  );
}

export function printReceipt(detail: SaleDetail) {
  const win = window.open("", "_blank", "width=420,height=720");
  if (!win) {
    window.print();
    return;
  }
  win.document.write(receiptHtml(detail));
  win.document.close();
  win.focus();
  win.print();
}

function receiptHtml(detail: SaleDetail) {
  const itemRows = detail.items.map((item) => `
    <div class="item">
      <div class="name">${escapeHtml(item.product_name)}</div>
      <div class="row small"><span>${item.quantity} x ${money(item.final_price)}</span><span>${money(item.line_total)}</span></div>
      ${item.discount_amount > 0 ? `<div class="row muted"><span>Discount ${item.discount_reason ? `(${escapeHtml(item.discount_reason)})` : ""}</span><span>-${money(item.discount_amount)}</span></div>` : ""}
    </div>
  `).join("");
  const payments = detail.payments.map((payment) => `<div class="row"><span>${escapeHtml(payment.payment_method)}</span><span>${money(payment.amount)}</span></div>`).join("");

  return `<!doctype html>
  <html>
    <head>
      <title>${escapeHtml(detail.receipt_number)}</title>
      <style>
        * { box-sizing: border-box; }
        body { margin: 0; background: #fff; color: #0f172a; font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", monospace; }
        .receipt { width: 320px; margin: 0 auto; padding: 20px 14px; font-size: 13px; }
        .center { text-align: center; }
        .brand { font-size: 16px; font-weight: 800; }
        .small, .muted { font-size: 11px; }
        .muted { color: #64748b; }
        .dash { border-top: 1px dashed #94a3b8; margin: 14px 0; }
        .item { margin-bottom: 10px; }
        .name { font-weight: 700; }
        .row { display: flex; justify-content: space-between; gap: 12px; margin: 4px 0; }
        .total { border-top: 1px solid #cbd5e1; margin-top: 8px; padding-top: 8px; font-size: 16px; font-weight: 800; }
        @page { size: 80mm auto; margin: 4mm; }
      </style>
    </head>
    <body>
      <main class="receipt">
        <div class="center">
          <div class="brand">Multi-Branch POS</div>
          <div>${escapeHtml(detail.branch_code)} · ${escapeHtml(detail.branch_name)}</div>
          <div class="small">${new Date(detail.created_at).toLocaleString()}</div>
          <div class="small">Receipt: ${escapeHtml(detail.receipt_number)}</div>
          <div class="small">Cashier: ${escapeHtml(detail.employee_name)}</div>
        </div>
        <div class="dash"></div>
        ${itemRows}
        <div class="dash"></div>
        <div class="row"><span>Subtotal</span><span>${money(detail.subtotal)}</span></div>
        <div class="row"><span>Discount</span><span>${money(detail.discount)}</span></div>
        <div class="row"><span>Tax</span><span>${money(detail.tax)}</span></div>
        <div class="row total"><span>Total</span><span>${money(detail.total)}</span></div>
        <div class="dash"></div>
        ${payments}
        <div class="center small" style="margin-top:18px;">Thank you</div>
      </main>
      <script>setTimeout(function(){ window.close(); }, 800);</script>
    </body>
  </html>`;
}

function escapeHtml(value: string) {
  return value.replace(/[&<>"']/g, (char) => ({
    "&": "&amp;",
    "<": "&lt;",
    ">": "&gt;",
    "\"": "&quot;",
    "'": "&#039;"
  }[char] ?? char));
}
