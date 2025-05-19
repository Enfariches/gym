import { defineStore } from 'pinia';
import {
  uploadMedia,
  validateVideoFile as validateVideoFileService,
  getVideoPresignedUrl,
  listMediaGrpc,
  deleteMediaGrpc
} from '../services/mediaService';

export interface Video {
  ID: number;
  Name: string;
  previewUrl?: string;
  id?: bigint; // для gRPC
  pressignedUrl?: string; // для gRPC
  createdAt?: string; // для gRPC
}

export const useMediaStore = defineStore('media', {
  state: () => ({
    videos: [] as Video[],
    uploadedFiles: [] as File[],
    isUploading: false,
    uploadProgress: 0,
    error: null as string | null,
    departmentId: 1, // Дефолтное значение, должно быть заменено на актуальное
    useGrpc: true, // переключатель для примера
  }),

  actions: {
    // Загрузка списка видео через gRPC
    async loadVideos() {
      this.error = null;
      try {
        const grpcVideos = await listMediaGrpc();
        this.videos = grpcVideos.map(v => ({
          ID: Number(v.id),
          Name: v.pressignedUrl || '', // или другое поле, если есть имя
          id: v.id,
          pressignedUrl: v.pressignedUrl,
          createdAt: v.createdAt
        }));
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Не удалось загрузить список видео';
        console.error('Ошибка при загрузке видео:', error);
      }
    },

    // Удаление видео через gRPC
    async deleteVideo(videoId: number) {
      this.error = null;
      try {
        await deleteMediaGrpc(videoId);
        await this.loadVideos();
        return true;
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Не удалось удалить видео';
        console.error('Ошибка при удалении видео:', error);
        return false;
      }
    },

    // Получение пресайн URL для превью видео (gRPC)
    async getVideoPreview(videoId: number) {
      try {
        const videoIndex = this.videos.findIndex(v => v.ID === videoId);
        if (videoIndex === -1) return;
        const video = this.videos[videoIndex];
        if (!video) return;
        const previewUrl = await getVideoPresignedUrl(videoId, this.departmentId);
        this.videos[videoIndex] = {
          ...video,
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

    // Загрузка файла на сервер (HTTP)
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
    getVideos: (state) => state.videos,
    hasFiles: (state) => state.uploadedFiles.length > 0,
    getUploadedFiles: (state) => state.uploadedFiles,
    getUploadStatus: (state) => ({
      isUploading: state.isUploading,
      progress: state.uploadProgress,
      error: state.error
    }),
    getCurrentDepartmentId: (state) => state.departmentId
  }
});
