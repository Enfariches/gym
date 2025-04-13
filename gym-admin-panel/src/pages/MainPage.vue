<template>
  <q-page class="bg-grey-1">
    <div class="content">
      <h1 class="page-title">Обзор системы</h1>

      <div class="stats-container">
        <SmallCard title="Всего пользователей" value="1,247" icon="people" />
        <SmallCard title="Просмотров сегодня" value="382" icon="visibility" />
        <SmallCard title="Завершено просмотров" value="76%" icon="check_circle" />
        <SmallCard title="Активных видео" value="24" icon="play_circle" />
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
            <!-- Add your chart component here -->
            <q-icon name="show_chart" size="100px" color="primary" class="q-mb-md" />
            <div class="text-grey-7">График просмотров видео по дням недели</div>
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
import SmallCard from '../components/SmallCard.vue';
import BigCard from 'src/components/mainPage/BigCard.vue';
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
  height: 300px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: rgba(108,117,125,1);
}
</style>
