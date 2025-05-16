import { defineStore } from 'pinia';
import { listSchedules, getSchedule, createSchedules, updateSchedule, deleteSchedule } from '../services/scheduleService';
import type { Schedule } from '../../protogen/v1/schedule/schedule';

export const useScheduleStore = defineStore('schedule', {
  state: () => ({
    schedules: [] as Schedule[],
    currentSchedule: null as Schedule | null,
    loading: false,
    error: null as string | null,
  }),

  actions: {
    // Получить список расписаний
    async fetchSchedules() {
      this.loading = true;
      this.error = null;

      try {
        const schedules = await listSchedules();
        this.schedules = schedules;
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Не удалось загрузить расписания';
        console.error('Error fetching schedules:', error);
      } finally {
        this.loading = false;
      }
    },

    // Получить конкретное расписание
    async fetchSchedule(scheduleId: string) {
      this.loading = true;
      this.error = null;

      try {
        const schedule = await getSchedule(scheduleId);
        this.currentSchedule = schedule;
        return schedule;
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Не удалось загрузить расписание';
        console.error('Error fetching schedule:', error);
      } finally {
        this.loading = false;
      }
    },

    // Создать новые расписания
    async createSchedules(scheduleData: Schedule[]) {
      this.loading = true;
      this.error = null;

      try {
        const newSchedules = await createSchedules(scheduleData);
        this.schedules = [...this.schedules, ...newSchedules];
        return newSchedules;
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Не удалось создать расписания';
        console.error('Error creating schedules:', error);
        throw error;
      } finally {
        this.loading = false;
      }
    },

    // Обновить расписание
    async updateSchedule(scheduleData: Schedule, fieldsToUpdate: string[]) {
      this.loading = true;
      this.error = null;

      try {
        const updatedSchedule = await updateSchedule(scheduleData, fieldsToUpdate);

        // Обновляем текущее расписание, если оно загружено
        if (this.currentSchedule && this.currentSchedule.id === updatedSchedule.id) {
          this.currentSchedule = updatedSchedule;
        }

        // Обновляем расписание в списке
        this.schedules = this.schedules.map(s =>
          s.id === updatedSchedule.id ? updatedSchedule : s
        );

        return updatedSchedule;
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Не удалось обновить расписание';
        console.error('Error updating schedule:', error);
        throw error;
      } finally {
        this.loading = false;
      }
    },

    // Удалить расписание
    async deleteSchedule(scheduleId: string) {
      this.loading = true;
      this.error = null;

      try {
        await deleteSchedule(scheduleId);

        // Удаляем расписание из списка
        this.schedules = this.schedules.filter(s => s.id.toString() !== scheduleId);

        // Если текущее расписание было удалено, сбрасываем его
        if (this.currentSchedule && this.currentSchedule.id.toString() === scheduleId) {
          this.currentSchedule = null;
        }
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Не удалось удалить расписание';
        console.error('Error deleting schedule:', error);
        throw error;
      } finally {
        this.loading = false;
      }
    },
  },

  getters: {
    // Получить все расписания
    allSchedules: (state) => state.schedules,

    // Получить текущее расписание
    schedule: (state) => state.currentSchedule,

    // Проверка загрузки
    isLoading: (state) => state.loading,

    // Получить ошибку
    getError: (state) => state.error,
  },
});
