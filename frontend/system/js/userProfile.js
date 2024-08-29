// Импорт функции из tokenUtils.js
import { fetchWithToken } from './tokenUtils.js';

async function GetUserProfile() {
    document.getElementById('backButton').addEventListener('click', () => {
        window.history.back();
    });

    const urlParams = new URLSearchParams(window.location.search);
    const userID = urlParams.get('userID');  

    
    try {
        console.log('userID:', userID); // Добавьте эту строку перед вызовом fetchWithToken

        // Использование fetchWithToken для выполнения запроса с обработкой токена
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

        document.getElementById('userProfile').innerHTML = renderUserProfileInfo(Data.email, Data.username);

    } catch(error){
        console.error('Ошибка получения профиля пользователя:', error);
        document.getElementById('userProfile').innerHTML = `
            <p>Не удалось загрузить профиль пользователя. Попробуйте позже.</p>
        `;
    }
}

function renderUserProfileInfo(email, username) {
    return `
        <h1>User Profile</h1>
        <p id="email">Email: ${email}</p>
        <p id="username">Username: ${username}</p>
    `;
}

// Ждем загрузки контента, после чего вызываем GetUserProfile
document.addEventListener('DOMContentLoaded', GetUserProfile);
