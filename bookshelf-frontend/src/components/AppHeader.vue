<script setup>
import { useRouter } from "vue-router";
import { useAuth } from "../lib/auth";
import { useTheme } from "../lib/theme";

const auth = useAuth();
const router = useRouter();
const { current, mode, setTheme } = useTheme();

function onLogout() {
  auth.logout();
  router.push("/login");
}

function cycleTheme() {
  // cycles: system -> light -> dark -> system
  const order = ["system", "light", "dark"];
  const idx = order.indexOf(current.value);
  const next = order[(idx + 1) % order.length];
  setTheme(next);
}

function themeLabel() {
  if (current.value === "system") return `System (${mode.value})`;
  return current.value.charAt(0).toUpperCase() + current.value.slice(1);
}
</script>

<template>
  <header class="site-header">
    <div class="wrap">
      <h1 class="brand">📚 Bookshelf</h1>

      <nav class="right">
        <button class="btn" @click="cycleTheme" title="Toggle theme">
          🌓 {{ themeLabel() }}
        </button>

        <template v-if="auth.isAuthed.value">
          <span class="hello">Hello, {{ auth.state.user?.name || 'Reader' }} 👋</span>
          <button class="btn btn-danger" @click="onLogout">Logout</button>
        </template>
        <template v-else>
          <router-link class="btn" to="/login">Login</router-link>
        </template>
      </nav>
    </div>
  </header>
</template>

<style scoped>
.site-header {
  position: sticky; top: 0; z-index: 50;
  background: var(--card);
  border-bottom: 1px solid var(--border);
}
.wrap {
  max-width: 1100px;
  margin: 0 auto;
  padding: 14px 20px;
  display: flex; align-items: center; justify-content: space-between; gap: 16px;
}
.brand { margin: 0; font-size: 24px; font-weight: 700; letter-spacing: .2px; }
.right { display: flex; align-items: center; gap: 12px; }
.hello { color: var(--muted); font-weight: 500; }
.btn {
  padding: 8px 12px; border-radius: 8px;
  background: var(--btn-bg); color: var(--text); border: 1px solid var(--btn-border);
  cursor: pointer; text-decoration: none; transition: background .15s ease, border-color .15s ease;
}
.btn:hover { filter: brightness(1.05); }
.btn-danger { background: var(--btn-danger-bg); border-color: var(--btn-danger-border); }
</style>
