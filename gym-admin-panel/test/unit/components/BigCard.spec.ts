import { describe, it, expect } from 'vitest';
import { mount } from '@vue/test-utils';
import BigCard from '../../../src/components/mainPage/BigCard.vue';

describe('BigCard.vue', () => {
  it('renders with correct props', () => {
    const title = 'Test Title';
    const icon = 'dashboard';

    const wrapper = mount(BigCard, {
      props: {
        title,
        icon,
      },
    });

    // Check if the title is rendered correctly
    expect(wrapper.find('.card-title').text()).toContain(title);

    // Check if the icon is rendered with correct props
    const iconElement = wrapper.find('.q-icon');
    expect(iconElement.exists()).toBe(true);
    expect(iconElement.classes()).toContain(icon);
  });

  it('renders slot content', () => {
    const slotContent = '<div class="test-slot">Slot Content</div>';

    const wrapper = mount(BigCard, {
      props: {
        title: 'Card with Slot',
        icon: 'info',
      },
      slots: {
        default: slotContent,
      },
    });

    // Check if the slot content is rendered
    expect(wrapper.html()).toContain('test-slot');
    expect(wrapper.html()).toContain('Slot Content');
  });
});
