// Report and Block Module for Uso App

class ReportBlock {
    constructor() {
        this.reportReasons = [
            { value: 'spam', label: 'スパム・勧誘行為' },
            { value: 'inappropriate', label: '不適切な内容' },
            { value: 'harassment', label: '嫌がらせ・ハラスメント' },
            { value: 'fake', label: '偽のプロフィール' },
            { value: 'scam', label: '詐欺・金銭要求' },
            { value: 'other', label: 'その他' }
        ];
    }

    // Show report modal
    showReportModal(userId, userName) {
        const modal = document.createElement('div');
        modal.className = 'report-modal';
        modal.innerHTML = `
            <div class="modal-overlay" onclick="reportBlock.closeModal()"></div>
            <div class="modal-content">
                <div class="modal-header">
                    <h3>ユーザーを通報</h3>
                    <button class="close-btn" onclick="reportBlock.closeModal()">×</button>
                </div>
                
                <div class="modal-body">
                    <p class="report-target">通報対象: <strong>${userName}</strong></p>
                    
                    <div class="form-group">
                        <label>通報理由を選択してください</label>
                        <select class="form-select" id="reportReason" required>
                            <option value="">選択してください</option>
                            ${this.reportReasons.map(reason => 
                                `<option value="${reason.value}">${reason.label}</option>`
                            ).join('')}
                        </select>
                    </div>
                    
                    <div class="form-group">
                        <label>詳細（任意）</label>
                        <textarea 
                            class="form-textarea" 
                            id="reportDescription" 
                            rows="4" 
                            placeholder="具体的な状況を記載してください"
                        ></textarea>
                    </div>
                    
                    <div class="report-info">
                        <p>📌 通報内容は運営チームが確認し、適切に対処いたします</p>
                        <p>📌 虚偽の通報は利用規約違反となります</p>
                    </div>
                </div>
                
                <div class="modal-footer">
                    <button class="btn btn-secondary" onclick="reportBlock.closeModal()">
                        キャンセル
                    </button>
                    <button class="btn btn-danger" onclick="reportBlock.submitReport(${userId})">
                        通報する
                    </button>
                </div>
            </div>
        `;
        
        document.body.appendChild(modal);
        setTimeout(() => modal.classList.add('show'), 10);
    }

    // Show block confirmation
    showBlockConfirmation(userId, userName) {
        const modal = document.createElement('div');
        modal.className = 'block-modal';
        modal.innerHTML = `
            <div class="modal-overlay" onclick="reportBlock.closeModal()"></div>
            <div class="modal-content">
                <div class="modal-header">
                    <h3>ユーザーをブロック</h3>
                    <button class="close-btn" onclick="reportBlock.closeModal()">×</button>
                </div>
                
                <div class="modal-body">
                    <p><strong>${userName}</strong>さんをブロックしますか？</p>
                    
                    <div class="block-effects">
                        <h4>ブロックすると：</h4>
                        <ul>
                            <li>相手からメッセージを受信しなくなります</li>
                            <li>相手のプロフィールが表示されなくなります</li>
                            <li>相手からあなたのプロフィールも見えなくなります</li>
                            <li>既存のマッチングは解除されます</li>
                        </ul>
                    </div>
                    
                    <p class="warning-text">※ブロックはいつでも解除できます</p>
                </div>
                
                <div class="modal-footer">
                    <button class="btn btn-secondary" onclick="reportBlock.closeModal()">
                        キャンセル
                    </button>
                    <button class="btn btn-danger" onclick="reportBlock.blockUser(${userId})">
                        ブロックする
                    </button>
                </div>
            </div>
        `;
        
        document.body.appendChild(modal);
        setTimeout(() => modal.classList.add('show'), 10);
    }

    // Submit report
    async submitReport(userId) {
        const reason = document.getElementById('reportReason').value;
        const description = document.getElementById('reportDescription').value;
        
        if (!reason) {
            showNotification('通報理由を選択してください', 'error');
            return;
        }
        
        try {
            const response = await apiCall('/api/report', {
                method: 'POST',
                body: JSON.stringify({
                    reported_user_id: userId,
                    reason: reason,
                    description: description
                })
            });
            
            showNotification('通報が送信されました', 'success');
            this.closeModal();
            
            // Optionally show block option
            setTimeout(() => {
                if (confirm('このユーザーをブロックしますか？')) {
                    this.blockUser(userId);
                }
            }, 1000);
            
        } catch (error) {
            showNotification('通報の送信に失敗しました', 'error');
        }
    }

