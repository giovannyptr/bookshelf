<script setup>
import { ref, onMounted } from "vue";
import api from "../lib/api";
import { useAuth } from "../lib/auth";
import { CATEGORY_OPTIONS } from "../lib/constants";
import { formatIDR } from "../lib/format";
import BookForm from "../components/BookForm.vue";

const { isAuthed } = useAuth();

// --- URL helper to fix broken cover paths ---
const API_BASE = (import.meta.env.VITE_API_BASE || "").replace(/\/$/, "");
const fullUrl = (p) => (!p ? "" : p.startsWith("http") ? p : `${API_BASE}${p}`);

// --- query state ---
const q = ref("");
const category = ref("");
const page = ref(1);
const limit = ref(10);
const total = ref(0);
const items = ref([]);
const loading = ref(false);
const error = ref("");

// --- create form model (used by BookForm) ---
const createModel = ref({
  title: "",
  author: "",
  category: "",
  price: "",
  stock: ""
});

async function fetchBooks() {
  loading.value = true; error.value = "";
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

async function createBook(fd) {
  try {
    await api.post("/books", fd, { headers: { "Content-Type": "multipart/form-data" } });
    createModel.value = { title:"", author:"", category:"", price:"", stock:"" };
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
      <select v-model="category" class="input w200">
        <option value="">All categories</option>
        <option v-for="opt in CATEGORY_OPTIONS" :key="opt" :value="opt">{{ opt }}</option>
      </select>
      <button @click="fetchBooks" class="btn">Search</button>
    </div>

    <!-- Create form (only if logged in) -->
    <details v-if="isAuthed" class="card">
      <summary class="summary">+ Add new book</summary>
      <BookForm
        v-model="createModel"
        submit-label="Create"
        @submit="createBook"
      />
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
            <img v-if="b.coverUrl" :src="fullUrl(b.coverUrl)" alt="" class="cover" />
          </td>
          <td><router-link :to="`/books/${b.id}`">{{ b.title }}</router-link></td>
          <td>{{ b.author }}</td>
          <td>{{ b.category }}</td>
          <td class="right">{{ formatIDR(b.price) }}</td>
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
.input { padding: 8px; border: 1px solid var(--line,#ddd); border-radius: 6px; background: var(--bg,white); color: var(--fg,#111);}
.flex { flex: 1; }
.w200 { width: 200px; }
.btn { padding: 8px 12px; border: 1px solid var(--line,#ddd); border-radius: 6px; background: var(--bg,white); color: var(--fg,#111); cursor: pointer; text-decoration: none; }
.btn:disabled { opacity: .5; cursor: not-allowed; }
.btn.danger { border-color: #ff9c9c; background: #fff3f3; }
.toolbar { display: flex; gap: 8px; align-items: center; margin-bottom: 16px; }
.card { margin: 12px 0; }
.summary { cursor: pointer; font-weight: 600; }
.table { width: 100%; border-collapse: collapse; }
.table th, .table td { border-bottom: 1px solid var(--line,#eee); padding: 8px; font-size: 14px; text-align: left; }
.table th.right, .table td.right { text-align: right; }
.cover { height: 48px; width: 48px; object-fit: cover; border: 1px solid var(--line,#eee); border-radius: 4px; }
.pager { display: flex; gap: 8px; align-items: center; justify-content: flex-end; margin-top: 12px; }
.muted { color: #666; }
.error { color: #b00020; margin: 8px 0; }
.pad12 { padding: 12px; }
</style>
