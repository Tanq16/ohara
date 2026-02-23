let timelineChart = null;

const CHART_COLORS = {
  blue:     "#89b4fa",
  mauve:    "#cba6f7",
  green:    "#a6e3a1",
  peach:    "#fab387",
  red:      "#f38ba8",
  yellow:   "#f9e2af",
  teal:     "#94e2d5",
  sky:      "#89dceb",
  sapphire: "#74c7ec",
  pink:     "#f5c2e7",
  lavender: "#b4befe",
  flamingo: "#f2cdcd",
};

const PALETTE = Object.values(CHART_COLORS);

const chartTheme = {
  text: "#a6adc8",
  grid: "#313244",
  gridLight: "#45475a",
};

Chart.defaults.color = chartTheme.text;
Chart.defaults.borderColor = chartTheme.grid;
Chart.defaults.font.family = "'Inter', system-ui, sans-serif";

async function initDashboard() {
  if (!allTouchpoints.length) {
    await loadTouchpoints();
  }

  renderTimelineChart(allTouchpoints);
  renderTouchpointList();
}

function renderTimelineChart(tps) {
  const ctx = document.getElementById("chart-timeline");
  if (timelineChart) timelineChart.destroy();

  const now = new Date();
  const labels = [];
  const monthKeys = [];
  for (let i = 11; i >= 0; i--) {
    const d = new Date(now.getFullYear(), now.getMonth() - i, 1);
    const key = d.toISOString().slice(0, 7);
    monthKeys.push(key);
    labels.push(d.toLocaleDateString(undefined, { month: "short", year: "2-digit" }));
  }

  const countByMonth = {};
  const categoriesByMonth = {};
  monthKeys.forEach((k) => {
    countByMonth[k] = 0;
    categoriesByMonth[k] = new Set();
  });

  tps.forEach((tp) => {
    const key = tp.date.slice(0, 7);
    if (countByMonth[key] !== undefined) {
      countByMonth[key]++;
      categoriesByMonth[key].add(tp.category);
    }
  });

  const counts = monthKeys.map((k) => countByMonth[k]);
  const diversity = monthKeys.map((k) => categoriesByMonth[k].size);

  timelineChart = new Chart(ctx, {
    type: "line",
    data: {
      labels,
      datasets: [
        {
          label: "Touchpoints",
          data: counts,
          borderColor: CHART_COLORS.blue,
          backgroundColor: CHART_COLORS.blue + "33",
          fill: true,
          tension: 0.3,
          yAxisID: "y",
        },
        {
          label: "Categories",
          data: diversity,
          borderColor: CHART_COLORS.green,
          backgroundColor: "transparent",
          fill: false,
          tension: 0.3,
          yAxisID: "y1",
        },
      ],
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      interaction: {
        mode: "index",
        intersect: false,
      },
      plugins: {
        legend: { display: true, labels: { color: chartTheme.text } },
        tooltip: {
          mode: "index",
          intersect: false,
          backgroundColor: "#181825",
          titleColor: "#cdd6f4",
          bodyColor: "#a6adc8",
          borderColor: "#45475a",
          borderWidth: 1,
          cornerRadius: 8,
          padding: 10,
        },
      },
      scales: {
        x: {
          ticks: { color: chartTheme.text, font: { size: 11 } },
          grid: { display: false },
        },
        y: {
          position: "left",
          title: { display: true, text: "Touchpoints", color: chartTheme.text },
          beginAtZero: true,
          ticks: { color: chartTheme.text, stepSize: 1, font: { size: 11 } },
          grid: { color: chartTheme.grid, drawBorder: false },
        },
        y1: {
          position: "right",
          title: { display: true, text: "Categories", color: chartTheme.text },
          beginAtZero: true,
          ticks: { color: chartTheme.text, stepSize: 1, font: { size: 11 } },
          grid: { drawOnChartArea: false },
        },
      },
    },
  });
}
