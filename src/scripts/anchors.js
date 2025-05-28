const initAnchors = () => {
    const headerAnchors = document.querySelectorAll('.header-anchor');
    
    if (!headerAnchors.length) return;

    const showCopyFeedback = (anchor) => {
        const originalIcon = anchor.innerHTML;
        anchor.innerHTML = '<i class="fas fa-check"></i>';
        anchor.style.color = 'var(--success, #28a745)';
        
        setTimeout(() => {
            anchor.innerHTML = originalIcon;
            anchor.style.color = '';
        }, 1500);
    };

    headerAnchors.forEach(anchor => {
        anchor.addEventListener('click', async (e) => {
            e.preventDefault();
            e.stopPropagation();
            
            const href = anchor.getAttribute('href');
            if (!href) return;
            
            const fullUrl = window.location.origin + window.location.pathname + href;
            
            try {
                await navigator.clipboard.writeText(fullUrl);
                showCopyFeedback(anchor);
            } catch (err) {
                try {
                    const textArea = document.createElement('textarea');
                    textArea.value = fullUrl;
                    textArea.style.position = 'fixed';
                    textArea.style.left = '-999999px';
                    textArea.style.top = '-999999px';
                    document.body.appendChild(textArea);
                    textArea.focus();
                    textArea.select();
                    document.execCommand('copy');
                    textArea.remove();
                    showCopyFeedback(anchor);
                } catch (fallbackErr) {
                    console.error('Failed to copy URL:', fallbackErr);
                }
            }
        });
    });
};

if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initAnchors);
} else {
    initAnchors();
}