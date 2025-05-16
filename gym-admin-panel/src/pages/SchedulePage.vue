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
          @edit="handleEditEvent"
          @delete="onDelete"
          @add="onAdd"
        />
      </div>
    </div>

    <ScheduleEditModal
      v-model:isOpen="isEditModalOpen"
      :videos="videos"
      v-bind="selectedScheduleForModal.id ? { scheduleId: selectedScheduleForModal.id } : {}"
      :initialVideoName="selectedScheduleForModal.videoName || ''"
      :initialTime="selectedScheduleForModal.time"
      :initialIsActive="selectedScheduleForModal.isActive"
      :dayOrder="selectedScheduleForModal.dayOrder.toString()"
      @save="onSave"
    />

    <!-- Индикатор загрузки -->
    <q-inner-loading :showing="scheduleStore.isLoading">
      <q-spinner size="50px" color="primary" />
    </q-inner-loading>
  </div>
</template>

<script lang="ts" setup>
import { onMounted, ref } from 'vue'
import { useQuasar } from 'quasar'
import ScheduleFilters from '../components/schedulePage/ScheduleFilters.vue'
import ScheduleSidebar from '../components/schedulePage/ScheduleSidebar.vue'
import ScheduleCard from '../components/schedulePage/ScheduleCard.vue'
import ScheduleEditModal from '../components/schedulePage/ScheduleEditModal.vue'
import { useScheduleStore } from '../stores/scheduleStore'
import type { Schedule as ProtoSchedule } from '../../protogen/v1/schedule/schedule'

const $q = useQuasar();
const scheduleStore = useScheduleStore();

interface VideoItem {
  ID: string;
  Name: string;
}

interface ScheduleItem {
  ID: string;
  Time: string;
  VideoID: string;
  Name: string;
  is_active: boolean;
  dayOrder: number;
  cron_expression: string;
}

// Data for preparing modal and for onSave to process modal's emitted data
interface ScheduleModalInputData {
  id?: string;
  videoName?: string; // videoName is for initial display, not strictly part of emitted save data
  time: string;
  isActive: boolean;
  dayOrder: number; // Internally, we might prefer number for dayOrder
  videoId?: string;
}

// Actual payload from ScheduleEditModal's @save event
interface ScheduleEditModalSavePayload {
  videoId: string;
  time: string;
  dayOrder: string; // Modal emits dayOrder as string
  scheduleId: string | undefined;
  isActive: boolean;
}

// Define ScheduleCardEditPayload based on what ScheduleCard emits
interface ScheduleCardEditPayload {
  videoName: string;
  videoId: string;
  scheduleId: string;
  dayOrder: number;
}

const days = ref([
  { items: [] as ScheduleItem[], name: "Понедельник", order: 1 },
  { items: [] as ScheduleItem[], name: "Вторник", order: 2 },
  { items: [] as ScheduleItem[], name: "Среда", order: 3 },
  { items: [] as ScheduleItem[], name: "Четверг", order: 4 },
  { items: [] as ScheduleItem[], name: "Пятница", order: 5 },
  { items: [] as ScheduleItem[], name: "Суббота", order: 6 },
  { items: [] as ScheduleItem[], name: "Воскресенье", order: 0 },
])

const videos = ref<VideoItem[]>([])
const isEditModalOpen = ref(false)
// Use ScheduleModalInputData for preparing data to pass to the modal
const selectedScheduleForModal = ref<ScheduleModalInputData>({
  videoName: '', // Ensure videoName is initialized as string
  time: '00:00',
  isActive: true,
  dayOrder: 0,
})

const departments = ref([
  { id: 'all', name: 'Все отделы', checked: true },
])
const selectedVideos = ref<string[]>([])
const allVideosChecked = ref(true)

const cronToTime = (cronExpression: string): string => {
  const parts = cronExpression.split(' ')
  if (parts.length >= 2) {
    const minute = parts[0]?.padStart(2, '0') || '00'
    const hour = parts[1]?.padStart(2, '0') || '00'
    return `${hour}:${minute}`
  }
  return '00:00'
}

const timeToCron = (time: string, dayOfWeek: number): string => {
  const [hourStr, minuteStr] = time.split(':')
  const minutes = minuteStr || '0'
  const hours = hourStr || '0'
  return `${minutes} ${hours} * * ${dayOfWeek}`
}

const getDayOfWeekFromCron = (cronExpression: string): number => {
  const parts = cronExpression.split(' ')
  if (parts.length >= 5 && parts[4] && parts[4] !== '*') {
    return parseInt(parts[4], 10)
  }
  return 0
}

const getVideoName = (videoId: string | undefined): string => {
  if (!videoId) return 'Неизвестное видео'
  const video = videos.value.find(v => v.ID === videoId)
  return video ? video.Name : 'Неизвестное видео'
}

