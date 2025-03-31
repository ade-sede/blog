initThemePicker("#inline-theme-picker");
initThemePicker("#conclusion-theme-picker");

function initContrastDemo() {
  const fgPicker = document.getElementById("fg-color");
  const bgPicker = document.getElementById("bg-color");
  const sampleText = document.getElementById("sample-text");
  const stepsContainer = document.getElementById("calculation-steps");

  const computedStyle = getComputedStyle(document.documentElement);
  const defaultFg = computedStyle.getPropertyValue("--fg").trim() || "#000000";
  const defaultBg = computedStyle.getPropertyValue("--bg").trim() || "#ffffff";

  fgPicker.value = defaultFg;
  bgPicker.value = defaultBg;

  function updateContrast() {
    const fgColor = hexToRgb(fgPicker.value);
    const bgColor = hexToRgb(bgPicker.value);

    sampleText.style.color = fgPicker.value;
    sampleText.style.backgroundColor = bgPicker.value;

    showCalculationSteps(fgColor, bgColor);
  }

  function showCalculationSteps(fgColor, bgColor) {
    stepsContainer.innerHTML = "";

    const stepNormalization = document.createElement("div");
    stepNormalization.classList.add("step");
    stepNormalization.innerHTML =
      "<strong>Step 1: Normalize RGB values</strong>";

    const fgNormalized = {
      r: (fgColor.r / 255).toFixed(3),
      g: (fgColor.g / 255).toFixed(3),
      b: (fgColor.b / 255).toFixed(3),
    };

    const bgNormalized = {
      r: (bgColor.r / 255).toFixed(3),
      g: (bgColor.g / 255).toFixed(3),
      b: (bgColor.b / 255).toFixed(3),
    };

    stepNormalization.innerHTML += `<div>FG: rgb(${fgColor.r}, ${fgColor.g}, ${fgColor.b}) → (${fgNormalized.r}, ${fgNormalized.g}, ${fgNormalized.b})</div>`;
    stepNormalization.innerHTML += `<div>BG: rgb(${bgColor.r}, ${bgColor.g}, ${bgColor.b}) → (${bgNormalized.r}, ${bgNormalized.g}, ${bgNormalized.b})</div>`;

    const stepLinear = document.createElement("div");
    stepLinear.classList.add("step");
    stepLinear.innerHTML =
      "<strong>Step 2: Convert to linear light (sRGB)</strong>";

    const fgLinear = {
      r: convertToLinear(fgColor.r / 255),
      g: convertToLinear(fgColor.g / 255),
      b: convertToLinear(fgColor.b / 255),
    };

    const bgLinear = {
      r: convertToLinear(bgColor.r / 255),
      g: convertToLinear(bgColor.g / 255),
      b: convertToLinear(bgColor.b / 255),
    };

    stepLinear.innerHTML += `<div>FG linear: (${fgLinear.r.toFixed(3)}, ${fgLinear.g.toFixed(3)}, ${fgLinear.b.toFixed(3)})</div>`;
    stepLinear.innerHTML += `<div>BG linear: (${bgLinear.r.toFixed(3)}, ${bgLinear.g.toFixed(3)}, ${bgLinear.b.toFixed(3)})</div>`;

    const stepLuminance = document.createElement("div");
    stepLuminance.classList.add("step");
    stepLuminance.innerHTML =
      "<strong>Step 3: Calculate relative luminance</strong>";

    const lumFg = computeLuminance(fgColor);
    const lumBg = computeLuminance(bgColor);

    stepLuminance.innerHTML += `<div>FG: 0.2126×${fgLinear.r.toFixed(3)} + 0.7152×${fgLinear.g.toFixed(3)} + 0.0722×${fgLinear.b.toFixed(3)} = ${lumFg.toFixed(3)}</div>`;
    stepLuminance.innerHTML += `<div>BG: 0.2126×${bgLinear.r.toFixed(3)} + 0.7152×${bgLinear.g.toFixed(3)} + 0.0722×${bgLinear.b.toFixed(3)} = ${lumBg.toFixed(3)}</div>`;

    const stepContrast = document.createElement("div");
    stepContrast.classList.add("step");
    stepContrast.innerHTML =
      "<strong>Step 4: Calculate contrast ratio</strong>";

    const lighter = Math.max(lumFg, lumBg);
    const darker = Math.min(lumFg, lumBg);
    const contrast = (lighter + 0.05) / (darker + 0.05);

    stepContrast.innerHTML += `<div>(${lighter.toFixed(3)} + 0.05) / (${darker.toFixed(3)} + 0.05) = ${contrast.toFixed(2)}</div>`;

    const stepVerdict = document.createElement("div");
    stepVerdict.classList.add("step");
    stepVerdict.innerHTML = "<strong>Result</strong>";

    let wcagLevel = "";
    if (contrast >= 7.0) {
      wcagLevel = "Passes AAA level (7:1)";
    } else if (contrast >= 4.5) {
      wcagLevel = "Passes AA level (4.5:1)";
    } else if (contrast >= 3.0) {
      wcagLevel = "Passes AA level for large text (3:1)";
    } else {
      wcagLevel = "Fails WCAG requirements";
    }

    stepVerdict.innerHTML += `<div>Contrast ratio: ${contrast.toFixed(2)}:1</div>`;
    stepVerdict.innerHTML += `<div>${wcagLevel}</div>`;

    document.getElementById("contrast-value").innerHTML =
      `Contrast: <span>${contrast.toFixed(2)}:1</span>`;

    stepsContainer.appendChild(stepNormalization);
    stepsContainer.appendChild(stepLinear);
    stepsContainer.appendChild(stepLuminance);
    stepsContainer.appendChild(stepContrast);
    stepsContainer.appendChild(stepVerdict);
  }

  fgPicker.addEventListener("input", updateContrast);
  bgPicker.addEventListener("input", updateContrast);

  updateContrast();
}

