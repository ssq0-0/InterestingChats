import { fetchWithToken } from './tokenUtils.js'; // Импортируем функцию fetchWithToken

document.addEventListener('DOMContentLoaded', () => {
    const urlParams = new URLSearchParams(window.location.search);
    const chatID = urlParams.get('chatID');

    if (chatID) {
        const members = JSON.parse(localStorage.getItem(chatID));
        console.log(members);
        if (members) {
            displayMembers(members, chatID);
        }
        addDeleteChatButton(chatID);
        fetchAndDisplayTags(chatID);
        setupChangeChatNameButton(chatID);
    }
});

function displayMembers(members, chatID) {
    const participantsListContainer = document.getElementById('participantsList');
    const membersArray = Object.values(members); // Преобразуем объект в массив
    const addMemberButton = document.getElementById('addMember');
    const tagsButton = document.getElementById('addTag');

    participantsListContainer.innerHTML = membersArray.map(member => 
        `<li>
            <a href="user_profile.html?userID=${member.id}">${member.username} (${member.email})</a>
            <button id="deleteMember-${member.id}">Удалить участника</button>
        </li>`
    ).join('');

    // Добавляем обработчики событий для каждой кнопки удаления
    membersArray.forEach(member => {
        const deleteButton = document.getElementById(`deleteMember-${member.id}`);
        deleteButton.addEventListener('click', () => {
            handleDeleteMember(member.id, chatID);
        });
    });

    // Добавляем обработчик для кнопки добавления участника
    addMemberButton.addEventListener('click', () => {
        initializeUserSearch(chatID);
    });

    tagsButton.addEventListener('click', () => {
        showTagInput(chatID);
    });
}

// Функция для получения и отображения тегов
async function fetchAndDisplayTags(chatID) {
    try {
        const response = await fetchWithToken(`http://localhost:8000/getTags?chatID=${chatID}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            },
        });

        if (!response.ok) throw new Error('Network response was not ok');

        const data = await response.json();

        // Проверяем, что в ответе есть поле Data и это массив
        if (data.Errors) {
            console.error('Errors fetching tags:', data.Errors);
            document.getElementById('error').textContent = `Ошибка при получении тегов: ${data.Errors.join(', ')}`;
        } else if (Array.isArray(data.Data)) {
            displayTags(data.Data); // Передаем массив тегов в displayTags
        } else {
            console.error('Unexpected data format:', data);
            document.getElementById('error').textContent = 'Неверный формат данных при получении тегов.';
        }
    } catch (error) {
        console.error('Error fetching tags:', error);
        document.getElementById('error').textContent = `Ошибка при получении тегов: ${error.message}`;
    }
}

// Функция для отображения тегов с возможностью удаления
function displayTags(tags) {
    const tagsContainer = document.getElementById('tags');

    if (!Array.isArray(tags)) {
        console.error('Tags is not an array:', tags);
        tagsContainer.innerHTML = '<p>Нет доступных тегов.</p>';
        return;
    }

    // Отображаем теги с кнопками удаления
    tagsContainer.innerHTML = tags.map(tag => `
        <li>
            ${tag.hashtag}
            <button class="delete-tag-btn" data-tag-id="${tag.id}">Удалить</button>
        </li>
    `).join('');

    // Добавляем обработчики для кнопок удаления
    document.querySelectorAll('.delete-tag-btn').forEach(button => {
        button.addEventListener('click', () => {
            const tagId = button.getAttribute('data-tag-id');
            handleDeleteTag(tagId);
        });
    });
}

// Функция для удаления хештега
async function handleDeleteTag(tagId) {
    const chatID = new URLSearchParams(window.location.search).get('chatID');
    if (!chatID || !tagId) {
        document.getElementById('error').textContent = 'Ошибка: не удалось получить идентификатор чата или тега.';
        return;
    }

    try {
        const response = await fetchWithToken(`http://localhost:8000/deleteTags?chatID=${chatID}`, {
            method: 'DELETE',
            body: JSON.stringify({
                chat_id: parseInt(chatID),
                options: [tagId] // Отправляем ID тега для удаления
            })
        });

        const data = await response.json();
        if (data.Errors) {
            document.getElementById('error').textContent = `Ошибка при удалении тега: ${data.Errors.join(', ')}`;
        } else {
            document.getElementById('error').textContent = 'Тег успешно удалён!';
            fetchAndDisplayTags(chatID); // Обновляем список тегов после удаления
        }
    } catch (error) {
        document.getElementById('error').textContent = `Ошибка: ${error.message}`;
    }
}


