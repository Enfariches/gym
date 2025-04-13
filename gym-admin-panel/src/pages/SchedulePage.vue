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
      v-model:isOpen="isEditModalOpen"
      :videos="videos"
      :scheduleId="index"
      :initialVideoName="newVideoName"
      :initialTime="timeWithSeconds"
      :dayOrder="videoToSchedule.dayofweek"
      @save="onSave"
    />
  </div>
</template>

<script lang="ts" setup>
import { onMounted, ref } from 'vue'
import ScheduleFilters from '../components/schedulePage/ScheduleFilters.vue'
import ScheduleSidebar from '../components/schedulePage/ScheduleSidebar.vue'
import ScheduleCard from '../components/schedulePage/ScheduleCard.vue'
import ScheduleEditModal from '../components/schedulePage/ScheduleEditModal.vue'

const mondayVideos: { ID: string, Time: string, VideoID: string, Name: string }[] = []
const tuesdayVideos: { ID: string, Time: string, VideoID: string, Name: string }[] = []
const wednesdayVideos: { ID: string, Time: string, VideoID: string, Name: string }[] = []
const thursdayVideos: { ID: string, Time: string, VideoID: string, Name: string }[] = []
const fridayVideos: { ID: string, Time: string, VideoID: string, Name: string }[] = []
const saturdayVideos: { ID: string, Time: string, VideoID: string, Name: string }[] = []
const sundayVideos: { ID: string, Time: string, VideoID: string, Name: string }[] = []

const days = ref([
  { items: mondayVideos, name: "Понедельник", bg: 'bg-primary', order:0 },
  { items: tuesdayVideos, name: "Вторник", bg: 'bg-primary',order:1 },
  { items: wednesdayVideos, name: "Среда", bg: 'bg-primary',order:2 },
  { items: thursdayVideos, name: "Четверг", bg: 'bg-primary',order:3 },
  { items: fridayVideos, name: "Пятница", bg: 'bg-primary',order:4 },
  { items: saturdayVideos, name: "Суббота", bg: 'bg-primary',order:5 },
  { items: sundayVideos, name: "Воскресенье", bg: 'bg-primary',order:6 },
])
const videos = ref()
const index = ref()
const isEditModalOpen = ref(false)
const newVideoName = ref('')
const timeWithSeconds = ref('')
const videoToSchedule = ref({
  dayofweek: '',
  time: '',
  videoid: ''
})
const timeError = ref('')

const departments = ref([
  { id: 'all', name: 'Все отделы', checked: true },
  { id: 'dev', name: 'Разработчики', checked: false },
  { id: 'finance', name: 'Финансовый отдел', checked: false },
  { id: 'marketing', name: 'Маркетинг', checked: false },
  { id: 'hr', name: 'HR', checked: false },
])

const selectedVideos = ref<string[]>([])
const allVideosChecked = ref(true)

const API_URL = process.env.QUASAR_API_URL || 'http://localhost:8083/api/v1'

const init = () => {
  const token = localStorage.getItem('token');

  days.value.forEach(async (day, idx) => {
    setTimeout(async () => {
      await fetch(`${API_URL}/schedule/${idx}`, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${token}`,
        },
        }).then(async (res) => {
          const response = await res.json()
          day.items = response
        })
    }, idx * 50)
  })
}

onMounted(async () => {
  init()
  await fetch(`${API_URL}/videos`, {
    method: 'GET',
  }).then(async (res) => {
    if (res.ok) {
      const response = await res.json()
      videos.value = response;
    }
  })
})

const _removeSchedule = async () => {
  const token = localStorage.getItem('token');
  await fetch(`${API_URL}/schedule/${index.value}`, {
    method: 'DELETE',
    headers: {
      'Authorization': `Bearer ${token}`,
    }
  }).then(async (res) => {
    if (res.ok) {
      init()
    }
  })
}

const updateSchedule = async () =>{
  const token = localStorage.getItem('token');

  if(videoToSchedule.value.videoid === '' || videoToSchedule.value.time === '' || videoToSchedule.value.dayofweek === '')
    return

  await fetch(`${API_URL}/schedule`, {
    method: 'PUT',
    headers: {
      'Authorization': `Bearer ${token}`,
    },
    body: JSON.stringify({
      dayofweek: Number(videoToSchedule.value.dayofweek),
      time: videoToSchedule.value.time,
      videoid: Number(videoToSchedule.value.videoid)
    }),
  }).then(async (res) => {
    if (res.ok) {
      init()
    }
  })
}

const validateTime = () => {
  const timeRegex = /^([0-1]?[0-9]|2[0-3]):[0-5][0-9]$/
  
  if (!timeWithSeconds.value) {
    timeError.value = 'Время не может быть пустым'
    return false
  }

  if (!timeRegex.test(timeWithSeconds.value)) {
    timeError.value = 'Неверный формат времени. Используйте формат ЧЧ:ММ (например, 09:00)'
    return false
  }

  timeError.value = ''
  videoToSchedule.value.time = timeWithSeconds.value
  return true
}

const onSave = () => {
  if (!validateTime()) return
  updateSchedule().then(
    () => onDecline()
  )
}

const onDecline = () => {
  videoToSchedule.value.dayofweek =''
  videoToSchedule.value.time =''
  videoToSchedule.value.videoid =''
  isEditModalOpen.value = false
  newVideoName.value=''
  timeWithSeconds.value = ''
}

const onEdit = ({ videoName, videoId, scheduleId, dayOrder }: { videoName: string; videoId: string; scheduleId: string; dayOrder: number }) => {
  newVideoName.value = videoName;
  videoToSchedule.value.videoid = videoId;
  videoToSchedule.value.dayofweek = dayOrder.toString();
  index.value = scheduleId;
  isEditModalOpen.value = true;
}

const onAdd = (dayOrder: number) => {
  videoToSchedule.value.dayofweek = dayOrder.toString();
  isEditModalOpen.value = true;
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
