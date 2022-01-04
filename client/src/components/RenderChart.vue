<template>
  <canvas id="chart"></canvas>
</template>

<script>
import { Chart } from 'chart.js'
import { defineComponent, onMounted } from 'vue'

export default defineComponent({
  setup() {
    const renderChart = () => {
      let ctx = 'chart'
      new Chart(ctx, {
        type: 'bar',
        data: {
          labels: ['赤', '橙', '緑', '青'],
          datasets: [
            {
              label: 'result',
              data: [40, 20, 50, 30],
              backgroundColor: [
                'rgba(255, 99, 132, 0.2)',
                'rgba(255, 159, 64, 0.2)',
                'rgba(75, 192, 192, 0.2)',
                'rgba(54, 162, 235, 0.2)'
              ],
              borderColor: [
                'rgb(255, 99, 132)',
                'rgb(255, 159, 64)',
                'rgb(75, 192, 192)',
                'rgb(54, 162, 235)'
              ],
              borderWidth: 1
            }
          ]
        },
        options: {
          indexAxis: 'y',
          scales: [
            {
              grid: {
                display: false,
                drawBorder: false,
                drawOnChartArea: false,
                drawTicks: false
              },
              ticks: {
                display: false
              }
            }
          ]
        }
      })
    }
    onMounted(() => {
      renderChart()
    })
    return {}
  }
})
</script>
