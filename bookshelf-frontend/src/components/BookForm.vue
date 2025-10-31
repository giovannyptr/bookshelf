<script setup>
import { ref, watch, computed } from "vue";
import { CATEGORY_OPTIONS } from "../lib/constants";

const props = defineProps({
  modelValue: {
    type: Object,
    default: () => ({ title: "", author: "", category: "", price: "", stock: "" }),
  },
  submitLabel: { type: String, default: "Save" },
  disabled: { type: Boolean, default: false },
  showCoverInput: { type: Boolean, default: true }, 
});

const emit = defineEmits(["submit", "update:modelValue"]);

const form = ref({ ...props.modelValue });
const coverRef = ref(null);

watch(
  () => props.modelValue,
  v => (form.value = { ...v }),
  { deep: true }
);

function update(k, v) {
  form.value[k] = v;
  emit("update:modelValue", { ...form.value });
}

function onSubmit() {
  const fd = new FormData();
  Object.entries(form.value).forEach(([k, v]) => fd.append(k, v ?? ""));
  if (props.showCoverInput && coverRef.value?.files?.[0]) {
    fd.append("cover", coverRef.value.files[0]);
  }
  emit("submit", fd);
}
</script>

<template>
  <div class="grid">
    <input class="input" placeholder="Title" :disabled="disabled"
           :value="form.title" @input="update('title', $event.target.value)" />
    <input class="input" placeholder="Author" :disabled="disabled"
           :value="form.author" @input="update('author', $event.target.value)" />

    <select class="input" :disabled="disabled"
            :value="form.category" @change="update('category', $event.target.value)">
      <option value="" disabled>Select category</option>
      <option v-for="opt in CATEGORY_OPTIONS" :key="opt" :value="opt">{{ opt }}</option>
    </select>

    <input class="input" type="number" step="0.01" placeholder="Price" :disabled="disabled"
           :value="form.price" @input="update('price', $event.target.value)" />
    <input class="input" type="number" placeholder="Stock" :disabled="disabled"
           :value="form.stock" @input="update('stock', $event.target.value)" />

    <input v-if="showCoverInput" ref="coverRef" class="input" type="file" accept="image/*" :disabled="disabled" />
  </div>

  <button class="btn mt8" :disabled="disabled" @click="onSubmit">{{ submitLabel }}</button>
</template>

<style scoped>
.grid { display:grid; grid-template-columns: repeat(2, minmax(0,1fr)); gap:8px; }
.input { padding:8px; border:1px solid #ddd; border-radius:6px; }
.btn { padding:8px 12px; border:1px solid #ddd; border-radius:6px; background:#fff; cursor:pointer; }
.btn:disabled { opacity:.5; cursor:not-allowed }
.mt8 { margin-top:8px; }
</style>
