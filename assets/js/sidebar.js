"use strict";

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

  testFilesContainer.innerHTML = "";

  for (const file of Object.values(details.params.files)) {
    const fileItem = document.createElement("div");
    fileItem.classList.add("sidebar__tests--file");

    const fileItemContainer = document.createElement("div");
    fileItemContainer.classList.add("sidebar__tests--file-container");

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
    fileItemCoverage.textContent = `${file.coverage}%`;
    fileItemContainer.appendChild(fileItemCoverage);

    fileItem.appendChild(fileItemContainer);

    const testsContainer = document.createElement("ul");
    testsContainer.classList.add("sidebar__tests--tests");

    for (const test of file.tests) {
      const testItem = document.createElement("li");
      testItem.classList.add("sidebar__tests--test");
      testItem.textContent = test.name;
      testItem.title = test.name;
      testsContainer.appendChild(testItem);

      const testItemCoverage = document.createElement("div");
      testItemCoverage.classList.add("sidebar__tests--test-coverage");
      testItemCoverage.textContent = `${test.result.coverage}%`;
      testItem.appendChild(testItemCoverage);
    }

    fileItem.appendChild(testsContainer);
    testFilesContainer.appendChild(fileItem);
  }
}

(() => {
  window.addEventListener("gofutz:init", handleGofutzInit);

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
