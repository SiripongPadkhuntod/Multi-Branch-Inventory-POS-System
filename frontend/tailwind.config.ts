import type { Config } from "tailwindcss";

const config: Config = {
  darkMode: "class",
  content: ["./src/**/*.{ts,tsx}"],
  theme: {
    extend: {
      colors: {
        ink: "#1f2933",
        field: "#f4f6f8",
        panel: "#ffffff",
        line: "#d8dee6",
        brand: "#0f766e",
        brandSoft: "#e6f4f1",
        accent: "#b45309",
        danger: "#dc2626"
      }
    }
  },
  plugins: []
};

export default config;
