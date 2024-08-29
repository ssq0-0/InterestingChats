import { fetchWithToken } from './tokenUtils.js'; // Импортируем функцию fetchWithToken

document.addEventListener('DOMContentLoaded', async function(){
    const createButton = document.getElementById('create');
    const backButton = document.getElementById('back');

    const chatNameElement = document.getElementById('chatname');
    const creator = Number(localStorage.getItem('id'));
    const username = localStorage.getItem('username');
    const email = localStorage.getItem('email');
    const membersElement = document.getElementById('members');

    // TODO: check token and go auth if user unauthorized
    createButton.addEventListener('click', async () => {
        const chatName = chatNameElement.value.trim();
        const members = membersElement.value.split(',')
            .map(id => Number(id.trim())) // Преобразование строк в числа
            .filter(id => !isNaN(id)); 

        if (!members.includes(creator)) {
            members.push(creator);
        }

        const data = {
            chat_name: chatName,
            creator: creator,
            messages: [
                {
                    body: `${username} created chat at ${new Date().toISOString()}`, // Более читабельное время
                    user: {
                        id: creator,
                        email: email,
                        username: username
                    },
                    time: new Date().toISOString()
                }
            ],
            members: members
        };

        try {
            const response = await fetchWithToken("http://localhost:8000/createChat", {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(data)
            });

            if (!response.ok) {
                console.error("Failed to create chat.");
                return;
            }

            const { Data, Errors } = await response.json();
            if (Errors) {
                console.error("Server Errors:", Errors);
                return;
            }

            // Переход к созданному чату
            window.location.href = `./chat.html?chatID=${Data.id}`;
        } catch(error) {
            console.error("Error during chat creation:", error);
        }
    });

    backButton.addEventListener('click', () => {
        window.history.back();
    });
});
