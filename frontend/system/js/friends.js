import { fetchWithToken } from "./tokenUtils.js";

const ITEMS_PER_PAGE = 20; // Количество элементов на странице
let currentPage = 1; // Текущая страница для друзей
let currentSubscriberPage = 1; // Текущая страница для подписчиков

// Функция загрузки друзей
export async function loadFriends(page = 1) {
    try {
        // Получаем данные друзей с сервера
        const response = await fetchWithToken(`http://localhost:8000/getFriends`, {
            method: 'GET'
        });

        if (!response.ok) {
            throw new Error('Ошибка при получении данных');
        }

        const data = await response.json();
        console.log("Список друзей:", data); // Выводим полный ответ для проверки

        if (data.Errors) {
            console.error('Ошибка при получении данных:', data.Errors);
            return;
        }

        if (data.Data === null || !Array.isArray(data.Data)) {
            console.error('Ожидался массив друзей, но получено:', data.Data);
            const friendsContainer = document.getElementById('friendsContainer');
            friendsContainer.innerHTML = `
                <p>Нет друзей</p>
                <button id="subscribersButton">Подписчики</button>
            `;
            return;
        }
        // Проверьте, что data.Data действительно является массивом
        if (Array.isArray(data.Data)) {
            console.log("Передача данных в loadUserData:", data.Data);
            localStorage.setItem('friends', JSON.stringify(data.Data));
            loadUserData(data.Data, page);
        } else {
            console.error('Ожидался массив друзей, но получено:', data.Data);
            const friendsContainer = document.getElementById('friendsContainer');
            friendsContainer.innerHTML += '<p>Не удалось загрузить список друзей.</p>';
        }
    } catch (error) {
        console.error('Ошибка при загрузке данных:', error);
    }
}

// Функция отображения друзей с учетом пагинации
export function loadUserData(friends, page) {
    console.log("Массив в функции loadUserData:", friends);
    const friendsContainer = document.getElementById('friendsContainer');
    friendsContainer.innerHTML = `
        <button id="subscribersButton">Подписчики</button>
    `;

    // Проверяем, что friends действительно является массивом и имеет элементы
    if (!Array.isArray(friends) || friends.length === 0) {
        console.error('Друзей не найдено:', friends);
        friendsContainer.innerHTML += '<p>Друзей не найдено</p>';
        return;
    }

    // Отображаем друзей на странице
    friends.forEach(user => {
        const userElement = document.createElement('div');
        userElement.classList.add('user-item');
        userElement.innerHTML = `
            <div class="profile-avatar" style="border-radius: 50%; overflow: hidden; width: 50px; height: 50px;">
                <img src="${user.avatar}" alt="${user.username}" style="width: 100%; height: 100%; object-fit: cover;" />
            </div>
            <div class="user-info">
                <a href="user_profile.html?userID=${user.id}" class="user-username" style="color: blue; text-decoration: underline;">${user.username}</a>
            </div>
        `;
        friendsContainer.appendChild(userElement);
    });

    // Добавляем кнопки пагинации только если количество друзей больше ITEMS_PER_PAGE
    if (friends.length > ITEMS_PER_PAGE) {
        const pagination = document.createElement('div');
        pagination.classList.add('pagination');
        pagination.innerHTML = `
            <button id="prevPage" ${page === 1 ? 'disabled' : ''}>Предыдущая</button>
            <span>Страница ${page}</span>
            <button id="nextPage">Следующая</button>
        `;
        friendsContainer.appendChild(pagination);

        // Обработчики событий для кнопок пагинации
        document.getElementById('prevPage').addEventListener('click', () => {
            if (page > 1) {
                loadFriends(page - 1);
            }
        });

        document.getElementById('nextPage').addEventListener('click', () => {
            loadFriends(page + 1);
        });
    }
}



