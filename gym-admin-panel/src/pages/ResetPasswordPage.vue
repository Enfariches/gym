<template>
  <div class="reset-password-page">
    <q-card class="reset-password-card">
      <q-card-section>
        <h4 class="text-h4 q-mb-md">Reset Password</h4>
        <q-form @submit="handleResetPassword" class="q-gutter-md">
          <q-input
            v-model="email"
            label="Email"
            type="email"
            :rules="[val => !!val || 'Email is required', val => /.+@.+\..+/.test(val) || 'Email must be valid']"
            required
          />
          <div class="row q-mt-lg">
            <q-btn
              type="submit"
              color="primary"
              label="Send Reset Instructions"
              :loading="loading"
              class="full-width"
            />
          </div>
          <div class="row q-mt-sm">
            <q-btn
              flat
              color="primary"
              label="Back to Login"
              to="/login"
            />
          </div>
        </q-form>
      </q-card-section>
    </q-card>
  </div>
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
      message: 'Password reset instructions have been sent to your email.',
      icon: 'check_circle'
    });
    router.push('/login');
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  } catch (error) {
    $q.notify({
      color: 'negative',
      message: 'Failed to send reset instructions. Please try again.',
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
  background-color: #f5f5f5;
}

.reset-password-card {
  width: 100%;
  max-width: 400px;
  padding: 20px;
}
</style>
