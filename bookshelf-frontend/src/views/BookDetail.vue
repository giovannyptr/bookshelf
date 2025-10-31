<script setup>
import { ref, onMounted, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import api from "../lib/api";
import { useAuth } from "../lib/auth";
import BookForm from "../components/BookForm.vue";
import { formatIDR } from "../lib/format"

const { isAuthed } = useAuth();
const API_BASE = import.meta.env.VITE_API_BASE || "";
const route = useRoute();
const router = useRouter();

const id = ref(route.params.id);
watch(() => route.params.id, v => { id.value = v; fetchBook(); });

const book = ref(null);
const loading = ref(false);
const error = ref("");

const editModel = ref({ title:"", author:"", category:"", price:"", stock:"" });

async function fetchBook() {
  loading.value = true; error.value = "";
  try {
    const { data } = await api.get(`/books/${id.value}`);
    const payload = data.data ?? data;
    book.value = payload;
    editModel.value = {
      title: payload.title ?? "",
      author: payload.author ?? "",
      category: payload.category ?? "",
      price: payload.price ?? "",
      stock: payload.stock ?? "",
    };
  } catch (e) {
    error.value = e?.response?.data?.error || e.message;
  } finally {
    loading.value = false;
  }
}

async function save(fd) {
  try {
    await api.put(`/books/${id.value}`, fd, { headers: { "Content-Type": "multipart/form-data" } });
    await fetchBook();
    alert("Saved!");
  } catch (e) {
    alert(e?.response?.data?.error || e.message);
  }
}

async function removeBook() {
  if (!confirm("Delete this book?")) return;
  try {
    await api.delete(`/books/${id.value}`);
    router.push("/books");
  } catch (e) {
    alert(e?.response?.data?.error || e.message);
  }
}

onMounted(fetchBook);
</script>

<template>
  <div>
    <button class="btn" @click="$router.back()">← Back</button>

    <div v-if="loading" class="muted" style="margin-top:8px;">Loading…</div>
    <div v-if="error" class="error">{{ error }}</div>

    <div v-if="book" class="wrap">
      <div class="left">
        <img v-if="book.coverUrl" :src="API_BASE + book.coverUrl" alt="" class="cover" />
        <div v-else class="placeholder">No cover</div>
      </div>

      <div class="right">
        <h2 style="margin:0 0 8px;">{{ book.title }}</h2>
        <div class="muted" style="margin-bottom:12px;">
          by <strong>{{ book.author || "Unknown" }}</strong> — <em>{{ book.category || "Uncategorized" }}</em>
        </div>

        <div class="meta">
          <div>Price: <strong>{{ formatIDR(book.price) }}</strong></div>
          <div>Stock: <strong>{{ book.stock }}</strong></div>
          <div>ID: <code>{{ book.id }}</code></div>
        </div>

        <hr style="margin:16px 0;">

        <template v-if="isAuthed">
          <h3 style="margin:0 0 8px;">Edit</h3>
          <BookForm
            v-model="editModel"
            submit-label="Save"
            @submit="save"
          />
          <div class="row">
            <button class="btn danger" @click="removeBook">Delete</button>
          </div>
        </template>
        <template v-else>
          <p class="muted">Login to edit or delete this book.</p>
        </template>
      </div>
    </div>
  </div>
</template>

<style scoped>
.wrap { display: grid; grid-template-columns: 180px 1fr; gap: 16px; margin-top: 12px; }
.left { display:flex; align-items:flex-start; justify-content:center; }
.cover { width: 160px; height: 220px; object-fit: cover; border: 1px solid #eee; border-radius: 6px; }
.placeholder { width:160px; height:220px; border:1px dashed #ccc; border-radius:6px; display:flex; align-items:center; justify-content:center; color:#888; }
.row { display:flex; gap:8px; margin-top:8px; }
.btn { padding:8px 12px; border:1px solid #ddd; border-radius:6px; background:white; cursor:pointer; }
.btn.danger { border-color:#ffb8b8; background:#fff4f4; }
.muted { color:#666; }
.error { color:#b00020; margin:8px 0; }
.meta { display:flex; gap:16px; flex-wrap:wrap; color:#333; }
</style>
