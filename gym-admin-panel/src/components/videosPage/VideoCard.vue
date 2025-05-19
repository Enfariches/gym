<template>
  <div class="video-card">
    <div class="video-thumbnail">
      <div class="video-duration">00:00</div>
    </div>
    <div class="video-info">
      <div class="video-title">{{ video.title }}</div>
      <div class="video-meta">
        <span>Просмотры: 0</span>
        <span>Добавлено: сегодня</span>
      </div>
      <div>
        <span class="tag">Видео</span>
      </div>
      <div class="video-actions">
        <button class="action-btn delete" @click="$emit('delete', video.ID)">
          Удалить
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps, defineEmits } from 'vue';
// eslint-disable-next-line @typescript-eslint/no-unused-vars
const props = defineProps<{
  video: {
    ID: number;
    title: string;
    pressignedUrl?: string;
  };
}>();
// eslint-disable-next-line @typescript-eslint/no-unused-vars
const emit = defineEmits(['delete']);
</script>

<style scoped>
.video-card {
  background: white;
  border: 2px solid rgba(227,230,240,1);
  border-radius: 6px;
  overflow: hidden;
}

.video-thumbnail {
  width: 100%;
  height: 160px;
  background-color: #ccc;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 14px;
  position: relative;
}

.video-duration {
  position: absolute;
  bottom: 8px;
  right: 8px;
  background: rgba(0,0,0,0.7);
  color: white;
  padding: 4px 8px;
  border-radius: 3px;
  font-size: 12px;
}

.video-info {
  padding: 15px;
}

.video-title {
  font-size: 16px;
  font-weight: bold;
  color: rgba(90,92,105,1);
  margin-bottom: 8px;
}

.video-meta {
  display: flex;
  justify-content: space-between;
  color: rgba(108,117,125,1);
  margin-bottom: 10px;
  font-size: 12px;
}

.tag {
  display: inline-block;
  background-color: rgba(78,115,223,0.1);
  color: rgba(78,115,223,1);
  font-size: 12px;
  padding: 4px 8px;
  border-radius: 12px;
  margin-right: 6px;
  margin-bottom: 6px;
}

.video-actions {
  display: flex;
  justify-content: space-between;
  margin-top: 10px;
  gap: 8px;
}

.action-btn {
  padding: 6px 12px;
  font-size: 12px;
  border-radius: 4px;
  border: none;
  cursor: pointer;
  flex: 1;
}

.action-btn.delete {
  background-color: rgba(231,74,59,1);
  color: white;
}
</style>
