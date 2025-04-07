## Some colors just don't fit together...
I recently built a theme picker for my blog and I am very happy with it. Readers can pick from a selection of various themes, some light, some dark. It is simple, it is light, and it instantly applies to all styling elements of the page ... except images.

An image with dominating tones of a pale yellow, easily readable on black background could be impossible to read on a white background.  
Don't trust me ? Look at the following image and play around with the theme picker.  
You can see the text clear as day when using _Nord_ or _Dracula_ but it is almost invisible when using _Solarized_.

<div style="display: flex; align-items: center; justify-content: center; gap: 20px;">
  <div id="inline-theme-picker" class="theme-toggle"></div>
  <img src="./images/unreadable-on-white.png" width="300px">
</div>

The thing is, I plan on writing articles containing various images or schematics containing text.  
And even if they are readable, I simply think it would be cool to be able to make sure images always match the general theme of the blog.
Considering I have no knowledge of how any of this works, let's start from the beginning and work our way up !


## Contrast is key

Contrast is one of the ways our brain identifies the boundaries between objects.  
Without it, we simply don't know where things start and where they end. This is particularly true for people with vision impairment. That is what makes contrast such an important aspect of _accessibility_.

And because contrast is so important, we have introduced standards to measure & enforce: meet [WCAG's Luminance Contrast Ratio](https://www.w3.org/WAI/WCAG21/Understanding/contrast-minimum.html).
The idea is simple: two colors with a similar luminance do not have enough contrast and will blend together. In order to guarantee the average person, and even some visually impaired individuals can see properly, the luminance of both colors must respect a certain ratio.  
Using this standard we should be able to compute whether or not two colors properly contrast with each other.  
But before we can do that, we need to build some fundamentals!  

### Luminance & Brightness

Light is a wave. It has an amplitude and a wavelength.  
The wavelength dictates the perceived color of the light.  
The amplitude dictates its intensity.  
_Luminance_ is an objective measurement of how intense a light is, measured in candela per square meter (cd/m²).[^1]
_Brightness_ is the subjective perception of that intensity by a given person.  

In digital systems, a pixel has three color channels: red, green, blue (RGB).  
Each channel is represented by an 8 bit value, ranging from 0 to 255.  
A common misconception is to think that changing the level of a color channel means changing the color that is emitted. But that is not actually what is going on.  
On OLED and AMOLED displays, a pixel is comprised of 3 LEDs: 1 green, 1 red, 1 blue.[^2]  
Each of these diodes emits light at a constant wavelength: changing the color of the light is impossible!  
When we increase the value of a color channel, we are in fact increasing the light's amplitude: we are increasing the brightness.[^3]  
Our eyes can interpet the mix of intensity of these 3 lights as a color.  

### Gamma 

But there is one important caveat: our eyes do not perceive a light twice as intense as twice as bright.  
The relationship between the intensity of a light that is emitted (_luminance_) and the intensity as perceived by the human eye (_brightness_) is exponential.  
On top of that, our eyes are better at spotting changes in darker tones than brigther ones.  
If we were to use our 8 bits to encode the luminance linearly, half of the range would be dedicated to the brightest 10%. Giving us the other half of the range for the 90% of darker tones where our eyes actualy perform much better.  
That is where _gamma_ (denoted as _γ_) and more generally the _gamma encoding_ function (denoted as _Γ_) come in.  
This _gamma encoding_ is meant to be an approximiation of the non-linear relationship between luminance and brightness.  
To increase the efficiency of our color encoding and make sure we get the most out of our 8 bits, the color channel is already gamma encoded!  
The value 255 is not twice as luminous as 128, but rather 4.5 times more luminous, which would make it roughly 2 times brighter to the human eye[^4] ! 

Gamma encoding and decoding works as follows:
\[
V_{linear} = V_{encoded}^{\gamma}
\]
\[
V_{encoded} = V_{linear}^{\frac{1}{\gamma}}
\]

Let's walk through a concrete example: suppose our screen's gamma is _2.2_. Our red channel is at 80 and we want to double the luminance.  
- Normalise the value of the red channel:
  \[
    \frac{80}{255} \approx 0.313
  \]
- Convert from gamma to linear representation:
  \[
    0.313^{2.2} \approx 0.097
  \]
- Double our linear representation which corresponds linearly to luminance:
  \[
    0.097 \times 2 = 0.194
  \]
- Convert back from linear to a gamma encoded value:
  \[
    0.194^{\frac{1}{2.2}} \approx 0.475
  \]
- Back to 8 bit representation:
  \[
    0.475 * 255 \approx 121
  \]

And there we have it: **Using a gamma of 2.2, doubling the luminance from 80 means increasing the channel's value to _121_**
To double the perceived brightness you would simply need to double the already gamma encoded value.  
That's the whole point of gamma: it helps us work with values that feel linear to humans.

### Relative luminance

We have built some understanding of what luminance is and how digital systems deal with it, but we still have one practical problem: gamma values can vary from one device to another. The gamma curve of an old CRT monitor is different from a modern OLED display.  

Remember, our immediate goal is to be able to measure and quantify the perception of an average human being. Variations between different type of screens is noise in our model.  

This is where the sRGB (standard RGB) color space comes in. Developed in the 90s, sRGB establishes a standard gamma curve that all devices should aim to reproduce. While we often approximate it as γ = 2.2, sRGB actually defines a more complex function.
I copy it here for completeness but don't ask me how this formula works exactly, my understanding does not run that deep.

\[
V_{linear} =
\begin{cases}
\frac{V_{sRGB}}{12.92} & \text{if } V_{sRGB} \leq 0.04045 \\
\left(\frac{V_{sRGB} + 0.055}{1.055}\right)^{2.4} & \text{if } V_{sRGB} > 0.04045
\end{cases}
\]

With a standardized way to convert encoded values to linear light, we can now tackle the question of perceived brightness in a device-independent way.

Relative luminance (denoted *L*) is a measure of the perceived brightness of a color to the human eye. It is "relative" because it accounts for the fact that human vision is more sensitive to certain wavelengths of light than others.

To compute relative luminance, we convert our RGB values to linear light using the sRGB transfer function, then apply a weighted sum:

\[
L = 0.2126 \times R_{linear} + 0.7152 \times G_{linear} + 0.0722 \times B_{linear}
\]

These weights reflect that we are most sensitive to green light, less sensitive to red light, and even less sensitive to blue light.

To build an even better understanding, let's compute the relative luminance of a vivid red `rgb(230, 30, 30)`:

- Normalise: 
   \[
     R = \frac{230}{255} \approx 0.902, \quad
     G = \frac{30}{255} \approx 0.118, \quad
     B = \frac{30}{255} \approx 0.118
   \]

- Convert to linear using the sRGB function:
   \[
     R_{linear} = \left(\frac{0.902 + 0.055}{1.055}\right)^{2.4} \approx 0.798 \quad \text{(since } 0.902 > 0.04045 \text{)}
   \]
   \[
     G_{linear} = \left(\frac{0.118 + 0.055}{1.055}\right)^{2.4} \approx 0.012 \quad \text{(since } 0.118 > 0.04045 \text{)}
   \]
   \[
     B_{linear} = \left(\frac{0.118 + 0.055}{1.055}\right)^{2.4} \approx 0.012 \quad \text{(since } 0.118 > 0.04045 \text{)}
   \]

- Calculate relative luminance:
   \[
     L = 0.2126 \times 0.798 + 0.7152 \times 0.012 + 0.0722 \times 0.012 \approx 0.178
   \]

This standardized measure of perceived brightness gives us exactly what we need to compute contrast ratios between colors, regardless of what device the colors are displayed on.

### Computing the luminance contrast ratio

We now have all required pieces to compute the _Luminance Contrast Ratio_.
- We understand what luminance & brightness are
- We understand what gamma is, how it differs from one device to another and how to make it a non-factor by working in the sRGB space
- We understand what _relative_ luminance is

\[
R = (L_{LighterColor}  + 0.05) / (L_{DarkerColor} + 0.05)
\]

Interpretation depends on two things: the ratio itself and the font-size it applies too.

| Level | Text Size | Minimum Contrast Ratio |
|-------|-----------|------------------------|
| AA (minimum compliance) | Regular text (< 18pt, or < 14pt bold) | 4.5:1 |
| AA (minimum compliance) | Large text (≥ 18pt, or ≥ 14pt bold) | 3:1 |
| AAA (enhanced compliance) | Regular text (< 18pt, or < 14pt bold) | 7:1 |
| AAA (enhanced compliance) | Large text (≥ 18pt, or ≥ 14pt bold) | 4.5:1 |

Level AA accommodates users with moderate visual impairments, approximately equivalent to 20/40 vision.  
Level AAA provides enhanced readability for users with more substantial vision loss, up to approximately 20/80 vision.


### Let's put it in practice !

Once you understand the model it simply comes down to applying the formulas!


```javascript:color_convert.js
function computeLuminance(rgb) {
  let r = rgb.r / 255;
  let g = rgb.g / 255;
  let b = rgb.b / 255;

  r = r <= 0.04045 ? r / 12.92 : Math.pow((r + 0.055) / 1.055, 2.4);
  g = g <= 0.04045 ? g / 12.92 : Math.pow((g + 0.055) / 1.055, 2.4);
  b = b <= 0.04045 ? b / 12.92 : Math.pow((b + 0.055) / 1.055, 2.4);

  return 0.2126 * r + 0.7152 * g + 0.0722 * b;
}
```

```javascript:color_convert.js
function computeContrastRatio(rgb1, rgb2) {
  const lum1 = computeLuminance(rgb1);
  const lum2 = computeLuminance(rgb2);

  const lighter = Math.max(lum1, lum2);
  const darker = lighter == lum1 ? lum2 : lum1;

  return (lighter + 0.05) / (darker + 0.05);
}
```

Just sprinkle a bit of UI on top, and you have an evaluator ready to experiment with various colors !

<div class="contrast-demo">
  <div class="controls">
    <div class="control-group">
      <label for="fg-color">Foreground:</label>
      <input type="color" id="fg-color" value="#000000">
    </div>
    <div class="control-group">
      <label for="bg-color">Background:</label>
      <input type="color" id="bg-color" value="#ffffff">
    </div>
  </div>
  
  <div id="sample-text" class="sample">
    Pick the colors of your choice, the color pickers are fully interactive 🎉 !
  </div>
  
  <div class="calculation">
    <div id="contrast-value">Contrast: <span></span></div>
    <div id="calculation-steps" class="steps">
      <!-- Steps updated dynamically via JS -->
    </div>
  </div>
</div>

## Dynamic color adjustments

Now that we are able to quantify contrast and understand what it actually is, we can work on modifying our images to make sure they properly contrast with our blog.  
The plan is as follows:
- Detect the dominant color in an image
- Detect whether or not it contrasts properly with our background
- If it does not, pick a color yielding a good contrast ratio and use it to replace the dominant color
- Replace other colors accordingly

That last step is definitely the hardest !
We want to shift the color tone of the image to have better contrast, but we want the image as a whole to remain coherent.  
And the most straightforward way to achieve that is most likely to use the _HSL Color Model_

### HSL Color Model

_**H**ue, **S**aturation, and **L**ightness_ is a color model that represents colors in a way that's more intuitive to how humans think about color than RGB.

_Hue_ represents the color type: red, green, blue, etc..
- it is measured as an angle around a color wheel (0-360°)
- 0° or 360° = red
- 120° = green
- 240° = blue, etc ...

<div class="hue-wheel-container">
  <div class="hue-wheel" id="hue-wheel">
<!-- Injected via JS -->
  </div>
  <div class="hue-indicator" id="hue-indicator"></div>
  <div class="hue-value" id="hue-value">Hue: 0°</div>
</div>

_Saturation_ represents the purity of the color:
- it is measured as a percentage (0-100%).
- 0% = grayscale (no color)
- 100% = fully saturated (pure color)
- think of it as how vibrant or muted a color appears: _how much of the hue is present_

<div class="saturation-container">
  <div class="saturation-demo">
    <div class="saturation-gradient" id="saturation-gradient"></div>
    <div class="saturation-indicator" id="saturation-indicator"></div>
  </div>
  <div class="saturation-value" id="saturation-value">Saturation: 100%</div>
  <div class="saturation-sample" id="saturation-sample"></div>
</div>

_Lightness_ represents how light or dark the color is:
- it is measured as a percentage (0-100%)
- 0% = black (no light)
- 50% = pure color
- 100% = white (full light)

<div class="lightness-container">
  <div class="lightness-demo">
    <div class="lightness-gradient" id="lightness-gradient"></div>
    <div class="lightness-indicator" id="lightness-indicator"></div>
  </div>
  <div class="lightness-value" id="lightness-value">Lightness: 50%</div>
  <div class="lightness-sample" id="lightness-sample"></div>
</div>

But this nothing more than a model. Everything we have discussed regarding the _Luminance Contrast Ratio_ and _brightness_ still applies. The _amplitude_ of the light remains the only thing our digital systems are able to act on !
But by using the HSL model, we have a much simpler paradigm. Once we've found a dominant color with better contrast to our background, we can compute that it's X° off our initial color. Therefore, we should shift all other colors in the image by X° as well!

Naturally, there are established ways to convert between RGB and HSL.  
I have copied them below, but I am unable to explain them in details and I accept them without fully understanding them.  
Maybe an opportunity for another article...

#### Conversion from 8-bit RGB to HSL

- Normalize the RGB values:
  \[
  R' = \frac{R}{255}, \quad G' = \frac{G}{255}, \quad B' = \frac{B}{255}
  \]

- Find the minimum and maximum values:
  \[
  C_{\text{min}} = \min(R', G', B')
  \]
  \[
  C_{\text{max}} = \max(R', G', B')
  \]
  \[
  \Delta = C_{\text{max}} - C_{\text{min}}
  \]

- Compute Lightness:
  \[
  L^{\prime} = \frac{C_{\text{max}} + C_{\text{min}}}{2}
  \]

- Compute Saturation:
  \[
  S^{\prime} =
    \begin{cases}
    0 & \text{if } \Delta = 0 \\
    \frac{\Delta}{C_{\text{max}} + C_{\text{min}}} & \text{if } L \leq 0.5 \\
    \frac{\Delta}{2 - C_{\text{max}} - C_{\text{min}}} & \text{if } L > 0.5
    \end{cases}
  \]

- Compute Hue:
  \[
  H^{\prime} =
    \begin{cases}
    \text{undefined} & \text{if } \Delta = 0 \\
    60^\circ \times \left( \frac{G' - B'}{\Delta} + (6 \text{ if } G' < B' \text{ else } 0) \right) & \text{if } C_{\text{max}} = R' \\
    60^\circ \times \left( \frac{B' - R'}{\Delta} + 2 \right) & \text{if } C_{\text{max}} = G' \\
    60^\circ \times \left( \frac{R' - G'}{\Delta} + 4 \right) & \text{if } C_{\text{max}} = B'
    \end{cases}
  \]

- Express in standard HSL units:
  \[
  H = H^{\prime} \times 360, \quad S = S^{\prime} \times 100, \quad L = L^{\prime} \times 100
  \]

```javascript:color_convert.js
function rgbToHsl(r, g, b) {
  r /= 255;
  g /= 255;
  b /= 255;

  const max = Math.max(r, g, b);
  const min = Math.min(r, g, b);
  let h, s, l = (max + min) / 2;

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
```

#### Conversion from HSL to 8-bit RGB 

Given HSL values H: 0-360°, S: 0-100%, L: 0-100%:

1. Normalize:
   \[
     S' = \frac{S}{100}, \quad L' = \frac{L}{100}
   \]

2. Compute intermediate values:
   \[
     C = (1 - |2L' - 1|) \times S'
   \]
   \[
     X = C \times (1 - |\frac{(H \bmod 60°)}{60°} - 1|)
   \]
   \[
     m = L' - \frac{C}{2}
   \]

3. Based on the hue H, assign RGB values:
   \[
     (R', G', B') = 
     \begin{cases}
       (C, X, 0) & \text{if } 0° \leq H < 60° \\
       (X, C, 0) & \text{if } 60° \leq H < 120° \\
       (0, C, X) & \text{if } 120° \leq H < 180° \\
       (0, X, C) & \text{if } 180° \leq H < 240° \\
       (X, 0, C) & \text{if } 240° \leq H < 300° \\
       (C, 0, X) & \text{if } 300° \leq H < 360°
     \end{cases}
   \]

4. Adjust with the offset m:
   \[
     R' = R' + m, \quad G' = G' + m, \quad B' = B' + m
   \]

5. Convert back to the 0-255 range:
   \[
     R = \text{round}(R' \times 255), \quad G = \text{round}(G' \times 255), \quad B = \text{round}(B' \times 255)
   \]

```javascript:color_convert.js
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
```

### Building a software solve

Now that we know how to shift color tones in way that feels coherent we have all the pieces, and all we need is to assemble them.

#### Detecting the dominant color

[Color Thief](https://lokeshdhakar.com/projects/color-thief/) does exaclty what we need.

```javascript:color_convert.js
function adjustImagesColors() {
  // Retrieve every image in the document with class 'dynamic-colors'
  const colorThief = new ColorThief();
  const images = document.querySelectorAll("img.dynamic-colors");

  for (const image of images) {
    // Returns an array of all the colors used within the image.
    // The first color in the array is the dominant color.
    // [[R, G, B], [R, G, B] ...]
    const colors = colorThief.getPalette(image);
  }
}
```

#### Detecting whether or not the dominant color contrasts enough with our background

First, we need to be able to detect our background color. The theme picker inject it as the `--bg` css variable, which we can access as follows.

```javascript
const rootStyles = window.getComputedStyle(document.documentElement);
const bgColor = rootStyles.getPropertyValue("--bg").trim();
```

This either returns a string with color in hex format: `#RRGGBB` or the name of the color such as `"white"`.  
We need a simple parsing function to have them available in the same format as what _Color Thief_ returns.

```javascript:color_convert.js
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
```

We can simply use our already existing functions to detect whether or not the contrast is sufficient !

```javascript:color_convert.js
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

    const contrastRatio = computeContrastRatio(dominantColor, bgColor);

    // We aim for AAA level at the minimum
    if (contrastRatio < 7.0) {
      // Contrast is not high enough, we need to adjust the colors !
    }
  }
}
```

#### Finding a color with a better contrast

We know we need to find a color with a better contrast.
One of the simplest solution is to probe various colors around the color wheel.
We can do this using a gradient ascent:
- pick a starting point on the color wheel
- try colors in both clockwise and counterclockwise directions of the wheel
- follow which ever direction yields the best improvements
- stop when we've found a color with a suitable contrast or when the improvements have stopped

To pick our starting point, we can use a very simple approximation: if the background is light use a complementary color, otherwise start of the original color hue

```javascript:color_convert.js
function findOptimalHueShift(dominantColor, bgColor) {
  const [originalHue, saturation, lightness] = rgbToHsl(
    dominantColor.r,
    dominantColor.g,
    dominantColor.b,
  );

  const originalContrast = computeContrastRatio(dominantColor, bgColor);

  // Calculate background luminance to inform starting point
  const bgLuminance = computeLuminance(bgColor);

  // Start with complementary color if background is light,
  // or keep similar hue if background is dark
  let currentHue = bgLuminance > 0.5 ? (originalHue + 180) % 360 : originalHue;

  let currentContrast = computeContrastWithHue(
    saturation,
    lightness,
    currentHue,
    bgColor,
  );

  let bestHue = currentHue;
  let bestContrast = currentContrast;

  // Start with larger steps, then refine
  let stepSize = 30; // Start with 30° steps
  let improving = true;

  while (improving && stepSize >= 1) {
    improving = false;

    // Try clockwise
    const clockwiseHue = (currentHue + stepSize) % 360;
    const clockwiseContrast = computeContrastWithHue(
      saturation,
      lightness,
      clockwiseHue,
      bgColor,
    );

    // Try counterclockwise
    const counterClockwiseHue = (currentHue - stepSize + 360) % 360;
    const counterClockwiseContrast = computeContrastWithHue(
      saturation,
      lightness,
      counterClockwiseHue,
      bgColor,
    );

    // Determine which direction is better
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
      // If we didn't improve, reduce step size and try again
      stepSize = Math.floor(stepSize / 2);
    }

    if (bestContrast >= 7.0) {
      break;
    }
  }

  // Calculate the actual shift amount (how much we moved from original)
  let hueShift = bestHue - originalHue;

  // Normalize to -180 to 180 range for minimal shift
  if (hueShift > 180) hueShift -= 360;
  if (hueShift < -180) hueShift += 360;

  // If there's less than 20% improvement, return zero shift
  if (bestContrast < originalContrast * 1.2) {
    hueShift = 0;
  }

  return hueShift;
}

function computeContrastWithHue(saturation, lightness, newHue, bgColor) {
  const [r, g, b] = hslToRgb(newHue, saturation, lightness);
  return computeContrastRatio({ r, g, b }, bgColor);
}
```

#### Apply the hue shift to our entire image

```javascript:color_convert.js
function applyHueShiftToImage(image, hueShift) {
  // Skip if no shift needed
  if (hueShift === 0) return;

  // Create a canvas element
  const canvas = document.createElement("canvas");
  const ctx = canvas.getContext("2d");

  // Set canvas dimensions to match the image
  canvas.width = image.naturalWidth;
  canvas.height = image.naturalHeight;

  // Draw the image on the canvas
  ctx.drawImage(image, 0, 0);

  // Get the image data
  const imageData = ctx.getImageData(0, 0, canvas.width, canvas.height);
  const data = imageData.data;

  // Process each pixel
  for (let i = 0; i < data.length; i += 4) {
    const r = data[i];
    const g = data[i + 1];
    const b = data[i + 2];
    const a = data[i + 3];

    // Skip transparent pixels
    if (a === 0) continue;

    // Convert to HSL
    const [h, s, l] = rgbToHsl(r, g, b);

    // Apply hue shift
    const newHue = (h + hueShift) % 360;

    // Convert back to RGB
    const [newR, newG, newB] = hslToRgb(newHue, s, l);

    // Update pixel data
    data[i] = newR;
    data[i + 1] = newG;
    data[i + 2] = newB;
  }

  // Put the modified image data back on the canvas
  ctx.putImageData(imageData, 0, 0);

  // Replace the image source with the canvas data
  const dataURL = canvas.toDataURL();

  // Store the original source for potential toggling
  if (!image.dataset.originalSrc) {
    image.dataset.originalSrc = image.src;
  }

  // Update the image
  image.src = dataURL;
}
```

All that is left is hooking our function to events emitted by the theme picker !


## The end

I learned a lot writing this article but I am not entirely satisfied.  
There are too many formulas that I am unable to breakdown and understand.
There are many rabbit holes that I have not crawled.
- How do you apply it to GIFs ?  
- Can we force the color shift to better match secondary and accent colors of theme ?
- We only shift hue, but what about shifting lightness and saturation ?
- Can we optimize performance by memoizing color shifts and building some kind of lookup table ?

This implementation is naive, amateurish and shows just how ignorant I am. But it is mine, and it will be good enough to adjust the colors of diagrams and schematics or even silly thank you notes written on a color that should be unreadable on a white background !


<div style="display: flex; align-items: center; justify-content: center; gap: 20px;">
  <div id="conclusion-theme-picker" class="theme-toggle"></div>
  <img crossorigin="anonymous" class="dynamic-colors" src="./images/thank-you-for-reading.png" width="300px">
</div>

[^1]: In physics, intensity is proportional to the square of the amplitude of a wave. Luminance is a measure of this intensity per unit area, so higher intensity light results in higher luminance.
[^2]: I chose to illustrate the concept using OLED & AMOLED because they are probably the easiest to gloss over. Other types of screens work very differently.
[^3]: Bit of an oversimplification. Amplitude is linked to intensity which in turn is linked to brightness but saying 'we are increasing brightness' is representative enough for the sake of this article.
[^4]: Using a standard gamma of about 2.2 which is often used as an approximation of the of sRGB
