<template>
  <q-page class="flex flex-center column">
    <div class="auth-container">
      <div class="text-center q-mb-md">
        <h4 class="q-my-none text-weight-bold">Смена пароля</h4>
      </div>

      <div class="q-pa-md">
        <div v-if="loading && !showForm" class="text-center">
          <q-spinner color="primary" size="3em" class="q-mb-md" />
          <p>Проверка токена...</p>
        </div>

        <div v-if="verified" class="text-center">
          <q-icon name="check_circle" color="positive" size="5em" />
          <h5 class="q-mt-md q-mb-xs">Пароль успешно изменен!</h5>
          <p>Теперь вы можете войти в систему с новым паролем.</p>
          <q-btn
            color="primary"
            to="/auth/login"
            label="Войти в систему"
            class="q-mt-md full-width"
          />
        </div>

        <div v-if="error && !showForm" class="text-center">
          <q-icon name="error" color="negative" size="5em" />
          <h5 class="q-mt-md q-mb-xs">Ошибка</h5>
          <p>{{ errorMessage }}</p>
          <q-btn
            color="primary"
            to="/auth/login"
            label="Вернуться на страницу входа"
            class="q-mt-md full-width"
          />
        </div>

        <div v-if="showForm">
          <q-form @submit="onSubmit" class="q-gutter-md">
            <q-input
              v-model="password"
              filled
              type="password"
              label="Новый пароль"
              lazy-rules
              :rules="[
                (val) => (val && val.length > 5) || 'Пароль должен содержать не менее 6 символов'
              ]"
            />

            <q-input
              v-model="confirmPassword"
              filled
              type="password"
              label="Подтвердите пароль"
              lazy-rules
              :rules="[
                (val) => (val && val === password) || 'Пароли не совпадают'
              ]"
            />

            <div class="q-mt-md">
              <q-btn
                label="Сохранить новый пароль"
                type="submit"
                color="primary"
                class="full-width"
                :loading="submitting"
              />
            </div>
          </q-form>
        </div>
      </div>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { changePassword } from 'src/services/authService';
import { useQuasar } from 'quasar';

const route = useRoute();
const router = useRouter();
const $q = useQuasar();

const resetToken = ref('');
const loading = ref(true);
const verified = ref(false);
const error = ref(false);
const errorMessage = ref('Произошла ошибка при проверке токена сброса пароля');
const showForm = ref(false);
const submitting = ref(false);

const password = ref('');
const confirmPassword = ref('');

onMounted(() => {
  const token = route.query.reset_token as string;

  if (!token) {
    error.value = true;
    errorMessage.value = 'Отсутствует токен сброса пароля';
    loading.value = false;
    return;
  }

  resetToken.value = token;
  loading.value = false;
  showForm.value = true;
});

const onSubmit = async () => {
  if (password.value !== confirmPassword.value) {
    $q.notify({
      type: 'negative',
      message: 'Пароли не совпадают'
    });
    return;
  }

  submitting.value = true;
  try {
    await changePassword(resetToken.value, password.value);
    verified.value = true;
    showForm.value = false;
    $q.notify({
      type: 'positive',
      message: 'Пароль успешно изменен!'
    });

    // Redirect to login page after a short delay
    setTimeout(() => {
      router.push('/auth/login');
    }, 1500);

  } catch (err) {
    console.error('Password change error:', err);
    error.value = true;
    errorMessage.value = 'Произошла ошибка при смене пароля';
    showForm.value = false;
  } finally {
    submitting.value = false;
  }
};
</script>

<style scoped>
.auth-container {
  width: 100%;
  max-width: 450px;
  padding: 2rem;
  background: white;
  border-radius: 8px;
  box-shadow: 0 1px 5px rgb(0 0 0 / 20%);
}
</style>
