import { createRouter, createWebHistory } from 'vue-router'
import Home from '../components/SearchEntrance.vue'
import Result from '../components/SearchResult.vue'

const routes = [
  { path: '/', component: Home },
  { path: '/result/:query', name: "Result", component: Result }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

export default router;