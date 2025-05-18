import { MediaServiceClient } from '../../protogen/v1/media/media.client';
import type { GetMediaRequest } from '../../protogen/v1/media/media';
import type { RpcTransport } from '@protobuf-ts/runtime-rpc';

// API URL configuration
const API_URL = process.env.QUASAR_API_URL || 'http://localhost:8083/api/v1';
const UPLOAD_URL = 'http://localhost:3000/api/upload';

/**
 * Создает клиент для gRPC-сервиса MediaService
 */
const createGrpcClient = () => {
  // Здесь должна быть корректная инициализация транспорта для gRPC
  const transport = { /* ... */ } as RpcTransport;
  return new MediaServiceClient(transport);
};

/**
 * Загружает видео файл на сервер с отслеживанием прогресса
 */
export const uploadMedia = async (
  file: File,
  onProgress?: (progress: number) => void
): Promise<void> => {
  const token = localStorage.getItem('auth_token');
  if (!token) {
    throw new Error('Требуется авторизация. Пожалуйста, войдите в систему.');
  }

  const formData = new FormData();
  formData.append('mediafile', file);

  // Создаем Promise для управления асинхронной загрузкой
  return new Promise((resolve, reject) => {
    const xhr = new XMLHttpRequest();

    // Отслеживаем прогресс загрузки
    if (onProgress) {
      xhr.upload.addEventListener('progress', (event) => {
        if (event.lengthComputable) {
          onProgress(event.loaded / event.total);
        }
      });
    }

    xhr.open('POST', UPLOAD_URL);
    xhr.setRequestHeader('Authorization', `Bearer ${token}`);

    // Важно: не устанавливаем Content-Type заголовок
    // браузер сам установит его с правильным boundary для multipart/form-data

    // Разрешаем отправлять куки для поддержки CORS
    xhr.withCredentials = true;

    xhr.onload = () => {
      if (xhr.status >= 200 && xhr.status < 300) {
        resolve();
      } else {
        let errorMsg = 'Ошибка загрузки';

        switch (xhr.status) {
          case 401:
            errorMsg = 'Необходима авторизация. Пожалуйста, войдите снова.';
            break;
          case 403:
            errorMsg = 'Доступ запрещен. У вас нет прав для загрузки файлов.';
            break;
          case 413:
            errorMsg = 'Файл слишком большой. Максимальный размер 200МБ.';
            break;
          case 415:
            errorMsg = 'Неподдерживаемый формат файла. Используйте только MP4.';
            break;
          case 500:
            errorMsg = 'Внутренняя ошибка сервера. Пожалуйста, попробуйте позже.';
            break;
          case 503:
            errorMsg = 'Сервис недоступен. Пожалуйста, попробуйте позже.';
            break;
          default:
            errorMsg = `Ошибка загрузки: ${xhr.statusText || 'Неизвестная ошибка'}`;
        }

        reject(new Error(errorMsg));
      }
    };

    xhr.ontimeout = () => {
      reject(new Error('Превышено время ожидания. Проверьте подключение к сети.'));
    };

    xhr.onerror = () => {
      reject(new Error('Ошибка сети при загрузке файла. Проверьте подключение и доступность сервера.'));
    };

    // Устанавливаем таймаут в 30 секунд
    xhr.timeout = 30000;

    xhr.send(formData);
  });
};

/**
 * Получает список всех видео
 */
export const fetchVideos = async (): Promise<{ID: number, Name: string}[]> => {
  const response = await fetch(`${API_URL}/videos`, {
    method: 'GET',
  });

  if (!response.ok) {
    throw new Error('Не удалось получить список видео');
  }

  return await response.json();
};

/**
 * Получает пресайн URL для видео с использованием gRPC
 */
export const getVideoPresignedUrl = async (
  mediaId: number,
  departmentId: number,
  expirySeconds: number = 3600
): Promise<string> => {
  const client = createGrpcClient();

  const request: GetMediaRequest = {
    mediaId: BigInt(mediaId),
    departamentId: BigInt(departmentId),
    expiry: { seconds: BigInt(expirySeconds), nanos: 0 }
  };

  try {
    const { response } = await client.getMedia(request);
    return response.pressignedUrl;
  } catch (error) {
    console.error('Error fetching presigned URL:', error);
    throw new Error('Не удалось получить ссылку для просмотра видео');
  }
};

/**
 * Изменяет название видео
 */
export const changeVideoName = async (videoId: number, newName: string): Promise<void> => {
  const token = localStorage.getItem('auth_token');
  if (!token) {
    throw new Error('Требуется авторизация. Пожалуйста, войдите в систему.');
  }

  const response = await fetch(`${API_URL}/videos/${videoId}`, {
    method: 'PUT',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ name: newName })
  });

  if (!response.ok) {
    throw new Error('Не удалось изменить название видео');
  }
};

/**
 * Удаляет видео
 */
export const deleteVideo = async (videoId: number): Promise<void> => {
  const token = localStorage.getItem('auth_token');
  if (!token) {
    throw new Error('Требуется авторизация. Пожалуйста, войдите в систему.');
  }

  const response = await fetch(`${API_URL}/videos/${videoId}`, {
    method: 'DELETE',
    headers: {
      'Authorization': `Bearer ${token}`,
    }
  });

  if (!response.ok) {
    throw new Error('Не удалось удалить видео');
  }
};

/**
 * Проверяет файл на соответствие требованиям (формат, размер)
 */
export const validateVideoFile = (file: File): { isValid: boolean; error?: string } => {
  // Проверка типа файла
  if (file.type !== 'video/mp4') {
    return {
      isValid: false,
      error: 'Только файлы в формате MP4.'
    };
  }

  // Проверка размера файла (200MB максимум)
  if (file.size > 200 * 1024 * 1024) {
    return {
      isValid: false,
      error: 'Файл должен быть не более 200 МБ.'
    };
  }

  return { isValid: true };
};
