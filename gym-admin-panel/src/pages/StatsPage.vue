<template>
  <div class="content">
    <h1 class="page-title">Статистика просмотров</h1>

    <div class="filter-bar">
      <div class="filter-group">
        <span class="filter-label">Период:</span>
        <select class="filter-select" v-model="selectedDateRange">
          <option value="week">Эта неделя</option>
          <option value="month">Этот месяц</option>
          <option value="30days">За 30 дней</option>
          <option value="90days">За 90 дней</option>
          <option value="all">За все время</option>
          <option value="custom">Выбрать период</option>
        </select>

        <span v-if="selectedDateRange === 'custom'" class="date-picker-container">
          <q-date
            v-model="customDateRange"
            range
            minimal
            flat
            class="date-picker"
          />
        </span>

        <span class="filter-label">Отдел:</span>
        <select class="filter-select" v-model="selectedDepartmentId">
          <option :value="null">Все отделы</option>
          <option value="1">Отдел разработки</option>
          <option value="2">Бухгалтерия</option>
          <option value="3">Маркетинг</option>
          <option value="4">HR</option>
        </select>

        <button class="filter-btn" @click="applyFilters">Применить</button>
      </div>

      <button class="export-btn" @click="exportToPDF">Экспорт в PDF</button>
    </div>

    <div class="stats-container">
      <SmallCard
        title="Всего просмотров"
        :value="statisticsStore.statisticsSummary.totalViews.toString()"
        icon="visibility"
      />
      <SmallCard
        title="Просмотрено до конца"
        :value="`${statisticsStore.statisticsSummary.completedPercentage}%`"
        icon="check_circle"
      />
      <SmallCard
        title="Уникальных пользователей"
        :value="statisticsStore.statisticsSummary.uniqueUsers.toString()"
        icon="people"
      />
      <SmallCard
        title="Среднее время просмотра"
        :value="statisticsStore.statisticsSummary.averageViewTime"
        icon="timer"
      />
    </div>

    <div class="dashboard-row">
      <div class="dashboard-card">
        <h2 class="card-title">Динамика просмотров</h2>
        <div class="chart-container">
          <StatisticsChartComponent :view-type="mapStoreViewModeToChartViewType(statisticsStore.viewMode)" />
        </div>
      </div>
    </div>

    <div class="dashboard-card">
      <div class="tab-nav">
        <div class="tab-item"
          v-for="option in selectableOptions"
          :key="option.type"
          :class="{ active: statisticsStore.viewMode === option.storeType }"
          @click="statisticsStore.setViewMode(option.storeType)"
        >
          {{ option.label }}
        </div>
      </div>

      <h2 class="card-title">Статистика просмотров</h2>

      <div v-if="statisticsStore.loading" class="loading-container">
        <q-spinner size="50px" color="primary" />
        <div class="loading-text">Загрузка данных...</div>
      </div>

      <div v-else-if="!statisticsStore.filteredStatistics.length" class="empty-data-message">
        Нет данных для отображения
      </div>

      <div v-else class="table-container">
        <table>
          <thead>
            <tr>
              <th>Название видео</th>
              <th>Сотрудник</th>
              <th>Дата просмотра</th>
              <th>Статус</th>
              <th>Процент просмотра</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(stat, index) in statisticsStore.filteredStatistics" :key="stat.id?.toString() || index">
              <td>{{ stat.mediaTitle }}</td>
              <td>{{ stat.employeeName }} {{ stat.employeeSurname }}</td>
              <td>{{ formatDate(stat.createdAt) }}</td>
              <td>
                <span :class="getProgressClass(stat.progress)">
                  {{ getProgressLabel(stat.progress) }}
                </span>
              </td>
              <td>
                <div class="progress-bar">
                  <div class="progress-fill" :style="`width: ${Number(stat.percentageView)}%;`"></div>
                </div>
                {{ stat.percentageView?.toString() }}%
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue';
import StatisticsChartComponent from 'src/components/statPage/StatisticsChartComponent.vue';
import SmallCard from '../components/SmallCard.vue';
import { useStatisticsStore } from '../stores/statisticsStore';
import { MediaProgress } from '../../protogen/v1/statistics/statistics';
import type { Timestamp } from '../../protogen/google/protobuf/timestamp';

const statisticsStore = useStatisticsStore();
const selectedDateRange = ref('all');
const customDateRange = ref({ from: '', to: '' });
const selectedDepartmentId = ref<number | null>(null);

