import { fetchWithToken } from './tokenUtils.js';

export async function fetchUserProfile() {
    try {
        // Используем fetchWithToken вместо fetch
        const response = await fetchWithToken('http://localhost:8000/my_profile', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            }
        });

        // Проверяем, что ответ успешен
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        
        const { Data: { email, username }, Errors } = await response.json();

        // Обработка ошибок от сервера
        if (Errors && Errors.length > 0) {
            document.getElementById('feedback').innerHTML = `
                <p>Error: ${Errors.join(', ') || 'Something went wrong'}</p>
            `;
            return;
        }

        // Обновляем UI с полученными данными
        document.getElementById('userInfoContainer').innerHTML = renderProfileInfo(email, username);
        addEventListeners();
    } catch (error) {
        console.error('Error:', error);
    }
}

export function renderProfileInfo(email, username) {
    return `
        <h1>User Profile</h1>
        <p id="email">Email: ${email}</p>
        <p id="username">Username: ${username}</p>
        <button type="submit" id="changeUsername">Change username</button>
        <button type="submit" id="changeEmail">Change email</button>
        <p id="feedback"></p>
    `;
}

async function addEventListeners() {
    const changeUsernameButton = document.getElementById('changeUsername');
    const changeEmailButton = document.getElementById('changeEmail');

    if (changeUsernameButton) {
        changeUsernameButton.addEventListener('click', function() {
            this.outerHTML = `
                <input type="text" id="newUsername" placeholder="Enter new username">
                <button type="submit" class="sendData" data-type="username">Send</button>
            `;
            addSendDataListener(); 
        });
    }
    
    if (changeEmailButton) {
        changeEmailButton.addEventListener('click', function() {
            this.outerHTML = `
                <input type="text" id="newEmail" placeholder="Enter new email">
                <button type="submit" class="sendData" data-type="email">Send</button>
            `;
            addSendDataListener(); 
        });
    }
}

function addSendDataListener() {
    const sendButtons = document.querySelectorAll('.sendData');

    sendButtons.forEach(button => {
        button.addEventListener('click', async function() {
            const dataType = this.dataset.type;
            const newValue = dataType === 'username'
                ? document.getElementById('newUsername').value
                : document.getElementById('newEmail').value;

            const userID = localStorage.getItem('id');

            try {
                const response = await fetchWithToken('http://localhost:8000/changeData', {
                    method: 'POST',
                    headers: { 
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        user_id: parseInt(userID),
                        Data: newValue,
                        type: dataType 
                    })
                });

                const { Data, Errors } = await response.json();

                if (!response.ok || (Errors && Errors.length > 0)) {
                    document.getElementById('feedback').innerHTML = `
                        <p>Error: ${Errors.join(', ') || 'Something went wrong'}</p>
                    `;
                    throw new Error('Network response was not ok');
                } else {
                    if (dataType === 'username') {
                        document.getElementById('username').innerText = `Username: ${newValue}`;
                        localStorage.setItem('username', newValue);
                    } else {
                        document.getElementById('email').innerText = `Email: ${newValue}`;
                        localStorage.setItem('email', newValue);
                    }
                    this.outerHTML = `<p>Updated successfully!</p>`;
                }
            } catch (error) {
                console.error('Error:', error);
                // this.outerHTML = `<p>Failed to update. Try again later.</p>`;
            }
        });
    });
}
