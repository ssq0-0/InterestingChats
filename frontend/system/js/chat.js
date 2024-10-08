
import { fetchWithToken } from './tokenUtils.js'; // Импортируем функцию fetchWithToken

document.addEventListener('DOMContentLoaded', () => {
    const urlParams = new URLSearchParams(window.location.search);
    const chatID = urlParams.get('chatID');  
    console.log('URL Chat ID:', chatID);  // Логируем chatID из URL

    if (chatID) {
        console.log("Calling showChatByID");
        showChatByID(chatID);
    } else {
        console.error('Chat ID not found in URL');
    }
});

async function showChatByID(chatID) {
    const token = localStorage.getItem('access_token');
    let socket;
    console.log("chat.js loaded");

    try {
        const chatData = await fetchChatData(chatID, token);
        console.log("data",chatData);
        displayChatData(chatData);

        socket = initializeWebSocket(chatID, token);
        setupSocketListeners(socket);

        const chatContainer = document.getElementById('chatContainer');
        chatContainer._socket = socket;

    } catch (error) {
        console.error('Error:', error);
    }

    return socket;
} 

// Fetch chat data from server
async function fetchChatData(chatID, token) {
    try {
        const response = await fetchWithToken(`http://localhost:8000/getChat?chatID=${chatID}&Authorization=${token}`, {
            method: 'GET',
            headers: { 'Content-Type': 'application/json' }
        });
        const result = await response.json();
        return result.Data;
    } catch (error) {
        console.error('Error fetching chat data:', error);
    }
}

// Обновленный метод для инициализации WebSocket
async function initializeWebSocket(chatID, token) {
    const userID = await authorizeUser(token);
    if (userID) {
        const wsUrl = `ws://localhost:8004/wsOpen?chatID=${chatID}&userID=${userID}`;
        const socket = new WebSocket(wsUrl);

        socket.onopen = function() {
            console.log('WebSocket connection established');
            // Теперь мы можем безопасно отправлять сообщения
            chatContainer._socket = socket; // Сохраняем сокет только после открытия
            setupSocketListeners(socket);
        };

        socket.onerror = function(event) {
            console.error('WebSocket error:', event);
        };

        return socket;
    } else {
        console.error('User authorization failed');
        return null;
    }
}

// Обновленный метод для отправки сообщения
function handleSendMessage() {
    const messageBody = document.getElementById('messageInput').value.trim();
    const chatContainer = document.getElementById('chatContainer');
    const socket = chatContainer._socket;
    const urlParams = new URLSearchParams(window.location.search);
    const chatID = urlParams.get('chatID');
    if (messageBody) {
        if (socket && socket.readyState === WebSocket.OPEN) {
            const messageData = {
                id: 1, // Возможно, стоит генерировать ID на сервере
                body: messageBody,
                chat_id: parseInt(chatID),
                user: {
                    id: parseInt(localStorage.getItem('id')),
                    email: localStorage.getItem('email'),
                    username: localStorage.getItem('username'),
                    avatar: localStorage.getItem('avatar') // добавляем аватар, если нужно
                },
                time: new Date().toISOString()
            };

            socket.send(JSON.stringify(messageData));
            document.getElementById('messageInput').value = ''; // Очищаем поле ввода
        } else {
            console.error('WebSocket is not ready. Current state:', socket ? socket.readyState : 'undefined');
        }
    } else {
        console.warn('Message input is empty');
    }
}


async function authorizeUser(token) {
    try {
        const response = await fetchWithToken(`http://localhost:8000/auth`, {
            method: 'POST',
            headers: { 'Authorization': token }
        });
        const result = await response.json();

        if (result.Data) {
            console.log('User authorized with ID:', result.Data);
            return result.Data;  // Вернем userID
        } else {
            console.error('Authorization failed:', result.Errors);
            return null;
        }
    } catch (error) {
        console.error('Error during authorization:', error);
        return null;
    }
}


