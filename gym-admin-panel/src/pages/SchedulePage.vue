<template>
  <div class="content">
    <h1 class="page-title">Управление расписанием</h1>

    <ScheduleFilters :videos="videos" @departmentChange="onDepartmentChange" @videoChange="onVideoChange" />

    <div class="calendar-container">
      <ScheduleSidebar
        :departments="departments"
        :videos="videos"
        :selectedVideos="selectedVideos"
        :allVideosChecked="allVideosChecked"
        @departmentToggle="onDepartmentToggle"
        @videoToggle="onVideoToggle"
        @allVideosToggle="onAllVideosToggle"
      />

      <div class="schedule-grid">
        <ScheduleCard
          v-for="(day, idx) in days"
          :key="idx"
          :dayName="day.name"
          :dayOrder="day.order"
          :items="day.items"
          :videos="videos"
          @edit="onEdit"
          @delete="_removeSchedule"
          @add="onAdd"
        />
      </div>
    </div>

    <ScheduleEditModal
      v-model:isOpen="isModalOpen"
      :videos="videos"
      :scheduleId="scheduleId"
      :initialVideoTitle="initialVideoTitle"
      :initialTime="initialTime"
      :dayOrder="selectedDayOrder"
      @save="onSave"
    />
  </div>
</template>

<script lang="ts" setup>
import { onMounted, ref, computed } from 'vue'
import ScheduleFilters from '../components/schedulePage/ScheduleFilters.vue'
import ScheduleSidebar from '../components/schedulePage/ScheduleSidebar.vue'
import ScheduleCard from '../components/schedulePage/ScheduleCard.vue'
import ScheduleEditModal from '../components/schedulePage/ScheduleEditModal.vue'
import { useScheduleStore } from '../stores/scheduleStore'
import { useMediaStore } from '../stores/mediaStore'
import type { Schedule } from '../../protogen/v1/schedule/schedule'

const scheduleStore = useScheduleStore()
const mediaStore = useMediaStore()

// Переменные для работы с модальным окном добавления/редактирования расписания
const isModalOpen = ref(false)
const selectedDayOrder = ref('')
const initialVideoTitle = ref('')
const initialTime = ref('')
const scheduleId = ref('')
const isEditing = ref(false)

// Остальные переменные от департаментов и фильтров
const departments = ref([
  { id: 'all', name: 'Все отделы', checked: true },
  { id: 'dev', name: 'Разработчики', checked: false },
  { id: 'finance', name: 'Финансовый отдел', checked: false },
  { id: 'marketing', name: 'Маркетинг', checked: false },
  { id: 'hr', name: 'HR', checked: false },
])

const selectedVideos = ref<string[]>([])
const allVideosChecked = ref(true)

onMounted(async () => {
  await mediaStore.loadVideos()
  await scheduleStore.loadSchedules()
})

const videos = computed(() => mediaStore.getVideos.map(v => ({
  title: v.title
})))

const schedules = computed(() => scheduleStore.getSchedules)

interface ScheduleItem {
  ID: string;
  Time: string;
  VideoTitle: string;
  schedule: Schedule;
}

const days = computed(() => {
  const daysArr = [
    { items: [] as ScheduleItem[], name: "Понедельник", bg: 'bg-primary', order: 0 },
    { items: [] as ScheduleItem[], name: "Вторник", bg: 'bg-primary', order: 1 },
    { items: [] as ScheduleItem[], name: "Среда", bg: 'bg-primary', order: 2 },
    { items: [] as ScheduleItem[], name: "Четверг", bg: 'bg-primary', order: 3 },
    { items: [] as ScheduleItem[], name: "Пятница", bg: 'bg-primary', order: 4 },
    { items: [] as ScheduleItem[], name: "Суббота", bg: 'bg-primary', order: 5 },
    { items: [] as ScheduleItem[], name: "Воскресенье", bg: 'bg-primary', order: 6 },
  ];
  const schList = schedules.value ?? [];
  schList.forEach(sch => {
    const cron = sch.cronExpression.split(' ');
    if (cron.length === 5) {
      const dayOfWeek = Number(cron[4]) - 1; // 1=Пн, 7=Вс
      if (dayOfWeek >= 0 && dayOfWeek < 7) {
        const day = daysArr[dayOfWeek];
        if (day) {
          day.items.push({
            ID: sch.id?.toString() ?? '',
            Time: `${cron[1]?.padStart(2, '0') ?? '00'}:${cron[0]?.padStart(2, '0') ?? '00'}`,
            VideoTitle: mediaStore.getVideos.find(v => v.id?.toString() === sch.mediaId?.toString())?.title ?? '',
            schedule: sch
          });
        }
      }
    }
  });
  return daysArr;
})

