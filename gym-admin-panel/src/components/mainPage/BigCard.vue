<template>
  <div class="dashboard-card">
    <h2 class="card-title">
      <q-icon :name="icon" class="q-mr-sm" />
      {{ title }}
    </h2>
    <slot></slot>
  </div>
</template>

<script setup lang="ts">
// eslint-disable-next-line @typescript-eslint/no-unused-vars
const props = defineProps<{
  title: string;
  icon: string;
}>();
</script>

<style scoped>
.dashboard-card {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
}

.card-title {
  color: rgba(90,92,105,1);
  font-weight: bold;
  font-size: 20px;
  margin-bottom: 20px;
  padding-bottom: 15px;
  border-bottom: 1px solid rgba(227,230,240,1);
  display: flex;
  align-items: center;
}
</style> 