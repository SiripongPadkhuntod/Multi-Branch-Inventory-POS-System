"use client";

import { AppShell } from "@/components/layout/app-shell";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { api } from "@/services/api";
import { useAuthStore } from "@/stores/auth-store";
import type { Branch, EmployeeSalesSummary, Role, User } from "@/types/domain";
import { BarChart3, Pencil, Save, UserPlus, X } from "lucide-react";
import { useEffect, useMemo, useState } from "react";

const money = (value: number) => `฿${(value / 100).toLocaleString("th-TH", { minimumFractionDigits: 2 })}`;

type FormState = {
  id: string;
  name: string;
  email: string;
  password: string;
  role: Role;
  branch_id: string;
  status: string;
};

const emptyForm: FormState = {
  id: "",
  name: "",
  email: "",
  password: "password123",
  role: "EMPLOYEE",
  branch_id: "",
  status: "ACTIVE"
};

export default function EmployeesPage() {
  const currentUser = useAuthStore((state) => state.user);
  const [branches, setBranches] = useState<Branch[]>([]);
  const [users, setUsers] = useState<User[]>([]);
  const [sales, setSales] = useState<EmployeeSalesSummary[]>([]);
  const [form, setForm] = useState<FormState>(emptyForm);
  const [formOpen, setFormOpen] = useState(false);
  const [notice, setNotice] = useState("");
  const [error, setError] = useState("");

  const branchMap = useMemo(() => new Map(branches.map((branch) => [branch.id, branch])), [branches]);

  async function load() {
    setError("");
    try {
      const [branchData, userData, salesData] = await Promise.all([api.myBranches(), api.users(), api.employeeSalesSummary()]);
      setBranches(Array.isArray(branchData) ? branchData : []);
      setUsers(Array.isArray(userData) ? userData : []);
      setSales(Array.isArray(salesData) ? salesData : []);
      setForm((value) => ({
        ...value,
        branch_id: value.branch_id || branchData[0]?.id || ""
      }));
    } catch (err) {
      setError(err instanceof Error ? err.message : "Cannot load employees");
    }
  }

  useEffect(() => {
    load();
  }, []);

  function createEmployee() {
    setNotice("");
    setError("");
    setForm({ ...emptyForm, branch_id: branches[0]?.id ?? "" });
    setFormOpen(true);
  }

  function edit(user: User) {
    setNotice("");
    setError("");
    setForm({
      id: user.id,
      name: user.name,
      email: user.email,
      password: "",
      role: user.role,
      branch_id: user.branch_id ?? "",
      status: user.status
    });
    setFormOpen(true);
  }

  function closeForm() {
    setFormOpen(false);
    setForm({ ...emptyForm, branch_id: branches[0]?.id ?? "" });
  }

  async function save() {
    setNotice("");
    setError("");
    if (!form.name || !form.email || !form.branch_id) {
      setError("Name, email, and branch are required");
      return;
    }
    if (!form.id && form.password.length < 8) {
      setError("Password must be at least 8 characters");
      return;
    }
    try {
      if (form.id) {
        await api.updateUser(form.id, {
          name: form.name,
          email: form.email,
          role: form.role,
          branch_id: form.branch_id,
          status: form.status
        });
        setNotice("Employee updated");
      } else {
        await api.createUser({
          name: form.name,
          email: form.email,
          password: form.password,
          role: form.role,
          branch_id: form.branch_id,
          status: form.status
        });
        setNotice("Employee created");
      }
      closeForm();
      await load();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Save failed");
    }
  }

  return (
    <AppShell>
      <div className="mb-5 flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold">Employees</h1>
          <p className="text-sm text-slate-500">
            {currentUser?.role === "MANAGER" ? "Manage employees in your assigned branches." : "Manage employees and review personal sales."}
          </p>
        </div>
        <Button onClick={createEmployee}>
          <UserPlus className="h-4 w-4" />
          New Employee
        </Button>
      </div>
      {notice ? <div className="mb-4 rounded-md border border-emerald-200 bg-emerald-50 p-3 text-sm text-emerald-800">{notice}</div> : null}
      {error ? <div className="mb-4 rounded-md border border-red-200 bg-red-50 p-3 text-sm text-red-700">{error}</div> : null}

      <div>
        <section className="space-y-4">
          <div className="overflow-hidden rounded-md border border-line bg-white">
            {users.length === 0 ? <div className="p-5 text-sm text-slate-500">No employees in your assigned branches.</div> : null}
            {users.map((user) => (
              <button key={user.id} className="flex w-full flex-wrap items-center justify-between gap-3 border-b border-line p-4 text-left last:border-b-0 hover:bg-field" onClick={() => edit(user)}>
                <div>
                  <div className="font-semibold">{user.name}</div>
                  <div className="text-xs text-slate-500">{user.email}</div>
                  <div className="mt-1 text-xs text-slate-500">{branchMap.get(user.branch_id ?? "")?.code ?? "No branch"} · {user.role}</div>
                </div>
                <div className="flex items-center gap-3 text-right text-sm">
                  <div>
                    <div className={user.status === "ACTIVE" ? "font-semibold text-emerald-700" : "font-semibold text-slate-500"}>{user.status}</div>
                    <div className="text-xs text-slate-500">Edit</div>
                  </div>
                  <Pencil className="h-4 w-4 text-brand" />
                </div>
              </button>
            ))}
          </div>

          <div className="rounded-md border border-line bg-white p-4">
            <h2 className="mb-3 flex items-center gap-2 font-semibold">
              <BarChart3 className="h-4 w-4 text-brand" />
              Sales By Employee
            </h2>
            <div className="space-y-2">
              {sales.length === 0 ? <div className="text-sm text-slate-500">No employee sales yet.</div> : null}
              {sales.map((item) => (
                <div key={item.user_id} className="grid gap-2 rounded-md bg-field p-3 text-sm md:grid-cols-[1fr_120px_120px] md:items-center">
                  <div>
                    <div className="font-semibold">{item.name}</div>
                    <div className="text-xs text-slate-500">{item.branch_code} · {item.email}</div>
                  </div>
                  <div>{item.sales_count} sale(s)</div>
                  <div className="font-bold md:text-right">{money(item.revenue)}</div>
                </div>
              ))}
            </div>
          </div>
        </section>
      </div>
      {formOpen ? (
        <div className="fixed inset-0 z-50 grid place-items-center bg-slate-950/50 px-4 py-6 backdrop-blur-sm">
          <div className="max-h-[calc(100vh-48px)] w-full max-w-xl overflow-y-auto rounded-md border border-line bg-white p-5 shadow-2xl">
            <div className="mb-4 flex items-start justify-between gap-3">
              <div>
                <h2 className="flex items-center gap-2 text-lg font-bold">
                  <UserPlus className="h-5 w-5 text-brand" />
                  {form.id ? "Edit Employee" : "New Employee"}
                </h2>
                <p className="mt-1 text-sm text-slate-500">
                  {form.id ? "Update employee profile, branch, role, and status." : "Create a new employee account for an assigned branch."}
                </p>
              </div>
              <button
                className="grid h-9 w-9 shrink-0 place-items-center rounded-md border border-line text-slate-600 hover:bg-field"
                onClick={closeForm}
                aria-label="Close employee form"
              >
                <X className="h-4 w-4" />
              </button>
            </div>
            <div className="grid gap-3 sm:grid-cols-2">
              <label className="block text-sm font-semibold sm:col-span-2">
                Name
                <Input className="mt-1" value={form.name} onChange={(event) => setForm({ ...form, name: event.target.value })} />
              </label>
              <label className="block text-sm font-semibold sm:col-span-2">
                Email
                <Input className="mt-1" value={form.email} onChange={(event) => setForm({ ...form, email: event.target.value })} />
              </label>
              {!form.id ? (
                <label className="block text-sm font-semibold sm:col-span-2">
                  Password
                  <Input className="mt-1" value={form.password} onChange={(event) => setForm({ ...form, password: event.target.value })} />
                </label>
              ) : null}
              <label className="block text-sm font-semibold">
                Branch
                <select className="mt-1 h-10 w-full rounded-md border border-line bg-white px-3 text-sm" value={form.branch_id} onChange={(event) => setForm({ ...form, branch_id: event.target.value })}>
                  {branches.map((branch) => (
                    <option key={branch.id} value={branch.id}>{branch.code} · {branch.name}</option>
                  ))}
                </select>
              </label>
              <label className="block text-sm font-semibold">
                Role
                <select className="mt-1 h-10 w-full rounded-md border border-line bg-white px-3 text-sm" value={form.role} onChange={(event) => setForm({ ...form, role: event.target.value as Role })}>
                  <option value="EMPLOYEE">Employee</option>
                  {currentUser?.role === "OWNER" ? <option value="MANAGER">Manager</option> : null}
                </select>
              </label>
              <label className="block text-sm font-semibold sm:col-span-2">
                Status
                <select className="mt-1 h-10 w-full rounded-md border border-line bg-white px-3 text-sm" value={form.status} onChange={(event) => setForm({ ...form, status: event.target.value })}>
                  <option value="ACTIVE">Active</option>
                  <option value="INACTIVE">Inactive</option>
                </select>
              </label>
            </div>
            <div className="mt-5 flex flex-col-reverse gap-2 sm:flex-row sm:justify-end">
              <Button className="!bg-white !text-slate-700 ring-1 ring-line hover:!bg-field" onClick={closeForm}>
                Cancel
              </Button>
              <Button onClick={save}>
                <Save className="h-4 w-4" />
                Save
              </Button>
            </div>
          </div>
        </div>
      ) : null}
    </AppShell>
  );
}
