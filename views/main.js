import "./base.scss";
import "bootstrap";
import agents from "./agent.js";

const cardGeo = document.getElementById("card_geo");
const inputCode = document.getElementById("input_code");
const cardSubmit = document.getElementById("card_submit");
const btnSubmit = document.getElementById("btn_submit");

let location = {
  latitude: 0,
  longitude: 0,
}
let photoUri = null;
let gotLocation = false;

function success(position) {
  location.latitude  = position.coords.latitude;
  location.longitude = position.coords.longitude;
  gotLocation = true;

  cardGeo.classList.add("visually-hidden")
}

function error() {
  cardGeo.classList.remove("visually-hidden")
}

const whereAmI = () => {
  if (!navigator.geolocation) {
    error();
  } else {
    navigator.geolocation.getCurrentPosition(success, error);
  }
}

const initWebcam = () => {
  Webcam.reset();
  Webcam.set({
    width: 800,
    height: 450,
    image_format: "jpeg",
    jpeg_quality: 90,
  });
  Webcam.attach("#camera");
}

document.addEventListener("DOMContentLoaded", () => {
  whereAmI();
  initWebcam();
  check();
  setInterval(check, 3000);
});

const btnShoot = document.getElementById("btn_shoot");
const btnReset = document.getElementById("btn_reset");
btnShoot.addEventListener("click", () => {
  Webcam.snap(function (data_uri) {
    photoUri = data_uri;
    document.getElementById("camera").innerHTML =
      `<img src="${data_uri}" id="photo"/>`
  });
  check();
});

btnReset.addEventListener("click", () => {
  initWebcam();
  check();
})

function check() {
  whereAmI();
  let code = inputCode.value;

  if (!(code !== "" && photoUri !== null && gotLocation)) {
    cardSubmit.classList.add("visually-hidden");
    return;
  }

  cardSubmit.classList.remove("visually-hidden");
}

btnSubmit.addEventListener("click", () => {
  let volunteer_code = inputCode.value;
  let coordinate = location
  let photo = photoUri;


  agents.Logs.store({
    volunteer_code,
    coordinate,
    photo,
  }).then(({res, raw}) => {
    if (raw.status !== 201) {
      alert("Terjadi kesalahan, Coba Lagi ya! Error:" + raw.statusText);
      return;
    }

    afterSubmit();
  })
})

function afterSubmit() {
  alert("Terimakasih! Datamu sudah tersimpan, kamu boleh meninggalkan laman ini.");
}
