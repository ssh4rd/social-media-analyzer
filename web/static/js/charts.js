document.addEventListener("DOMContentLoaded", () => {
    const labels = data.map(d => d.group);

    // График 1: Средние лайки на пост
    const avgLikes = data.map(d => d.avgLikes);

    const ctxLine1 = document.getElementById('lineChart').getContext('2d');
    new Chart(ctxLine1, {
        type: 'line',
        data: {
            labels: labels,
            datasets: [{
                label: 'Сред. лайки на пост',
                data: avgLikes,
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

    // График 2: Средние комментарии на пост
    const avgComments = data.map(d => d.avgComments);

    const ctxLine2 = document.getElementById('barChart').getContext('2d');
    new Chart(ctxLine2, {
        type: 'line',
        data: {
            labels: labels,
            datasets: [{
                label: 'Сред. кол-во комментов',
                data: avgComments,
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
