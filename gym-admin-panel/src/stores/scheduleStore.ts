import { defineStore } from 'pinia';
import { listSchedulesGrpc, createSchedulesGrpc, updateScheduleGrpc, deleteScheduleGrpc } from '../services/scheduleService';
import type { Schedule } from '../../protogen/v1/schedule/schedule';

export const useScheduleStore = defineStore('schedule', {
  state: () => ({
    schedules: [] as Schedule[],
    loading: false,
    error: null as string | null,
  }),
  actions: {
    async loadSchedules() {
      this.loading = true;
      this.error = null;
      try {
        this.schedules = await listSchedulesGrpc();
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Не удалось загрузить расписания';
      } finally {
        this.loading = false;
      }
    },
    async addSchedules(schedules: Schedule[]) {
      this.loading = true;
      this.error = null;
      try {
        await createSchedulesGrpc(schedules);
        await this.loadSchedules();
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Не удалось добавить расписание';
      } finally {
        this.loading = false;
      }
    },
    async updateSchedule(schedule: Schedule, fieldMask?: string[]) {
      this.loading = true;
      this.error = null;
      try {
        await updateScheduleGrpc(schedule, fieldMask);
        await this.loadSchedules();
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Не удалось обновить расписание';
      } finally {
        this.loading = false;
      }
    },
    async deleteSchedule(scheduleId: bigint) {
      this.loading = true;
      this.error = null;
      try {
        await deleteScheduleGrpc(scheduleId);
        await this.loadSchedules();
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Не удалось удалить расписание';
      } finally {
        this.loading = false;
      }
    },
    clearError() {
      this.error = null;
    }
  },
  getters: {
    getSchedules: (state) => state.schedules,
    isLoading: (state) => state.loading,
    getError: (state) => state.error,
  }
});
