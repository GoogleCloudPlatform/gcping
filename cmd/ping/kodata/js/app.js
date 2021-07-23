// TODO: Show regions on a map, with lines overlayed according to ping times.
// TODO: Add an option to contribute times and JS geolocation info to a public BigQuery dataset.

// These regions won't be plotted on the map
const IGNORE_REGIONS_PLOT={
  "global":true
};

const GLOBAL_REGION_KEY="global",
  PING_TEST_RUNNING_STATUS="running",
  PING_TEST_STOPPED_STATUS="stopped";

let map,
  zones = {},
  fastestZone = '',
  locations = getLocations(),
  markers = {},
  sortKey = 'latency',
  sortOrder = "asc",
  pingTestStatus = PING_TEST_RUNNING_STATUS;

function initMap() {
  map = new google.maps.Map(document.getElementById("map"), {
    center: { lat: 17.2667283, lng: 30.0585942 },
    zoom: 3,
    gestureHandling: "cooperative",
    mapTypeControl: false,
    streetViewControl: false
  });

  fetchZones();
}

function fetchZones() {
  fetch("/endpoints").then((resp) => {
    return resp.json();
  }).then(async (endpoints) => {
    for (zone of Object.values(endpoints)) {
      let gcpZone = { region: zone.Region, label: zone.RegionName, pingUrl: zone.URL };
      zones[gcpZone.region] = gcpZone;
    }
    fetchPingData();
  });
}

function clearData(){
  // clear the markers
  Object.values(zones).forEach((zone)=>{
    // clear the markers
    removeMarker(zone.region);

    //clear the latency
    delete zone.latency;
  });

  // clear the data on the list view
  updateZoneList();
}

async function fetchPingData(){
  let analyzedRegions=0,
      zoneArr=Object.values(zones),
      totalRegions=zoneArr.length;

  for (zone of zoneArr) {
    if(zone.region!==GLOBAL_REGION_KEY)
      await updateRegionOnMap(zone.region);

    await fetchZoneLatency(zone.region);

    if(pingTestStatus==="stopped"){
      // remove the intermediate marker
      removeMarker(zone.region);
      break;
    }

    if(zone.region!==GLOBAL_REGION_KEY){
      updateRegionOnMap(zone.region);
      addRegionToList(zone.region);
    }
    

    if(zone.region===GLOBAL_REGION_KEY){
      document.getElementById("globalRegion").innerText=`${zone.latency} ms`;
    }

    if(zone.region!==GLOBAL_REGION_KEY && (fastestZone==='' || zones[fastestZone].latency>zone.latency)){
      fastestZone=zone.region;
      document.getElementById("fastestRegion").innerText=`${zone.region} (${zone.latency} ms)`;
    }

    analyzedRegions++;
    document.getElementById("analyzedRegions").innerText=analyzedRegions;
    document.getElementById("remainingRegions").innerText=totalRegions-analyzedRegions;

    if(pingTestStatus==="stopped"){
      break;
    }
  }

  // failsafe
  document.getElementById("remainingRegions").innerText=0;

  updatePingTestState(PING_TEST_STOPPED_STATUS);
}

function fetchZoneLatency(region) {
  return new Promise((resolve) => {
    const gcpZone = zones[region],
      start = new Date().getTime();

    fetch(gcpZone.pingUrl,{
      mode: 'no-cors',
      cache: 'no-cache'
    }).then((resp) => {
      const latency = new Date().getTime() - start;
      zones[region].latency = latency;

      resolve(latency);
    });
  });
}

async function updateRegionOnMap(region) {

  // if the region is in the ignore list we don't plot it
  if(IGNORE_REGIONS_PLOT[region]!==undefined)
    return;

  removeMarker(region);

  const image=getMarkerImage(region),
    title=zones[region].latency === undefined ? zones[region].label : `${zones[region].label} ${zones[region].latency}ms`;

  //we have the location cached
  if (locations[region] !== undefined) {
    const marker = new google.maps.Marker({
      position: locations[region],
      map: map,
      title:title,
      icon:image
    });

    markers[region] = marker;
  }
  //we fetch from the places API
  else {
    const location = await getLocationFromPlace(zones[region].label);
    if (location) {
      const marker = new google.maps.Marker({
        position: location,
        map: map,
        title:title,
        icon:image
      });

      markers[region] = marker;
    }
  }
}

function removeMarker(region){
  // if there is a marker present remove it
  if(markers[region]!==undefined){
    markers[region].setMap(null);
    delete markers[region];
  }
}