const _removeSchedule = async (scheduleId: string) => {
  await scheduleStore.deleteSchedule(BigInt(scheduleId))
}

// Реализуем логику открытия модалки для добавления расписания
const onAdd = (dayOrder: number) => {
  selectedDayOrder.value = dayOrder.toString()
  initialVideoTitle.value = ''
  initialTime.value = ''
  scheduleId.value = ''
  isEditing.value = false
  isModalOpen.value = true
}

// Реализуем логику открытия модалки для редактирования
const onEdit = ({ videoName, scheduleId: schedId, dayOrder }: { videoName: string; scheduleId: string; dayOrder: number }) => {
  const schedule = schedules.value?.find(s => s.id?.toString() === schedId);
  if (schedule) {
    const cron = schedule.cronExpression?.split(' ') || [];
    if (cron.length === 5) {
      const hours = cron[1]?.padStart(2, '0') || '00';
      const minutes = cron[0]?.padStart(2, '0') || '00';
      initialTime.value = `${hours}:${minutes}`;
    }
    selectedDayOrder.value = dayOrder.toString();
    initialVideoTitle.value = videoName;
    scheduleId.value = schedId;
    isEditing.value = true;
    isModalOpen.value = true;
  }
}

// Реализуем логику сохранения расписания после модалки
const onSave = async (data: { videoTitle: string; time: string; dayOrder: string; scheduleId?: string }) => {
  const [hh, mm] = data.time.split(':')
  const cronExpression = `${mm} ${hh} * * ${Number(data.dayOrder) + 1}`

  // Находим ID медиа по его названию
  const video = mediaStore.getVideos.find(v => v.title === data.videoTitle);
  const mediaId = video?.id ?? BigInt(0);

  if (data.scheduleId) {
    const schedule: Schedule = {
      id: BigInt(data.scheduleId),
      cronExpression,
      isActive: true,
      mediaId: mediaId,
      adminId: BigInt(0),
      createdAt: ''
    }
    await scheduleStore.updateSchedule(schedule, ['cron_expression', 'media_id'])
  } else {
    const schedule: Schedule = {
      id: BigInt(0),
      cronExpression,
      isActive: true,
      mediaId: mediaId,
      adminId: BigInt(0),
      createdAt: ''
    }
    await scheduleStore.addSchedules([schedule])
  }
}

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const onDepartmentChange = (value: string) => {
  // Логика изменения отдела
}

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const onVideoChange = (value: string) => {
  // Логика изменения видео
}

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const onDepartmentToggle = (id: string) => {
  // Логика переключения отдела
}

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const onVideoToggle = (id: string) => {
  // Логика переключения видео
}

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const onAllVideosToggle = () => {
  // Логика переключения всех видео
}
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
}

.page-title {
  color: rgba(90,92,105,1);
  font-weight: bold;
  font-size: 32px;
  margin-bottom: 20px;
}