async function handleDeleteMember(memberId, chatID) {
    console.log(`Member with ID ${memberId} will be deleted.`);

    const data = [{
        user_id: memberId,
        chat_id: parseInt(chatID)
    }];

    try {
        const response = await fetchWithToken(`http://localhost:8000/deleteMember?chatID=${chatID}`, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        });

        if (!response.ok) throw new Error('Network response was not ok');

        const result = await response.json();
        const { Errors } = result;

        if (Errors && Errors.length > 0) {
            console.error('Errors:', Errors);
            const err = document.getElementById('error');
            err.innerHTML = `Failed to delete member: ${Errors.join(', ')}`;
        } else {
            console.log('Member deleted successfully');
            // Обновите список участников после успешного удаления
            const members = JSON.parse(localStorage.getItem(chatID));
            if (members) {
                // Удалите участника из localStorage и обновите отображение
                delete members[memberId];
                localStorage.setItem(chatID, JSON.stringify(members));
                displayMembers(members, chatID);
            }
        }
    } catch (error) {
        console.error('Error:', error);
        const err = document.getElementById('error');
        err.innerHTML = `Failed to delete member: ${error.message}`;
    }
}

function addDeleteChatButton(chatID) {
    const deleteChatButton = document.getElementById('deleteChat');
    
    if (deleteChatButton) { // Проверяем, существует ли элемент
        deleteChatButton.addEventListener('click', () => {
            handleDeleteChat(chatID);
        });
    } else {
        console.error('Элемент с id="deleteChat" не найден.');
    }
}

async function handleDeleteChat(chatID) {
    console.log(`Чат с ID ${chatID} будет удален.`);

    try {
        const response = await fetchWithToken(`http://localhost:8000/deleteChat?chatID=${chatID}`, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json'
            }
        });

        const isJson = response.headers.get('content-type')?.includes('application/json');
        if (isJson) {
            const { Data, Errors } = await response.json();
            console.error('Data:', Data);

            if (Errors && Errors.length > 0) {
                console.error('Errors:', Errors);
                const err = document.getElementById('error');
                err.innerHTML = `Failed to delete chat: ${Errors.join(', ')}`;
            } else {
                console.log('Chat deleted successfully');
                window.location.href = 'main.html'; // Переход на главную страницу после успешного удаления
            }
        } else {
            console.log('Empty or non-JSON response, redirecting to main page.');
            window.location.href = 'main.html'; // Перенаправление на главную страницу
        }
    } catch (error) {
        console.error('Error:', error);
        const err = document.getElementById('error');
        err.innerHTML = `Failed to delete chat: ${error.message}`;
    }
}

document.getElementById('back').addEventListener('click', () => {
    window.history.back();
});

