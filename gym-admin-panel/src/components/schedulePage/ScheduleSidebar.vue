<template>
  <div class="schedule-sidebar">
    <h3 class="sidebar-title">Отделы</h3>
    <div class="department-list">
      <div class="department-item" v-for="dept in departments" :key="dept.id">
        <label class="checkbox-container">
          <input 
            type="checkbox" 
            :checked="dept.checked"
            @change="$emit('departmentToggle', dept.id)"
          >
          <span class="checkmark"></span>
          <span class="department-name">{{ dept.name }}</span>
        </label>
      </div>
    </div>
    
    <h3 class="sidebar-title">Видео</h3>
    <div class="department-list">
      <div class="department-item">
        <label class="checkbox-container">
          <input 
            type="checkbox" 
            :checked="allVideosChecked"
            @change="$emit('allVideosToggle')"
          >
          <span class="checkmark"></span>
          <span class="department-name">Все видео</span>
        </label>
      </div>
      <div v-for="video in videos" :key="video.ID" class="department-item">
        <label class="checkbox-container">
          <input 
            type="checkbox"
            :checked="selectedVideos.includes(video.ID)"
            @change="$emit('videoToggle', video.ID)"
          >
          <span class="checkmark"></span>
          <span class="department-name">{{ video.Name }}</span>
        </label>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
defineProps<{
  departments: Array<{
    id: string;
    name: string;
    checked: boolean;
  }>;
  videos: Array<{
    ID: string;
    Name: string;
  }>;
  selectedVideos: string[];
  allVideosChecked: boolean;
}>();

defineEmits<{
  (e: 'departmentToggle', id: string): void;
  (e: 'videoToggle', id: string): void;
  (e: 'allVideosToggle'): void;
}>();
</script>

<style scoped>
.schedule-sidebar {
  width: 250px;
  padding: 20px;
  background-color: #f5f5f5;
  border-radius: 8px;
}

.sidebar-title {
  font-size: 1.2em;
  margin-bottom: 15px;
  color: #333;
}

.department-list {
  margin-bottom: 20px;
}

.department-item {
  margin-bottom: 10px;
}

.checkbox-container {
  display: flex;
  align-items: center;
  cursor: pointer;
}

.checkbox-container input {
  margin-right: 8px;
}

.department-name {
  font-size: 0.9em;
  color: #666;
}

.checkmark {
  margin-right: 8px;
}
</style> 