// Маппинг типов представления
const selectableOptions = [
  { storeType: 'media', type: 'full', label: 'По видео' },
  { storeType: 'department', type: 'half', label: 'По отделам' },
  { storeType: 'employee', type: 'declined', label: 'По сотрудникам' },
  { storeType: 'time', type: 'time', label: 'По времени дня' }
];

// Маппинг из типа представления стора в тип для компонента диаграммы
const mapStoreViewModeToChartViewType = (storeViewMode: string) => {
  const option = selectableOptions.find(opt => opt.storeType === storeViewMode);
  return option ? option.type : 'full';
};

// Получение класса стиля в зависимости от статуса прогресса
const getProgressClass = (progress: MediaProgress) => {
  switch (progress) {
    case MediaProgress.COMPLETED:
      return 'status-completed';
    case MediaProgress.INCOMPLETE:
      return 'status-incomplete';
    case MediaProgress.SKIPPED:
      return 'status-skipped';
    default:
      return '';
  }
};

// Получение текстовой метки в зависимости от статуса прогресса
const getProgressLabel = (progress: MediaProgress) => {
  switch (progress) {
    case MediaProgress.COMPLETED:
      return 'Завершено';
    case MediaProgress.INCOMPLETE:
      return 'Не завершено';
    case MediaProgress.SKIPPED:
      return 'Пропущено';
    default:
      return 'Не определено';
  }
};

// Форматирование даты для отображения
const formatDate = (timestamp?: Timestamp) => {
  if (!timestamp) return '-';

  const date = new Date(Number(timestamp.seconds) * 1000);
  return date.toLocaleDateString('ru-RU', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  });
};

// Применение фильтров
const applyFilters = () => {
  let startDate: Date | null = null;
  let endDate: Date | null = null;

  const now = new Date();

  switch (selectedDateRange.value) {
    case 'week': {
      // Начало текущей недели (понедельник)
      const day = now.getDay();
      const diff = now.getDate() - day + (day === 0 ? -6 : 1);
      startDate = new Date(now.setDate(diff));
      startDate.setHours(0, 0, 0, 0);
      endDate = new Date();
      break;
    }
    case 'month': {
      // Начало текущего месяца
      startDate = new Date(now.getFullYear(), now.getMonth(), 1);
      endDate = new Date();
      break;
    }
    case '30days': {
      // 30 дней назад
      startDate = new Date();
      startDate.setDate(startDate.getDate() - 30);
      endDate = new Date();
      break;
    }
    case '90days': {
      // 90 дней назад
      startDate = new Date();
      startDate.setDate(startDate.getDate() - 90);
      endDate = new Date();
      break;
    }
    case 'custom': {
      // Пользовательский период
      if (customDateRange.value.from) {
        startDate = new Date(customDateRange.value.from);
      }
      if (customDateRange.value.to) {
        endDate = new Date(customDateRange.value.to);
        endDate.setHours(23, 59, 59, 999); // До конца дня
      }
      break;
    }
    default: {
      // За все время
      startDate = null;
      endDate = null;
    }
  }

  // Применяем фильтры к стору
  statisticsStore.setDateRange({
    startDate,
    endDate
  });

  statisticsStore.setDepartmentFilter(selectedDepartmentId.value);
};

// Экспорт в PDF
const exportToPDF = async () => {
  await statisticsStore.exportToPDF();
};

onMounted(async () => {
  // Загружаем статистику при монтировании компонента
  await statisticsStore.fetchDepartmentStatistics();
});
</script>

<style scoped>
* {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
  font-family: 'Inter', sans-serif;
}

.content {
  padding: 20px 40px;
  background-color: #f8f9fc;
  max-width: 1600px;
  margin: 0 auto;
}

.page-title {
  color: rgba(90,92,105,1);
  font-weight: bold;
  font-size: 32px;
  margin-bottom: 20px;
}

.filter-bar {
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
  align-items: center;
  background: white;
  border: 2px solid rgba(227,230,240,1);
  border-radius: 6px;
  padding: 15px;
}

.filter-group {
  display: flex;
  align-items: center;
  gap: 15px;
  flex-wrap: wrap;
}

.filter-label {
  font-size: 14px;
  color: rgba(90,92,105,1);
  font-weight: bold;
}

.filter-select {
  padding: 8px 12px;
  font-size: 14px;
  border: 2px solid rgba(227,230,240,1);
  border-radius: 6px;
  color: rgba(90,92,105,1);
  background-color: white;
  min-width: 150px;
}

