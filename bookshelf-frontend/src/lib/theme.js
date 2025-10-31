import { reactive, computed } from "vue";

const STORAGE_KEY = "theme"; // 'light' | 'dark' | 'system'

const state = reactive({
  choice: localStorage.getItem(STORAGE_KEY) || "system", // user choice
});

function systemPrefersDark() {
  if (typeof window === "undefined" || !window.matchMedia) return false;
  return window.matchMedia("(prefers-color-scheme: dark)").matches;
}

function effectiveTheme() {
  return state.choice === "system" ? (systemPrefersDark() ? "dark" : "light") : state.choice;
}

function applyTheme() {
  const t = effectiveTheme();
  // put theme on <html data-theme="dark|light">
  document.documentElement.setAttribute("data-theme", t);
}

function setTheme(choice) {
  state.choice = choice; // 'light' | 'dark' | 'system'
  localStorage.setItem(STORAGE_KEY, choice);
  applyTheme();
}

function initTheme() {
  // react to system changes when in "system"
  if (typeof window !== "undefined" && window.matchMedia) {
    const mq = window.matchMedia("(prefers-color-scheme: dark)");
    const handler = () => { if (state.choice === "system") applyTheme(); };
    if (mq.addEventListener) mq.addEventListener("change", handler);
    else mq.addListener(handler); // Safari old
  }
  applyTheme();
}

export function useTheme() {
  const current = computed(() => state.choice);          // 'light' | 'dark' | 'system'
  const mode = computed(() => effectiveTheme());         // 'light' | 'dark'
  return { current, mode, setTheme, initTheme };
}
