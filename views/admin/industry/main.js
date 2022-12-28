import "../../base.scss";
import "bootstrap";
import agents from "../../agent.js";

const newIndustryEl = document.getElementById("link_new");
const newFormCardEl = document.getElementById("form_new_card");
const addPlaceEl = document.getElementById("btn_add_place");
const placesEl = document.getElementById("places");
const formNewEl = document.getElementById("form_new");
const formSubmitEl = document.getElementById("btn_submit");

newIndustryEl.addEventListener("click", function () {
  newFormCardEl.classList.remove("visually-hidden")
})

addPlaceEl.addEventListener("click", function () {
  const stub = `<div class="row mb-2">
  <div class="col-auto latitude">
    <label class="form-label">Latitude</label>
    <input type="number" class="form-control form-control-sm">
  </div>
  <div class="col-auto longitude">
    <label class="form-label">Longitude</label>
    <input type="number" class="form-control form-control-sm">
  </div>
  <div class="col-auto wide">
    <label class="form-label">Wide</label>
    <input type="number" class="form-control form-control-sm">
  </div>
</div>
`;
  let place = document.createElement("div")
  place.innerHTML = stub;
  place = place.firstChild;

  placesEl.appendChild(place);
})

formNewEl.addEventListener("submit", function (ev) {
  ev.preventDefault();
})

formSubmitEl.addEventListener("click", function (ev) {
  let name = document.getElementById("input_name").value;
  let places = [];

  for (let i = 0; i < placesEl.children.length; i++) {
    const placeEl = placesEl.children.item(i);
    let latitude = placeEl.querySelector(".latitude").querySelector("input").value;
    let longitude = placeEl.querySelector(".longitude").querySelector("input").value;
    let wide = placeEl.querySelector(".wide").querySelector("input").value;

    latitude = parseInt(latitude)
    longitude = parseInt(longitude)
    wide = parseInt(wide)

    places.push({
      latitude,
      longitude,
      wide
    })
  }

  agents.Industries.store({
    name,
    places
  }).then(({res, raw}) => {
    if (res === null) {
      alert(raw.statusText)
      return;
    }

    console.log(res);

    if (raw.status === 201) {
      window.location.reload();
    }
  })
})

const industryListEl = document.getElementById("industry_list");

const getIndustryListItemEl = (id, industry) => {
  return `<li class="list-group-item">
  <a class="text-decoration-none text-dark d-block" data-bs-toggle="collapse" href="#${id}" role="button" aria-expanded="false" aria-controls="collapseExample">${industry.name}</a>
  <div class="collapse mt-2" id="${id}">
    <pre>${JSON.stringify(industry)}</pre>
  </div>
</li>`
}

document.addEventListener("DOMContentLoaded", function () {
  agents.Industries.getAll().then(({res, raw}) => {
    if (raw.status === 200 && res !== null) {
      res.forEach((ind, index) => {
        const tmp = document.createElement("div")
        tmp.innerHTML = getIndustryListItemEl(`ind-${index}`, ind)
        industryListEl.appendChild(tmp.firstChild)
      })
    }
  })
})
