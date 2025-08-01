// App State
const app = {
    currentUser: null,
    currentPage: 'home',
    stripe: null,
    api: {
        baseURL: '/api',
        token: localStorage.getItem('token')
    }
};

// Initialize Stripe
if (window.Stripe) {
    app.stripe = Stripe('pk_test_your_stripe_publishable_key');
}

// API Helper
async function apiCall(endpoint, options = {}) {
    const headers = {
        'Content-Type': 'application/json',
        ...options.headers
    };

    if (app.api.token) {
        headers['Authorization'] = `Bearer ${app.api.token}`;
    }

    try {
        const response = await fetch(`${app.api.baseURL}${endpoint}`, {
            ...options,
            headers
        });

        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.error || 'API Error');
        }

        return data;
    } catch (error) {
        console.error('API Error:', error);
        throw error;
    }
}

// Authentication
async function login(email, password) {
    try {
        const data = await apiCall('/login', {
            method: 'POST',
            body: JSON.stringify({ email, password })
        });

        app.api.token = data.token;
        localStorage.setItem('token', data.token);
        app.currentUser = data.user;
        
        updateNavigation();
        closeModal();
        showNotification('Login successful!', 'success');
        
        // Redirect based on user type
        if (data.user.user_type === 'cast') {
            navigateTo('cast-dashboard');
        } else {
            navigateTo('search');
        }
    } catch (error) {
        showNotification(error.message, 'error');
    }
}

async function register(userData) {
    try {
        const data = await apiCall('/register', {
            method: 'POST',
            body: JSON.stringify(userData)
        });

        app.api.token = data.token;
        localStorage.setItem('token', data.token);
        app.currentUser = data.user;
        
        updateNavigation();
        closeModal();
        showNotification('Registration successful!', 'success');
        
        if (userData.user_type === 'cast') {
            navigateTo('cast-profile-setup');
        } else {
            navigateTo('search');
        }
    } catch (error) {
        showNotification(error.message, 'error');
    }
}

function logout() {
    app.api.token = null;
    app.currentUser = null;
    localStorage.removeItem('token');
    updateNavigation();
    navigateTo('home');
    showNotification('Logged out successfully', 'success');
}

// Navigation
function navigateTo(page) {
    app.currentPage = page;
    renderPage();
}

function updateNavigation() {
    const navMenu = document.getElementById('navMenu');
    
    if (app.currentUser) {
        navMenu.innerHTML = `
            <a href="#" class="nav-link" data-page="home">Home</a>
            <a href="#" class="nav-link" data-page="search">Find Cast</a>
            ${app.currentUser.user_type === 'guest' ? 
                '<a href="#" class="nav-link" data-page="bookings">My Bookings</a>' :
                '<a href="#" class="nav-link" data-page="cast-dashboard">Dashboard</a>'
            }
            <a href="#" class="nav-link" data-page="profile">Profile</a>
            <a href="#" class="nav-link login-btn" onclick="logout()">Logout</a>
        `;
    } else {
        navMenu.innerHTML = `
            <a href="#" class="nav-link" data-page="home">Home</a>
            <a href="#" class="nav-link" data-page="search">Find Cast</a>
            <a href="#" class="nav-link" data-page="about">About</a>
            <a href="#" class="nav-link login-btn" data-page="login">Login</a>
            <a href="#" class="nav-link register-btn" data-page="register">Join</a>
        `;
    }
    
    attachNavListeners();
}

function attachNavListeners() {
    document.querySelectorAll('.nav-link[data-page]').forEach(link => {
        link.addEventListener('click', (e) => {
            e.preventDefault();
            const page = e.target.getAttribute('data-page');
            if (page === 'login' || page === 'register') {
                showAuthModal(page);
            } else {
                navigateTo(page);
            }
        });
    });
}

