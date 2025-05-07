<template>
  <q-page class="bg-grey-1 login-page">
    <div class="content">
      <div class="login-container">
        <q-card class="login-card">
          <q-card-section class="q-pb-none">
            <h1 class="page-title text-center">Вход в систему</h1>
            <p class="text-subtitle text-center">Войдите в свой аккаунт</p>
          </q-card-section>

          <q-card-section>
            <q-form @submit="handleLogin" class="q-gutter-md">
              <q-input
                v-model="email"
                label="Email"
                type="email"
                outlined
                class="login-input"
                :rules="[val => !!val || 'Email обязателен', val => /.+@.+\..+/.test(val) || 'Email должен быть действительным']"
                required
              />
              <q-input
                v-model="password"
                label="Пароль"
                type="password"
                outlined
                class="login-input"
                :rules="[val => !!val || 'Пароль обязателен']"
                required
              />

              <div class="row justify-between items-center q-mb-md">
                <q-checkbox v-model="rememberMe" label="Запомнить меня" color="primary" />
                <q-btn
                  flat
                  color="primary"
                  label="Забыли пароль?"
                  @click="showResetPassword = true"
                  class="no-padding"
                />
              </div>

              <div class="row q-mt-lg">
                <q-btn
                  type="submit"
                  color="primary"
                  label="Войти"
                  :loading="loading"
                  class="full-width"
                  size="lg"
                />
              </div>

              <div class="row q-mt-md justify-center">
                <p class="q-mr-sm no-margin">Нет аккаунта?</p>
                <q-btn
                  flat
                  color="primary"
                  label="Зарегистрироваться"
                  to="/register"
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

    <q-dialog v-model="showResetPassword">
      <q-card style="min-width: 400px">
        <q-card-section>
          <h4 class="text-h5 q-mb-md">Сброс пароля</h4>
          <q-form @submit="handleResetPassword" class="q-gutter-md">
            <q-input
              v-model="resetEmail"
              label="Email"
              type="email"
              outlined
              :rules="[val => !!val || 'Email обязателен', val => /.+@.+\..+/.test(val) || 'Email должен быть действительным']"
              required
            />
            <div class="row q-mt-lg">
              <q-btn
                type="submit"
                color="primary"
                label="Отправить инструкции по сбросу"
                :loading="resetting"
                class="full-width"
              />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>
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
const password = ref('');
const resetEmail = ref('');
const showResetPassword = ref(false);
const loading = ref(false);
const resetting = ref(false);
const rememberMe = ref(false);

const handleLogin = async () => {
  try {
    loading.value = true;
    await authStore.login(email.value, password.value);
    $q.notify({
      color: 'positive',
      message: 'Вход выполнен успешно!',
      icon: 'check_circle'
    });
    router.push('/');
  } catch (error) {
    $q.notify({
      color: 'negative',
      message: error instanceof Error ? error.message : 'Ошибка входа. Проверьте свои учетные данные.',
      icon: 'report_problem'
    });
  } finally {
    loading.value = false;
  }
};

const handleResetPassword = async () => {
  try {
    resetting.value = true;
    await authStore.resetPassword(resetEmail.value);
    $q.notify({
      color: 'positive',
      message: 'Инструкции по сбросу пароля отправлены на вашу электронную почту.',
      icon: 'check_circle'
    });
    showResetPassword.value = false;
  } catch (error) {
    $q.notify({
      color: 'negative',
      message: error instanceof Error ? error.message : 'Не удалось сбросить пароль. Пожалуйста, попробуйте снова.',
      icon: 'report_problem'
    });
  } finally {
    resetting.value = false;
  }
};
</script>

<style scoped>
.login-page {
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

.login-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.login-card {
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

.login-input {
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
