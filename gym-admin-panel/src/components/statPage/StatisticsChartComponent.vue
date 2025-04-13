<template>
  <div>
        <!-- Диаграмма -->
    <div class="col-12 col-sm-6" style="width: 80%">
      <q-card flat bordered style="min-height: 300px">
        <q-card-section style="min-height: 300px">
          <canvas ref="canvasRef" style="margin: 0 auto; height: 250px"/>
        </q-card-section>
      </q-card>
    </div>
  </div>
</template>

<script setup lang="ts">

import {nextTick, onMounted, ref, watch, shallowRef } from 'vue'
import { Chart, registerables } from 'chart.js'
import { getStatisticByViewType } from 'src/helpers/StatisticsApi'

  export interface StatisticChartComponentProps {
    viewType: string;
    dateStart?: string;
    dateEnd?: string;
  }
  interface ChartDataObject{
    xValue: string,
    yValue: number,
  }
  interface ResponseItem{
    created_at: string,
    id: number,
    type: string,
    user_id: number,
    video_id: number
  }
  const props = withDefaults(defineProps<StatisticChartComponentProps>(), {
    viewType:'half'
  });

  const canvasRef = ref<HTMLCanvasElement | null>(null)
  const chartInstance = shallowRef<Chart | null>(null);
  const dataset = ref<Array<ChartDataObject> | null>(null)

  onMounted(async () => {
    watch(
      () => props.viewType, async () => {
        await prepareDataSet().then((data:Array<ChartDataObject> | null) => {
          if(chartInstance.value && data){

            const chart = chartInstance.value.data
            chart.labels = data.map((data) => data.xValue)
            if(chart.datasets[0])
              chart.datasets[0].data = data.map((data) => data.yValue)
            chartInstance.value.update()
          }
        })
      }
    );
    dataset.value = await prepareDataSet()
    await createChart()
  })

  const prepareDataSet = async (): Promise<Array<ChartDataObject> | null> => {
    const data = await getStatisticByViewType(props.viewType)
    const map = new Map()
    const arr: Array<ChartDataObject> = []

    data?.forEach((el:ResponseItem)=>{
        const date = new Date(el.created_at)
        const formattedDate = date.toLocaleDateString("ru-RU", {
          day: "2-digit",
          month: "2-digit"
        })

        if(map.has(formattedDate)){
          const current = map.get(formattedDate)
          map.set(formattedDate, current + 1)
        }
        else{
          map.set(formattedDate, 1)
        }
    })

    map.forEach((value, key) => {
      arr.push({
        xValue:key,
        yValue:value,
      })
    })

    return arr
  }

  const createChart = async () => {
    await nextTick(); // Ждём, пока Vue обновит DOM

    if (!canvasRef.value) {
      console.error("Canvas not found!");
      return;
    }

    const ctx = canvasRef.value.getContext("2d");

    if (chartInstance.value) {
      chartInstance.value.destroy();
      chartInstance.value = null; // Очищаем ссылку, чтобы избежать конфликтов
      await new Promise((resolve) => setTimeout(resolve, 50));
    }
    if (!ctx) {
      console.error("Failed to get 2D context!");
      return;
    }

    if(!dataset.value) return

    Chart.register(...registerables)

    chartInstance.value = new Chart(ctx,{
      type: 'line',
      data: {
        labels: dataset.value.map((data) => data.xValue),
        datasets: [
          {
            label: 'Просмотры',
            backgroundColor: '#1976D2',
            data: dataset.value.map((data) => data.yValue),
          },
        ],
      },
      options: {
        scales: {
          y: {
            beginAtZero: true,
            title: {
              display: true,
              text: 'Количество просмотров',
            },
          },
          x: {
            title: {
              display: true,
              text: 'Дни недели',
            },
          },
        },
      },
    })
  }


</script>
