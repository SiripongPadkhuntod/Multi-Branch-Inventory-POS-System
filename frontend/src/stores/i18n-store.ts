"use client";

import { translations, type Language } from "@/i18n/translations";
import { create } from "zustand";

type I18nState = {
  language: Language;
  hydrate: () => void;
  setLanguage: (language: Language) => void;
  toggleLanguage: () => void;
  t: (key: string) => string;
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

function createTranslator(language: Language) {
  return (key: string) => {
    const languageMap = translations[language] as Record<string, string>;
    const fallbackMap = translations.en as Record<string, string>;
    return languageMap[key] ?? fallbackMap[key] ?? key;
  };
}

export const useI18nStore = create<I18nState>((set, get) => ({
  language: "en",
  hydrate: () => {
    const language = storedLanguage();
    applyLanguage(language);
    if (get().language === language) {
      return;
    }
    set({ language, t: createTranslator(language) });
  },
  setLanguage: (language) => {
    applyLanguage(language);
    if (get().language === language) {
      return;
    }
    set({ language, t: createTranslator(language) });
  },
  toggleLanguage: () => {
    const next = get().language === "en" ? "th" : "en";
    applyLanguage(next);
    set({ language: next, t: createTranslator(next) });
  },
  t: createTranslator("en")
}));
