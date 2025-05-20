<template>
  <div class="schedule-element">
    <div class="schedule-element__head">
      <div>{{ dayName }}</div>
    </div>
    <div class="schedule-element__body">
      <div v-if="items.length === 0" class="body__empty">
        <div>Нет видео на этот день</div>
      </div>
      <div v-else class="body__list">
        <div
          v-for="(schedule, idx) in items"
          :key="idx"
          class="schedule-item"
          :class="{'schedule-item-odd': idx % 2 !== 0}"
        >
          <div class="schedule-item__content">
            <div class="schedule-item__title">
              <div class="video__label">Видео</div>
              <span>{{ schedule.VideoTitle }}</span>
            </div>
            <div class="schedule-item__time">
              <div class="time__label">Время</div>
              <span>{{ schedule.Time }}</span>
            </div>
            <div class="schedule-item__actions">
              <button
                class="action-btn edit"
                title="Редактировать"
                @click="$emit('edit', {
                  videoName: schedule.VideoTitle,
                  scheduleId: schedule.ID,
                  dayOrder: dayOrder
                })"
              >
                <q-icon name="edit" color="primary" />
                <span class="action-tooltip">Редактировать</span>
              </button>
              <button
                class="action-btn delete"
                title="Удалить"
                @click="$emit('delete', schedule.ID)"
              >
                <q-icon name="delete" color="negative" />
                <span class="action-tooltip">Удалить</span>
              </button>
            </div>
          </div>
        </div>
      </div>
      <button
        class="btn btn-primary add-schedule"
        @click="$emit('add', dayOrder)"
      >
        <q-icon name="add" color="primary" /> Добавить видео
      </button>
    </div>
  </div>
</template>

<script lang="ts" setup>
interface ScheduleItem {
  ID: string;
  Time: string;
  VideoTitle: string;
}

interface Video {
  title: string;
}

defineProps<{
  dayName: string;
  dayOrder: number;
  items: ScheduleItem[];
  videos: Video[];
}>();

defineEmits<{
  (e: 'edit', data: { videoName: string; scheduleId: string; dayOrder: number }): void;
  (e: 'delete', id: string): void;
  (e: 'add', dayOrder: number): void;
}>();
</script>

<style scoped>
.schedule-element {
  border: 1px solid #ddd;
  border-radius: 8px;
  overflow: hidden;
  background: white;
}

.schedule-element__head {
  background: #f5f5f5;
  padding: 12px;
  font-weight: bold;
  border-bottom: 1px solid #ddd;
}

.schedule-element__body {
  padding: 15px;
}

.body__empty {
  text-align: center;
  padding: 20px;
  color: #666;
}

.schedule-item {
  padding: 10px;
  border-radius: 4px;
  margin-bottom: 8px;
}

.schedule-item-odd {
  background-color: #f9f9f9;
}

.schedule-item__content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.video__label, .time__label {
  font-size: 0.8em;
  color: #666;
  margin-bottom: 4px;
}

.schedule-item__actions {
  display: flex;
  gap: 8px;
}

.action-btn {
  background: none;
  border: none;
  cursor: pointer;
  padding: 4px 8px;
  position: relative;
}

.action-btn.edit {
  color: #2196F3;
}

.action-btn.delete {
  color: #F44336;
}

.action-tooltip {
  display: none;
  position: absolute;
  background: rgba(0,0,0,0.8);
  color: white;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 0.8em;
  bottom: 100%;
  left: 50%;
  transform: translateX(-50%);
}

.action-btn:hover .action-tooltip {
  display: block;
}

.add-schedule {
  width: 100%;
  margin-top: 15px;
}
</style>
