import { defineStore } from 'pinia';
import { getAdmin, updateAdmin } from '../services/adminService';
import type { Admin } from '../../protogen/v1/users/admin';

export const useAdminStore = defineStore('admin', {
  state: () => ({
    currentAdmin: null as Admin | null,
    loading: false,
    error: null as string | null,
  }),

  actions: {
    // Получить информацию об администраторе
    async fetchAdmin() {
      this.loading = true;
      this.error = null;

      try {
        const admin = await getAdmin();
        this.currentAdmin = admin;
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Не удалось загрузить данные администратора';
        console.error('Error fetching admin:', error);
      } finally {
        this.loading = false;
      }
    },

    // Обновить данные администратора
    async updateAdmin(fieldsToUpdate: string[]) {
      if (!this.currentAdmin) {
        throw new Error('Администратор не загружен');
      }

      this.loading = true;
      this.error = null;

      try {
        const updatedAdmin = await updateAdmin(this.currentAdmin, fieldsToUpdate);
        this.currentAdmin = updatedAdmin;
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Не удалось обновить данные администратора';
        console.error('Error updating admin:', error);
      } finally {
        this.loading = false;
      }
    },

    // Очистить данные администратора
    clearAdmin() {
      this.currentAdmin = null;
      this.error = null;
    },
  },

  getters: {
    // Получить текущего администратора
    admin: (state) => state.currentAdmin,

    // Проверка загрузки
    isLoading: (state) => state.loading,

    // Получить ошибку
    getError: (state) => state.error,
  },
});
