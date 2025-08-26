import { createRouter, createWebHistory } from "vue-router";
import type { RouteRecordRaw } from "vue-router"; // <- type-only import

import Home from "../pages/Home.vue";
import Admin from "../pages/Admin.vue";

const routes: Array<RouteRecordRaw> = [
  { path: "/", name: "Home", component: Home },
  { path: "/admin", name: "Admin", component: Admin },
  // { path: "/queue/:id", name: "QueueDetail", component: Queuedetail },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;