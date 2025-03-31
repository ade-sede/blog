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

  atom: {
    name: "atom",
    displayName: "Atom Dark",
    bg: "#282c34",
    fg: "#abb2bf",
    primary: "#61afef",
    primaryRgb: "97, 175, 239",
    secondary: "#c678dd",
    accent: "#98c379",
    trackColor: "#3a3f4b",
    ballColor: "#61afef",
    codeBg: "#21252b",
    codeFg: "#abb2bf",
    codeBorder: "#3a3f4b",
    linkHover: "#61afef",
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
  },
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
      // Handle migration from previous high constract themes to the new theme picker
      if (themeObject.name === "light" || themeObject.name === "dark") {
        if (!themeObject.displayName) {
          themeObject = availableThemes[themeObject.name];
        }
      }
      themeObject = setTheme(themeObject);
    }
  } else {
    themeObject = setTheme(availableThemes.light);
  }

  try {
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
 * @param {string} [selector=".theme-toggle"] - CSS selector for the theme toggle container(s)
 */
function initThemePicker(selector = ".theme-toggle") {
  // Find all theme toggle containers matching the selector
  const themeToggleContainers = document.querySelectorAll(selector);
  if (!themeToggleContainers.length) return;

  // Counter for unique IDs
  let instanceCounter = 0;

  // Process each theme toggle container
  themeToggleContainers.forEach((container) => {
    // Create a unique ID for this instance
    const instanceId = `theme-instance-${instanceCounter++}`;
    container.dataset.themeInstance = instanceId;
    container.innerHTML = "";

    const themeButton = document.createElement("button");
    themeButton.id = `theme-button-${instanceId}`;
    themeButton.className = "theme-button";
    themeButton.innerHTML =
      '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512" width="16" height="16" fill="currentColor"><path d="M512 256c0 141.4-114.6 256-256 256S0 397.4 0 256 114.6 0 256 0s256 114.6 256 256zM256 48C141.1 48 48 141.1 48 256s93.1 208 208 208 208-93.1 208-208S370.9 48 256 48zm0 384c-97.2 0-176-78.8-176-176S158.8 80 256 80s176 78.8 176 176-78.8 176-176 176zm-80-176c0 44.2 35.8 80 80 80s80-35.8 80-80-35.8-80-80-80-80 35.8-80 80z"/></svg>';
    themeButton.setAttribute("aria-label", "Change theme");
    themeButton.setAttribute("title", "Change theme");

    const themeMenu = document.createElement("div");
    themeMenu.id = `theme-menu-${instanceId}`;
    themeMenu.className = "theme-menu";

    const themeCount = Object.keys(availableThemes).length;
    if (themeCount > 4) {
      themeMenu.style.width = "200px";
    }

    for (const themeName in availableThemes) {
      const theme = availableThemes[themeName];

      const themeOption = document.createElement("div");
      themeOption.className = "theme-option";
      themeOption.setAttribute("data-theme", themeName);
      themeOption.setAttribute("title", theme.displayName);
      themeOption.style.backgroundColor = theme.bg;

      if (themeName === "light" || themeName === "dark") {
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

      const themeTitleSpan = document.createElement("span");
      themeTitleSpan.className = "theme-option-name";
      themeTitleSpan.textContent = theme.displayName.split(" ")[0];
      themeOption.appendChild(themeTitleSpan);

      themeOption.addEventListener("click", function () {
        setThemeByName(themeName);
        updateActiveTheme(themeName);
        toggleThemeMenu(instanceId);
        const event = new CustomEvent("themeLoaded", { themeName: themeName });
        document.dispatchEvent(event);
      });

      themeMenu.appendChild(themeOption);
    }

    themeButton.addEventListener("click", function () {
      toggleThemeMenu(instanceId);
    });

    document.addEventListener("click", function (event) {
      if (!container.contains(event.target)) {
        themeMenu.classList.remove("open");
      }
    });

    container.appendChild(themeButton);
    container.appendChild(themeMenu);

    const currentTheme = safeStorage.getItem("theme");
    if (currentTheme) {
      const parsedTheme = JSON.parse(currentTheme);
      updateActiveTheme(parsedTheme.name);
    }
  });
}

/**
 * Updates the active theme in the UI for all theme pickers
 * @param {string} themeName - The name of the active theme
 */
function updateActiveTheme(themeName) {
  // Update all theme options across all instances
  const themeOptions = document.querySelectorAll(".theme-option");
  themeOptions.forEach((option) => {
    if (option.getAttribute("data-theme") === themeName) {
      option.classList.add("active");
    } else {
      option.classList.remove("active");
    }
  });

  // Update all theme buttons across all instances
  const themeButtons = document.querySelectorAll(".theme-button");
  if (availableThemes[themeName]) {
    themeButtons.forEach((button) => {
      button.style.backgroundColor = availableThemes[themeName].primary;
      button.style.borderColor = availableThemes[themeName].primary;
      button.style.color = availableThemes[themeName].bg;
    });
  }
}

/**
 * Toggles the theme menu open/closed
 * @param {string} [instanceId] - Optional ID of the specific instance to toggle
 */
function toggleThemeMenu(instanceId) {
  if (instanceId) {
    // Toggle specific instance
    const themeMenu = document.getElementById(`theme-menu-${instanceId}`);
    if (themeMenu) {
      themeMenu.classList.toggle("open");
    }
  } else {
    // Legacy support for the main theme picker
    const themeMenu = document.getElementById("theme-menu");
    if (themeMenu) {
      themeMenu.classList.toggle("open");
    }
  }
}

/**
 * Calculates the relative luminance of a color according to WCAG 2.0.
 *
 * @param {Object} rgb - The RGB color object.
 * @param {number} rgb.r - The red component (0-255).
 * @param {number} rgb.g - The green component (0-255).
 * @param {number} rgb.b - The blue component (0-255).
 * @returns {number} The relative luminance value between 0 and 1.
 */
function computeLuminance(rgb) {
  // First, we normalize:
  let r = rgb.r / 255;
  let g = rgb.g / 255;
  let b = rgb.b / 255;
  // Second: we transform our RGB to sRGB
  r = r <= 0.04045 ? r / 12.92 : Math.pow((r + 0.055) / 1.055, 2.4);
  g = g <= 0.04045 ? g / 12.92 : Math.pow((g + 0.055) / 1.055, 2.4);
  b = b <= 0.04045 ? b / 12.92 : Math.pow((b + 0.055) / 1.055, 2.4);
  // Third: we apply the formula for luminance
  return 0.2126 * r + 0.7152 * g + 0.0722 * b;
}

/**
 * Calculates the contrast ratio between two colors according to WCAG 2.0.
 *
 * @param {Object} rgb1 - The first RGB color object.
 * @param {number} rgb1.r - The red component (0-255).
 * @param {number} rgb1.g - The green component (0-255).
 * @param {number} rgb1.b - The blue component (0-255).
 * @param {Object} rgb2 - The second RGB color object.
 * @param {number} rgb2.r - The red component (0-255).
 * @param {number} rgb2.g - The green component (0-255).
 * @param {number} rgb2.b - The blue component (0-255).
 * @returns {number} The contrast ratio, ranging from 1 to 21.
 */
function computeContrastRatio(rgb1, rgb2) {
  const lum1 = computeLuminance(rgb1);
  const lum2 = computeLuminance(rgb2);
  const lighter = Math.max(lum1, lum2);
  const darker = lighter == lum1 ? lum2 : lum1;
  return (lighter + 0.05) / (darker + 0.05);
}

document.addEventListener("DOMContentLoaded", function () {
  initThemePicker();
  loadThemeFromLocalStorage();
});

function hexToRgb(hex) {
  hex = hex.replace(/^#/, "");
  const bigint = parseInt(hex, 16);
  return {
    r: (bigint >> 16) & 255,
    g: (bigint >> 8) & 255,
    b: bigint & 255,
  };
}

function rgbToHex(r, g, b) {
  return "#" + ((1 << 24) + (r << 16) + (g << 8) + b).toString(16).slice(1);
}

function convertToLinear(val) {
  return val <= 0.04045 ? val / 12.92 : Math.pow((val + 0.055) / 1.055, 2.4);
}

function rgbToHsl(r, g, b) {
  r /= 255;
  g /= 255;
  b /= 255;

  const max = Math.max(r, g, b);
  const min = Math.min(r, g, b);
  let h,
    s,
    l = (max + min) / 2;

  if (max === min) {
    h = s = 0;
  } else {
    const d = max - min;
    s = l > 0.5 ? d / (2 - max - min) : d / (max + min);

    switch (max) {
      case r:
        h = (g - b) / d + (g < b ? 6 : 0);
        break;
      case g:
        h = (b - r) / d + 2;
        break;
      case b:
        h = (r - g) / d + 4;
        break;
    }

    h /= 6;
  }

  return [h * 360, s * 100, l * 100];
}

function hslToRgb(h, s, l) {
  h /= 360;
  s /= 100;
  l /= 100;

  let r, g, b;

  if (s === 0) {
    r = g = b = l;
  } else {
    const hue2rgb = (p, q, t) => {
      if (t < 0) t += 1;
      if (t > 1) t -= 1;
      if (t < 1 / 6) return p + (q - p) * 6 * t;
      if (t < 1 / 2) return q;
      if (t < 2 / 3) return p + (q - p) * (2 / 3 - t) * 6;
      return p;
    };

    const q = l < 0.5 ? l * (1 + s) : l + s - l * s;
    const p = 2 * l - q;

    r = hue2rgb(p, q, h + 1 / 3);
    g = hue2rgb(p, q, h);
    b = hue2rgb(p, q, h - 1 / 3);
  }

  return [Math.round(r * 255), Math.round(g * 255), Math.round(b * 255)];
}

function shiftHue(color, degrees) {
  const hsl = rgbToHsl(color.r, color.g, color.b);
  hsl[0] = (hsl[0] + degrees) % 360;
  if (hsl[0] < 0) hsl[0] += 360;

  const rgb = hslToRgb(hsl[0], hsl[1], hsl[2]);
  return { r: rgb[0], g: rgb[1], b: rgb[2] };
}

/**
 * Converts a color string to an RGB array.
 *
 * @param {string} colorStr - Color in hex format (#RRGGBB or #RGB) or a named color
 * @returns {number[]} Array of RGB values as integers [R, G, B]
 * @throws {Error} If the color format is not recognized
 */
function parseColorToRGB(colorStr) {
  if (colorStr.startsWith("#")) {
    let hex = colorStr.slice(1);

    if (hex.length === 3) {
      hex = hex
        .split("")
        .map((char) => char + char)
        .join("");
    }

    return [
      parseInt(hex.substring(0, 2), 16),
      parseInt(hex.substring(2, 4), 16),
      parseInt(hex.substring(4, 6), 16),
    ];
  }

  const colorMap = {
    black: [0, 0, 0],
    white: [255, 255, 255],
    red: [255, 0, 0],
    green: [0, 128, 0],
    blue: [0, 0, 255],
    yellow: [255, 255, 0],
    cyan: [0, 255, 255],
    magenta: [255, 0, 255],
    gray: [128, 128, 128],
    purple: [128, 0, 128],
    orange: [255, 165, 0],
  };

  color = colorMap[colorStr.toLowerCase()];

  if (!color) {
    throw new Error(`Unexpected color: ${colorStr.toLowerCase()}`);
  }

  return color;
}

function adjustImagesColors() {
  const colorThief = new ColorThief();

  const rootStyles = window.getComputedStyle(document.documentElement);
  let bgColor = rootStyles.getPropertyValue("--bg").trim();
  bgColor = parseColorToRGB(bgColor);
  bgColor = {
    r: bgColor[0],
    g: bgColor[1],
    b: bgColor[2],
  };

  const images = document.querySelectorAll("img.dynamic-colors");

  for (const image of images) {
    const colors = colorThief.getPalette(image);
    const dominantColor = {
      r: colors[0][0],
      g: colors[0][1],
      b: colors[0][2],
    };

    console.log(dominantColor);

    const contrastRatio = computeContrastRatio(dominantColor, bgColor);

    if (contrastRatio < 7.0) {
      const optimalShift = findOptimalHueShift(dominantColor, bgColor);

      applyHueShiftToImage(image, optimalShift);
    }
  }
}

/**
 * Finds a color with better contrast using gradient ascent
 * @param {Object} dominantColor - RGB object {r, g, b} for the dominant color
 * @param {Object} bgColor - RGB object {r, g, b} for the background color
 * @returns {number} The optimal hue shift in degrees
 */
function findOptimalHueShift(dominantColor, bgColor) {
  const [originalHue, saturation, lightness] = rgbToHsl(
    dominantColor.r,
    dominantColor.g,
    dominantColor.b,
  );

  const originalContrast = computeContrastRatio(dominantColor, bgColor);

  const bgLuminance = computeLuminance(bgColor);

  let currentHue = bgLuminance > 0.5 ? (originalHue + 180) % 360 : originalHue;

  let currentContrast = computeContrastWithHue(
    saturation,
    lightness,
    currentHue,
    bgColor,
  );

  let bestHue = currentHue;
  let bestContrast = currentContrast;

  let stepSize = 30;
  let improving = true;

  while (improving && stepSize >= 1) {
    improving = false;

    const clockwiseHue = (currentHue + stepSize) % 360;
    const clockwiseContrast = computeContrastWithHue(
      saturation,
      lightness,
      clockwiseHue,
      bgColor,
    );

    const counterClockwiseHue = (currentHue - stepSize + 360) % 360;
    const counterClockwiseContrast = computeContrastWithHue(
      saturation,
      lightness,
      counterClockwiseHue,
      bgColor,
    );

    if (
      clockwiseContrast > bestContrast &&
      clockwiseContrast >= counterClockwiseContrast
    ) {
      bestContrast = clockwiseContrast;
      bestHue = clockwiseHue;
      currentHue = clockwiseHue;
      improving = true;
    } else if (counterClockwiseContrast > bestContrast) {
      bestContrast = counterClockwiseContrast;
      bestHue = counterClockwiseHue;
      currentHue = counterClockwiseHue;
      improving = true;
    } else {
      stepSize = Math.floor(stepSize / 2);
    }

    if (bestContrast >= 7.0) {
      break;
    }
  }

  let hueShift = bestHue - originalHue;

  if (hueShift > 180) hueShift -= 360;
  if (hueShift < -180) hueShift += 360;

  if (bestContrast < originalContrast * 1.2) {
    hueShift = 0;
  }

  return hueShift;
}

function computeContrastWithHue(saturation, lightness, newHue, bgColor) {
  const [r, g, b] = hslToRgb(newHue, saturation, lightness);
  return computeContrastRatio({ r, g, b }, bgColor);
}

/**
 * Apply the hue shift to an image using canvas
 * @param {HTMLImageElement} image - The image element to modify
 * @param {number} hueShift - The hue shift in degrees to apply
 */
function applyHueShiftToImage(image, hueShift) {
  if (hueShift === 0) return;

  const canvas = document.createElement("canvas");
  const ctx = canvas.getContext("2d");

  canvas.width = image.naturalWidth;
  canvas.height = image.naturalHeight;

  ctx.drawImage(image, 0, 0);

  const imageData = ctx.getImageData(0, 0, canvas.width, canvas.height);
  const data = imageData.data;

  for (let i = 0; i < data.length; i += 4) {
    const r = data[i];
    const g = data[i + 1];
    const b = data[i + 2];
    const a = data[i + 3];

    if (a === 0) continue;

    const [h, s, l] = rgbToHsl(r, g, b);
    const newHue = (h + hueShift) % 360;
    const [newR, newG, newB] = hslToRgb(newHue, s, l);

    data[i] = newR;
    data[i + 1] = newG;
    data[i + 2] = newB;
  }

  ctx.putImageData(imageData, 0, 0);

  const dataURL = canvas.toDataURL();

  if (!image.dataset.originalSrc) {
    image.dataset.originalSrc = image.src;
  }

  image.src = dataURL;
}
