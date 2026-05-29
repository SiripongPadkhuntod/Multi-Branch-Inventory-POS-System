"use client";

import { useToastStore } from "@/stores/toast-store";
import { CheckCircle2, Info, X, XCircle } from "lucide-react";

const tone = {
  success: "border-emerald-200 bg-emerald-50 text-emerald-900",
  error: "border-red-200 bg-red-50 text-red-900",
  info: "border-slate-200 bg-white text-slate-900"
};

const icons = {
  success: CheckCircle2,
  error: XCircle,
  info: Info
};

export function Toaster() {
  const toasts = useToastStore((state) => state.toasts);
  const dismiss = useToastStore((state) => state.dismiss);

  return (
    <div className="fixed right-3 top-3 z-[60] flex w-[calc(100vw-24px)] max-w-sm flex-col gap-2 sm:right-5 sm:top-5">
      {toasts.map((toast) => {
        const Icon = icons[toast.type];
        return (
          <div key={toast.id} className={`flex gap-3 rounded-md border p-3 shadow-lg ${tone[toast.type]}`}>
            <Icon className="mt-0.5 h-5 w-5 shrink-0" />
            <div className="min-w-0 flex-1">
              <div className="font-semibold">{toast.title}</div>
              {toast.message ? <div className="mt-0.5 text-sm opacity-80">{toast.message}</div> : null}
            </div>
            <button className="rounded-md p-1 hover:bg-white/60" onClick={() => dismiss(toast.id)} title="Dismiss">
              <X className="h-4 w-4" />
            </button>
          </div>
        );
      })}
    </div>
  );
}
