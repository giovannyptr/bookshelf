import axios from "axios";

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE, // e.g. http://localhost:8080
});

// Attach token from localStorage on every request
api.interceptors.request.use((config) => {
  const token = localStorage.getItem("token");
  if (token) config.headers.Authorization = `Bearer ${token}`;
  return config;
});

// Optional: auto-logout on 401
api.interceptors.response.use(
  (res) => res,
  (err) => {
    if (err.response && err.response.status === 401) {
      localStorage.removeItem("token");
      localStorage.removeItem("user");
      // You can also redirect to /login here if you want
      // window.location.href = "/login";
    }
    return Promise.reject(err);
  }
);

export default api;
