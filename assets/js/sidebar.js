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
    initializeHomeButton();
    initializeSortOption();
    initializeSearch();

    renderFiles();
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

    renderFiles();
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

  function initializeHomeButton() {
    /** @type {HTMLButtonElement | null} */
    const homeBtn = document.querySelector(".btn__home");

    if (!homeBtn) {
      console.error("Could not find home button");

      return;
    }

    homeBtn.addEventListener("click", () => {
      globalThis.dispatchEvent(
        new CustomEvent("gofutz:toggle-file", { detail: undefined }),
      );
    });
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

      renderFiles();
    });
  }

  function initializeSearch() {
    /** @type {HTMLInputElement | null} */
    const searchInput = document.querySelector(".sidebar__search input");

    if (!searchInput) {
      console.error("Could not find search input element");

      return;
    }

    const urlParams = new URLSearchParams(globalThis.location.search);
    const searchParam = urlParams.get("q");

    if (searchParam) {
      searchInput.value = searchParam;
    }

    searchInput.addEventListener("input", () => {
      const url = new URL(globalThis.location.href);

      if (searchInput.value) {
        url.searchParams.set("q", searchInput.value);
      } else {
        url.searchParams.delete("q");
      }

      globalThis.history.replaceState({}, "", url);
      renderFiles();
    });

    /** @type {HTMLButtonElement | null} */
    const clearBtn = document.querySelector(".btn__search-clear");

    if (clearBtn) {
      clearBtn.addEventListener("click", () => {
        searchInput.value = "";
        const url = new URL(globalThis.location.href);
        url.searchParams.delete("q");
        globalThis.history.replaceState({}, "", url);
        renderFiles();
      });
    }
  }

  function renderFiles() {
    const testFilesContainer = document.querySelector(".sidebar__tests");

    if (!testFilesContainer) {
      console.error("Could not find test files container");

      return;
    }

    if (!globalThis.testData.files) {
      return;
    }

    const urlParams = new URLSearchParams(globalThis.location.search);
    const sortOption = urlParams.get("sort") || "name-asc";
    const searchQuery = urlParams.get("q")?.toLowerCase() || "";

    let files = sortFiles(globalThis.testData.files, sortOption);

    if (searchQuery) {
      files = files.filter((file) => {
        const fileNameMatch = file.name.toLowerCase().includes(searchQuery);
        const functionMatch = file.functions.some((func) =>
          func.name.toLowerCase().includes(searchQuery),
        );
        return fileNameMatch || functionMatch;
      });
    }

    testFilesContainer.innerHTML = "";

    for (const file of files) {
      renderTestFile(file, testFilesContainer);
    }
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
   * @param {number} coverage
   * @returns {string}
   */
  function renderCoverageCircle(coverage) {
    const strokeDasharray = Math.PI * 2 * 16;
    let strokeDashoffset = strokeDasharray;
    let coverageText = "…%";

    if (coverage >= 0) {
      coverageText = `${coverage.toFixed(1)}%`;

      strokeDashoffset = strokeDasharray - (coverage / 100) * strokeDasharray;
    }

    return `
      <span class="sidebar__coverage-text">${coverageText}</span>
      <div class="sidebar__coverage-circle">
        <svg viewBox="0 0 36 36">
          <circle
            class="sidebar__coverage-circle-bg"
            cx="18"
            cy="18"
            r="16"
          ></circle>
          <circle
            class="sidebar__coverage-circle-fill"
            cx="18"
            cy="18"
            r="16"
            stroke="var(--color-coverage-bar-covered)"
            stroke-dasharray="${strokeDasharray}"
            stroke-dashoffset="${strokeDashoffset}"
          ></circle>
        </svg>
      </div>
    `;
  }

  /**
   * @param {File} file
   * @param {Element} testFilesContainer
   */
  function renderTestFile(file, testFilesContainer) {
    const fileItem = document.createElement("div");
    fileItem.classList.add("sidebar__tests--file");
    fileItem.dataset.name = file.name;
    fileItem.tabIndex = 0;
    fileItem.setAttribute("role", "button");
    fileItem.setAttribute("aria-label", `Toggle ${file.name} test file`);

    const toggleFile = () => {
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
    };

    fileItem.addEventListener("click", toggleFile);
    fileItem.addEventListener("keydown", (e) => {
      if (e.key === "Enter" || e.key === " ") {
        e.preventDefault();
        toggleFile();
      }
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
    fileItemStatus.classList.add(
      "sidebar__tests--file-status",
      `status-${file.status}`,
    );
    fileItemContainer.appendChild(fileItemStatus);

    const fileItemTitle = document.createElement("div");
    fileItemTitle.classList.add("sidebar__tests--file-title");
    fileItemTitle.textContent = file.name;
    fileItemTitle.title = file.name;
    fileItemContainer.appendChild(fileItemTitle);

    const fileItemCoverage = document.createElement("div");
    fileItemCoverage.classList.add("sidebar__tests--file-coverage");
    fileItemCoverage.innerHTML = renderCoverageCircle(file.coverage);

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
      funcItemCoverage.innerHTML = renderCoverageCircle(func.result.coverage);

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
