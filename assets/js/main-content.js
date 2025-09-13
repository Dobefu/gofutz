"use strict";

(() => {
  function handleGofutzInit() {
    if (window.location.hash) {
      const fileName = decodeURIComponent(window.location.hash.slice(1));
      const file = globalThis.testData.files[fileName];

      if (file) {
        window.dispatchEvent(
          new CustomEvent("gofutz:toggle-file", { detail: file }),
        );
      }
    }
  }

  /**
   * @param {number} lineNumber
   * @param {Line[]} coveredLines
   * @returns {string}
   */
  function getLineCoverageStatus(lineNumber, coveredLines) {
    const matchingLine = coveredLines.find((line) => {
      return lineNumber >= line.startLine && lineNumber <= line.endLine;
    });

    if (!matchingLine) {
      return "";
    }

    return matchingLine.executionCount > 0 ? "covered" : "uncovered";
  }

  /**
   * @param {File} file
   * @param {HTMLDivElement} mainContentContainer
   */
  function renderFileContent(file, mainContentContainer) {
    const code = file.highlightedCode
      .split("\n")
      .map((line, idx) => {
        const processedLine = line.replace(/^<\/span>/g, "");
        const lineNumber = idx + 1;
        const coverageStatus = getLineCoverageStatus(
          lineNumber,
          file.coveredLines,
        );

        return `<div class="main-content__code--line ${coverageStatus}">
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

  /**
   * @param {CustomEvent} e
   */
  function handleGofutzToggleFile(e) {
    /** @type {File} */
    const file = e.detail;
    /** @type {HTMLDivElement | null} */
    const mainContentContainer = document.querySelector("#main-content");

    if (!mainContentContainer) {
      console.error("Could not find main content container");

      return;
    }

    /** @type {HTMLPreElement | null} */
    const currentCodeContainer = mainContentContainer.querySelector(
      ".main-content__code",
    );

    if (
      currentCodeContainer &&
      currentCodeContainer.dataset.file === file.name
    ) {
      currentCodeContainer.remove();

      return;
    }

    renderFileContent(file, mainContentContainer);
  }

  function handleGofutzUpdate() {
    /** @type {HTMLPreElement | null} */
    const currentCodeContainer = document.querySelector(".main-content__code");
    /** @type {HTMLDivElement | null} */
    const mainContentContainer = document.querySelector("#main-content");

    if (!currentCodeContainer || !mainContentContainer) {
      return;
    }

    const currentFileName = currentCodeContainer.dataset.file;

    if (!globalThis.testData.files) {
      return;
    }

    for (const file of Object.values(globalThis.testData.files)) {
      if (file.name === currentFileName) {
        renderFileContent(file, mainContentContainer);

        break;
      }
    }
  }

  window.addEventListener("gofutz:init", handleGofutzInit);
  window.addEventListener("gofutz:toggle-file", handleGofutzToggleFile);
  window.addEventListener("gofutz:update", handleGofutzUpdate);
})();
