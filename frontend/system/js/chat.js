// document.addEventListener('DOMContentLoaded', async function(){
//     const token = localStorage.getItem('access_token');
//     const urlParams = new URLSearchParams(window.location.search);
//     const chatID = urlParams.get('chatID');  
//     try {
//         const response = await fetch(`http://localhost:8000/getChat?chatID=${chatID}&Authorization=${token}`, {
//             method:'GET',
//             headers:{
//                 'Content-Type':'application:json'
//             }
//         });

//         if (!response.ok) {
//             console.log(error);
//         }

//         const {Data, Errors } = await response.json();
//         console.log(Data);
//         showInfo(Data);
//         const chatContainer = document.getElementById('chatContainer');
//         chatContainer._socket
//     } catch(error) {
//         console.log(error);
//     }
// })

// function showInfo(data) {
//     const chatContainer = document.getElementById('chatContainer');
//     chatContainer.innerHTML =  `
//     <h1>${data.chat_name}</h1>
//     <h2>Members:</h2>
//     <ul>
//         ${Object.values(data.members).map(member => `
//             <li><a href="user_profile.html?userID=${member.id}">${member.username} (${member.email})</a></li>
//         `).join('')}
//     </ul>
//     <h2>Messages:</h2>
//     <ul id="messagesList">
//         ${data.messages && data.messages.length > 0 ? 
//             data.messages.map(message => `
//                 <li>
//                     <strong>User ${message.user.username}:</strong> ${message.body} <em>(${new Date(message.time).toLocaleString()})</em>
//                 </li>
//             `).join('') :
//             '<li>No messages in this chat.</li>'}
//     </ul>
//     <label for="message">Write message...</label>
//     <input type="text" id="messageInput" name="message">
//     <button type="submit" id="sendMsg">Send message</button>
// `;
// }















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
        const response = await fetch(`http://localhost:8000/getChat?chatID=${chatID}&Authorization=${token}`, {
            method: 'GET',
            headers: { 'Content-Type': 'application/json' }
        });
        const result = await response.json();
        return result.Data;
    } catch (error) {
        console.error('Error fetching chat data:', error);
    }
}

// Initialize WebSocket connection
function initializeWebSocket(chatID, token) {
    const wsUrl = `ws://localhost:8004/wsOpen?chatID=${chatID}&Authorization=${token}`;
    console.log('WebSocket URL:', wsUrl); // Логируем WebSocket URL
    return new WebSocket(wsUrl);
}

// Setup WebSocket event listeners
function setupSocketListeners(socket) {
    socket.addEventListener('error', error => {
        console.error('WebSocket error:', error);
    });

    socket.addEventListener('close', () => {
        console.log('WebSocket connection closed');
    });

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
}

function generateChatHTML(data) {
    return `
        <h1>${data.chat_name}</h1>
        <h2>Members:</h2>
        <ul>
            ${Object.values(data.members).map(member => `
                <li><a href="user_profile.html?userID=${member.id}">${member.username} (${member.email})</a></li>
            `).join('')}
        </ul>
        <h2>Messages:</h2>
        <ul id="messagesList">
            ${data.messages && data.messages.length > 0 ? 
                data.messages.map(message => `
                    <li>
                        <strong>User ${message.user.username}:</strong> ${message.body} <em>(${new Date(message.time).toLocaleString()})</em>
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
function handleSendMessage() {
    const messageBody = document.getElementById('messageInput').value.trim();
    const chatContainer = document.getElementById('chatContainer');
    const socket = chatContainer._socket;

    if (messageBody) {
        if (socket && socket.readyState === WebSocket.OPEN) {
            // Log WebSocket state and URL before sending the message
            console.log('WebSocket state before sending message:', socket.readyState);
            console.log('WebSocket URL:', socket.url || 'URL not available'); // Если URL не доступен, используем заглушку

            const messageData = {
                body: messageBody,
                user: {
                    id: parseInt(localStorage.getItem('id')), 
                    email: localStorage.getItem('email'), 
                    username: localStorage.getItem('username')
                },
                time: new Date().toISOString()
            };

            socket.send(JSON.stringify(messageData));
            document.getElementById('messageInput').value = '';
        } else {
            console.error('WebSocket is not ready. Current state:', socket ? socket.readyState : 'undefined');
        }
    } else {
        console.warn('Message input is empty');
    }
}

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
            <strong>User ${message.user.username}:</strong> ${message.body} <em>(${new Date(message.time).toLocaleString()})</em>
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
