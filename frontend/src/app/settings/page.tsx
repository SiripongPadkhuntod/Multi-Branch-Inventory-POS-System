"use client";

import { AppShell } from "@/components/layout/app-shell";
import { Button } from "@/components/ui/button";
import { ConfirmModal } from "@/components/ui/confirm-modal";
import { Input } from "@/components/ui/input";
import { api } from "@/services/api";
import { useI18nStore } from "@/stores/i18n-store";
import type { Branch } from "@/types/domain";
import { Building2, Edit3, Plus, Save, X } from "lucide-react";
import { useEffect, useState } from "react";

const emptyForm = { id: "", code: "", name: "", address: "", phone: "", status: "ACTIVE" };

export default function SettingsPage() {
  const [branches, setBranches] = useState<Branch[]>([]);
  const [form, setForm] = useState(emptyForm);
  const [formOpen, setFormOpen] = useState(false);
  const [confirmUpdateOpen, setConfirmUpdateOpen] = useState(false);
  const [notice, setNotice] = useState("");
  const [error, setError] = useState("");
  const t = useI18nStore((state) => state.t);

  async function load() {
    setError("");
    try {
      const data = await api.myBranches();
      setBranches(Array.isArray(data) ? data : []);
    } catch (err) {
      setError(err instanceof Error ? err.message : t("settings.loadFailed"));
    }
  }

  useEffect(() => {
    load();
  }, []);

  function edit(branch: Branch) {
    setNotice("");
    setError("");
    setForm({
      id: branch.id,
      code: branch.code,
      name: branch.name,
      address: branch.address,
      phone: branch.phone,
      status: branch.status
    });
    setFormOpen(true);
  }

  function resetForm() {
    setNotice("");
    setError("");
    setForm(emptyForm);
    setFormOpen(true);
  }

  function closeForm() {
    setFormOpen(false);
    setConfirmUpdateOpen(false);
    setForm(emptyForm);
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
    if (!form.code || !form.name) {
      setError(t("settings.branchRequired"));
      return;
    }
    const payload = {
      code: form.code,
      name: form.name,
      address: form.address,
      phone: form.phone,
      status: form.status
    };
    try {
      if (form.id) {
        await api.updateBranch(form.id, payload);
        setNotice(t("settings.branchUpdated"));
      } else {
        await api.createBranch(payload);
        setNotice(t("settings.branchCreated"));
      }
      closeForm();
      await load();
    } catch (err) {
      setError(err instanceof Error ? err.message : t("products.saveFailed"));
    }
  }

  return (
    <AppShell>
      <div className="mb-5 flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold">{t("settings.title")}</h1>
          <p className="text-sm text-slate-500">{t("settings.description")}</p>
        </div>
        <Button onClick={resetForm}>
          <Plus className="h-4 w-4" />
          {t("settings.newBranch")}
        </Button>
      </div>
      {notice ? <div className="mb-4 rounded-md border border-emerald-200 bg-emerald-50 p-3 text-sm text-emerald-800">{notice}</div> : null}
      {error ? <div className="mb-4 rounded-md border border-red-200 bg-red-50 p-3 text-sm text-red-700">{error}</div> : null}

      <div>
        <section className="overflow-hidden rounded-md border border-line bg-white">
          <div className="border-b border-line p-4 font-semibold">{t("settings.branches")}</div>
          {branches.length === 0 ? <div className="p-5 text-sm text-slate-500">{t("settings.emptyBranches")}</div> : null}
          {branches.map((branch) => (
            <div key={branch.id} className="flex flex-wrap items-center justify-between gap-3 border-b border-line p-4 last:border-b-0">
              <div className="flex min-w-0 items-center gap-3">
                <Building2 className="h-5 w-5 shrink-0 text-brand" />
                <div className="min-w-0">
                  <div className="truncate font-semibold">{branch.code} · {branch.name}</div>
                  <div className="text-xs text-slate-500">{branch.address || t("common.noAddress")} · {branch.phone || t("common.noPhone")}</div>
                </div>
              </div>
              <div className="flex items-center gap-3">
                <div className={branch.status === "ACTIVE" ? "text-sm font-semibold text-emerald-700" : "text-sm font-semibold text-slate-500"}>{branch.status}</div>
                <Button className="bg-slate-800 hover:bg-slate-700" onClick={() => edit(branch)}>
                  <Edit3 className="h-4 w-4" />
                  {t("common.edit")}
                </Button>
              </div>
            </div>
          ))}
        </section>
      </div>
      {formOpen ? (
        <div className="fixed inset-0 z-50 grid place-items-center bg-slate-950/50 px-4 py-6 backdrop-blur-sm">
          <div className="max-h-[calc(100vh-48px)] w-full max-w-xl overflow-y-auto rounded-md border border-line bg-white p-5 shadow-2xl">
            <div className="mb-4 flex items-start justify-between gap-3">
              <div>
                <h2 className="flex items-center gap-2 text-lg font-bold">
                  <Building2 className="h-5 w-5 text-brand" />
                  {form.id ? t("settings.editBranch") : t("settings.newBranch")}
                </h2>
                <p className="mt-1 text-sm text-slate-500">
                  {form.id ? t("settings.editBranchDescription") : t("settings.createBranchDescription")}
                </p>
              </div>
              <button className="grid h-9 w-9 shrink-0 place-items-center rounded-md border border-line text-slate-600 hover:bg-field" onClick={closeForm} aria-label={t("common.close")}>
                <X className="h-4 w-4" />
              </button>
            </div>
            <div className="space-y-3">
              <div className="grid gap-3 sm:grid-cols-2">
                <label className="block text-sm font-semibold">{t("field.code")}<Input className="mt-1" value={form.code} onChange={(event) => setForm({ ...form, code: event.target.value })} /></label>
                <label className="block text-sm font-semibold">{t("field.phone")}<Input className="mt-1" value={form.phone} onChange={(event) => setForm({ ...form, phone: event.target.value })} /></label>
              </div>
              <label className="block text-sm font-semibold">{t("field.name")}<Input className="mt-1" value={form.name} onChange={(event) => setForm({ ...form, name: event.target.value })} /></label>
              <label className="block text-sm font-semibold">{t("field.address")}<Input className="mt-1" value={form.address} onChange={(event) => setForm({ ...form, address: event.target.value })} /></label>
              <label className="block text-sm font-semibold">{t("field.status")}
                <select className="mt-1 h-10 w-full rounded-md border border-line bg-white px-3 text-sm" value={form.status} onChange={(event) => setForm({ ...form, status: event.target.value })}>
                  <option value="ACTIVE">{t("status.active")}</option>
                  <option value="INACTIVE">{t("status.inactive")}</option>
                </select>
              </label>
            </div>
            <div className="mt-5 flex flex-col-reverse gap-2 sm:flex-row sm:justify-end">
              <Button className="!bg-white !text-slate-700 ring-1 ring-line hover:!bg-field" onClick={closeForm}>{t("common.cancel")}</Button>
              <Button onClick={requestSave}><Save className="h-4 w-4" />{t("settings.saveBranch")}</Button>
            </div>
          </div>
        </div>
      ) : null}
      <ConfirmModal
        open={confirmUpdateOpen}
        title={t("settings.confirmTitle")}
        description={t("settings.confirmDescription")}
        confirmLabel={t("settings.confirmButton")}
        onCancel={() => setConfirmUpdateOpen(false)}
        onConfirm={save}
      >
        <div className="font-semibold">{form.code || t("settings.branchFallback")} · {form.name || t("settings.unnamed")}</div>
        <div className="text-xs text-slate-500">{form.address || t("common.noAddress")} · {form.phone || t("common.noPhone")}</div>
      </ConfirmModal>
    </AppShell>
  );
}
