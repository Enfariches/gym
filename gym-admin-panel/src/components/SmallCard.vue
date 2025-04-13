<template>
  <div class="stat-card">
    <div class="stat-title">{{ title }}</div>
    <div class="stat-value">{{ value }}</div>
    <q-icon :name="icon" class="stat-icon" />
  </div>
</template>

<script setup lang="ts">
// eslint-disable-next-line @typescript-eslint/no-unused-vars
const props = defineProps<{
  title: string;
  value: string | number;
  icon: string;
}>();
</script>

<style scoped>
.stat-card {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
  position: relative;
  overflow: hidden;
  transition: transform 0.2s;
}

.stat-card:hover {
  transform: translateY(-4px);
}

.stat-title {
  color: rgba(78,115,223,1);
  font-size: 18px;
  font-weight: bold;
  margin-bottom: 8px;
}

.stat-value {
  font-size: 36px;
  font-weight: bold;
  color: rgba(90,92,105,1);
}

.stat-icon {
  position: absolute;
  right: 20px;
  top: 50%;
  transform: translateY(-50%);
  font-size: 48px;
  opacity: 0.1;
  color: rgba(78,115,223,1);
}
</style> 