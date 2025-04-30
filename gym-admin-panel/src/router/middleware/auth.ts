import type { RouteLocationNormalized, NavigationGuardNext } from 'vue-router';
import { useAuthStore } from '../../stores/auth';

export const authGuard = async (
  to: RouteLocationNormalized,
  from: RouteLocationNormalized,
  next: NavigationGuardNext
) => {
  const authStore = useAuthStore();
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth);

  if (requiresAuth && !authStore.isAuthenticated) {
    next('/login');
  } else if (!requiresAuth && authStore.isAuthenticated) {
    next('/admin');
  } else {
    next();
  }
};
