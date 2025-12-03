const rowsPerPage = 5;
let currentPage = 1;
let data = [];

function extractDataFromTable() {
    const tbody = document.getElementById("data-body");
    const rows = tbody.querySelectorAll("tr");
    
    data = Array.from(rows).map(row => {
        const cells = row.querySelectorAll("td");
        return {
            group: cells[0]?.textContent?.trim() || "",
            members: cells[1]?.textContent?.trim() || "-",
            postsLastWeek: parseInt(cells[2]?.textContent?.trim()) || 0,
            totalLikes: parseInt(cells[3]?.textContent?.trim()) || 0,
            avgLikes: parseFloat(cells[4]?.textContent?.trim()) || 0,
            maxLikes: parseInt(cells[5]?.textContent?.trim()) || 0,
            avgComments: parseFloat(cells[6]?.textContent?.trim()) || 0,
            avgImages: cells[7]?.textContent?.trim() || "-"
        };
    });
}

function renderTable(page = 1) {
    const start = (page - 1) * rowsPerPage;
    const end = start + rowsPerPage;
    const pageData = data.slice(start, end);

    const tbody = document.getElementById("data-body");
    tbody.innerHTML = "";

    if (pageData.length === 0) {
        tbody.innerHTML = "<tr><td colspan='8' class='text-center'>Нет данных</td></tr>";
        return;
    }

    pageData.forEach(item => {
        const row = `
            <tr>
                <td>${item.group}</td>
                <td>${item.members}</td>
                <td>${item.postsLastWeek}</td>
                <td>${item.totalLikes}</td>
                <td>${typeof item.avgLikes === 'number' ? item.avgLikes.toFixed(2) : item.avgLikes}</td>
                <td>${item.maxLikes}</td>
                <td>${typeof item.avgComments === 'number' ? item.avgComments.toFixed(2) : item.avgComments}</td>
                <td>${item.avgImages}</td>
            </tr>`;
        tbody.insertAdjacentHTML("beforeend", row);
    });
}

function renderPagination() {
    const totalPages = Math.ceil(data.length / rowsPerPage);
    const pagination = document.getElementById("pagination");
    pagination.innerHTML = "";

    if (totalPages <= 1) return;

    const pageLimit = 10;
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
    extractDataFromTable();
    renderTable();
    renderPagination();
});
