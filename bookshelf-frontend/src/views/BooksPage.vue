<script setup>
import { ref, onMounted } from "vue";
import api from "../lib/api";
import { useAuth } from "../lib/auth";

const { isAuthed } = useAuth();
const API_BASE = import.meta.env.VITE_API_BASE || "";

// --- query state ---
const q = ref("");
const category = ref("");
const page = ref(1);
const limit = ref(10);
const total = ref(0);
const items = ref([]);
const loading = ref(false);
const error = ref("");

// --- create form state (authed only) ---
const form = ref({ title: "", author: "", category: "", price: "", stock: "" });
const cover = ref(null);

// ---- actions ----
async function fetchBooks() {
  loading.value = true;
  error.value = "";
  try {
    const { data } = await api.get("/books", {
      params: { q: q.value, category: category.value, page: page.value, limit: limit.value },
    });
    const payload = data.data ?? data;
    items.value = payload.items ?? [];
    total.value = Number(payload.total ?? 0);
  } catch (e) {
    error.value = e?.response?.data?.error || e.message;
  } finally {
    loading.value = false;
  }
}

function nextPage() {
  if (page.value * limit.value < total.value) {
    page.value++;
    fetchBooks();
  }
}
function prevPage() {
  if (page.value > 1) {
    page.value--;
    fetchBooks();
  }
}

async function createBook() {
  const fd = new FormData();
  Object.entries(form.value).forEach(([k, v]) => fd.append(k, v));
  if (cover.value?.files?.[0]) fd.append("cover", cover.value.files[0]);

  try {
    await api.post("/books", fd, { headers: { "Content-Type": "multipart/form-data" } });
    form.value = { title: "", author: "", category: "", price: "", stock: "" };
    if (cover.value) cover.value.value = "";
    await fetchBooks();
  } catch (e) {
    alert(e?.response?.data?.error || e.message);
  }
}

async function removeBook(id) {
  if (!confirm("Delete this book?")) return;
  try {
    await api.delete(`/books/${id}`);
    await fetchBooks();
  } catch (e) {
    alert(e?.response?.data?.error || e.message);
  }
}

onMounted(fetchBooks);
</script>

<template>
  <div>
    <!-- Search / filter -->
    <div class="toolbar">
      <input v-model="q" placeholder="Search title/author..." class="input flex" />
      <input v-model="category" placeholder="Category" class="input w200" />
      <button @click="fetchBooks" class="btn">Search</button>
    </div>

    <!-- Create form (only if logged in) -->
    <details v-if="isAuthed" class="card">
      <summary class="summary">+ Add new book</summary>
      <div class="grid">
        <input v-model="form.title" placeholder="Title" class="input" />
        <input v-model="form.author" placeholder="Author" class="input" />
        <input v-model="form.category" placeholder="Category" class="input" />
        <input v-model="form.price" type="number" step="0.01" placeholder="Price" class="input" />
        <input v-model="form.stock" type="number" placeholder="Stock" class="input" />
        <input ref="cover" type="file" accept="image/*" class="input" />
      </div>
      <button @click="createBook" class="btn mt8">Create</button>
    </details>

    <div v-if="error" class="error">{{ error }}</div>
    <div v-if="loading" class="muted">Loading…</div>

    <!-- List -->
    <table class="table">
      <thead>
        <tr>
          <th>Cover</th>
          <th>Title</th>
          <th>Author</th>
          <th>Category</th>
          <th class="right">Price</th>
          <th class="right">Stock</th>
          <th class="right"></th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="b in items" :key="b.id">
          <td>
            <img v-if="b.coverUrl" :src="API_BASE + b.coverUrl" alt="" class="cover" />
          </td>
          <td>
            <router-link :to="`/books/${b.id}`">{{ b.title }}</router-link>
          </td>
          <td>{{ b.author }}</td>
          <td>{{ b.category }}</td>
          <td class="right">{{ Number(b.price || 0).toFixed(2) }}</td>
          <td class="right">{{ b.stock }}</td>
          <td class="right">
            <router-link :to="`/books/${b.id}`" class="btn">Detail</router-link>
            <button v-if="isAuthed" @click="removeBook(b.id)" class="btn danger">Delete</button>
          </td>
        </tr>
        <tr v-if="!loading && items.length === 0">
          <td colspan="7" class="muted pad12">No books.</td>
        </tr>
      </tbody>
    </table>

    <!-- Pagination -->
    <div class="pager">
      <button @click="prevPage" :disabled="page === 1" class="btn">Prev</button>
      <span>Page {{ page }}</span>
      <button @click="nextPage" :disabled="page * limit >= total" class="btn">Next</button>
    </div>
  </div>
</template>

<style scoped>
.input { padding: 8px; border: 1px solid #ddd; border-radius: 6px; }
.flex { flex: 1; }
.w200 { width: 200px; }
.btn { padding: 8px 12px; border: 1px solid #ddd; border-radius: 6px; background: white; cursor: pointer; text-decoration: none; }
.btn:disabled { opacity: .5; cursor: not-allowed; }
.btn.danger { border-color: #ff9c9c; background: #fff3f3; }
.toolbar { display: flex; gap: 8px; align-items: center; margin-bottom: 16px; }
.card { margin: 12px 0; }
.summary { cursor: pointer; font-weight: 600; }
.grid { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 8px; margin-top: 8px; }
.mt8 { margin-top: 8px; }
.table { width: 100%; border-collapse: collapse; }
.table th, .table td { border-bottom: 1px solid #eee; padding: 8px; font-size: 14px; text-align: left; }
.table th.right, .table td.right { text-align: right; }
.cover { height: 48px; width: 48px; object-fit: cover; border: 1px solid #eee; border-radius: 4px; }
.pager { display: flex; gap: 8px; align-items: center; justify-content: flex-end; margin-top: 12px; }
.muted { color: #666; }
.error { color: #b00020; margin: 8px 0; }
.pad12 { padding: 12px; }
</style>
