<script setup>
import AppHeader from "./components/AppHeader.vue";
import { useTheme } from "./lib/theme";

const { initTheme } = useTheme();
initTheme(); // apply + watch system changes
</script>

<template>
  <AppHeader />
  <main class="page">
    <router-view />
  </main>
</template>

<style>
.page {
  max-width: 1100px;
  margin: 24px auto;
  padding: 0 20px;
}
</style>
