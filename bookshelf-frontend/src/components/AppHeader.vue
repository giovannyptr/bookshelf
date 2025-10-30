<script setup>
import { useAuth } from "../lib/auth";
import { useRouter } from "vue-router";

const { state, isAuthed, logout } = useAuth();
const router = useRouter();

function onLogout() {
  logout();
  router.push("/books");
}
</script>

<template>
  <header class="wrap">
    <router-link to="/books" class="brand">📚 Bookshelf by Gio</router-link>

    <nav class="nav">
      <router-link to="/books">Browse</router-link>
      <span class="sep">·</span>

      <template v-if="isAuthed">
        <span class="hello">Hello, {{ state.user?.name || state.user?.email || "User" }} 👋</span>
        <button class="btn" @click="onLogout">Logout</button>
      </template>
      <template v-else>
        <router-link to="/login" class="btn">Login</router-link>
      </template>
    </nav>
  </header>
</template>

<style scoped>
.wrap {
  display:flex; align-items:center; justify-content:space-between;
  gap:12px; padding:12px 16px; border-bottom:1px solid #dc9fed; margin-bottom:16px;
}
.brand { font-weight:800; text-decoration:none; color:#111; }
.nav { display:flex; align-items:center; gap:10px; }
.sep { color:#aaa; }
.btn { padding:6px 10px; border:1px solid #ddd; border-radius:6px; background:#ffffff; cursor:pointer; text-decoration:none; }
.hello { color:#333; }
</style>
