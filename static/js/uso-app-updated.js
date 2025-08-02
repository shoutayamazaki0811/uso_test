// Uso App - Complete Frontend Application with Swipe Feature

// Global state
const app = {
    currentUser: null,
    currentPage: 'home',
    token: localStorage.getItem('uso_token'),
    apiBase: '/api',
    casts: [],
    messages: [],
    bookings: [],
    favorites: new Set()
};

// API Helper Functions
async function apiCall(endpoint, options = {}) {
    const headers = {
        'Content-Type': 'application/json',
        ...options.headers
    };

    if (app.token) {
        headers['Authorization'] = `Bearer ${app.token}`;
    }

    try {
        const response = await fetch(`${app.apiBase}${endpoint}`, {
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
        showNotification(error.message, 'error');
        throw error;
    }
}

// Mock Data with extended properties for swipe feature
const mockCasts = [
    {
        id: 1,
        name: 'ゆか',
        age: 25,
        image: window.IMAGE_BASE_PATH ? `${window.IMAGE_BASE_PATH}/casts/yuka.png` : 'static/images/casts/yuka.png',
        rating: 4.9,
        reviewCount: 47,
        hourlyRate: 15000,
        rank: 'premium',
        bio: '初めまして、ゆかです。楽しい時間を一緒に過ごしましょう！',
        introduction: 'こんにちは！ゆかです。楽しい時間を過ごしましょう！',
        serviceAreas: ['六本木', '渋谷', '恵比寿'],
        languages: ['日本語', 'English'],
        tags: ['モデル経験あり', 'ワイン好き', '会話上手'],
        interests: ['映画鑑賞', '料理', '旅行'],
        zodiacSign: '乙女座',
        bloodType: 'A型',
        occupation: 'モデル',
        personality: '明るい',
        verified: true
    },
    {
        id: 2,
        name: 'かな',
        age: 27,
        image: window.IMAGE_BASE_PATH ? `${window.IMAGE_BASE_PATH}/casts/kana.png` : 'static/images/casts/kana.png',
        rating: 5.0,
        reviewCount: 32,
        hourlyRate: 20000,
        rank: 'vip',
        bio: '素敵な思い出作りのお手伝いをさせていただきます。',
        introduction: 'かなです♪ 一緒に素敵な時間を過ごしませんか？',
        serviceAreas: ['銀座', '丸の内', '六本木'],
        languages: ['日本語', 'English', '中文'],
        tags: ['元CA', '高級レストラン', 'ビジネス会食'],
        interests: ['ショッピング', 'カフェ巡り', '音楽'],
        zodiacSign: '双子座',
        bloodType: 'O型',
        occupation: '元CA',
        personality: '優しい',
        verified: true
    },
    {
        id: 3,
        name: 'なつき',
        age: 24,
        image: window.IMAGE_BASE_PATH ? `${window.IMAGE_BASE_PATH}/casts/natsuki.png` : 'static/images/casts/natsuki.png',
        rating: 4.8,
        reviewCount: 29,
        hourlyRate: 12000,
        rank: 'standard',
        bio: 'カジュアルな雰囲気で楽しくお話しできます♪',
        introduction: 'VIPキャストのなつきです。最高のおもてなしをお約束します。',
        serviceAreas: ['新宿', '池袋', '渋谷'],
        languages: ['日本語'],
        tags: ['学生', 'カフェ巡り', 'アニメ好き'],
        interests: ['ワイン', 'アート', 'ゴルフ'],
        zodiacSign: '蠍座',
        bloodType: 'B型',
        occupation: '大学生',
        personality: '大人っぽい',
        verified: true
    },
    {
        id: 4,
        name: 'りな',
        age: 26,
        image: window.IMAGE_BASE_PATH ? `${window.IMAGE_BASE_PATH}/casts/rina.png` : 'static/images/casts/rina.png',
        rating: 4.7,
        reviewCount: 18,
        hourlyRate: 18000,
        rank: 'premium',
        bio: '個性的なスタイルで特別な時間をお届けします。',
        introduction: 'りなです！楽しくお話ししましょう♪',
        serviceAreas: ['原宿', '表参道', '青山'],
        languages: ['日本語', 'English'],
        tags: ['アーティスト', 'ファッション', 'パーティー'],
        interests: ['ファッション', 'ダンス', 'カラオケ'],
        zodiacSign: '牡羊座',
        bloodType: 'AB型',
        occupation: 'アパレル',
        personality: '元気',
        verified: false
    },
    {
        id: 5,
        name: 'みお',
        age: 28,
        image: window.IMAGE_BASE_PATH ? `${window.IMAGE_BASE_PATH}/casts/mio.png` : 'static/images/casts/mio.png',
        rating: 4.9,
        reviewCount: 55,
        hourlyRate: 25000,
        rank: 'vip',
        bio: 'ビジネスシーンでも安心してご利用いただけます。',
        introduction: 'みおです。素敵な思い出作りのお手伝いをします♡',
        serviceAreas: ['銀座', '六本木', '赤坂'],
        languages: ['日本語', 'English', 'Français'],
        tags: ['秘書経験', 'ビジネス英語', '高級クラブ'],
        interests: ['美容', 'ヨガ', 'スイーツ巡り'],
        zodiacSign: '魚座',
        bloodType: 'A型',
        occupation: '美容師',
        personality: '癒し系',
        verified: true
    }
];

// Initialize Application
function initApp() {
    console.log('Initializing uso app...');
    
    // Check authentication
    if (app.token) {
        // In real app, verify token validity
        app.currentUser = {
            id: 1,
            name: '山田 太郎',
            email: 'yamada@example.com',
            userType: 'guest',
            membershipType: 'premium'
        };
    }

    // Load mock data
    app.casts = mockCasts;
    console.log('Loaded casts:', app.casts.length);

    // Initialize page
    renderApp();
    attachEventListeners();
    console.log('App initialization complete');
}

// Render Application
function renderApp() {
    const appContainer = document.getElementById('app');
    
    appContainer.innerHTML = `
        <div class="app-container">
            ${renderSidebar()}
            <main class="main-content">
                ${renderHeader()}
                <div class="page-content" id="pageContent">
                    ${renderPage(app.currentPage)}
                </div>
            </main>
            ${renderBottomNav()}
        </div>
    `;

    // Update active states
    updateActiveStates();
}

// Render Sidebar
function renderSidebar() {
    return `
        <aside class="sidebar" id="sidebar">
            <div class="sidebar-header">
                <h1 class="logo">uso</h1>
                <button class="sidebar-toggle" onclick="toggleSidebar()">
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M18 6L6 18M6 6l12 12"/>
                    </svg>
                </button>
            </div>

            <nav class="sidebar-nav">
                <div class="nav-section">
                    <div class="nav-section-title">メイン</div>
                    <a href="#" class="nav-link" data-page="home">
                        <span class="nav-icon">🏠</span>
                        <span class="nav-text">ホーム</span>
                    </a>
                    <a href="#" class="nav-link" data-page="swipe">
                        <span class="nav-icon">💖</span>
                        <span class="nav-text">スワイプ</span>
                    </a>
                    <a href="#" class="nav-link" data-page="search">
                        <span class="nav-icon">🔍</span>
                        <span class="nav-text">キャストを探す</span>
                    </a>
                    <a href="#" class="nav-link" data-page="likes">
                        <span class="nav-icon">❤️</span>
                        <span class="nav-text">いいね！</span>
                        <span class="nav-badge">3</span>
                    </a>
                </div>

                <div class="nav-section">
                    <div class="nav-section-title">コミュニケーション</div>
                    <a href="#" class="nav-link" data-page="messages">
                        <span class="nav-icon">💬</span>
                        <span class="nav-text">メッセージ</span>
                        <span class="nav-badge">2</span>
                    </a>
                    <a href="#" class="nav-link" data-page="bookings">
                        <span class="nav-icon">📅</span>
                        <span class="nav-text">予約管理</span>
                    </a>
                </div>

                <div class="nav-section">
                    <div class="nav-section-title">アカウント</div>
                    <a href="#" class="nav-link" data-page="profile">
                        <span class="nav-icon">👤</span>
                        <span class="nav-text">プロフィール</span>
                    </a>
                    <a href="#" class="nav-link" data-page="settings">
                        <span class="nav-icon">⚙️</span>
                        <span class="nav-text">設定</span>
                    </a>
                    <a href="#" class="nav-link" data-page="billing">
                        <span class="nav-icon">💳</span>
                        <span class="nav-text">お支払い</span>
                    </a>
                </div>
            </nav>
        </aside>
    `;
}

// Render Header
function renderHeader() {
    return `
        <header class="header">
            <button class="mobile-menu-toggle mobile-only" onclick="toggleMobileSidebar()">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M3 12h18M3 6h18M3 18h18"/>
                </svg>
            </button>

            <div class="header-search desktop-only">
                <input type="text" class="search-input" placeholder="キャスト名やエリアで検索..." onkeypress="handleSearch(event)">
            </div>

            <div class="header-actions">
                <button class="header-button" onclick="showNotifications()">
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"/>
                        <path d="M13.73 21a2 2 0 0 1-3.46 0"/>
                    </svg>
                    <span class="notification-dot"></span>
                </button>

                <div class="user-menu" onclick="toggleUserMenu()">
                    <div class="user-avatar">山</div>
                    <div class="user-info">
                        <div class="user-name">山田 太郎</div>
                        <div class="user-role">プレミアム会員</div>
                    </div>
                </div>
            </div>
        </header>
    `;
}

// Render Bottom Navigation
function renderBottomNav() {
    return `
        <nav class="bottom-nav">
            <div class="bottom-nav-items">
                <a href="#" class="bottom-nav-item" data-page="home">
                    <span class="bottom-nav-icon">🏠</span>
                    <span class="bottom-nav-label">ホーム</span>
                </a>
                <a href="#" class="bottom-nav-item" data-page="swipe">
                    <span class="bottom-nav-icon">💖</span>
                    <span class="bottom-nav-label">スワイプ</span>
                </a>
                <a href="#" class="bottom-nav-item" data-page="search">
                    <span class="bottom-nav-icon">🔍</span>
                    <span class="bottom-nav-label">探す</span>
                </a>
                <a href="#" class="bottom-nav-item" data-page="messages">
                    <span class="bottom-nav-icon">💬</span>
                    <span class="bottom-nav-label">メッセージ</span>
                    <span class="nav-badge" style="position: absolute; top: 0; right: 20%;">2</span>
                </a>
                <a href="#" class="bottom-nav-item" data-page="bookings">
                    <span class="bottom-nav-icon">📅</span>
                    <span class="bottom-nav-label">予約</span>
                </a>
                <a href="#" class="bottom-nav-item" data-page="profile">
                    <span class="bottom-nav-icon">👤</span>
                    <span class="bottom-nav-label">マイページ</span>
                </a>
            </div>
        </nav>
    `;
}

// Render Page Content
function renderPage(pageName) {
    switch (pageName) {
        case 'home':
            return renderHomePage();
        case 'swipe':
            return renderSwipePage();
        case 'search':
            return renderSearchPage();
        case 'cast-detail':
            return renderCastDetailPage();
        case 'likes':
            return renderLikesPage();
        case 'messages':
            return renderMessagesPage();
        case 'bookings':
            return renderBookingsPage();
        case 'profile':
            return renderProfilePage();
        case 'settings':
            return renderSettingsPage();
        case 'billing':
            return renderBillingPage();
        default:
            return renderHomePage();
    }
}

// Home Page
function renderHomePage() {
    const newCasts = app.casts.filter(cast => true).slice(0, 4);
    
    return `
        <div class="page-header">
            <h1 class="page-title">ダッシュボード</h1>
            <p class="page-description">あなたへのおすすめキャストをご紹介</p>
        </div>

        <div class="stats-grid">
            <div class="stat-card">
                <div class="stat-icon">❤️</div>
                <div class="stat-content">
                    <div class="stat-value">12</div>
                    <div class="stat-label">いいね！</div>
                </div>
            </div>
            <div class="stat-card">
                <div class="stat-icon">💬</div>
                <div class="stat-content">
                    <div class="stat-value">8</div>
                    <div class="stat-label">マッチング</div>
                </div>
            </div>
            <div class="stat-card">
                <div class="stat-icon">📅</div>
                <div class="stat-content">
                    <div class="stat-value">3</div>
                    <div class="stat-label">予約中</div>
                </div>
            </div>
            <div class="stat-card">
                <div class="stat-icon">⭐</div>
                <div class="stat-content">
                    <div class="stat-value">4.8</div>
                    <div class="stat-label">平均評価</div>
                </div>
            </div>
        </div>

        <div class="promo-card">
            <div class="promo-content">
                <h2>本日の無料いいね！ 🎁</h2>
                <p>今日はあと3回無料でいいね！を送れます</p>
                <button class="btn btn-primary" onclick="navigateTo('search')">
                    キャストを探す
                </button>
            </div>
        </div>

        <div class="section">
            <div class="section-header">
                <h2>新着キャスト ✨</h2>
                <a href="#" onclick="navigateTo('search')">すべて見る →</a>
            </div>
            <div class="cast-grid">
                ${newCasts.map(cast => renderCastCard(cast)).join('')}
            </div>
        </div>
    `;
}

// Swipe Page
function renderSwipePage() {
    console.log('Rendering swipe page...');
    return `
        <div class="swipe-page" id="swipePage">
            <!-- Swipe cards will be rendered here -->
        </div>
    `;
}

// Search Page
function renderSearchPage() {
    return `
        <div class="page-header">
            <h1 class="page-title">キャストを探す</h1>
            <p class="page-description">条件に合うキャストを見つけましょう</p>
        </div>

        <div class="search-filters">
            <div class="filter-group">
                <label>エリア</label>
                <select class="form-select" id="areaFilter">
                    <option value="">すべてのエリア</option>
                    <option value="六本木">六本木</option>
                    <option value="銀座">銀座</option>
                    <option value="渋谷">渋谷</option>
                    <option value="新宿">新宿</option>
                    <option value="池袋">池袋</option>
                </select>
            </div>
            <div class="filter-group">
                <label>料金</label>
                <select class="form-select" id="priceFilter">
                    <option value="">すべての料金</option>
                    <option value="0-10000">¥10,000以下</option>
                    <option value="10000-20000">¥10,000 - ¥20,000</option>
                    <option value="20000-">¥20,000以上</option>
                </select>
            </div>
            <div class="filter-group">
                <label>ランク</label>
                <select class="form-select" id="rankFilter">
                    <option value="">すべてのランク</option>
                    <option value="standard">スタンダード</option>
                    <option value="premium">プレミアム</option>
                    <option value="vip">VIP</option>
                </select>
            </div>
            <button class="btn btn-primary" onclick="applyFilters()">
                検索
            </button>
        </div>

        <div class="search-results">
            <p class="results-count">検索結果: ${app.casts.length}名</p>
            <div class="cast-grid">
                ${app.casts.map(cast => renderCastCard(cast)).join('')}
            </div>
        </div>
    `;
}

// Cast Detail Page
function renderCastDetailPage() {
    const castId = app.selectedCastId || 1;
    const cast = app.casts.find(c => c.id === castId) || app.casts[0];
    
    return `
        <div class="cast-detail">
            <button class="back-button" onclick="navigateTo('search')">
                ← 戻る
            </button>
            
            <div class="cast-detail-content">
                <div class="cast-detail-images">
                    <img src="${cast.image}" alt="${cast.name}" class="cast-main-image">
                    <div class="cast-gallery">
                        <!-- Additional images would go here -->
                    </div>
                </div>
                
                <div class="cast-detail-info">
                    <div class="cast-detail-header">
                        <h1>${cast.name}</h1>
                        <span class="cast-badge ${cast.rank}">${cast.rank.toUpperCase()}</span>
                    </div>
                    
                    <div class="cast-stats">
                        <div class="stat">
                            <span class="stat-icon">⭐</span>
                            <span>${cast.rating} (${cast.reviewCount}件)</span>
                        </div>
                        <div class="stat">
                            <span class="stat-icon">💰</span>
                            <span>¥${cast.hourlyRate.toLocaleString()}/時間</span>
                        </div>
                        <div class="stat">
                            <span class="stat-icon">🎂</span>
                            <span>${cast.age}歳</span>
                        </div>
                    </div>
                    
                    <div class="cast-bio">
                        <h3>自己紹介</h3>
                        <p>${cast.bio}</p>
                    </div>
                    
                    <div class="cast-details-section">
                        <h3>詳細情報</h3>
                        <div class="detail-item">
                            <span class="detail-label">対応エリア:</span>
                            <span>${cast.serviceAreas.join(', ')}</span>
                        </div>
                        <div class="detail-item">
                            <span class="detail-label">言語:</span>
                            <span>${cast.languages.join(', ')}</span>
                        </div>
                        <div class="detail-item">
                            <span class="detail-label">タグ:</span>
                            <div class="tag-list">
                                ${cast.tags.map(tag => `<span class="tag">${tag}</span>`).join('')}
                            </div>
                        </div>
                    </div>
                    
                    <div class="action-buttons">
                        <button class="btn btn-primary btn-large" onclick="bookCast(${cast.id})">
                            予約する
                        </button>
                        <button class="btn btn-secondary" onclick="toggleFavorite(${cast.id})">
                            ${app.favorites.has(cast.id) ? '❤️ いいね！済み' : '🤍 いいね！'}
                        </button>
                    </div>
                </div>
            </div>
            
            <div class="reviews-section">
                <h2>レビュー</h2>
                <div class="reviews-list">
                    <div class="review-card">
                        <div class="review-header">
                            <div class="reviewer-info">
                                <div class="reviewer-avatar">T</div>
                                <div>
                                    <div class="reviewer-name">田中さん</div>
                                    <div class="review-date">2024年1月10日</div>
                                </div>
                            </div>
                            <div class="review-rating">⭐⭐⭐⭐⭐</div>
                        </div>
                        <p class="review-text">とても楽しい時間を過ごせました。会話も弾み、素敵な思い出になりました。</p>
                    </div>
                </div>
            </div>
        </div>
    `;
}

// Likes Page
function renderLikesPage() {
    const likedCasts = Array.from(app.favorites).map(id => 
        app.casts.find(cast => cast.id === id)
    ).filter(Boolean);
    
    return `
        <div class="page-header">
            <h1 class="page-title">いいね！リスト</h1>
            <p class="page-description">あなたがいいね！したキャスト</p>
        </div>

        <div class="likes-tabs">
            <button class="tab-button active">いいかも (${likedCasts.length})</button>
            <button class="tab-button">相手から (0)</button>
        </div>

        <div class="likes-content">
            ${likedCasts.length > 0 ? `
                <div class="cast-grid">
                    ${likedCasts.map(cast => renderCastCard(cast)).join('')}
                </div>
            ` : `
                <div class="empty-state">
                    <p>まだいいね！したキャストはいません</p>
                    <button class="btn btn-primary" onclick="navigateTo('swipe')">
                        スワイプを始める
                    </button>
                </div>
            `}
        </div>
    `;
}

// Messages Page
function renderMessagesPage() {
    const messages = [
        { id: 1, name: 'ゆか', lastMessage: 'こんにちは！今週末はいかがですか？', time: '5分前', unread: 2, image: window.IMAGE_BASE_PATH ? `${window.IMAGE_BASE_PATH}/casts/yuka.png` : 'static/images/casts/yuka.png' },
        { id: 2, name: 'かな', lastMessage: 'ありがとうございました😊', time: '1時間前', unread: 0, image: window.IMAGE_BASE_PATH ? `${window.IMAGE_BASE_PATH}/casts/kana.png` : 'static/images/casts/kana.png' }
    ];
    
    return `
        <div class="page-header">
            <h1 class="page-title">メッセージ</h1>
            <p class="page-description">マッチしたキャストとチャット</p>
        </div>

        <div class="messages-list">
            ${messages.map(msg => `
                <div class="message-item" onclick="openChat(${msg.id})">
                    <img src="${msg.image}" alt="${msg.name}" class="message-avatar">
                    <div class="message-content">
                        <div class="message-header">
                            <h4>${msg.name}</h4>
                            <span class="message-time">${msg.time}</span>
                        </div>
                        <p class="message-preview">${msg.lastMessage}</p>
                    </div>
                    ${msg.unread > 0 ? `<span class="unread-badge">${msg.unread}</span>` : ''}
                </div>
            `).join('')}
        </div>
    `;
}

// Bookings Page
function renderBookingsPage() {
    return `
        <div class="page-header">
            <h1 class="page-title">予約管理</h1>
            <p class="page-description">予約の確認と管理</p>
        </div>

        <div class="bookings-tabs">
            <button class="tab-button active" onclick="showBookingTab('upcoming')">予約中</button>
            <button class="tab-button" onclick="showBookingTab('past')">過去の予約</button>
        </div>

        <div id="upcomingBookings" class="bookings-content">
            <div class="booking-card">
                <div class="booking-header">
                    <img src="${window.IMAGE_BASE_PATH ? `${window.IMAGE_BASE_PATH}/casts/kana.png` : 'static/images/casts/kana.png'}" alt="かな" class="booking-avatar">
                    <div class="booking-info">
                        <h3>かな さん</h3>
                        <div class="booking-details">
                            <span>📅 2024年1月20日（土）</span>
                            <span>🕐 19:00 - 22:00</span>
                            <span>📍 銀座</span>
                        </div>
                    </div>
                    <span class="booking-status confirmed">確定</span>
                </div>
                <div class="booking-actions">
                    <button class="btn btn-secondary" onclick="showBookingDetails(1)">詳細</button>
                    <button class="btn btn-primary" onclick="openChat(2)">メッセージ</button>
                </div>
            </div>
        </div>

        <div id="pastBookings" class="bookings-content hidden">
            <div class="booking-card">
                <div class="booking-header">
                    <img src="${window.IMAGE_BASE_PATH ? `${window.IMAGE_BASE_PATH}/casts/mio.png` : 'static/images/casts/mio.png'}" alt="みお" class="booking-avatar">
                    <div class="booking-info">
                        <h3>みお さん</h3>
                        <div class="booking-details">
                            <span>📅 2024年1月5日（金）</span>
                            <span>🕐 20:00 - 23:00</span>
                            <span>📍 六本木</span>
                        </div>
                    </div>
                    <span class="booking-status completed">完了</span>
                </div>
                <div class="booking-actions">
                    <button class="btn btn-primary" onclick="writeReview(5)">レビューを書く</button>
                </div>
            </div>
        </div>
    `;
}

// Profile Page
function renderProfilePage() {
    return `
        <div class="page-header">
            <h1 class="page-title">マイページ</h1>
            <p class="page-description">プロフィールの確認と編集</p>
        </div>

        <div class="profile-section">
            <div class="profile-header">
                <div class="profile-avatar">
                    <div class="avatar-placeholder">山</div>
                    <button class="avatar-edit">編集</button>
                </div>
                <div class="profile-info">
                    <h2>山田 太郎</h2>
                    <p class="profile-status">プレミアム会員</p>
                    <div class="profile-stats">
                        <div class="stat">
                            <span class="stat-value">12</span>
                            <span class="stat-label">マッチ</span>
                        </div>
                        <div class="stat">
                            <span class="stat-value">8</span>
                            <span class="stat-label">予約</span>
                        </div>
                        <div class="stat">
                            <span class="stat-value">4.8</span>
                            <span class="stat-label">評価</span>
                        </div>
                    </div>
                </div>
            </div>

            <div class="profile-actions">
                <button class="btn btn-secondary" onclick="navigateTo('settings')">
                    プロフィール編集
                </button>
                <button class="btn btn-primary" onclick="showUpgradePlan()">
                    プランをアップグレード
                </button>
            </div>
        </div>

        <div class="profile-menu">
            <a href="#" class="menu-item" onclick="navigateTo('billing')">
                <span class="menu-icon">💳</span>
                <span class="menu-text">支払い方法</span>
                <span class="menu-arrow">›</span>
            </a>
            <a href="#" class="menu-item" onclick="navigateTo('settings')">
                <span class="menu-icon">⚙️</span>
                <span class="menu-text">設定</span>
                <span class="menu-arrow">›</span>
            </a>
            <a href="#" class="menu-item" onclick="showHelp()">
                <span class="menu-icon">❓</span>
                <span class="menu-text">ヘルプ</span>
                <span class="menu-arrow">›</span>
            </a>
            <a href="#" class="menu-item" onclick="showTerms()">
                <span class="menu-icon">📄</span>
                <span class="menu-text">利用規約</span>
                <span class="menu-arrow">›</span>
            </a>
        </div>
    `;
}

// Billing Page
function renderBillingPage() {
    return `
        <div class="page-header">
            <h1 class="page-title">お支払い</h1>
            <p class="page-description">プランと支払い方法の管理</p>
        </div>

        <div class="billing-section">
            <h3>現在のプラン</h3>
            <div class="plan-card active">
                <div class="plan-header">
                    <h4>プレミアムプラン</h4>
                    <span class="plan-price">¥9,800/月</span>
                </div>
                <ul class="plan-features">
                    <li>✅ 無制限のいいね！</li>
                    <li>✅ プロフィール優先表示</li>
                    <li>✅ 高度な検索フィルター</li>
                    <li>✅ リワインド機能</li>
                </ul>
                <button class="btn btn-secondary">プラン変更</button>
            </div>
        </div>

        <div class="billing-section">
            <h3>支払い方法</h3>
            <div class="payment-method">
                <div class="payment-icon">💳</div>
                <div class="payment-info">
                    <p>•••• •••• •••• 4242</p>
                    <p class="payment-expire">有効期限: 12/25</p>
                </div>
                <button class="btn btn-small">変更</button>
            </div>
        </div>

        <div class="billing-section">
            <h3>請求履歴</h3>
            <div class="invoice-list">
                <div class="invoice-item">
                    <div class="invoice-date">2024年1月1日</div>
                    <div class="invoice-amount">¥9,800</div>
                    <a href="#" class="invoice-download">ダウンロード</a>
                </div>
            </div>
        </div>
    `;
}

// Settings Page
function renderSettingsPage() {
    return `
        <div class="page-header">
            <h1 class="page-title">設定</h1>
            <p class="page-description">アカウントと通知の設定</p>
        </div>

        <div class="settings-section">
            <h3>通知設定</h3>
            <div class="setting-item">
                <div class="setting-info">
                    <h4>新着メッセージ</h4>
                    <p>新しいメッセージを受信したときに通知</p>
                </div>
                <label class="switch">
                    <input type="checkbox" checked>
                    <span class="slider"></span>
                </label>
            </div>
            <div class="setting-item">
                <div class="setting-info">
                    <h4>予約リマインダー</h4>
                    <p>予約の1時間前に通知</p>
                </div>
                <label class="switch">
                    <input type="checkbox" checked>
                    <span class="slider"></span>
                </label>
            </div>
        </div>

        <div class="settings-section">
            <h3>プライバシー設定</h3>
            <div class="setting-item">
                <div class="setting-info">
                    <h4>プロフィール公開</h4>
                    <p>他のユーザーにプロフィールを表示</p>
                </div>
                <label class="switch">
                    <input type="checkbox" checked>
                    <span class="slider"></span>
                </label>
            </div>
        </div>

        <div class="settings-section">
            <h3>アカウント</h3>
            <button class="btn btn-secondary" onclick="changePassword()">パスワード変更</button>
            <button class="btn btn-danger" onclick="logout()">ログアウト</button>
        </div>
    `;
}

// Cast Card Component
function renderCastCard(cast) {
    const isFavorite = app.favorites.has(cast.id);
    
    return `
        <div class="cast-card" onclick="viewCastDetail(${cast.id})">
            <div class="cast-image">
                <img src="${cast.image}" alt="${cast.name}" style="width: 100%; height: 100%; object-fit: cover;">
                ${cast.rank === 'vip' ? '<span class="cast-badge vip">VIP</span>' : ''}
                ${cast.rank === 'premium' ? '<span class="cast-badge premium">PREMIUM</span>' : ''}
                <button class="cast-favorite ${isFavorite ? 'active' : ''}" onclick="event.stopPropagation(); toggleFavorite(${cast.id})">
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="${isFavorite ? 'currentColor' : 'none'}" stroke="currentColor" stroke-width="2">
                        <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"/>
                    </svg>
                </button>
            </div>
            <div class="cast-info">
                <h3 class="cast-name">${cast.name}</h3>
                <div class="cast-details">
                    <span>⭐ ${cast.rating}</span>
                    <span>•</span>
                    <span>${cast.age}歳</span>
                    <span>•</span>
                    <span class="cast-price">¥${cast.hourlyRate.toLocaleString()}/h</span>
                </div>
                <div class="cast-tags">
                    ${cast.serviceAreas.slice(0, 2).map(area => 
                        `<span class="cast-tag">${area}</span>`
                    ).join('')}
                </div>
            </div>
        </div>
    `;
}

// Navigation Functions
function navigateTo(page) {
    app.currentPage = page;
    renderApp();
    window.scrollTo(0, 0);
    
    // Initialize swipe cards if navigating to swipe page
    if (page === 'swipe') {
        setTimeout(() => {
            const swipeContainer = document.getElementById('swipePage');
            console.log('Swipe container:', swipeContainer);
            console.log('SwipeCards class available:', typeof SwipeCards !== 'undefined');
            console.log('Available casts:', app.casts.length);
            
            if (swipeContainer && typeof SwipeCards !== 'undefined') {
                window.swipeCard = new SwipeCards(swipeContainer, app.casts);
                console.log('SwipeCards initialized successfully');
            } else {
                console.error('SwipeCards class not found or container not available');
            }
        }, 100);
    }
}

function viewCastDetail(castId) {
    app.selectedCastId = castId;
    navigateTo('cast-detail');
}

// Toggle Functions
function toggleSidebar() {
    const sidebar = document.getElementById('sidebar');
    sidebar.classList.toggle('collapsed');
}

function toggleMobileSidebar() {
    const sidebar = document.getElementById('sidebar');
    sidebar.classList.toggle('open');
}

function toggleFavorite(castId) {
    if (app.favorites.has(castId)) {
        app.favorites.delete(castId);
        showNotification('いいね！を取り消しました', 'info');
    } else {
        app.favorites.add(castId);
        showNotification('いいね！しました', 'success');
    }
    renderApp();
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

// Update Active States
function updateActiveStates() {
    // Update sidebar links
    document.querySelectorAll('.nav-link').forEach(link => {
        const page = link.getAttribute('data-page');
        if (page === app.currentPage) {
            link.classList.add('active');
        } else {
            link.classList.remove('active');
        }
    });

    // Update bottom nav
    document.querySelectorAll('.bottom-nav-item').forEach(item => {
        const page = item.getAttribute('data-page');
        if (page === app.currentPage) {
            item.classList.add('active');
        } else {
            item.classList.remove('active');
        }
    });
}

// Event Listeners
function attachEventListeners() {
    // Navigation clicks
    document.addEventListener('click', (e) => {
        if (e.target.closest('.nav-link')) {
            e.preventDefault();
            const page = e.target.closest('.nav-link').getAttribute('data-page');
            navigateTo(page);
        }
        
        if (e.target.closest('.bottom-nav-item')) {
            e.preventDefault();
            const page = e.target.closest('.bottom-nav-item').getAttribute('data-page');
            navigateTo(page);
        }
    });

    // Close mobile sidebar on outside click
    document.addEventListener('click', (e) => {
        const sidebar = document.getElementById('sidebar');
        const toggle = document.querySelector('.mobile-menu-toggle');
        
        if (sidebar && sidebar.classList.contains('open') && 
            !sidebar.contains(e.target) && 
            !toggle.contains(e.target)) {
            sidebar.classList.remove('open');
        }
    });
}

// Utility Functions
function handleSearch(event) {
    if (event.key === 'Enter') {
        navigateTo('search');
    }
}

function bookCast(castId) {
    showNotification('予約画面を開きます...', 'info');
    // In real app, open booking modal
}

function openChat(userId) {
    console.log('Opening chat for user:', userId);
    
    // Find the cast by ID
    const cast = app.casts.find(c => c.id === userId);
    if (!cast) {
        showNotification('キャストが見つかりません', 'error');
        return;
    }
    
    // Create chat modal
    const chatModal = document.createElement('div');
    chatModal.className = 'chat-modal';
    chatModal.innerHTML = `
        <div class="chat-container">
            <div class="chat-header">
                <img src="${cast.image}" alt="${cast.name}" class="chat-avatar">
                <div class="chat-user-info">
                    <h3>${cast.name}</h3>
                    <span class="chat-status">オンライン</span>
                </div>
                <button class="chat-close" onclick="closeChat()">×</button>
            </div>
            <div class="chat-messages" id="chatMessages">
                <div class="message received">
                    <div class="message-content">
                        <p>こんにちは！${cast.name}です 😊</p>
                        <span class="message-time">${new Date().toLocaleTimeString()}</span>
                    </div>
                </div>
            </div>
            <div class="chat-input">
                <input type="text" id="chatInput" placeholder="メッセージを入力..." onkeypress="sendMessage(event)">
                <button onclick="sendMessage()">送信</button>
            </div>
        </div>
    `;
    
    document.body.appendChild(chatModal);
    
    // Focus on input
    setTimeout(() => {
        document.getElementById('chatInput').focus();
    }, 100);
}

function closeChat() {
    const chatModal = document.querySelector('.chat-modal');
    if (chatModal) {
        chatModal.remove();
    }
}

function sendMessage(event) {
    if (event && event.key !== 'Enter') return;
    
    const input = document.getElementById('chatInput');
    const message = input.value.trim();
    
    if (!message) return;
    
    const messagesContainer = document.getElementById('chatMessages');
    const messageElement = document.createElement('div');
    messageElement.className = 'message sent';
    messageElement.innerHTML = `
        <div class="message-content">
            <p>${message}</p>
            <span class="message-time">${new Date().toLocaleTimeString()}</span>
        </div>
    `;
    
    messagesContainer.appendChild(messageElement);
    input.value = '';
    
    // Scroll to bottom
    messagesContainer.scrollTop = messagesContainer.scrollHeight;
    
    // Simulate reply after 1 second
    setTimeout(() => {
        const replyElement = document.createElement('div');
        replyElement.className = 'message received';
        replyElement.innerHTML = `
            <div class="message-content">
                <p>ありがとうございます！😊</p>
                <span class="message-time">${new Date().toLocaleTimeString()}</span>
            </div>
        `;
        messagesContainer.appendChild(replyElement);
        messagesContainer.scrollTop = messagesContainer.scrollHeight;
    }, 1000);
}

function showBookingTab(tab) {
    document.querySelectorAll('.tab-button').forEach(btn => {
        btn.classList.remove('active');
    });
    event.target.classList.add('active');
    
    document.querySelectorAll('.bookings-content').forEach(content => {
        content.classList.add('hidden');
    });
    
    if (tab === 'upcoming') {
        document.getElementById('upcomingBookings').classList.remove('hidden');
    } else {
        document.getElementById('pastBookings').classList.remove('hidden');
    }
}

function showUpgradePlan() {
    showNotification('プランアップグレード画面を開きます...', 'info');
}

function showHelp() {
    showNotification('ヘルプセンターを開きます...', 'info');
}

function showTerms() {
    showNotification('利用規約を開きます...', 'info');
}

function writeReview(castId) {
    showNotification('レビュー画面を開きます...', 'info');
}

function showBookingDetails(bookingId) {
    showNotification('予約詳細を開きます...', 'info');
}

function changePassword() {
    showNotification('パスワード変更画面を開きます...', 'info');
}

function showNotifications() {
    showNotification('通知センターを開きます...', 'info');
}

function toggleUserMenu() {
    showNotification('ユーザーメニューを開きます...', 'info');
}

function applyFilters() {
    showNotification('フィルターを適用しました', 'success');
    // In real app, apply filters and refresh results
}

function logout() {
    if (confirm('ログアウトしますか？')) {
        localStorage.removeItem('uso_token');
        app.token = null;
        app.currentUser = null;
        showNotification('ログアウトしました', 'success');
        window.location.reload();
    }
}

// Initialize app when DOM is ready
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initApp);
} else {
    initApp();
}