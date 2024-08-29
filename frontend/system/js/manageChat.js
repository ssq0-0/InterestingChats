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
    }
});

function displayMembers(members, chatID) {
    const participantsListContainer = document.getElementById('participantsList');
    const membersArray = Object.values(members); // Преобразуем объект в массив

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
}

async function handleDeleteMember(memberId, chatID) {
    console.log(`Member with ID ${memberId} will be deleted.`);

    const chatIdInt = parseInt(chatID, 10);
    const data = [{
        user_id: memberId,
        chat_id: chatIdInt
    }];

    try {
        const response = await fetchWithToken('http://localhost:8000/deleteMember', {
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
