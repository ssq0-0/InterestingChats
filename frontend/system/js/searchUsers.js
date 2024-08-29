import { fetchWithToken } from './tokenUtils.js'; // Импортируем функцию fetchWithToken

export function initializeUserSearch() {
    console.log("DOM fully loaded and parsed");
    const userSearchContainer = document.getElementById('userSearchContainer');

    userSearchContainer.innerHTML = `
        <h1>Search users</h1>
        <input type="text" id="searchInput" placeholder="Search users...">
        <ul id="searchResults"></ul>
        <ul id="usersList"></ul>
    `;

    const searchInput = document.getElementById('searchInput');
    searchInput.addEventListener('input', debounce(() => {
        const query = searchInput.value.trim();
        if (query.length > 0) {
            searchUsers(query);
        } else {
            clearSearchResults();
        }
    }, 500)); // Задержка в 500 мс
}

function debounce(func, delay) {
    let timeoutId;
    return function (...args) {
        if (timeoutId) clearTimeout(timeoutId);
        timeoutId = setTimeout(() => func.apply(this, args), delay);
    };
}

function clearSearchResults() {
    document.getElementById('searchResults').innerHTML = '';
}

export async function searchUsers(query) {
    try {
        const response = await fetchWithToken(`http://localhost:8000/searchUsers?symbols=${encodeURIComponent(query)}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            },
        });

        if (!response.ok) throw new Error('Network response was not ok');

        const result = await response.json();
        
        // Проверка на наличие ошибок в ответе и данных
        if (result.Errors) {
            console.error('Server Errors:', result.Errors);
            renderSearchResults([]); // Показать, что пользователи не найдены
        } else if (Array.isArray(result.Data)) {
            renderSearchResults(result.Data);
        } else {
            console.error('Invalid data format:', result);
            renderSearchResults([]); // Показать, что пользователи не найдены
        }
    } catch (error) {
        console.error('Error during user search:', error);
        renderSearchResults([]); // Показать, что пользователи не найдены в случае ошибки
    }
}

function renderSearchResults(users) {
    const searchResults = document.getElementById('searchResults');

    if (users && users.length > 0) {
        searchResults.innerHTML = users.map(user => `
            <li>
                ${user.username} (${user.email})
                <button class="view-profile-btn" data-user-id="${user.id}">View Profile</button>
            </li>
        `).join('');
    } else {
        searchResults.innerHTML = '<li>No users found.</li>';
    }

    document.querySelectorAll('.view-profile-btn').forEach(button => {
        button.addEventListener('click', () => {
            const userId = button.getAttribute('data-user-id');
            viewUserProfile(userId);
        });
    });
}

function viewUserProfile(userId) {
    // Реализуйте переход на страницу профиля пользователя
    window.location.href = `./user_profile.html?userID=${userId}`;
}
