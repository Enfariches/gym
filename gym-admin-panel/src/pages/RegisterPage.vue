<template>
  <q-page class="bg-grey-1 register-page">
    <div class="content">
      <div class="register-container">
        <q-card class="register-card">
          <q-card-section class="q-pb-none">
            <h1 class="page-title text-center">Регистрация</h1>
            <p class="text-subtitle text-center">Создайте новый аккаунт</p>
          </q-card-section>

          <q-card-section>
            <q-form @submit="handleRegister" class="q-gutter-md">
              <q-input
                v-model="email"
                label="Email"
                type="email"
                outlined
                class="register-input"
                :rules="[val => !!val || 'Email обязателен', val => /.+@.+\..+/.test(val) || 'Email должен быть действительным']"
                required
              />
              <q-input
                v-model="password"
                label="Пароль"
                type="password"
                outlined
                class="register-input"
                :rules="[val => !!val || 'Пароль обязателен', val => val.length >= 8 || 'Пароль должен содержать минимум 8 символов']"
                required
              />
              <q-input
                v-model="confirmPassword"
                label="Подтвердите пароль"
                type="password"
                outlined
                class="register-input"
                :rules="[val => !!val || 'Подтвердите пароль', val => val === password || 'Пароли не совпадают']"
                required
              />

              <div class="row q-mt-lg">
                <q-btn
                  type="submit"
                  color="primary"
                  label="Зарегистрироваться"
                  :loading="loading"
                  class="full-width"
                  size="lg"
                />
              </div>

              <div class="row q-mt-md justify-center">
                <p class="q-mr-sm no-margin">Уже есть аккаунт?</p>
                <q-btn
                  flat
                  color="primary"
                  label="Войти"
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

    <q-dialog v-model="showVerificationDialog">
      <q-card style="min-width: 400px">
        <q-card-section>
          <h4 class="text-h5 q-mb-md">Подтверждение регистрации</h4>
          <p class="q-mb-md">Пожалуйста, проверьте ваш email для подтверждения регистрации. Если вы не получили письмо, введите токен подтверждения ниже:</p>
          <q-form @submit="handleVerification" class="q-gutter-md">
            <q-input
              v-model="verificationToken"
              label="Токен подтверждения"
              outlined
              :rules="[val => !!val || 'Токен обязателен']"
              required
            />
            <div class="row q-mt-lg">
              <q-btn
                type="submit"
                color="primary"
                label="Подтвердить"
                :loading="verifying"
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
const confirmPassword = ref('');
const loading = ref(false);
const showVerificationDialog = ref(false);
const verificationToken = ref('');
const verifying = ref(false);

const handleRegister = async () => {
  try {
    loading.value = true;
    const response = await authStore.register(email.value, password.value);
    console.log(response);
    showVerificationDialog.value = true;
    $q.notify({
      color: 'positive',
      message: 'Регистрация успешна! Пожалуйста, проверьте email для подтверждения аккаунта.',
      icon: 'check_circle'
    });
  } catch (error) {
    console.error(error);
    $q.notify({
      color: 'negative',
      message: 'Ошибка регистрации. Пожалуйста, попробуйте снова.',
      icon: 'report_problem'
    });
  } finally {
    loading.value = false;
  }
};

const handleVerification = async () => {
  try {
    verifying.value = true;
    await authStore.verifyRegister(verificationToken.value);
    $q.notify({
      color: 'positive',
      message: 'Аккаунт успешно подтвержден! Теперь вы можете войти в систему.',
      icon: 'check_circle'
    });
    router.push('/login');
  } catch (error) {
    console.error(error);
    $q.notify({
      color: 'negative',
      message: 'Ошибка подтверждения. Пожалуйста, попробуйте снова.',
      icon: 'report_problem'
    });
  } finally {
    verifying.value = false;
  }
};
</script>

<style scoped>
.register-page {
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

.register-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.register-card {
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

.register-input {
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
