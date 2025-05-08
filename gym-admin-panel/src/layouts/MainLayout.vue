<template>
  <q-layout view="lHh Lpr lFf">
    <q-header elevated class="header-gradient">
      <q-toolbar>
        <q-btn
          flat
          dense
          round
          icon="menu"
          aria-label="Menu"
          @click="toggleLeftDrawer"
          class="q-mr-sm"
        />

        <q-toolbar-title class="text-weight-bold">
          Производственная Гимнастика
        </q-toolbar-title>
      </q-toolbar>
    </q-header>

    <q-drawer
      v-model="leftDrawerOpen"
      show-if-above
      :width="280"
      :breakpoint="400"
      class="drawer-background no-border-right"
    >
      <div class="column fit">
        <q-list padding class="col scroll">
          <EssentialLink
            v-for="link in linksList"
            :key="link.title"
            v-bind="link"
            :active="isLinkActive(link.link)"
          />
        </q-list>

        <q-item
          clickable
          v-ripple
          class="logout-item q-mt-md"
          @click="logout"
        >
          <q-item-section avatar>
            <q-icon name="logout" />
          </q-item-section>

          <q-item-section>
            <q-item-label>Выйти из системы</q-item-label>
          </q-item-section>
        </q-item>
      </div>
    </q-drawer>

    <q-page-container class="content-background">
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import EssentialLink, { type EssentialLinkProps } from 'components/EssentialLink.vue';
import { useRouter, useRoute } from 'vue-router';

const router = useRouter();
const route = useRoute();

const linksList: EssentialLinkProps[] = [
  {
    title: 'Панель',
    caption: 'Главная страница',
    icon: 'dashboard',
    link: '/admin',
  },
  {
    title: 'Видео',
    caption: 'Добавьте видео для сотрудников',
    icon: 'video_library',
    link: '/admin/videos/',
  },
  {
    title: 'Расписание',
    caption: 'Создайте / отредактируйте расписание',
    icon: 'schedule',
    link: '/admin/schedule/',
  },
  {
    title: 'Статистика',
    caption: 'Статистика по просмотренным видео',
    icon: 'trending_up',
    link: '/admin/stats/',
  },
  {
    title: 'Профиль',
    caption: 'Управление личным профилем',
    icon: 'person',
    link: '/admin/profile/',
  },
];

const leftDrawerOpen = ref(false);

// Загружаем состояние drawer'а при монтировании
onMounted(() => {
  const savedState = localStorage.getItem('drawerOpen');
  if (savedState !== null) {
    leftDrawerOpen.value = JSON.parse(savedState);
  }
});

// Сохраняем состояние drawer'а при изменении
watch(leftDrawerOpen, (newValue) => {
  localStorage.setItem('drawerOpen', JSON.stringify(newValue));
});

function toggleLeftDrawer() {
  leftDrawerOpen.value = !leftDrawerOpen.value;
}

function isLinkActive(link: string): boolean {
  if (link === '/admin') {
    // Для главной панели активна только если путь точно /admin
    return route.path === '/admin';
  }
  // Для остальных страниц проверяем, начинается ли путь с link
  return route.path.startsWith(link);
}

function logout(): void {
  localStorage.removeItem('drawerOpen'); // Очищаем состояние drawer'а при выходе
  localStorage.clear();
  router.push('/login');
}
</script>

<style lang="scss">
.header-gradient {
  background: linear-gradient(135deg, #4e73df 0%, #224abe 100%);
}

.drawer-background {
  background: #1e1e2d;
  color: white;

  .scroll {
    height: 100%;
    overflow-y: auto;
  }

  .q-item {
    border-radius: 8px;
    margin: 4px 8px;
    color: #ffffff;
    transition: all 0.3s ease;
    min-height: 56px;
    padding: 8px 16px;

    &:hover {
      background: rgba(78,115,223,0.1);
      color: #4e73df;
      transform: translateX(4px);

      .q-icon {
        color: #4e73df;
      }
    }

    &.q-router-link--active {
      background: rgba(78,115,223,0.2);
      color: #ffffff;
      font-weight: 600;
      transform: translateX(4px);

      .q-icon {
        color: #ffffff;
      }

      .q-item__label--caption {
        color: rgba(255, 255, 255, 0.9);
      }
    }

    .q-item__section--avatar {
      min-width: 40px;
      padding-right: 12px;

      .q-icon {
        font-size: 24px;
        color: #ffffff;
      }
    }

    .q-item__label {
      font-size: 0.9rem;
      font-weight: 500;
      letter-spacing: 0.01em;
    }

    .q-item__label--caption {
      font-size: 0.75rem;
      color: rgba(255, 255, 255, 0.7);
      margin-top: 4px;
    }
  }
}

.no-border-right {
  .q-drawer {
    border-right: 0 !important;
  }
}

.logout-item {
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  color: #e74a3b !important;

  .q-icon {
    color: #e74a3b !important;
  }

  &:hover {
    background: rgba(231, 74, 59, 0.1) !important;
    color: #e74a3b !important;
  }
}

.content-background {
  background: #f8f9fc;
}

// Smooth transitions
.q-drawer--left {
  transition: transform 0.3s ease;
}

.q-layout-drawer {
  &--shown {
    transform: translateX(0);
  }

  &--hidden {
    transform: translateX(-100%);
  }
}
</style>
