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

const _MAX_RESULTS = 20; // Only consider most recent results.
let _ENDPOINTS = {};
let _REGIONS = [];
let IDX = 0;
let RESULTS = {};

function median(arr) {
  if (arr.length == 0) { return 0; }
  let copy = arr.slice(0);
  copy.sort();
  return copy[Math.floor(copy.length/2)];
}

let GLOBAL_GOT = '';
document.addEventListener('nextping', function() {
  let r = _REGIONS[IDX];
  IDX = (IDX+1)%_REGIONS.length; // wrap around
  let url = _ENDPOINTS[r].URL + '/ping';

  let start = new Date().getTime();
  fetch(url).then((value) => {
    if (!value.ok) {
      console.log('fetch', value.url, value.ok, value.status);
    }
    value.text().then((resp) => {
      if (r == 'global') {
        console.log('global got', resp);
        GLOBAL_GOT = resp.trim();
      }
    });
    if (value.headers['X-First-Request'] == 'true') {
      console.log('Discarding first request for', url);
    } else {
      let took = new Date().getTime()-start;
      RESULTS[r].push(took);
      if (RESULTS[r].length > _MAX_RESULTS) {
        RESULTS[r].shift();
      }
      let a = median(RESULTS[r]);
      let out = document.getElementById(r+'-result');
      updateTable();
    }
    if (!stopped) {
      document.dispatchEvent(new Event('nextping'));
    }
  });
});

function updateTable() {
  let medians = [];
  for (k in RESULTS) {
    medians.push([k, median(RESULTS[k])]);
  }
  medians.sort(function(a, b) {
    if (a[1] < b[1]) { return -1; }
    if (a[1] > b[1]) { return 1; }
    return 0;
  });

  let html = '';
  let tweet = 'My lowest-latency #GCP regions via gcping.com:';
  let place = 0;
  let top = 0;
  for (var i = 0; i < medians.length; i++) {
    let region = medians[i][0];
    let latency = medians[i][1];
    if (latency == 0) { continue; }
    if (i == 0 && region === 'global') { top++; }
    if (region != 'global') { place++; }
    let regionsub = region;
    if (region == 'global' && RESULTS[GLOBAL_GOT]) { regionsub = '<i>→'+GLOBAL_GOT+'</i>'; }
    var cls = (i == top) ? 'top' : '';
    var desc = (_ENDPOINTS[region] || {}).RegionName || '';
    html += '<tr class="'+cls+'"><td class="regiondesc">'+desc+'<div class="region">'+regionsub+'</div></td>' +
      '<td class="result" id="'+region+'"><div>'+latency+' ms</div></td></tr>';

    if (place <= 3 && region != 'global') {
      tweet += '\n'+region+' ('+latency+' ms)';
    }
  }
  document.getElementsByTagName('tbody')[0].innerHTML = html;
  document.getElementById('tweet-link').href = 'https://twitter.com/share?text='+encodeURIComponent(tweet);
}

let stopped = true;
let ss = document.getElementById('stopstart');
let st;
ss.onclick = function() {
  if (stopped) {
    console.log('starting');
    stopped = false;
    ss.children[0].innerText = 'stop';
    document.dispatchEvent(new Event('nextping'));
    st = setTimeout(ss.onclick, 30000); // stop after 30s.
  } else {
    console.log('stopping');
    stopped = true;
    clearTimeout(st); // cancel the stop timeout, as it is stopped.
    ss.children[0].innerText = 'play_arrow';
  }
};

// Start pinging.
fetch("/endpoints").then(function(resp) { 
  return resp.json();
}).then(function(endpoints) { 
  _ENDPOINTS = endpoints;
  Object.entries(endpoints).forEach(function (entry) {
    let key = entry[0];
    let value = entry[1];
    _REGIONS.push(value.Region);
    RESULTS[key] = [];
  })
  ss.onclick(); 
});