/**
 * Copyright 2021 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// TODO: Show regions on a map, with lines overlayed according to ping times.
// TODO: Add an option to contribute times and JS geolocation info to a public BigQuery dataset.

import { MDCDialog } from "@material/dialog";
import { MDCDataTable } from "@material/data-table";
import { MDCTooltip } from "@material/tooltip";

const GLOBAL_REGION_KEY = "global";
const PING_TEST_RUNNING_STATUS = "running";
const PING_TEST_STOPPED_STATUS = "stopped";
const INITIAL_ITERATIONS = 10;
const btnCtrl = document.getElementById("stopstart");

/**
 * The `regions` obj is of the following format:
 * {
 *  "us-east1": {
 *    "key": "",
 *    "label": "",
 *    "pingUrl": "",
 *    "latencies": [],
 *    "median": ""
 *  }
 * }
 */
const regions = {};
const results = []; // this will always be sorted according to sortKey and sortDir
let pingTestStatus = PING_TEST_RUNNING_STATUS;
let fastestRegionVisible = false;
let fastestRegion = null;
let globalRegionProxy = "";
let sortKey = "median"; // column to sort the data with
let sortDir = "ascending"; // sorting direction(ascending/descending)

/**
 * Fetches the endpoints for different Cloud Run regions.
 * We will later send a request to these endpoints and measure the latency.
 */
function getEndpoints() {
  fetch("/api/endpoints")
    .then(function (resp) {
      return resp.json();
    })
    .then(function (endpoints) {
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

      // once we're done fetching all endpoints, let's start pinging
      pingAllRegions(INITIAL_ITERATIONS);
    });
}

/**
 * Ping all regions to fetch their latency
 *
 * @param {number} iter
 */
async function pingAllRegions(iter) {
  const regionsArr = Object.values(regions);

  for (let i = 0; i < iter; i++) {
    for (const region of regionsArr) {
      // Takes care of the stopped button
      if (pingTestStatus === PING_TEST_STOPPED_STATUS) {
        break;
      }

      const latency = await pingSingleRegion(region.key);

      // add the latency to the array of latencies
      // from where we can compute the median and populate the table
      regions[region.key]["latencies"].push(latency);
      regions[region.key]["median"] = getMedian(
        regions[region.key]["latencies"],
      );

      // update fastest region
      if (
        fastestRegion === null ||
        regions[region.key]["median"] < regions[fastestRegion]["median"]
      ) {
        fastestRegion = region.key;
      }

      addResult(region.key);
      updateList();
    }

    // start displaying the fastest region after at least 1 iteration is over.
    // subsequent calls to this won't change anything
    displayFastest(true);
  }

  // when all the region latencies have been fetched, let's update our status flag
  updatePingTestState(PING_TEST_STOPPED_STATUS);
}

/**
 * Computes the ping time for a single GCP region
 * @param {string} regionKey The key of the GCP region, ex: us-east1
 * @return {Promise} Promise
 */
function pingSingleRegion(regionKey) {
  return new Promise((resolve) => {
    const gcpZone = regions[regionKey];
    const start = new Date().getTime();

    fetch(gcpZone.pingUrl, {
      mode: "no-cors",
      cache: "no-cache",
    }).then(async (resp) => {
      const latency = new Date().getTime() - start;

      // if we just pinged the global region, the response should contain
      // the region that the Global Load Balancer uses to route the traffic.
      if (regionKey === GLOBAL_REGION_KEY) {
        resp.text().then((val) => {
          globalRegionProxy = val.trim();
        });
      }

      resolve(latency);
    });
  });
}

/**
 * Function to update the current status of pinging
 * @param {string} status
 */
function updatePingTestState(status) {
  pingTestStatus = status;
  if (status === PING_TEST_RUNNING_STATUS) {
    btnCtrl.classList.add("running");
  } else if (status === PING_TEST_STOPPED_STATUS) {
    btnCtrl.classList.remove("running");
  }
}

/**
 * Updates the list view with the result set of regions and their latencies.
 */
function updateList() {
  let html = "";
  let cls = "";
  let regionKey = "";

  for (let i = 0; i < results.length; i++) {
    cls =
      results[i] === fastestRegion && fastestRegionVisible
        ? "fastest-region"
        : "";
    regionKey = getDisplayedRegionKey(results[i]);
    html +=
      '<tr class="mdc-data-table__row ' +
      cls +
      '"><td class="mdc-data-table__cell regiondesc">' +
      regions[results[i]]["label"] +
      '<div class="embedded-region d-none d-md-block">' +
      regionKey +
      "</div>" +
      '</td><td class="mdc-data-table__cell region d-md-none">' +
      regionKey +
      "</td>" +
      '<td class="mdc-data-table__cell result"><div>' +
      regions[results[i]]["median"] +
      " ms</div></td></tr>";
  }

  document.getElementsByTagName("tbody")[0].innerHTML = html;
}

