const buildTOCTree = (tocItems) => {
    const tree = [];
    const stack = [];
    
    tocItems.forEach(item => {
        const level = getItemLevel(item);
        const link = item.querySelector('a');
        const node = {
            element: item,
            level: level,
            text: link ? link.textContent : '',
            href: link ? link.getAttribute('href') : '',
            children: [],
            hasChildren: false
        };
        
        while (stack.length > 0 && stack[stack.length - 1].level >= level) {
            stack.pop();
        }
        
        if (stack.length === 0) {
            tree.push(node);
        } else {
            stack[stack.length - 1].children.push(node);
            stack[stack.length - 1].hasChildren = true;
        }
        
        stack.push(node);
    });
    
    return tree;
};

const renderTOCTree = (tree, container) => {
    const ul = document.createElement('ul');
    ul.className = 'toc-tree';
    
    tree.forEach(node => {
        const li = document.createElement('li');
        li.className = `toc-level-${node.level}`;
        
        const link = document.createElement('a');
        link.href = node.href;
        link.textContent = node.text;
        li.appendChild(link);
        
        if (node.hasChildren) {
            const toggleBtn = document.createElement('button');
            toggleBtn.className = 'toc-toggle-btn';
            toggleBtn.innerHTML = '<i class="fas fa-chevron-right"></i>';
            toggleBtn.setAttribute('aria-expanded', 'false');
            li.appendChild(toggleBtn);
        }
        
        if (node.children.length > 0) {
            const childContainer = document.createElement('div');
            childContainer.className = 'toc-children';
            childContainer.style.display = 'none';
            renderTOCTree(node.children, childContainer);
            li.appendChild(childContainer);
        }
        
        ul.appendChild(li);
    });
    
    container.appendChild(ul);
};

const getItemLevel = (item) => {
    for (let i = 1; i <= 6; i++) {
        if (item.classList.contains(`toc-level-${i}`)) {
            return i;
        }
    }
    return 1;
};

const initTOC = () => {
    const tocItems = document.querySelectorAll('.table-of-contents li');
    const headings = document.querySelectorAll('.article h1[id], .article h2[id], .article h3[id], .article h4[id]');
    const tocList = document.getElementById('toc-list');
    
    if (!tocItems.length || !headings.length) return;
    
    const tree = buildTOCTree(tocItems);
    tocList.innerHTML = '';
    renderTOCTree(tree, tocList);

    const updateActiveLink = () => {
        let current = '';

        headings.forEach(heading => {
            const rect = heading.getBoundingClientRect();
            if (rect.top <= 100) {
                current = heading.id;
            }
        });

        const tocLinks = document.querySelectorAll('.table-of-contents a');
        tocLinks.forEach(link => {
            const href = link.getAttribute('href');
            const linkId = href ? href.substring(1) : '';
            
            link.classList.remove('active');
            if (linkId === current) {
                link.classList.add('active');
                
                expandParentsIfHidden(link);
            }
        });
    };

    const expandParentsIfHidden = (activeLink) => {
        const li = activeLink.closest('li');
        if (!li) return;

        const isVisible = li.offsetParent !== null;
        if (isVisible) return;

        let current = li;

        while (current) {
            const parentContainer = current.closest('.toc-children');
            if (!parentContainer) break;

            const parentLi = parentContainer.previousElementSibling?.closest('li') || 
                             parentContainer.parentElement?.closest('li');
            
            if (parentLi) {
                const toggleBtn = parentLi.querySelector('.toc-toggle-btn');
                
                if (toggleBtn) {
                    const isExpanded = toggleBtn.getAttribute('aria-expanded') === 'true';
                    
                    if (!isExpanded) {
                        parentContainer.style.display = 'block';
                        const icon = toggleBtn.querySelector('i');
                        if (icon) {
                            icon.className = 'fas fa-chevron-down';
                        }
                        toggleBtn.setAttribute('aria-expanded', 'true');
                    }
                }
            }

            current = parentContainer.parentElement?.closest('li');
        }
    };


    const setupToggleListeners = () => {
        const toggleButtons = document.querySelectorAll('.toc-toggle-btn');
        
        toggleButtons.forEach(button => {
            button.addEventListener('click', (e) => {
                e.preventDefault();
                e.stopPropagation();
                
                const li = button.parentElement;
                const childContainer = li.querySelector('.toc-children');
                const icon = button.querySelector('i');
                const isExpanded = button.getAttribute('aria-expanded') === 'true';
                
                if (childContainer) {
                    if (isExpanded) {
                        childContainer.style.display = 'none';
                        icon.className = 'fas fa-chevron-right';
                        button.setAttribute('aria-expanded', 'false');
                    } else {
                        childContainer.style.display = 'block';
                        icon.className = 'fas fa-chevron-down';
                        button.setAttribute('aria-expanded', 'true');
                    }
                }
            });
        });
    };

    const setupLinkListeners = () => {
        const tocLinks = document.querySelectorAll('.table-of-contents a');
        
        tocLinks.forEach(link => {
            link.addEventListener('click', (e) => {
                e.preventDefault();
                
                const targetId = link.getAttribute('href').substring(1);
                const targetElement = document.getElementById(targetId);
                
                if (targetElement) {
                    const navbar = document.querySelector('nav.navbar');
                    const navbarHeight = navbar ? navbar.offsetHeight : 0;
                    const elementPosition = targetElement.getBoundingClientRect().top;
                    const offsetPosition = elementPosition + window.pageYOffset - navbarHeight - 20;
                    
                    window.scrollTo({
                        top: offsetPosition,
                        behavior: 'smooth'
                    });
                }
            });
        });
    };

    setupToggleListeners();
    setupLinkListeners();

    window.addEventListener('scroll', updateActiveLink, { passive: true });
    updateActiveLink();
};

if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initTOC);
} else {
    initTOC();
}