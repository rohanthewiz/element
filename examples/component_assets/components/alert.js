// Alert Component JavaScript - Dismissible functionality
(function() {
    'use strict';

    // Handle alert dismissal
    function initAlerts() {
        const closeButtons = document.querySelectorAll('.ele-alert-close');

        closeButtons.forEach(function(button) {
            button.addEventListener('click', function(e) {
                e.preventDefault();
                const alert = this.closest('.ele-alert');

                // Add fade out animation
                alert.style.transition = 'opacity 0.3s ease-out';
                alert.style.opacity = '0';

                setTimeout(function() {
                    alert.classList.add('dismissed');
                }, 300);
            });
        });
    }

    // Initialize on DOM ready
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', initAlerts);
    } else {
        initAlerts();
    }
})();
