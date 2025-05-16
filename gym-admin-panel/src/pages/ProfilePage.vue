<template>
  <q-page class="bg-grey-1">
    <div class="content">
      <h1 class="page-title">Профиль администратора</h1>

      <div class="row q-col-gutter-md">
        <div class="col-12 col-md-6">
          <ProfileCard
            :admin="adminStore.admin"
            :edit-mode="editMode"
            @edit="startEditMode"
            @save="saveProfile"
            @cancel="cancelEdit"
            @update:formData="updateFormData"
          />
        </div>

        <div class="col-12 col-md-6">
          <q-card class="password-card">
            <q-card-section class="bg-primary text-white">
              <div class="text-h6">Изменение пароля</div>
            </q-card-section>

            <q-card-section>
              <q-form @submit="changePassword" class="q-pa-md">
                <q-input
                  v-model="passwordData.currentPassword"
                  label="Текущий пароль"
                  type="password"
                  outlined
                  dense
                  class="q-mb-md"
                  :rules="[val => !!val || 'Пожалуйста, введите текущий пароль']"
                />
                <q-input
                  v-model="passwordData.newPassword"
                  label="Новый пароль"
                  type="password"
                  outlined
                  dense
                  class="q-mb-md"
                  :rules="[
                    val => !!val || 'Пожалуйста, введите новый пароль',
                    val => val.length >= 6 || 'Пароль должен содержать минимум 6 символов'
                  ]"
                />
                <q-input
                  v-model="passwordData.confirmPassword"
                  label="Подтвердите новый пароль"
                  type="password"
                  outlined
                  dense
                  class="q-mb-md"
                  :rules="[
                    val => !!val || 'Пожалуйста, подтвердите новый пароль',
                    val => val === passwordData.newPassword || 'Пароли не совпадают'
                  ]"
                />

                <div class="row justify-end">
                  <q-btn
                    type="submit"
                    color="primary"
                    label="Сменить пароль"
                    :loading="passwordLoading"
                  />
                </div>
              </q-form>
            </q-card-section>
          </q-card>
        </div>
      </div>
    </div>

    <!-- Индикатор загрузки -->
    <q-inner-loading :showing="adminStore.isLoading">
      <q-spinner size="50px" color="primary" />
    </q-inner-loading>
  </q-page>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useQuasar } from 'quasar';
import ProfileCard from 'src/components/ProfileCard.vue';
import { useAdminStore } from '../stores/adminStore';
import type { Admin } from '../../protogen/v1/users/admin';

const $q = useQuasar();
const adminStore = useAdminStore();

const editMode = ref(false);
const formData = ref<Partial<Admin>>({});
const passwordData = ref({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
});
const passwordLoading = ref(false);

onMounted(async () => {
  try {
    await adminStore.fetchAdmin();
  } catch {
    $q.notify({
      color: 'negative',
      message: 'Не удалось загрузить данные профиля.',
      icon: 'error'
    });
  }
});

const startEditMode = () => {
  editMode.value = true;
};

const cancelEdit = () => {
  editMode.value = false;
  formData.value = {};
};

const updateFormData = (newData: Partial<Admin>) => {
  formData.value = newData;
};

const saveProfile = async () => {
  try {
    // Определяем, какие поля были изменены
    const fieldsToUpdate: string[] = [];
    if (formData.value.name) fieldsToUpdate.push('name');
    if (formData.value.surname) fieldsToUpdate.push('surname');
    if (formData.value.departament) fieldsToUpdate.push('departament');

    // Если есть изменения, обновляем профиль
    if (fieldsToUpdate.length > 0) {
      await adminStore.updateAdmin(formData.value, fieldsToUpdate);

      $q.notify({
        color: 'positive',
        message: 'Профиль успешно обновлен',
        icon: 'check_circle'
      });

      editMode.value = false;
    } else {
      $q.notify({
        color: 'warning',
        message: 'Нет изменений для сохранения',
        icon: 'info'
      });
    }
  } catch {
    $q.notify({
      color: 'negative',
      message: 'Не удалось обновить профиль',
      icon: 'error'
    });
  }
};

const changePassword = async () => {
  if (passwordData.value.newPassword !== passwordData.value.confirmPassword) {
    $q.notify({
      color: 'negative',
      message: 'Пароли не совпадают',
      icon: 'error'
    });
    return;
  }

  passwordLoading.value = true;

  try {
    await new Promise(resolve => setTimeout(resolve, 1000));

    $q.notify({
      color: 'positive',
      message: 'Пароль успешно изменен',
      icon: 'check_circle'
    });

    // Сбрасываем поля формы
    passwordData.value = {
      currentPassword: '',
      newPassword: '',
      confirmPassword: ''
    };
  } catch {
    $q.notify({
      color: 'negative',
      message: 'Не удалось изменить пароль',
      icon: 'error'
    });
  } finally {
    passwordLoading.value = false;
  }
};
</script>

<style scoped>
.content {
  padding: 0 40px 40px;
}

.page-title {
  color: rgba(90, 92, 105, 1);
  font-weight: bold;
  font-size: 32px;
  margin-bottom: 30px;
}

.password-card {
  width: 100%;
  border-radius: 8px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}
</style>
