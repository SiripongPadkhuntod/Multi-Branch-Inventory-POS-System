"use client";

import type { User } from "@/types/domain";
import { create } from "zustand";

type AuthState = {
	user: User | null;
	hydrate: () => void;
	setSession: (user: User, accessToken: string) => void;
	clear: () => void;
};

function storedUser() {
	if (typeof window === "undefined") {
		return null;
	}
	const value = localStorage.getItem("user");
	if (!value) {
		return null;
	}
	try {
		return JSON.parse(value) as User;
	} catch {
		localStorage.removeItem("user");
		return null;
	}
}

export const useAuthStore = create<AuthState>((set) => ({
	user: null,
	hydrate: () => set({ user: storedUser() }),
	setSession: (user, accessToken) => {
		localStorage.setItem("access_token", accessToken);
		localStorage.setItem("user", JSON.stringify(user));
		set({ user });
	},
	clear: () => {
		localStorage.removeItem("access_token");
		localStorage.removeItem("refresh_token");
		localStorage.removeItem("user");
		set({ user: null });
	}
}));
