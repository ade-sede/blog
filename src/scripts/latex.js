document.addEventListener("DOMContentLoaded", function() {
  const displayElements = document.querySelectorAll(".katex-display");
  const inlineElements = document.querySelectorAll(".katex-inline");
  
  displayElements.forEach(function(element) {
    const latex = element.getAttribute("data-latex");
    if (latex) {
      try {
        katex.render(latex, element, {
          displayMode: true,
          throwOnError: false
        });
      } catch (e) {
        console.error("KaTeX rendering error:", e);
        element.textContent = "Error rendering LaTeX: " + e.message;
      }
    }
  });
  
  inlineElements.forEach(function(element) {
    const latex = element.getAttribute("data-latex");
    if (latex) {
      try {
        katex.render(latex, element, {
          displayMode: false,
          throwOnError: false
        });
      } catch (e) {
        console.error("KaTeX rendering error:", e);
        element.textContent = "Error rendering LaTeX: " + e.message;
      }
    }
  });
});