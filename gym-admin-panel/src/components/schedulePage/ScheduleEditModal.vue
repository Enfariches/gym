<template>
  <q-dialog v-model="localIsOpen">
    <q-card class="edit-modal">
      <h3 class="modal-title">{{ scheduleId ? 'Редактировать расписание' : 'Добавить расписание' }}</h3>
      <div class="form-group">
        <label>Выберите видео</label>
        <q-select
          v-model="selectedVideoName"
          class="video-select"
          :options="videos.map(el => el.Name)"
          @update:model-value="onSelectedVideoChange"
          :display-value="selectedVideoName || 'Выберите видео'"
        />
      </div>
      <div class="form-group">
        <label>Выберите время показа</label>
        <input
          type="text"
          v-model="time"
          class="time-input"
          placeholder="ЧЧ:ММ (например, 09:00)"
          @input="validateTime"
        />
        <div v-if="timeError" class="error-message">
          {{ timeError }}
        </div>
      </div>
      <div class="modal-actions">
        <button class="btn btn-outline" @click="closeModal">Отмена</button>
        <button class="btn btn-primary" @click="onSave" :disabled="!!timeError">Сохранить</button>
      </div>
    </q-card>
  </q-dialog>
</template>

<script lang="ts" setup>
import { ref, watch } from 'vue';

interface Video {
  ID: string;
  Name: string;
}

const props = defineProps<{
  isOpen: boolean;
  videos: Video[];
  scheduleId?: string;
  initialVideoName?: string;
  initialTime?: string;
  dayOrder: string;
}>();

const emit = defineEmits<{
  (e: 'update:isOpen', value: boolean): void;
  (e: 'save', data: { videoId: string; time: string; dayOrder: string; scheduleId: string | undefined }): void;
}>();

const localIsOpen = ref(props.isOpen);

const selectedVideoName = ref(props.initialVideoName || '');
const time = ref(props.initialTime || '');
const timeError = ref('');

watch(() => props.isOpen, (newValue) => {
  localIsOpen.value = newValue;
  if (newValue) {
    selectedVideoName.value = props.initialVideoName || '';
    time.value = props.initialTime || '';
    timeError.value = '';
  }
});

const validateTime = () => {
  const timeRegex = /^([0-1]?[0-9]|2[0-3]):[0-5][0-9]$/;
  if (!timeRegex.test(time.value)) {
    timeError.value = 'Неверный формат времени. Используйте формат ЧЧ:ММ';
  } else {
    timeError.value = '';
  }
};

const onSelectedVideoChange = (newValue: string) => {
  selectedVideoName.value = newValue;
};

const closeModal = () => {
  localIsOpen.value = false;
  emit('update:isOpen', false);
};

const onSave = () => {
  if (timeError.value) return;

  const video = props.videos.find(v => v.Name === selectedVideoName.value);
  if (!video) return;

  emit('save', {
    videoId: video.ID,
    time: time.value,
    dayOrder: props.dayOrder,
    scheduleId: props.scheduleId || undefined
  });

  closeModal();
};
</script>

<style scoped>
.edit-modal {
  min-width: 400px;
  padding: 20px;
}

.modal-title {
  margin-top: 0;
  margin-bottom: 20px;
  font-size: 1.5em;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
}

.video-select {
  width: 100%;
}

.time-input {
  width: 100%;
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.error-message {
  color: #F44336;
  font-size: 0.8em;
  margin-top: 4px;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 20px;
}

.btn {
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  font-weight: 500;
}

.btn-outline {
  background: none;
  border: 1px solid #ddd;
}

.btn-primary {
  background: #2196F3;
  color: white;
  border: none;
}

.btn-primary:disabled {
  background: #ccc;
  cursor: not-allowed;
}
</style> 