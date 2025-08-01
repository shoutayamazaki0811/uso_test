// Profile Verification Module for Uso App

class ProfileVerification {
    constructor() {
        this.verificationSteps = [
            'photo_id',
            'selfie',
            'profile_info'
        ];
        this.currentStep = 0;
        this.verificationData = {};
    }

    render() {
        return `
            <div class="verification-container">
                <div class="verification-header">
                    <h2>プロフィール認証</h2>
                    <p>本人確認を行うことで、信頼性が向上します</p>
                </div>

                <div class="verification-progress">
                    <div class="progress-bar">
                        <div class="progress-fill" style="width: ${((this.currentStep + 1) / this.verificationSteps.length) * 100}%"></div>
                    </div>
                    <div class="progress-steps">
                        ${this.verificationSteps.map((step, index) => `
                            <div class="progress-step ${index <= this.currentStep ? 'active' : ''} ${index < this.currentStep ? 'completed' : ''}">
                                <div class="step-number">${index + 1}</div>
                                <div class="step-label">${this.getStepLabel(step)}</div>
                            </div>
                        `).join('')}
                    </div>
                </div>

                <div class="verification-content">
                    ${this.renderCurrentStep()}
                </div>

                <div class="verification-actions">
                    ${this.currentStep > 0 ? `
                        <button class="btn btn-secondary" onclick="profileVerification.previousStep()">
                            戻る
                        </button>
                    ` : ''}
                    <button class="btn btn-primary" onclick="profileVerification.nextStep()">
                        ${this.currentStep === this.verificationSteps.length - 1 ? '提出する' : '次へ'}
                    </button>
                </div>
            </div>
        `;
    }

    renderCurrentStep() {
        const step = this.verificationSteps[this.currentStep];
        
        switch (step) {
            case 'photo_id':
                return this.renderPhotoIdStep();
            case 'selfie':
                return this.renderSelfieStep();
            case 'profile_info':
                return this.renderProfileInfoStep();
            default:
                return '';
        }
    }

    renderPhotoIdStep() {
        return `
            <div class="step-content">
                <h3>身分証明書のアップロード</h3>
                <p>以下のいずれかの身分証明書の写真をアップロードしてください：</p>
                
                <ul class="id-types">
                    <li>運転免許証</li>
                    <li>パスポート</li>
                    <li>マイナンバーカード（表面のみ）</li>
                    <li>在留カード</li>
                </ul>

                <div class="upload-area" id="idUploadArea">
                    <input type="file" id="idUpload" accept="image/*" hidden>
                    <div class="upload-placeholder" onclick="document.getElementById('idUpload').click()">
                        <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
                            <circle cx="8.5" cy="8.5" r="1.5"/>
                            <polyline points="21 15 16 10 5 21"/>
                        </svg>
                        <p>クリックして画像を選択</p>
                        <small>JPG, PNG形式（最大5MB）</small>
                    </div>
                    <div class="upload-preview hidden" id="idPreview">
                        <img src="" alt="ID preview">
                        <button class="remove-btn" onclick="profileVerification.removeUpload('id')">×</button>
                    </div>
                </div>

                <div class="security-notice">
                    <p>🔒 アップロードされた画像は暗号化され、安全に保管されます</p>
                </div>
            </div>
        `;
    }

    renderSelfieStep() {
        return `
            <div class="step-content">
                <h3>自撮り写真の撮影</h3>
                <p>本人確認のため、現在の自撮り写真を撮影してください</p>

                <div class="selfie-instructions">
                    <div class="instruction-item">
                        <span class="icon">💡</span>
                        <p>明るい場所で撮影してください</p>
                    </div>
                    <div class="instruction-item">
                        <span class="icon">😊</span>
                        <p>正面を向いて、顔全体が写るようにしてください</p>
                    </div>
                    <div class="instruction-item">
                        <span class="icon">🚫</span>
                        <p>サングラスやマスクは外してください</p>
                    </div>
                </div>

                <div class="camera-container">
                    <video id="selfieVideo" autoplay playsinline></video>
                    <canvas id="selfieCanvas" hidden></canvas>
                    <div class="camera-controls">
                        <button class="btn btn-primary" id="captureBtn" onclick="profileVerification.captureSelfie()">
                            撮影する
                        </button>
                    </div>
                </div>

                <div class="selfie-preview hidden" id="selfiePreview">
                    <img src="" alt="Selfie preview">
                    <button class="btn btn-secondary" onclick="profileVerification.retakeSelfie()">
                        撮り直す
                    </button>
                </div>
            </div>
        `;
    }

    renderProfileInfoStep() {
        return `
            <div class="step-content">
                <h3>プロフィール情報の確認</h3>
                <p>以下の情報が正しいことを確認してください</p>

                <form class="verification-form">
                    <div class="form-group">
                        <label>氏名（本名）</label>
                        <input type="text" class="form-input" placeholder="山田 太郎" required>
                        <small>身分証明書と同じ氏名を入力してください</small>
                    </div>

                    <div class="form-group">
                        <label>生年月日</label>
                        <input type="date" class="form-input" required>
                    </div>

                    <div class="form-group">
                        <label>性別</label>
                        <select class="form-select" required>
                            <option value="">選択してください</option>
                            <option value="male">男性</option>
                            <option value="female">女性</option>
                            <option value="other">その他</option>
                        </select>
                    </div>

                    <div class="form-group">
                        <label>
                            <input type="checkbox" required>
                            <span>提供された情報が正確であることを確認します</span>
                        </label>
                    </div>

                    <div class="form-group">
                        <label>
                            <input type="checkbox" required>
                            <span>利用規約とプライバシーポリシーに同意します</span>
                        </label>
                    </div>
                </form>

                <div class="verification-benefits">
                    <h4>認証完了後のメリット</h4>
                    <ul>
                        <li>✓ プロフィールに認証バッジが表示されます</li>
                        <li>✓ マッチング率が向上します</li>
                        <li>✓ 信頼性の高いユーザーとして認識されます</li>
                    </ul>
                </div>
            </div>
        `;
    }

