"use client";

import { Button } from "@/components/ui/button";
import type { ReactNode } from "react";

type ConfirmModalProps = {
  open: boolean;
  title: string;
  description: string;
  confirmLabel?: string;
  cancelLabel?: string;
  danger?: boolean;
  children?: ReactNode;
  onConfirm: () => void;
  onCancel: () => void;
};

export function ConfirmModal({
  open,
  title,
  description,
  confirmLabel = "Confirm",
  cancelLabel = "Cancel",
  danger = false,
  children,
  onConfirm,
  onCancel
}: ConfirmModalProps) {
  if (!open) {
    return null;
  }

  return (
    <div className="fixed inset-0 z-50 grid place-items-center bg-slate-950/50 px-4 backdrop-blur-sm">
      <div className="w-full max-w-md rounded-md border border-line bg-white p-5 shadow-2xl">
        <h2 className="text-lg font-bold">{title}</h2>
        <p className="mt-2 text-sm text-slate-600">{description}</p>
        {children ? <div className="mt-4 rounded-md border border-line bg-field p-3 text-sm text-slate-700">{children}</div> : null}
        <div className="mt-5 flex flex-col-reverse gap-2 sm:flex-row sm:justify-end">
          <Button className="!bg-white !text-slate-700 ring-1 ring-line hover:!bg-field" onClick={onCancel}>
            {cancelLabel}
          </Button>
          <Button className={danger ? "bg-red-600 hover:bg-red-700" : ""} onClick={onConfirm}>
            {confirmLabel}
          </Button>
        </div>
      </div>
    </div>
  );
}
