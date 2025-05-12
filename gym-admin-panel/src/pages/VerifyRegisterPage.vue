<template>
  <q-page class="flex flex-center column">
    <div class="auth-container">
      <div class="text-center q-mb-md">
        <h4 class="q-my-none text-weight-bold">Подтверждение регистрации</h4>
      </div>

      <div class="q-pa-md">
        <div v-if="loading" class="text-center">
          <q-spinner color="primary" size="3em" class="q-mb-md" />
          <p>Проверка токена...</p>
        </div>

        <div v-if="verified" class="text-center">
          <q-icon name="check_circle" color="positive" size="5em" />
          <h5 class="q-mt-md q-mb-xs">Регистрация подтверждена!</h5>
          <p>Ваш email успешно подтвержден.</p>
          <q-btn
            color="primary"
            to="/auth/login"
            label="Войти в систему"
            class="q-mt-md full-width"
          />
        </div>

        <div v-if="error" class="text-center">
          <q-icon name="error" color="negative" size="5em" />
          <h5 class="q-mt-md q-mb-xs">Ошибка подтверждения</h5>
          <p>{{ errorMessage }}</p>
          <q-btn
            color="primary"
            to="/auth/login"
            label="Вернуться на страницу входа"
            class="q-mt-md full-width"
          />
        </div>
      </div>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { verifyRegister } from 'src/services/authService';
import { useQuasar } from 'quasar';

const route = useRoute();
const router = useRouter();
const $q = useQuasar();

const loading = ref(true);
const verified = ref(false);
const error = ref(false);
const errorMessage = ref('Произошла ошибка при подтверждении регистрации');

onMounted(async () => {
  const authToken = route.query.auth_token as string;

  if (!authToken) {
    error.value = true;
    errorMessage.value = 'Отсутствует токен подтверждения';
    loading.value = false;
    return;
  }

  try {
    await verifyRegister(authToken);
    verified.value = true;
    $q.notify({
      type: 'positive',
      message: 'Email успешно подтвержден!'
    });

    // Redirect to login page after a short delay
    setTimeout(() => {
      router.push('/auth/login');
    }, 1500);

  } catch (err) {
    console.error('Verification error:', err);
    error.value = true;
    errorMessage.value = 'Произошла ошибка при подтверждении регистрации';
  } finally {
    loading.value = false;
  }
});
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
