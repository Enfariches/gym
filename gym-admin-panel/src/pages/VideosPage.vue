<template>
  <div>
    <div class="content">
      <h1 class="page-title">Управление видео</h1>
      <div class="upload-section">
        <h2 style="font-size: 32px; color: rgba(90,92,105,1); margin-bottom: 20px;">Загрузить новое видео</h2>
        <div class="upload-area" @click="handleContainerClick" :class="{ 'drag-over': isDragging }">
          <div class="upload-icon">+</div>
          <div class="upload-text">Перетащите файл видео сюда или нажмите для выбора</div>
          <button class="btn btn-primary">Выбрать файл</button>
          <input
            type="file"
            ref="fileInput"
            accept="video/mp4"
            @change="handleFileChange"
            style="display: none"
          />
        </div>

        <div v-if="errorMessage" class="error-message q-mt-md">
          {{ errorMessage }}
        </div>

        <div v-if="isUploading" class="progress-container q-mt-md">
          <div class="progress-label">Загрузка: {{ Math.round(uploadProgress * 100) }}%</div>
          <q-linear-progress
            :value="uploadProgress"
            color="primary"
            size="15px"
          />
        </div>

        <div v-if="uploadedFiles.length" class="uploaded-files q-mt-lg">
          <div
            v-for="(file, index) in uploadedFiles"
            :key="index"
            class="uploaded-file"
          >
            <span>{{ file.name }}</span>
            <q-btn
              flat
              icon="cancel"
              color="negative"
              size="sm"
              @click.stop="removeFile(index)"
            />
          </div>
        </div>

        <button
          class="btn btn-primary"
          style="margin-top: 20px;"
          :disabled="isUploading || uploadedFiles.length === 0"
          @click="uploadFiles"
        >
          {{ isUploading ? 'Загрузка...' : 'Загрузить видео' }}
        </button>
      </div>

      <div class="video-toolbar">
        <div class="search-container">
          <input type="text" class="search-input" placeholder="Поиск видео по названию, категории...">
        </div>
        <div class="filter-container">
          <select class="btn btn-outline">
            <option>Все категории</option>
            <option>Утренняя</option>
            <option>Офисная</option>
            <option>Спина</option>
            <option>Глаза</option>
          </select>
          <select class="btn btn-outline">
            <option>Сортировка</option>
            <option>По просмотрам</option>
            <option>По дате</option>
            <option>По названию</option>
          </select>
        </div>
      </div>

      <div class="video-grid" v-if="videos.length">
        <VideoCard
          v-for="video in videos"
          :key="video.title"
          :video="{
            ID: Number(video.id || 0),
            title: video.title,
            pressignedUrl: video.pressignedUrl ?? ''
          }"
          @edit="isModalChangeNameVideo = true; index = $event"
          @delete="index = $event; removeVideo()"
        />
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, onMounted, ref } from 'vue';
import { useQuasar } from 'quasar';
import VideoCard from 'src/components/videosPage/VideoCard.vue';
import { useMediaStore } from 'src/stores/mediaStore';
import { storeToRefs } from 'pinia';

