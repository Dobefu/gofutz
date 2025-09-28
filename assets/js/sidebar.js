"use strict";

(() => {
  function handleGofutzInit() {
    const testFilesContainer = document.querySelector(".sidebar__tests");

    if (!testFilesContainer) {
      console.error("Could not find test files container");

      return;
    }

    renderCoverage();
    updateRunButtonState();
    initializeSortOption();

    testFilesContainer.innerHTML = "";

    const urlParams = new URLSearchParams(globalThis.location.search);
    const sortOption = urlParams.get("sort") || "name-asc";
    const files = sortFiles(globalThis.testData.files, sortOption);

    for (const file of files) {
      renderTestFile(file, testFilesContainer);
    }

    handleGofutzToggleFile();
  }

  function handleGofutzUpdate() {
    const testFilesContainer = document.querySelector(".sidebar__tests");

    if (!testFilesContainer) {
      console.error("Could not find test files container");

      return;
    }

    renderCoverage();
    updateRunButtonState();

    if (!globalThis.testData.files) {
      return;
    }

    const urlParams = new URLSearchParams(globalThis.location.search);
    const sortOption = urlParams.get("sort") || "name-asc";
    const files = sortFiles(globalThis.testData.files, sortOption);

    testFilesContainer.innerHTML = "";

    for (const file of files) {
      renderTestFile(file, testFilesContainer);
    }

    handleGofutzToggleFile();
  }

  function handleGofutzToggleFile() {
    const fileName =
      globalThis.testData.files[
        decodeURIComponent(globalThis.location.hash.slice(1))
      ]?.name ?? "";

    /** @type {HTMLElement | null} */
    const testFilesContainer = document.querySelector(".sidebar__tests");

    if (!testFilesContainer) {
      console.error("Could not find test files container");

      return;
    }

    /** @type {NodeListOf<HTMLElement> | null} */
    const files = testFilesContainer.querySelectorAll(".sidebar__tests--file");

    if (!files) {
      console.error("Could not find any files in the test files container");

      return;
    }

    for (const file of files) {
      if (file.dataset.name !== fileName) {
        file.classList.remove("open");

        continue;
      }

      if (file.classList.contains("open")) {
        file.classList.remove("open");
        globalThis.location.hash = "";

        continue;
      }

      file.classList.add("open");
    }
  }

  function updateRunButtonState() {
    /** @type {HTMLButtonElement | null} */
    const runBtn = document.querySelector(".btn__run-tests");

    if (!runBtn) {
      console.error("Could not find run button");

      return;
    }

    if (globalThis.testData.isRunning) {
      runBtn.classList.add("running");
      runBtn.title = "Stop tests";
    } else {
      runBtn.classList.remove("running");
      runBtn.title = "Run all tests";
    }

    if (navigator.userAgent.includes("Mac")) {
      runBtn.title += " (Cmd+Enter)";

      return;
    }

    runBtn.title += " (Ctrl+Enter)";
  }

  function renderCoverage() {
    /** @type {Element | null} */
    const coverageContainer = document.querySelector(".coverage");

    if (!coverageContainer) {
      console.error("Could not find coverage container");

      return;
    }

    if (globalThis.testData.coverage < 0) {
      coverageContainer.textContent = "…%";

      return;
    }

    coverageContainer.textContent = `${globalThis.testData.coverage.toFixed(1)}%`;
  }

  function initializeSortOption() {
    /** @type {HTMLSelectElement | null} */
    const sortSelect = document.querySelector(".btn__sort-tests");

    if (!sortSelect) {
      console.error("Could not find sort select element");

      return;
    }

    const urlParams = new URLSearchParams(globalThis.location.search);
    const sortParam = urlParams.get("sort");

    if (sortParam) {
      sortSelect.value = sortParam;
    }

    sortSelect.addEventListener("change", () => {
      const url = new URL(globalThis.location.href);

      url.searchParams.set("sort", sortSelect.value);
      globalThis.history.replaceState({}, "", url);

      const testFilesContainer = document.querySelector(".sidebar__tests");

      if (!testFilesContainer) {
        console.error("Could not find test files container");

        return;
      }

      if (!globalThis.testData.files) {
        return;
      }

      const files = sortFiles(globalThis.testData.files, sortSelect.value);
      testFilesContainer.innerHTML = "";

      for (const file of files) {
        renderTestFile(file, testFilesContainer);
      }
    });
  }

  /**
   * @param {Record<string, File>} files
   * @param {string} sortOption
   *
   * @returns {File[]}
   */
  function sortFiles(files, sortOption) {
    switch (sortOption) {
      case "name-asc":
        return Object.values(files).sort((a, b) => {
          return a.name.localeCompare(b.name);
        });

      case "name-desc":
        return Object.values(files).sort((a, b) => {
          return b.name.localeCompare(a.name);
        });

      case "coverage-asc":
        return Object.values(files).sort((a, b) => {
          return a.coverage - b.coverage;
        });

      case "coverage-desc":
        return Object.values(files).sort((a, b) => {
          return b.coverage - a.coverage;
        });

      default:
        return Object.values(files);
    }
  }

  /**
   * @param {File} file
   * @param {Element} testFilesContainer
   */
  function renderTestFile(file, testFilesContainer) {
    const fileItem = document.createElement("div");
    fileItem.classList.add("sidebar__tests--file");
    fileItem.dataset.name = file.name;
    fileItem.addEventListener("click", () => {
      /** @type {HTMLDivElement | null} */
      const mainContentContainer = document.querySelector("#main-content");

      // If the main content container is not visible, disable file toggling.
      if (mainContentContainer && mainContentContainer.clientWidth <= 0) {
        return;
      }

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

    buildFileContent(file, fileItem);
    testFilesContainer.appendChild(fileItem);
  }

  /**
   * @param {File} file
   * @param {Element} fileItem
   */
  function buildFileContent(file, fileItem) {
    const fileItemContainer = document.createElement("div");
    fileItemContainer.classList.add("sidebar__tests--file-container");

    const fileItemStatus = document.createElement("div");
    fileItemStatus.classList.add(`sidebar__tests--file-status status-${file.status}`);
    fileItemContainer.appendChild(fileItemStatus);

    const fileItemTitle = document.createElement("div");
    fileItemTitle.classList.add("sidebar__tests--file-title");
    fileItemTitle.textContent = file.name;
    fileItemTitle.title = file.name;
    fileItemContainer.appendChild(fileItemTitle);

    const fileItemCoverage = document.createElement("div");
    fileItemCoverage.classList.add("sidebar__tests--file-coverage");
    if (file.coverage >= 0) {
      fileItemCoverage.textContent = `${file.coverage.toFixed(1)}%`;
    } else {
      fileItemCoverage.textContent = "…%";
    }
    fileItemContainer.appendChild(fileItemCoverage);

    fileItem.appendChild(fileItemContainer);

    const testsContainer = document.createElement("ul");
    testsContainer.classList.add("sidebar__tests--tests");

    for (const func of file.functions) {
      const funcItemContainer = document.createElement("li");
      funcItemContainer.classList.add("sidebar__tests--test");
      testsContainer.appendChild(funcItemContainer);

      const funcItem = document.createElement("div");
      funcItem.classList.add("sidebar__tests--test-item");
      funcItem.textContent = func.name;
      funcItem.title = func.name;
      funcItemContainer.appendChild(funcItem);

      const funcItemCoverage = document.createElement("div");
      funcItemCoverage.classList.add("sidebar__tests--test-coverage");
      if (func.result.coverage >= 0) {
        funcItemCoverage.textContent = `${func.result.coverage.toFixed(1)}%`;
      } else {
        funcItemCoverage.textContent = "…%";
      }
      funcItemContainer.appendChild(funcItemCoverage);
    }

    fileItem.appendChild(testsContainer);
  }

  globalThis.addEventListener("gofutz:init", handleGofutzInit);
  globalThis.addEventListener("gofutz:update", handleGofutzUpdate);
  globalThis.addEventListener("gofutz:toggle-file", handleGofutzToggleFile);

  /** @type {NodeListOf<HTMLButtonElement> | null} */
  const btnRunAllTests = document.querySelectorAll(".btn__run-tests");

  if (!btnRunAllTests.length) {
    console.error("Could not find any buttons to run all tests");

    return;
  }

  for (const btn of btnRunAllTests) {
    btn.addEventListener("click", () => {
      if (btn.classList.contains("running")) {
        globalThis.dispatchEvent(new CustomEvent("gofutz:stop-tests"));

        return;
      }

      globalThis.dispatchEvent(new CustomEvent("gofutz:run-all-tests"));
    });
  }

  document.addEventListener("keydown", (e) => {
    if ((e.ctrlKey || e.metaKey) && e.key === "Enter") {
      e.preventDefault();
      /** @type {HTMLButtonElement | null} */
      const runBtn = document.querySelector(".btn__run-tests");

      if (!runBtn) {
        console.error("Could not find run button");

        return;
      }

      if (runBtn.classList.contains("running")) {
        globalThis.dispatchEvent(new CustomEvent("gofutz:stop-tests"));

        return;
      }

      globalThis.dispatchEvent(new CustomEvent("gofutz:run-all-tests"));
    }
  });
})();
