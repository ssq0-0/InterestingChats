import { fetchWithToken } from './tokenUtils.js'; // Импортируем функцию fetchWithToken

let currentPage = 1;
const itemsPerPage = 10;
const debounceDelay = 500; // Задержка в миллисекундах

// Функция дебаунсинга
function debounce(func, delay) {
    let timeoutId;
    return function (...args) {
        if (timeoutId) clearTimeout(timeoutId);
        timeoutId = setTimeout(() => func.apply(this, args), delay);
    };
}

export async function fetchChats() {
    try {
        const response = await fetchWithToken('http://localhost:8000/getAllChats', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            },
        });

        if (!response.ok) throw new Error('Network response was not ok');

        const { Data: chats } = await response.json();
        renderChats(chats);
    } catch (error) {
        console.error('Error loading chats:', error);
    }
}

export function renderChats(chats) {
    const chatContainer = document.getElementById('chatsContainer');
    const totalPages = Math.ceil(chats.length / itemsPerPage);

    // Если нет чатов
    if (!chats || chats.length === 0) {
        chatContainer.innerHTML = `
            <p>No chats found.</p>
            <button id="createChat">Create new chat</button>
        `;
        return;
    }

    const start = (currentPage - 1) * itemsPerPage;
    const end = start + itemsPerPage;
    const chatsToDisplay = chats.slice(start, end);

    chatContainer.innerHTML = `
        <h1>Chats</h1>
        <input type="text" id="searchInput" placeholder="Search chats...">
        <ul id="searchResults"></ul>
        <button id="createChat">Create new chat</button>
        <ul id="chatsList">
            ${chatsToDisplay.map(chat => `
                <li>
                    ${chat.chat_name}
                    <button class="join-chat-btn" data-chat-id="${chat.id}">Join chat</button>
                </li>
            `).join('')}
        </ul>
        ${totalPages > 1 ? `
            <button id="prevPage" ${currentPage === 1 ? 'disabled' : ''}>Prev</button>
            <button id="nextPage" ${currentPage === totalPages ? 'disabled' : ''}>Next</button>
        ` : ''}
    `;

    // Добавляем обработчики событий для пагинации
    if (totalPages > 1) {
        document.getElementById('prevPage').addEventListener('click', () => {
            if (currentPage > 1) {
                currentPage--;
                fetchChats(); // Загрузка данных для новой страницы
            }
        });

        document.getElementById('nextPage').addEventListener('click', () => {
            if (currentPage < totalPages) {
                currentPage++;
                fetchChats(); // Загрузка данных для новой страницы
            }
        });
    }

    // Добавляем обработчик события для поля поиска после его добавления в DOM
    const searchInput = document.getElementById('searchInput');
    if (searchInput) {
        searchInput.addEventListener('input', debounce((e) => {
            const query = e.target.value.trim();

            if (query.length > 0) {
                searchChats(query);
            } else {
                clearSearchResults(); // Очистить результаты, если поле поиска пустое
            }
        }, debounceDelay));
    }

    // Добавляем обработчик для создания нового чата
    const createChatButton = document.getElementById('createChat');
    if (createChatButton) {
        createChatButton.addEventListener('click', () => {
            window.location.href = './chatCreateForm.html';
        });
    }
}

async function searchChats(query) {
    try {
        const response = await fetchWithToken(`http://localhost:8000/getChatBySymbol?chatName=${encodeURIComponent(query)}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            },
        });

        if (!response.ok) throw new Error('Network response was not ok');

        const { Data: chats } = await response.json();
        renderSearchResults(chats);
    } catch (error) {
        console.error('Error during search:', error);
    }
}

function renderSearchResults(chats) {
    const searchResults = document.getElementById('searchResults');
    const chatsList = document.getElementById('chatsList');

    if (chats.length > 0) {
        searchResults.innerHTML = chats.map(chat => `
            <li>
                ${chat.chat_name}
                ${chat.members && chat.members.length > 0 ? chat.members.length : ''}
                <button class="join-chat-btn" data-chat-id="${chat.id}">Join chat</button>
            </li>
        `).join('');
        // Скрываем основной список чатов
        chatsList.style.display = 'none';
    } else {
        // Если ничего не найдено, показываем основной список чатов и очищаем результаты поиска
        searchResults.innerHTML = '<li>No results found.</li>';
    }
}

function clearSearchResults() {
    const searchResults = document.getElementById('searchResults');
    const chatsList = document.getElementById('chatsList');

    searchResults.innerHTML = '';
    // Показываем основной список чатов
    chatsList.style.display = 'block';
}