// Setup WebSocket event listeners
function setupSocketListeners(socket) {
    // Обработчик ошибок
    socket.addEventListener('error', error => {
        console.error('WebSocket error occurred:', error);

        // Дополнительное логирование для отладки
        if (socket.readyState === WebSocket.CLOSED) {
            console.error('WebSocket connection closed unexpectedly');
        } else if (socket.readyState === WebSocket.CLOSING) {
            console.warn('WebSocket connection is closing');
        } else if (socket.readyState === WebSocket.CONNECTING) {
            console.warn('WebSocket connection is still in the process of connecting');
        }
    });

    // Обработчик закрытия соединения
    socket.addEventListener('close', event => {
        if (event.wasClean) {
            console.log(`WebSocket connection closed cleanly, code=${event.code}, reason=${event.reason}`);
        } else {
            console.error('WebSocket connection died unexpectedly');
        }
    });

    // Обработчик сообщений
    socket.addEventListener('message', event => {
        handleIncomingMessage(event.data);
    });
}

// Handle incoming WebSocket messages
function handleIncomingMessage(data) {
    try {
        const message = JSON.parse(data);
        console.log('Raw message received:', data);

        if (validateMessage(message)) {
            appendMessage(message);
        } else {
            console.error('Received message data is invalid:', message);
        }
    } catch (error) {
        console.error('Error parsing message data:', error);
    }
}

// Validate message structure
function validateMessage(message) {
    return message && message.user && message.user.id && message.user.username && message.user.email && message.body && message.time;
}

// Display chat data in the UI
function displayChatData(data) {
    console.log('Displaying chat data:', data);

    const chatContainer = document.getElementById('chatContainer');

    // Ensure that `messagesList` is created if needed
    let messageList = document.getElementById('messagesList');
    if (!messageList) {
        messageList = document.createElement('ul');
        messageList.id = 'messagesList';
        chatContainer.appendChild(messageList);
    }

    const chatHTML = generateChatHTML(data);
    chatContainer.innerHTML = chatHTML;

    // Call setupMsgButton after chatContainer is updated
    setupMsgButton();

    if (data.creator == localStorage.getItem('id')) {
        addManageChatButton(data.id, data.members);
    }

    const leaveChatButton = document.getElementById('leaveChatBtn');
    
    if (leaveChatButton) {
        leaveChatButton.addEventListener('click', async () => {
            const chatID = new URLSearchParams(window.location.search).get('chatID');
            const userID = localStorage.getItem('id');
            
            if (chatID && userID) {
                await leaveChat(chatID, userID);
            } else {
                console.error('Chat ID or User ID not found');
            }
        });
    }
}

function getInitials(username) {
    const names = username.split(' ');
    return names.map(name => name.charAt(0)).join('').toUpperCase();
}

function generateChatHTML(data) {
    const currentUserId = parseInt(localStorage.getItem('id')); // Получаем текущего пользователя из localStorage
    
    return `
        <h1>${data.chat_name}</h1>
        <h2>Members:</h2>
        <ul>
            ${Object.values(data.members).map(member => `
                <li>
                    <div class="user-avatar" title="${member.username}">
                        ${member.avatar ? `<img src="${member.avatar}" alt="${member.username}">` : getInitials(member.username)}
                    </div>
                    <a href="user_profile.html?userID=${member.id}">
                        ${member.username} (${member.email})
                    </a>
                </li>
            `).join('')}
        </ul>
        ${data.members[currentUserId] ? `<button id="leaveChatBtn">Выйти из чата</button>` : ''}
        <h2>Messages:</h2>
        <ul id="messagesList">
            ${data.messages && data.messages.length > 0 ? 
                data.messages.map(message => `
                    <li>
                        <div class="message-avatar" title="${message.user.username}">
                            ${message.user.avatar ? `<img src="${message.user.avatar}" alt="${message.user.username}">` : getInitials(message.user.username)}
                        </div>
                        <strong>${message.user.username}:</strong> ${message.body} <em>(${new Date(message.time).toLocaleString()})</em>
                    </li>
                `).join('') :
                '<li>No messages in this chat.</li>'}
        </ul>
        <label for="message">Write message...</label>
        <input type="text" id="messageInput" name="message">
        <button type="submit" id="sendMsg">Send message</button>
    `;
}




