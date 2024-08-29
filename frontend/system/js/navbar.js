import { switchToContainer } from './main.js';
import { fetchUserProfile } from './profile.js';

document.addEventListener('DOMContentLoaded', () => {
    const chatsListButton = document.getElementById('chatsList');
    const profileButton = document.getElementById('profile');

    if (chatsListButton) {
        chatsListButton.addEventListener('click', () => {
            closeWebSocket();
            window.history.pushState({}, '', './main.html'); // Обновляем URL, но без перезагрузки
            switchToContainer('chats');
            fetchChats();
        });
    }

    if (profileButton) {
        profileButton.addEventListener('click', () => {
            closeWebSocket();
            window.history.pushState({}, '', './main.html'); // Обновляем URL, но без перезагрузки
            switchToContainer('profile');
            fetchUserProfile();
        });
    }
});

function closeWebSocket() {
    const chatContainer = document.getElementById('chatContainer');
    if (chatContainer && chatContainer._socket) {
        chatContainer._socket.close();
        chatContainer._socket = null;
    }
}
