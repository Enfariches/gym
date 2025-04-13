<template>
  <div class="content">
    <h1 class="page-title">Статистика просмотров</h1>
    
    <div class="filter-bar">
      <div class="filter-group">
        <span class="filter-label">Период:</span>
        <select class="filter-select">
          <option>Эта неделя</option>
          <option>Этот месяц</option>
          <option>За 30 дней</option>
          <option>За 90 дней</option>
          <option>За все время</option>
        </select>
        
        <span class="filter-label">Отдел:</span>
        <select class="filter-select">
          <option>Все отделы</option>
          <option>Отдел разработки</option>
          <option>Бухгалтерия</option>
          <option>Маркетинг</option>
          <option>HR</option>
        </select>
        
        <button class="filter-btn">Применить</button>
      </div>
      
      <button class="export-btn">Экспорт в Excel</button>
    </div>
    
    <div class="stats-container">
      <SmallCard title="Всего просмотров" value="3,842" icon="visibility" />
      <SmallCard title="Просмотрено до конца" value="76%" icon="check_circle" />
      <SmallCard title="Уникальных пользователей" value="945" icon="people" />
      <SmallCard title="Среднее время просмотра" value="8:24" icon="timer" />
    </div>
    
    <div class="dashboard-row">
      <div class="dashboard-card">
        <h2 class="card-title">Динамика просмотров</h2>
        <div class="chart-container">
          <StatisticsChartComponent :view-type="state.activeState" />
        </div>
      </div>
    </div>
    
    <div class="dashboard-card">
      <div class="tab-nav">
        <div class="tab-item" 
          v-for="option in selectableOptions" 
          :key="option.type"
          :class="{ active: option.type === state.activeState }"
          @click="state.activeState = option.type"
        >
          {{ option.label }}
        </div>
      </div>
      
      <h2 class="card-title">Статистика просмотров по видео</h2>
      
      <div class="table-container">
        <table>
          <thead>
            <tr>
              <th>Название видео</th>
              <th>Продолжительность</th>
              <th>Всего просмотров</th>
              <th>Завершено</th>
              <th>Процент завершения</th>
              <th>Ср. время просмотра</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td>Разминка для рук и плеч</td>
              <td>5:20</td>
              <td>842</td>
              <td>721</td>
              <td>
                <div class="progress-bar">
                  <div class="progress-fill" style="width: 85%;"></div>
                </div>
                85.6%
              </td>
              <td>4:55</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue'
import StatisticsChartComponent from 'src/components/statPage/StatisticsChartComponent.vue'
import SmallCard from '../components/SmallCard.vue'

interface errorWithMessage {
  isError: boolean
  errorMessage: string
}

const state = ref({
  activeState: 'full',
  selectedStartDate: '',
  selectedEndDate: '',
  dateRangeError: { isError: false, errorMessage: '' } as errorWithMessage
})

const selectableOptions = [
  { type: 'full', label: 'По видео' },
  { type: 'half', label: 'По отделам' },
  { type: 'declined', label: 'По сотрудникам' },
  { type: 'time', label: 'По времени дня' }
]
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
  background-color: rgba(28,200,138,1);
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-weight: bold;
  transition: all 0.3s ease;
}

.export-btn:hover {
  background-color: rgba(28,200,138,0.8);
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