// Modal Functions
function showAuthModal(type) {
    const modal = document.getElementById('authModal');
    const content = document.getElementById('authContent');
    
    if (type === 'login') {
        content.innerHTML = `
            <div class="auth-form">
                <h2>Welcome Back</h2>
                <p class="auth-subtitle">Login to your account</p>
                <form id="loginForm">
                    <div class="form-group">
                        <label>Email</label>
                        <input type="email" name="email" required class="form-input">
                    </div>
                    <div class="form-group">
                        <label>Password</label>
                        <input type="password" name="password" required class="form-input">
                    </div>
                    <button type="submit" class="btn btn-primary btn-block">Login</button>
                </form>
                <p class="auth-switch">
                    Don't have an account? 
                    <a href="#" onclick="showAuthModal('register')">Register</a>
                </p>
            </div>
        `;
        
        document.getElementById('loginForm').addEventListener('submit', async (e) => {
            e.preventDefault();
            const formData = new FormData(e.target);
            await login(formData.get('email'), formData.get('password'));
        });
    } else {
        content.innerHTML = `
            <div class="auth-form">
                <h2>Join uso</h2>
                <p class="auth-subtitle">Create your premium account</p>
                <form id="registerForm">
                    <div class="form-group">
                        <label>Account Type</label>
                        <div class="radio-group">
                            <label class="radio-label">
                                <input type="radio" name="user_type" value="guest" checked>
                                <span>Guest - Book premium entertainment</span>
                            </label>
                            <label class="radio-label">
                                <input type="radio" name="user_type" value="cast">
                                <span>Cast - Provide entertainment services</span>
                            </label>
                        </div>
                    </div>
                    <div class="form-group">
                        <label>Name</label>
                        <input type="text" name="name" required class="form-input">
                    </div>
                    <div class="form-group">
                        <label>Email</label>
                        <input type="email" name="email" required class="form-input">
                    </div>
                    <div class="form-group">
                        <label>Password</label>
                        <input type="password" name="password" required minlength="8" class="form-input">
                    </div>
                    <div class="form-group">
                        <label>Birth Date</label>
                        <input type="date" name="birth_date" required class="form-input">
                    </div>
                    <div class="form-group">
                        <label>Phone (Optional)</label>
                        <input type="tel" name="phone" class="form-input">
                    </div>
                    <button type="submit" class="btn btn-primary btn-block">Create Account</button>
                </form>
                <p class="auth-switch">
                    Already have an account? 
                    <a href="#" onclick="showAuthModal('login')">Login</a>
                </p>
            </div>
        `;
        
        document.getElementById('registerForm').addEventListener('submit', async (e) => {
            e.preventDefault();
            const formData = new FormData(e.target);
            const userData = {
                user_type: formData.get('user_type'),
                name: formData.get('name'),
                email: formData.get('email'),
                password: formData.get('password'),
                birth_date: formData.get('birth_date'),
                phone: formData.get('phone') || null
            };
            await register(userData);
        });
    }
    
    modal.style.display = 'block';
}

function closeModal() {
    document.getElementById('authModal').style.display = 'none';
}

// Notification System
function showNotification(message, type = 'info') {
    const notification = document.createElement('div');
    notification.className = `notification notification-${type}`;
    notification.textContent = message;
    
    document.body.appendChild(notification);
    
    setTimeout(() => {
        notification.classList.add('show');
    }, 10);
    
    setTimeout(() => {
        notification.classList.remove('show');
        setTimeout(() => {
            notification.remove();
        }, 300);
    }, 3000);
}

// Page Rendering
function renderPage() {
    const mainContent = document.getElementById('mainContent');
    
    switch (app.currentPage) {
        case 'home':
            renderHomePage();
            break;
        case 'search':
            renderSearchPage();
            break;
        case 'cast-detail':
            renderCastDetailPage();
            break;
        case 'bookings':
            renderBookingsPage();
            break;
        case 'cast-dashboard':
            renderCastDashboard();
            break;
        case 'profile':
            renderProfilePage();
            break;
        default:
            renderHomePage();
    }
}

function renderHomePage() {
    // Home page is already in the HTML
    location.reload();
}

async function renderSearchPage() {
    const mainContent = document.getElementById('mainContent');
    mainContent.innerHTML = `
        <div class="search-page">
            <div class="search-hero">
                <h1>Find Your Perfect Cast</h1>
                <p>Discover premium entertainers for your special occasions</p>
            </div>
            
            <div class="search-container">
                <div class="search-filters">
                    <h3>Filters</h3>
                    <form id="searchForm">
                        <div class="form-group">
                            <label>Location</label>
                            <select name="location" class="form-input">
                                <option value="">All Locations</option>
                            </select>
                        </div>
                        <div class="form-group">
                            <label>Date</label>
                            <input type="date" name="date" class="form-input">
                        </div>
                        <div class="form-group">
                            <label>Time</label>
                            <input type="time" name="start_time" class="form-input">
                        </div>
                        <div class="form-group">
                            <label>Price Range</label>
                            <div class="price-range">
                                <input type="number" name="min_price" placeholder="Min" class="form-input">
                                <span>-</span>
                                <input type="number" name="max_price" placeholder="Max" class="form-input">
                            </div>
                        </div>
                        <div class="form-group">
                            <label>Rank</label>
                            <select name="rank" class="form-input">
                                <option value="">All Ranks</option>
                                <option value="standard">Standard ($60/hr)</option>
                                <option value="premium">Premium ($100/hr)</option>
                                <option value="vip">VIP ($150/hr)</option>
                            </select>
                        </div>
                        <button type="submit" class="btn btn-primary btn-block">Search</button>
                    </form>
                </div>
                
                <div class="search-results">
                    <div id="castGrid" class="cast-grid">
                        <!-- Cast cards will be loaded here -->
                    </div>
                </div>
            </div>
        </div>
    `;
    
    // Load service areas
    try {
        const areas = await apiCall('/service-areas');
        const select = document.querySelector('select[name="location"]');
        areas.forEach(area => {
            const option = document.createElement('option');
            option.value = area.name;
            option.textContent = area.name;
            select.appendChild(option);
        });
    } catch (error) {
        console.error('Failed to load service areas');
    }
    
    // Attach search form listener
    document.getElementById('searchForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        await searchCasts();
    });
    
    // Initial search
    await searchCasts();
}