function initHueWheel() {
  const wheel = document.getElementById("hue-wheel");
  const indicator = document.getElementById("hue-indicator");
  const valueDisplay = document.getElementById("hue-value");

  if (!wheel || !indicator || !valueDisplay) {
    console.error("Missing required elements for hue wheel");
    return;
  }

  wheel.innerHTML = "";

  const canvas = document.createElement("canvas");
  canvas.width = wheel.clientWidth;
  canvas.height = wheel.clientHeight;
  wheel.appendChild(canvas);

  const ctx = canvas.getContext("2d");
  const centerX = canvas.width / 2;
  const centerY = canvas.height / 2;
  const radius = Math.min(centerX, centerY) - 5;

  for (let angle = 0; angle < 360; angle++) {
    const startAngle = ((angle - 1) * Math.PI) / 180;
    const endAngle = ((angle + 1) * Math.PI) / 180;

    ctx.beginPath();
    ctx.moveTo(centerX, centerY);
    ctx.arc(centerX, centerY, radius, startAngle, endAngle);
    ctx.closePath();

    ctx.fillStyle = `hsl(${angle}, 100%, 50%)`;
    ctx.fill();
  }

  const keyColors = [
    { angle: 0, label: "Red (0°)" },
    { angle: 60, label: "Yellow (60°)" },
    { angle: 120, label: "Green (120°)" },
    { angle: 180, label: "Cyan (180°)" },
    { angle: 240, label: "Blue (240°)" },
    { angle: 300, label: "Magenta (300°)" },
  ];

  ctx.strokeStyle = "rgba(255, 255, 255, 0.7)";
  ctx.lineWidth = 2;

  keyColors.forEach((color) => {
    const radians = (color.angle * Math.PI) / 180;

    ctx.beginPath();
    ctx.moveTo(centerX, centerY);
    ctx.lineTo(
      centerX + Math.cos(radians) * radius,
      centerY + Math.sin(radians) * radius,
    );
    ctx.stroke();

    ctx.fillStyle = "white";
    ctx.beginPath();
    ctx.arc(
      centerX + Math.cos(radians) * radius,
      centerY + Math.sin(radians) * radius,
      3,
      0,
      Math.PI * 2,
    );
    ctx.fill();
  });

  updateIndicator(0);

  function updateIndicator(hue) {
    const radians = (hue * Math.PI) / 180;
    const x = centerX + Math.cos(radians) * (radius - 10);
    const y = centerY + Math.sin(radians) * (radius - 10);

    indicator.style.left = `${x}px`;
    indicator.style.top = `${y}px`;
    valueDisplay.textContent = `Hue: ${Math.round(hue)}°`;

    const colorSample = document.getElementById("color-sample");
    if (colorSample) {
      colorSample.style.backgroundColor = `hsl(${hue}, 100%, 50%)`;
    }
  }

  function calculateHue(e) {
    const rect = canvas.getBoundingClientRect();
    const x = e.clientX - rect.left - centerX;
    const y = e.clientY - rect.top - centerY;

    let angle = (Math.atan2(y, x) * 180) / Math.PI;
    return (angle + 360) % 360;
  }

  canvas.addEventListener("mousemove", function (e) {
    const hue = calculateHue(e);
    updateIndicator(hue);
  });

  canvas.addEventListener("click", function (e) {
    const hue = calculateHue(e);
    updateIndicator(hue);
  });
}