    getStepLabel(step) {
        const labels = {
            'photo_id': '身分証明書',
            'selfie': '自撮り写真',
            'profile_info': '情報確認'
        };
        return labels[step] || '';
    }

    nextStep() {
        if (this.validateCurrentStep()) {
            if (this.currentStep < this.verificationSteps.length - 1) {
                this.currentStep++;
                this.updateUI();
            } else {
                this.submitVerification();
            }
        }
    }

    previousStep() {
        if (this.currentStep > 0) {
            this.currentStep--;
            this.updateUI();
        }
    }

    validateCurrentStep() {
        const step = this.verificationSteps[this.currentStep];
        
        switch (step) {
            case 'photo_id':
                if (!this.verificationData.photoId) {
                    showNotification('身分証明書の画像をアップロードしてください', 'error');
                    return false;
                }
                return true;
                
            case 'selfie':
                if (!this.verificationData.selfie) {
                    showNotification('自撮り写真を撮影してください', 'error');
                    return false;
                }
                return true;
                
            case 'profile_info':
                const form = document.querySelector('.verification-form');
                if (!form.checkValidity()) {
                    showNotification('すべての必須項目を入力してください', 'error');
                    return false;
                }
                return true;
                
            default:
                return true;
        }
    }

    async captureSelfie() {
        const video = document.getElementById('selfieVideo');
        const canvas = document.getElementById('selfieCanvas');
        const context = canvas.getContext('2d');
        
        canvas.width = video.videoWidth;
        canvas.height = video.videoHeight;
        context.drawImage(video, 0, 0);
        
        const dataURL = canvas.toDataURL('image/jpeg');
        this.verificationData.selfie = dataURL;
        
        // Show preview
        document.getElementById('selfiePreview').classList.remove('hidden');
        document.querySelector('#selfiePreview img').src = dataURL;
        document.querySelector('.camera-container').classList.add('hidden');
    }

    retakeSelfie() {
        document.getElementById('selfiePreview').classList.add('hidden');
        document.querySelector('.camera-container').classList.remove('hidden');
        this.verificationData.selfie = null;
    }

    async initCamera() {
        try {
            const stream = await navigator.mediaDevices.getUserMedia({ 
                video: { facingMode: 'user' } 
            });
            const video = document.getElementById('selfieVideo');
            if (video) {
                video.srcObject = stream;
            }
        } catch (err) {
            console.error('Camera access denied:', err);
            showNotification('カメラへのアクセスが拒否されました', 'error');
        }
    }

    removeUpload(type) {
        if (type === 'id') {
            this.verificationData.photoId = null;
            document.getElementById('idPreview').classList.add('hidden');
            document.getElementById('idUpload').value = '';
        }
    }

    updateUI() {
        const container = document.querySelector('.verification-container');
        if (container) {
            container.innerHTML = this.render();
            
            // Initialize camera if on selfie step
            if (this.verificationSteps[this.currentStep] === 'selfie') {
                this.initCamera();
            }
            
            // Add file upload listener
            const idUpload = document.getElementById('idUpload');
            if (idUpload) {
                idUpload.addEventListener('change', (e) => this.handleFileUpload(e));
            }
        }
    }

    handleFileUpload(event) {
        const file = event.target.files[0];
        if (file) {
            if (file.size > 5 * 1024 * 1024) {
                showNotification('ファイルサイズは5MB以下にしてください', 'error');
                return;
            }
            
            const reader = new FileReader();
            reader.onload = (e) => {
                this.verificationData.photoId = e.target.result;
                document.getElementById('idPreview').classList.remove('hidden');
                document.querySelector('#idPreview img').src = e.target.result;
                document.querySelector('.upload-placeholder').style.display = 'none';
            };
            reader.readAsDataURL(file);
        }
    }

    async submitVerification() {
        // Collect form data
        const form = document.querySelector('.verification-form');
        const formData = new FormData(form);
        
        this.verificationData.fullName = formData.get('fullName');
        this.verificationData.birthDate = formData.get('birthDate');
        this.verificationData.gender = formData.get('gender');
        
        // Show loading state
        showNotification('認証情報を送信中...', 'info');
        
        try {
            // In real app, send to backend
            const response = await apiCall('/api/verification/submit', {
                method: 'POST',
                body: JSON.stringify(this.verificationData)
            });
            
            showNotification('認証申請が完了しました！審査結果は24時間以内にお知らせします', 'success');
            
            // Redirect or update UI
            setTimeout(() => {
                window.location.href = '/profile';
            }, 2000);
            
        } catch (error) {
            showNotification('認証の送信に失敗しました。もう一度お試しください', 'error');
        }
    }
}

// Global instance
const profileVerification = new ProfileVerification();