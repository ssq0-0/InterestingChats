import { fetchWithToken } from './tokenUtils.js'; // Импортируем функцию fetchWithToken

export async function fetchMyChats() {
    try {
        const response = await fetchWithToken("http://localhost:8000/getUserChats", {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            },
        });

        const { Errors, Data } = await response.json();
        if (!response.ok || (Errors && Errors.length > 0)) {
            document.getElementById('errorContainer').innerHTML = `
                <p>Error: ${Errors.join(', ') || 'Something went wrong'}</p>
            `;
            throw new Error('Network response was not ok');
        }
        console.log("Received chats:", Data);
        renderMyChats(Data);
    } catch (error) {
        console.error('Error:', error);
        document.getElementById('errorContainer').innerHTML = `
            <p>Failed to load chats. Please try again later.</p>
        `;
    }
}

export function renderMyChats(chats) {
    const chatContainer = document.getElementById('myChatsContainer');
    if (!chatContainer) {
        console.error('Error: Element with id "myChatsContainer" not found.');
        return;
    }

    if (!chats || chats.length === 0) {
        chatContainer.innerHTML = `<p>No chats found.</p>`;
        return;
    }

    chatContainer.innerHTML = `
        <h1>Your Chats</h1>
        <ul>
            ${chats.map(chat => `
                <li>
                    <span>${chat.chat_name}</span>
                    <button class="show-chat-btn" data-chat-id="${chat.id}">Показать чат</button>
                </li>
            `).join('')}
        </ul>
    `;
}
