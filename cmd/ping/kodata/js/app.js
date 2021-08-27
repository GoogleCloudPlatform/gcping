// TODO: Add an option to contribute times and JS geolocation info to a public BigQuery dataset.

// These regions won't be plotted on the map
const IGNORE_REGIONS_PLOT={
  "global":true
},
MARKER_PATHS={
  "default":"/images/marker.svg",
  "slow":"/images/marker-red.svg",
  "medium":"/images/marker-orange.svg",
  "fast":"/images/marker-green.svg",
  "user":"/images/marker-user.svg"
};

const GLOBAL_REGION_KEY="global",
  PING_TEST_RUNNING_STATUS="running",
  PING_TEST_STOPPED_STATUS="stopped";

let map,
  zones = {},
  fastestZone = '',
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

  if (navigator.geolocation) {
    navigator.geolocation.getCurrentPosition(function (position) {
        const currentLocation = new google.maps.LatLng(position.coords.latitude, position.coords.longitude);

        map.setCenter(currentLocation);
        addCurrentUserToMap(currentLocation);
    });
}

  fetchZones();
}

function fetchZones() {
  fetch("/endpoints").then((resp) => {
    return resp.json();
  }).then(async (endpoints) => {
    for (zone of Object.values(endpoints)) {
      let gcpZone = { 
        region: zone.Region, 
        label: zone.RegionName, 
        pingUrl: zone.URL,
        lat: zone.Lat,
        lng: zone.Lng
      };
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
      document.getElementById("fastestRegion").innerText=`${zone.region} (${zone.label})(${zone.latency} ms)`;
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
    title=zones[region].latency === undefined ? zones[region].label : `${zones[region].label} ${zones[region].latency}ms`,
    lat=zones[region].lat,
    lng=zones[region].lng;

  if (lat !== undefined && lng!==undefined) {
    const marker = new google.maps.Marker({
      position: {lat,lng},
      map: map,
      title:title,
      icon:image
    });

    markers[region] = marker;
  }
}

function removeMarker(region){
  // if there is a marker present remove it
  if(markers[region]!==undefined){
    markers[region].setMap(null);
    delete markers[region];
  }
}

function getMarkerImage(region){
  const latency=zones[region].latency;

  if(latency===undefined){
    return MARKER_PATHS.default;
  }
  else if(latency<=100){
    return MARKER_PATHS.fast;
  }
  else if(latency>100 && latency<300){
    return MARKER_PATHS.medium;
  }
  else{
    return MARKER_PATHS.slow;
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

function addCurrentUserToMap(loc){
  // display the current user's location
  new google.maps.Marker({
    position: loc,
    icon: MARKER_PATHS.user,
    title: "This device",
    map: map,
  });
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