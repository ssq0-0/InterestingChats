export async function logout() {
    const email = localStorage.getItem('email')
    try {
        const response = await fetch(`http://localhost:8000/deleteTokens?email=${email}`, {
            method: 'DELETE'
        });
    } catch(error) {
        console.log(error)
    }

    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    // TODO: добавьте запрос на удаление токенов с сервера

    window.location.href = '../Auth/auth.html';
}