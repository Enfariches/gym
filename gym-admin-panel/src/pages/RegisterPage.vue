<template>
  <div class="register-page">
    <q-card class="register-card">
      <q-card-section>
        <h4 class="text-h4 q-mb-md">Register</h4>
        <q-form @submit="handleRegister" class="q-gutter-md">
          <q-input
            v-model="email"
            label="Email"
            type="email"
            :rules="[val => !!val || 'Email is required', val => /.+@.+\..+/.test(val) || 'Email must be valid']"
            required
          />
          <q-input
            v-model="password"
            label="Password"
            type="password"
            :rules="[val => !!val || 'Password is required', val => val.length >= 8 || 'Password must be at least 8 characters']"
            required
          />
          <q-input
            v-model="confirmPassword"
            label="Confirm Password"
            type="password"
            :rules="[val => !!val || 'Please confirm your password', val => val === password || 'Passwords do not match']"
            required
          />
          <div class="row q-mt-lg">
            <q-btn
              type="submit"
              color="primary"
              label="Register"
              :loading="loading"
              class="full-width"
            />
          </div>
          <div class="row q-mt-sm">
            <q-btn
              flat
              color="primary"
              label="Already have an account? Login"
              to="/login"
            />
          </div>
        </q-form>
      </q-card-section>
    </q-card>

    <q-dialog v-model="showVerificationDialog">
      <q-card>
        <q-card-section>
          <h4 class="text-h4 q-mb-md">Verify Registration</h4>
          <p class="q-mb-md">Please check your email for the verification link. If you haven't received it, you can enter the verification token below:</p>
          <q-form @submit="handleVerification" class="q-gutter-md">
            <q-input
              v-model="verificationToken"
              label="Verification Token"
              :rules="[val => !!val || 'Token is required']"
              required
            />
            <div class="row q-mt-lg">
              <q-btn
                type="submit"
                color="primary"
                label="Verify"
                :loading="verifying"
                class="full-width"
              />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>
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
      message: 'Registration successful! Please check your email to verify your account.',
      icon: 'check_circle'
    });
  } catch (error) {
    console.error(error);
    $q.notify({
      color: 'negative',
      message: 'Registration failed. Please try again.',
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
      message: 'Account verified successfully! You can now login.',
      icon: 'check_circle'
    });
    router.push('/login');
  } catch (error) {
    console.error(error);
    $q.notify({
      color: 'negative',
      message: 'Verification failed. Please try again.',
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
  background-color: #f5f5f5;
}

.register-card {
  width: 100%;
  max-width: 400px;
  padding: 20px;
}
</style>
