const CHROME_ALARM_ID = 'gcping_endpoint_alarm';
const CHROME_STORAGE_ENDPOINTS_KEY = 'gcping_endpoints';

// when the extension is installed, add an alarm to refresh our endpoints
chrome.runtime.onInstalled.addListener(details => {
  if (details.reason === chrome.runtime.OnInstalledReason.INSTALL) {
    chrome.alarms.create(CHROME_ALARM_ID,{
      delayInMinutes: 0,
      periodInMinutes: 60
    });
  }
});

/**
 * Event listener for the alarm
 */
chrome.alarms.onAlarm.addListener(function(alarm){
  if(alarm.name === CHROME_ALARM_ID){
    fetchAndSaveEndpoints();
  }
});

/**
 * Event listener on click on the extension's action
 */
chrome.action.onClicked.addListener(async (tab) => {
  pingAllRegions();
});

/**
 * Helper to fetch the different endpoints that we need to ping
 * and save it in the chrome localstorage
 */
async function fetchAndSaveEndpoints() {
  return new Promise((resolve, reject)=>{
    fetch("https://gcping.com/api/endpoints")
    .then(function (resp) {
      return resp.json();
    })
    .then(function (endpoints) {
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
  let regions = await getRegionsToPing();

  // fallback in case the regions have never been fetched
  if(!regions){
    await fetchAndSaveEndpoints();
    regions = await getRegionsToPing();
  }

  let numRegions = Object.keys(regions).length,
    counter = 1,
    results = {},
    fastestRegion;

  chrome.action.setBadgeText({ text: `0/${numRegions}` });

  for (let region of Object.values(regions)) {
    let ping = await pingSingleRegion(region['pingUrl']);

    results[region['key']] = ping;
    if(fastestRegion === undefined || ping < results[fastestRegion]){
      fastestRegion = region['key'];
    }

    chrome.action.setBadgeText({ text: `${counter}/${numRegions}` });
    counter++;
  }

  chrome.action.setBadgeText({ text: '' });
  displayPingResults(fastestRegion, results[fastestRegion]);
}

/**
 * Helper function to ping a single URL and return the result
 * @param {string} url
 */
async function pingSingleRegion(url) {
  return new Promise((resolve) => {
    const start = new Date().getTime();

    fetch(url, {
      mode: "no-cors",
      cache: "no-cache",
    }).then(async (resp) => {
      const latency = new Date().getTime() - start;

      resolve(latency);
    });
  });
}

/**
 * Displays the results in one ping test
 */
function displayPingResults(region, ping) {
  chrome.notifications.create('gcping-notif', {
    type: 'basic',
    title: 'Gcping',
    iconUrl: '../images/icon.png',
    message: `Gcping run complete. ${region} is the fastest region for you with a median ping of ${ping}ms.`
  });
}

/**
 * Helper function that fetches the current regions stored in chrome.storage
 */
async function getRegionsToPing(){
  return new Promise((resolve, reject)=>{
    chrome.storage.local.get([CHROME_STORAGE_ENDPOINTS_KEY],function(result){
      resolve(result[CHROME_STORAGE_ENDPOINTS_KEY]);
    });
  });
}