export default defineComponent({
  name: 'VideoUploadPage',
  components: {
    VideoCard
  },
  setup() {
    const $q = useQuasar();
    const fileInput = ref<HTMLInputElement | null>(null);
    const errorMessage = ref<string>('');
    const isModalChangeNameVideo = ref(false);
    const index = ref(0);
    const isDragging = ref(false);
    const mediaStore = useMediaStore();
    const {
      videos,
      uploadedFiles,
      isUploading,
      uploadProgress
    } = storeToRefs(mediaStore);

    /** Триггерим выбор файла */
    const triggerFileInput = () => {
      if (fileInput.value) {
        fileInput.value.click();
      }
    };

    /** Обработка клика по контейнеру */
    const handleContainerClick = (event: MouseEvent) => {
      if (event.target instanceof HTMLElement && !event.target.closest('.uploaded-file')) {
        triggerFileInput();
      }
    };

    /** Обработка выбранного файла */
    const handleFileChange = (event: Event) => {
      const files = (event.target as HTMLInputElement).files;
      if (files && files.length > 0) {
        const file = files[0];
        if (file) {
          const success = mediaStore.addFile(file);
          if (!success && mediaStore.error) {
            errorMessage.value = mediaStore.error;
          } else {
            errorMessage.value = '';
          }
        }
      }
    };

    /** Удаление файла из списка */
    const removeFile = (index: number) => {
      mediaStore.removeFile(index);
    };

    /** Отправка файлов на сервер */
    const uploadFiles = async () => {
      if (uploadedFiles.value.length === 0) {
        errorMessage.value = 'Пожалуйста, выберите файл для загрузки.';
        return;
      }
      try {
        const success = await mediaStore.uploadFile();
        if (success) {
          $q.notify({ type: 'positive', message: 'Видео успешно загружено', position: 'top', timeout: 2000 });
          errorMessage.value = '';
        } else if (mediaStore.error) {
          errorMessage.value = mediaStore.error;
          $q.notify({ type: 'negative', message: mediaStore.error, position: 'top', timeout: 3000 });
        }
      } catch (error) {
        console.error('Ошибка при загрузке видео:', error);
        const errorMsg = error instanceof Error ? error.message : 'Ошибка при загрузке файла';
        errorMessage.value = errorMsg;
        $q.notify({ type: 'negative', message: errorMsg, position: 'top', timeout: 3000 });
      }
    };

    const removeVideo = async () => {
      try {
        const success = await mediaStore.deleteVideo(index.value);
        if (success) {
          $q.notify({ type: 'positive', message: 'Видео успешно удалено', position: 'top', timeout: 2000 });
        } else if (mediaStore.error) {
          $q.notify({ type: 'negative', message: mediaStore.error, position: 'top', timeout: 3000 });
        }
      } catch (error) {
        console.error('Ошибка при удалении видео:', error);
        $q.notify({ type: 'negative', message: 'Не удалось удалить видео', position: 'top', timeout: 3000 });
      }
    };

    // Поддержка Drag and Drop
    const handleDragOver = (event: DragEvent) => {
      event.preventDefault();
      isDragging.value = true;
    };
    const handleDragLeave = () => { isDragging.value = false; };
    const handleDrop = (event: DragEvent) => {
      event.preventDefault();
      isDragging.value = false;
      const dataTransfer = event.dataTransfer;
      if (!dataTransfer?.files || dataTransfer.files.length === 0) return;
      const file = dataTransfer.files[0];
      if (!file) return;
      const success = mediaStore.addFile(file);
      if (!success && mediaStore.error) {
        errorMessage.value = mediaStore.error;
      } else {
        errorMessage.value = '';
      }
    };


    onMounted(() => {
      mediaStore.loadVideos();
      const uploadArea = document.querySelector('.upload-area');
      if (uploadArea) {
        uploadArea.addEventListener('dragover', (event: Event) => {
          handleDragOver(event as DragEvent);
        });
        uploadArea.addEventListener('dragleave', () => {
          handleDragLeave();
        });
        uploadArea.addEventListener('drop', (event: Event) => {
          handleDrop(event as DragEvent);
        });
      }
    });

    return {
      fileInput,
      uploadedFiles,
      isUploading,
      uploadProgress,
      errorMessage,
      triggerFileInput,
      handleFileChange,
      removeFile,
      uploadFiles,
      handleContainerClick,
      videos,
      isModalChangeNameVideo,
      index,
      removeVideo,
      isDragging,
      mediaStore,
    };
  },
});
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

.video-toolbar {
  background: white;
  border: 2px solid rgba(227,230,240,1);
  border-radius: 6px;
  padding: 15px;
  margin-bottom: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-container {
  flex: 0 0 60%;
  position: relative;
}

.search-input {
  width: 100%;
  padding: 8px 15px;
  font-size: 14px;
  border: 1px solid rgba(227,230,240,1);
  border-radius: 4px;
}

.btn {
  padding: 8px 15px;
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
  border: 1px solid rgba(78,115,223,1);
  color: rgba(78,115,223,1);
}

.filter-container {
  display: flex;
  gap: 10px;
}

.video-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
}

.upload-section {
  background: white;
  border: 2px solid rgba(227,230,240,1);
  border-radius: 6px;
  padding: 20px;
  margin-bottom: 20px;
  text-align: center;
}

.upload-section h2 {
  font-size: 24px !important;
  color: rgba(90,92,105,1);
  margin-bottom: 15px !important;
}

.upload-area {
  border: 2px dashed rgba(78,115,223,0.5);
  border-radius: 6px;
  padding: 20px;
  margin: 15px 0;
  background-color: rgba(78,115,223,0.05);
  cursor: pointer;
}

.upload-icon {
  font-size: 40px;
  color: rgba(78,115,223,0.7);
  margin-bottom: 15px;
}

.upload-text {
  font-size: 16px;
  color: rgba(90,92,105,1);
  margin-bottom: 15px;
}

.edit-modal {
  min-width: 300px;
  padding: 20px;
  background: white;
}

.form-group {
  margin-bottom: 15px;
  text-align: left;
}

.form-group label {
  display: block;
  font-size: 14px;
  font-weight: bold;
  color: rgba(90,92,105,1);
  margin-bottom: 6px;
}

.form-group input {
  width: 100%;
  padding: 8px;
  font-size: 14px;
  border: 1px solid rgba(227,230,240,1);
  border-radius: 4px;
}

.modal-actions {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
}

.error-message {
  color: red;
  font-size: 12px;
  margin-top: 8px;
}

.progress-container {
  margin-top: 15px;
}

.uploaded-files {
  margin-top: 15px;
  text-align: left;
}

.uploaded-file {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 6px;
  background-color: rgba(78,115,223,0.1);
  border-radius: 4px;
  margin-bottom: 8px;
  font-size: 12px;
}

.progress-label {
  font-size: 14px;
  color: rgba(78,115,223,1);
  margin-bottom: 5px;
  text-align: right;
}

.drag-over {
  border-color: rgba(78,115,223,1);
  background-color: rgba(78,115,223,0.1);
}
</style>
