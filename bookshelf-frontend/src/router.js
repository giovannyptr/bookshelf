import { createRouter, createWebHistory } from "vue-router";
import BooksPage from "./views/BooksPage.vue";
import BookDetail from "./views/BookDetail.vue"; 

const routes = [
  { path: "/", redirect: "/books" },
  { path: "/books", component: BooksPage },
  { path: "/books/:id", component: BookDetail, props: true }, 
];

export const router = createRouter({
  history: createWebHistory(),
  routes,
});