onMounted(async () => {
  if (videos.value.length === 0) {
    videos.value = [
      { ID: "1", Name: "Утреннняя зарядка" },
      { ID: "2", Name: "Вечерняя йога" },
    ]
  }
  await loadSchedules()
})

const loadSchedules = async () => {
  try {
    days.value.forEach(day => day.items = [])
    await scheduleStore.fetchSchedules()
    const protoSchedules = scheduleStore.allSchedules

    protoSchedules.forEach(protoSchedule => {
      // Assuming protobuf-ts converts snake_case to camelCase
      const dayOrder = getDayOfWeekFromCron(protoSchedule.cronExpression)
      const scheduleItem: ScheduleItem = {
        ID: protoSchedule.id.toString(),
        Time: cronToTime(protoSchedule.cronExpression),
        VideoID: protoSchedule.videoId.toString(),
        Name: getVideoName(protoSchedule.videoId.toString()),
        is_active: protoSchedule.isActive, // camelCase
        dayOrder: dayOrder,
        cron_expression: protoSchedule.cronExpression, // Store original cron
      }

      const day = days.value.find(d => d.order === dayOrder)
      if (day) {
        day.items.push(scheduleItem)
      }
    })
  } catch (error) {
    console.error('Error loading schedules:', error)
    $q.notify({
      color: 'negative',
      message: 'Не удалось загрузить расписания',
      icon: 'error'
    });
  }
}

const onAdd = (dayOrder: number) => {
  selectedScheduleForModal.value = {
    // id is undefined for new item
    videoName: '', // Provide a default videoName or leave for modal to handle
    time: '00:00',
    isActive: true,
    dayOrder: dayOrder,
    // videoId is undefined, modal needs to select one
  }
  isEditModalOpen.value = true
}

// New handler for the @edit event from ScheduleCard
const handleEditEvent = (payload: ScheduleCardEditPayload) => {
  const day = days.value.find(d => d.order === payload.dayOrder)
  if (!day) return
  const itemToEdit = day.items.find(item => item.ID === payload.scheduleId)
  if (!itemToEdit) return

  selectedScheduleForModal.value = {
    id: itemToEdit.ID,
    videoId: itemToEdit.VideoID,
    videoName: getVideoName(itemToEdit.VideoID),
    time: itemToEdit.Time,
    isActive: itemToEdit.is_active,
    dayOrder: itemToEdit.dayOrder,
  }
  isEditModalOpen.value = true
}

const onDelete = async (scheduleId: string) => {
  try {
    await scheduleStore.deleteSchedule(scheduleId)
    await loadSchedules()
    $q.notify({
      color: 'positive',
      message: 'Расписание успешно удалено',
      icon: 'check_circle'
    });
  } catch (error) {
    console.error('Error deleting schedule:', error)
    $q.notify({
      color: 'negative',
      message: 'Не удалось удалить расписание',
      icon: 'error'
    });
  }
}

// onSave now expects the payload from ScheduleEditModalSavePayload
const onSave = async (data: ScheduleEditModalSavePayload) => {
  try {
    // Convert dayOrder from string (emitted by modal) to number for timeToCron
    const dayOrderNum = parseInt(data.dayOrder, 10)
    const cronExpression = timeToCron(data.time, dayOrderNum)

    if (!data.videoId) {
      console.error("Video ID is missing in modal data")
      $q.notify({
        color: 'negative',
        message: 'Необходимо выбрать видео',
        icon: 'error'
      });
      return
    }
    const videoIdBigInt = BigInt(data.videoId)

    const scheduleData: ProtoSchedule = {
      id: data.scheduleId ? BigInt(data.scheduleId) : 0n, // Use scheduleId from modal payload
      cronExpression: cronExpression,
      isActive: data.isActive,
      videoId: videoIdBigInt,
      adminId: 0n,
      createdAt: '',
    }

    if (data.scheduleId) { // If scheduleId (from modal) exists, it's an update
      await scheduleStore.updateSchedule(scheduleData, ['cron_expression', 'is_active', 'video_id'])
      $q.notify({
        color: 'positive',
        message: 'Расписание успешно обновлено',
        icon: 'check_circle'
      });
    } else {
      await scheduleStore.createSchedules([scheduleData])
      $q.notify({
        color: 'positive',
        message: 'Расписание успешно создано',
        icon: 'check_circle'
      });
    }

    await loadSchedules()
    isEditModalOpen.value = false
  } catch (error) {
    console.error('Error saving schedule:', error)
    $q.notify({
      color: 'negative',
      message: 'Не удалось сохранить расписание',
      icon: 'error'
    });
  }
}

const onDepartmentChange = (value: string) => {
  console.log('Department changed:', value)
}

const onVideoChange = (value: string) => {
  console.log('Video changed:', value)
}

const onDepartmentToggle = (id: string) => {
  console.log('Department toggled:', id)
}

const onVideoToggle = (id: string) => {
  console.log('Video toggled:', id)
}

const onAllVideosToggle = () => {
  console.log('All videos toggled')
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
