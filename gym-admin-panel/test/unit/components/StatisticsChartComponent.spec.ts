/* eslint-disable @typescript-eslint/no-explicit-any */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { mount, flushPromises } from '@vue/test-utils';
import StatisticsChartComponent from '../../../src/components/statPage/StatisticsChartComponent.vue';
import { createTestingPinia } from '@pinia/testing';
import { useStatisticsStore } from '../../../src/stores/statisticsStore';

// Mock the statisticsStore
vi.mock('../../../src/stores/statisticsStore', () => ({
  useStatisticsStore: vi.fn(() => ({
    loading: false,
    statistics: [],
    chartData: [{ date: '19.05', count: 1 }, { date: '20.05', count: 1 }],
    viewMode: 'media',
    setViewMode: vi.fn(),
    fetchDepartmentStatistics: vi.fn().mockResolvedValue([])
  }))
}));

describe('StatisticsChartComponent.vue', () => {
  beforeEach(() => {
    vi.resetAllMocks();
    // Clear any global window event listeners to prevent test pollution
    vi.stubGlobal('addEventListener', vi.fn());
    vi.stubGlobal('removeEventListener', vi.fn());
  });

  it('renders loading state when loading', async () => {
    // Override mock for this test
    vi.mocked(useStatisticsStore).mockReturnValueOnce({
      loading: true,
      statistics: [],
      chartData: [],
      viewMode: 'media',
      setViewMode: vi.fn(),
      fetchDepartmentStatistics: vi.fn().mockResolvedValue([])
    } as any);

    const wrapper = mount(StatisticsChartComponent, {
      props: {
        viewType: 'full',
      },
      global: {
        plugins: [createTestingPinia({ createSpy: vi.fn })],
      },
    });

    expect(wrapper.find('.loading-container').exists()).toBe(true);
    expect(wrapper.find('.loading-text').text()).toBe('Загрузка данных...');
  });

  it('shows no data message when chartData is empty', async () => {
    // Override mock for this test
    vi.mocked(useStatisticsStore).mockReturnValueOnce({
      loading: false,
      statistics: [],
      chartData: [],
      viewMode: 'media',
      setViewMode: vi.fn(),
      fetchDepartmentStatistics: vi.fn().mockResolvedValue([])
    } as any);

    const wrapper = mount(StatisticsChartComponent, {
      props: {
        viewType: 'full',
      },
      global: {
        plugins: [createTestingPinia({ createSpy: vi.fn })],
      },
    });

    expect(wrapper.find('.no-data-container').exists()).toBe(true);
    expect(wrapper.find('.no-data-text').text()).toBe('Нет данных для отображения');
  });

  it('renders the chart when data is available', async () => {
    const wrapper = mount(StatisticsChartComponent, {
      props: {
        viewType: 'full',
      },
      global: {
        plugins: [createTestingPinia({ createSpy: vi.fn })],
      },
    });

    await flushPromises();

    // Chart canvas should be rendered
    expect(wrapper.find('.canvas-container').exists()).toBe(true);
    expect(wrapper.find('canvas').exists()).toBe(true);
  });

  it('changes viewMode when viewType prop changes', async () => {
    const setViewModeSpy = vi.fn();

    // Override mock for this test
    vi.mocked(useStatisticsStore).mockReturnValueOnce({
      loading: false,
      statistics: [],
      chartData: [{ date: '19.05', count: 1 }, { date: '20.05', count: 1 }],
      viewMode: 'media',
      setViewMode: setViewModeSpy,
      fetchDepartmentStatistics: vi.fn().mockResolvedValue([])
    } as any);

    const wrapper = mount(StatisticsChartComponent, {
      props: {
        viewType: 'full',
      },
      global: {
        plugins: [createTestingPinia({ createSpy: vi.fn })],
      },
    });

    await flushPromises();

    // Change the viewType prop
    await wrapper.setProps({ viewType: 'half' });

    // The setViewMode should have been called with 'department' (based on the mapping)
    expect(setViewModeSpy).toHaveBeenCalledWith('department');
  });

  it('adds resize event listener on mount and removes it on unmount', async () => {
    const addEventListenerSpy = vi.fn();
    const removeEventListenerSpy = vi.fn();

    window.addEventListener = addEventListenerSpy;
    window.removeEventListener = removeEventListenerSpy;

    const wrapper = mount(StatisticsChartComponent, {
      props: {
        viewType: 'full',
      },
      global: {
        plugins: [createTestingPinia({ createSpy: vi.fn })],
      },
    });

    await flushPromises();

    // Check that addEventListener was called with 'resize'
    expect(addEventListenerSpy).toHaveBeenCalledWith('resize', expect.any(Function));

    // Unmount the component
    wrapper.unmount();

    // Check that removeEventListener was called with 'resize'
    expect(removeEventListenerSpy).toHaveBeenCalledWith('resize', expect.any(Function));
  });
});
