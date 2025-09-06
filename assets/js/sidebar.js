"use strict";

(() => {
  /**
   * @param {CustomEvent} e
   */
  function handleGofutzInit(e) {
    /** @type {UpdateMessage} */
    const details = e.detail;
    const testFilesContainer = document.querySelector(".sidebar__tests");

    if (!testFilesContainer) {
      console.error("Could not find test files container");

      return;
    }

    renderCoverage(details);
    updateRunButtonState(details.params.isRunning);

    testFilesContainer.innerHTML = "";

    const files = Object.values(details.params.files).sort((a, b) => {
      return a.name.localeCompare(b.name);
    });

    for (const file of files) {
      renderTestFile(file, testFilesContainer);
    }
  }

  /**
   * @param {CustomEvent} e
   */
  function handleGofutzUpdate(e) {
    /** @type {UpdateMessage} */
    const details = e.detail;
    const testFilesContainer = document.querySelector(".sidebar__tests");

    if (!testFilesContainer) {
      console.error("Could not find test files container");

      return;
    }

    renderCoverage(details);
    updateRunButtonState(details.params.isRunning);

    if (!details.params.files) {
      return;
    }

    const files = Object.values(details.params.files).sort((a, b) => {
      return a.name.localeCompare(b.name);
    });

    for (const file of files) {
      const fileItem = testFilesContainer.querySelector(
        `[data-name="${file.name}"]`,
      );

      if (fileItem) {
        fileItem.innerHTML = "";
        buildFileContent(file, fileItem);

        continue;
      }

      renderTestFile(file, testFilesContainer);
    }
  }

  /**
   * @param {CustomEvent} e
   */
  function handleGofutzToggleFile(e) {
    const fileName = e.detail.name;
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
      const hashFileName = decodeURIComponent(window.location.hash.slice(1));

      if (
        file.dataset.name !== fileName &&
        file.dataset.name !== hashFileName
      ) {
        file.classList.remove("open");

        continue;
      }

      if (file.classList.contains("open")) {
        file.classList.remove("open");

        continue;
      }

      file.classList.add("open");
    }
  }

  /**
   * @param {boolean} isRunning
   */
  function updateRunButtonState(isRunning) {
    /** @type {HTMLButtonElement | null} */
    const runBtn = document.querySelector(".btn__run-tests");

    if (!runBtn) {
      console.error("Could not find run button");

      return;
    }

    runBtn.disabled = isRunning;
  }

  /**
   * @param {UpdateMessage} details
   */
  function renderCoverage(details) {
    /** @type {Element | null} */
    const coverageContainer = document.querySelector(".coverage");

    if (!coverageContainer) {
      console.error("Could not find coverage container");

      return;
    }

    if (details.params.coverage < 0) {
      coverageContainer.textContent = "…%";

      return;
    }

    coverageContainer.textContent = `${details.params.coverage.toFixed(1)}%`;
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

      window.location.hash = encodeURIComponent(file.name);

      window.dispatchEvent(
        new CustomEvent("gofutz:toggle-file", { detail: file }),
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
    fileItemStatus.classList.add("sidebar__tests--file-status");
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

  window.addEventListener("gofutz:init", handleGofutzInit);
  window.addEventListener("gofutz:update", handleGofutzUpdate);
  window.addEventListener("gofutz:toggle-file", handleGofutzToggleFile);

  const btnRunAllTests = document.querySelectorAll(".btn__run-tests");

  if (!btnRunAllTests.length) {
    console.error("Could not find any buttons to run all tests");

    return;
  }

  for (const btn of btnRunAllTests) {
    btn.addEventListener("click", () => {
      window.dispatchEvent(new CustomEvent("gofutz:run-all-tests"));
    });
  }
})();
