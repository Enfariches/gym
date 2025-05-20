import { defineStore } from 'pinia';
import {
  createStatistics,
  getEmployeeStatistics,
  listMediaStatistics,
  listEmployeeStatistics,
  listDepartmentStatistics,
  exportStatisticsToPDF
} from '../services/statisticsService';
import { MediaProgress, type Statistic } from '../../protogen/v1/statistics/statistics';

interface StatisticSummary {
  totalViews: number;
  completedPercentage: number;
  uniqueUsers: number;
  averageViewTime: string;
}

interface DateRange {
  startDate: Date | null;
  endDate: Date | null;
}

export const useStatisticsStore = defineStore('statistics', {
  state: () => ({
    // Filtered statistics
    statistics: [] as Statistic[],
    // Current view mode
    viewMode: 'media' as string,
    // Filters
    filters: {
      dateRange: {
        startDate: null,
        endDate: null
      } as DateRange,
      departmentId: null as number | null,
    },
    // Currently selected IDs
    selectedMediaId: null as number | null,
    selectedEmployeeId: null as number | null,
    // Loading state
    loading: false,
    error: null as string | null
  }),

  actions: {
    // Записать статистику просмотра
    async recordStatistics(progress: MediaProgress, percentageView: number, mediaId: number) {
      this.error = null;
      this.loading = true;
      try {
        await createStatistics(progress, percentageView, mediaId);
        return true;
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Ошибка при записи статистики';
        console.error('Ошибка при записи статистики:', error);
        return false;
      } finally {
        this.loading = false;
      }
    },

    // Получить статистику сотрудника по конкретному видео
    async fetchEmployeeStatistics(employeeId: number, mediaId: number) {
      this.error = null;
      this.loading = true;
      try {
        const data = await getEmployeeStatistics(employeeId, mediaId);
        this.statistics = [data];
        this.selectedEmployeeId = Number(employeeId);
        this.selectedMediaId = Number(mediaId);
        return data;
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Ошибка при получении статистики сотрудника';
        console.error('Ошибка при получении статистики сотрудника:', error);
        return null;
      } finally {
        this.loading = false;
      }
    },

    // Получить статистику по всем сотрудникам для конкретного медиа
    async fetchMediaStatistics(mediaId: number) {
      this.error = null;
      this.loading = true;
      try {
        const data = await listMediaStatistics(mediaId);
        this.statistics = data;
        this.selectedMediaId = Number(mediaId);
        this.selectedEmployeeId = null;
        this.viewMode = 'media';
        return data;
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Ошибка при получении статистики по видео';
        console.error('Ошибка при получении статистики по видео:', error);
        return [];
      } finally {
        this.loading = false;
      }
    },

    // Получить статистику сотрудника по всем видео
    async fetchEmployeeAllStatistics(employeeId: number) {
      this.error = null;
      this.loading = true;
      try {
        const data = await listEmployeeStatistics(employeeId);
        this.statistics = data;
        this.selectedEmployeeId = Number(employeeId);
        this.selectedMediaId = null;
        this.viewMode = 'employee';
        return data;
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Ошибка при получении статистики сотрудника';
        console.error('Ошибка при получении статистики сотрудника:', error);
        return [];
      } finally {
        this.loading = false;
      }
    },

    // Получить статистику по отделу
    async fetchDepartmentStatistics() {
      this.error = null;
      this.loading = true;
      try {
        const data = await listDepartmentStatistics();
        this.statistics = data;
        this.selectedEmployeeId = null;
        this.selectedMediaId = null;
        this.viewMode = 'department';
        return data;
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Ошибка при получении статистики отдела';
        console.error('Ошибка при получении статистики отдела:', error);
        return [];
      } finally {
        this.loading = false;
      }
    },

    // Экспортировать статистику в PDF
    async exportToPDF() {
      this.error = null;
      this.loading = true;
      try {
        const blob = await exportStatisticsToPDF(
          this.filters.dateRange.startDate,
          this.filters.dateRange.endDate,
          this.filters.departmentId
        );

        // Создаем ссылку на PDF и скачиваем его
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `statistics-export-${new Date().toISOString().split('T')[0]}.pdf`;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        window.URL.revokeObjectURL(url);

        return true;
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Ошибка при экспорте статистики в PDF';
        console.error('Ошибка при экспорте статистики в PDF:', error);
        return false;
      } finally {
        this.loading = false;
      }
    },

    // Изменить режим отображения статистики
    setViewMode(mode: string) {
      this.viewMode = mode;

      // Загружаем соответствующие данные при изменении режима
      switch (mode) {
        case 'media':
          if (this.selectedMediaId) {
            this.fetchMediaStatistics(this.selectedMediaId);
          }
          break;
        case 'employee':
          if (this.selectedEmployeeId) {
            this.fetchEmployeeAllStatistics(this.selectedEmployeeId);
          }
          break;
        case 'department':
          this.fetchDepartmentStatistics();
          break;
        case 'time':
          // Для отображения по времени используем все текущие данные
          break;
      }
    },

    // Установить фильтр по диапазону дат
    setDateRange(range: DateRange) {
      this.filters.dateRange = range;
    },

    // Установить фильтр по отделу
    setDepartmentFilter(departmentId: number | null) {
      this.filters.departmentId = departmentId;
    },

    // Сбросить все фильтры
    resetFilters() {
      this.filters = {
        dateRange: {
          startDate: null,
          endDate: null
        },
        departmentId: null
      };
    }
  },

  getters: {
    // Получить отфильтрованную статистику
    filteredStatistics(): Statistic[] {
      let result = this.statistics;

      // Фильтр по дате
      if (this.filters.dateRange.startDate || this.filters.dateRange.endDate) {
        result = result.filter(stat => {
          if (!stat.createdAt) return true;

          const statDate = new Date(stat.createdAt.seconds.toString());

          if (this.filters.dateRange.startDate && statDate < this.filters.dateRange.startDate) {
            return false;
          }

          if (this.filters.dateRange.endDate && statDate > this.filters.dateRange.endDate) {
            return false;
          }

          return true;
        });
      }

      return result;
    },

    // Общая сводка по статистике
    statisticsSummary(): StatisticSummary {
      const statistics = this.filteredStatistics;

      if (!statistics.length) {
        return {
          totalViews: 0,
          completedPercentage: 0,
          uniqueUsers: 0,
          averageViewTime: '0:00'
        };
      }

      // Общее количество просмотров
      const totalViews = statistics.length;

      // Процент завершенных просмотров
      const completedViews = statistics.filter(s => s.progress === MediaProgress.COMPLETED).length;
      const completedPercentage = totalViews > 0 ? Math.round((completedViews / totalViews) * 100) : 0;

      // Уникальные пользователи
      const uniqueUserIds = new Set(statistics.map(s => `${s.employeeName} ${s.employeeSurname}`));

      // Среднее время просмотра (в секундах)
      // Здесь используем процент просмотра, предполагая, что видео имеет среднюю продолжительность 10 минут
      const totalPercentage = statistics.reduce((sum, stat) => sum + Number(stat.percentageView), 0);
      const averagePercentage = totalViews > 0 ? totalPercentage / totalViews : 0;
      // Предполагаем, что 100% = 10 минут = 600 секунд
      const averageSeconds = (averagePercentage / 100) * 600;
      const minutes = Math.floor(averageSeconds / 60);
      const seconds = Math.floor(averageSeconds % 60);

      return {
        totalViews,
        completedPercentage,
        uniqueUsers: uniqueUserIds.size,
        averageViewTime: `${minutes}:${seconds.toString().padStart(2, '0')}`
      };
    },

    // Данные для графика в зависимости от текущего режима просмотра
    chartData(): Array<{date: string, count: number}> {
      const statistics = this.filteredStatistics;

      // Создаем карту для агрегации данных
      const dataMap = new Map<string, number>();

      statistics.forEach(stat => {
        if (!stat.createdAt) return;

        let key;
        const date = new Date(Number(stat.createdAt.seconds) * 1000);

        switch (this.viewMode) {
          case 'time': {
            // Группировка по времени дня (часам)
            const hour = date.getHours();
            key = `${hour}:00`;
            break;
          }
          case 'media':
          case 'employee':
          case 'department':
          default: {
            // Группировка по дням для других режимов
            key = date.toLocaleDateString('ru-RU', { day: '2-digit', month: '2-digit' });
            break;
          }
        }

        if (dataMap.has(key)) {
          dataMap.set(key, dataMap.get(key)! + 1);
        } else {
          dataMap.set(key, 1);
        }
      });

      // Преобразуем карту в массив для отображения
      return Array.from(dataMap.entries()).map(([date, count]) => ({
        date,
        count
      })).sort((a, b) => a.date.localeCompare(b.date));
    }
  }
});
