"use client";

import { Button } from "@/components/ui/button";
import { ConfirmModal } from "@/components/ui/confirm-modal";
import { LanguageToggle } from "@/components/ui/language-toggle";
import { ThemeToggle } from "@/components/ui/theme-toggle";
import { Toaster } from "@/components/ui/toaster";
import { api } from "@/services/api";
import { useAuthStore } from "@/stores/auth-store";
import { useI18nStore } from "@/stores/i18n-store";
import type { TranslationKey } from "@/i18n/translations";
import { ArrowRightLeft, BarChart3, Boxes, ClipboardList, ClipboardPlus, FileClock, LayoutDashboard, LogOut, Menu, Package, Receipt, Settings, ShoppingCart, Users, X, type LucideIcon } from "lucide-react";
import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import type { ReactNode } from "react";
import { useEffect, useState } from "react";

type NavLink = {
  href: string;
  labelKey: TranslationKey;
  icon: LucideIcon;
};

const ownerLinks: NavLink[] = [
  { href: "/dashboard", labelKey: "nav.dashboard", icon: LayoutDashboard },
  { href: "/products", labelKey: "nav.products", icon: Package },
  { href: "/categories", labelKey: "nav.categories", icon: ClipboardPlus },
  { href: "/inventory", labelKey: "nav.inventory", icon: Boxes },
  { href: "/all-stock", labelKey: "nav.allStock", icon: ClipboardList },
  { href: "/stock-movements", labelKey: "nav.movements", icon: ArrowRightLeft },
  { href: "/transfers", labelKey: "nav.transfers", icon: Receipt },
  { href: "/employees", labelKey: "nav.employees", icon: Users },
  { href: "/reports", labelKey: "nav.reports", icon: BarChart3 },
  { href: "/audit-logs", labelKey: "nav.auditLogs", icon: FileClock },
  { href: "/settings", labelKey: "nav.settings", icon: Settings }
];

const employeeLinks: NavLink[] = [
  { href: "/pos", labelKey: "nav.pos", icon: ShoppingCart },
  { href: "/my-sales", labelKey: "nav.mySales", icon: Receipt },
  { href: "/branch-inventory", labelKey: "nav.branchStock", icon: Boxes }
];

