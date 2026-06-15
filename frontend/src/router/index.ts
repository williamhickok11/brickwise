import { createRouter, createWebHistory } from 'vue-router'
import PropertiesView from '../views/PropertiesView.vue'
import DashboardView from '../views/DashboardView.vue'
import ResidencesView from '../views/ResidencesView.vue'
import REPSView from '../views/REPSView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/dashboard',
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: DashboardView,
    },
    {
      path: '/properties',
      name: 'properties',
      component: PropertiesView,
    },
    {
      path: '/residences',
      name: 'residences',
      component: ResidencesView,
    },
    {
      path: '/reps',
      name: 'reps',
      component: REPSView,
    },
    {
      path: '/reps-v2',
      redirect: { name: 'reps', query: { tab: 'capture' } },
    },
  ],
})

export default router
