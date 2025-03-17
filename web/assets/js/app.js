function closeLoading() {
    document.getElementById('loading').classList.add('hidden');
}

const axiosInstance = axios.create({
    baseURL: '/', // Your API base URL
    timeout: 10000,
    headers: {
      'Content-Type': 'application/json',
    },
  });

axiosInstance.interceptors.response.use(
  response => {
    return response;
  },
  async error => {
    const { response } = error;

    if (response && response.status === 403 && response.data.code === 'INVALID_CSRF_TOKEN') {
      alert('操作逾時，請重新操作。');

      try {
        await getCSRFToken();
      } catch (csrfError) {
        console.error('Failed to get CSRF token:', csrfError);
      }
    }

    return Promise.reject(error);
  }
);

function getCSRFToken() {
    return axiosInstance.get('/csrf-token')
        .then(response => {
            csrfToken = response.headers['x-csrf-token'];

            sessionStorage.setItem('csrfToken', csrfToken);

            axiosInstance.defaults.headers.common['X-Csrf-Token'] = csrfToken;
        })
        .catch(error => {
            console.error("Error fetching CSRF token:", error);
            throw error;
        });
}

function register(userData) {
    return axiosInstance.post('/register', userData)
        .then(response => response.data)
        .catch(error => {
            console.error("Error registering user:", error);
            throw error;
        });
}

function loginWithEmail(email, password) {
    return axiosInstance.post('/login', { email, password })
        .then(response => response.data)
        .catch(error => {
            console.error("Error logging in:", error);
            throw error;
        });
}

function logout() {
    axiosInstance.post('/logout')
        .then(() => window.location.href = '/template/login')
        .catch(error => {
            console.error("Error logging out:", error);
        });
}

async function getUser() {
    return axiosInstance.get('/user');
}

function getSocialAuthUrl(provider) {
    return axiosInstance.get('/login/social/${provider}')
        .then(response => response.data)
        .catch(error => {
            console.error(`Error getting social auth URL for ${provider}:`, error);
            throw error;
        });
}

function handleSocialAuthCallback(callbackData) {
    return axiosInstance.post('/login/social/callback', callbackData)
        .then(response => response.data)
        .catch(error => {
            console.error("Error handling social auth callback:", error);
            throw error;
        });
}

function getSocialAuthUrlForLinkingExistingUser(provider) {
    return axiosInstance.get('/auth/social/${provider}/link/url')
        .then(response => response.data)
        .catch(error => {
            console.error(`Error getting social auth URL for linking existing user with ${provider}:`, error);
            throw error;
        });
}

function linkSocialAccount(socialAccountData) {
    return axiosInstance.post('/auth/social/link', socialAccountData)
        .then(response => response.data)
        .catch(error => {
            console.error("Error linking social account:", error);
            throw error;
        });
}

function unlinkSocialAccount(provider) {
    return axiosInstance.delete('/user/unlink/${provider}')
        .then(response => response.data)
        .catch(error => {
            console.error(`Error unlinking social account with ${provider}:`, error);
            throw error;
        });
}