/**
 * Helper function to return median from a given array
 * @param {*} arr Array of latencies
 * @return {*}
 */
function getMedian(arr) {
  if (arr.length == 0) {
    return 0;
  }
  const copy = arr.slice(0);
  copy.sort();
  return copy[Math.floor(copy.length / 2)];
}

/**
 * Helper that adds the regionKey to it's proper position keeping the results array sorted
 * This means we don't always have to sort the whole results array
 * TODO: Try and use an ordered map here to simply this
 * @param {string} regionKey
 */
function addResult(regionKey) {
  if (!results.length) {
    results.push(regionKey);
    return;
  }

  // remove any current values with the same regionKey
  for (let i = 0; i < results.length; i++) {
    if (results[i] === regionKey) {
      results.splice(i, 1);
      break;
    }
  }

  // TODO: Probably use Binary search here to merge the following 2 blocks
  // if new region is at 0th position
  if (compareTwoRegions(regionKey, results[0]) < 0) {
    results.unshift(regionKey);
    return;
  }
  // if new region is at last position
  else if (compareTwoRegions(regionKey, results[results.length - 1]) > 0) {
    results.push(regionKey);
    return;
  }

  // add the region to it's proper position
  for (let i = 0; i < results.length - 1; i++) {
    // if the region to be added is b/w i and i+1 elements
    if (
      compareTwoRegions(regionKey, results[i]) >= 0 &&
      compareTwoRegions(regionKey, results[i + 1]) < 0
    ) {
      results.splice(i + 1, 0, regionKey);
      return;
    }
  }
}

/**
 * Sets the visiblity for the fastest region indicator on the list(the green cell)
 * @param {bool} isVisible Indicator to toggle visibility for the fastest region indicator
 */
function displayFastest(isVisible) {
  fastestRegionVisible = true;
  updateList();
}

/**
 * Helper function to deduce the region to be displayed in the list
 * @param {string} regionKey
 * @return {string}
 */
function getDisplayedRegionKey(regionKey) {
  // if the region is not global, return it as it is.
  if (regionKey !== GLOBAL_REGION_KEY) return regionKey;

  // if the region is global and we have received the region that is used by the Gloabl Load Balancer
  // we display that
  if (globalRegionProxy.length > 0)
    return "<em>→" + globalRegionProxy + "</em>";

  // if the region is global and we don't have the routing region, we show "gloabl"
  return "global";
}

/**
 * Sort the table data based on a column(defined in sortKey) and direction(sortDir)
 */
function sortResults() {
  results.sort(compareTwoRegions);
}

/**
 * Function to compare order of 2 regions based on the current sort options
 * @param {string} a Region key for first region to be compared
 * @param {string} b Region key for second region to be compared
 * @return {int}
 */
function compareTwoRegions(a, b) {
  const multiplier = sortDir === "ascending" ? 1 : -1;

  a = regions[a][sortKey];
  b = regions[b][sortKey];

  if (a == b) {
    return 0;
  }

  return multiplier * (a > b ? 1 : -1);
}

/**
 * Event listener for the button to start/stop the pinging
 */
btnCtrl.addEventListener("click", function () {
  const newStatus =
    pingTestStatus === PING_TEST_STOPPED_STATUS
      ? PING_TEST_RUNNING_STATUS
      : PING_TEST_STOPPED_STATUS;
  updatePingTestState(newStatus);

  if (newStatus === PING_TEST_RUNNING_STATUS) pingAllRegions(1);
});

// start the process by fetching the endpoints
getEndpoints();

window.onload = function () {
  // How it works btn
  const dialog = new MDCDialog(document.querySelector(".mdc-dialog"));
  document
    .querySelector(".how-it-works-link")
    .addEventListener("click", function (e) {
      e.preventDefault();
      dialog.open();
    });

  // init data-table
  new MDCDataTable(document.querySelector(".mdc-data-table"));

  document
    .querySelector(".mdc-data-table")
    .addEventListener("MDCDataTable:sorted", function (data) {
      const detail = data.detail;

      // update the sorting options according to the requested values
      (sortKey = detail.columnId), (sortDir = detail.sortValue);

      sortResults();
      updateList();
    });

  // init tooltips
  [].map.call(document.querySelectorAll(".mdc-tooltip"), function (el) {
    return new MDCTooltip(el);
  });
};
