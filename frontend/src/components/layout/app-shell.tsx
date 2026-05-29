"use client";

import { Button } from "@/components/ui/button";
import { ConfirmModal } from "@/components/ui/confirm-modal";
import { ThemeToggle } from "@/components/ui/theme-toggle";
import { Toaster } from "@/components/ui/toaster";
import { api } from "@/services/api";
import { useAuthStore } from "@/stores/auth-store";
import { ArrowRightLeft, BarChart3, Boxes, ClipboardList, ClipboardPlus, LayoutDashboard, LogOut, Package, Receipt, Settings, ShoppingCart, Users } from "lucide-react";
import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import type { ReactNode } from "react";
import { useEffect, useState } from "react";

const ownerLinks = [
  { href: "/dashboard", label: "Dashboard", icon: LayoutDashboard },
  { href: "/products", label: "Products", icon: Package },
  { href: "/categories", label: "Categories", icon: ClipboardPlus },
  { href: "/inventory", label: "Inventory", icon: Boxes },
  { href: "/all-stock", label: "All Stock", icon: ClipboardList },
  { href: "/stock-movements", label: "Movements", icon: ArrowRightLeft },
  { href: "/transfers", label: "Transfers", icon: Receipt },
  { href: "/employees", label: "Employees", icon: Users },
  { href: "/reports", label: "Reports", icon: BarChart3 },
  { href: "/settings", label: "Settings", icon: Settings }
];

const employeeLinks = [
  { href: "/pos", label: "POS", icon: ShoppingCart },
  { href: "/my-sales", label: "My Sales", icon: Receipt },
  { href: "/branch-inventory", label: "Branch Stock", icon: Boxes }
];

const managerLinks = [
  { href: "/dashboard", label: "Dashboard", icon: LayoutDashboard },
  { href: "/pos", label: "POS", icon: ShoppingCart },
  { href: "/my-sales", label: "My Sales", icon: Receipt },
  { href: "/branch-sales", label: "Branch Sales", icon: BarChart3 },
  { href: "/branch-inventory", label: "Branch Stock", icon: Boxes },
  { href: "/stock-receive", label: "Receive Stock", icon: ClipboardPlus },
  { href: "/employees", label: "Employees", icon: Users }
];

const defaultRouteByRole = {
  OWNER: "/dashboard",
  MANAGER: "/dashboard",
  EMPLOYEE: "/pos"
};

export function AppShell({ children }: { children: ReactNode }) {
  const path = usePathname();
  const router = useRouter();
  const user = useAuthStore((state) => state.user);
  const hydrate = useAuthStore((state) => state.hydrate);
  const clear = useAuthStore((state) => state.clear);
  const [mounted, setMounted] = useState(false);
  const [logoutModalOpen, setLogoutModalOpen] = useState(false);
  const links = user?.role === "EMPLOYEE" ? employeeLinks : user?.role === "MANAGER" ? managerLinks : ownerLinks;
  const canAccessCurrentPage = Boolean(user && links.some((link) => link.href === path));

  useEffect(() => {
    hydrate();
    setMounted(true);
  }, [hydrate]);

  useEffect(() => {
    if (!mounted) {
      return;
    }
    if (!user) {
      router.replace("/login");
      return;
    }
    if (!links.some((link) => link.href === path)) {
      router.replace(defaultRouteByRole[user.role]);
    }
  }, [links, mounted, path, router, user]);

  async function logout() {
    try {
      await api.logout();
    } finally {
      clear();
      window.location.assign("/login");
    }
  }

  if (!mounted || !user) {
    return (
      <main className="grid min-h-screen place-items-center bg-field px-4 text-sm font-medium text-slate-600">
        Redirecting to login...
      </main>
    );
  }

  if (!canAccessCurrentPage) {
    return (
      <main className="grid min-h-screen place-items-center bg-field px-4 text-sm font-medium text-slate-600">
        Redirecting to an allowed page...
      </main>
    );
  }

  const navItems = links.map((link) => {
    const Icon = link.icon;
    const active = path === link.href;
    return (
      <Link
        key={link.href}
        href={link.href}
        className={`flex min-w-fit items-center gap-2 rounded-md px-3 py-2 text-sm font-semibold transition ${
          active ? "bg-brand text-white shadow-sm" : "text-slate-600 hover:bg-brandSoft hover:text-brand"
        }`}
      >
        <Icon className="h-4 w-4" />
        {link.label}
      </Link>
    );
  });

  return (
    <div className="min-h-screen bg-field lg:flex">
      <Toaster />
      <aside className="hidden border-r border-line bg-white lg:sticky lg:top-0 lg:flex lg:h-screen lg:w-72 lg:flex-col">
        <div className="flex min-h-20 items-center border-b border-line px-5">
          <div>
            <div className="text-lg font-bold tracking-tight">Multi-Branch POS</div>
            <div className="mt-1 text-xs font-medium text-slate-500">{user?.name ?? "Not signed in"} · {user?.role ?? "GUEST"}</div>
          </div>
        </div>
        <nav className="space-y-1.5 p-3">
          {navItems}
        </nav>
        <div className="mt-auto space-y-2 border-t border-line p-3">
          <ThemeToggle />
          <Button className="w-full bg-slate-800 hover:bg-slate-700" onClick={() => setLogoutModalOpen(true)}>
            <LogOut className="h-4 w-4" />
            Logout
          </Button>
        </div>
      </aside>
      <div className="min-w-0 flex-1">
        <header className="sticky top-0 z-20 border-b border-line bg-white/95 backdrop-blur lg:hidden">
          <div className="flex min-h-16 items-center justify-between gap-3 px-4 lg:px-6">
            <div>
              <div className="text-base font-bold tracking-tight">Multi-Branch POS</div>
              <div className="text-xs text-slate-500">{user?.name ?? "Not signed in"} · {user?.role ?? "GUEST"}</div>
            </div>
            <div className="flex items-center gap-2">
              <ThemeToggle compact />
              <Button className="bg-slate-800 hover:bg-slate-700" onClick={() => setLogoutModalOpen(true)}>
                <LogOut className="h-4 w-4" />
                Logout
              </Button>
            </div>
          </div>
        </header>
        <main className="min-w-0 px-4 py-5 pb-28 sm:px-5 lg:p-7">{children}</main>
        <nav className="fixed inset-x-0 bottom-0 z-30 flex gap-2 overflow-x-auto border-t border-line bg-white/95 p-2 shadow-[0_-8px_20px_rgba(15,23,42,0.08)] backdrop-blur lg:hidden">
          {navItems}
        </nav>
      </div>
      <ConfirmModal
        open={logoutModalOpen}
        title="Confirm Logout"
        description="Are you sure you want to sign out of this POS session?"
        confirmLabel="Logout"
        danger
        onCancel={() => setLogoutModalOpen(false)}
        onConfirm={logout}
      />
    </div>
  );
}
