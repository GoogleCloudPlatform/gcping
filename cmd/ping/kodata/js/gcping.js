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

const GLOBAL_REGION_KEY="global",
  PING_TEST_RUNNING_STATUS="running",
  PING_TEST_STOPPED_STATUS="stopped",
  btnCtrl = document.getElementById('stopstart');

/**
 * The `regions` obj is of the following format:
 * {
 *  "us-east1": {
 *    "key": "",
 *    "label": "",
 *    "pingUrl": "",
 *    "latencies": []
 *  }
 * }
 */
let regions = {},
  results = [],
  pingTestStatus = PING_TEST_RUNNING_STATUS;

/**
 * Fetches the endpoints for different Cloud Run regions.
 * We will later send a request to these endpoints and measure the latency.
 */
function getEndpoints(){
  fetch("/endpoints").then(function(resp) { 
    return resp.json();
  }).then(function(endpoints) {
    for (zone of Object.values(endpoints)) {
      let gcpZone = {
        key: zone.Region, 
        label: zone.RegionName,
        pingUrl: zone.URL,
        latencies: []
      };

      regions[gcpZone.key] = gcpZone;
      results[gcpZone.key] = {'median':''};
    }

    // once we're done fetching all endpoints, let's start pinging
    pingAllRegions();
  });
}

/**
 * Ping all regions to fetch their latency
 */
async function pingAllRegions(){
  let regionsArr=Object.values(regions);

  // reset the results
  results=[];

  for (region of regionsArr) {
    let latency = await pingSingleRegion(region.key);

    // add the latency to the array of latencies
    // from where we can compute the median and populate the table
    regions[region.key]['latencies'].push(latency);
    results.push({"key": region.key, "median": getMedian(regions[region.key]['latencies'])});

    sortResults();
    updateList();
    updateTweetLink();

    // Takes care of the stopped button
    if(pingTestStatus === PING_TEST_STOPPED_STATUS){
      break;
    }
  }

  // when all the region latencies have been fetched, let's update our status flag
  updatePingTestState(PING_TEST_STOPPED_STATUS);
}

/**
 * Computes the ping time for a single GCP region
 * @param {string} regionKey The key of the GCP region, ex: us-east1
 * @returns Promise
 */
function pingSingleRegion(regionKey){
  return new Promise((resolve) => {
    const gcpZone = regions[regionKey],
      start = new Date().getTime();

    fetch(gcpZone.pingUrl,{
      mode: 'no-cors',
      cache: 'no-cache'
    }).then((resp) => {
      const latency = new Date().getTime() - start;

      resolve(latency);
    });
  });
}

/**
 * Function to update the current status of pinging
 * @param {string} status 
 */
function updatePingTestState(status){
  pingTestStatus = status;
  if(status === PING_TEST_RUNNING_STATUS){
    btnCtrl.querySelector('.material-icons').innerText = 'stop'
  }
  else if( status === PING_TEST_STOPPED_STATUS){
    btnCtrl.querySelector('.material-icons').innerText = 'play_arrow';
  }
}

/**
 * Updates the list view with the result set of regions and their latencies.
 */
function updateList(){
  let html = '',
    cls ='';

  for (let i = 0; i < results.length; i++) {
    cls = i ===0 ? 'top' : '';
    html += '<tr class="'+cls+'"><td class="regiondesc">'+regions[results[i]['key']]['label']+'<div class="region">'+results[i]['key']+'</div></td>' +
      '<td class="result" id="'+results[i]['key']+'"><div>'+results[i]['median']+' ms</div></td></tr>';
  }

  document.getElementsByTagName('tbody')[0].innerHTML = html;
}

/**
 * Helper function to return median from a given array
 * @param {*} arr Array of latencies
 * @returns 
 */
function getMedian(arr) {
  if (arr.length == 0) { return 0; }
  let copy = arr.slice(0);
  copy.sort();
  return copy[Math.floor(copy.length/2)];
}

/**
 * Simple sorting helper for the current result set
 */
function sortResults(){
  results = results.sort((a, b) =>{
    return a['median'] - b['median'];
  });
}

/**
 * Updates the tweet link to contain `numRegions` num of fastest regions.
 * @param {int} numRegions 
 */
function updateTweetLink(numRegions = 3){
  let tweet = 'My lowest-latency #GCP regions via gcping.com:';

  for(let i = 0; i < results.length; i++){
    if(results[i]['key'] !== 'global'){
      tweet += '\n'+results[i]['key']+' ('+results[i]['median']+' ms)';

      if(--numRegions === 0)
        break;
    }
  }

  document.getElementById('tweet-link').href = 'https://twitter.com/share?text='+encodeURIComponent(tweet);
}

getEndpoints();

btnCtrl.addEventListener('click',function(){
  let newStatus = pingTestStatus === PING_TEST_STOPPED_STATUS ? PING_TEST_RUNNING_STATUS : PING_TEST_STOPPED_STATUS;
  updatePingTestState(newStatus);

  if(newStatus === PING_TEST_RUNNING_STATUS)
    pingAllRegions();
});