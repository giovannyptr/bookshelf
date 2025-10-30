import { createRouter, createWebHistory } from "vue-router";
import BooksPage from "./views/BooksPage.vue";
import BookDetail from "./views/BookDetail.vue";
import LoginPage from "./views/LoginPage.vue";
import { useAuth } from "./lib/auth";

const routes = [
  { path: "/", redirect: "/books" },
  { path: "/login", component: LoginPage, meta: { guest: true } },
  { path: "/books", component: BooksPage },
  { path: "/books/:id", component: BookDetail, props: true },
];

export const router = createRouter({
  history: createWebHistory(),
  routes,
});

// keep it simple: if already logged in, don't show /login
router.beforeEach((to, _from, next) => {
  const { isAuthed } = useAuth();
  if (to.meta.guest && isAuthed.value) return next("/books");
  return next();
});
