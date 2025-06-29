let footnoteRefs = [];
let footnoteItems = [];
const sidebarTop = 80;

document.addEventListener("DOMContentLoaded", function () {
  initSidebarFootnotes();
});

function initSidebarFootnotes() {
  if (window.innerWidth <= 1200) {
    return;
  }

  const footnotesSection = document.querySelector(".footnotes.desktop-hidden");
  const footnotesSidebar = document.getElementById("footnotes-sidebar");

  if (!footnotesSection || !footnotesSidebar) {
    return;
  }

  const footnotesList = footnotesSection.querySelector("ol");
  if (!footnotesList) {
    return;
  }

  createSidebarFootnotes(footnotesList, footnotesSidebar);
  setupHoverBehavior();

  window.addEventListener(
    "resize",
    debounce(function () {
      if (window.innerWidth <= 1200) {
        footnotesSidebar.innerHTML = "";
        footnotesSection.classList.remove("desktop-hidden");
      } else {
        footnotesSection.classList.add("desktop-hidden");
        createSidebarFootnotes(footnotesList, footnotesSidebar);
        setupHoverBehavior();
      }
    }, 150),
  );

  let isScrolling = false;
  window.addEventListener("scroll", () => {
    if (!isScrolling) {
      requestAnimationFrame(() => {
        updateFootnotePositions();
        isScrolling = false;
      });
      isScrolling = true;
    }
  });
}

function createSidebarFootnotes(footnotesList, sidebar) {
  sidebar.innerHTML = "";

  const footnoteItems = footnotesList.querySelectorAll("li");

  footnoteItems.forEach((item, index) => {
    const footnoteNumber = index + 1;
    const footnoteContent = item.innerHTML;

    const sidebarItem = document.createElement("div");
    sidebarItem.className = "footnote-item";
    sidebarItem.setAttribute("data-footnote", footnoteNumber);
    sidebarItem.innerHTML = `<span class="footnote-number">${footnoteNumber}</span>${footnoteContent}`;

    sidebar.appendChild(sidebarItem);
  });

  updateFootnotePositions();
}

function setupHoverBehavior() {
  footnoteRefs = document.querySelectorAll(".footnote-ref");
  footnoteItems = document.querySelectorAll(".footnote-item");

  footnoteRefs.forEach((ref, index) => {
    const correspondingFootnote = footnoteItems[index];
    if (!correspondingFootnote) return;

    ref.addEventListener("mouseenter", () => {
      correspondingFootnote.classList.add("show");
    });

    ref.addEventListener("mouseleave", () => {
      correspondingFootnote.classList.remove("show");
    });
  });
}

function updateFootnotePositions() {
  if (footnoteRefs.length === 0 || footnoteItems.length === 0) {
    return;
  }

  footnoteRefs.forEach((ref, index) => {
    const footnoteItem = footnoteItems[index];
    if (!footnoteItem) return;

    const rect = ref.getBoundingClientRect();
    const targetTop = rect.top - sidebarTop - 10;

    footnoteItem.style.top = `${targetTop}px`;
  });
}

function debounce(func, wait) {
  let timeout;
  return function executedFunction(...args) {
    const later = () => {
      clearTimeout(timeout);
      func(...args);
    };
    clearTimeout(timeout);
    timeout = setTimeout(later, wait);
  };
}
