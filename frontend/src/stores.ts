import { writable, type Writable } from "svelte/store";
import type { Toast } from "./types";

export const toasts: Writable<Toast[]> = writable([]);
