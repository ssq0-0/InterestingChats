import { fetchWithToken } from './tokenUtils.js'; // Импортируем функцию fetchWithToken


let unreadNotifications = [];
let debounceTimeout = null;
const readNotifications = new Set(); 

export async function fetchNotifications() {
    try {
        const userID = localStorage.getItem('id');
        const response = await fetchWithToken(`http://localhost:8000/getNotification?userID=${userID}`);
        const data = await response.json();

        if (data.Errors) {
            console.error('Ошибка при получении уведомлений:', data.Errors);
            return;
        }

        const notificationsContainer = document.getElementById('notificationsContainer');
        notificationsContainer.innerHTML = '';

        if (!data.Data || data.Data.length === 0) {
            const noNotificationsElement = document.createElement('div');
            noNotificationsElement.classList.add('no-notifications');
            noNotificationsElement.textContent = 'Уведомлений нет';
            notificationsContainer.appendChild(noNotificationsElement);
            return;
        }

        data.Data.forEach(notification => {
            const notificationElement = document.createElement('div');
            notificationElement.classList.add('notification-item');
            notificationElement.dataset.id = notification.ID; 

            const messageParts = notification.Message.split(' ');
            const username = messageParts[1]; 
            const userLink = `<a href="user_profile.html?userID=${notification.SenderID}">${username}</a>`;

            const messageWithLink = notification.Message.replace(username, userLink);

            notificationElement.innerHTML = `
                <div class="notification-type">${notification.Type}</div>
                <div class="notification-message">${messageWithLink}</div>
                <div class="notification-time">${new Date(notification.Time).toLocaleString()}</div>
                <div class="notification-status">${notification.IsRead ? 'Прочитано' : 'Не прочитано'}</div>
            `;

            if (notification.IsRead) {
                readNotifications.add(notification.ID);
                notificationElement.querySelector('.notification-status').textContent = 'Прочитано';
            } else {
                notificationElement.addEventListener('click', () => {
                    handleNotificationClick(notification.ID, notificationElement);
                });
            }

            notificationsContainer.appendChild(notificationElement);
        });
    } catch (error) {
        console.error('Ошибка при получении уведомлений:', error);
    }
}

function handleNotificationClick(notificationID, element) {
    if (!readNotifications.has(notificationID)) {
        unreadNotifications.push(notificationID);
        readNotifications.add(notificationID); 

        element.querySelector('.notification-status').textContent = 'Прочитано';

        clearTimeout(debounceTimeout);
        debounceTimeout = setTimeout(() => {
            markAsRead(unreadNotifications);
            unreadNotifications = []; 
        }, 3000); 
    }
}

async function markAsRead(notificationIDs) {
    try {
        const response = await fetchWithToken(`http://localhost:8000/readNotification`, {
            method: 'PATCH',
            body: JSON.stringify(notificationIDs.map(id => ({ ID: id }))) 
        });

        if (!response.ok) {
            throw new Error('Не удалось пометить уведомления как прочитанные');
        }

        console.log(`Уведомления с ID ${notificationIDs.join(', ')} помечены как прочитанные`);
    } catch (error) {
        console.error('Ошибка при пометке уведомлений как прочитанных:', error);
    }
}