async function searchCasts() {
    const form = document.getElementById('searchForm');
    const formData = new FormData(form);
    const params = new URLSearchParams();
    
    for (const [key, value] of formData.entries()) {
        if (value) params.append(key, value);
    }
    
    try {
        const data = await apiCall(`/casts/search?${params}`);
        renderCastGrid(data.casts);
    } catch (error) {
        showNotification('Failed to search casts', 'error');
    }
}

function renderCastGrid(casts) {
    const grid = document.getElementById('castGrid');
    
    if (casts.length === 0) {
        grid.innerHTML = '<div class="no-results">No casts found matching your criteria</div>';
        return;
    }
    
    grid.innerHTML = casts.map(cast => `
        <div class="cast-card" onclick="viewCastDetail(${cast.id})">
            <div class="cast-image">
                <img src="${cast.gallery_image || cast.profile_image || 'static/images/default-avatar.jpg'}" alt="${cast.name}">
                <div class="cast-rank ${cast.rank}">${cast.rank.toUpperCase()}</div>
            </div>
            <div class="cast-info">
                <h3>${cast.name}</h3>
                <div class="cast-meta">
                    <span class="cast-rate">$${cast.hourly_rate}/hr</span>
                    <span class="cast-rating">
                        ⭐ ${cast.rating.toFixed(1)} (${cast.review_count})
                    </span>
                </div>
                <p class="cast-bio">${cast.bio || 'Premium entertainment professional'}</p>
                <div class="cast-areas">
                    ${cast.service_areas.slice(0, 3).map(area => 
                        `<span class="area-tag">${area}</span>`
                    ).join('')}
                    ${cast.service_areas.length > 3 ? 
                        `<span class="area-tag">+${cast.service_areas.length - 3}</span>` : ''
                    }
                </div>
            </div>
        </div>
    `).join('');
}

function viewCastDetail(castId) {
    app.currentCastId = castId;
    navigateTo('cast-detail');
}

// Initialize App
document.addEventListener('DOMContentLoaded', () => {
    // Check if user is logged in
    if (app.api.token) {
        // Verify token and get user info
        apiCall('/profile').then(profile => {
            app.currentUser = profile;
            updateNavigation();
        }).catch(() => {
            // Token invalid
            logout();
        });
    } else {
        updateNavigation();
    }
    
    // Mobile menu toggle
    const navToggle = document.getElementById('navToggle');
    const navMenu = document.getElementById('navMenu');
    
    navToggle.addEventListener('click', () => {
        navMenu.classList.toggle('active');
    });
    
    // Modal close
    document.querySelector('.close').addEventListener('click', closeModal);
    
    window.addEventListener('click', (e) => {
        const modal = document.getElementById('authModal');
        if (e.target === modal) {
            closeModal();
        }
    });
    
    // Attach initial nav listeners
    attachNavListeners();
});

