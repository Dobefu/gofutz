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

      return;
    }

    showDashboard();
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
    const totalTestsElement = dashboardContainer.querySelector(
      ".stat__value--total",
    );

    /** @type {HTMLElement | null} */
    const coverageElement = dashboardContainer.querySelector(
      ".stat__value--coverage",
    );

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

    updateHeatmap();
  }

  function updateHeatmap() {
    /** @type {HTMLDivElement | null} */
    const heatmapGrid = document.querySelector("#coverage-heatmap");

    if (!heatmapGrid) {
      return;
    }

    const files = globalThis.testData.files;
    const fileEntries = Object.entries(files);

    heatmapGrid.innerHTML = "";

    for (const [fileName, file] of fileEntries) {
      const cell = document.createElement("div");
      cell.className = "heatmap-cell";
      cell.dataset.fileName = fileName;
      cell.dataset.coverage = file.coverage.toFixed(1);

      const coverage = file.coverage;
      let backgroundColor;

      if (coverage <= 0) {
        backgroundColor = "var(--color-heatmap-none)";
      } else if (coverage < 20) {
        const mixPercentage = (coverage / 20) * 100;
        backgroundColor = `color-mix(in srgb, var(--color-heatmap-none) 100%, var(--color-heatmap-low) ${mixPercentage}%)`;
      } else if (coverage < 50) {
        const mixPercentage = ((coverage - 20) / 30) * 100;
        backgroundColor = `color-mix(in srgb, var(--color-heatmap-low) 100%, var(--color-heatmap-medium) ${mixPercentage}%)`;
      } else if (coverage < 80) {
        const mixPercentage = ((coverage - 50) / 30) * 100;
        backgroundColor = `color-mix(in srgb, var(--color-heatmap-medium) 100%, var(--color-heatmap-high) ${mixPercentage}%)`;
      } else {
        backgroundColor = "var(--color-heatmap-high)";
      }

      cell.style.backgroundColor = backgroundColor;

      cell.addEventListener("mouseenter", (e) => {
        if (!(e.target instanceof HTMLElement)) {
          return;
        }

        showTooltip(e.target, fileName, file.coverage);
      });

      cell.addEventListener("mouseleave", () => {
        hideTooltip();
      });

      cell.addEventListener("click", () => {
        const currentFile = globalThis.testData.files[file.name];

        if (!currentFile) {
          console.error("Could not find file:", file.name);

          return;
        }

        globalThis.location.hash = encodeURIComponent(currentFile.name);

        globalThis.dispatchEvent(
          new CustomEvent("gofutz:toggle-file", { detail: currentFile }),
        );
      });

      heatmapGrid.appendChild(cell);
    }
  }

  /**
   * @param {HTMLElement} element
   * @param {string} fileName
   * @param {number} coverage
   */
  function showTooltip(element, fileName, coverage) {
    hideTooltip();

    const tooltip = document.createElement("div");
    tooltip.className = "heatmap-tooltip";
    tooltip.textContent = `${fileName}: ${coverage.toFixed(1)}%`;

    element.appendChild(tooltip);

    const rect = element.getBoundingClientRect();
    const tooltipRect = tooltip.getBoundingClientRect();
    const rem = Number.parseFloat(
      getComputedStyle(document.documentElement).fontSize,
    );

    let left = rect.left + rect.width / 2 - tooltipRect.width / 2;
    let top = rect.top - tooltipRect.height - rem;

    left = Math.max(
      rem,
      Math.min(left, window.innerWidth - tooltipRect.width - rem),
    );
    top = top < rem ? rect.bottom + rem : top;

    tooltip.style.left = `${left / rem}rem`;
    tooltip.style.top = `${top / rem}rem`;
    tooltip.style.transform = "none";

    tooltip.classList.add("heatmap-tooltip--visible");
  }

  function hideTooltip() {
    /** @type {HTMLElement | null} */
    const existingTooltip = document.querySelector(".heatmap-tooltip");

    if (existingTooltip) {
      existingTooltip.remove();
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
    /** @type {File | undefined} */
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

    /** @type {HTMLDivElement | null} */
    const fileContainer = document.querySelector("#file-container");

    /** @type {HTMLDivElement | null} */
    const dashboardContainer = document.querySelector("#dashboard-container");

    if (
      !file ||
      (currentCodeContainer && currentCodeContainer.dataset.file === file.name)
    ) {
      if (dashboardContainer) {
        dashboardContainer.style.display = "";
      }

      if (fileContainer) {
        fileContainer.innerHTML = "";
        updateDashboard();
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
