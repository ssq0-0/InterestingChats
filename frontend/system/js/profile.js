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
        
        const { Data: { email, username, avatar }, Errors } = await response.json();

        // Обработка ошибок от сервера
        if (Errors && Errors.length > 0) {
            document.getElementById('feedback').innerHTML = `
                <p>Error: ${Errors.join(', ') || 'Something went wrong'}</p>
            `;
            return;
        }

        // Обновляем UI с полученными данными
        document.getElementById('userInfoContainer').innerHTML = renderProfileInfo(email, username, avatar);
        addEventListeners();
    } catch (error) {
        console.error('Error:', error);
    }
}

export function renderProfileInfo(email, username, avatar) {
    const avatarSection = avatar
        ? `
            <img src="${avatar}" alt="User Avatar" id="userAvatar" style="width: 150px; height: 150px; border-radius: 50%;"/>
            <button type="submit" id="changeAvatar">Change photo</button>
          `
        : `
            <button type="submit" id="uploadAvatar">Upload photo</button>
          `;

    return `
        <h1>User Profile</h1>
        ${avatarSection}
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
    const changeAvatarButton = document.getElementById('changeAvatar');
    const uploadAvatarButton = document.getElementById('uploadAvatar');


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

    if (changeAvatarButton) {
        changeAvatarButton.addEventListener('click', function() {
            handlePhotoUpload();
        });
    }

    if (uploadAvatarButton) {
        uploadAvatarButton.addEventListener('click', function() {
            handlePhotoUpload();
        });
    }
}

// TODO: change func to change photo too
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
                    method: 'PATCH',
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

async function handlePhotoUpload() {
    // Создаем элемент <input type="file"> для выбора файла
    const fileInput = document.createElement('input');
    fileInput.type = 'file';
    fileInput.accept = 'image/*';
    fileInput.name = 'image'; // Имя должно соответствовать тому, что ожидается на сервере
    fileInput.style.display = 'none'; // Скрываем элемент

    // Добавляем элемент на страницу
    document.body.appendChild(fileInput);

    fileInput.addEventListener('change', async (event) => {
        const file = event.target.files[0];
        if (file) {
            // Создаем объект FormData
            const formData = new FormData();
            formData.append('image', file);

            try {
                // Отправляем файл на сервер
                const response = await fetch('http://localhost:8000/saveImage', {
                    method: 'POST',
                    body: formData,
                    headers: {
                        'Authorization': `Bearer ${localStorage.getItem('access_token')}` // Получаем токен из localStorage
                    }
                });

                const result = await response.json();

                // Проверяем, есть ли ошибки в ответе
                if (result.Errors) {
                    throw new Error(result.Errors);
                }

                // Обновляем аватар на странице
                const newAvatarUrl = result.Data; // URL аватара

                // Обновляем URL изображения, добавляя уникальный параметр для обхода кеша
                const timestamp = new Date().getTime(); // Получаем текущий timestamp
                document.getElementById('userAvatar').src = `${newAvatarUrl}?t=${timestamp}`;
                document.getElementById('feedback').innerHTML = `<p>Avatar updated successfully!</p>`;

            } catch (error) {
                console.error('Error:', error);
                document.getElementById('feedback').innerHTML = `<p>Failed to update avatar. Try again later.</p>`;
            }
        }
    });

    fileInput.click(); // Открываем диалог выбора файла
}
