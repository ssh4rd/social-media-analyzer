const data = [
    { id: 1, date: "01.01.2025", type: "Тип A", duration: "2 дня", destination: "Куда A", money: 500 },
    { id: 2, date: "05.01.2025", type: "Тип B", duration: "3 дня", destination: "Куда B", money: 800 },
    { id: 3, date: "10.01.2025", type: "Тип C", duration: "1 день", destination: "Куда C",  money: 300 },
    { id: 4, date: "12.01.2025", type: "Тип D", duration: "4 дня", destination: "Куда D", money: 700 },
    { id: 5, date: "15.01.2025", type: "Тип E", duration: "5 дней", destination: "Куда E",  money: 1200 },
    { id: 6, date: "18.01.2025", type: "Тип F", duration: "2 дня", destination: "Куда F", money: 400 },
    { id: 7, date: "20.01.2025", type: "Тип G", duration: "3 дня", destination: "Куда G", money: 900 },
    { id: 8, date: "22.01.2025", type: "Тип H", duration: "2 дня", destination: "Куда H", money: 600 },
    { id: 9, date: "25.01.2025", type: "Тип I", duration: "1 день", destination: "Куда I",  money: 250 },
    { id: 10, date: "28.01.2025", type: "Тип J", duration: "4 дня", destination: "Куда J", money: 1000 }
];

const rowsPerPage = 5;
let currentPage = 1;

function renderTable(page = 1) {
    const start = (page - 1) * rowsPerPage;
    const end = start + rowsPerPage;
    const pageData = data.slice(start, end);

    const tbody = document.getElementById("data-body");
    tbody.innerHTML = "";

    pageData.forEach(item => {
        const row = `
            <tr>
                <td>${item.type}</td>
                <td>${item.date}</td>
                <td>${item.duration}</td>
                <td>${item.destination}</td>
                <td>${item.money}</td>
                <td>
                    <a href="/employee/${item.id}">
                        <button class="btn btn-success btn-sm")">
                            Открыть
                        </button>
                    </a>
                </td>
            </tr>`;
        tbody.insertAdjacentHTML("beforeend", row);
    });
}

function renderPagination() {
    const totalPages = Math.ceil(data.length / rowsPerPage);
    const pagination = document.getElementById("pagination");
    pagination.innerHTML = "";

    const pageLimit = 10; // макс отображаемых страниц
    let startPage = Math.max(1, currentPage - Math.floor(pageLimit / 2));
    let endPage = startPage + pageLimit - 1;
    if (endPage > totalPages) {
        endPage = totalPages;
        startPage = Math.max(1, endPage - pageLimit + 1);
    }

    // Кнопка Назад
    pagination.insertAdjacentHTML("beforeend", `
        <li class="page-item ${currentPage === 1 ? 'disabled' : ''}">
            <button class="page-link">&laquo;</button>
        </li>
    `);

    // Левая многоточие
    if (startPage > 1) {
        pagination.insertAdjacentHTML("beforeend", `
            <li class="page-item"><button class="page-link">1</button></li>
            <li class="page-item disabled"><span class="page-link">...</span></li>
        `);
    }

    // Основные страницы
    for (let i = startPage; i <= endPage; i++) {
        pagination.insertAdjacentHTML("beforeend", `
            <li class="page-item ${i === currentPage ? 'active' : ''}">
                <button class="page-link">${i}</button>
            </li>
        `);
    }

    // Правая многоточие
    if (endPage < totalPages) {
        pagination.insertAdjacentHTML("beforeend", `
            <li class="page-item disabled"><span class="page-link">...</span></li>
            <li class="page-item"><button class="page-link">${totalPages}</button></li>
        `);
    }

    // Кнопка Вперёд
    pagination.insertAdjacentHTML("beforeend", `
        <li class="page-item ${currentPage === totalPages ? 'disabled' : ''}">
            <button class="page-link">&raquo;</button>
        </li>
    `);

    // Обработчики кликов
    const buttons = pagination.querySelectorAll(".page-link");
    buttons.forEach(btn => {
        btn.addEventListener("click", () => {
            const text = btn.textContent;
            if (text === '«' && currentPage > 1) currentPage--;
            else if (text === '»' && currentPage < totalPages) currentPage++;
            else if (!isNaN(text)) currentPage = Number(text);

            renderTable(currentPage);
            renderPagination();
        });
    });
}

document.addEventListener("DOMContentLoaded", () => {
    renderTable();
    renderPagination();
});