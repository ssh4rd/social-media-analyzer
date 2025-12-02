const data = [
    { group: "Группа A", members: 10, postsLastWeek: 25, totalLikes: 500, avgLikes: 20, maxLikes: 50, avgComments: 5, avgImages: 2 },
    { group: "Группа B", members: 15, postsLastWeek: 30, totalLikes: 800, avgLikes: 26.7, maxLikes: 70, avgComments: 8, avgImages: 3 },
    { group: "Группа C", members: 8, postsLastWeek: 20, totalLikes: 300, avgLikes: 15, maxLikes: 35, avgComments: 4, avgImages: 1 },
    { group: "Группа D", members: 12, postsLastWeek: 28, totalLikes: 700, avgLikes: 25, maxLikes: 60, avgComments: 6, avgImages: 2 },
    { group: "Группа E", members: 20, postsLastWeek: 40, totalLikes: 1200, avgLikes: 30, maxLikes: 90, avgComments: 10, avgImages: 4 },
    { group: "Группа F", members: 9, postsLastWeek: 22, totalLikes: 400, avgLikes: 18.2, maxLikes: 45, avgComments: 5, avgImages: 2 },
    { group: "Группа G", members: 14, postsLastWeek: 35, totalLikes: 900, avgLikes: 25.7, maxLikes: 65, avgComments: 7, avgImages: 3 },
    { group: "Группа H", members: 11, postsLastWeek: 27, totalLikes: 600, avgLikes: 22.2, maxLikes: 55, avgComments: 6, avgImages: 2 },
    { group: "Группа I", members: 7, postsLastWeek: 18, totalLikes: 250, avgLikes: 13.9, maxLikes: 30, avgComments: 3, avgImages: 1 },
    { group: "Группа J", members: 16, postsLastWeek: 38, totalLikes: 1000, avgLikes: 26.3, maxLikes: 80, avgComments: 9, avgImages: 3 }
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
                <td>${item.group}</td>
                <td>${item.members}</td>
                <td>${item.postsLastWeek}</td>
                <td>${item.totalLikes}</td>
                <td>${item.avgLikes}</td>
                <td>${item.maxLikes}</td>
                <td>${item.avgComments}</td>
                <td>${item.avgImages}</td>
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