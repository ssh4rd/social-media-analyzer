document.addEventListener("DOMContentLoaded", () => {
    // Use chart data from backend
    const labels = chartData.subscribers.map(sub => `${sub} подписчиков`);

    // Chart 1: Dependence of average likes on subscribers
    const ctxLine1 = document.getElementById('lineChart').getContext('2d');
    new Chart(ctxLine1, {
        type: 'line',
        data: {
            labels: labels,
            datasets: [{
                label: 'Сред. лайки на пост',
                data: chartData.avgLikes,
                borderColor: '#3399ff',
                backgroundColor: 'rgba(51,153,255,0.2)',
                tension: 0.3,
                pointStyle: 'circle',
                pointRadius: 6,
                pointHoverRadius: 10,
                pointBackgroundColor: '#66b3ff'
            }]
        },
        options: {
            responsive: true,
            plugins: { legend: { display: true } },
            scales: { y: { beginAtZero: true } }
        }
    });

    // Chart 2: Dependence of average comments on subscribers
    const ctxLine2 = document.getElementById('barChart').getContext('2d');
    new Chart(ctxLine2, {
        type: 'line',
        data: {
            labels: labels,
            datasets: [{
                label: 'Сред. кол-во комментов',
                data: chartData.avgComments,
                borderColor: '#66ccff',
                backgroundColor: 'rgba(102,204,255,0.2)',
                tension: 0.3,
                pointStyle: 'rectRot',
                pointRadius: 6,
                pointHoverRadius: 10,
                pointBackgroundColor: '#99ddff'
            }]
        },
        options: {
            responsive: true,
            plugins: { legend: { display: true } },
            scales: { y: { beginAtZero: true } }
        }
    });
});
