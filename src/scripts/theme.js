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

const availableThemes = {
  light: { name: "light", displayName: "Light (High Contrast)" },
  dark: { name: "dark", displayName: "Dark (High Contrast)" },
  atom: { name: "atom", displayName: "Atom Dark" },
  nord: { name: "nord", displayName: "Nord" },
  solarized: { name: "solarized", displayName: "Solarized Light" },
  dracula: { name: "dracula", displayName: "Dracula" },
  github: { name: "github", displayName: "GitHub" },
  monokai: { name: "monokai", displayName: "Monokai" },
};

/**
 * Loads theme from localStorage and applies it if present
 * @returns {string} The loaded theme name
 */
function loadThemeFromLocalStorage() {
  const themeString = safeStorage.getItem("theme");
  let themeName = "light";

  if (themeString) {
    try {
      const themeObject = JSON.parse(themeString);
      if (themeObject && themeObject.name && availableThemes[themeObject.name]) {
        themeName = themeObject.name;
      }
    } catch (err) {
      console.warn("Failed to parse stored theme, using default");
    }
  }

  setTheme(themeName);
  return themeName;
}

/**
 * Toggles between light and dark themes (for backward compatibility)
 * @returns {string} The new theme name after toggle
 */
function toggleTheme() {
  const themeString = safeStorage.getItem("theme");
  let currentTheme = "light";
  
  if (themeString) {
    try {
      const themeObject = JSON.parse(themeString);
      if (themeObject && themeObject.name) {
        currentTheme = themeObject.name;
      }
    } catch (err) {}
  }

  const newTheme = currentTheme === "dark" ? "light" : "dark";
  return setTheme(newTheme);
}

/**
 * Sets theme by name
 * @param {string} themeName - Name of the theme to apply
 * @returns {string} The applied theme name
 */
function setThemeByName(themeName) {
  if (availableThemes[themeName]) {
    return setTheme(themeName);
  } else {
    console.warn(`Theme "${themeName}" not found, using default light theme`);
    return setTheme("light");
  }
}

/**
 * Sets dark theme (for backward compatibility)
 * @returns {string} The dark theme name
 */
function setDarkTheme() {
  return setTheme("dark");
}

/**
 * Sets light theme (for backward compatibility)
 * @returns {string} The light theme name
 */
function setLightTheme() {
  return setTheme("light");
}

/**
 * Applies theme by setting the data-theme attribute.
 * CSS handles all color changes through attribute selectors for better performance.
 * @param {string} themeName - The theme name to apply
 * @returns {string} The applied theme name
 */
function setTheme(themeName) {
  if (!availableThemes[themeName]) {
    console.warn(`Theme "${themeName}" not found, using light theme`);
    themeName = "light";
  }

  try {
    const root = document.documentElement;
    
    root.style.removeProperty("--bg");
    root.style.removeProperty("--fg");
    root.style.removeProperty("--primary");
    root.style.removeProperty("--primary-rgb");
    root.style.removeProperty("--secondary");
    root.style.removeProperty("--accent");
    root.style.removeProperty("--code-bg");
    root.style.removeProperty("--code-fg");
    root.style.removeProperty("--code-border");
    root.style.removeProperty("--link-hover");
    
    root.setAttribute("data-theme", themeName);
  } catch (err) {
    console.warn("Failed to set theme:", err);
  }

  safeStorage.setItem("theme", JSON.stringify({ 
    name: themeName, 
    displayName: availableThemes[themeName].displayName 
  }));

  return themeName;
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
    themeButton.innerHTML = '<i class="fas fa-palette"></i>';
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
      
      themeOption.classList.add(`theme-preview-${themeName}`);

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

  const themeButtons = document.querySelectorAll(".theme-button");
  themeButtons.forEach((button) => {
    button.setAttribute("data-current-theme", themeName);
    button.style.backgroundColor = "";
    button.style.borderColor = "";
    button.style.color = "";
  });
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

  const color = colorMap[colorStr.toLowerCase()];

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
