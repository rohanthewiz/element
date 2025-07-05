/* Debug Table */

function switchTab(tabName) {
    // Remove active class from all tabs and contents
    const tabs = document.querySelectorAll('.tab');
    const contents = document.querySelectorAll('.tab-content');

    tabs.forEach(tab => tab.classList.remove('active'));
    contents.forEach(content => content.classList.remove('active'));

    // Add active class to selected tab and content
    document.getElementById(tabName + '-tab').classList.add('active');
    document.getElementById(tabName + '-content').classList.add('active');
}

function copyToClipboard(text) {
    navigator.clipboard.writeText(text).then(function() {
        showNotification('Copied: ' + text);
    }).catch(function(err) {
        console.error('Failed to copy: ', err);
        // Fallback for older browsers
        const textArea = document.createElement("textarea");
        textArea.value = text;
        textArea.style.position = "fixed";
        textArea.style.left = "-999999px";
        document.body.appendChild(textArea);
        textArea.focus();
        textArea.select();
        try {
            document.execCommand('copy');
            showNotification('Copied: ' + text);
        } catch (err) {
            console.error('Fallback copy failed: ', err);
        }
        document.body.removeChild(textArea);
    });
}

function showNotification(message) {
    const notification = document.createElement('div');
    notification.className = 'notification';
    notification.textContent = message;
    document.body.appendChild(notification);

    // Trigger reflow to enable transition
    notification.offsetHeight;
    notification.classList.add('show');

    setTimeout(function() {
        notification.classList.remove('show');
        setTimeout(function() {
            document.body.removeChild(notification);
        }, 300);
    }, 2000);
}

function clearIssues() {
    fetch('/debug/clear-issues')
        .then(response => {
            if (response.ok) {
                showNotification('Issues cleared successfully');
                setTimeout(function() {
                    window.location.reload();
                }, 1000);
            } else {
                showNotification('Failed to clear issues');
            }
        })
        .catch(error => {
            console.error('Error clearing issues:', error);
            showNotification('Error clearing issues');
        });
}

function copyMarkdownContent() {
    const markdownContent = document.getElementById('markdown-content-pre').textContent;
    copyToClipboard(markdownContent);
}
