<template>
  <q-page class="bg-grey-1">
    <div class="content">
      <h1 class="page-title">Обзор системы</h1>

      <div class="stats-container">
        <SmallCard title="Всего пользователей" :value="statisticsStore.statisticsSummary.uniqueUsers.toString()" icon="people" />
        <SmallCard title="Просмотров сегодня" :value="todayViews.toString()" icon="visibility" />
        <SmallCard title="Завершено просмотров" :value="`${statisticsStore.statisticsSummary.completedPercentage}%`" icon="check_circle" />
        <SmallCard title="Активных видео" :value="activeVideos.toString()" icon="play_circle" />
      </div>

      <div class="dashboard-row">
        <BigCard title="Активность системы" icon="history">
          <div class="activity-item">
            <q-icon name="circle" color="primary" size="xs" class="activity-indicator" />
            <div class="activity-text">Новое видео добавлено: "Утренняя зарядка"</div>
            <div class="activity-time">10:24</div>
          </div>
          <div class="activity-item">
            <q-icon name="circle" color="positive" size="xs" class="activity-indicator" />
            <div class="activity-text">Расписание обновлено для Отдела разработки</div>
            <div class="activity-time">09:41</div>
          </div>
          <div class="activity-item">
            <q-icon name="circle" color="warning" size="xs" class="activity-indicator" />
            <div class="activity-text">Пользователь Иванов И. достиг 5 достижений</div>
            <div class="activity-time">Вчера</div>
          </div>
          <div class="activity-item">
            <q-icon name="circle" color="negative" size="xs" class="activity-indicator" />
            <div class="activity-text">Проблема с загрузкой видео "Растяжка шеи"</div>
            <div class="activity-time">Вчера</div>
          </div>
        </BigCard>

        <BigCard title="Статистика просмотров" icon="insert_chart">
          <div class="chart-container">
            <template v-if="!statisticsLoaded">
              <q-spinner size="50px" color="primary" />
              <div class="loading-text">Загрузка данных...</div>
            </template>
            <StatisticsChartComponent v-else view-type="full" />
          </div>
        </BigCard>
      </div>

      <div class="dashboard-row">
        <BigCard title="Популярные видео" icon="trending_up">
          <div class="activity-item">
            <q-icon name="play_circle" color="primary" size="sm" class="activity-indicator" />
            <div class="activity-text">Разминка для рук и плеч</div>
            <div class="activity-time">842 просмотра</div>
          </div>
          <div class="activity-item">
            <q-icon name="play_circle" color="primary" size="sm" class="activity-indicator" />
            <div class="activity-text">Растяжка спины для офиса</div>
            <div class="activity-time">753 просмотра</div>
          </div>
          <div class="activity-item">
            <q-icon name="play_circle" color="primary" size="sm" class="activity-indicator" />
            <div class="activity-text">Гимнастика для глаз</div>
            <div class="activity-time">687 просмотров</div>
          </div>
        </BigCard>

        <BigCard title="Системные уведомления" icon="notifications">
          <div class="activity-item">
            <q-icon name="circle" color="negative" size="xs" class="activity-indicator" />
            <div class="activity-text">Необходимо обновить JWT-токены</div>
            <div class="activity-time">Важно</div>
          </div>
          <div class="activity-item">
            <q-icon name="circle" color="warning" size="xs" class="activity-indicator" />
            <div class="activity-text">Высокая нагрузка на сервер</div>
            <div class="activity-time">Предупреждение</div>
          </div>
          <div class="activity-item">
            <q-icon name="circle" color="positive" size="xs" class="activity-indicator" />
            <div class="activity-text">Резервное копирование завершено</div>
            <div class="activity-time">Успех</div>
          </div>
        </BigCard>
      </div>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { QPage } from 'quasar';
import { onMounted, computed, ref } from 'vue';
import SmallCard from '../components/SmallCard.vue';
import BigCard from 'src/components/mainPage/BigCard.vue';
import StatisticsChartComponent from 'src/components/statPage/StatisticsChartComponent.vue';
import { useStatisticsStore } from 'src/stores/statisticsStore';
import { useMediaStore } from 'src/stores/mediaStore';

const statisticsStore = useStatisticsStore();
const mediaStore = useMediaStore();
const statisticsLoaded = ref(false);

// Вычисляемые значения для карточек статистики
const activeVideos = computed(() => mediaStore.videos.length);
const todayViews = computed(() => {
  // Вычисляем просмотры за сегодня
  const today = new Date();
  today.setHours(0, 0, 0, 0);

  return statisticsStore.filteredStatistics.filter(stat => {
    if (!stat.createdAt) return false;
    const statDate = new Date(Number(stat.createdAt.seconds) * 1000);
    return statDate >= today;
  }).length;
});

onMounted(async () => {
  // Загружаем данные при загрузке страницы
  if (mediaStore.videos.length === 0) {
    await mediaStore.loadVideos();
  }

  await statisticsStore.fetchDepartmentStatistics();
  statisticsLoaded.value = true;
});
</script>

<style scoped>
.content {
  padding: 0 40px 40px;
}

.page-title {
  color: rgba(90,92,105,1);
  font-weight: bold;
  font-size: 32px;
  margin-bottom: 30px;
}

.stats-container {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 24px;
  margin-bottom: 30px;
}

.dashboard-row {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24px;
  margin-bottom: 24px;
}

.activity-item {
  display: flex;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid rgba(227,230,240,0.5);
}

.activity-indicator {
  margin-right: 12px;
}

.activity-text {
  font-size: 14px;
  color: rgba(90,92,105,1);
  flex: 1;
}

.activity-time {
  font-size: 12px;
  color: rgba(108,117,125,1);
  margin-left: 12px;
}

.chart-container {
  width: 100%;
  height: 300px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: rgba(108,117,125,1);
}

.loading-text {
  margin-top: 16px;
  color: #666;
  font-size: 14px;
}
</style>