// Функция загрузки подписчиков
async function loadSubscribers(page = 1) {
    try {
        // Получаем данные подписчиков с сервера
        const response = await fetchWithToken(`http://localhost:8000/getSubscribers`, {
            method: 'GET',
        });

        if (!response.ok) {
            throw new Error('Ошибка при получении данных');
        }

        const data = await response.json();
        console.log("Список подписчиков:", data); // Выводим полный ответ для проверки

        if (data.Errors) {
            console.error('Ошибка при получении данных:', data.Errors);
            return;
        }
        if (data.Data === null || !Array.isArray(data.Data)) {
            console.error('Ожидался массив подписчиков, но получено:', data.Data);
            const friendsContainer = document.getElementById('friendsContainer');
            friendsContainer.innerHTML = `
                <p>Нет подписчиков</p>
            `;
            return;
        }
        // Проверьте, что data.Data действительно является массивом
        if (Array.isArray(data.Data)) {
            loadSubscriberData(data.Data, page);
            console.log("subs: ", data.Data);
        } else {
            console.error('Ожидался массив подписчиков, но получено:', data.Data);
            const friendsContainer = document.getElementById('friendsContainer');
            friendsContainer.innerHTML = '<p>Не удалось загрузить список подписчиков.</p>';
        }
    } catch (error) {
        console.error('Ошибка при загрузке подписчиков:', error);
    }
}

// Функция отображения подписчиков с учетом пагинации
function loadSubscriberData(subscribers, page) {
    const friendsContainer = document.getElementById('friendsContainer');
    friendsContainer.innerHTML = '';

    // Проверяем, что подписчики действительно являются массивом и имеют элементы
    if (!Array.isArray(subscribers) || subscribers.length === 0) {
        console.error('Подписчиков нет:', subscribers);
        friendsContainer.innerHTML = '<p>Подписчиков нет</p>';
        return;
    }

    // Отображаем подписчиков на странице
    subscribers.forEach(user => {
        const userElement = document.createElement('div');
        userElement.classList.add('user-item');
        userElement.innerHTML = `
            <div class="profile-avatar" style="border-radius: 50%; overflow: hidden; width: 50px; height: 50px;">
                <img src="${user.avatar}" alt="${user.username}" style="width: 100%; height: 100%; object-fit: cover;" />
            </div>
            <div class="user-info">
                <a href="user_profile.html?userID=${user.id}" class="user-username" style="color: blue; text-decoration: underline;">${user.username}</a>
                <button class="acceptSubscription" data-user-id="${user.id}">Принять запрос</button>
            </div>
        `;
        friendsContainer.appendChild(userElement);
    });

    // Добавляем кнопки пагинации если количество подписчиков больше ITEMS_PER_PAGE
    if (subscribers.length > ITEMS_PER_PAGE) {
        const pagination = document.createElement('div');
        pagination.classList.add('pagination');
        pagination.innerHTML = `
            <button id="prevSubscriberPage" ${page === 1 ? 'disabled' : ''}>Предыдущая</button>
            <span>Страница ${page}</span>
            <button id="nextSubscriberPage">Следующая</button>
        `;
        friendsContainer.appendChild(pagination);

        // Обработчики событий для кнопок пагинации
        document.getElementById('prevSubscriberPage').addEventListener('click', () => {
            if (page > 1) {
                loadSubscribers(page - 1);
            }
        });

        document.getElementById('nextSubscriberPage').addEventListener('click', () => {
            loadSubscribers(page + 1);
        });
    }

    // Обработчики для принятия запросов на подписку
    friendsContainer.querySelectorAll('.acceptSubscription').forEach(button => {
        button.addEventListener('click', async (event) => {
            const userId = event.target.getAttribute('data-user-id');
            if (isNaN(userId)) {
                console.error('Ошибка: неверный userId');
                return;
            }
            const data = {
                friend_id: parseInt(userId, 10)
            };
            try {
                const response = await fetchWithToken(`http://localhost:8000/acceptFriendShip?requestID=${userId}`, {
                    method: 'POST',
                    headers: { 
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(data)
                });
                const result = await response.json();
                if (result.Errors) {
                    console.error('Ошибка при принятии запроса:', result.Errors);
                } else {
                    alert('Запрос принят!');
                    // Перезагружаем списки друзей и подписчиков
                    loadFriends(currentPage);
                    loadSubscribers(currentSubscriberPage);
                }
            } catch (error) {
                console.error('Ошибка при принятии запроса:', error);
            }
        });
    });
}

// Обработчик для кнопки "Подписчики"
document.getElementById('friendsContainer').addEventListener('click', (event) => {
    if (event.target.id === 'subscribersButton') {
        loadSubscribers();
    }
});
