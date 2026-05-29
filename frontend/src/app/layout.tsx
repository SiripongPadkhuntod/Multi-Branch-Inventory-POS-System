import "./globals.css";
import type { Metadata } from "next";
import type { ReactNode } from "react";

const themeScript = `
try {
  var theme = localStorage.getItem("theme");
  if (!theme && window.matchMedia("(prefers-color-scheme: dark)").matches) theme = "dark";
  if (theme === "dark") document.documentElement.classList.add("dark");
} catch (_) {}
`;

export const metadata: Metadata = {
  title: "Multi-Branch POS",
  description: "Inventory and POS management"
};

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body suppressHydrationWarning>
        <script dangerouslySetInnerHTML={{ __html: themeScript }} />
        {children}
      </body>
    </html>
  );
}
