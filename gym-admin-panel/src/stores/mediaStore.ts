import { defineStore } from 'pinia';
import {
  fetchVideos,
  uploadMedia,
  changeVideoName as changeVideoNameService,
  deleteVideo as deleteVideoService,
  validateVideoFile as validateVideoFileService,
  getVideoPresignedUrl
} from '../services/mediaService';

export interface Video {
  ID: number;
  Name: string;
  previewUrl?: string;
}

export const useMediaStore = defineStore('media', {
  state: () => ({
    videos: [] as Video[],
    uploadedFiles: [] as File[],
    isUploading: false,
    uploadProgress: 0,
    error: null as string | null,
    departmentId: 1, // Дефолтное значение, должно быть заменено на актуальное
  }),

  actions: {
    // Загрузка списка видео
    async loadVideos() {
      this.error = null;
      try {
        const fetchedVideos = await fetchVideos();
        this.videos = fetchedVideos;
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Не удалось загрузить список видео';
        console.error('Ошибка при загрузке видео:', error);
      }
    },

    // Получение пресайн URL для превью видео
    async getVideoPreview(videoId: number) {
      try {
        const videoIndex = this.videos.findIndex(v => v.ID === videoId);
        if (videoIndex === -1) return;

        // Избегаем ошибок с undefined, проверяя наличие видео
        const video = this.videos[videoIndex];
        if (!video) return;

        const previewUrl = await getVideoPresignedUrl(videoId, this.departmentId);

        // Обновляем видео с URL для превью
        this.videos[videoIndex] = {
          ID: video.ID,
          Name: video.Name,
          previewUrl
        };

        return previewUrl;
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Не удалось получить превью видео';
        console.error('Ошибка при получении превью видео:', error);
        return null;
      }
    },

    // Добавление файла в список для загрузки
    addFile(file: File) {
      const validationResult = validateVideoFileService(file);
      if (!validationResult.isValid) {
        this.error = validationResult.error || 'Файл не соответствует требованиям';
        return false;
      }

      this.error = null;
      this.uploadedFiles = [file]; // Сохраняем только один файл
      return true;
    },

    // Удаление файла из списка для загрузки
    removeFile(index: number) {
      this.uploadedFiles.splice(index, 1);
    },

    // Очистка списка файлов
    clearFiles() {
      this.uploadedFiles = [];
    },

    // Загрузка файла на сервер
    async uploadFile() {
      if (this.uploadedFiles.length === 0) {
        this.error = 'Пожалуйста, выберите файл для загрузки.';
        return false;
      }

      this.isUploading = true;
      this.uploadProgress = 0;
      this.error = null;

      try {
        const file = this.uploadedFiles[0];

        if (file) {
          await uploadMedia(file, (progress) => {
            this.uploadProgress = progress;
          });

          this.uploadProgress = 1;
          this.uploadedFiles = [];

          // Обновляем список видео после успешной загрузки
          await this.loadVideos();

          return true;
        } else {
          throw new Error('Файл не найден');
        }
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Ошибка при загрузке файла';
        console.error('Ошибка при загрузке видео:', error);
        return false;
      } finally {
        this.isUploading = false;
      }
    },

    // Изменение названия видео
    async changeVideoName(videoId: number, newName: string) {
      this.error = null;

      try {
        await changeVideoNameService(videoId, newName);
        // Обновляем список видео после изменения названия
        await this.loadVideos();
        return true;
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Не удалось изменить название видео';
        console.error('Ошибка при изменении названия:', error);
        return false;
      }
    },

    // Удаление видео
    async deleteVideo(videoId: number) {
      this.error = null;

      try {
        await deleteVideoService(videoId);
        // Обновляем список видео после удаления
        await this.loadVideos();
        return true;
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Не удалось удалить видео';
        console.error('Ошибка при удалении видео:', error);
        return false;
      }
    },

    // Задание ID департамента
    setDepartmentId(departmentId: number) {
      this.departmentId = departmentId;
    },

    // Сброс ошибки
    clearError() {
      this.error = null;
    }
  },

  getters: {
    // Получить список всех видео
    getVideos: (state) => state.videos,

    // Проверить, есть ли файлы для загрузки
    hasFiles: (state) => state.uploadedFiles.length > 0,

    // Получить загруженные файлы
    getUploadedFiles: (state) => state.uploadedFiles,

    // Статус загрузки
    getUploadStatus: (state) => ({
      isUploading: state.isUploading,
      progress: state.uploadProgress,
      error: state.error
    }),

    // Текущий департамент
    getCurrentDepartmentId: (state) => state.departmentId
  }
});
