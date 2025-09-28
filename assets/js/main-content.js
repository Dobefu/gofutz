"use strict";

(() => {
  function handleGofutzInit() {
    if (globalThis.location.hash) {
      const fileName = decodeURIComponent(globalThis.location.hash.slice(1));
      const file = globalThis.testData.files[fileName];

      if (file) {
        /** @type {HTMLDivElement | null} */
        const mainContentContainer = document.querySelector("#main-content");

        if (mainContentContainer) {
          renderFileContent(file);
        }
      }
    } else {
      showDashboard();
    }
  }

  function updateDashboard() {
    /** @type {HTMLDivElement | null} */
    const dashboardContainer = document.querySelector("#dashboard-container");

    if (!dashboardContainer) {
      console.error("Could not find dashboard container");

      return;
    }

    const coverage = globalThis.testData.coverage;
    const numFiles = Object.keys(globalThis.testData.files).length;
    const coveragePercentage = coverage > 0 ? coverage.toFixed(1) : "…";
    const isRunning = globalThis.testData.isRunning;

    /** @type {HTMLElement | null} */
    const totalTestsElement = dashboardContainer.querySelector(".stat__value--total");
    /** @type {HTMLElement | null} */
    const coverageElement = dashboardContainer.querySelector(".stat__value--coverage");
    /** @type {HTMLElement | null} */
    const statusElement = dashboardContainer.querySelector(
      ".status__running, .status__idle",
    );

    if (totalTestsElement) {
      totalTestsElement.textContent = numFiles.toString();
    }

    if (coverageElement) {
      coverageElement.textContent = `${coveragePercentage}%`;
    }

    if (!statusElement) {
      console.error("Could not find status element");

      return;
    }

    if (isRunning) {
      statusElement.textContent = "🔄 Running...";
    } else {
      statusElement.textContent = "⏸️ Idle";
    }
  }

  function showDashboard() {
    /** @type {HTMLDivElement | null} */
    const dashboardContainer = document.querySelector("#dashboard-container");

    /** @type {HTMLDivElement | null} */
    const fileContainer = document.querySelector("#file-container");

    if (dashboardContainer) {
      dashboardContainer.style.display = "";
    }

    if (fileContainer) {
      fileContainer.innerHTML = "";
    }

    updateDashboard();
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
   */
  function renderFileContent(file) {
    const code = file.highlightedCode
      .split("\n")
      .map((line, idx) => {
        const processedLine = line.replaceAll(/^<\/span>/g, "");
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

    /** @type {HTMLDivElement | null} */
    const dashboardContainer = document.querySelector("#dashboard-container");

    /** @type {HTMLDivElement | null} */
    const fileContainer = document.querySelector("#file-container");

    if (dashboardContainer) {
      dashboardContainer.style.display = "none";
    }

    if (fileContainer) {
      fileContainer.innerHTML = "";

      const filePathHeader = document.createElement("div");
      filePathHeader.classList.add("file-path");
      filePathHeader.textContent = file.name;
      fileContainer.appendChild(filePathHeader);

      const codeContainer = document.createElement("pre");
      codeContainer.classList.add("main-content__code");
      codeContainer.dataset.file = file.name;
      codeContainer.innerHTML = code;
      fileContainer.appendChild(codeContainer);
    }
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
      /** @type {HTMLDivElement | null} */
      const dashboardContainer = document.querySelector("#dashboard-container");

      /** @type {HTMLDivElement | null} */
      const fileContainer = document.querySelector("#file-container");

      if (dashboardContainer) {
        dashboardContainer.style.display = "";
      }

      if (fileContainer) {
        fileContainer.innerHTML = "";
      }

      return;
    }

    renderFileContent(file);
  }

  function handleGofutzUpdate() {
    updateDashboard();

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
        renderFileContent(file);

        break;
      }
    }
  }

  globalThis.addEventListener("gofutz:init", handleGofutzInit);
  globalThis.addEventListener("gofutz:toggle-file", handleGofutzToggleFile);
  globalThis.addEventListener("gofutz:update", handleGofutzUpdate);
})();
