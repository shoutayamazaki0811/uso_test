// Swipe Cards Module for Uso App

class SwipeCards {
    constructor(container, cards) {
        this.container = container;
        this.cards = cards;
        this.currentIndex = 0;
        this.isAnimating = false;
        this.dailyLikes = 30;
        this.remainingTime = this.getTimeUntilMidnight();
        this.init();
    }

    init() {
        this.render();
        this.attachEventListeners();
        this.startTimer();
    }

    render() {
        this.container.innerHTML = `
            <div class="swipe-container">
                <div class="swipe-header">
                    <div class="swipe-stats">
                        <div class="points-badge">
                            <span class="points-icon">🌟</span>
                            <span>ポイントGETまであと${this.dailyLikes}人</span>
                        </div>
                        <div class="timer-badge">
                            <span>${this.formatTime(this.remainingTime)}</span>
                        </div>
                    </div>
                </div>

                <div class="cards-stack" id="cardsStack">
                    ${this.renderCards()}
                </div>

                <div class="swipe-actions">
                    <button class="action-button skip" onclick="swipeCard.skip()">
                        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <path d="M18 6L6 18M6 6l12 12"/>
                        </svg>
                    </button>
                    <button class="action-button superlike" onclick="swipeCard.superlike()">
                        <svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
                            <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
                        </svg>
                    </button>
                    <button class="action-button like" onclick="swipeCard.like()">
                        <svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
                            <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"/>
                        </svg>
                    </button>
                    <button class="action-button boost" onclick="swipeCard.openBoost()">
                        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z"/>
                        </svg>
                    </button>
                    <button class="action-button rewind" onclick="swipeCard.rewind()">
                        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <path d="M1 4v6h6M23 20v-6h-6"/>
                            <path d="M20.49 9A9 9 0 0 0 5.64 5.64L1 10m22 4l-4.64 4.36A9 9 0 0 1 3.51 15"/>
                        </svg>
                    </button>
                </div>

                <div class="chance-time-popup hidden" id="chanceTimePopup">
                    <div class="chance-time-content">
                        <h2>チャンスタイム残り</h2>
                        <div class="chance-timer">29:14</div>
                        <p>チャンスタイム中です</p>
                        <p>今だけ「いいかも」30回するまで無料でフリック！大量マッチのチャンス！</p>
                        <button class="btn btn-primary" onclick="swipeCard.closeChanceTime()">とじる</button>
                    </div>
                </div>
            </div>
        `;
    }

    renderCards() {
        return this.cards.slice(this.currentIndex, this.currentIndex + 3)
            .reverse()
            .map((cast, index) => this.renderCard(cast, index))
            .join('');
    }

    renderCard(cast, stackIndex) {
        const zIndex = 3 - stackIndex;
        const scale = 1 - (stackIndex * 0.03);
        const translateY = stackIndex * 10;
        
        return `
            <div class="swipe-card" 
                 data-cast-id="${cast.id}" 
                 style="z-index: ${zIndex}; transform: scale(${scale}) translateY(${translateY}px);">
                <div class="card-image-container">
                    <img src="${cast.image}" alt="${cast.name}" class="card-image">
                    ${cast.verified ? '<div class="verified-badge">✓ 本人確認済み</div>' : ''}
                    <div class="card-gradient"></div>
                </div>
                
                <div class="card-info">
                    <div class="card-header">
                        <h2>${cast.name} <span class="age">${cast.age}歳</span></h2>
                        <span class="location">${cast.serviceAreas[0]}</span>
                    </div>
                    
                    <div class="card-tags">
                        ${cast.interests.slice(0, 3).map(interest => 
                            `<span class="interest-tag">${interest}</span>`
                        ).join('')}
                    </div>
                    
                    <div class="card-details">
                        <div class="detail-item">
                            <span class="detail-icon">💼</span>
                            <span>${cast.occupation || '会社員'}</span>
                        </div>
                        <div class="detail-item">
                            <span class="detail-icon">🎯</span>
                            <span>${cast.personality || '明るい'}</span>
                        </div>
                        ${cast.zodiacSign ? `
                        <div class="detail-item">
                            <span class="detail-icon">♈</span>
                            <span>${cast.zodiacSign}</span>
                        </div>
                        ` : ''}
                    </div>
                    
                    <p class="card-bio">${cast.introduction}</p>
                </div>
                
                <div class="card-action-hint nope">NOPE</div>
                <div class="card-action-hint like">LIKE</div>
                <div class="card-action-hint super">SUPER LIKE</div>
            </div>
        `;
    }

