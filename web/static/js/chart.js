document.addEventListener("DOMContentLoaded", () => {

    // Метки — типы
    const labels = data.map(d => d.type);

    // Данные для графика 1 — деньги
    const money = data.map(d => d.money);

    // Данные для графика 2 — длительность в числах
    const duration = data.map(d => parseInt(d.duration)); // "3 дня" → 3

    // === ГРАФИК 1: Деньги ===
    const ctxLine1 = document.getElementById('lineChart').getContext('2d');
    new Chart(ctxLine1, {
        type: 'line',
        data: {
            labels: labels,
            datasets: [{
                label: 'Сколько денег потрачено',
                data: money,
                borderColor: '#2e8b57',
                backgroundColor: 'rgba(46,139,87,0.2)',
                tension: 0.3,
                pointStyle: 'circle',
                pointRadius: 6,
                pointHoverRadius: 10,
                pointBackgroundColor: '#3cb371'
            }]
        },
        options: {
            responsive: true,
            plugins: { legend: { display: true } },
            scales: { y: { beginAtZero: true } }
        }
    });

    // === ГРАФИК 2: Длительность ===
    const ctxLine2 = document.getElementById('barChart').getContext('2d');
    new Chart(ctxLine2, {
        type: 'line',
        data: {
            labels: labels,
            datasets: [{
                label: 'Длительность поездки (дни)',
                data: duration,
                borderColor: '#3cb371',
                backgroundColor: 'rgba(60,179,113,0.2)',
                tension: 0.3,
                pointStyle: 'rectRot',
                pointRadius: 6,
                pointHoverRadius: 10,
                pointBackgroundColor: '#90ee90'
            }]
        },
        options: {
            responsive: true,
            plugins: { legend: { display: true } },
            scales: { y: { beginAtZero: true } }
        }
    });

});
