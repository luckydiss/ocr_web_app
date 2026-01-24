const TelegramApp = {
    webapp: null,
    isAvailable: false,

    init() {
        if (window.Telegram && window.Telegram.WebApp) {
            this.webapp = window.Telegram.WebApp;
            this.isAvailable = true;
            this.webapp.ready();
            this.webapp.expand();
            this.applyTheme();
        }
    },

    applyTheme() {
        if (!this.isAvailable) return;
        document.body.style.backgroundColor = this.webapp.backgroundColor;
    },

    showMainButton(text, callback) {
        if (!this.isAvailable) return;

        const mainButton = this.webapp.MainButton;
        mainButton.setText(text);
        mainButton.onClick(callback);
        mainButton.show();
    },

    hideMainButton() {
        if (!this.isAvailable) return;
        this.webapp.MainButton.hide();
    },

    setMainButtonLoading(loading) {
        if (!this.isAvailable) return;

        if (loading) {
            this.webapp.MainButton.showProgress();
        } else {
            this.webapp.MainButton.hideProgress();
        }
    },

    hapticFeedback(type = 'impact') {
        if (!this.isAvailable || !this.webapp.HapticFeedback) return;

        switch (type) {
            case 'impact':
                this.webapp.HapticFeedback.impactOccurred('medium');
                break;
            case 'success':
                this.webapp.HapticFeedback.notificationOccurred('success');
                break;
            case 'error':
                this.webapp.HapticFeedback.notificationOccurred('error');
                break;
        }
    },

    showAlert(message) {
        if (this.isAvailable) {
            this.webapp.showAlert(message);
        } else {
            alert(message);
        }
    }
};

TelegramApp.init();
