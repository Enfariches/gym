import { boot } from 'quasar/wrappers';
import { Screen } from 'quasar';

// Определяем типы для размеров экрана
export interface ScreenSizes {
  xs: boolean;
  sm: boolean;
  md: boolean;
  lg: boolean;
  xl: boolean;
}

export default boot(({ app }) => {
  // Добавляем глобальные методы для проверки размера экрана
  app.config.globalProperties.$isScreenSize = (size: keyof ScreenSizes): boolean => {
    return Screen[size];
  };
}); 