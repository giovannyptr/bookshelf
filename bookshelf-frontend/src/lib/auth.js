import { reactive, computed } from "vue";

const state = reactive({
  token: localStorage.getItem("token") || "",
  user: JSON.parse(localStorage.getItem("user") || "null"),
});

function setAuth(token, user) {
  state.token = token;
  state.user = user || null;
  localStorage.setItem("token", token);
  localStorage.setItem("user", JSON.stringify(user || null));
}

function logout() {
  state.token = "";
  state.user = null;
  localStorage.removeItem("token");
  localStorage.removeItem("user");
}

export function useAuth() {
  const isAuthed = computed(() => !!state.token);
  return { state, isAuthed, setAuth, logout };
}
