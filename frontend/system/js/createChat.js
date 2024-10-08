import { fetchWithToken } from './tokenUtils.js'; // Импортируем функцию fetchWithToken

document.addEventListener('DOMContentLoaded', async function() {
    const createButton = document.getElementById('create');
    const backButton = document.getElementById('back');
    const chatNameElement = document.getElementById('chatname');
    const membersContainer = document.getElementById('membersContainer');
    const dropdownButton = document.querySelector('.dropdown-btn');
    const hashtagsInput = document.getElementById('hashtags');

    const creator = Number(localStorage.getItem('id'));
    const username = localStorage.getItem('username');
    const email = localStorage.getItem('email');

    // Функция для загрузки друзей и заполнения выпадающего списка
    async function loadFriends() {
        try {
            const friends = JSON.parse(localStorage.getItem('friends') || '[]');
            membersContainer.innerHTML = '';

            friends.forEach(friend => {
                const label = document.createElement('label');
                const checkbox = document.createElement('input');
                checkbox.type = 'checkbox';
                checkbox.value = friend.id;
                checkbox.name = 'member';
                checkbox.className = 'member-checkbox';

                label.appendChild(checkbox);
                label.appendChild(document.createTextNode(friend.username));
                membersContainer.appendChild(label);
            });

        } catch (error) {
            console.error("Ошибка при загрузке друзей:", error);
        }
    }

    await loadFriends();

    dropdownButton.addEventListener('click', () => {
        membersContainer.style.display = membersContainer.style.display === 'block' ? 'none' : 'block';
    });

    createButton.addEventListener('click', async () => {
        const chatName = chatNameElement.value.trim();

        // Получаем выбранных пользователей
        const checkboxes = membersContainer.querySelectorAll('input[name="member"]:checked');
        const members = Array.from(checkboxes).map(checkbox => Number(checkbox.value));

        // Добавляем создателя чата в список участников, если его нет
        if (!members.includes(creator)) {
            members.push(creator);
        }

        // Получаем теги из поля и преобразуем их в массив
        const hashtagsArray = hashtagsInput.value.split(',').map(tag => tag.trim()).filter(tag => tag);
        const hashtags = hashtagsArray.map(tag => ({ hashtag: tag }));
        const data = {
            chat_name: chatName,
            creator: creator,
            members: members,
            hashtags: hashtags // Добавляем теги в данные запроса
        };
        console.log("data: ", data);
        try {
            const response = await fetchWithToken("http://localhost:8000/createChat", {
                method: 'POST',
                body: JSON.stringify(data)
            });

            if (!response.ok) {
                console.error("Не удалось создать чат.");
                return;
            }

            const { Data, Errors } = await response.json();
            if (Errors) {
                console.error("Ошибки на сервере:", Errors);
                return;
            }

            window.location.href = `./chat.html?chatID=${Data.id}`;
        } catch (error) {
            console.error("Ошибка при создании чата:", error);
        }
    });

    backButton.addEventListener('click', () => {
        window.history.back();
    });
});
