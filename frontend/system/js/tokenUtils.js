export async function fetchWithToken(url, options = {}) {
    let accessToken = localStorage.getItem('access_token');
    const refreshToken = localStorage.getItem('refresh_token');

    if (!accessToken || !refreshToken) {
        // Если токенов нет, перенаправляем пользователя на страницу входа
        window.location.href = '/frontend/auth/auth.html';
        return;
    }

    options.headers = {
        ...options.headers,
        'Authorization': `Bearer ${accessToken}`,
        'Content-Type': 'application/json',
    };
    console.log('Access Token:', accessToken);

    let response = await fetch(url, options);

    // Если аксес токен истек, обновляем его
    if (response.status === 401) {
        console.log("токен истекк, рефрешим")
        const newAccessToken = await refreshAccessToken(refreshToken);
        console.log("рефрешим")

        if (newAccessToken) {
            // Обновляем аксес токен в localStorage
            localStorage.setItem('access_token', newAccessToken);

            // Повторяем запрос с новым аксес токеном
            options.headers['Authorization'] = `Bearer ${newAccessToken}`;
            response = await fetch(url, options);
        } else {
            // Если обновление не удалось, перенаправляем на страницу входа
            window.location.href = '/frontend/auth/auth.html';
            return;
        }
    }

    return response;
}

async function refreshAccessToken(refreshToken) {
    try {
        console.log("do refresh")
        const response = await fetch(`http://localhost:8000/refreshToken?refreshToken=${refreshToken}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
        });

        if (!response.ok) {
            throw new Error('Ошибка обновления токена');
        }

        const { Data, Errors } = await response.json();
        console.log('Refreshed Token:', Data);
        console.log("Errors", Errors)

        return Data; // Предполагается, что сервер возвращает новый аксес токен в `access_token`
    } catch (error) {
        console.error('Ошибка обновления токена:', error);
        return null;
    }
}
