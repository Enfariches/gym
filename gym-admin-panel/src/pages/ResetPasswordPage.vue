<template>
  <q-page class="bg-grey-1 reset-password-page">
    <div class="content">
      <div class="reset-password-container">
        <q-card class="reset-password-card">
          <q-card-section class="q-pb-none">
            <h1 class="page-title text-center">Сброс пароля</h1>
            <p class="text-subtitle text-center">Введите ваш email для получения инструкций по сбросу пароля</p>
          </q-card-section>

          <q-card-section>
            <q-form @submit="handleResetPassword" class="q-gutter-md">
              <q-input
                v-model="email"
                label="Email"
                type="email"
                outlined
                class="reset-input"
                :rules="[val => !!val || 'Email обязателен', val => /.+@.+\..+/.test(val) || 'Email должен быть действительным']"
                required
              />

              <div class="row q-mt-lg">
                <q-btn
                  type="submit"
                  color="primary"
                  label="Отправить инструкции"
                  :loading="loading"
                  class="full-width"
                  size="lg"
                />
              </div>

              <div class="row q-mt-md justify-center">
                <q-btn
                  flat
                  color="primary"
                  label="Вернуться на страницу входа"
                  to="/login"
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

const email = ref('');
const loading = ref(false);

const handleResetPassword = async () => {
  try {
    loading.value = true;
    await authStore.resetPassword(email.value);
    $q.notify({
      color: 'positive',
      message: 'Инструкции по сбросу пароля отправлены на вашу электронную почту',
      icon: 'check_circle'
    });
    router.push('/login');
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  } catch (error) {
    $q.notify({
      color: 'negative',
      message: 'Не удалось отправить инструкции. Пожалуйста, проверьте email и попробуйте снова',
      icon: 'report_problem'
    });
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.reset-password-page {
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

.reset-password-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.reset-password-card {
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

.reset-input {
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
