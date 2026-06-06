"use client";

import { AppShell } from "@/components/layout/app-shell";
import { Button } from "@/components/ui/button";
import { ConfirmModal } from "@/components/ui/confirm-modal";
import { Input } from "@/components/ui/input";
import { api } from "@/services/api";
import { useI18nStore } from "@/stores/i18n-store";
import type { Category } from "@/types/domain";
import { Edit3, Plus, Save, Tags, Trash2, X } from "lucide-react";
import { useEffect, useState } from "react";

const emptyForm = { id: "", name: "", description: "" };

export default function CategoriesPage() {
  const [categories, setCategories] = useState<Category[]>([]);
  const [form, setForm] = useState(emptyForm);
  const [formOpen, setFormOpen] = useState(false);
  const [confirmUpdateOpen, setConfirmUpdateOpen] = useState(false);
  const [notice, setNotice] = useState("");
  const [error, setError] = useState("");
  const t = useI18nStore((state) => state.t);

  async function load() {
    setError("");
    try {
      const data = await api.categories();
      setCategories(Array.isArray(data) ? data : []);
    } catch (err) {
      setError(err instanceof Error ? err.message : t("categories.loadFailed"));
    }
  }

  useEffect(() => {
    load();
  }, []);

  function resetForm() {
    setNotice("");
    setError("");
    setForm(emptyForm);
    setFormOpen(true);
  }

  function edit(category: Category) {
    setNotice("");
    setError("");
    setForm(category);
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
    if (!form.name.trim()) {
      setError(t("categories.nameRequired"));
      return;
    }
    try {
      if (form.id) {
        await api.updateCategory(form.id, { name: form.name, description: form.description });
        setNotice(t("categories.updated"));
      } else {
        await api.createCategory({ name: form.name, description: form.description });
        setNotice(t("categories.created"));
      }
      closeForm();
      await load();
    } catch (err) {
      setError(err instanceof Error ? err.message : t("products.saveFailed"));
    }
  }

  async function remove() {
    if (!form.id) return;
    setNotice("");
    setError("");
    try {
      await api.deleteCategory(form.id);
      setNotice(t("categories.deleted"));
      closeForm();
      await load();
    } catch (err) {
      setError(err instanceof Error ? err.message : t("categories.deleteFailed"));
    }
  }

  return (
    <AppShell>
      <div className="mb-5 flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold">{t("categories.title")}</h1>
          <p className="text-sm text-slate-500">{t("categories.description")}</p>
        </div>
        <Button onClick={resetForm}>
          <Plus className="h-4 w-4" />
          {t("categories.new")}
        </Button>
      </div>
      {notice ? <div className="mb-4 rounded-md border border-emerald-200 bg-emerald-50 p-3 text-sm text-emerald-800">{notice}</div> : null}
      {error ? <div className="mb-4 rounded-md border border-red-200 bg-red-50 p-3 text-sm text-red-700">{error}</div> : null}
      <div>
        <section className="overflow-hidden rounded-md border border-line bg-white">
          {categories.length === 0 ? <div className="p-5 text-sm text-slate-500">{t("categories.empty")}</div> : null}
          {categories.map((category) => (
            <div key={category.id} className="flex flex-wrap items-center justify-between gap-3 border-b border-line p-4 last:border-b-0">
              <div className="flex min-w-0 items-center gap-3">
                <Tags className="h-5 w-5 shrink-0 text-brand" />
                <div className="min-w-0">
                <div className="font-semibold">{category.name}</div>
                <div className="text-xs text-slate-500">{category.description || t("common.noDescription")}</div>
                </div>
              </div>
              <Button className="bg-slate-800 hover:bg-slate-700" onClick={() => edit(category)}>
                <Edit3 className="h-4 w-4" />
                {t("common.edit")}
              </Button>
            </div>
          ))}
        </section>
      </div>
      {formOpen ? (
        <div className="fixed inset-0 z-50 grid place-items-center bg-slate-950/50 px-4 py-6 backdrop-blur-sm">
          <div className="max-h-[calc(100vh-48px)] w-full max-w-lg overflow-y-auto rounded-md border border-line bg-white p-5 shadow-2xl">
            <div className="mb-4 flex items-start justify-between gap-3">
              <div>
                <h2 className="flex items-center gap-2 text-lg font-bold">
                  <Tags className="h-5 w-5 text-brand" />
                  {form.id ? t("categories.edit") : t("categories.new")}
                </h2>
                <p className="mt-1 text-sm text-slate-500">
                  {form.id ? t("categories.editDescription") : t("categories.createDescription")}
                </p>
              </div>
              <button className="grid h-9 w-9 shrink-0 place-items-center rounded-md border border-line text-slate-600 hover:bg-field" onClick={closeForm} aria-label={t("common.close")}>
                <X className="h-4 w-4" />
              </button>
            </div>
            <div className="space-y-3">
              <label className="block text-sm font-semibold">{t("field.name")}<Input className="mt-1" value={form.name} onChange={(event) => setForm({ ...form, name: event.target.value })} /></label>
              <label className="block text-sm font-semibold">{t("field.description")}<Input className="mt-1" value={form.description} onChange={(event) => setForm({ ...form, description: event.target.value })} /></label>
            </div>
            <div className="mt-5 flex flex-col-reverse gap-2 sm:flex-row sm:justify-between">
              <div>
                {form.id ? <Button className="w-full bg-red-600 hover:bg-red-700 sm:w-auto" onClick={remove}><Trash2 className="h-4 w-4" />{t("common.delete")}</Button> : null}
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
        title={t("categories.confirmTitle")}
        description={t("categories.confirmDescription")}
        confirmLabel={t("categories.confirmButton")}
        onCancel={() => setConfirmUpdateOpen(false)}
        onConfirm={save}
      >
        <div className="font-semibold">{form.name || t("categories.selected")}</div>
        <div className="text-xs text-slate-500">{form.description || t("common.noDescription")}</div>
      </ConfirmModal>
    </AppShell>
  );
}
