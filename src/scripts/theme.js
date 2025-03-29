/**
 * Safe wrapper for localStorage to handle environments where it's not available
 * @typedef {Object} SafeStorage
 * @property {function(string): (string|null)} getItem - Gets item from storage
 * @property {function(string, string): void} setItem - Sets item in storage
 */

/** @type {SafeStorage} */
const safeStorage = {
  /**
   * Safely get item from localStorage with fallback
   * @param {string} key - The key to retrieve from storage
   * @returns {string|null} The stored value or null if not found/unavailable
   */
  getItem(key) {
    try {
      return localStorage.getItem(key);
    } catch (err) {
      console.warn("localStorage not available:", err);
      return null;
    }
  },

  /**
   * Safely set item in localStorage with fallback
   * @param {string} key - The key to store
   * @param {string} value - The value to store
   * @returns {void}
   */
  setItem(key, value) {
    try {
      localStorage.setItem(key, value);
    } catch (err) {
      console.warn("localStorage not available:", err);
    }
  },
};

/**
 * @typedef {Object} Theme
 * @property {string} name - Theme name
 * @property {string} displayName - User-friendly theme name
 * @property {string} bg - Background color value
 * @property {string} fg - Foreground color value
 * @property {string} primary - Primary color for accents and highlights
 * @property {string} primaryRgb - Primary color in RGB format (r, g, b)
 * @property {string} secondary - Secondary color for subtle elements
 * @property {string} accent - Accent color for emphasis
 * @property {string} trackColor - Color of theme change track
 * @property {string} ballColor - Color of theme change ball
 * @property {string} codeBg - Code block background
 * @property {string} codeFg - Code text color
 * @property {string} codeBorder - Code block border color
 * @property {string} linkHover - Link hover color
 */

/**
 * Available themes collection
 * @type {Object.<string, Theme>}
 */
const availableThemes = {
  light: {
    name: "light",
    displayName: "Light (High Contrast)",
    bg: "white",
    fg: "black",
    primary: "black",
    primaryRgb: "0, 0, 0",
    secondary: "#333",
    accent: "#666",
    trackColor: "#e0e0e0",
    ballColor: "#ffffff",
    codeBg: "#f8f8f8",
    codeFg: "#333333",
    codeBorder: "#ddd",
    linkHover: "black",
  },
  
  dark: {
    name: "dark",
    displayName: "Dark (High Contrast)",
    bg: "black",
    fg: "white",
    primary: "white",
    primaryRgb: "255, 255, 255",
    secondary: "#ddd",
    accent: "#aaa",
    trackColor: "#333333",
    ballColor: "#000000",
    codeBg: "#272822",
    codeFg: "#f8f8f2", 
    codeBorder: "#444",
    linkHover: "white",
  },
  
  nord: {
    name: "nord",
    displayName: "Nord",
    bg: "#2E3440",
    fg: "#ECEFF4",
    primary: "#88C0D0",
    primaryRgb: "136, 192, 208",
    secondary: "#81A1C1",
    accent: "#EBCB8B",
    trackColor: "#4C566A",
    ballColor: "#88C0D0",
    codeBg: "#3B4252",
    codeFg: "#E5E9F0",
    codeBorder: "#4C566A",
    linkHover: "#88C0D0",
  },
  
  solarized: {
    name: "solarized",
    displayName: "Solarized Light",
    bg: "#FDF6E3",
    fg: "#657B83",
    primary: "#268BD2",
    primaryRgb: "38, 139, 210",
    secondary: "#2AA198",
    accent: "#CB4B16",
    trackColor: "#EEE8D5",
    ballColor: "#268BD2",
    codeBg: "#EEE8D5",
    codeFg: "#586E75",
    codeBorder: "#D3CAA8",
    linkHover: "#268BD2",
  },
  
  dracula: {
    name: "dracula",
    displayName: "Dracula",
    bg: "#282A36",
    fg: "#F8F8F2",
    primary: "#BD93F9",
    primaryRgb: "189, 147, 249",
    secondary: "#6272A4",
    accent: "#FF79C6",
    trackColor: "#44475A",
    ballColor: "#BD93F9",
    codeBg: "#44475A",
    codeFg: "#F8F8F2",
    codeBorder: "#6272A4",
    linkHover: "#FF79C6",
  },
  
  github: {
    name: "github",
    displayName: "GitHub",
    bg: "#ffffff",
    fg: "#24292e",
    primary: "#0366d6",
    primaryRgb: "3, 102, 214",
    secondary: "#586069",
    accent: "#d73a49",
    trackColor: "#e1e4e8",
    ballColor: "#0366d6",
    codeBg: "#f6f8fa",
    codeFg: "#24292e",
    codeBorder: "#e1e4e8",
    linkHover: "#0366d6",
  },
  
  monokai: {
    name: "monokai",
    displayName: "Monokai",
    bg: "#272822",
    fg: "#F8F8F2",
    primary: "#F92672",
    primaryRgb: "249, 38, 114",
    secondary: "#66D9EF",
    accent: "#A6E22E",
    trackColor: "#49483E",
    ballColor: "#F92672",
    codeBg: "#3E3D32",
    codeFg: "#F8F8F2",
    codeBorder: "#75715E",
    linkHover: "#F92672",
  }
};

