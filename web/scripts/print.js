function setToPrintableDimensions() {
  const resumeContainer = document.querySelector(".resume-container");
  resumeContainer.style.setProperty("width", "800px");
  resumeContainer.style.setProperty("font-size", "8pt");

  const resumeActions = document.querySelector(".resume-actions");
  resumeActions.remove();

  const resume = document.querySelector(".resume");
  resume.style.setProperty("border", "none");

  const img = document.querySelector(".resume img");
  img.style.setProperty("width", "100px");
  img.style.setProperty("height", "100px");
}
