/**
 * @typedef {Object} Theme
 * @property {string} name - Theme name
 * @property {string} bg - Background color value
 * @property {string} fg - Foreground color value
 * @property {string} trackColor - Color of theme change track
 * @property {string} ballColor - Color of theme change ball
 */

/**
 * Loads theme from localStorage and applies it if present
 * @returns {Theme}
 */
function loadThemeFromLocalStorage() {
  let themeObject;
  const themeString = localStorage.getItem("theme");

  if (themeString) {
    themeObject = JSON.parse(themeString);

    if (themeObject) {
      themeObject = setTheme(themeObject);
    }
  } else {
    themeObject = setDarkTheme();
  }

  let slider = document.getElementById("theme-slider");

  if (themeObject.name === "dark") {
    try {
      slider.checked = true;
    } catch (err) {}
  } else if (themeObject.name === "light") {
    try {
      slider.checked = false;
    } catch (err) {}
  } else {
    throw new Error(`Unknown theme name: ${themeObject.name}`);
  }

  return themeObject;
}

function toggleTheme() {
  const themeString = localStorage.getItem("theme");

  if (!themeString) {
    throw new Error("Can't reverse theme, no theme is currently set");
  }

  const themeObject = JSON.parse(themeString);

  if (themeObject.name === "dark") {
    return setLightTheme();
  } else if (themeObject.name === "light") {
    return setDarkTheme();
  } else {
    throw new Error(`Unknown theme name: ${themeObject.name}`);
  }
}

/**
 * Sets dark theme and saves to localStorage
 * @returns {Theme}
 */
function setDarkTheme() {
  const themeObject = {
    name: "dark",
    bg: "black",
    fg: "white",
    trackColor: "#333333",
    ballColor: "#000000",
  };
  setTheme(themeObject);
  localStorage.setItem("theme", JSON.stringify(themeObject));

  return themeObject;
}

/**
 * Sets light theme and saves to localStorage
 * @returns {Theme}
 */
function setLightTheme() {
  const themeObject = {
    name: "light",
    bg: "white",
    fg: "black",
    trackColor: "#e0e0e0",
    ballColor: "#ffffff",
  };
  setTheme(themeObject);
  localStorage.setItem("theme", JSON.stringify(themeObject));

  return themeObject;
}

/**
 * Applies theme by setting CSS custom properties
 * Note setting the variables is not enough.
 * The stylesheet is loaded after this script is run
 * Before the style sheet is loaded, the default background of
 * the html and body elements is the browser's default which can cause flickering.
 * @param {Theme} themeObject - The theme object containing color values
 * @returns {Theme}
 */
function setTheme(themeObject) {
  try {
    const body = document.querySelector("body");
    body.style.setProperty("background-color", "var(--bg)");
    body.style.setProperty("color", "var(--fg)");
  } catch (err) {}

  try {
    const navbar = document.querySelector("nav.navbar");
    navbar.style.setProperty("--track-color", themeObject.trackColor);
    navbar.style.setProperty("--ball-color", themeObject.ballColor);
  } catch (err) {}

  try {
    const root = document.querySelector(":root");
    root.style.setProperty("--bg", themeObject.bg);
    root.style.setProperty("--fg", themeObject.fg);
  } catch (err) {}

  return themeObject;
}