// Add styles for new components
const additionalStyles = `
<style>
/* Auth Forms */
.auth-form {
    padding: 48px;
}

.auth-form h2 {
    font-size: 36px;
    margin-bottom: 8px;
}

.auth-subtitle {
    color: var(--text-gray);
    margin-bottom: 32px;
}

.form-group {
    margin-bottom: 24px;
}

.form-group label {
    display: block;
    margin-bottom: 8px;
    font-weight: 500;
}

.form-input {
    width: 100%;
    padding: 16px;
    background: var(--glass-bg);
    border: 1px solid var(--glass-border);
    border-radius: 8px;
    color: var(--text-white);
    font-size: 16px;
    transition: all 0.3s ease;
}

.form-input:focus {
    outline: none;
    border-color: var(--accent-gold);
    box-shadow: 0 0 0 3px rgba(212, 175, 55, 0.1);
}

.radio-group {
    display: flex;
    flex-direction: column;
    gap: 12px;
}

.radio-label {
    display: flex;
    align-items: center;
    padding: 16px;
    background: var(--glass-bg);
    border: 1px solid var(--glass-border);
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.3s ease;
}

.radio-label:hover {
    border-color: var(--accent-gold);
}

.radio-label input[type="radio"] {
    margin-right: 12px;
}

.btn-block {
    width: 100%;
}

.auth-switch {
    text-align: center;
    margin-top: 24px;
    color: var(--text-gray);
}

.auth-switch a {
    color: var(--accent-gold);
    text-decoration: none;
}

/* Notifications */
.notification {
    position: fixed;
    top: 100px;
    right: 32px;
    padding: 16px 24px;
    background: var(--glass-bg);
    backdrop-filter: blur(10px);
    border: 1px solid var(--glass-border);
    border-radius: 8px;
    color: var(--text-white);
    font-weight: 500;
    opacity: 0;
    transform: translateX(100%);
    transition: all 0.3s ease;
    z-index: 3000;
}

.notification.show {
    opacity: 1;
    transform: translateX(0);
}

.notification-success {
    border-color: var(--success-green);
    background: rgba(16, 185, 129, 0.1);
}

.notification-error {
    border-color: #EF4444;
    background: rgba(239, 68, 68, 0.1);
}

/* Search Page */
.search-page {
    padding-top: 120px;
    min-height: 100vh;
}

.search-hero {
    text-align: center;
    padding: 80px 32px;
    background: linear-gradient(135deg, var(--secondary-blue), var(--primary-black));
}

.search-hero h1 {
    font-size: 48px;
    margin-bottom: 16px;
}

.search-container {
    display: grid;
    grid-template-columns: 300px 1fr;
    gap: 48px;
    max-width: 1400px;
    margin: 0 auto;
    padding: 0 32px 80px;
}

.search-filters {
    background: var(--glass-bg);
    backdrop-filter: blur(10px);
    border: 1px solid var(--glass-border);
    border-radius: 16px;
    padding: 32px;
    height: fit-content;
    position: sticky;
    top: 120px;
}

.search-filters h3 {
    margin-bottom: 24px;
}

.price-range {
    display: grid;
    grid-template-columns: 1fr auto 1fr;
    gap: 12px;
    align-items: center;
}

.cast-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
    gap: 32px;
}

.cast-card {
    background: var(--glass-bg);
    backdrop-filter: blur(10px);
    border: 1px solid var(--glass-border);
    border-radius: 16px;
    overflow: hidden;
    cursor: pointer;
    transition: all 0.3s ease;
}

.cast-card:hover {
    transform: translateY(-5px);
    box-shadow: 0 20px 40px rgba(212, 175, 55, 0.2);
    border-color: var(--accent-gold);
}

.cast-image {
    position: relative;
    height: 300px;
    overflow: hidden;
}

.cast-image img {
    width: 100%;
    height: 100%;
    object-fit: cover;
}

.cast-rank {
    position: absolute;
    top: 16px;
    right: 16px;
    padding: 8px 16px;
    border-radius: 8px;
    font-weight: 600;
    font-size: 14px;
    backdrop-filter: blur(10px);
}

.cast-rank.standard {
    background: rgba(107, 114, 128, 0.8);
}

.cast-rank.premium {
    background: rgba(212, 175, 55, 0.8);
}

.cast-rank.vip {
    background: linear-gradient(135deg, var(--rose-gold-start), var(--rose-gold-end));
}

.cast-info {
    padding: 24px;
}

.cast-info h3 {
    font-size: 24px;
    margin-bottom: 12px;
}

.cast-meta {
    display: flex;
    justify-content: space-between;
    margin-bottom: 16px;
    color: var(--text-gray);
}

.cast-rate {
    color: var(--accent-gold);
    font-weight: 600;
}

.cast-bio {
    color: var(--text-gray);
    margin-bottom: 16px;
    line-height: 1.6;
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
}

.cast-areas {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
}

.area-tag {
    padding: 4px 12px;
    background: var(--glass-bg);
    border: 1px solid var(--glass-border);
    border-radius: 16px;
    font-size: 14px;
    color: var(--text-gray);
}

.no-results {
    grid-column: 1 / -1;
    text-align: center;
    padding: 80px 32px;
    color: var(--text-gray);
    font-size: 20px;
}

@media (max-width: 1024px) {
    .search-container {
        grid-template-columns: 1fr;
    }
    
    .search-filters {
        position: relative;
        top: auto;
    }
}

@media (max-width: 768px) {
    .cast-grid {
        grid-template-columns: 1fr;
    }
}
</style>
`;

document.head.insertAdjacentHTML('beforeend', additionalStyles);