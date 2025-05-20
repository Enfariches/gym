import { config } from '@vue/test-utils';
import { vi } from 'vitest';

// Mock Quasar components and plugins globally
config.global.mocks = {
  $q: {
    notify: vi.fn(),
    loading: {
      show: vi.fn(),
      hide: vi.fn(),
    },
    screen: {
      gt: {
        sm: true,
        md: true,
      },
      lt: {
        sm: false,
        md: false,
      },
    },
  },
};

// Mock Quasar components
config.global.stubs = {
  'q-page': {
    template: '<div class="q-page"><slot /></div>'
  },
  'q-icon': {
    props: ['name', 'color', 'size'],
    template: '<i :class="name" :data-color="color" :data-size="size" class="q-icon"><slot /></i>'
  },
  'q-spinner': {
    props: ['size', 'color'],
    template: '<div class="q-spinner" :data-size="size" :data-color="color"></div>'
  },
  'q-card': {
    props: ['flat', 'bordered'],
    template: '<div class="q-card"><slot /></div>'
  },
  'q-card-section': {
    template: '<div class="q-card-section"><slot /></div>'
  },
  'q-date': {
    props: ['modelValue', 'range', 'minimal', 'flat'],
    template: '<div class="q-date"></div>'
  }
};

// Mock Chart.js to prevent canvas errors
vi.mock('chart.js', () => {
  return {
    Chart: class Chart {
      static register = vi.fn();
      destroy = vi.fn();
    },
    registerables: [],
  };
});
