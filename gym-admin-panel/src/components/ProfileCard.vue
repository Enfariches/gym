<template>
  <q-card class="profile-card">
    <q-card-section class="bg-primary text-white">
      <div class="row items-center no-wrap">
        <div class="col">
          <div class="text-h6">{{ admin?.name || '' }} {{ admin?.surname || '' }}</div>
          <div class="text-subtitle2">{{ admin?.email || '' }}</div>
        </div>
        <div class="col-auto">
          <q-btn
            v-if="!editMode"
            flat
            round
            icon="edit"
            @click="$emit('edit')"
            size="md"
          />
          <div v-else class="row q-gutter-x-sm">
            <q-btn
              flat
              round
              icon="check"
              color="positive"
              @click="$emit('save')"
              size="md"
            />
            <q-btn
              flat
              round
              icon="close"
              color="negative"
              @click="$emit('cancel')"
              size="md"
            />
          </div>
        </div>
      </div>
    </q-card-section>

    <q-card-section>
      <div v-if="!editMode" class="q-pa-md">
        <div class="row q-mb-md items-center">
          <div class="col-3 text-subtitle1 text-weight-medium">Имя:</div>
          <div class="col-9">{{ admin?.name || 'Не указано' }}</div>
        </div>
        <div class="row q-mb-md items-center">
          <div class="col-3 text-subtitle1 text-weight-medium">Фамилия:</div>
          <div class="col-9">{{ admin?.surname || 'Не указано' }}</div>
        </div>
        <div class="row q-mb-md items-center">
          <div class="col-3 text-subtitle1 text-weight-medium">Email:</div>
          <div class="col-9">{{ admin?.email || 'Не указано' }}</div>
        </div>
        <div class="row items-center">
          <div class="col-3 text-subtitle1 text-weight-medium">Отдел:</div>
          <div class="col-9">{{ admin?.departament || 'Не указано' }}</div>
        </div>
      </div>

      <div v-else class="q-pa-md">
        <q-input
          v-model="formData.name"
          label="Имя"
          outlined
          dense
          class="q-mb-md"
        />
        <q-input
          v-model="formData.surname"
          label="Фамилия"
          outlined
          dense
          class="q-mb-md"
        />
        <q-input
          v-model="formData.email"
          label="Email"
          type="email"
          outlined
          dense
          class="q-mb-md"
          readonly
        />
        <q-input
          v-model="formData.departament"
          label="Отдел"
          outlined
          dense
        />
      </div>
    </q-card-section>
  </q-card>
</template>

<script setup lang="ts">
import { ref, watch, defineProps, defineEmits } from 'vue';
import type { Admin } from '../../protogen/v1/users/admin';

const props = defineProps<{
  admin: Admin | null;
  editMode: boolean;
}>();

const emit = defineEmits<{
  (e: 'edit'): void;
  (e: 'save'): void;
  (e: 'cancel'): void;
  (e: 'update:formData', formData: Partial<Admin>): void;
}>();

const formData = ref({
  name: '',
  surname: '',
  email: '',
  departament: ''
});

// Обновляем форму при изменении данных администратора
watch(() => props.admin, (newAdmin) => {
  if (newAdmin) {
    formData.value = {
      name: newAdmin.name,
      surname: newAdmin.surname,
      email: newAdmin.email,
      departament: newAdmin.departament
    };
  }
}, { immediate: true });

// Уведомляем родительский компонент об изменениях в форме
watch(formData, (newFormData) => {
  emit('update:formData', newFormData);
}, { deep: true });
</script>

<style scoped>
.profile-card {
  width: 100%;
  border-radius: 8px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}
</style>
