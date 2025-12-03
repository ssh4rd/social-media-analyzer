document.getElementById('addGroupBtn').addEventListener('click', async function() {
    const groupLink = document.getElementById('groupLink').value.trim();
    const loadingSpinner = document.getElementById('loadingSpinner');
    const errorAlert = document.getElementById('errorAlert');
    const successAlert = document.getElementById('successAlert');
    const errorMessage = document.getElementById('errorMessage');
    const addBtn = document.getElementById('addGroupBtn');

    // Reset alerts
    errorAlert.style.display = 'none';
    successAlert.style.display = 'none';

    // Validate input
    if (!groupLink) {
        errorMessage.textContent = 'Пожалуйста, введите ссылку на группу.';
        errorAlert.style.display = 'block';
        return;
    }

    // Show loader and disable button
    loadingSpinner.style.display = 'flex';
    addBtn.disabled = true;

    try {
        const response = await fetch('/api/groups', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ link: groupLink })
        });

        if (!response.ok) {
            const data = await response.json().catch(() => ({}));
            throw new Error(data.message || 'Ошибка при добавлении группы. Попробуйте снова.');
        }

        // Success
        loadingSpinner.style.display = 'none';
        successAlert.style.display = 'block';
        document.getElementById('groupLink').value = '';

        // Reload table data
        setTimeout(() => {
            location.reload();
        }, 1500);
    } catch (error) {
        loadingSpinner.style.display = 'none';
        errorMessage.textContent = error.message;
        errorAlert.style.display = 'block';
    } finally {
        addBtn.disabled = false;
    }
});

// Allow Enter key to submit
document.getElementById('groupLink').addEventListener('keypress', function(event) {
    if (event.key === 'Enter') {
        document.getElementById('addGroupBtn').click();
    }
});
