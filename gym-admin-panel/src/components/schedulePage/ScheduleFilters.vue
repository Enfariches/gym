<template>
  <div class="schedule-toolbar">
    <div class="filter-container">
      <select class="btn btn-outline" @change="$emit('departmentChange', ($event.target as HTMLInputElement).value)">
        <option value="">Все отделы</option>
        <option value="dev">Разработчики</option>
        <option value="finance">Финансовый отдел</option>
        <option value="marketing">Маркетинг</option>
        <option value="hr">HR</option>
      </select>
      <select class="btn btn-outline" @change="$emit('videoChange', ($event.target as HTMLInputElement).value)">
        <option value="">Все видео</option>
        <option v-for="video in videos" :key="video.title" :value="video.title">
          {{ video.title }}
        </option>
      </select>
    </div>
  </div>
</template>

<script lang="ts" setup>
defineProps<{
  videos: Array<{
    title: string;
  }>;
}>();

defineEmits<{
  (e: 'departmentChange', value: string): void;
  (e: 'videoChange', value: string): void;
}>();
</script>

<style scoped>
.schedule-toolbar {
  margin-bottom: 20px;
}

.filter-container {
  display: flex;
  gap: 10px;
}

select.btn-outline {
  padding: 8px 12px;
  border-radius: 4px;
  border: 1px solid #ddd;
  background-color: white;
  min-width: 200px;
}
</style>
