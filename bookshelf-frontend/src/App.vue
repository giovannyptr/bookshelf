<template>
  <div style="max-width: 960px; margin: 0 auto; padding: 0 16px;">
    <AppHeader />
    <router-view />
  </div>
</template>

<script setup>
import AppHeader from "./components/AppHeader.vue";
</script>
