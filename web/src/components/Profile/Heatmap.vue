<script setup>
import { computed } from 'vue';

const props = defineProps({
  userSubmitRecords: Array
});

const rows = 5;

// 图表数据计算
const sparklineData = computed(() => {
  const counts = {};
  const now = new Date();
  const startDate = new Date();
  startDate.setDate(now.getDate() - 119);

  for (let d = new Date(startDate); d <= now; d.setDate(d.getDate() + 1)) {
    const dateString = d.toISOString().split('T')[0];
    counts[dateString] = 0;
  }

  props.userSubmitRecords.forEach(record => {
    const date = new Date(record.Time);
    if (date >= startDate && date <= now) {
      const dateString = date.toISOString().split('T')[0];
      counts[dateString] = (counts[dateString] || 0) + 1;
    }
  });

  return Object.keys(counts)
    .map(date => ({ date, count: counts[date] }))
    .sort((a, b) => new Date(a.date) - new Date(b.date));
});

const maxCount = computed(() => 
  Math.max(...sparklineData.value.map(d => d.count), 1)
);

// 渲染热力图
const heatColor = (count) => {
  if (count === 0) {
    return '#e0e0e0';
  }
  const ratio = count / maxCount.value;
  const hue = 120;
  const saturation = 70;
  const lightness = 85 - ratio * 50; 
  return `hsl(${hue}, ${saturation}%, ${lightness}%)`;
};

const cols = computed(() =>
  Math.ceil(sparklineData.value.length / rows)
);

const gridData = computed(() => {
  const arr = [];
  for (let c = 0; c < cols.value; c++) {
    const colArr = [];
    for (let r = 0; r < rows; r++) {
      const idx = c * rows + r;
      colArr.push(sparklineData.value[idx] || { date: '', count: 0 });
    }
    arr.push(colArr);
  }
  return arr;
});
</script>

<template>
    <div class="heatmap-grid">
        <template v-for="(colArr, cIdx) in gridData" :key="cIdx">
            <template v-for="(cell, rIdx) in colArr" :key="cIdx + '-' + rIdx">
                <v-tooltip bottom>
                    <template #activator="{ props }">
                        <div v-bind="props" class="heat-cell" :style="{ backgroundColor: heatColor(cell.count) }"></div>
                    </template>
                    <span>{{ cell.date }}: {{ cell.count }}</span>
                </v-tooltip>
            </template>
        </template>
    </div>
</template>

<style scoped>
.heatmap-grid {
    display: grid;
    grid-template-rows: repeat(5, 1fr);
    grid-auto-columns: minmax(24px, 1fr);
    grid-auto-flow: column;
    gap: 2px;
    width: 100%;
    padding-bottom: 4px;
    padding-top: 10px;
    overflow-x: auto;
}

.heat-cell {
    width: 100%;
    aspect-ratio: 1;
}
</style>