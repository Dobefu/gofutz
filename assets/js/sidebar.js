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
    fileItem.classList.add("sidebar__tests--file");
    fileItem.open = true;

    const fileItemSummary = document.createElement("summary");
    fileItemSummary.classList.add("sidebar__tests--file-summary");
    fileItemSummary.textContent = file.name;
    fileItemSummary.title = file.name;
    fileItemSummary.onclick = (e) => {
      e.preventDefault();

      window.dispatchEvent(
        new CustomEvent("gofutz:toggle-file", { detail: file }),
      );
    };
    fileItem.appendChild(fileItemSummary);

    const testsContainer = document.createElement("ul");
    testsContainer.classList.add("sidebar__tests--tests");

    for (const test of file.tests) {
      const testItem = document.createElement("li");
      testItem.classList.add("sidebar__tests--test");
      testItem.textContent = test.name;
      testItem.title = test.name;

      testsContainer.appendChild(testItem);
    }

    fileItem.appendChild(testsContainer);
    testFilesContainer.appendChild(fileItem);
  }
}

(() => {
  window.addEventListener("gofutz:init", handleGofutzInit);
})();