/**
 * Loads theme from localStorage and applies it if present
 * @returns {Theme} The loaded or default theme object
 */
function loadThemeFromLocalStorage() {
  let themeObject;
  const themeString = safeStorage.getItem("theme");
  
  if (themeString) {
    themeObject = JSON.parse(themeString);
    if (themeObject) {
      themeObject = setTheme(themeObject);
    }
  } else {
    themeObject = setTheme(availableThemes.light);
  }
  
    const themeSelector = document.getElementById("theme-selector");
    if (themeSelector) {
      themeSelector.value = themeObject.name;
    }
  } catch (err) {}
  
  return themeObject;
}

/**
 * Toggles between light and dark themes (for backward compatibility)
 * @returns {Theme} The new theme object after toggle
 */
function toggleTheme() {
  const themeString = safeStorage.getItem("theme");
  if (!themeString) {
    return setTheme(availableThemes.light);
  }
  
  const themeObject = JSON.parse(themeString);
  if (themeObject.name === "dark") {
    return setTheme(availableThemes.light);
  } else {
    return setTheme(availableThemes.dark);
  }
}

/**
 * Sets theme by name
 * @param {string} themeName - Name of the theme to apply
 * @returns {Theme} The applied theme object
 */
function setThemeByName(themeName) {
  if (availableThemes[themeName]) {
    return setTheme(availableThemes[themeName]);
  } else {
    console.warn(`Theme "${themeName}" not found, using default light theme`);
    return setTheme(availableThemes.light);
  }
}

/**
 * Sets dark theme (for backward compatibility)
 * @returns {Theme} The dark theme object
 */
function setDarkTheme() {
  return setTheme(availableThemes.dark);
}

/**
 * Sets light theme (for backward compatibility)
 * @returns {Theme} The light theme object
 */
function setLightTheme() {
  return setTheme(availableThemes.light);
}

/**
 * Applies theme by setting CSS custom properties.
 * Note: setting the variables is not enough.
 * The stylesheet is loaded after this script is run.
 * Before the style sheet is loaded, the default background of
 * the html and body elements is the browser's default which can cause flickering.
 * @param {Theme} themeObject - The theme object containing color values
 * @returns {Theme} The applied theme object
 */
function setTheme(themeObject) {
  try {
    const body = document.querySelector("body");
    body.style.setProperty("background-color", "var(--bg)");
    body.style.setProperty("color", "var(--fg)");
  } catch (err) {}
  
  try {
    const navbar = document.querySelector("nav.navbar");
    if (navbar) {
      navbar.style.setProperty("--track-color", themeObject.trackColor);
      navbar.style.setProperty("--ball-color", themeObject.ballColor);
    }
  } catch (err) {}
  
  try {
    const root = document.querySelector(":root");
    root.setAttribute("data-theme", themeObject.name);
    
    root.style.setProperty("--bg", themeObject.bg);
    root.style.setProperty("--fg", themeObject.fg);
    root.style.setProperty("--primary", themeObject.primary);
    root.style.setProperty("--primary-rgb", themeObject.primaryRgb);
    root.style.setProperty("--secondary", themeObject.secondary);
    root.style.setProperty("--accent", themeObject.accent);
    root.style.setProperty("--code-bg", themeObject.codeBg);
    root.style.setProperty("--code-fg", themeObject.codeFg);
    root.style.setProperty("--code-border", themeObject.codeBorder);
    root.style.setProperty("--link-hover", themeObject.linkHover);
  } catch (err) {}
  
  safeStorage.setItem("theme", JSON.stringify(themeObject));
  
  return themeObject;
}

