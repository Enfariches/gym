import type { RouteRecordRaw } from 'vue-router';
import { authGuard } from './middleware/auth';

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/admin',
    beforeEnter: authGuard
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('pages/LoginPage.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/register',
    name: 'register',
    component: () => import('pages/RegisterPage.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/reset-password',
    name: 'reset-password',
    component: () => import('pages/ResetPasswordPage.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/change-password',
    name: 'change-password',
    component: () => import('pages/ChangePasswordPage.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/admin',
    component: () => import('layouts/MainLayout.vue'),
    meta: { requiresAuth: true },
    beforeEnter: authGuard,
    children: [
      { path: '', component: () => import('pages/MainPage.vue') },
      { path: 'videos', component: () => import('pages/VideosPage.vue') },
      { path: 'schedule', component: () => import('pages/SchedulePage.vue') },
      { path: 'stats', component: () => import('pages/StatsPage.vue') },
      { path: 'profile', component: () => import('pages/ProfilePage.vue') },
    ],
  },
  {
    path: '/:catchAll(.*)*',
    component: () => import('pages/ErrorNotFound.vue'),
  }
];

export default routes;