.date-picker-container {
  position: relative;
  display: inline-block;
}

.date-picker {
  position: absolute;
  top: 100%;
  left: 0;
  z-index: 10;
  background: white;
  box-shadow: 0px 3px 8px rgba(0, 0, 0, 0.15);
  border-radius: 6px;
  min-width: 300px;
}

.filter-btn {
  padding: 8px 15px;
  font-size: 14px;
  background-color: rgba(78,115,223,1);
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-weight: bold;
  transition: all 0.3s ease;
}

.filter-btn:hover {
  background-color: rgba(78,115,223,0.8);
}

.export-btn {
  padding: 8px 15px;
  font-size: 14px;
  background-color: rgba(231,74,59,1);
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-weight: bold;
  transition: all 0.3s ease;
}

.export-btn:hover {
  background-color: rgba(231,74,59,0.8);
}

.stats-container {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 20px;
  margin-bottom: 20px;
}

.stat-card {
  background: white;
  border: 2px solid rgba(227,230,240,1);
  border-radius: 6px;
  padding: 15px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  min-height: 140px;
}

.stat-title {
  color: rgba(78,115,223,1);
  font-size: 16px;
  font-weight: bold;
  margin-bottom: 8px;
}

.stat-value {
  font-size: 32px;
  font-weight: bold;
  color: rgba(90,92,105,1);
}

.stat-description {
  font-size: 14px;
  color: rgba(108,117,125,1);
  margin-top: 8px;
}

.dashboard-row {
  margin-bottom: 20px;
}

.dashboard-card {
  background: white;
  border: 2px solid rgba(227,230,240,1);
  border-radius: 6px;
  padding: 15px;
}

.card-title {
  color: rgba(90,92,105,1);
  font-weight: bold;
  font-size: 18px;
  margin-bottom: 15px;
  padding-bottom: 10px;
  border-bottom: 1px solid rgba(227,230,240,1);
}

.chart-container {
  height: 250px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: rgba(108,117,125,1);
}

.tab-nav {
  display: flex;
  margin-bottom: 15px;
  border-bottom: 1px solid rgba(227,230,240,1);
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

.tab-item {
  padding: 10px 15px;
  font-size: 14px;
  color: rgba(108,117,125,1);
  cursor: pointer;
  margin-right: 8px;
  white-space: nowrap;
  border-bottom: 2px solid transparent;
  transition: all 0.3s ease;
}

.tab-item.active {
  color: rgba(78,115,223,1);
  font-weight: bold;
  border-bottom: 2px solid rgba(78,115,223,1);
}

.loading-container {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  height: 200px;
}

.loading-text {
  margin-top: 16px;
  color: #666;
  font-size: 14px;
}

.empty-data-message {
  text-align: center;
  padding: 40px;
  color: #666;
  font-size: 16px;
  background-color: #f8f9fc;
  border-radius: 6px;
}

.table-container {
  width: 100%;
  margin-top: 15px;
  overflow-x: auto;
  border: 1px solid rgba(227,230,240,1);
  border-radius: 6px;
}

table {
  width: 100%;
  border-collapse: collapse;
}

th, td {
  padding: 10px 15px;
  text-align: left;
  border-bottom: 1px solid rgba(227,230,240,1);
  font-size: 14px;
}

th {
  font-weight: bold;
  color: rgba(90,92,105,1);
  background-color: rgba(248,249,252,1);
}

td {
  color: rgba(90,92,105,1);
}

.progress-bar {
  width: 100%;
  height: 18px;
  background-color: rgba(227,230,240,1);
  border-radius: 6px;
  overflow: hidden;
  margin: 4px 0;
}

.progress-fill {
  height: 100%;
  background-color: rgba(78,115,223,1);
  transition: width 0.3s ease;
}

.status-completed {
  color: rgba(28,200,138,1);
  font-weight: bold;
}

.status-incomplete {
  color: rgba(246,194,62,1);
  font-weight: bold;
}

.status-skipped {
  color: rgba(231,74,59,1);
  font-weight: bold;
}

@media (max-width: 768px) {
  .content {
    padding: 15px;
  }

  .filter-group {
    flex-direction: column;
    align-items: stretch;
  }

  .filter-select {
    width: 100%;
  }

  .stats-container {
    grid-template-columns: 1fr;
  }
}
</style>
