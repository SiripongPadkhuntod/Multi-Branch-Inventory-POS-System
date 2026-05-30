"use client";

import { useI18nStore } from "@/stores/i18n-store";
import { Languages } from "lucide-react";
import { useEffect, useState } from "react";

export function LanguageToggle({ compact = false }: { compact?: boolean }) {
  const language = useI18nStore((state) => state.language);
  const hydrate = useI18nStore((state) => state.hydrate);
  const toggleLanguage = useI18nStore((state) => state.toggleLanguage);
  const t = useI18nStore((state) => state.t);
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    hydrate();
    setMounted(true);
  }, [hydrate]);

  const label = mounted ? language.toUpperCase() : "EN";

  return (
    <button
      type="button"
      className="inline-flex h-10 items-center justify-center gap-2 rounded-md border border-line bg-white px-3 text-sm font-semibold text-slate-700 transition hover:bg-field dark:bg-slate-900 dark:text-slate-100"
      onClick={toggleLanguage}
      aria-label={t("language.switchTo")}
      title={t("language.switchTo")}
    >
      <Languages className="h-4 w-4" />
      {compact ? <span className="text-xs">{label}</span> : <span>{label}</span>}
    </button>
  );
}