    // Block user
    async blockUser(userId) {
        try {
            const response = await apiCall('/api/block', {
                method: 'POST',
                body: JSON.stringify({
                    blocked_user_id: userId
                })
            });
            
            showNotification('ユーザーをブロックしました', 'success');
            this.closeModal();
            
            // Update UI to reflect block
            this.updateUIAfterBlock(userId);
            
        } catch (error) {
            showNotification('ブロックに失敗しました', 'error');
        }
    }

    // Unblock user
    async unblockUser(userId) {
        try {
            const response = await apiCall(`/api/block/${userId}`, {
                method: 'DELETE'
            });
            
            showNotification('ブロックを解除しました', 'success');
            
            // Update UI
            this.updateUIAfterUnblock(userId);
            
        } catch (error) {
            showNotification('ブロック解除に失敗しました', 'error');
        }
    }

    // Get blocked users list
    async getBlockedUsers() {
        try {
            const response = await apiCall('/api/blocked-users');
            return response.blocked_users || [];
        } catch (error) {
            console.error('Failed to get blocked users:', error);
            return [];
        }
    }

    // Render blocked users list
    renderBlockedUsersList() {
        return `
            <div class="blocked-users-section">
                <h3>ブロックリスト</h3>
                <div class="blocked-users-list" id="blockedUsersList">
                    <div class="loading">読み込み中...</div>
                </div>
            </div>
        `;
    }

    // Load and display blocked users
    async loadBlockedUsers() {
        const container = document.getElementById('blockedUsersList');
        if (!container) return;
        
        try {
            const blockedUsers = await this.getBlockedUsers();
            
            if (blockedUsers.length === 0) {
                container.innerHTML = '<p class="empty-state">ブロックしているユーザーはいません</p>';
                return;
            }
            
            container.innerHTML = blockedUsers.map(user => `
                <div class="blocked-user-item">
                    <img src="${user.image || '/static/images/default-avatar.png'}" alt="${user.name}" class="user-avatar">
                    <div class="user-info">
                        <h4>${user.name}</h4>
                        <p class="blocked-date">ブロック日: ${new Date(user.blocked_at).toLocaleDateString('ja-JP')}</p>
                    </div>
                    <button class="btn btn-small btn-secondary" onclick="reportBlock.unblockUser(${user.id})">
                        ブロック解除
                    </button>
                </div>
            `).join('');
            
        } catch (error) {
            container.innerHTML = '<p class="error-state">ブロックリストの読み込みに失敗しました</p>';
        }
    }

    // Add report/block buttons to user profiles
    addReportBlockButtons(userId, userName) {
        return `
            <div class="user-actions">
                <button class="btn btn-small btn-outline" onclick="reportBlock.showReportModal(${userId}, '${userName}')">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M4 15s1-1 4-1 5 2 8 2 4-1 4-1V3s-1 1-4 1-5-2-8-2-4 1-4 1z"/>
                        <line x1="4" y1="22" x2="4" y2="15"/>
                    </svg>
                    通報
                </button>
                <button class="btn btn-small btn-outline" onclick="reportBlock.showBlockConfirmation(${userId}, '${userName}')">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <circle cx="12" cy="12" r="10"/>
                        <line x1="4.93" y1="4.93" x2="19.07" y2="19.07"/>
                    </svg>
                    ブロック
                </button>
            </div>
        `;
    }

    // Update UI after block
    updateUIAfterBlock(userId) {
        // Remove from current view
        const userElements = document.querySelectorAll(`[data-user-id="${userId}"]`);
        userElements.forEach(el => el.remove());
        
        // Update counts
        const matchCount = document.querySelector('.match-count');
        if (matchCount) {
            const currentCount = parseInt(matchCount.textContent) || 0;
            matchCount.textContent = Math.max(0, currentCount - 1);
        }
    }

    // Update UI after unblock
    updateUIAfterUnblock(userId) {
        // Reload blocked users list if visible
        if (document.getElementById('blockedUsersList')) {
            this.loadBlockedUsers();
        }
    }

    // Close modal
    closeModal() {
        const modals = document.querySelectorAll('.report-modal, .block-modal');
        modals.forEach(modal => {
            modal.classList.remove('show');
            setTimeout(() => modal.remove(), 300);
        });
    }
}

// Global instance
const reportBlock = new ReportBlock();

// API call helper (if not already defined)
if (typeof apiCall === 'undefined') {
    async function apiCall(endpoint, options = {}) {
        const headers = {
            'Content-Type': 'application/json',
            ...options.headers
        };

        const token = localStorage.getItem('uso_token');
        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }

        try {
            const response = await fetch(endpoint, {
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
}

// Notification helper (if not already defined)
if (typeof showNotification === 'undefined') {
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
}