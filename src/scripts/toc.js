const initTOC = () => {
    const tocLinks = document.querySelectorAll('.table-of-contents a');
    const tocItems = document.querySelectorAll('.table-of-contents li');
    const headings = document.querySelectorAll('.article h1[id], .article h2[id], .article h3[id], .article h4[id]');
    
    if (!tocLinks.length || !headings.length) return;

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