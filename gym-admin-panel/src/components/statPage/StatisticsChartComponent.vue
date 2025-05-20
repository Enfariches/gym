<template>
  <div>
    <!-- Диаграмма -->
    <div class="chart-wrapper">
      <q-card flat bordered style="min-height: 300px; width: 100%">
        <q-card-section style="min-height: 300px; position: relative;">
          <div v-if="statisticsStore.loading" class="loading-container">
            <q-spinner size="50px" color="primary" />
            <div class="loading-text">Загрузка данных...</div>
          </div>
          <div v-else-if="!chartData.length" class="no-data-container">
            <div class="no-data-text">Нет данных для отображения</div>
          </div>
          <div v-else class="canvas-container">
            <canvas ref="canvasRef" />
          </div>
        </q-card-section>
      </q-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { nextTick, onMounted, ref, watch, computed, shallowRef, onBeforeUnmount } from 'vue';
import { Chart, registerables } from 'chart.js';
import { useStatisticsStore } from 'src/stores/statisticsStore';

export interface StatisticChartComponentProps {
  viewType: string;
  dateStart?: string;
  dateEnd?: string;
}

const props = withDefaults(defineProps<StatisticChartComponentProps>(), {
  viewType: 'full'
});

const canvasRef = ref<HTMLCanvasElement | null>(null);
const chartInstance = shallowRef<Chart | null>(null);
const statisticsStore = useStatisticsStore();
const chartUpdatePending = ref(false);

// Маппинг viewType из props в viewMode стора
const viewModeMap = {
  'full': 'media',
  'half': 'department',
  'declined': 'employee',
  'time': 'time'
};

// Данные для графика из стора
const chartData = computed(() => {
  return statisticsStore.chartData;
});

// Наблюдаем за изменением типа представления
watch(() => props.viewType, async (newType) => {
  const storeViewMode = viewModeMap[newType as keyof typeof viewModeMap] || 'media';
  statisticsStore.setViewMode(storeViewMode);
}, { immediate: false });

// Наблюдаем за изменениями данных в сторе
watch(() => chartData.value, () => {
  if (!chartUpdatePending.value && chartData.value.length > 0) {
    chartUpdatePending.value = true;
    nextTick(() => {
      updateOrCreateChart();
      chartUpdatePending.value = false;
    });
  }
}, { deep: true });

onMounted(async () => {
  // При монтировании компонента загружаем данные по умолчанию (статистика по отделу)
  const storeViewMode = viewModeMap[props.viewType as keyof typeof viewModeMap] || 'media';

  if (statisticsStore.statistics.length === 0) {
    await statisticsStore.fetchDepartmentStatistics();
  }

  statisticsStore.setViewMode(storeViewMode);
  window.addEventListener('resize', handleResize);
});

onBeforeUnmount(() => {
  // Уничтожаем график перед удалением компонента
  if (chartInstance.value) {
    chartInstance.value.destroy();
    chartInstance.value = null;
  }
  window.removeEventListener('resize', handleResize);
});

const handleResize = () => {
  if (!chartUpdatePending.value) {
    chartUpdatePending.value = true;
    setTimeout(() => {
      updateOrCreateChart();
      chartUpdatePending.value = false;
    }, 200);
  }
};

// Обновляем или создаем график
const updateOrCreateChart = async () => {
  if (chartData.value.length === 0) {
    if (chartInstance.value) {
      chartInstance.value.destroy();
      chartInstance.value = null;
    }
    return;
  }

  await nextTick();

  if (!canvasRef.value) {
    // Если canvas не найден, не пытаемся обновить график
    return;
  }

  const ctx = canvasRef.value.getContext("2d");

  if (!ctx) {
    console.error("Failed to get 2D context!");
    return;
  }

  if (chartInstance.value) {
    chartInstance.value.destroy();
    chartInstance.value = null;
    await new Promise((resolve) => setTimeout(resolve, 50));
  }

  Chart.register(...registerables);

  const labels = chartData.value.map(item => item.date);
  const data = chartData.value.map(item => item.count);

  const yAxisTitle = statisticsStore.viewMode === 'time'
    ? 'Просмотры по часам суток'
    : 'Количество просмотров';

  const xAxisTitle = statisticsStore.viewMode === 'time'
    ? 'Время суток'
    : 'Даты';

  try {
    chartInstance.value = new Chart(ctx, {
      type: 'line',
      data: {
        labels,
        datasets: [
          {
            label: 'Просмотры',
            backgroundColor: '#1976D2',
            borderColor: '#1976D2',
            borderWidth: 2,
            pointBackgroundColor: '#1976D2',
            data,
            tension: 0.4,
            fill: false,
          },
        ],
      },
      options: {
        responsive: true,
        maintainAspectRatio: true,
        plugins: {
          legend: {
            display: true,
            position: 'top'
          },
          tooltip: {
            mode: 'index',
            intersect: false,
          },
        },
        scales: {
          y: {
            beginAtZero: true,
            title: {
              display: true,
              text: yAxisTitle,
            },
          },
          x: {
            title: {
              display: true,
              text: xAxisTitle,
            },
          },
        },
      },
    });
  } catch (error) {
    console.error("Error creating chart:", error);
  }
};
</script>

<style scoped>
.chart-wrapper {
  width: 100%;
  display: flex;
  justify-content: center;
}

.canvas-container {
  position: relative;
  width: 100%;
  height: 250px;
}

.loading-container {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  height: 100%;
  min-height: 250px;
}

.loading-text {
  margin-top: 16px;
  color: #666;
  font-size: 14px;
}

.no-data-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  min-height: 250px;
  background-color: #f8f9fc;
  border-radius: 4px;
}

.no-data-text {
  color: #666;
  font-size: 16px;
}
</style>
