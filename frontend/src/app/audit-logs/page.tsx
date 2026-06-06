"use client";

import { AppShell } from "@/components/layout/app-shell";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { ListSkeleton } from "@/components/ui/skeleton";
import { api } from "@/services/api";
import { useI18nStore } from "@/stores/i18n-store";
import type { AuditLog } from "@/types/domain";
import { ChevronDown, ChevronUp, FileClock, Search, X } from "lucide-react";
import { useEffect, useMemo, useState } from "react";

const actionOptions = ["", "REFUND", "CREATE", "UPDATE", "DELETE", "LOGIN", "TRANSFER", "ADJUSTMENT"];
const entityOptions = ["", "sale", "product", "inventory", "user", "branch", "category", "transfer"];

export default function AuditLogsPage() {
  const [logs, setLogs] = useState<AuditLog[]>([]);
  const [expandedId, setExpandedId] = useState("");
  const [query, setQuery] = useState("");
  const [action, setAction] = useState("");
  const [entityType, setEntityType] = useState("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const t = useI18nStore((state) => state.t);

  const actionCounts = useMemo(() => {
    return logs.reduce<Record<string, number>>((acc, log) => {
      acc[log.action] = (acc[log.action] ?? 0) + 1;
      return acc;
    }, {});
  }, [logs]);

  async function load() {
    setLoading(true);
    setError("");
    try {
      const data = await api.auditLogs({ q: query, action, entity_type: entityType, limit: 150 });
      setLogs(Array.isArray(data) ? data : []);
      setExpandedId("");
    } catch (err) {
      setLogs([]);
      setError(err instanceof Error ? err.message : t("audit.empty"));
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    load();
  }, []);

  function clearFilters() {
    setQuery("");
    setAction("");
    setEntityType("");
    setTimeout(() => load(), 0);
  }

  return (
    <AppShell>
      <div className="mb-5 flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold">{t("audit.title")}</h1>
          <p className="text-sm text-slate-500">{t("audit.description")}</p>
        </div>
        <div className="rounded-md border border-line bg-white px-4 py-3 text-sm">
          <span className="text-slate-500">{t("audit.loaded")}</span> <span className="font-bold">{logs.length}</span>
        </div>
      </div>

      <section className="mb-4 rounded-md border border-line bg-white p-4 shadow-sm">
        <div className="grid gap-3 lg:grid-cols-[1fr_180px_180px_auto_auto]">
          <label className="block text-sm font-semibold">
            {t("common.search")}
            <Input className="mt-1" placeholder="Actor, action, entity id..." value={query} onChange={(event) => setQuery(event.target.value)} onKeyDown={(event) => event.key === "Enter" && load()} />
          </label>
          <label className="block text-sm font-semibold">
            {t("field.action")}
            <select className="mt-1 h-10 w-full rounded-md border border-line bg-white px-3 text-sm" value={action} onChange={(event) => setAction(event.target.value)}>
              {actionOptions.map((option) => <option key={option || "all"} value={option}>{option || "All actions"}</option>)}
            </select>
          </label>
          <label className="block text-sm font-semibold">
            Entity
            <select className="mt-1 h-10 w-full rounded-md border border-line bg-white px-3 text-sm" value={entityType} onChange={(event) => setEntityType(event.target.value)}>
              {entityOptions.map((option) => <option key={option || "all"} value={option}>{option || "All entities"}</option>)}
            </select>
          </label>
          <Button className="self-end" onClick={load}>
            <Search className="h-4 w-4" />
            {t("common.apply")}
          </Button>
          <Button className="self-end !bg-white !text-slate-700 ring-1 ring-line hover:!bg-field" onClick={clearFilters}>
            <X className="h-4 w-4" />
            {t("common.clear")}
          </Button>
        </div>
      </section>

      {error ? <div className="mb-4 rounded-md border border-red-200 bg-red-50 p-3 text-sm text-red-700">{error}</div> : null}

      {Object.keys(actionCounts).length > 0 ? (
        <section className="mb-4 grid gap-2 sm:grid-cols-2 xl:grid-cols-4">
          {Object.entries(actionCounts).slice(0, 4).map(([name, count]) => (
            <div key={name} className="rounded-md border border-line bg-white p-3 text-sm">
              <div className="text-slate-500">{name}</div>
              <div className="mt-1 text-xl font-bold">{count}</div>
            </div>
          ))}
        </section>
      ) : null}

      {loading ? <ListSkeleton rows={5} /> : (
        <section className="overflow-hidden rounded-md border border-line bg-white shadow-sm">
          {logs.length === 0 ? <div className="p-5 text-sm text-slate-500">{t("audit.empty")}</div> : null}
          {logs.map((log) => {
            const expanded = expandedId === log.id;
            return (
              <div key={log.id} className="border-b border-line last:border-b-0">
                <button className="grid w-full gap-3 p-4 text-left hover:bg-field md:grid-cols-[1fr_160px_160px_32px] md:items-center" onClick={() => setExpandedId(expanded ? "" : log.id)}>
                  <div className="flex min-w-0 items-start gap-3">
                    <FileClock className="mt-0.5 h-5 w-5 shrink-0 text-brand" />
                    <div className="min-w-0">
                      <div className="flex flex-wrap items-center gap-2">
                        <span className="rounded-md bg-brandSoft px-2 py-1 text-xs font-bold text-brand">{log.action}</span>
                        <span className="font-semibold">{log.entity_type}</span>
                      </div>
                      <div className="mt-1 truncate text-xs text-slate-500">{log.entity_id}</div>
                    </div>
                  </div>
                  <div className="text-sm">
                    <div className="font-semibold">{log.user_name || "System"}</div>
                    <div className="truncate text-xs text-slate-500">{log.user_email || log.user_id}</div>
                  </div>
                  <div className="text-sm md:text-right">
                    <div className="font-medium">{new Date(log.created_at).toLocaleDateString()}</div>
                    <div className="text-xs text-slate-500">{new Date(log.created_at).toLocaleTimeString()}</div>
                  </div>
                  {expanded ? <ChevronUp className="h-4 w-4" /> : <ChevronDown className="h-4 w-4" />}
                </button>
                {expanded ? (
                  <div className="grid gap-3 border-t border-line bg-slate-50 p-4 lg:grid-cols-2">
                    <Payload title="Old Data" value={log.old_data} />
                    <Payload title={t("audit.newData")} value={log.new_data} />
                    {log.ip_address ? <div className="text-xs text-slate-500 lg:col-span-2">IP address: {log.ip_address}</div> : null}
                  </div>
                ) : null}
              </div>
            );
          })}
        </section>
      )}
    </AppShell>
  );
}

function Payload({ title, value }: { title: string; value: string }) {
  return (
    <div className="rounded-md border border-line bg-white p-3">
      <div className="mb-2 text-xs font-semibold text-slate-500">{title}</div>
      <pre className="max-h-56 overflow-auto whitespace-pre-wrap break-words rounded-md bg-field p-3 text-xs">{formatPayload(value)}</pre>
    </div>
  );
}

function formatPayload(value: string) {
  if (!value) {
    return "No data";
  }
  try {
    return JSON.stringify(JSON.parse(value), null, 2);
  } catch {
    return value;
  }
}