async function initializeUserSearch(chatID) {
    const userSearchContainer = document.getElementById('userSearchContainer');
    userSearchContainer.innerHTML = `
        <h1>Search users</h1>
        <input type="text" id="searchInput" placeholder="Search users...">
        <ul id="searchResults"></ul>
    `;

    const searchInput = document.getElementById('searchInput');
    searchInput.addEventListener('input', debounce(() => {
        const query = searchInput.value.trim();
        if (query.length > 0) {
            searchUsers(query, chatID);
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

async function searchUsers(query, chatID) {
    const friends = JSON.parse(localStorage.getItem('friends')) || [];
    const friendsInChat = Object.values(JSON.parse(localStorage.getItem(chatID)) || {}).map(member => member.id);
    
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
            renderSearchResults([], friends, chatID); // Показать, что пользователи не найдены
        } else if (Array.isArray(result.Data)) {
            const nonChatFriends = friends.filter(friend => !friendsInChat.includes(friend.id));
            const recommendedFriends = nonChatFriends.filter(friend => 
                result.Data.some(user => user.id === friend.id)
            );
            const otherUsers = result.Data.filter(user => !nonChatFriends.some(friend => friend.id === user.id));
            renderSearchResults(recommendedFriends.concat(otherUsers), friends, chatID);
        } else {
            console.error('Invalid data format:', result);
            renderSearchResults([], friends, chatID); // Показать, что пользователи не найдены
        }
    } catch (error) {
        console.error('Error during user search:', error);
        renderSearchResults([], [], chatID); // Показать, что пользователи не найдены в случае ошибки
    }
}

function renderSearchResults(users, friends, chatID) {
    const searchResults = document.getElementById('searchResults');

    if (users && users.length > 0) {
        searchResults.innerHTML = users.map(user => `
            <li>
                ${user.username} (${user.email})
                <button class="add-to-chat-btn" data-user-id="${user.id}" data-user-username="${user.username}" data-user-email="${user.email}">Добавить в чат</button>
            </li>
        `).join('');
    } else {
        searchResults.innerHTML = '<li>Пользователи не найдены.</li>';
    }

    document.querySelectorAll('.add-to-chat-btn').forEach(button => {
        button.addEventListener('click', () => {
            const userId = button.getAttribute('data-user-id');
            const username = button.getAttribute('data-user-username');
            const user_email = button.getAttribute('data-user-email');
            addUserToChat(userId, username, user_email, chatID);
        });
    });
}

async function addUserToChat(userId, username, user_email, chatID) {
    try {
        // Отправляем запрос на сервер для добавления пользователя
        const response = await fetchWithToken(`http://localhost:8000/addMember?chatID=${chatID}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify([{
                user_id: parseInt(userId),
                chat_id: parseInt(chatID)
            }])
        });

        if (!response.ok) throw new Error('Network response was not ok');

        const result = await response.json();
        const { Errors } = result;

        if (Errors && Errors.length > 0) {
            console.error('Server Errors:', Errors);
            const errorElement = document.getElementById('error');
            if (errorElement) {
                errorElement.innerHTML = `Failed to add user: ${Errors.join(', ')}`;
            }
        } else {
            console.log('User added successfully');

            // Обновляем список участников в локальном хранилище
            const members = JSON.parse(localStorage.getItem(chatID)) || {};
            
            // Создаем объект пользователя для добавления в localStorage
            const newUser = {
                id: parseInt(userId),
                username: username,
                email: user_email,
                avatar: null // Убедитесь, что у вас есть информация о пользователе для заполнения других полей
            };

            members[userId] = newUser; // Обновляем или добавляем пользователя в список
            localStorage.setItem(chatID, JSON.stringify(members));
            
            // Логируем текущий список участников
            console.log('Updated members in localStorage:', JSON.parse(localStorage.getItem(chatID)));

            // Обновляем отображение списка участников
            displayMembers(members, chatID);
        }
    } catch (error) {
        console.error('Error:', error);
        const errorElement = document.getElementById('error');
        if (errorElement) {
            errorElement.innerHTML = `Failed to add user: ${error.message}`;
        }
    }
}

function showTagInput(chatID) {
    const container = document.getElementById('userSearchContainer');
    container.innerHTML = `
        <div id="tagForm">
            <input type="text" id="tagInput" placeholder="Введите теги, разделенные запятыми" />
            <button id="submitTags">Добавить теги</button>
        </div>
    `;

    document.getElementById('submitTags').addEventListener('click', () => {
        const tagInput = document.getElementById('tagInput').value;
        const tags = tagInput.split(',').map(tag => tag.trim()).filter(tag => tag);
        submitTags(tags, chatID);
    });
}

async function submitTags(tags, chatID) {
    console.log("tags: ", tags)
    try {
        const response = await fetchWithToken('http://localhost:8000/setTag', {
            method: 'PATCH',
            body: JSON.stringify({
                chat_id: parseInt(chatID),
                options: tags
            })
        });

        const data = await response.json();
        if (data.Errors) {
            document.getElementById('error').textContent = `Ошибка: ${data.Errors}`;
        } else {
            document.getElementById('error').textContent = `Теги успешно добавлены!`;
            // Очистка формы после успешной отправки
            document.getElementById('userSearchContainer').innerHTML = '';
        }
    } catch (error) {
        document.getElementById('error').textContent = `Ошибка: ${error.message}`;
    }
}

function setupChangeChatNameButton(chatID) {
    const changeChatNameButton = document.getElementById('changeChatName');

    changeChatNameButton.addEventListener('click', () => {
        // Проверяем, существует ли контейнер
        let chatNameContainer = document.getElementById('changeChatNameContainer');
        
        // Если контейнера нет, создаем его
        if (!chatNameContainer) {
            chatNameContainer = document.createElement('div');
            chatNameContainer.id = 'changeChatNameContainer';
            
            const input = document.createElement('input');
            input.type = 'text';
            input.id = 'newChatName';
            input.placeholder = 'Enter new chat name';

            const submitButton = document.createElement('button');
            submitButton.id = 'submitChatName';
            submitButton.textContent = 'Submit';

            chatNameContainer.appendChild(input);
            chatNameContainer.appendChild(submitButton);
            document.querySelector('body').appendChild(chatNameContainer);
        }

        // Показываем контейнер для ввода нового названия
        chatNameContainer.style.display = 'block';

        // Обработчик для кнопки отправки
        document.getElementById('submitChatName').addEventListener('click', async () => {
            const newChatName = document.getElementById('newChatName').value.trim();
            if (newChatName) {
                try {
                    const response = await fetchWithToken(`http://localhost:8000/changeChatName?chatID=${chatID}&chatName=${newChatName}`, {
                        method: 'PATCH'
                    });

                    const data = await response.json();
                    if (data.Errors) {
                        document.getElementById('error').textContent = `Error: ${data.Errors.join(', ')}`;
                    } else {
                        document.getElementById('error').textContent = 'Chat name changed successfully!';
                        // Обновляем заголовок или имя чата на странице
                        document.querySelector('h1').textContent = newChatName;
                        // Скрываем форму после успешного изменения
                        document.getElementById('changeChatNameContainer').style.display = 'none';
                    }
                } catch (error) {
                    document.getElementById('error').textContent = `Error: ${error.message}`;
                }
            } else {
                document.getElementById('error').textContent = 'Chat name cannot be empty.';
            }
        });
    });
}