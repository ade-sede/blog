const initTOC = () => {
    const tocLinks = document.querySelectorAll('.table-of-contents a');
    const tocItems = document.querySelectorAll('.table-of-contents li');
    const headings = document.querySelectorAll('.article h1[id], .article h2[id], .article h3[id], .article h4[id]');
    const tocToggle = document.getElementById('toc-toggle');
    const tocList = document.getElementById('toc-list');
    const tocIcon = document.getElementById('toc-toggle-icon');
    
    if (!tocLinks.length || !headings.length) return;

    const TOC_STORAGE_KEY = 'toc-visible';
    
    const getTocVisibility = () => {
        const stored = localStorage.getItem(TOC_STORAGE_KEY);
        return stored !== null ? stored === 'true' : true;
    };
    
    const setTocVisibility = (visible) => {
        localStorage.setItem(TOC_STORAGE_KEY, visible);
        if (visible) {
            tocList.style.display = 'block';
            tocIcon.className = 'fas fa-eye-slash';
        } else {
            tocList.style.display = 'none';
            tocIcon.className = 'fas fa-eye';
        }
    };
    
    setTocVisibility(getTocVisibility());
    
    if (tocToggle) {
        tocToggle.addEventListener('click', () => {
            const isVisible = tocList.style.display !== 'none';
            setTocVisibility(!isVisible);
        });
    }

    tocLinks.forEach(link => {
        link.addEventListener('click', (e) => {
            e.preventDefault();
            
            const targetId = link.getAttribute('href').substring(1);
            const targetElement = document.getElementById(targetId);
            
            if (targetElement) {
                targetElement.scrollIntoView({
                    behavior: 'smooth',
                    block: 'start'
                });
            }
        });
    });

    const getItemLevel = (item) => {
        for (let i = 1; i <= 6; i++) {
            if (item.classList.contains(`toc-level-${i}`)) {
                return i;
            }
        }
        return 1;
    };

    const updateActiveLink = () => {
        let current = '';

        headings.forEach(heading => {
            const rect = heading.getBoundingClientRect();
            if (rect.top <= 100) {
                current = heading.id;
            }
        });

        tocItems.forEach(item => {
            const link = item.querySelector('a');
            if (!link) return;
            
            const href = link.getAttribute('href');
            const linkId = href ? href.substring(1) : '';
            
            link.classList.remove('active');
            if (linkId === current) {
                link.classList.add('active');
            }
            
            item.style.display = 'block';
        });
    };

    window.addEventListener('scroll', updateActiveLink, { passive: true });
    updateActiveLink();
};

if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initTOC);
} else {
    initTOC();
}