// Setup message button event listener
async function setupMsgButton() {
    const sendMsgButton = document.getElementById('sendMsg');
    const messageInput = document.getElementById('messageInput'); // Исправление: добавлено

    if (sendMsgButton) {
        sendMsgButton.addEventListener('click', () => {
            handleSendMessage();
        });
    } else {
        console.error('Send message button not found');
    }

    // Add event listener for Enter key press
    if (messageInput) {
        messageInput.addEventListener('keypress', (event) => {
            if (event.key === 'Enter') {
                event.preventDefault(); // Prevent form submission if inside a form
                handleSendMessage();
            }
        });
    } else {
        console.error('Message input field not found');
    }
}

// Handle sending a message
// function handleSendMessage() {
//     const messageBody = document.getElementById('messageInput').value.trim();
//     const chatContainer = document.getElementById('chatContainer');
//     const socket = chatContainer._socket;

//     const urlParams = new URLSearchParams(window.location.search);
//     const chatID = urlParams.get('chatID');

//     if (messageBody) {
//         if (socket && socket.readyState === WebSocket.OPEN) {
//             // Log WebSocket state and URL before sending the message
//             console.log('WebSocket state before sending message:', socket.readyState);
//             console.log('WebSocket URL:', socket.url || 'URL not available'); // Если URL не доступен, используем заглушку

//             const messageData = {
//                 body: messageBody,
//                 chat_id: parseInt(chatID),
//                 user: {
//                     id: parseInt(localStorage.getItem('id')), 
//                     email: localStorage.getItem('email'), 
//                     username: localStorage.getItem('username')
//                 },
//                 time: new Date().toISOString()
//             };

//             socket.send(JSON.stringify(messageData));
//             document.getElementById('messageInput').value = '';
//         } else {
//             console.error('WebSocket is not ready. Current state:', socket ? socket.readyState : 'undefined');
//         }
//     } else {
//         console.warn('Message input is empty');
//     }
// }

// Append a message to the chat
function appendMessage(message) {
    console.log('Appending message:', message);
    let messageList = document.getElementById('messagesList');

    if (!messageList) {
        messageList = document.createElement('ul');
        messageList.id = 'messagesList';
        const chatContainer = document.getElementById('chatContainer');
        chatContainer.appendChild(messageList);
    }

    const noMessagesNotice = Array.from(messageList.querySelectorAll('li')).find(li => li.textContent.includes('No messages in this chat.'));
    if (noMessagesNotice) {
        noMessagesNotice.remove();
    }

    const messageHTML = `
        <li>
            <img src="${message.user.avatar || 'default-avatar.png'}" alt="${message.user.username}" class="message-avatar">
            <strong>${message.user.username}:</strong> ${message.body} <em>(${new Date(message.time).toLocaleString()})</em>
        </li>
    `;
    console.log('Message to insert:', messageHTML);

    messageList.insertAdjacentHTML('beforeend', messageHTML);
}


document.getElementById('back').addEventListener('click', () => {
    window.location.href = 'main.html';
});

// Add manage chat button if the user is the creator
function addManageChatButton(chatID, members) {
    const chatContainer = document.getElementById('chatContainer');
    const manageButton = document.createElement('button');
    manageButton.textContent = 'Управление чатом';
    manageButton.id = 'manageChatBtn';
    chatContainer.appendChild(manageButton);

    // Setup event listener for the manage chat button
    manageButton.addEventListener('click', () => {
        console.log('Manage chat button clicked');
        localStorage.setItem(chatID, JSON.stringify(members))
        window.location.href = `./manageChat.html?chatID=${chatID}`;
    });
}

async function leaveChat(chatID, userID) {
    try {
        const response = await fetchWithToken('http://localhost:8000/leaveChat', {
            method: 'DELETE',
            body: JSON.stringify([{ chat_id: parseInt(chatID), user_id: parseInt(userID) }])
        });
        
        const result = await response.json();
        
        if (response.ok) {
            console.log('Successfully left the chat');
            // Перенаправляем пользователя на главную страницу после успешного выхода
            window.location.href = 'main.html';
        } else {
            console.error('Error leaving chat:', result.Errors);
        }
    } catch (error) {
        console.error('Error while leaving chat:', error);
    }
}
