"use strict";

/**
 * @param {CustomEvent} e
 */
function handleGofutzInit(e) {
  /** @type {Message} */
  const details = e.detail;
  const testFilesContainer = document.querySelector(".sidebar__tests");

  if (!testFilesContainer) {
    return;
  }

  testFilesContainer.innerHTML = "";

  for (const file of Object.values(details.params.files)) {
    const fileItem = document.createElement("details");
    const fileItemSummary = document.createElement("summary");
    fileItemSummary.textContent = file.name;
    fileItem.appendChild(fileItemSummary);

    const testsContainer = document.createElement("ul");

    for (const test of file.tests) {
      const testItem = document.createElement("li");

      testItem.textContent = test.name;
      testsContainer.appendChild(testItem);
    }

    fileItem.appendChild(testsContainer);
    testFilesContainer.appendChild(fileItem);
  }
}

(() => {
  window.addEventListener("gofutz:init", handleGofutzInit);
})();
