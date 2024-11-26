import { createRouter, createWebHistory } from "vue-router";
import LoginForm from "../components/LoginForm.vue";
import App from "../App.vue";

const routes = [
  {
    path: "/login",
    name: "Login",
    component: LoginForm,
  },
  {
    path: "/",
    name: "Home",
    component: App,
  },
  {
    path: "/:pathMatch(.*)*",
    redirect: "/login",
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem("token");

  const isTokenValid = token && token !== "null";

  if (!isTokenValid && to.name !== "Login") {
    next({ name: "Login" });
  } else if (isTokenValid && to.name === "Login") {
    next({ name: "Home" });
  } else {
    next();
  }
});

export default router;