const managerLinks: NavLink[] = [
  { href: "/dashboard", labelKey: "nav.dashboard", icon: LayoutDashboard },
  { href: "/pos", labelKey: "nav.pos", icon: ShoppingCart },
  { href: "/my-sales", labelKey: "nav.mySales", icon: Receipt },
  { href: "/branch-sales", labelKey: "nav.branchSales", icon: BarChart3 },
  { href: "/branch-inventory", labelKey: "nav.branchStock", icon: Boxes },
  { href: "/stock-receive", labelKey: "nav.receiveStock", icon: ClipboardPlus },
  { href: "/employees", labelKey: "nav.employees", icon: Users }
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
  const hydrateLanguage = useI18nStore((state) => state.hydrate);
  const language = useI18nStore((state) => state.language);
  const t = useI18nStore((state) => state.t);
  const [mounted, setMounted] = useState(false);
  const [logoutModalOpen, setLogoutModalOpen] = useState(false);
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const links = user?.role === "EMPLOYEE" ? employeeLinks : user?.role === "MANAGER" ? managerLinks : ownerLinks;
  const canAccessCurrentPage = Boolean(user && links.some((link) => link.href === path));
  const currentLink = links.find((link) => link.href === path);
  const quickLinks = links.length > 4 ? links.slice(0, 4) : links;

  useEffect(() => {
    hydrate();
    hydrateLanguage();
    setMounted(true);
  }, [hydrate, hydrateLanguage]);

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
        {t("common.redirectLogin")}
      </main>
    );
  }

  if (!canAccessCurrentPage) {
    return (
      <main className="grid min-h-screen place-items-center bg-field px-4 text-sm font-medium text-slate-600">
        {t("common.redirectAllowed")}
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
        {t(link.labelKey)}
      </Link>
    );
  });

  const mobileDrawerItems = links.map((link) => {
    const Icon = link.icon;
    const active = path === link.href;
    return (
      <Link
        key={link.href}
        href={link.href}
        onClick={() => setMobileMenuOpen(false)}
        className={`flex min-h-12 items-center gap-3 rounded-md px-3 text-sm font-semibold transition ${
          active ? "bg-brand text-white shadow-sm" : "text-slate-600 hover:bg-brandSoft hover:text-brand"
        }`}
      >
        <Icon className="h-5 w-5" />
        {t(link.labelKey)}
      </Link>
    );
  });

  const quickNavItems = quickLinks.map((link) => {
    const Icon = link.icon;
    const active = path === link.href;
    return (
      <Link
        key={link.href}
        href={link.href}
        className={`flex min-w-0 flex-1 flex-col items-center justify-center gap-1 rounded-md px-1 py-2 text-[11px] font-semibold transition ${
          active ? "bg-brand text-white shadow-sm" : "text-slate-600 hover:bg-brandSoft hover:text-brand"
        }`}
      >
        <Icon className="h-5 w-5 shrink-0" />
        <span className="max-w-full truncate">{t(link.labelKey)}</span>
      </Link>
    );
  });

  return (
    <div key={language} className="min-h-screen bg-field lg:flex">
      <Toaster />
      <aside className="hidden border-r border-line bg-white lg:sticky lg:top-0 lg:flex lg:h-screen lg:w-72 lg:flex-col">
        <div className="flex min-h-20 items-center border-b border-line px-5">
          <div>
            <div className="text-lg font-bold tracking-tight">{t("app.name")}</div>
            <div className="mt-1 text-xs font-medium text-slate-500">{user?.name ?? "Not signed in"} · {user?.role ?? "GUEST"}</div>
          </div>
        </div>
        <nav className="space-y-1.5 p-3">
          {navItems}
        </nav>
        <div className="mt-auto space-y-2 border-t border-line p-3">
          <LanguageToggle />
          <ThemeToggle />
          <Button className="w-full bg-slate-800 hover:bg-slate-700" onClick={() => setLogoutModalOpen(true)}>
            <LogOut className="h-4 w-4" />
            {t("common.logout")}
          </Button>
        </div>
      </aside>
      <div className="min-w-0 flex-1">
        <header className="sticky top-0 z-20 border-b border-line bg-white/95 backdrop-blur lg:hidden">
          <div className="flex min-h-16 items-center justify-between gap-2 px-3">
            <button
              className="grid h-10 w-10 shrink-0 place-items-center rounded-md border border-line bg-white text-slate-700"
              onClick={() => setMobileMenuOpen(true)}
              aria-label={t("shell.openMenu")}
            >
              <Menu className="h-5 w-5" />
            </button>
            <div className="min-w-0 flex-1">
              <div className="truncate text-base font-bold tracking-tight">{currentLink ? t(currentLink.labelKey) : t("app.name")}</div>
              <div className="truncate text-xs text-slate-500">{user?.name ?? "Not signed in"} · {user?.role ?? "GUEST"}</div>
            </div>
            <div className="flex shrink-0 items-center gap-2">
              <ThemeToggle compact />
              <LanguageToggle compact />
              <Button className="h-10 w-10 bg-slate-800 px-0 hover:bg-slate-700" onClick={() => setLogoutModalOpen(true)} title="Logout">
                <LogOut className="h-4 w-4" />
              </Button>
            </div>
          </div>
        </header>
        <main className="min-w-0 px-3 py-4 pb-24 sm:px-5 lg:p-7">{children}</main>
        <nav className={`fixed inset-x-0 bottom-0 z-30 gap-1 border-t border-line bg-white/95 p-2 shadow-[0_-8px_20px_rgba(15,23,42,0.08)] backdrop-blur lg:hidden ${links.length > 4 ? "grid grid-cols-5" : "flex"}`}>
          {quickNavItems}
          {links.length > 4 ? (
            <button
              className="flex min-w-0 flex-1 flex-col items-center justify-center gap-1 rounded-md px-1 py-2 text-[11px] font-semibold text-slate-600 hover:bg-brandSoft hover:text-brand"
              onClick={() => setMobileMenuOpen(true)}
            >
              <Menu className="h-5 w-5 shrink-0" />
              <span>{t("common.more")}</span>
            </button>
          ) : null}
        </nav>
      </div>
      {mobileMenuOpen ? (
        <div className="fixed inset-0 z-50 bg-slate-950/50 backdrop-blur-sm lg:hidden">
          <div className="flex h-full w-[86vw] max-w-sm flex-col border-r border-line bg-white shadow-2xl">
            <div className="flex min-h-16 items-center justify-between gap-3 border-b border-line px-4">
              <div className="min-w-0">
                <div className="truncate text-base font-bold">{t("app.name")}</div>
                <div className="truncate text-xs text-slate-500">{user?.name} · {user?.role}</div>
              </div>
              <button
                className="grid h-10 w-10 shrink-0 place-items-center rounded-md border border-line text-slate-600 hover:bg-field"
                onClick={() => setMobileMenuOpen(false)}
                aria-label={t("shell.closeMenu")}
              >
                <X className="h-5 w-5" />
              </button>
            </div>
            <nav className="min-h-0 flex-1 space-y-1 overflow-y-auto p-3">
              {mobileDrawerItems}
            </nav>
            <div className="space-y-2 border-t border-line p-3">
              <LanguageToggle />
              <ThemeToggle />
              <Button className="w-full bg-slate-800 hover:bg-slate-700" onClick={() => setLogoutModalOpen(true)}>
                <LogOut className="h-4 w-4" />
                {t("common.logout")}
              </Button>
            </div>
          </div>
        </div>
      ) : null}
      <ConfirmModal
        open={logoutModalOpen}
        title={t("logout.title")}
        description={t("logout.description")}
        confirmLabel={t("common.logout")}
        danger
        onCancel={() => setLogoutModalOpen(false)}
        onConfirm={logout}
      />
    </div>
  );
}
