window.onload = function () {
  initComponents();
  showCurrentRuns();
};

/**
 * Loads and shows the latest runs in the options page
 */
async function showCurrentRuns() {
  const runs = await getCurrentRuns();
  const container = document.getElementById("runsListContainer");

  runs.forEach(async (runStartTime) => {
    const row = document.createElement("tr");
    row.classList.add("mdc-data-table__header-row");

    const cell1 = document.createElement("td");
    cell1.setAttribute("data-run-id", runStartTime);
    cell1.innerText = getFormattedTime(new Date(runStartTime));
    cell1.classList.add("mdc-data-table__cell");

    const cell2 = document.createElement("td");
    const fastestRun = await getFastestPingInRun(runStartTime);
    cell2.innerText = `${fastestRun.region} (${fastestRun.ping} ms)`;
    cell2.setAttribute("data-run-id", runStartTime);
    cell2.classList.add("mdc-data-table__cell");

    row.appendChild(cell1);
    row.appendChild(cell2);

    container.appendChild(row);
  });
}

/**
 * Returns the information for the fastest region for a saved run recognized by a runId
 * @param {int} runId
 * @return {Object}
 */
async function getFastestPingInRun(runId) {
  const data = await getRunData(runId);
  let fastestRegion = null;

  // loop therough the results of the run to find the fastest region
  for (const region of Object.keys(data.results)) {
    if (
      fastestRegion === null ||
      data.results[fastestRegion] > data.results[region]
    ) {
      fastestRegion = region;
    }
  }

  return {
    region: fastestRegion,
    ping: data.results[fastestRegion],
    runId: runId,
  };
}

/**
 * Fetches the current stored runs from chrome.storage.local
 * @return {Promise}
 */
async function getCurrentRuns() {
  return new Promise((resolve, reject) => {
    chrome.storage.local.get("runs", (result) => {
      // return an empty array by default
      resolve(result["runs"] ?? []);
    });
  });
}

/**
 * Fetches the details of a single historical run from chrome.storage.local
 * @param {object} runId
 */
async function getRunData(runId) {
  runId = "run-" + runId;
  return new Promise((resolve, reject) => {
    chrome.storage.local.get(runId, (result) => {
      // return an empty array by default
      resolve(result[runId] ?? false);
    });
  });
}

/**
 * Handler for when a tab is clicked
 * @param {int} selectedIndex
 */
function focusTab(selectedIndex) {
  const tabs = document
    .querySelector(".tabDetails")
    .querySelectorAll(".tabSingle");
  [...tabs].forEach((tab, index) => {
    if (index === selectedIndex) {
      tab.classList.add("tabActive");
    } else {
      tab.classList.remove("tabActive");
    }
  });
}

/**
 * Helper to initialize all the material components
 */
function initComponents() {
  const tabBar = new mdc.tabBar.MDCTabBar(
    document.querySelector(".mdc-tab-bar")
  );

  // Handle tab click
  tabBar[0].addEventListener("MDCTabBar:activated", function (ev) {
    const tabIndex = ev?.detail?.index ?? 0;

    focusTab(tabIndex);
  });
}

/**
 * Helper to return formatted time given a Date object
 * @param {Date} dt
 * @return {string}
 */
function getFormattedTime(dt) {
  return (
    new Intl.DateTimeFormat("en-US", { dateStyle: "long" }).format(dt) +
    ", " +
    new Intl.DateTimeFormat("en-US", { timeStyle: "medium" }).format(dt)
  );
}