    attachEventListeners() {
        const stack = document.getElementById('cardsStack');
        let startX = null;
        let startY = null;
        let currentX = null;
        let currentY = null;
        let card = null;
        let isDragging = false;

        const handleStart = (e) => {
            if (this.isAnimating) return;
            
            const topCard = stack.querySelector('.swipe-card');
            if (!topCard || topCard !== e.target.closest('.swipe-card')) return;
            
            card = topCard;
            isDragging = true;
            
            const touch = e.touches ? e.touches[0] : e;
            startX = touch.clientX;
            startY = touch.clientY;
            currentX = startX;
            currentY = startY;
            
            card.style.transition = 'none';
        };

        const handleMove = (e) => {
            if (!isDragging || !card) return;
            
            const touch = e.touches ? e.touches[0] : e;
            currentX = touch.clientX;
            currentY = touch.clientY;
            
            const deltaX = currentX - startX;
            const deltaY = currentY - startY;
            const rotation = deltaX * 0.1;
            
            card.style.transform = `translate(${deltaX}px, ${deltaY}px) rotate(${rotation}deg)`;
            
            // Show action hints
            const threshold = 50;
            card.querySelectorAll('.card-action-hint').forEach(hint => {
                hint.style.opacity = '0';
            });
            
            if (deltaX > threshold) {
                card.querySelector('.card-action-hint.like').style.opacity = Math.min(1, deltaX / 150);
            } else if (deltaX < -threshold) {
                card.querySelector('.card-action-hint.nope').style.opacity = Math.min(1, Math.abs(deltaX) / 150);
            } else if (deltaY < -threshold) {
                card.querySelector('.card-action-hint.super').style.opacity = Math.min(1, Math.abs(deltaY) / 150);
            }
        };

        const handleEnd = (e) => {
            if (!isDragging || !card) return;
            
            isDragging = false;
            card.style.transition = '';
            
            const deltaX = currentX - startX;
            const deltaY = currentY - startY;
            const threshold = 100;
            
            if (Math.abs(deltaX) > threshold || Math.abs(deltaY) > threshold) {
                if (deltaY < -threshold) {
                    this.superlike();
                } else if (deltaX > threshold) {
                    this.like();
                } else if (deltaX < -threshold) {
                    this.skip();
                } else {
                    this.resetCard(card);
                }
            } else {
                this.resetCard(card);
            }
            
            card = null;
        };

        // Touch events
        stack.addEventListener('touchstart', handleStart, { passive: true });
        stack.addEventListener('touchmove', handleMove, { passive: true });
        stack.addEventListener('touchend', handleEnd, { passive: true });
        
        // Mouse events
        stack.addEventListener('mousedown', handleStart);
        document.addEventListener('mousemove', handleMove);
        document.addEventListener('mouseup', handleEnd);
    }

    resetCard(card) {
        card.style.transform = '';
        card.querySelectorAll('.card-action-hint').forEach(hint => {
            hint.style.opacity = '0';
        });
    }

    async animateCard(direction) {
        if (this.isAnimating) return;
        
        const card = document.querySelector('.swipe-card');
        if (!card) return;
        
        this.isAnimating = true;
        
        let transform = '';
        let hint = '';
        
        switch (direction) {
            case 'left':
                transform = 'translateX(-150%) rotate(-30deg)';
                hint = '.nope';
                break;
            case 'right':
                transform = 'translateX(150%) rotate(30deg)';
                hint = '.like';
                break;
            case 'up':
                transform = 'translateY(-150%) rotate(10deg) scale(0.8)';
                hint = '.super';
                break;
        }
        
        if (hint) {
            card.querySelector(hint).style.opacity = '1';
        }
        
        card.style.transform = transform;
        
        await new Promise(resolve => setTimeout(resolve, 300));
        
        card.remove();
        this.currentIndex++;
        
        if (this.currentIndex < this.cards.length) {
            this.render();
            this.attachEventListeners();
        } else {
            this.showNoMoreCards();
        }
        
        this.isAnimating = false;
        this.dailyLikes--;
        this.updateStats();
    }

