"use strict";

/**
 * @param {CustomEvent} e
 */
function handleGofutzToggleFile(e) {
  /** @type {File} */
  const file = e.detail;
  const mainContentContainer = document.querySelector("#main-content");

  if (!mainContentContainer) {
    console.error("Could not find main content container");

    return;
  }

  /** @type {HTMLPreElement | null} */
  const currentCodeContainer = mainContentContainer.querySelector(
    ".main-content__code",
  );

  if (currentCodeContainer && currentCodeContainer.dataset.file === file.name) {
    currentCodeContainer.remove();

    return;
  }

  const coveredLines = new Set();

  for (const line of file.coveredLines) {
    coveredLines.add(line.number);
  }

  const code = file.highlightedCode
    .split("\n")
    .map((line, idx) => {
      const processedLine = line.replace(/^<\/span>/g, "");
      const lineNumber = idx + 1;
      const isCovered = coveredLines.has(lineNumber);
      const coveredClass = isCovered ? "covered" : "uncovered";

      return `<div class="main-content__code--line ${coveredClass}">
        <span class="main-content__code--line-number">${lineNumber}</span>
        <span class="main-content__code--line-content">${processedLine}</span></span>
      </div>`;
    })
    .join("");

  mainContentContainer.innerHTML = "";

  const codeContainer = document.createElement("pre");
  codeContainer.classList.add("main-content__code");
  codeContainer.dataset.file = file.name;
  codeContainer.innerHTML = code;
  mainContentContainer.appendChild(codeContainer);
}

(() => {
  window.addEventListener("gofutz:toggle-file", handleGofutzToggleFile);
})();
