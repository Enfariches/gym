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
    async fetchAdmin(adminId: number) {
      this.loading = true;
      this.error = null;

      try {
        const admin = await getAdmin(adminId);
        this.currentAdmin = admin;
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Не удалось загрузить данные администратора';
        console.error('Error fetching admin:', error);
      } finally {
        this.loading = false;
      }
    },

    // Обновить данные администратора
    async updateAdmin(adminData: Partial<Admin>, fieldsToUpdate: string[]) {
      if (!this.currentAdmin) {
        this.error = 'Нет данных администратора для обновления';
        return;
      }

      // Создаем обновленный объект администратора
      const updatedAdmin: Admin = {
        ...this.currentAdmin,
        ...adminData,
        // Преобразуем id в bigint, если он был предоставлен как number
        id: adminData.id ? BigInt(adminData.id.toString()) : this.currentAdmin.id
      };

      this.loading = true;
      this.error = null;

      try {
        const admin = await updateAdmin(updatedAdmin, fieldsToUpdate);
        this.currentAdmin = admin;
        return admin;
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Не удалось обновить данные администратора';
        console.error('Error updating admin:', error);
        throw error;
      } finally {
        this.loading = false;
      }
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