    skip() {
        this.animateCard('left');
    }

    like() {
        this.animateCard('right');
        this.checkForMatch();
    }

    superlike() {
        this.animateCard('up');
        this.checkForMatch(true);
    }

    checkForMatch(isSuperLike = false) {
        // Simulate match probability
        const matchChance = isSuperLike ? 0.7 : 0.3;
        if (Math.random() < matchChance) {
            setTimeout(() => this.showMatch(), 500);
        }
    }

    showMatch() {
        const currentCast = this.cards[this.currentIndex - 1];
        const popup = document.createElement('div');
        popup.className = 'match-popup';
        popup.innerHTML = `
            <div class="match-content">
                <h1>マッチしました！</h1>
                <div class="match-avatars">
                    <img src="${currentCast.image}" alt="${currentCast.name}">
                </div>
                <p>${currentCast.name}さんとマッチしました！</p>
                <button class="btn btn-primary" onclick="swipeCard.closeMatch(this)">メッセージを送る</button>
                <button class="btn btn-secondary" onclick="swipeCard.closeMatch(this)">後で</button>
            </div>
        `;
        document.body.appendChild(popup);
        
        setTimeout(() => popup.classList.add('show'), 10);
    }

    closeMatch(button) {
        const popup = button.closest('.match-popup');
        popup.classList.remove('show');
        setTimeout(() => popup.remove(), 300);
    }

    openBoost() {
        showNotification('ブースト機能は準備中です', 'info');
    }

    rewind() {
        if (this.currentIndex > 0) {
            this.currentIndex--;
            this.render();
            this.attachEventListeners();
            showNotification('1つ前のカードに戻りました', 'info');
        }
    }

    closeChanceTime() {
        document.getElementById('chanceTimePopup').classList.add('hidden');
    }

    updateStats() {
        const pointsBadge = document.querySelector('.points-badge span:last-child');
        if (pointsBadge) {
            pointsBadge.textContent = `ポイントGETまであと${this.dailyLikes}人`;
        }
        
        if (this.dailyLikes === 0) {
            this.showDailyLimitReached();
        }
    }

    showDailyLimitReached() {
        const popup = document.createElement('div');
        popup.className = 'limit-popup';
        popup.innerHTML = `
            <div class="limit-content">
                <h2>本日の無料いいね！が終了しました</h2>
                <p>明日また30回無料でいいね！ができます</p>
                <p>今すぐ続けたい場合は、プレミアムプランをご検討ください</p>
                <button class="btn btn-primary">プレミアムプランを見る</button>
                <button class="btn btn-secondary" onclick="this.closest('.limit-popup').remove()">閉じる</button>
            </div>
        `;
        document.body.appendChild(popup);
    }

    showNoMoreCards() {
        const stack = document.getElementById('cardsStack');
        stack.innerHTML = `
            <div class="no-more-cards">
                <h2>本日の新しいキャストは以上です</h2>
                <p>明日また新しいキャストが追加されます</p>
                <button class="btn btn-primary" onclick="navigateTo('search')">
                    他のキャストを探す
                </button>
            </div>
        `;
    }

    startTimer() {
        setInterval(() => {
            this.remainingTime--;
            if (this.remainingTime <= 0) {
                this.remainingTime = this.getTimeUntilMidnight();
                this.dailyLikes = 30;
                this.updateStats();
            }
            this.updateTimer();
        }, 1000);
    }

    updateTimer() {
        const timerBadge = document.querySelector('.timer-badge span');
        if (timerBadge) {
            timerBadge.textContent = this.formatTime(this.remainingTime);
        }
    }

    getTimeUntilMidnight() {
        const now = new Date();
        const midnight = new Date(now);
        midnight.setHours(24, 0, 0, 0);
        return Math.floor((midnight - now) / 1000);
    }

    formatTime(seconds) {
        const hours = Math.floor(seconds / 3600);
        const minutes = Math.floor((seconds % 3600) / 60);
        return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}`;
    }
}

// Global instance
let swipeCard = null;