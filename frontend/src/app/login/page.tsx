"use client";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { LanguageToggle } from "@/components/ui/language-toggle";
import { ThemeToggle } from "@/components/ui/theme-toggle";
import { api } from "@/services/api";
import { useAuthStore } from "@/stores/auth-store";
import { useI18nStore } from "@/stores/i18n-store";
import { zodResolver } from "@hookform/resolvers/zod";
import { LockKeyhole, Mail, Store } from "lucide-react";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

const schema = z.object({
  email: z.string().email(),
  password: z.string().min(8)
});

type LoginForm = z.infer<typeof schema>;

export default function LoginPage() {
  const router = useRouter();
  const setSession = useAuthStore((state) => state.setSession);
  const hydrateLanguage = useI18nStore((state) => state.hydrate);
  const language = useI18nStore((state) => state.language);
  const t = useI18nStore((state) => state.t);
  const [error, setError] = useState("");
  const form = useForm<LoginForm>({
    resolver: zodResolver(schema),
    defaultValues: { email: "owner@example.com", password: "password123" }
  });

  useEffect(() => {
    hydrateLanguage();
    if (window.location.search.includes("session=expired")) {
      setError(t("login.expired"));
      window.history.replaceState(null, "", "/login");
    }
  }, [hydrateLanguage, t]);

  async function onSubmit(values: LoginForm) {
    setError("");
    try {
      const session = await api.login(values.email, values.password);
      setSession(session.user, session.access_token);
      router.push(session.user.role === "EMPLOYEE" ? "/pos" : "/dashboard");
    } catch (err) {
      setError(err instanceof Error ? err.message : t("login.failed"));
    }
  }

  return (
    <main key={language} className="grid min-h-screen place-items-center bg-field px-4">
      <div className="fixed right-4 top-4 flex gap-2">
        <LanguageToggle compact />
        <ThemeToggle compact />
      </div>
      <form onSubmit={form.handleSubmit(onSubmit)} className="w-full max-w-sm rounded-md border border-line bg-white p-6 shadow-xl shadow-slate-200/70">
        <div className="mb-6">
          <div className="mb-4 grid h-12 w-12 place-items-center rounded-md bg-brandSoft text-brand">
            <Store className="h-6 w-6" />
          </div>
          <h1 className="text-2xl font-bold tracking-tight">{t("app.name")}</h1>
          <p className="mt-1 text-sm text-slate-500">{t("login.subtitle")}</p>
        </div>
        <label className="mb-4 block text-sm font-medium">
          {t("login.email")}
          <div className="mt-1 flex items-center gap-2">
            <Mail className="h-4 w-4 text-slate-500" />
            <Input {...form.register("email")} />
          </div>
        </label>
        <label className="mb-4 block text-sm font-medium">
          {t("login.password")}
          <div className="mt-1 flex items-center gap-2">
            <LockKeyhole className="h-4 w-4 text-slate-500" />
            <Input type="password" {...form.register("password")} />
          </div>
        </label>
        {error ? <p className="mb-4 rounded-md bg-red-50 p-3 text-sm text-red-700">{error}</p> : null}
        <Button className="w-full" disabled={form.formState.isSubmitting}>
          {t("login.signIn")}
        </Button>
      </form>
    </main>
  );
}
