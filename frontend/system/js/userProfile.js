import { fetchWithToken } from './tokenUtils.js';

async function GetUserProfile() {
    const backButton = document.getElementById('backButton');
    if (backButton) {
        backButton.addEventListener('click', () => {
            window.history.back();
        });
    }

    const urlParams = new URLSearchParams(window.location.search);
    const userID = urlParams.get('userID');
    const currentUserId = localStorage.getItem('id'); // Получаем текущий ID пользователя из localStorage
    const friends = JSON.parse(localStorage.getItem('friends') || '[]'); // Достаем список друзей

    try {
        // Приводим currentUserId и userID к числу для корректного сравнения
        const currentUserIdNumber = parseInt(currentUserId, 10);
        const userIdNumber = parseInt(userID, 10);

        // Проверка, если профиль свой — перенаправляем на основную страницу
        if (currentUserIdNumber && currentUserIdNumber === userIdNumber) {
            console.log('Перенаправление на main.html'); // Проверка перед редиректом
            window.location.href = './main.html'; // Убедитесь, что путь к main.html правильный
            return;
        }

        const response = await fetchWithToken(`http://localhost:8000/user_profile?userID=${userID}`, {
            method: 'GET',
            headers: { 
                'Content-Type': 'application/json'
            }
        });

        if (!response.ok) {
            throw new Error('Network response was not ok');
        }

        const { Data, Errors } = await response.json();

        if (Errors && Errors.length > 0) {
            document.getElementById('userProfile').innerHTML = `
                <p>Error: ${Errors.join(', ') || 'Something went wrong'}</p>
            `;
            return;
        }

        // Проверка, является ли пользователь другом
        const isFriend = friends.some(friend => friend.id === userIdNumber);

        // Отображаем информацию профиля и назначаем обработчики после рендеринга
        document.getElementById('userProfile').innerHTML = renderUserProfileInfo(
            Data.email,
            Data.username,
            Data.avatar,
            isFriend
        );

        // Назначаем обработчик события для кнопки добавления в друзья
        const addFriendButton = document.getElementById('addFriend');
        if (addFriendButton) {
            addFriendButton.addEventListener('click', handleAddFriend);
        }

        // Назначаем обработчик события для кнопки удаления из друзей, если она есть
        const removeFriendButton = document.getElementById('removeFriend');
        if (removeFriendButton) {
            removeFriendButton.addEventListener('click', handleRemoveFriend);
        }

    } catch (error) {
        console.error('Ошибка получения профиля пользователя:', error);
        document.getElementById('userProfile').innerHTML = `
            <p>Не удалось загрузить профиль пользователя. Попробуйте позже.</p>
        `;
    }
}

function renderUserProfileInfo(email, username, avatarUrl, isFriend) {
    const avatarHtml = avatarUrl ? `<img src="${avatarUrl}" alt="${username}'s avatar" class="profile-avatar">` : '';
    const friendButton = isFriend 
        ? `<button id="removeFriend">Удалить из друзей</button>`
        : `<button id="addFriend">Добавить в друзья</button>`;

    return `
        <h1>Профиль пользователя</h1>
        <div class="profile-info">
            ${avatarHtml}
            <p id="email">Email: ${email}</p>
            <p id="username">Username: ${username}</p>
            ${friendButton}
        </div>
    `;
}

async function handleAddFriend() {
    try {
        const urlParams = new URLSearchParams(window.location.search);
        const userID = urlParams.get('userID');  
        const data = {
            friend_id: parseInt(userID, 10)
        };
        const response = await fetchWithToken(`http://localhost:8000/requestToFriendShip?receiverID=${userID}`, {
            method: 'POST',
            headers: { 
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        });

        const { Data, Errors } = await response.json();
        if (Errors && Errors.length > 0) {
            console.error('Ошибка добавления в друзья:', Errors.join(', '));
        } else {
            console.log('Успешно добавлено в друзья:', Data);
            // Можно обновить интерфейс здесь, если нужно
        }
    } catch (error) {
        console.error('Ошибка при добавлении в друзья:', error);
    }
}

async function handleRemoveFriend() {
    try {
        const urlParams = new URLSearchParams(window.location.search);
        const userID = urlParams.get('userID');
        const userIdNumber = parseInt(userID, 10);

        // Получаем текущие списки друзей и подписчиков из localStorage
        let friends = JSON.parse(localStorage.getItem('friends') || '[]');
        let subscribers = JSON.parse(localStorage.getItem('subscribers') || '[]');

        const data = {
            friend_id:userIdNumber
        }
        // Отправляем запрос на удаление друга
        const response = await fetchWithToken(`http://localhost:8000/deleteFriend?removing=${userID}`, {
            method: 'DELETE',
            body: JSON.stringify(data)
        });

        const { Data, Errors } = await response.json();
        if (Errors && Errors.length > 0) {
            console.error('Ошибка удаления из друзей:', Errors.join(', '));
            return;
        }

        // console.log('Успешно удален из друзей:', Data);
        console.log('Друзья до удаления:', friends);
        console.log('Удаляем друга с ID:', userIdNumber);

        // Находим удаленного друга
        const removedFriend = friends.find(friend => friend.id === userIdNumber);
        // Удаляем пользователя из списка друзей
        const updatedFriends = friends.filter(friend => friend.id !== userIdNumber);

        // Добавляем пользователя в список подписчиков, если его там нет
        if (removedFriend) {
            if (!Array.isArray(subscribers)) {
                subscribers = [];
            }

            if (!subscribers.some(subscriber => subscriber.id === userIdNumber)) {
                subscribers.push(removedFriend);
            }
        }

        // Обновляем localStorage
        localStorage.setItem('friends', JSON.stringify(updatedFriends));
        localStorage.setItem('subscribers', JSON.stringify(subscribers));

        // Перезагружаем списки после обновления
        // loadUserData(currentPage); // Если находитесь на странице друзей
        // loadSubscriberData(currentSubscriberPage); // Если находитесь на странице подписчиков

        console.log('Друзья после удаления:', updatedFriends);
        console.log('Подписчики до добавления:', subscribers);

    } catch (error) {
        console.error('Ошибка при удалении из друзей:', error);
    }
}

// Добавляем обработчик события, чтобы отслеживать изменения в localStorage
window.addEventListener('storage', (event) => {
    if (event.key === 'friends' || event.key === 'subscribers') {
        // Перезагружаем данные друзей и подписчиков при изменении localStorage
        loadUserData(currentPage); // Перезагружает список друзей
        loadSubscriberData(currentSubscriberPage); // Перезагружает список подписчиков
    }
});

document.addEventListener('DOMContentLoaded', GetUserProfile);
