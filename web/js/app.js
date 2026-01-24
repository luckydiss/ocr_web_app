const App = {
    elements: {},
    currentMarkdown: '',
    selectedFile: null,

    init() {
        this.cacheElements();
        this.bindEvents();
        this.setupMainButton();
    },

    cacheElements() {
        this.elements = {
            imageInput: document.getElementById('imageInput'),
            selectBtn: document.getElementById('selectBtn'),
            preview: document.getElementById('preview'),
            previewImage: document.getElementById('previewImage'),
            clearBtn: document.getElementById('clearBtn'),
            extractBtn: document.getElementById('extractBtn'),
            resultSection: document.getElementById('resultSection'),
            renderedView: document.getElementById('renderedView'),
            rawView: document.getElementById('rawView'),
            renderedMarkdown: document.getElementById('renderedMarkdown'),
            rawMarkdown: document.getElementById('rawMarkdown'),
            copyBtn: document.getElementById('copyBtn'),
            loading: document.getElementById('loading'),
            error: document.getElementById('error'),
            tabs: document.querySelectorAll('.tab')
        };
    },

    bindEvents() {
        this.elements.selectBtn.addEventListener('click', () => {
            this.elements.imageInput.click();
        });

        this.elements.imageInput.addEventListener('change', (e) => {
            this.handleFileSelect(e.target.files[0]);
        });

        this.elements.clearBtn.addEventListener('click', () => {
            this.clearImage();
        });

        this.elements.extractBtn.addEventListener('click', () => {
            this.extractText();
        });

        this.elements.copyBtn.addEventListener('click', () => {
            this.copyToClipboard();
        });

        this.elements.tabs.forEach(tab => {
            tab.addEventListener('click', () => {
                this.switchTab(tab.dataset.view);
            });
        });
    },

    setupMainButton() {
        if (TelegramApp.isAvailable) {
            TelegramApp.hideMainButton();
        }
    },

    handleFileSelect(file) {
        if (!file) return;

        const validTypes = ['image/jpeg', 'image/png', 'image/gif', 'image/webp'];
        if (!validTypes.includes(file.type)) {
            this.showError('Please select a valid image (JPEG, PNG, GIF, or WebP)');
            return;
        }

        if (file.size > 10 * 1024 * 1024) {
            this.showError('Image is too large. Maximum size is 10MB.');
            return;
        }

        this.selectedFile = file;
        this.hideError();

        const reader = new FileReader();
        reader.onload = (e) => {
            this.elements.previewImage.src = e.target.result;
            this.elements.preview.classList.remove('hidden');
        };
        reader.readAsDataURL(file);

        if (TelegramApp.isAvailable) {
            TelegramApp.showMainButton('Extract Text', () => this.extractText());
        }

        this.elements.extractBtn.disabled = false;

        TelegramApp.hapticFeedback('impact');
    },

    clearImage() {
        this.selectedFile = null;
        this.elements.imageInput.value = '';
        this.elements.preview.classList.add('hidden');
        this.elements.resultSection.classList.add('hidden');
        this.elements.extractBtn.disabled = true;
        this.currentMarkdown = '';

        TelegramApp.hideMainButton();

        TelegramApp.hapticFeedback('impact');
    },

    async extractText() {
        if (!this.selectedFile) {
            this.showError('Please select an image first');
            return;
        }

        this.showLoading();
        this.hideError();
        this.elements.resultSection.classList.add('hidden');

        TelegramApp.setMainButtonLoading(true);

        try {
            const base64 = await this.fileToBase64(this.selectedFile);
            const base64Data = base64.split(',')[1];

            const response = await fetch('/api/v1/extract', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    image_base64: base64Data,
                    mime_type: this.selectedFile.type
                })
            });

            const data = await response.json();

            if (!data.success) {
                throw new Error(data.error || 'Failed to extract text');
            }

            this.currentMarkdown = data.payload.markdown;
            this.displayResult();

            TelegramApp.hapticFeedback('success');

        } catch (error) {
            console.error('Extraction error:', error);
            this.showError(error.message || 'Failed to extract text from image');
            TelegramApp.hapticFeedback('error');
        } finally {
            this.hideLoading();
            TelegramApp.setMainButtonLoading(false);
        }
    },

    fileToBase64(file) {
        return new Promise((resolve, reject) => {
            const reader = new FileReader();
            reader.onload = () => resolve(reader.result);
            reader.onerror = reject;
            reader.readAsDataURL(file);
        });
    },

    displayResult() {
        this.elements.rawMarkdown.textContent = this.currentMarkdown;
        MarkdownRenderer.render(this.currentMarkdown, this.elements.renderedMarkdown);
        this.elements.resultSection.classList.remove('hidden');
        this.switchTab('rendered');
    },

    switchTab(view) {
        this.elements.tabs.forEach(tab => {
            tab.classList.toggle('active', tab.dataset.view === view);
        });

        if (view === 'rendered') {
            this.elements.renderedView.classList.remove('hidden');
            this.elements.rawView.classList.add('hidden');
        } else {
            this.elements.renderedView.classList.add('hidden');
            this.elements.rawView.classList.remove('hidden');
        }
    },

    async copyToClipboard() {
        try {
            await navigator.clipboard.writeText(this.currentMarkdown);

            const originalText = this.elements.copyBtn.innerHTML;
            this.elements.copyBtn.innerHTML = '<span class="icon">✓</span> Copied!';

            TelegramApp.hapticFeedback('success');

            setTimeout(() => {
                this.elements.copyBtn.innerHTML = originalText;
            }, 2000);

        } catch (error) {
            const textarea = document.createElement('textarea');
            textarea.value = this.currentMarkdown;
            document.body.appendChild(textarea);
            textarea.select();
            document.execCommand('copy');
            document.body.removeChild(textarea);

            TelegramApp.showAlert('Copied to clipboard!');
        }
    },

    showLoading() {
        this.elements.loading.classList.remove('hidden');
    },

    hideLoading() {
        this.elements.loading.classList.add('hidden');
    },

    showError(message) {
        this.elements.error.textContent = message;
        this.elements.error.classList.remove('hidden');
    },

    hideError() {
        this.elements.error.classList.add('hidden');
    }
};

document.addEventListener('DOMContentLoaded', () => {
    App.init();
});
