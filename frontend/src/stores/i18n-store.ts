"use client";

import { translations, type Language, type TranslationKey } from "@/i18n/translations";
import { create } from "zustand";

type I18nState = {
  language: Language;
  hydrate: () => void;
  setLanguage: (language: Language) => void;
  toggleLanguage: () => void;
  t: (key: TranslationKey) => string;
};

function storedLanguage(): Language {
  if (typeof window === "undefined") {
    return "en";
  }
  const value = localStorage.getItem("language");
  return value === "th" ? "th" : "en";
}

function applyLanguage(language: Language) {
  if (typeof document !== "undefined") {
    document.documentElement.lang = language;
  }
  if (typeof window !== "undefined") {
    localStorage.setItem("language", language);
  }
}

export const useI18nStore = create<I18nState>((set, get) => ({
  language: "en",
  hydrate: () => {
    const language = storedLanguage();
    applyLanguage(language);
    set({ language });
  },
  setLanguage: (language) => {
    applyLanguage(language);
    set({ language });
  },
  toggleLanguage: () => {
    const next = get().language === "en" ? "th" : "en";
    applyLanguage(next);
    set({ language: next });
  },
  t: (key) => translations[get().language][key] ?? translations.en[key] ?? key
}));