/**
 * Initializes the theme picker UI
 */
function initThemePicker() {
  const themeToggleContainer = document.querySelector('.theme-toggle');
  if (!themeToggleContainer) return;
  
  themeToggleContainer.innerHTML = '';
  
  const themeButton = document.createElement('button');
  themeButton.id = 'theme-button';
  themeButton.className = 'theme-button';
  themeButton.innerHTML = '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512" width="16" height="16" fill="currentColor"><path d="M512 256c0 141.4-114.6 256-256 256S0 397.4 0 256 114.6 0 256 0s256 114.6 256 256zM256 48C141.1 48 48 141.1 48 256s93.1 208 208 208 208-93.1 208-208S370.9 48 256 48zm0 384c-97.2 0-176-78.8-176-176S158.8 80 256 80s176 78.8 176 176-78.8 176-176 176zm-80-176c0 44.2 35.8 80 80 80s80-35.8 80-80-35.8-80-80-80-80 35.8-80 80z"/></svg>';
  themeButton.setAttribute('aria-label', 'Change theme');
  themeButton.setAttribute('title', 'Change theme');
  
  const themeMenu = document.createElement('div');
  themeMenu.id = 'theme-menu';
  themeMenu.className = 'theme-menu';
  
  const themeCount = Object.keys(availableThemes).length;
  if (themeCount > 4) {
    themeMenu.style.width = '200px';
  }
  
  for (const themeName in availableThemes) {
    const theme = availableThemes[themeName];
    
    const themeOption = document.createElement('div');
    themeOption.className = 'theme-option';
    themeOption.setAttribute('data-theme', themeName);
    themeOption.setAttribute('title', theme.displayName);
    themeOption.style.backgroundColor = theme.bg;
    
    if (themeName === 'light' || themeName === 'dark') {
      themeOption.style.background = `linear-gradient(135deg, ${theme.bg} 0%, ${theme.bg} 50%, ${theme.primary} 50%, ${theme.primary} 100%)`;
    } else {
      themeOption.style.background = `linear-gradient(135deg, 
        ${theme.bg} 0%, 
        ${theme.bg} 40%, 
        ${theme.primary} 40%, 
        ${theme.primary} 60%,
        ${theme.accent} 60%, 
        ${theme.accent} 100%)`;
    }
    
    const themeTitleSpan = document.createElement('span');
    themeTitleSpan.className = 'theme-option-name';
    themeTitleSpan.textContent = theme.displayName.split(' ')[0];
    themeOption.appendChild(themeTitleSpan);
    
    themeOption.addEventListener('click', function() {
      setThemeByName(themeName);
      updateActiveTheme(themeName);
      toggleThemeMenu();
    });
    
    themeMenu.appendChild(themeOption);
  }
  
  themeButton.addEventListener('click', toggleThemeMenu);
  
  document.addEventListener('click', function(event) {
    if (!themeToggleContainer.contains(event.target)) {
      themeMenu.classList.remove('open');
    }
  });
  
  themeToggleContainer.appendChild(themeButton);
  themeToggleContainer.appendChild(themeMenu);
  
  const currentTheme = safeStorage.getItem("theme");
  if (currentTheme) {
    const parsedTheme = JSON.parse(currentTheme);
    updateActiveTheme(parsedTheme.name);
  }
}

/**
 * Updates the active theme in the UI
 * @param {string} themeName - The name of the active theme
 */
function updateActiveTheme(themeName) {
  const themeOptions = document.querySelectorAll('.theme-option');
  themeOptions.forEach(option => {
    if (option.getAttribute('data-theme') === themeName) {
      option.classList.add('active');
    } else {
      option.classList.remove('active');
    }
  });
  
  const themeButton = document.getElementById('theme-button');
  if (themeButton && availableThemes[themeName]) {
    themeButton.style.backgroundColor = availableThemes[themeName].primary;
    themeButton.style.borderColor = availableThemes[themeName].primary;
    themeButton.style.color = availableThemes[themeName].bg;
  }
}

/**
 * Toggles the theme menu open/closed
 */
function toggleThemeMenu() {
  const themeMenu = document.getElementById('theme-menu');
  if (themeMenu) {
    themeMenu.classList.toggle('open');
  }
}

document.addEventListener('DOMContentLoaded', function() {
  initThemePicker();
  loadThemeFromLocalStorage();
});
