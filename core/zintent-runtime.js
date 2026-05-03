/**
 * ZINTENT Runtime Helper v2.1.0
 * Provides lightweight theme management and interaction hooks.
 */
(function(window) {
    const zintent = {
        /**
         * Set the current theme
         * @param {string} themeName - e.g., 'midnight', 'forest', 'light', 'nordic'
         * @param {boolean} persist - whether to save to localStorage
         */
        setTheme(themeName, persist = true) {
            document.documentElement.setAttribute('data-theme', themeName);
            if (persist) {
                localStorage.setItem('zi-theme', themeName);
            }
            // Dispatch event for components to react
            window.dispatchEvent(new CustomEvent('zi-theme-change', { detail: { theme: themeName } }));
        },

        /**
         * Initialize theme from storage or system preference
         */
        initTheme() {
            const saved = localStorage.getItem('zi-theme');
            if (saved) {
                this.setTheme(saved, false);
            } else {
                // Check system preference
                const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
                // Default to midnight if dark, light if light
                this.setTheme(prefersDark ? 'midnight' : 'light', false);
            }
        },

        /**
         * Get current active theme
         */
        getTheme() {
            return document.documentElement.getAttribute('data-theme') || 'light';
        },

        /**
         * Cycle through a list of themes
         * @param {string[]} themeList 
         */
        cycleThemes(themeList = ['light', 'midnight', 'nordic', 'forest', 'vibrant']) {
            const current = this.getTheme();
            let nextIndex = (themeList.indexOf(current) + 1) % themeList.length;
            if (nextIndex < 0) nextIndex = 0;
            this.setTheme(themeList[nextIndex]);
        }
    };

    // Export to window
    window.zintent = zintent;

    // Auto-init on DOM load
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', () => zintent.initTheme());
    } else {
        zintent.initTheme();
    }

})(window);
