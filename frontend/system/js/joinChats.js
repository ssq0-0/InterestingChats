import { fetchWithToken } from './tokenUtils.js'; // Импортируем функцию fetchWithToken

export async function joinToChat(chatID) {
    try {
        const response = await fetchWithToken(`http://localhost:8000/joinToChat?chatID=${chatID}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            }
        });

        if (!response.ok) throw new Error('Network response was not ok');
        
        // Дополнительная обработка успешного присоединения к чату, если нужно
        console.log(`Successfully joined chat with ID ${chatID}`);
        
    } catch (error) {
        console.error('Error joining chat:', error);
        const err = document.getElementById('error');
        if (err) {
            err.innerHTML = `Failed to join chat: ${error.message}`;
        }
    }
}
