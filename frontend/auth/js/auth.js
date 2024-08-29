const showLoginButton = document.getElementById('showLogin');
const showRegistrationButton = document.getElementById('showRegistration');
const loginFormContainer = document.getElementById('loginFormContainer');
const registrationFormContainer = document.getElementById('registrationFormContainer');

showLoginButton.addEventListener('click', () => {
    loginFormContainer.classList.add('active');
    registrationFormContainer.classList.remove('active');
});

showRegistrationButton.addEventListener('click', () => {
    registrationFormContainer.classList.add('active');
    loginFormContainer.classList.remove('active');
});

document.getElementById('loginForm').addEventListener('submit', async function(e) {
    e.preventDefault();

    const email = document.getElementById('emailLogin').value;
    const pass = document.getElementById('passwordLogin').value;
    
    const data = {
        email: email,
        password: pass
    };

    const jsonData = JSON.stringify(data);

    console.log(jsonData);

    try {
        const response = await fetch('http://localhost:8000/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: jsonData
        });

        if (!response.ok) {
            throw new Error('Network response was not ok');
        }

        const result = await response.json();

        console.log(result); 

        // TODO: create error message on client 
        if (result.Errors) {
            throw new Error(result.Errors);
        }

        if (!result.Data || !result.Data.tokens) {
            throw new Error('Tokens are missing from the response');
        }

        const { accessToken, refreshToken } = result.Data.tokens;

        if (!accessToken || !refreshToken) {
            throw new Error('Tokens are missing from the response');
        }

        localStorage.setItem('access_token', accessToken);
        localStorage.setItem('refresh_token', refreshToken);
        localStorage.setItem('email', result.Data.user.email);
        localStorage.setItem('username', result.Data.user.username);
        localStorage.setItem('id', result.Data.user.id);

        console.log('Login successful');
        console.log('Access token:', accessToken);
        console.log('Refresh token:', refreshToken);

        window.location.href = '../system/main.html';
    } catch (error) {
        console.error('Error:', error);
    }
});

document.getElementById('registrationForm').addEventListener('submit', async function(e) {
    e.preventDefault();

    const email = document.getElementById('email').value;
    const username = document.getElementById('username').value;
    const pass = document.getElementById('password').value;
    
    const data = {
        email: email,
        username: username,
        password: pass
    };

    const jsonData = JSON.stringify(data);

    console.log(jsonData);

    try {
        const response = await fetch('http://localhost:8000/registration', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: jsonData
        });

        if (!response.ok) {
            throw new Error('Network response was not ok');
        }

        const result = await response.json();

        console.log(result); 

        // TODO: create error message on client 
        if (result.Errors) {
            throw new Error(result.Errors);
        }

        if (!result.Data || !result.Data.tokens) {
            throw new Error('Tokens are missing from the response');
        }

        const { accessToken, refreshToken } = result.Data.tokens;

        if (!accessToken || !refreshToken) {
            throw new Error('Tokens are missing from the response');
        }

        localStorage.setItem('access_token', accessToken);
        localStorage.setItem('refresh_token', refreshToken);
        
        console.log('Registration successful');
        console.log('Access token:', accessToken);
        console.log('Refresh token:', refreshToken);

        window.location.href = '../system/main.html';
    } catch (error) {
        console.error('Error:', error);
    }
});
