"use strict";

(() => {
  /**
   * @param {CustomEvent} e
   */
  function handleGofutzInit(e) {
    /** @type {Message} */
    const details = e.detail;
    const testFilesContainer = document.querySelector(".sidebar__tests");

    if (!testFilesContainer) {
      console.error("Could not find test files container");

      return;
    }

    renderCoverage(details);
    updateRunButtonState(details.params.isRunning);

    testFilesContainer.innerHTML = "";

    for (const file of Object.values(details.params.files)) {
      renderTestFile(file, testFilesContainer);
    }
  }

  /**
   * @param {CustomEvent} e
   */
  function handleGofutzUpdate(e) {
    /** @type {Message} */
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

    for (const file of Object.values(details.params.files)) {
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
   * @param {Message} details
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
    fileItemTitle.addEventListener("click", () => {
      window.dispatchEvent(
        new CustomEvent("gofutz:toggle-file", { detail: file }),
      );
    });
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
      const funcItem = document.createElement("li");
      funcItem.classList.add("sidebar__tests--test");
      funcItem.textContent = func.name;
      funcItem.title = func.name;
      testsContainer.appendChild(funcItem);

      const funcItemCoverage = document.createElement("div");
      funcItemCoverage.classList.add("sidebar__tests--test-coverage");
      if (func.result.coverage >= 0) {
        funcItemCoverage.textContent = `${func.result.coverage.toFixed(1)}%`;
      } else {
        funcItemCoverage.textContent = "…%";
      }
      funcItem.appendChild(funcItemCoverage);
    }

    fileItem.appendChild(testsContainer);
  }

  window.addEventListener("gofutz:init", handleGofutzInit);
  window.addEventListener("gofutz:update", handleGofutzUpdate);

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
