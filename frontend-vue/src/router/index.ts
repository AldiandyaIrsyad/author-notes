import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import RegisterView from '../views/RegisterView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/login' // Redirect root path to login
    },
    {
      path: '/login',
      name: 'login',
      component: LoginView
    },
    {
      path: '/register',
      name: 'register',
      component: RegisterView
    },
    // Template =>
    // {
    //   path: '/dashboard',
    //   name: 'dashboard',
    //   // route level code-splitting
    //   // this generates a separate chunk (About.[hash].js) for this route
    //   // which is lazy-loaded when the route is visited.
    //   component: () => import('../views/DashboardView.vue') // Example
    // }

    
    // Protected Template =>
    // {
    //   path: '/dashboard',
    //   name: 'dashboard',
    //
    //   component: () => import('../views/DashboardView.vue'), 
    //   beforeEnter: (to, from, next) => {
    //     const isAuthenticated = false; // Replace with actual authentication check
    //     if (isAuthenticated) {
    //       next(); // Allow access to the route
    //     } else {
    //       next({ name: 'login' }); // Redirect to login if not authenticated
    //     }
    //   }
  ]
})

export default router
