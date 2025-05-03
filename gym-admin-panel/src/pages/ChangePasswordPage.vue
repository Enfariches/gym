<template>
  <div class="change-password-page">
    <q-card class="change-password-card">
      <q-card-section>
        <h4 class="text-h4 q-mb-md">Change Password</h4>
        <q-form @submit="handleChangePassword" class="q-gutter-md">
          <q-input
            v-model="newPassword"
            label="New Password"
            type="password"
            :rules="[val => !!val || 'Password is required', val => val.length >= 8 || 'Password must be at least 8 characters']"
            required
          />
          <q-input
            v-model="confirmPassword"
            label="Confirm New Password"
            type="password"
            :rules="[val => !!val || 'Please confirm your password', val => val === newPassword || 'Passwords do not match']"
            required
          />
          <div class="row q-mt-lg">
            <q-btn
              type="submit"
              color="primary"
              label="Change Password"
              :loading="loading"
              class="full-width"
            />
          </div>
        </q-form>
      </q-card-section>
    </q-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import { useQuasar } from 'quasar';

const $q = useQuasar();
const router = useRouter();
const route = useRoute();
const authStore = useAuthStore();

const newPassword = ref('');
const confirmPassword = ref('');
const loading = ref(false);
const resetToken = ref('');

onMounted(() => {
  resetToken.value = route.query.token as string;
  if (!resetToken.value) {
    $q.notify({
      color: 'negative',
      message: 'Invalid reset token. Please request a new password reset.',
      icon: 'report_problem'
    });
    router.push('/reset-password');
  }
});

const handleChangePassword = async () => {
  try {
    loading.value = true;
    await authStore.changePassword(resetToken.value, newPassword.value);
    $q.notify({
      color: 'positive',
      message: 'Password changed successfully!',
      icon: 'check_circle'
    });
    router.push('/login');
  } catch (error) {
    console.error(error);
    $q.notify({
      color: 'negative',
      message: 'Failed to change password. Please try again.',
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
  background-color: #f5f5f5;
}

.change-password-card {
  width: 100%;
  max-width: 400px;
  padding: 20px;
}
</style>
