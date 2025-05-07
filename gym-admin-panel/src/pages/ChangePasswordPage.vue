<template>
  <q-page class="bg-grey-1 change-password-page">
    <div class="content">
      <div class="change-password-container">
        <q-card class="change-password-card">
          <q-card-section class="q-pb-none">
            <h1 class="page-title text-center">Изменение пароля</h1>
            <p class="text-subtitle text-center">Введите новый пароль для вашего аккаунта</p>
          </q-card-section>

          <q-card-section>
            <q-form @submit="handleChangePassword" class="q-gutter-md">
              <q-input
                v-model="currentPassword"
                label="Текущий пароль"
                type="password"
                outlined
                class="password-input"
                :rules="[val => !!val || 'Текущий пароль обязателен']"
                required
              />
              <q-input
                v-model="newPassword"
                label="Новый пароль"
                type="password"
                outlined
                class="password-input"
                :rules="[val => !!val || 'Новый пароль обязателен', val => val.length >= 8 || 'Пароль должен содержать минимум 8 символов']"
                required
              />
              <q-input
                v-model="confirmPassword"
                label="Подтвердите новый пароль"
                type="password"
                outlined
                class="password-input"
                :rules="[val => !!val || 'Подтвердите новый пароль', val => val === newPassword || 'Пароли не совпадают']"
                required
              />

              <div class="row q-mt-lg">
                <q-btn
                  type="submit"
                  color="primary"
                  label="Изменить пароль"
                  :loading="loading"
                  class="full-width"
                  size="lg"
                />
              </div>

              <div class="row q-mt-md justify-center">
                <q-btn
                  flat
                  color="primary"
                  label="Вернуться на главную страницу"
                  to="/"
                  padding="none"
                  class="no-padding"
                />
              </div>
            </q-form>
          </q-card-section>
        </q-card>

        <div class="copyright-text">
          © {{ new Date().getFullYear() }} Производственная Гимнастика. Все права защищены.
        </div>
      </div>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import { useQuasar } from 'quasar';

const $q = useQuasar();
const router = useRouter();
const authStore = useAuthStore();

const currentPassword = ref('');
const newPassword = ref('');
const confirmPassword = ref('');
const loading = ref(false);

const handleChangePassword = async () => {
  try {
    loading.value = true;
    await authStore.changePassword(currentPassword.value, newPassword.value);
    $q.notify({
      color: 'positive',
      message: 'Пароль успешно изменен',
      icon: 'check_circle'
    });
    router.push('/');
  } catch (error) {
    $q.notify({
      color: 'negative',
      message: error instanceof Error ? error.message : 'Не удалось изменить пароль. Пожалуйста, попробуйте снова.',
      icon: 'report_problem'
    });
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.change-password-page {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
}

.content {
  width: 100%;
  max-width: 1200px;
  padding: 40px;
}

.change-password-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.change-password-card {
  width: 100%;
  max-width: 450px;
  padding: 20px;
}

.page-title {
  color: rgba(90,92,105,1);
  font-weight: bold;
  font-size: 32px;
  margin-bottom: 10px;
}

.text-subtitle {
  color: rgba(108,117,125,1);
  font-size: 16px;
  margin-bottom: 20px;
}

.password-input {
  margin-bottom: 20px;
}

.no-padding {
  padding: 0;
}

.copyright-text {
  margin-top: 30px;
  color: rgba(108,117,125,1);
  font-size: 14px;
  text-align: center;
}
</style>
