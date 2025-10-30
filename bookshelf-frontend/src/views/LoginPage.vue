<script setup>
import { ref } from "vue";
import { useRouter } from "vue-router";
import api from "../lib/api";
import { useAuth } from "../lib/auth";

const email = ref("");
const password = ref("");
const loading = ref(false);
const error = ref("");
const router = useRouter();
const { setAuth } = useAuth();

async function submit() {
  loading.value = true; error.value = "";
  try {
    const { data } = await api.post("/auth/login", { email: email.value, password: password.value });
    const payload = data.data ?? data; // supports both {ok,data:{...}} or flat
    setAuth(payload.token, payload.user);
    router.push("/books");
  } catch (e) {
    error.value = e?.response?.data?.error || e.message;
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <div style="max-width:420px; margin:32px auto;">
    <h2 style="margin:0 0 12px;">Login</h2>
    <p class="muted" style="margin:0 0 16px;">Use the account you created with <code>/auth/register</code> (Postman ok).</p>

    <div class="grid">
      <input v-model="email" type="email" placeholder="Email" class="input" />
      <input v-model="password" type="password" placeholder="Password" class="input" />
      <button :disabled="loading" class="btn primary" @click="submit">Login</button>
    </div>

    <div v-if="error" class="error">{{ error }}</div>
  </div>
</template>

<style scoped>
.grid { display:grid; gap:10px; }
.input { padding:10px; border:1px solid #ddd; border-radius:6px; }
.btn { padding:10px 12px; border:1px solid #ddd; border-radius:6px; background:#fff; cursor:pointer; }
.btn.primary { border-color:#bfe3ff; background:#f5faff; }
.error { color:#b00020; margin-top:10px; }
.muted { color:#666; }
</style>
