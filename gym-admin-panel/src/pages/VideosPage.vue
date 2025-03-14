<template>
  <div>
    <div class="content">
      <h1 class="page-title">Управление видео</h1>

      <div class="upload-section">
        <h2 style="font-size: 32px; color: rgba(90,92,105,1); margin-bottom: 20px;">Загрузить новое видео</h2>
        <div class="upload-area" @click="handleContainerClick">
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
          :disabled="isUploading"
          @click="uploadFiles"
        >
          Загрузить видео
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
        <div class="video-card" v-for="video in videos" :key="video.ID">
          <div class="video-thumbnail">
            <div class="video-duration">00:00</div>
          </div>
          <div class="video-info">
            <div class="video-title">{{ video.Name }}</div>
            <div class="video-meta">
              <span>Просмотры: 0</span>
              <span>Добавлено: сегодня</span>
            </div>
            <div>
              <span class="tag">Видео</span>
            </div>
            <div class="video-actions">
              <button class="action-btn edit" @click="isModalChangeNameVideo = true; index = video.ID">
                Редактировать
              </button>
              <button class="action-btn delete" @click="index = video.ID; removeVideo()">
                Удалить
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <q-dialog v-model="isModalChangeNameVideo">
    <q-card class="edit-modal">
      <div class="form-group">
        <label>Название видео</label>
        <input
          v-model="newVideoName"
          type="text"
          placeholder="Введите новое название видео"
        />
      </div>
      <div class="modal-actions">
        <button class="btn btn-outline" @click="isModalChangeNameVideo = false">
          Отмена
        </button>
        <button 
          class="btn btn-primary" 
          @click="changeVideoName(index); isModalChangeNameVideo = false"
        >
          Сохранить
        </button>
      </div>
    </q-card>
  </q-dialog>
</template>

<script lang="ts">
import { defineComponent, onMounted, ref } from 'vue';

export default defineComponent({
  name: 'VideoUploadPage',
  setup() {
    const fileInput = ref<HTMLInputElement | null>(null);
    const uploadedFiles = ref<File[]>([]);
    const isUploading = ref(false);
    const uploadProgress = ref(0);
    const errorMessage = ref<string>('');
    const videos = ref<{ID: number, Name: string}[]>([])
    const isModalChangeNameVideo = ref(false)
    const newVideoName = ref('')
    const index = ref()
    const API_URL = process.env.QUASAR_API_URL || 'http://localhost:8083/api/v1'

    /** Триггерим выбор файла */
    const triggerFileInput = () => {
      if (fileInput.value) {
        fileInput.value.click();
      }
    };

    /** Обработка клика по контейнеру */
    const handleContainerClick = (event: MouseEvent) => {
      /** Проверяем, что клик не был на кнопке "Отправить" или на крестике */
      if (event.target instanceof HTMLElement && !event.target.closest('.uploaded-file')) {
        triggerFileInput();
      }
    };

    /** Обработка выбранного файла */
    const handleFileChange = (event: Event) => {
      const files = (event.target as HTMLInputElement).files;
      if (files && files.length > 0) {
        const file = files[0]!;
        if (file.type !== 'video/mp4') {
          errorMessage.value = 'Только файлы в формате MP4.';
          return;
        }
        if (file.size > 500 * 1024 * 1024) {
          errorMessage.value = 'Файл должен быть не более 500 МБ.';
          return;
        }
        errorMessage.value = '';

        uploadedFiles.value.push(file);
      }
    };

    /** Удаление файла из списка */
    const removeFile = (index: number) => {
      uploadedFiles.value.splice(index, 1);
    };

    /** Отправка файлов на сервер */
    const uploadFiles = async () => {
      if (uploadedFiles.value.length === 0) return;

      isUploading.value = true;
      const formData = new FormData();
      uploadedFiles.value.forEach((file) => {
        formData.append('file', file);
      });

      const token = localStorage.getItem('token');
      if (!token) {
        errorMessage.value = 'Токен не найден. Пожалуйста, войдите в систему.';
        isUploading.value = false;
        return;
      }

      try {
        await fetch(`${API_URL}/videos`, {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${token}`,
          },
          body: formData
        });

        isUploading.value = false;
        uploadedFiles.value = [];
        init()
      } catch (error) {
        console.error(error);

        errorMessage.value = 'Ошибка при загрузке файла';
        isUploading.value = false;
      }
      isUploading.value = false;
    };
    const init = async () => {
      await fetch(`${API_URL}/videos`, {
        method: 'GET',
      }).then(async (res) => {
        if (res.ok) {
          const response = await res.json()
          videos.value = response
        }
      })
    }

    const changeVideoName = async (idx: number) => {
      const token = localStorage.getItem('token');
      await fetch(`${API_URL}/videos/${idx}`, {
        method: 'PUT',
        headers: {
          'Authorization': `Bearer ${token}`,
        }, body: JSON.stringify({ name: newVideoName.value })
      }).then(async (res) => {
        if (res.ok) {
          init()
        }
      })
    }

    const removeVideo = async () => {
      const token = localStorage.getItem('token');
      await fetch(`${API_URL}/videos/${index.value}`, {
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

    onMounted(async () => {
      init()
    })

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
      newVideoName,
      changeVideoName,
      index,
      removeVideo
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

.action-btn.edit {
  background-color: rgba(78,115,223,1);
  color: white;
}

.action-btn.delete {
  background-color: rgba(231,74,59,1);
  color: white;
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
</style>
