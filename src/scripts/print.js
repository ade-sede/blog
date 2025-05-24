function setToPrintableDimensions() {
  const resumeContainer = document.querySelector(".resume-container");
  resumeContainer.style.setProperty("width", "800px");
  resumeContainer.style.setProperty("font-size", "7pt");

  const resumeActions = document.querySelector(".resume-actions");
  resumeActions.remove();

  const resume = document.querySelector(".resume");
  resume.style.setProperty("border", "none");

  const header = document.querySelector(".resume .header");
  header.style.setProperty("padding", "1rem");
  header.style.setProperty("gap", "1rem");

  const img = document.querySelector(".resume img");
  img.style.setProperty("width", "80px");
  img.style.setProperty("height", "80px");

  const name = document.querySelector(".resume .info h1.name");
  name.style.setProperty("font-size", "1.5rem");

  const sections = document.querySelectorAll(".resume .section");
  sections.forEach(section => {
    section.style.setProperty("padding-left", "1rem");
    section.style.setProperty("padding-right", "1rem");
  });

  const sectionTitles = document.querySelectorAll(".resume .section h1");
  sectionTitles.forEach(title => {
    title.style.setProperty("font-size", "1.2rem");
    title.style.setProperty("margin", "0.5rem 0");
  });

  const hrs = document.querySelectorAll(".resume hr");
  hrs.forEach(hr => {
    hr.style.setProperty("margin", "0.5rem 0");
  });

  const entries = document.querySelectorAll(".resume .entry");
  entries.forEach(entry => {
    entry.style.setProperty("margin-bottom", "0.5rem");
  });

  const links = document.querySelector(".resume .section .links");
  if (links) {
    links.style.setProperty("padding", "0.5rem");
    links.style.setProperty("gap", "2rem");
  }
}
