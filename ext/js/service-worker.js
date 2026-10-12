const CHROME_ALARM_ID = "gcping_endpoint_alarm";
const CHROME_STORAGE_ENDPOINTS_KEY = "gcping_endpoints";
const PING_STATUS_RUNNING = "running";
const PING_STATUS_NOT_RUNNING = "not running";

const currentStatus = {
  status: PING_STATUS_NOT_RUNNING,
  completed: 0,
  total: 0,
};

// when the extension is installed, add an alarm to refresh our endpoints
chrome.runtime.onInstalled.addListener((details) => {
  if (details.reason === chrome.runtime.OnInstalledReason.INSTALL) {
    // Create an alarm to run every hour without any delay
    chrome.alarms.create(CHROME_ALARM_ID, {
      delayInMinutes: 0,
      periodInMinutes: 60,
    });
  }
});

/**
 * Event listener for the alarm
 */
chrome.alarms.onAlarm.addListener(function (alarm) {
  if (alarm.name === CHROME_ALARM_ID) {
    fetchAndSaveEndpoints();
  }
});

/**
 * Event listener on click on the extension's action
 */
chrome.action.onClicked.addListener(pingAllRegions);

/**
 * Message received from other parts of the extension
 */
chrome.runtime.onMessage.addListener(function (request, sender, sendResponse) {
  if (request.action === "fetch_current_status") {
    fetchCurrentStatus().then((data) => {
      sendResponse(data);
    });
  } else if (request.action === "run_test") {
    pingAllRegions();
  } else if (request.action === "stop_test") {
    stopRunningTest();
  }

  return true;
});

/**
 * Helper to fetch the different endpoints that we need to ping
 * and save it in the chrome localstorage
 */
async function fetchAndSaveEndpoints() {
  return new Promise((resolve, reject) => {
    fetch("https://gcping.com/api/endpoints")
      .then((resp) => resp.json())
      .then((endpoints) => {
        const regions = {};

        for (const zone of Object.values(endpoints)) {
          const gcpZone = {
            key: zone.Region,
            label: zone.RegionName,
            pingUrl: zone.URL + "/api/ping",
            latencies: [],
            median: "",
          };

          regions[gcpZone.key] = gcpZone;
        }

        const data = {};
        data[CHROME_STORAGE_ENDPOINTS_KEY] = regions;

        chrome.storage.local.set(data);
        resolve();
      });
  });
}

/**
 * Ping all regions to get results
 */
async function pingAllRegions() {
  // Don't do anything if the test is already running
  if (currentStatus.status === PING_STATUS_RUNNING) {
    return;
  }

  let regions = await getRegionsToPing();

  // fallback in case the regions have never been fetched
  if (!regions) {
    await fetchAndSaveEndpoints();
    regions = await getRegionsToPing();
  }

  const numRegions = Object.keys(regions).length;
  let counter = 1;
  const results = {};
  let fastestRegion;
  const runData = {
    startTime: Date.now(),
  };

  currentStatus.status = PING_STATUS_RUNNING;
  currentStatus.total = numRegions;

  chrome.action.setBadgeText({ text: `0/${numRegions}` });

  for (const region of Object.values(regions)) {
    const ping = await pingSingleRegion(region["pingUrl"]);

    results[region["key"]] = ping;
    if (fastestRegion === undefined || ping < results[fastestRegion]) {
      fastestRegion = region["key"];
    }

    // This may be changed in b/w test by the stopRunningTest func
    // called from the options page
    if (currentStatus.status !== PING_STATUS_RUNNING) {
      return;
    }

    chrome.action.setBadgeText({ text: `${counter}/${numRegions}` });
    currentStatus.completed = counter;
    counter++;
    syncCurrentStatus();
  }

  currentStatus.status = PING_STATUS_NOT_RUNNING;
  syncCurrentStatus();
  chrome.action.setBadgeText({ text: "" });
  displayPingResults(fastestRegion, results[fastestRegion]);

  runData["endTime"] = Date.now();
  runData["results"] = results;

  await saveRunData(runData);
}

/**
 * Helper to stop the current running test
 */
function stopRunningTest() {
  currentStatus.status = PING_STATUS_NOT_RUNNING;
  currentStatus.completed = 0;
  currentStatus.total = 0;

  chrome.action.setBadgeText({ text: "" });
  syncCurrentStatus();
}

/**
 * Helper function to ping a single URL and return the result
 * @param {string} url
 */
async function pingSingleRegion(url) {
  return new Promise((resolve) => {
    const start = new Date().getTime();

    fetch(url, {
      cache: "no-cache",
    }).then(async (resp) => {
      const latency = new Date().getTime() - start;

      resolve(latency);
    });
  });
}

/**
 * Displays the results in one ping test
 * @param {string} region
 * @param {string} ping
 */
function displayPingResults(region, ping) {
  chrome.notifications.create("gcping-notif", {
    type: "basic",
    title: "Gcping",
    iconUrl: "../images/icon.png",
    message: `Gcping run complete. ${region} is the fastest region for you with a median ping of ${ping}ms.`,
  });
}

/**
 * Helper function that fetches the current regions stored in chrome.storage
 */
async function getRegionsToPing() {
  return new Promise((resolve, reject) => {
    chrome.storage.local.get([CHROME_STORAGE_ENDPOINTS_KEY], function (result) {
      resolve(result[CHROME_STORAGE_ENDPOINTS_KEY]);
    });
  });
}

/**
 * Saves the run data to chrome local storage
 * @param {object} runData
 */
async function saveRunData(runData) {
  const currentRuns = await getCurrentRuns();
  currentRuns.push(runData["startTime"]);

  const localData = {
    runs: currentRuns,
  };

  localData[`run-${runData["startTime"]}`] = runData;

  return new Promise((resolve, reject) => {
    chrome.storage.local.set(localData, resolve);
  });
}

/**
 * Fetches the past runs from chrome.storage
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
 * Helper to return the current status of the ping test
 * @return {Object}
 */
async function fetchCurrentStatus() {
  return currentStatus;
}

/**
 * Function that sends the current ping test status to the options page(s).
 */
function syncCurrentStatus() {
  chrome.runtime.sendMessage(
    { action: "sync_ping_status", currentStatus },
    function (response) {}
  );
}
