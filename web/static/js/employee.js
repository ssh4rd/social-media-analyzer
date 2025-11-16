// Данные для одной строки статистики
const data = [
    {
        totalMoney: 7500,
        tripsCount: 25,
        avgMoney: 300,
        medianMoney: 280
    }
];

document.addEventListener("DOMContentLoaded", () => {
    // Рендер таблицы (одна строка)
    const tbody = document.getElementById("data-body");
    tbody.innerHTML = "";

    data.forEach(item => {
        const row = `
            <tr>
                <td>${item.totalMoney}</td>
                <td>${item.tripsCount}</td>
                <td>${item.avgMoney}</td>
                <td>${item.medianMoney}</td>
            </tr>`;
        tbody.insertAdjacentHTML("beforeend", row);
    });

    // === ГРАФИК: Общие траты ===
    const ctxBar = document.getElementById('barChart').getContext('2d');
    new Chart(ctxBar, {
        type: 'bar',
        data: {
            labels: ['Общие траты', 'Средние траты', 'Медиана'],
            datasets: [{
                label: 'Деньги',
                data: [data[0].totalMoney, data[0].avgMoney, data[0].medianMoney],
                backgroundColor: ['#2e8b57', '#3cb371', '#90ee90'],
                borderRadius: 8
            }]
        },
        options: {
            responsive: true,
            plugins: { legend: { display: false } },
            scales: { y: { beginAtZero: true } }
        }
    });
});