function getLocations() {
  return {
    "asia-east1": {lat: 23.4817418, lng: 118.9632941},
    "asia-east2": {lat: 22.3526632, lng: 113.9876185},
    "asia-northeast1": {lat: 35.5079447, lng: 139.2094288},
    "asia-northeast2": {lat: 34.66229, lng: 135.4807797},
    "asia-northeast3": {lat: 37.5638354, lng: 126.9040472},
    "asia-south1": { lat: 19.0822375, lng: 72.8111468 },
    "asia-southeast1": { lat: 1.3139991, lng: 103.7742106 },
    "asia-southeast2": { lat: -6.2297419, lng: 106.7594786 },
    "australia-southeast1": { lat: -33.8481647, lng: 150.7918939 },
    "europe-north1": { lat: 64.8255751, lng: 21.5433516 },
    "europe-west1": { lat: 50.499734, lng: 3.9057533 },
    "europe-west2": { lat: 51.5285582, lng: -0.2416781 },
    "europe-west3": { lat: 50.1211909, lng: 8.566525 },
    "europe-west4": { lat: 52.2076832, lng: 4.1585844 },
    "europe-west6": { lat: 47.3774497, lng: 8.5016958 },
    "northamerica-northeast1": { lat: 45.5016889, lng: -73.567256 },
    "southamerica-east1": { lat: -23.5505199, lng: -46.63330939999999 },
    "us-central1": { lat: 41.8780025, lng: -93.097702 },
    "us-east1": { lat: 33.836081, lng: -81.1637245 },
    "us-east4": { lat: 37.4315734, lng: -78.6568942 },
    "us-west1": { lat: 43.8041334, lng: -120.5542012 },
    "us-west2": { lat: 34.0522342, lng: -118.2436849 },
    "us-west3": { lat: 40.7607793, lng: -111.8910474 },
    "us-west4": { lat: 36.1699412, lng: -115.1398296 }
  };
}

function getLocationFromPlace(placeLabel) {
  return new Promise((resolve) => {
    var request = {
      query: placeLabel,
      fields: ['name', 'geometry'],
    };

    let service = new google.maps.places.PlacesService(map);

    service.findPlaceFromQuery(request, function (results, status) {
      if (status === google.maps.places.PlacesServiceStatus.OK) {
        resolve({ lat: results[0].geometry.location.lat(), lng: results[0].geometry.location.lng() });
      }

      resolve(false);
    });
  });
}

function getMarkerImage(region){
  const latency=zones[region].latency;

  if(latency===undefined){
    return "/images/marker.svg";
  }
  else if(latency<=100){
    return "/images/marker-green.svg";
  }
  else if(latency>100 && latency<300){
    return "/images/marker-orange.svg";
  }
  else{
    return "/images/marker-red.svg";
  }
}

function getZoneClass(region){
  const latency=zones[region].latency;

  if(latency===undefined){
    return "";
  }
  else if(latency<=100){
    return "fast";
  }
  else if(latency>100 && latency<300){
    return "average";
  }
  else{
    return "slow";
  }
}

function addRegionToList(region){
  updateZoneList();
}

function updateZoneList(){
  const parent=document.getElementById("listContainer"),
    list=getSortedListItems();
  // clear
  parent.querySelectorAll("li:not(.heading)").forEach((node)=>{
    parent.removeChild(node);
  });


  list.forEach(el => {
    if(el.region===GLOBAL_REGION_KEY)
      return;

    const cls=getZoneClass(el.region);
    parent.innerHTML=parent.innerHTML+`
    <li class="mdl-list__item ${cls}">
      <span class="mdl-list__item-primary-content list-zone-container">
        <span class="region-name">${el.region}</span>
        <span class="region-latency">${zones[el.region].latency ?? '-'}</span>
      </span>
    </li>
    `;
  });
  
}

// manual interupts to ongoing/stopped pings
function updatePingTestState(curState){
  pingTestStatus=curState;

  // UI changes based on the current status
  document.getElementById("runningCtrls").querySelector("button.visible").classList.remove("visible");
  if(pingTestStatus===PING_TEST_RUNNING_STATUS){
    document.getElementById("stopTest").classList.add("visible");
  }
  else if(pingTestStatus===PING_TEST_STOPPED_STATUS){
    document.getElementById("rerunTest").classList.add("visible");
  }
}

function getSortedListItems(){
  let curZones=Object.values(zones);
  curZones.sort((a,b)=>{
    return sortOrder==="asc" ? a[sortKey] - b[sortKey] : b[sortKey] - a[sortKey];
  });

  return curZones;
}

function registerDialog(){
  let dialog = document.querySelector('dialog'),
  showModalButton = document.querySelector('#how_it_works_link');

  if (! dialog.showModal) {
    dialogPolyfill.registerDialog(dialog);
  }
  showModalButton.addEventListener('click', function() {
    dialog.showModal();
  });
  dialog.querySelector('.close').addEventListener('click', function() {
    dialog.close();
  });
}

document.querySelector("body").addEventListener("click",function(e){
  if(e.target.classList.contains("toggle-sort-order")){
    sortOrder = (sortOrder === "asc" ? "desc" : "asc");
    e.target.setAttribute("data-order",sortOrder);
    updateZoneList();
  }
  
});

// stop the ongoing test
document.getElementById("stopTest").addEventListener("click",function(){
  updatePingTestState(PING_TEST_STOPPED_STATUS);
});

// rerun the test
document.getElementById("rerunTest").addEventListener("click",function(){
  updatePingTestState(PING_TEST_RUNNING_STATUS);

  // clear the currently fetched data
  clearData();

  // clear the Global region score
  document.getElementById("globalRegion").innerText=`- ms`;

  // reset the fastest region
  fastestZone='';
  document.getElementById("fastestRegion").innerText="-";

  // restart the pinging process
  fetchPingData();
});

registerDialog();