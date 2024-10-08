// main.js
import { fetchChats } from './chats.js';
import { initializeUserSearch } from './searchUsers.js';
import { joinToChat } from './joinChats.js';
import { fetchMyChats } from './my_chats.js';
import { fetchUserProfile } from './profile.js';
import { logout } from './exitSystem.js';
import {fetchNotifications} from './notification.js';
import {loadFriends} from './friends.js';
import './navbar.js';

document.addEventListener('DOMContentLoaded', () => {
    setupEventListeners();
    handleInitialView();
});

function setupEventListeners() {
    document.getElementById('chatShorts').addEventListener('click', () => {
        switchToContainer('chats');
        fetchChats();
    });

    document.getElementById('profile').addEventListener('click', () => {
        switchToContainer('profile');
        fetchUserProfile();
    });

    document.getElementById('userSearch').addEventListener('click', () =>{
        switchToContainer('userSearch');
        initializeUserSearch();

    })

    document.getElementById('exit').addEventListener('click', () => {
        logout();
    });

    document.getElementById('chatsContainer').addEventListener('click', event => {
        const button = event.target.closest('.join-chat-btn');
        if (button) {
            const chatID = button.getAttribute('data-chat-id');
            joinToChat(chatID);
            // window.location.href = `./chat.html?chatID=${chatID}`;
            // 
            // TODO: create in func redirect to chat
            // 
        }
    });

    document.getElementById('myChatsContainer').addEventListener('click', event => {
        const button = event.target.closest('.show-chat-btn');
        if (button) {
            const chatID = button.getAttribute('data-chat-id');
            window.location.href = `./chat.html?chatID=${chatID}`;
        }
    });

    document.getElementById('myChats').addEventListener('click', ()=>{
        switchToContainer('myChat');
        fetchMyChats();
    })

    document.getElementById('notifications').addEventListener('click', () => {
        switchToContainer('notifications');
        fetchNotifications();
    });

    document.getElementById('friends').addEventListener('click', ()=>{
        switchToContainer('friends');
        loadFriends(1);
    })
}

function handleInitialView() {
    const activeContainer = localStorage.getItem('activeContainer') || 'chats';
    switchToContainer(activeContainer);
    if (activeContainer === 'profile') {
        fetchUserProfile();
    } else {
        fetchChats();
    }
}

export function switchToContainer(activeContainer) {
    const containers = {
        chats: document.getElementById('chatsContainer'),
        profile: document.getElementById('userInfoContainer'),
        chat: document.getElementById('chatContainer'),
        myChat: document.getElementById('myChatsContainer'),
        userSearch: document.getElementById('userSearchContainer'),
        notifications: document.getElementById('notificationsContainer'),
        friends: document.getElementById('friendsContainer')
    };

    Object.keys(containers).forEach(key => {
        const container = containers[key];
        if (container) {
            if (key === activeContainer) {
                container.classList.add('active');
                container.style.display = 'block';
            } else {
                container.classList.remove('active');
                container.style.display = 'none';
            }
        }
    });
}