function initSaturationSlider(hue = 0, lightness = 50) {
  const gradient = document.getElementById("saturation-gradient");
  const indicator = document.getElementById("saturation-indicator");
  const valueDisplay = document.getElementById("saturation-value");
  const sample = document.getElementById("saturation-sample");

  if (!gradient || !indicator || !valueDisplay || !sample) {
    console.error("Missing required elements for saturation slider");
    return;
  }

  gradient.style.background = `linear-gradient(to right,
    hsl(${hue}, 0%, ${lightness}%),
    hsl(${hue}, 100%, ${lightness}%))`;

  updateSaturation(100);

  function updateSaturation(saturation) {
    const position = (saturation / 100) * gradient.offsetWidth;
    indicator.style.left = `${position}px`;

    valueDisplay.textContent = `Saturation: ${Math.round(saturation)}%`;
    sample.style.backgroundColor = `hsl(${hue}, ${saturation}%, ${lightness}%)`;
  }

  function handleClick(e) {
    const rect = gradient.getBoundingClientRect();
    const position = e.clientX - rect.left;
    const saturation = Math.max(
      0,
      Math.min(100, (position / rect.width) * 100),
    );
    updateSaturation(saturation);
  }

  gradient.addEventListener("click", handleClick);

  gradient.addEventListener("mousedown", function (e) {
    const handleDrag = function (e) {
      const rect = gradient.getBoundingClientRect();
      const position = e.clientX - rect.left;
      const saturation = Math.max(
        0,
        Math.min(100, (position / rect.width) * 100),
      );
      updateSaturation(saturation);
    };

    document.addEventListener("mousemove", handleDrag);
    document.addEventListener(
      "mouseup",
      function () {
        document.removeEventListener("mousemove", handleDrag);
      },
      { once: true },
    );

    handleDrag(e);
  });

  gradient.addEventListener("click", function (e) {
    const rect = gradient.getBoundingClientRect();
    const position = e.clientX - rect.left;
    const saturation = Math.max(
      0,
      Math.min(100, (position / rect.width) * 100),
    );
    updateSaturation(saturation);
  });

  function updateHue(newHue) {
    hue = newHue;
    gradient.style.background = `linear-gradient(to right,
      hsl(${hue}, 0%, ${lightness}%),
      hsl(${hue}, 100%, ${lightness}%))`;

    const position = parseFloat(indicator.style.left) || 0;
    const saturation = (position / gradient.offsetWidth) * 100;

    sample.style.backgroundColor = `hsl(${hue}, ${saturation}%, ${lightness}%)`;
  }

  return {
    updateHue,
  };
}

function initLightnessSlider() {
  const gradient = document.getElementById("lightness-gradient");
  const indicator = document.getElementById("lightness-indicator");
  const valueDisplay = document.getElementById("lightness-value");
  const sample = document.getElementById("lightness-sample");

  if (!gradient || !indicator || !valueDisplay || !sample) {
    console.error("Missing required elements for lightness slider");
    return;
  }

  const hue = 210; // Blue
  const saturation = 100;

  gradient.style.background = `linear-gradient(to right,
    hsl(${hue}, ${saturation}%, 0%),
    hsl(${hue}, ${saturation}%, 50%),
    hsl(${hue}, ${saturation}%, 100%))`;

  updateLightness(50);

  function updateLightness(lightness) {
    const position = (lightness / 100) * gradient.offsetWidth;
    indicator.style.left = `${position}px`;

    valueDisplay.textContent = `Lightness: ${Math.round(lightness)}%`;
    sample.style.backgroundColor = `hsl(${hue}, ${saturation}%, ${lightness}%)`;
  }

  function handleClick(e) {
    const rect = gradient.getBoundingClientRect();
    const position = e.clientX - rect.left;
    const lightness = Math.max(0, Math.min(100, (position / rect.width) * 100));
    updateLightness(lightness);
  }

  gradient.addEventListener("click", handleClick);

  gradient.addEventListener("mousedown", function (e) {
    const handleDrag = function (e) {
      const rect = gradient.getBoundingClientRect();
      const position = e.clientX - rect.left;
      const lightness = Math.max(
        0,
        Math.min(100, (position / rect.width) * 100),
      );
      updateLightness(lightness);
    };

    document.addEventListener("mousemove", handleDrag);
    document.addEventListener(
      "mouseup",
      function () {
        document.removeEventListener("mousemove", handleDrag);
      },
      { once: true },
    );

    handleDrag(e);
  });

  gradient.addEventListener("click", function (e) {
    const rect = gradient.getBoundingClientRect();
    const position = e.clientX - rect.left;
    const lightness = Math.max(0, Math.min(100, (position / rect.width) * 100));
    updateLightness(lightness);
  });
}

document.addEventListener("themeLoaded", function () {
  try {
    adjustImagesColors();
  } catch {}
});

document.addEventListener("DOMContentLoaded", function () {
  if (document.getElementById("fg-color")) {
    initContrastDemo();
  }

  if (document.getElementById("hue-wheel")) {
    initHueWheel();
  }

  if (document.getElementById("saturation-gradient")) {
    initSaturationSlider();
  }

  if (document.getElementById("lightness-gradient")) {
    initLightnessSlider();
  }

  try {
    adjustImagesColors();
  } catch {}
});
