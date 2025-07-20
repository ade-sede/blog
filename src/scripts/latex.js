function waitForKaTeX(callback) {
  if (typeof katex !== "undefined") {
    callback();
  } else {
    setTimeout(() => waitForKaTeX(callback), 50);
  }
}

function renderLatexElements() {
  const displayElements = document.querySelectorAll(".katex-display");
  const inlineElements = document.querySelectorAll(".katex-inline");

  displayElements.forEach(function (element) {
    const latex = element.getAttribute("data-latex");
    if (latex && !element.classList.contains("katex-rendered")) {
      try {
        katex.render(latex, element, {
          displayMode: true,
          throwOnError: false,
        });
        element.classList.add("katex-rendered");
      } catch (e) {
        console.error("KaTeX rendering error:", e);
        element.textContent = "Error rendering LaTeX: " + e.message;
      }
    }
  });

  inlineElements.forEach(function (element) {
    const latex = element.getAttribute("data-latex");
    if (latex && !element.classList.contains("katex-rendered")) {
      try {
        katex.render(latex, element, {
          displayMode: false,
          throwOnError: false,
        });
        element.classList.add("katex-rendered");
      } catch (e) {
        console.error("KaTeX rendering error:", e);
        element.textContent = "Error rendering LaTeX: " + e.message;
      }
    }
  });
}

document.addEventListener("DOMContentLoaded", function () {
  waitForKaTeX(renderLatexElements);
});