.schedule-toolbar {
  background: white;
  border: 2px solid rgba(227,230,240,1);
  border-radius: 6px;
  padding: 15px;
  margin-bottom: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.btn {
  padding: 10px 20px;
  font-size: 14px;
  font-weight: bold;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  display: flex;
  align-items: center;
}

.btn-primary {
  background-color: rgba(78,115,223,1);
  color: white;
}

.btn-outline {
  background-color: white;
  border: 2px solid rgba(78,115,223,1);
  color: rgba(78,115,223,1);
}

.filter-container {
  display: flex;
  gap: 15px;
}

.calendar-container {
  display: flex;
  gap: 20px;
}

.schedule-sidebar {
  flex: 0 0 250px;
  background: white;
  border: 2px solid rgba(227,230,240,1);
  border-radius: 6px;
  padding: 15px;
}

.sidebar-title {
  font-size: 18px;
  font-weight: bold;
  color: rgba(90,92,105,1);
  margin-bottom: 15px;
  padding-bottom: 10px;
  border-bottom: 1px solid rgba(227,230,240,1);
}

.department-list {
  margin-bottom: 30px;
}

.department-item {
  display: flex;
  align-items: center;
  padding: 8px 0;
  cursor: pointer;
}

.checkbox-container {
  display: flex;
  align-items: center;
  position: relative;
  padding-left: 30px;
  cursor: pointer;
  width: 100%;
}

.checkbox-container input {
  position: absolute;
  opacity: 0;
  cursor: pointer;
  height: 0;
  width: 0;
}

.checkmark {
  position: absolute;
  left: 0;
  height: 18px;
  width: 18px;
  border: 2px solid rgba(78,115,223,1);
  border-radius: 3px;
  background-color: white;
}

.checkbox-container:hover input ~ .checkmark {
  background-color: rgba(78,115,223,0.1);
}

.checkbox-container input:checked ~ .checkmark {
  background-color: rgba(78,115,223,1);
}

.checkmark:after {
  content: "";
  position: absolute;
  display: none;
}

.checkbox-container input:checked ~ .checkmark:after {
  display: block;
}

.checkbox-container .checkmark:after {
  left: 5px;
  top: 1px;
  width: 4px;
  height: 9px;
  border: solid white;
  border-width: 0 2px 2px 0;
  transform: rotate(45deg);
}

.department-name {
  font-size: 14px;
  color: rgba(90,92,105,1);
  margin-left: 5px;
}

.schedule-grid {
  flex: 1;
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
  align-items: start;
}

.schedule-element {
  flex: 1;
  background: white;
  border: 2px solid rgba(227,230,240,1);
  border-radius: 6px;
  overflow: hidden;
  height: 100%;
}

.schedule-element__head {
  background: rgba(78,115,223,1);
  color: white;
  padding: 12px;
  font-size: 16px;
  font-weight: bold;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.schedule-element__body {
  padding: 15px;
  min-height: 150px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.body__empty {
  text-align: center;
  color: rgba(108,117,125,1);
  margin: auto;
}

.schedule-item {
  background-color: white;
  border: 1px solid rgba(227,230,240,1);
  padding: 12px;
  border-radius: 4px;
  margin-bottom: 8px;
  transition: all 0.3s ease;
}

.schedule-item:hover {
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.schedule-item-odd {
  background-color: rgba(78,115,223,0.02);
}

.schedule-item__content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.schedule-item__title {
  flex: 1;
}

.video__label, .time__label {
  font-size: 12px;
  color: rgba(108,117,125,1);
  margin-bottom: 5px;
}

.schedule-item__time {
  margin: 0 15px;
  min-width: 80px;
}

.schedule-item__actions {
  display: flex;
  gap: 10px;
}

.action-btn {
  position: relative;
  width: 32px;
  height: 32px;
  border-radius: 4px;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  transition: all 0.3s ease;
}

.action-tooltip {
  position: absolute;
  background: rgba(0,0,0,0.8);
  color: white;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  bottom: 100%;
  left: 50%;
  transform: translateX(-50%);
  white-space: nowrap;
  opacity: 0;
  visibility: hidden;
  transition: all 0.2s ease;
  margin-bottom: 5px;
}

.action-btn:hover .action-tooltip {
  opacity: 1;
  visibility: visible;
}

.action-btn.edit {
  background-color: rgba(78,115,223,1);
  color: white;
}

.action-btn.edit:hover {
  background-color: rgba(78,115,223,0.8);
}

.action-btn.delete {
  background-color: rgba(231,74,59,1);
  color: white;
}

.action-btn.delete:hover {
  background-color: rgba(231,74,59,0.8);
}

.edit-modal {
  background: white;
  padding: 20px;
  min-width: 350px;
}

.modal-title {
  font-size: 20px;
  font-weight: bold;
  color: rgba(90,92,105,1);
  margin-bottom: 15px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  font-size: 14px;
  font-weight: bold;
  color: rgba(90,92,105,1);
  margin-bottom: 8px;
}

.video-select {
  width: 100%;
}

.time-input {
  width: 100%;
  padding: 8px;
  font-size: 14px;
  border: 2px solid rgba(227,230,240,1);
  border-radius: 4px;
}

.time-input:focus {
  outline: none;
  border-color: rgba(78,115,223,1);
}

.error-message {
  color: rgba(231,74,59,1);
  font-size: 12px;
  margin-top: 4px;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 20px;
}

.add-schedule {
  width: 100%;
  justify-content: center;
  gap: 8px;
  margin-top: auto;
  background-color: rgba(78,115,223,1);
  color: white;
  transition: all 0.3s ease;
}

.add-schedule:hover {
  background-color: rgba(78,115,223,0.8);
}
</style>
