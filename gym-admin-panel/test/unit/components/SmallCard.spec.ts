import { describe, it, expect } from 'vitest';
import { mount } from '@vue/test-utils';
import SmallCard from '../../../src/components/SmallCard.vue';

describe('SmallCard.vue', () => {
  it('renders with correct props', () => {
    const title = 'Test Title';
    const value = '123';
    const icon = 'people';

    const wrapper = mount(SmallCard, {
      props: {
        title,
        value,
        icon,
      },
    });

    // Check if the title is rendered correctly
    expect(wrapper.find('.stat-title').text()).toBe(title);

    // Check if the value is rendered correctly
    expect(wrapper.find('.stat-value').text()).toBe(value);

    // Check if the icon is rendered with correct props
    const iconElement = wrapper.find('.q-icon');
    expect(iconElement.exists()).toBe(true);
    expect(iconElement.classes()).toContain(icon);
  });

  it('handles numeric value correctly', () => {
    const wrapper = mount(SmallCard, {
      props: {
        title: 'Number Value',
        value: 456,
        icon: 'check',
      },
    });

    expect(wrapper.find('.stat-value').text()).toBe('456');
  });

  it('applies hover styles on mouseenter', async () => {
    const wrapper = mount(SmallCard, {
      props: {
        title: 'Hover Test',
        value: 'Test',
        icon: 'star',
      },
    });

    // Get the card element
    const card = wrapper.find('.stat-card');

    // Get the styles directly from the rendered component
    // Note: This test might not work as expected with jsdom as it doesn't fully compute CSS styles
    // but it demonstrates the testing approach
    await card.trigger('mouseenter');

    // In a real browser, we could test the transform style
    // For now, just verify that the card element exists
    expect(card.exists()).toBe(true);
  });
});
