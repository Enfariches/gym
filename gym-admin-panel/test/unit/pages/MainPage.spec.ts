/* eslint-disable @typescript-eslint/no-explicit-any */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { mount, flushPromises } from '@vue/test-utils';
import MainPage from '../../../src/pages/MainPage.vue';
import { createTestingPinia } from '@pinia/testing';
import { useStatisticsStore } from '../../../src/stores/statisticsStore';
import { useMediaStore } from '../../../src/stores/mediaStore';

// Mock the stores
vi.mock('../../../src/stores/statisticsStore', () => ({
  useStatisticsStore: vi.fn(() => ({
    statisticsSummary: {
      uniqueUsers: 42,
      completedPercentage: 75,
      totalViews: 100,
      averageViewTime: '10:30'
    },
    filteredStatistics: [
      { createdAt: { seconds: Math.floor(Date.now() / 1000) } }
    ],
    fetchDepartmentStatistics: vi.fn().mockResolvedValue([])
  }))
}));

vi.mock('../../../src/stores/mediaStore', () => ({
  useMediaStore: vi.fn(() => ({
    videos: [{ id: 1, title: 'Test Video' }],
    loadVideos: vi.fn().mockResolvedValue([])
  }))
}));

// Mock the components used in MainPage
vi.mock('../../../src/components/SmallCard.vue', () => ({
  default: {
    name: 'SmallCard',
    props: ['title', 'value', 'icon'],
    template: '<div class="mock-small-card">{{ title }}: {{ value }}</div>'
  }
}));

vi.mock('../../../src/components/mainPage/BigCard.vue', () => ({
  default: {
    name: 'BigCard',
    props: ['title', 'icon'],
    template: '<div class="mock-big-card"><div>{{ title }}</div><slot /></div>'
  }
}));

vi.mock('../../../src/components/statPage/StatisticsChartComponent.vue', () => ({
  default: {
    name: 'StatisticsChartComponent',
    props: ['viewType'],
    template: '<div class="mock-chart">Chart: {{ viewType }}</div>'
  }
}));

describe('MainPage.vue', () => {
  beforeEach(() => {
    vi.resetAllMocks();
  });

  it('renders the page title correctly', async () => {
    const wrapper = mount(MainPage, {
      global: {
        plugins: [createTestingPinia({ createSpy: vi.fn })],
      }
    });

    expect(wrapper.find('.page-title').text()).toBe('Обзор системы');
  });

  it('displays statistics data', async () => {
    const wrapper = mount(MainPage, {
      global: {
        plugins: [createTestingPinia({ createSpy: vi.fn })],
      }
    });

    await flushPromises();

    // Check for page title
    expect(wrapper.find('.page-title').text()).toBe('Обзор системы');

    // Check for chart container
    expect(wrapper.find('.chart-container').exists()).toBe(true);
  });

  it('loads data on mount', async () => {
    const loadVideosSpy = vi.fn().mockResolvedValue([]);
    const fetchStatsSpy = vi.fn().mockResolvedValue([]);

    vi.mocked(useMediaStore).mockReturnValue({
      videos: [],
      loadVideos: loadVideosSpy
    } as any);

    vi.mocked(useStatisticsStore).mockReturnValue({
      statisticsSummary: {
        uniqueUsers: 0,
        completedPercentage: 0,
        totalViews: 0,
        averageViewTime: '0:00'
      },
      filteredStatistics: [],
      fetchDepartmentStatistics: fetchStatsSpy
    } as any);

    mount(MainPage, {
      global: {
        plugins: [createTestingPinia({ createSpy: vi.fn })],
      }
    });

    await flushPromises();

    expect(loadVideosSpy).toHaveBeenCalled();
    expect(fetchStatsSpy).toHaveBeenCalled();
  });

  it('shows loading spinner initially', async () => {
    // Create a promise that doesn't resolve immediately
    const loadingPromise = new Promise(resolve => {
      setTimeout(() => resolve([]), 100);
    });

    // Mock the statistics loading state
    vi.mocked(useStatisticsStore).mockReturnValue({
      statisticsSummary: {
        uniqueUsers: 0,
        completedPercentage: 0,
        totalViews: 0,
        averageViewTime: '0:00'
      },
      filteredStatistics: [],
      fetchDepartmentStatistics: vi.fn().mockReturnValue(loadingPromise)
    } as any);

    const wrapper = mount(MainPage, {
      global: {
        plugins: [createTestingPinia({ createSpy: vi.fn })],
      }
    });

    // Initially there should be a spinner
    expect(wrapper.find('.q-spinner').exists()).toBe(true);

    // Wait for data loading to complete
    await flushPromises();
    await new Promise(resolve => setTimeout(resolve, 150));
  });
});
