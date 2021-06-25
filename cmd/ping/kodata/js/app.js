// TODO: Show regions on a map, with lines overlayed according to ping times.
// TODO: Add an option to contribute times and JS geolocation info to a public BigQuery dataset.

// These regions won't be plotted on the map
const IGNORE_REGIONS_PLOT={
  "global":true
};

const GLOBAL_REGION_KEY="global";

let map,
  zones = {},
  fastestZone = '',
  locations = getLocations(),
  markers = {};

function initMap() {
  map = new google.maps.Map(document.getElementById("map"), {
    center: { lat: 17.2667283, lng: 30.0585942 },
    zoom: 3,
    gestureHandling: "cooperative",
  });

  fetchZones();
}

function fetchZones() {
  fetch("/endpoints").then((resp) => {
    return resp.json();
  }).then(async (endpoints) => {
    let analyzedRegions=0,
      totalRegions=Object.values(endpoints).length;

    for (zone of Object.values(endpoints)) {
      let gcpZone = { region: zone.Region, label: zone.RegionName, pingUrl: zone.URL };
      zones[gcpZone.region] = gcpZone;

      if(gcpZone.region!==GLOBAL_REGION_KEY)
        await updateRegionOnMap(gcpZone.region);
      await fetchZoneLatency(gcpZone.region);
      if(gcpZone.region!==GLOBAL_REGION_KEY){
        updateRegionOnMap(gcpZone.region);
        addRegionToList(gcpZone.region);
      }
        

      if(gcpZone.region===GLOBAL_REGION_KEY){
        document.getElementById("globalRegion").innerText=`${zones[gcpZone.region].latency} ms`;
      }

      if(gcpZone.region!==GLOBAL_REGION_KEY && (fastestZone==='' || zones[fastestZone].latency>gcpZone.latency)){
        fastestZone=gcpZone.region;
        document.getElementById("fastestRegion").innerText=`${gcpZone.region} (${zones[gcpZone.region].latency} ms)`;
      }

      analyzedRegions++;
      document.getElementById("analyzedRegions").innerText=analyzedRegions;
      document.getElementById("remainingRegions").innerText=totalRegions-analyzedRegions;
    }

    // failsafe
    document.getElementById("remainingRegions").innerText=0;
  });
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

  // if there is a marker present remove it
  if(markers[region]!==undefined){
    markers[region].setMap=null;
    markers[region]=null;
  }

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
  const cls=getZoneClass(region);
  document.getElementById("listContainer").innerHTML=document.getElementById("listContainer").innerHTML+`
  <li class="mdl-list__item ${cls}">
    <span class="mdl-list__item-primary-content list-zone-container">
      <span class="region-name">${region}</span>
      <span class="region-latency">${zones[region].latency}</span>
    </span>
  </li>
  `;
}

document.getElementById("viewSwitch").addEventListener("change",function(e){
  // map view
  if(e.target.checked===true){
    document.getElementById("map").style.display="block";
    document.getElementById("list").style.display="none";
  }
  // list view
  else{
    document.getElementById("map").style.display="none";
    document.getElementById("list").style.display="block";
  }
})